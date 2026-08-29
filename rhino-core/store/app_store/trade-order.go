package app_store

import (
	"database/sql"
	"log"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"sort"
	"strings"
	"time"

	"github.com/linchunquan/sqlgen/db"
)

const (
	UpdateTradeOrderStatusByIdStmt         = `UPDATE trade_orders SET f_ord_status=?, f_ord_status_update_time=? WHERE f_id=?`
	UpdateTradeOrderOnSubmitFailedByIdStmt = `UPDATE trade_orders SET f_ord_status=?, f_ord_status_update_time=?, f_order_submit_fail_reason=? WHERE f_id=?`
	UpdateTradeOrderClientIDByIdStmt       = `UPDATE trade_orders SET f_cl_ord_id=? WHERE f_id=?`

	UpdateTradeOrderFromTradeActionRespByClOrdIDStmt = `UPDATE trade_orders SET 
    f_ord_id=?
   ,f_ord_status_update_time=?
   ,f_ord_status=?
   ,f_last_shares=?
   ,f_last_px=?
   ,f_leaves_qty=?
   ,f_cum_qty=?
   ,f_avg_px=?
   ,f_msg_seq=?
   WHERE f_cl_ord_id=? and f_msg_seq<?`

	UpdateTradeOrderFromTradeActionRespByOrdIDStmt = `UPDATE trade_orders SET 
    f_ord_status_update_time=?
   ,f_ord_status=?
   ,f_last_shares=?
   ,f_last_px=?
   ,f_leaves_qty=?
   ,f_cum_qty=?
   ,f_avg_px=?
   ,f_msg_seq=?
   WHERE f_order_id=? and f_channel_code=? and f_msg_seq<?`

	GetMaxMsgSeqOfTradeOrderStmt = `SELECT COALESCE(MAX(f_msg_seq),-1) FROM trade_orders`

	FindStreamTradeOrderBySystemCodeAndBusinessCodeStmt = SelectTradeOrderStmt + ` WHERE f_msg_seq>0 and f_system_code=? and f_business_code=?`

	UpdateTradeOrderCancelCountByAppOrdIdStmt = `UPDATE trade_orders SET f_ord_cancel_count=? WHERE f_app_ord_id=?`

	GetMsgSeqsOfTradeOrderStmt = `SELECT f_msg_seq FROM trade_orders where f_msg_seq >= ?`

	UpdateTradeOrderExtendAttrStmt = `UPDATE trade_orders SET f_extend_attr=? WHERE f_app_ord_id=?`
)

// 在mysql，当更新行的内容和原始内容是一样的，mysql是不执行实际更新操作的，这时候影响行数是0
// 在golden db for mysql，影响行数也是对的
func UpdateTradeOrderStatusById(db db.SimpleDB, id int64, ordStatus enum.OrdStatus) (affectRows int64, err error) {
	args := []interface{}{string(ordStatus), timeutil.ConvertTimeToMilliseconds(time.Now()), id}
	var result sql.Result
	result, err = db.Exec(UpdateTradeOrderStatusByIdStmt, args...)
	if err != nil {
		return
	}
	affectRows, err = result.RowsAffected()
	return
}

func UpdateTradeOrderOnSubmitFailedById(db db.SimpleDB, id int64, submitFailReason string) (affectRows int64, ordStatus string, ordStatusUpdateTime int64, err error) {
	ordStatus = string(enum.OrdStatus_InternalSubmitFailed)
	ordStatusUpdateTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	args := []interface{}{ordStatus, ordStatusUpdateTime, submitFailReason, id}
	log.Printf("UpdateTradeOrderOnSubmitFailedById, sql:%s, args:%+v\n", UpdateTradeOrderOnSubmitFailedByIdStmt, args)
	var result sql.Result
	result, err = db.Exec(UpdateTradeOrderOnSubmitFailedByIdStmt, args...)
	if err != nil {
		return
	}
	affectRows, err = result.RowsAffected()
	return
}

func UpdateTradeOrderClientIDById(db db.SimpleDB, id int64, clOrdID string) (err error) {
	args := []interface{}{clOrdID, id}
	_, err = db.Exec(UpdateTradeOrderClientIDByIdStmt, args...)
	return
}

func UpdateTradeOrderFromTradeActionRespByClOrdID(db db.SimpleDB, tradeActionResp *schema.TradeActionResp) (err error) {
	args := []interface{}{tradeActionResp.OrderID, tradeActionResp.TransactTime, tradeActionResp.OrdStatus, tradeActionResp.LastShares, tradeActionResp.LastPx, tradeActionResp.LeavesQty, tradeActionResp.CumQty, tradeActionResp.AvgPx, tradeActionResp.MsgSeq, tradeActionResp.ClOrdID, tradeActionResp.MsgSeq}
	_, err = db.Exec(UpdateTradeOrderFromTradeActionRespByClOrdIDStmt, args...)
	//log.Println(">>>finish UpdateTradeOrderFromTradeActionRespByClOrdID")
	return
}

func UpdateTradeOrderFromTradeActionRespByOrdID(db db.SimpleDB, tradeActionResp *schema.TradeActionResp) (err error) {
	args := []interface{}{tradeActionResp.TransactTime, tradeActionResp.OrdStatus, tradeActionResp.LastShares, tradeActionResp.LastPx, tradeActionResp.LeavesQty, tradeActionResp.CumQty, tradeActionResp.AvgPx, tradeActionResp.MsgSeq, tradeActionResp.OrderID, tradeActionResp.ChannelCode, tradeActionResp.MsgSeq}
	_, err = db.Exec(UpdateTradeOrderFromTradeActionRespByOrdIDStmt, args...)
	//log.Println(">>>finish UpdateTradeOrderFromTradeActionRespByOrdID")
	return
}

func GetMaxMsgSeqOfTradeOrder(db db.SimpleDB, systemCode, businessCode string) (msgSeq int, err error) {
	// row := db.QueryRow(GetMaxMsgSeqOfTradeOrderStmt)
	// err = row.Scan(&msgSeq)
	// if dbutil.IsDbRecordEmptyError(err) {
	// 	err = nil
	// 	msgSeq = -1
	// }

	//
	// 【背景】
	//
	// 因为用了并发连接池，在异常中断时，在一批同时执行的事务中，可能会有一些不能成功完成任务，从而导致序号不连续。因此，需要进一步调整逻辑。
	//
	// 1、根据systemcode、businesscode、channelcode找出全部 MsgSeq>0（表明是基于stream api流入的订单） 的tradeOrder。
	// 2、按 DBInsertTime 从小到大排序。
	// 3、找到 DBInsertTime 最大的记录。
	// 4、找到 DBInsertTime 仅次于或等于 DBInsertTime - 10s的记录
	// 5、如果第4步找不到任何记录，直接返回0
	// 6、从第4步中得到的记录中，找到 MsgSeq 最大的那一条返回
	//

	args := []interface{}{systemCode, businessCode}
	v, err := genericSelectTradeOrders(db, FindStreamTradeOrderBySystemCodeAndBusinessCodeStmt, args...)
	if dbutil.IsDbRecordEmptyError(err) || err == nil && len(v) == 0 {
		err = nil
		msgSeq = 0
		return
	}

	if err != nil {
		return
	}

	sort.Slice(v, func(i, j int) bool {
		return v[i].DBInsertTime < v[j].DBInsertTime
	})

	n := len(v)
	maxMsgTime := v[n-1].DBInsertTime
	targetMsgTime := maxMsgTime - 10*1000 // maxMsgTime的前10秒

	for i := n - 1; i >= 0; i-- {
		if v[i].DBInsertTime < targetMsgTime {
			return int(v[i].MsgSeq), nil
		}
	}

	return
}

// 与tx配合使用，锁表，确保获得正确的cancelCount和updateCount的值
func GetTradeOrderByAppOrdIdForUpdate(db db.SimpleDB, appOrdID string) (*schema.TradeOrder, error) {
	args := []interface{}{appOrdID}
	v, err := genericSelectTradeOrder(db, SelectTradeOrderByAppOrdIdStmt+" for update", args...)
	return v, err

}

// 更新订单撤销次数
func UpdateTradeOrderCancelCountByAppOrdId(db db.SimpleDB, appOrdID string, ordCancelCount int) error {
	args := []interface{}{ordCancelCount, appOrdID}
	_, err := db.Exec(UpdateTradeOrderCancelCountByAppOrdIdStmt, args...)
	return err
}

// 大于等于某个值的订单的全部msg序号
func GetMsgSeqsOfTradeOrder(db db.SimpleDB, minSeq int64) ([]int64, error) {
	args := []interface{}{minSeq}
	v, err := genericSelectNumbers(db, GetMsgSeqsOfTradeOrderStmt, args...)
	return v, err
}

func genericSelectNumbers(db db.SimpleDB, query string, args ...interface{}) (numbers []int64, err error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v0 sql.NullInt64
	for rows.Next() {
		err = rows.Scan(
			&v0,
		)
		if err != nil {
			return
		}

		if v0.Valid {
			numbers = append(numbers, v0.Int64)
		} else {
			numbers = append(numbers, 0)
		}
	}

	err = rows.Err()

	return
}

func FindAllTradeOrdersInRangeOrderById(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeOrder, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeOrders(db, strings.Replace(SelectTradeOrderRangeStmt, "FROM trade_orders", "FROM trade_orders ORDER BY f_id ASC", 1), args...)
	return v, err
}

func UpdateTradeOrderExtendAttr(db db.SimpleDB, tradeOrder*schema.TradeOrder) error {
	args := []interface{}{tradeOrder.ExtendAttr, tradeOrder.AppOrdID}
	_, err := db.Exec(UpdateTradeOrderExtendAttrStmt, args...)
	return err
}