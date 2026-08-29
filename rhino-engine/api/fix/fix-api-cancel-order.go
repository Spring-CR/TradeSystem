package fix

import (
	"encoding/json"
	"fmt"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"rhino-core/types"
	"strconv"
	"time"

	fixenum "github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
)

func (a *FixApi) cancelOrder(message *quickfix.Message, sessionID quickfix.SessionID, rateLimitErr *domain_error.Error) quickfix.MessageRejectError {
	orderCancelRequest, rejErr := a.fixApiAdapter.DecodeForOrderCancelRequest(message, sessionID)
	if rejErr != nil {
		return rejErr
	}

	// 触发限流
	if rateLimitErr != nil {
		tradeOrder, ok := a.engine.GetOrderByAppOrdID(orderCancelRequest.AppOrdID)
		if !ok {
			return a.processDomainError(rateLimitErr, message)
		} else {

			senderCompID, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "senderCompID", enum.AttrValueType_STRING)
			targetCompID, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "targetCompID", enum.AttrValueType_STRING)
			if senderCompID == "" || targetCompID == "" {
				return a.processDomainError(rateLimitErr, message)
			}

			a.processCancelOrderError(orderCancelRequest, tradeOrder, rateLimitErr)
			return nil
		}
	}

	tradeOrder, de := a.engine.GetOrderAcceptor().AcceptOrderCancelRequest(orderCancelRequest)
	if rejErr := a.processDomainError(de, message); rejErr != nil {
		if tradeOrder == nil {
			return rejErr
		}
		a.processCancelOrderError(orderCancelRequest, tradeOrder, de)
		return nil
	}
	return nil
}

func (a *FixApi) processCancelOrderError(ordCxlReq *types.ApplicationOrderCancelRequest, tradeOrder *schema.TradeOrder, de *domain_error.Error) {

	msg := make(map[string]interface{})
	if len(tradeOrder.ExtendAttrMap) > 0 {
		msg = tradeOrder.ExtendAttrMap
	} else if len(tradeOrder.ExtendAttr) > 0 {
		err := json.Unmarshal([]byte(tradeOrder.ExtendAttr), &msg)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("ConvertTradeResponseMessage:: fail to parse tradeOrder.ExtendAttr:%s", tradeOrder.ExtendAttr), "fail to Unmarshal tradeOrder.ExtendAttr")
			return
		}
	}
	if len(msg) == 0 {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, fmt.Errorf("trade order ExtendAttrMap is empty"))
		domain_error.ProcessSevereError(false, 0, de, nil, "trade order ExtendAttrMap is empty")
		return
	}

	message := quickfix.NewMessage()
	senderCompID, _, _ := attrutil.GetAttrValue(msg, "senderCompID", enum.AttrValueType_STRING)
	if senderCompID == "" {
		return
	}
	targetCompID, _, _ := attrutil.GetAttrValue(msg, "targetCompID", enum.AttrValueType_STRING)
	if targetCompID == "" {
		return
	}
	message.Header.SetString(tag.MsgType, string(fixenum.MsgType_ORDER_CANCEL_REJECT))
	message.Header.SetString(tag.BeginString, "FIX.4.2")
	message.Header.SetString(tag.SenderCompID, senderCompID.(string))
	message.Header.SetString(tag.TargetCompID, targetCompID.(string))

	if tradeOrder.AppOrdID != "" {
		message.Body.SetString(tag.OrigClOrdID, tradeOrder.AppOrdID)
	}

	if tradeOrder.ClOrdID != "" {
		message.Body.SetString(tag.OrderID, tradeOrder.ClOrdID)
	} else {
		message.Body.SetString(tag.OrderID, "OrderID_for_"+tradeOrder.AppOrdID)
	}

	// 看看能否改成原来的状态
	message.Body.SetString(tag.OrdStatus, tradeOrder.OrdStatus)
	message.Body.SetString(tag.CxlRejResponseTo, "1")

	ordRejReason, _ := strconv.Atoi(de.Code)
	if ordRejReason > 0 {
		message.Body.SetString(tag.CxlRejReason, de.Code)
	}

	message.Body.SetString(tag.Text, de.Msg)

	if tradeOrder.Account != "" {
		message.Body.SetString(tag.Account, tradeOrder.Account)
	}

	message.Body.SetString(tag.TransactTime, time.Now().In(time.UTC).Format(timeutil.TransactTimeLayout))

	err := quickfix.Send(message)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to send FIX message")
	}
}
