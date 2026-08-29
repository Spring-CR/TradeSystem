package fix

import (
	"darkpool-common/bean"
	"fmt"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"rhino-core/store/app_store"
	"rhino-core/types"
	"rhino-trade-channel/adapter/data_convert/fix_convert"
	"strings"
	"time"

	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
)

// 撤销订单的算法步骤：
// 1、插入数据库记录TradeActionLatestResp。需要考虑cancel count的问题，需要用到select for update来锁表。
// 2、调用 afterTradeActionLatestRespInsert（更新内存模型）。
// 3、调用 FIX 接口发起撤单动作。
func (c *GenericFIXChannel) OrderCancelRequest(actionUser string, actionTime int64, actionKey string, targetClOrdID string, streamInputMsgSeq int64, order *schema.TradeOrder, afterTradeActionLatestRespInsert func(tradeActionLatestResp *schema.TradeActionLatestResp), syncTradeActionLatestRespAfterError func(tradeActionLatestResp *schema.TradeActionLatestResp)) (de *domain_error.Error) {

	// 获得正确的actionTime
	nowTime := timeutil.ConvertTimeToMilliseconds(time.Now())
	if actionTime <= 0 || actionTime > nowTime {
		actionTime = nowTime
	}

	// 先强制把autoTx全部提交，确保全部之前的TradeActionLatestResp记录已经入库
	c.cfg.GetAutoTx().Flush()

	// 开始事务，select for update必须与事务配合使用
	tx, de := dbutil.BeginTx(c.cfg.GetAppDB())
	if de != nil {
		return de
	}

	// 锁表取得TradeOrder记录，避免并发引起的脏读问题
	order2, err := app_store.GetTradeOrderByAppOrdIdForUpdate(tx, order.AppOrdID)
	if err != nil {
		dbutil.RollbackTx(tx)
		de = domain_error.Build(domain_error.CANNOT_GET_TRADE_ORDER_BY_APP_ORD_ID_ERR_CODE, err, order.AppOrdID)
		return
	}

	// 更新订单撤销次数
	cancelCount := order2.OrdCancelCount + 1
	err = app_store.UpdateTradeOrderCancelCountByAppOrdId(tx, order.AppOrdID, cancelCount)
	if err != nil {
		dbutil.RollbackTx(tx)
		de = domain_error.Build(domain_error.CANNOT_UPDATE_TRADE_ORDER_CANCEL_COUNT_ERR_CODE, err, order.AppOrdID, cancelCount)
		return
	}

	// 创建TradeActionLatestResp记录
	cancelAction := &schema.TradeActionLatestResp{
		ActionUser:             actionUser,
		ActionTime:             actionTime,
		ActionMsgTime:          nowTime,
		ActionType:             string(enum.ActionType_Withdraw),
		ActionKey:              actionKey,
		OrderID:                order.OrdID,
		AppOrdID:               order.AppOrdID,
		RootClOrdID:            order.ClOrdID,
		ClOrdID:                fmt.Sprintf("%s-cancel-%d", order.ClOrdID, cancelCount), // {原始订单ID}-cancel-{撤单次数}
		OrigClOrdID:            targetClOrdID,
		ChannelCode:            order.ChannelCode,
		SecurityExchange:       order.SecurityExchange,
		SecurityExchangeRegion: order.SecurityExchangeRegion,
		Account:                order.Account,
		Symbol:                 order.Symbol,
		Side:                   order.Side,
		OrdType:                order.OrdType,
		Price:                  order.Price,
		OrderQty:               order.OrderQty,
		CashOrderQty:           order.CashOrderQty,
		TransactTime:           actionTime,
		StreamInputMsgSeq:      streamInputMsgSeq,
	}
	if actionKey == "" {
		cancelAction.ActionKey = cancelAction.GetCacheKey()
	}

	// 插入TradeActionLatestResp记录
	err = app_store.InsertTradeActionLatestResp(tx, cancelAction)
	if err != nil {
		dbutil.RollbackTx(tx)

		de = domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
		if dbutil.IsMysqlDuplicateEntryError(err) {
			// 如果是违反唯一性约束，可以什么都不做
			de = nil
		}

		return
	}

	// 提交数据库事务
	de = dbutil.CommitTx(tx)
	if de != nil {
		dbutil.RollbackTx(tx)
		return
	}

	// 数据库插入完成，立即更新内存模型
	if afterTradeActionLatestRespInsert != nil {
		afterTradeActionLatestRespInsert(cancelAction)
	}

	// Todo：考虑在发送撤销指令前，系统出现重启的情况（同样的情形适用于下单操作）

	// 发送FIX撤单指令
	de = c.doOrderCancelRequestByFIX(cancelAction, order)
	if de != nil { // 如果发送FIX失败了，需要立即更新TradeActionLatestResp，设置取消被拒绝的状态和信息，拒绝理由可以直接取domain error的文本
		cancelAction.OrdStatus = string(enum.OrdStatus_Rejected)
		cancelAction.OrdRejReason = types.OrderCancelRejectPrefix + de.ErrorString()
		cancelAction.CxlRejResponseTo = string(enum.CxlRejResponseTo_Cancel)
		cancelAction.TransactTime = timeutil.ConvertTimeToMilliseconds(time.Now())
		cancelAction.MsgTime = timeutil.ConvertTimeToMilliseconds(time.Now())
		app_store.UpdateTradeActionLatestRespById(c.cfg.GetAppDB(), cancelAction)

		if syncTradeActionLatestRespAfterError != nil {
			syncTradeActionLatestRespAfterError(cancelAction)
		}
	}

	return
}

// 实现FIX撤单指令
func (c *GenericFIXChannel) doOrderCancelRequestByFIX(cancelAction *schema.TradeActionLatestResp, order*schema.TradeOrder) (de *domain_error.Error) {

	newCancelAction := &schema.TradeActionLatestResp{}
	bean.Copy(cancelAction).To(newCancelAction)

	if c.refineOrderCancelID != nil {
		newCancelAction.OrigClOrdID, newCancelAction.ClOrdID = c.refineOrderCancelID(cancelAction.OrigClOrdID, cancelAction.ClOrdID, order)
		cancelAction = newCancelAction
	}

	var msg *quickfix.Message

	switch enum.ChannelProtocolType(c.cfg.GetTradeChannel().ChannelProtocolType) {
	case enum.ChannelProtocolType_FIX42:
		cancelRequest := fix_convert.DomainTradeAction2OrderCancelRequest42(cancelAction)

		// 使用无后缀的symbol
		indexDot := strings.Index(cancelAction.Symbol, ".")
		if indexDot > 0 {
			cancelRequest.SetSymbol(cancelAction.Symbol[:indexDot])
		}
		// 设置Account
		if cancelAction.Account!=""{
			cancelRequest.SetAccount(cancelAction.Account)
		}
		
		// 设置OrderQty
		if cancelAction.OrderQty > 0 {
			cancelRequest.SetOrderQty(decimal.NewFromFloat(cancelAction.OrderQty), 0)
		}
		
		// 设置交易所
		if cancelAction.SecurityExchange != "" {
			cancelRequest.SetSecurityExchange(cancelAction.SecurityExchange)
		}

		msg = cancelRequest.ToMessage()

	case enum.ChannelProtocolType_FIX44:
		cancelRequest := fix_convert.DomainTradeAction2OrderCancelRequest44(cancelAction)

		// 使用无后缀的symbol
		indexDot := strings.Index(cancelAction.Symbol, ".")
		if indexDot > 0 {
			cancelRequest.SetSymbol(cancelAction.Symbol[:indexDot])
		}
		// 设置Account
		if cancelAction.Account!=""{
			cancelRequest.SetAccount(cancelAction.Account)
		}
		
		// 设置OrderQty
		if cancelAction.OrderQty > 0 {
			cancelRequest.SetOrderQty(decimal.NewFromFloat(cancelAction.OrderQty), 0)
		}
		
		// 设置交易所
		if cancelAction.SecurityExchange != "" {
			cancelRequest.SetSecurityExchange(cancelAction.SecurityExchange)
		}

		msg = cancelRequest.ToMessage()
	}

	c.QueryHeader(&msg.Header)

	err := quickfix.Send(msg)
	if err != nil {
		de = domain_error.Build(domain_error.FIX_SEND_MSG_ERR_CODE, err)
		return de
	}

	return
}
