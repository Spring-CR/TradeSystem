package order_status

import (
	"bytes"
	"darkpool-common/bean"
	"errors"
	"fmt"
	"log"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-core/schema"

	"github.com/linchunquan/sqlgen/db"
)

func (s *OrderStatusReplica) getOrderArgs(order *schema.TradeOrder, clOrdID string) (args []interface{}, argStatusOnly []interface{}, err error) {

	if clOrdID == "" {
		clOrdID = order.ClOrdID
	}

	extendAttrMap := order.ExtendAttrMap
	if extendAttrMap == nil {
		return nil, nil, errors.New("extendAttrMap is empty for order " + order.AppOrdID)
	}

	//jsData, _:=json.MarshalIndent(extendAttrMap, "", "  ")
	//log.Printf("======> extendAttrMap:%s\n", jsData)

	// 先加入id（仅在order insert时用到，order resp采用id自增的策略）
	args = append(args, order.ID)

	// 加入应用自定义字段
	n := len(s.extendAttrItems)
	for i := 0; i < n; i++ {
		extendAttrItem := s.extendAttrItems[i]
		if extendAttrItem.AttrName == "f_ord_status_update_time" {
			break
		}
		val, ok, err := attrutil.GetAttrValue(extendAttrMap, extendAttrItem.AttrName, enum.AttrValueType(extendAttrItem.AttrValueType))
		if err != nil {
			log.Printf("error occur while get value of %s, ok = %v\n", extendAttrItem.AttrName, ok)
			return nil, nil, err
		}
		args = append(args, val)
	}

	// 加入系统字段（对照initMemDb的字段次序，位置不能乱）
	// args 是用于插入新记录的参数列表
	args = append(args, order.OrdStatusUpdateTime)
	args = append(args, order.OrdStatus)
	args = append(args, order.DBInsertTime)
	args = append(args, order.TransactTime)
	args = append(args, order.Reviewer)
	args = append(args, order.ApproveStatus)
	args = append(args, order.LastShares)
	args = append(args, order.LastPx)
	args = append(args, order.LeavesQty)
	args = append(args, order.CumQty)
	args = append(args, order.AvgPx)
	args = append(args, order.OrdRejReason)
	// argStatusOnly 是用于更新订单状态的参数列表
	argStatusOnly = append(argStatusOnly, order.OrdStatusUpdateTime)
	argStatusOnly = append(argStatusOnly, order.OrdStatus)
	argStatusOnly = append(argStatusOnly, order.DBInsertTime)
	argStatusOnly = append(argStatusOnly, order.TransactTime)
	argStatusOnly = append(argStatusOnly, order.Reviewer)
	argStatusOnly = append(argStatusOnly, order.ApproveStatus)
	argStatusOnly = append(argStatusOnly, order.LastShares)
	argStatusOnly = append(argStatusOnly, order.LastPx)
	argStatusOnly = append(argStatusOnly, order.LeavesQty)
	argStatusOnly = append(argStatusOnly, order.CumQty)
	argStatusOnly = append(argStatusOnly, order.AvgPx)
	argStatusOnly = append(argStatusOnly, order.OrdRejReason)
	// args相比argStatusOnly还要多出AppOrdID和ClOrdID
	args = append(args, order.AppOrdID)
	args = append(args, clOrdID)
	args = append(args, order.OrdCreateTime)
	args = append(args, order.AlgParams)
	args = append(args, order.TradeDate)

	// 增加了四个用户
	// args = append(args, order.OrdCreator)
	// args = append(args, order.OrdDraftUpdateUser)
	// args = append(args, order.OrdDraftDelUser)
	// args = append(args, order.OrdExecUser)

	return
}

func (s *OrderStatusReplica) insertOrder(db db.SimpleDB, order *schema.TradeOrder) error {
	args, _, err := s.getOrderArgs(order, "")
	if err != nil {
		return err
	}
	_, err = db.Exec(s.insertOrderSql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderStatusReplica) updateOrder(db db.SimpleDB, order *schema.TradeOrder) error {

	_, argStatusOnly, err := s.getOrderArgs(order, "")
	if err != nil {
		return err
	}

	// 按AppOrdID来更新
	argStatusOnly = append(argStatusOnly, order.AppOrdID)

	_, err = db.Exec(s.updateOrderSql, argStatusOnly...)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderStatusReplica) fullUpdateOrder(db db.SimpleDB, order *schema.TradeOrder) error {

	args, _, err := s.getOrderArgs(order, "")
	if err != nil {
		return err
	}
	args = args[1:]
	// 按AppOrdID来更新
	args = append(args, order.AppOrdID)

	_, err = db.Exec(s.fullUpdateOrderSql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderStatusReplica) updateOrderAttributes(db db.SimpleDB, appOrdID string, updateAttrs map[string]interface{}) error {

	var args []interface{}
	sqlBuf := bytes.NewBufferString(fmt.Sprintf(`UPDATE trade_orders_%s_%s SET `, s.systemCode, s.businessCode))
	first := true
	for k := range updateAttrs {
		extendAttrItem, ok := s.extendAttrMap[k]
		if !ok {
			continue
		}
		val, ok, err := attrutil.GetAttrValue(updateAttrs, k, enum.AttrValueType(extendAttrItem.AttrValueType))
		if err != nil {
			log.Printf("error occur while get value of %s, ok = %v\n", extendAttrItem.AttrName, ok)
			return err
		}
		args = append(args, val)
		if !first {
			sqlBuf.WriteString(",")
		}
		sqlBuf.WriteString(extendAttrItem.AttrName + "=?\n")

		first = false
	}

	sqlBuf.WriteString("WHERE f_app_ord_id=?")
	args = append(args, appOrdID)

	_, err := db.Exec(sqlBuf.String(), args...)
	return err
}

func (s *OrderStatusReplica) insertOrderResp(db db.SimpleDB, order *schema.TradeOrder, tradeActionResp *schema.TradeActionResp) error {

	// 复制一个订单
	tradeOrder := &schema.TradeOrder{}
	err := bean.Copy(order).To(tradeOrder)
	if err != nil {
		return err
	}
	// 跟traceable-order.go的逻辑保持一致
	tradeOrder.OrdStatus = tradeActionResp.OrdStatus
	tradeOrder.OrdStatusUpdateTime = tradeActionResp.TransactTime
	tradeOrder.LastShares = tradeActionResp.LastShares
	tradeOrder.LastPx = tradeActionResp.LastPx
	tradeOrder.LeavesQty = tradeActionResp.LeavesQty
	tradeOrder.CumQty = tradeActionResp.CumQty
	tradeOrder.AvgPx = tradeActionResp.AvgPx
	tradeOrder.OrdRejReason = tradeActionResp.OrdRejReason

	args, _, err := s.getOrderArgs(tradeOrder, tradeActionResp.ClOrdID)
	if err != nil {
		return err
	}

	// 加入系统字段（对照createOrderResponseTable的字段次序，位置不能乱）
	args = append(args, tradeActionResp.OrigClOrdID)
	args = append(args, tradeActionResp.ExecID)
	args = append(args, tradeActionResp.ExecRefID)
	args = append(args, tradeActionResp.ExecTransType)
	args = append(args, tradeActionResp.ExecType)
	args = append(args, tradeActionResp.MsgTime)
	args = append(args, tradeActionResp.ChannelCode)

	// 加入成交回报的扩展字段
	if tradeActionResp.ExtendAttrMap == nil {
		tradeActionResp.RecoverExtendAttrMap()
	}
	for _, extendAttrItem := range s.tradeRespAttrItems {
		val, ok, err := attrutil.GetAttrValue(tradeActionResp.ExtendAttrMap, extendAttrItem.AttrName, enum.AttrValueType(extendAttrItem.AttrValueType))
		if err != nil {
			log.Printf("error occur while get value of %s, ok = %v\n", extendAttrItem.AttrName, ok)
		}
		args = append(args, val)
	}

	// 去掉第一个id字段，因为使用的是id自增策略
	args = args[1:]

	_, err = db.Exec(s.insertOrderRespSql, args...)
	if err != nil {
		return err
	}

	return nil
}
