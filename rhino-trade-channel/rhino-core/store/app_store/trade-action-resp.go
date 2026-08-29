package app_store

import (
	"rhino-common/utils/dbutil"
	"rhino-core/schema"
	"sort"
	"strings"

	"github.com/linchunquan/sqlgen/db"
)

var (
	UpdateTradeActionLatestRespFromTradeActionRespByClOrdIDStmt = `UPDATE trade_action_latest_resps SET 
    f_order_id=?
   ,f_orig_cl_ord_id=?
   ,f_exec_id=?
   ,f_exec_type=?
   ,f_ord_status=?
   ,f_ord_rej_reason=?
   ,f_exec_restatement_reason=?
   ,f_account=?
   ,f_symbol=?
   ,f_symbol_sfx=?
   ,f_security_id=?
   ,f_id_source=?
   ,f_security_type=?
   ,f_side=?
   ,f_open_close=?
   ,f_order_qty=?
   ,f_cash_order_qty=?
   ,f_ord_type=?
   ,f_price=?
   ,f_currency=?
   ,f_effective_time=?
   ,f_expire_time=?
   ,f_last_shares=?
   ,f_last_px=?
   ,f_leaves_qty=?
   ,f_cum_qty=?
   ,f_avg_px=?
   ,f_transact_time=?
   ,f_msg_seq=?
   WHERE f_cl_ord_id=? and f_msg_seq<?
   `

	UpdateTradeActionLatestRespFromTradeActionRespByOrdIDStmt = `UPDATE trade_action_latest_resps SET 
    f_orig_cl_ord_id=?
   ,f_exec_id=?
   ,f_exec_type=?
   ,f_ord_status=?
   ,f_ord_rej_reason=?
   ,f_exec_restatement_reason=?
   ,f_account=?
   ,f_symbol=?
   ,f_symbol_sfx=?
   ,f_security_id=?
   ,f_id_source=?
   ,f_security_type=?
   ,f_side=?
   ,f_open_close=?
   ,f_order_qty=?
   ,f_cash_order_qty=?
   ,f_ord_type=?
   ,f_price=?
   ,f_currency=?
   ,f_effective_time=?
   ,f_expire_time=?
   ,f_last_shares=?
   ,f_last_px=?
   ,f_leaves_qty=?
   ,f_cum_qty=?
   ,f_avg_px=?
   ,f_transact_time=?
   ,f_msg_seq=?
   WHERE f_order_id=? and f_channel_code=? and f_msg_seq<?
   `

	//GetMaxMsgSeqOfTradeActionRespStmt = `SELECT COALESCE(MAX(f_msg_seq),0) FROM trade_action_resps`
	GetMaxMsgSeqOfTradeActionRespStmt = `SELECT f_msg_seq FROM trade_action_resps WHERE f_id = (SELECT MAX(f_id) FROM trade_action_resps)`

	FindTradeActionRespBySystemCodeAndBusinessCodeAndChannelCodeStmt = func(systemCode, businessCode string) string {
		return SelectTradeActionRespStmt + ` WHERE f_cl_ord_id LIKE '%-` + systemCode + `-` + businessCode + `-%' AND f_channel_code=? AND f_msg_time > ?`
	}

	GetStreamInputMsgSeqsOfTradeActionLatestRespStmt = `SELECT f_stream_input_msg_seq FROM trade_action_latest_resps where f_stream_input_msg_seq >= ?`
)

func UpdateTradeActionLatestRespFromTradeActionRespByClOrdID(db db.SimpleDB, tradeActionResp *schema.TradeActionResp) (err error) {
	args := []interface{}{
		tradeActionResp.OrderID, tradeActionResp.OrigClOrdID, tradeActionResp.ExecID, tradeActionResp.ExecType, tradeActionResp.OrdStatus, tradeActionResp.OrdRejReason, tradeActionResp.ExecRestatementReason,
		tradeActionResp.Account, tradeActionResp.Symbol, tradeActionResp.SymbolSfx, tradeActionResp.SecurityID, tradeActionResp.IDSource, tradeActionResp.SecurityType, tradeActionResp.Side, tradeActionResp.OpenClose,
		tradeActionResp.OrderQty, tradeActionResp.CashOrderQty, tradeActionResp.OrdType, tradeActionResp.Price, tradeActionResp.Currency, tradeActionResp.EffectiveTime, tradeActionResp.ExpireTime, tradeActionResp.LastShares,
		tradeActionResp.LastPx, tradeActionResp.LeavesQty, tradeActionResp.CumQty, tradeActionResp.AvgPx, tradeActionResp.TransactTime, tradeActionResp.MsgSeq, tradeActionResp.ClOrdID, tradeActionResp.MsgSeq,
	}
	_, err = db.Exec(UpdateTradeActionLatestRespFromTradeActionRespByClOrdIDStmt, args...)
	//log.Println(">>>finish UpdateTradeActionLatestRespFromTradeActionRespByClOrdID")
	return
}

func UpdateTradeActionLatestRespFromTradeActionRespByOrdID(db db.SimpleDB, tradeActionResp *schema.TradeActionResp) (err error) {
	args := []interface{}{
		tradeActionResp.OrigClOrdID, tradeActionResp.ExecID, tradeActionResp.ExecType, tradeActionResp.OrdStatus, tradeActionResp.OrdRejReason, tradeActionResp.ExecRestatementReason,
		tradeActionResp.Account, tradeActionResp.Symbol, tradeActionResp.SymbolSfx, tradeActionResp.SecurityID, tradeActionResp.IDSource, tradeActionResp.SecurityType, tradeActionResp.Side, tradeActionResp.OpenClose,
		tradeActionResp.OrderQty, tradeActionResp.CashOrderQty, tradeActionResp.OrdType, tradeActionResp.Price, tradeActionResp.Currency, tradeActionResp.EffectiveTime, tradeActionResp.ExpireTime, tradeActionResp.LastShares,
		tradeActionResp.LastPx, tradeActionResp.LeavesQty, tradeActionResp.CumQty, tradeActionResp.AvgPx, tradeActionResp.TransactTime, tradeActionResp.MsgSeq, tradeActionResp.OrderID, tradeActionResp.ChannelCode, tradeActionResp.MsgSeq,
	}
	_, err = db.Exec(UpdateTradeActionLatestRespFromTradeActionRespByOrdIDStmt, args...)
	//log.Println(">>>finish UpdateTradeActionLatestRespFromTradeActionRespByOrdID")
	return
}

// 获取原始回报的最大消息序号
func GetMaxMsgSeqOfTradeActionResp(db db.SimpleDB, systemCode, businessCode, tradeChannelCode string, latestRestMsgSeqTime int64) (msgSeq int, err error) {
	//func GetMaxMsgSeqOfTradeActionResp(db db.SimpleDB) (msgSeq int, err error) {
	// row := db.QueryRow(GetMaxMsgSeqOfTradeActionRespStmt)
	// err = row.Scan(&msgSeq)
	// if dbutil.IsDbRecordEmptyError(err) {
	// 	err = nil
	// 	msgSeq = 0
	// }

	//
	// 【背景】
	//
	// 因为用了并发tx，在异常中断时，在一批同时执行的tx中，可能会有一些不能成功完成执行，这样就会导致序号不连续。因此，需要进一步调整逻辑。
	// 【算法1】（有漏洞，可能会丢失部分成交的回报）
	// 1、根据systemcode、businesscode、channelcode找出全部tradeResp
	// 2、按照ClOrdID的第4节，即order的数据库id，以及transacttime，从小到大双重排序（第一层，如果时间差大于1秒，时间小的在前面；如果小于1秒，按id确认先后）
	// 3、从右往左，找到第一个id是1的元素
	// 4、如果第3步找不到，则立即返回0
	// 5、从第3步找到的元素开始，如果下一个元素减去当前元素的id大于1，则立即返回当前元素的MsgSeq
	//
	// 【算法2】（不会丢任何消息）
	// 1、利用autoTx每秒自动flush一次的特性。
	// 2、根据systemcode、businesscode、channelcode找出全部tradeResp。
	// 3、按 DBInsertTime 从小到大排序。
	// 4、找到 DBInsertTime 最大的记录。
	// 5、找到 DBInsertTime 仅次于或等于 DBInsertTime - 10s的记录
	// 6、如果第5步找不到任何记录，直接返回0
	// 7、从第5步中得到的记录中，找到 MsgSeq 最大的那一条返回
	//

	args := []interface{}{tradeChannelCode, latestRestMsgSeqTime}
	// Todo：全部取数太长了，可以优化为按 DBInsertTime排序 分页获取，如果时间间隔没有达到10秒，继续渐进式取获取
	// Todo：注意：因为可能回出现session reset，原来tradeActionResps保留的seqnum就是之前session的序号，就没有意义了（解决reset问题的方案，增加一个标识为在发生reset时，dbInsertTime最大的那个resp的标识位设置为true，代表从它之后的tradeActionResp已经被reset过了）; 还要增加一个to_admin的心跳包的msgSeq计数，避免因seqnum差得太多而产生过多的ToAdmin: 8=FIX.4.29=10335=534=2749=TW52=20250222-03:02:40.00156=ISLD58=MsgSeqNum too low, expecting 4255 but received 2610=103
	// 提示：
	// 1. logon时同步设置重置序号：ToAdmin: 8=FIX.4.29=6635=A34=149=TW52=20250224-00:57:33.00156=ISLD98=0108=5141=Y10=202
	// 2. FIX协议参考：https://www.onixs.biz/fix-dictionary/4.2/msgType_A_65.html、https://www.onixs.biz/fix-dictionary/4.2/msgType_4_4.html
	v, err := genericSelectTradeActionResps(db, FindTradeActionRespBySystemCodeAndBusinessCodeAndChannelCodeStmt(systemCode, businessCode), args...)
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

	// 按 DBInsertTime 倒序遍历
	for i := n - 1; i >= 0; i-- {
		if v[i].DBInsertTime < targetMsgTime {
			return int(v[i].MsgSeq), nil
		}
	}

	return
}

// 大于等于某个值的tradeActionLatestResp的全部StreamInputMsgSeq
func GetStreamInputMsgSeqsOfTradeActionLatestResp(db db.SimpleDB, minSeq int64) ([]int64, error) {
	args := []interface{}{minSeq}
	v, err := genericSelectNumbers(db, GetStreamInputMsgSeqsOfTradeActionLatestRespStmt, args...)
	return v, err
}

func FindAllTradeActionLatestRespsInRangeOrderById(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeActionLatestResp, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeActionLatestResps(db, strings.Replace(SelectTradeActionLatestRespRangeStmt, "FROM trade_action_latest_resps", "FROM trade_action_latest_resps ORDER BY f_id ASC", 1), args...)
	return v, err
}

func FindAllTradeActionRespsInRangeOrderById(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeActionResp, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeActionResps(db, strings.Replace(SelectTradeActionRespRangeStmt, "FROM trade_action_resps", "FROM trade_action_resps ORDER BY f_id ASC", 1), args...)
	return v, err
}