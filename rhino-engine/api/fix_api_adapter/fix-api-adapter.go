package fix_api_adapter

import (
	"rhino-common/domain_error"
	"rhino-core/schema"
	"rhino-core/types"
	"time"

	"github.com/quickfixgo/quickfix"
)

type FixApiAdapter interface {
	GetConfigPath() string
	LoginValidate(username, password string, sessionID quickfix.SessionID) bool
	DecodeForNewOrderSingle(message *quickfix.Message, sessionID quickfix.SessionID) (msgProps map[string]interface{}, rejErr quickfix.MessageRejectError)
	AutoTurnToReviewForErrors(order *schema.TradeOrder, de *domain_error.Error) (turnToReview bool, displayErr *domain_error.Error)
	DecodeForOrderCancelRequest(message *quickfix.Message, sessionID quickfix.SessionID) (applicationOrderCancelRequest *types.ApplicationOrderCancelRequest, rejErr quickfix.MessageRejectError)
	ConvertTradeResponseMessage(tradeResp *types.TradeActionRespReturn) (message *quickfix.Message, de *domain_error.Error)
	GetFixPortOpenTimeRange() (begin string, end string, layout string, local *time.Location) // 格式：08:00-20:00
}
