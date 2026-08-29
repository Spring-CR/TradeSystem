package handlers

import (
	"rhino-api/api/api_const"
	"rhino-api/api/options"
	"rhino-common/context"
	"rhino-common/server/middleware"
	"rhino-common/utils/request"
	trade_channel "rhino-instr/domain/trade-channel"
	domain_trade_instr "rhino-instr/domain/trade-instr"

	"github.com/gin-gonic/gin"
)

// 执行单笔交易指令
func ExecuteSingleTradeInstr(c *gin.Context) {
	opt := &options.ExecuteSingleTradeInstrOption{}
	if !middleware.BindInputOption(c, opt) {
		return
	}
	tradeInstr, de := domain_trade_instr.ExecuteSingleTradeInstr(opt.Date, opt.DailyInstrNo, opt.IndexDailyModify, opt.StockSerialNo, opt.OrdType, opt.Price, opt.Amount, opt.InstrOperator, context.DefaultATPUser, trade_channel.DefaultTradeChannel)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	middleware.ResponseJson(c, tradeInstr)
}

// 查询交易台母单编号
func FindTradeDeskOrderIdsByParentKey(c *gin.Context) {
	parentKey, de := request.GetQueryAsString(c, api_const.ParamKey, true)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	ids, de := domain_trade_instr.FindTradeDeskOrderIdsByParentKey(parentKey)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	if ids == nil {
		ids = make([]string, 0)
	}
	middleware.ResponseJson(c, ids)
}