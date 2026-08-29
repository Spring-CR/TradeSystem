package types

import (
	"rhino-core/schema"
	"sync"
)

// 应用层输入的OrderCancelRequest数据结构
type ApplicationOrderCancelRequest struct {
	ActionUser        string // 动作执行人。（参考TradeActionLatestResp）
	ActionTime        int64  // 动作发起时到毫秒时间戳。（参考TradeActionLatestResp）
	ActionKey         string // 确保幂等性
	AppOrdID          string // 由对接中央交易台的应用方设置。一方面，该指定允许应用方根据业务情况自定义对其更友好的订单ID；另一方面，需要保障AppOrdID的全局唯一性，系统在底层数据库对该字段添加唯一性约束，以防止重复下单。（参考TradeOrder）
	StreamInputMsgSeq int64  // 当使用Stream Api传入订单时，这个字段设置为传入参数的消息序号。
}

type ApplicationOrderDraftDeleteRequest struct {
	ActionUser string // 动作执行人。（参考TradeActionLatestResp）
	ActionTime int64  // 动作发起时到毫秒时间戳。（参考TradeActionLatestResp）
	AppOrdID   string
}

type ApplicationOrderAttributeUpdateRequest struct {
	ActionUser       string                 // 动作执行人。（参考TradeActionLatestResp）
	ActionTime       int64                  // 动作发起时到毫秒时间戳。（参考TradeActionLatestResp）
	AppOrdID         string                 // 订单ID
	UpdateAttributes map[string]interface{} // 待更新的属性
}

type TradeActionRespReturn struct {
	*TraceableTradeActionResp
	CurrentTradeActionResp *schema.TradeActionResp
	WaitGroup              *sync.WaitGroup
}

func NewTradeActionRespReturn(traceableTradeActionResp *TraceableTradeActionResp, currentTradeActionResp *schema.TradeActionResp) *TradeActionRespReturn {
	inst := &TradeActionRespReturn{}
	inst.TraceableTradeActionResp = traceableTradeActionResp
	inst.CurrentTradeActionResp = currentTradeActionResp
	return inst
}
