package fix

import (
	"encoding/json"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"strconv"
	"time"

	fixenum "github.com/quickfixgo/enum"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
)

func (a *FixApi) newOrderSingle(message *quickfix.Message, sessionID quickfix.SessionID, rateLimitErr *domain_error.Error) quickfix.MessageRejectError {

	msgProps, rejErr := a.fixApiAdapter.DecodeForNewOrderSingle(message, sessionID)
	if rejErr != nil {
		return rejErr
	}

	order, de := a.apiAdapter.ConvertNewOrderSingleMessage(nil, msgProps)
	if rejErr := a.processDomainError(de, message); rejErr != nil {
		return rejErr
	}

	if order.AppOrdID == "" {
		de := domain_error.Build(domain_error.APPORDID_EMPTY_ERR_CODE, nil)
		if rejErr := a.processDomainError(de, message); rejErr != nil {
			return rejErr
		}
	}

	if rateLimitErr != nil {
		log.Printf("ratelimit achived for order %s\n", order.AppOrdID)
		a.processNewOrderSingleError(order, rateLimitErr)
		return nil
	}

	order.SystemCode = a.systemCode
	order.BusinessCode = a.businessCode

	// 订单属性精化和校验
	de = a.apiAdapter.RefineAndValidate(order, true)
	if de != nil {
		turnToReview, displayErr := a.fixApiAdapter.AutoTurnToReviewForErrors(order, de)
		if turnToReview {

			// 转json文本
			if len(order.ExtendAttrMap) > 0 && len(order.ExtendAttr) == 0 {
				data, _ := json.Marshal(order.ExtendAttrMap)
				order.ExtendAttr = string(data)
			}

			if len(order.AlgParamsMap) > 0 && len(order.AlgParamsMap) == 0 {
				data, _ := json.Marshal(order.AlgParamsMap)
				order.AlgParams = string(data)
			}

			// 更新ExtendAttrMap
			_, de = a.apiAdapter.ConvertTradeOrderParams(order)
			if rejErr := a.processDomainError(de, message); rejErr != nil {
				//return rejErr
				a.processNewOrderSingleError(order, de)
				return nil
			}

			// 转为待审批
			de = a.engine.GetOrderAcceptor().AcceptOrderDraft(order, enum.ActionType_SubmitForReview)
			if rejErr := a.processDomainError(de, message); rejErr != nil {
				//return rejErr
				a.processNewOrderSingleError(order, de)
				return nil
			}

			// 拒绝订单并提示订单已经转入待审批
			if rejErr := a.processDomainError(displayErr, message); rejErr != nil {
				return rejErr
			}

			return nil
		} else {
			a.processNewOrderSingleError(order, de)
			return nil
		}
	}

	// 转json文本
	if len(order.ExtendAttrMap) > 0 && len(order.ExtendAttr) == 0 {
		data, _ := json.Marshal(order.ExtendAttrMap)
		order.ExtendAttr = string(data)
	}

	if len(order.AlgParamsMap) > 0 && len(order.AlgParamsMap) == 0 {
		data, _ := json.Marshal(order.AlgParamsMap)
		order.AlgParams = string(data)
	}

	// 更新ExtendAttrMap
	_, de = a.apiAdapter.ConvertTradeOrderParams(order)
	if rejErr := a.processDomainError(de, message); rejErr != nil {
		//return rejErr
		a.processNewOrderSingleError(order, de)
		return nil
	}

	duplicatedOrder, de := a.engine.GetOrderAcceptor().AcceptNewOrderSingleRequest(order)
	if de == nil && duplicatedOrder {
		de = domain_error.Build(domain_error.DUPLICATE_ORDER_ERR_CODE, nil, order.AppOrdID)
	}

	if de != nil {
		turnToReview, displayErr := a.fixApiAdapter.AutoTurnToReviewForErrors(order, de)
		if turnToReview {

			// 转为待审批
			de = a.engine.GetOrderAcceptor().AcceptOrderDraft(order, enum.ActionType_SubmitForReview)
			if rejErr := a.processDomainError(de, message); rejErr != nil {
				//return rejErr
				a.processNewOrderSingleError(order, de)
				return nil
			}

			// 拒绝订单并提示订单已经转入待审批
			if rejErr := a.processDomainError(displayErr, message); rejErr != nil {
				return rejErr
			}

			return nil
		} else {
			a.processNewOrderSingleError(order, de)
			return nil
		}
	}

	return nil
}

func (a *FixApi) processNewOrderSingleError(tradeOrder *schema.TradeOrder, de *domain_error.Error) {

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
	message.Header.SetString(tag.MsgType, string(fixenum.MsgType_EXECUTION_REPORT))
	message.Header.SetString(tag.BeginString, "FIX.4.2")
	message.Header.SetString(tag.SenderCompID, senderCompID.(string))
	message.Header.SetString(tag.TargetCompID, targetCompID.(string))

	if tradeOrder.OrdType != "" {
		message.Body.SetString(tag.OrdType, tradeOrder.OrdType)
	}

	if tradeOrder.SecurityExchange != "" {
		message.Body.SetString(tag.SecurityExchange, tradeOrder.SecurityExchange)
	}

	if tradeOrder.HandlInst != "" {
		message.Body.SetString(tag.HandlInst, tradeOrder.HandlInst)
	}

	if tradeOrder.AppOrdID != "" {
		message.Body.SetString(tag.ClOrdID, tradeOrder.AppOrdID)
	}

	if tradeOrder.ClOrdID != "" {
		message.Body.SetString(tag.OrderID, tradeOrder.ClOrdID)
	} else {
		message.Body.SetString(tag.OrderID, "OrderID_for_"+tradeOrder.AppOrdID)
	}

	message.Body.SetString(tag.ExecID, fmt.Sprintf("ExecID_for_%s_%d", tradeOrder.AppOrdID, timeutil.ConvertTimeToMilliseconds(time.Now())))
	message.Body.SetString(tag.ExecType, string(fixenum.ExecType_REJECTED))
	message.Body.SetString(tag.ExecTransType, "0")
	message.Body.SetString(tag.OrdStatus, string(fixenum.OrdStatus_REJECTED))

	ordRejReason, _ := strconv.Atoi(de.Code)
	if ordRejReason > 0 {
		message.Body.SetString(tag.OrdRejReason, de.Code)
	}

	message.Body.SetString(tag.Text, de.Msg)

	if tradeOrder.Account != "" {
		message.Body.SetString(tag.Account, tradeOrder.Account)
	}

	if tradeOrder.Symbol != "" {
		message.Body.SetString(tag.Symbol, tradeOrder.Symbol)
	}

	if tradeOrder.Side != "" {
		message.Body.SetString(tag.Side, tradeOrder.Side)
	}

	if tradeOrder.OpenClose != "" {
		message.Body.SetString(tag.OpenClose, tradeOrder.OpenClose)
	}

	if tradeOrder.OrderQty > 0 {
		message.Body.SetInt(tag.OrderQty, int(tradeOrder.OrderQty))
	}

	price := strconv.FormatFloat(tradeOrder.Price, 'f', -1, 64)
	if price != "" {
		message.Body.SetString(tag.Price, price)
	}

	message.Body.SetInt(tag.LeavesQty, 0)
	message.Body.SetInt(tag.CumQty, 0)
	message.Body.SetInt(tag.AvgPx, 0)

	message.Body.SetString(tag.TransactTime, time.Now().In(time.UTC).Format(timeutil.TransactTimeLayout))

	err := quickfix.Send(message)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to send FIX message")
	}
}
