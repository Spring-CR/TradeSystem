package open_order

import (
	"encoding/json"
	"errors"
	"ficc-utils/api/options"
	"ficc-utils/api/server/router"
	"ficc-utils/common/utils/orderutils"
	"ficc-utils/common/utils/wechat"
	"fmt"
	"log"
)

const (
	CheckOpenOrdersResultHead = "【盘后在途订单提示】\n"
)

func Task_CheckOpenOrders(tradeOrdersServiceUrl, WebhookUrl string) error {
	passed, checkInfo, err := checkOpenOrders(tradeOrdersServiceUrl)
	if err != nil {
		log.Printf("Run Task Task_CheckOpenOrders error: %v", err)
		return err
	}

	if !passed {
		msg := CheckOpenOrdersResultHead + checkInfo
		log.Printf("Run Task Task_CheckOpenOrders not passed:\n%s", msg)
		err = wechat.SendToWeChat(WebhookUrl, msg)
		if err != nil {
			log.Printf("Run Task Task_CheckOpenOrders SendToWeChat error: %v", err)
			return err
		}
		return nil
	}
	log.Println("Run Task Task_CheckOpenOrders check passed!")
	return nil
}

func checkOpenOrders(tradeOrdersServiceUrl string) (passed bool, checkInfo string, err error) {
	openOrders, err := qryOpenOrders(tradeOrdersServiceUrl)
	if err != nil {
		return false, "", err
	}
	if len(openOrders) > 0 {
		passed = false
		for _, order := range openOrders {
			checkInfo += fmt.Sprintf("订单号：%s，账户名称：%s，标的：%s %s，状态：%s\n", order.ClOrdID, order.Counterparty, order.Symbol, orderutils.ConvertSide(order.Side), orderutils.ConvertTrsOrderStatus(order.F_ord_status, order.F_approve_status, order.F_cum_qty))
		}
		return passed, checkInfo, nil
	}
	return true, "", nil
}

// 查询在途订单
func qryOpenOrders(tradeOrdersServiceUrl string) ([]*options.Order, error) {
	input := map[string]any{}
	fieldConditions := []map[string]any{}
	ordStatusCondition := map[string]any{
		"field": "f_ord_status",
		"field_type": 3, //字符串
		"value_type": 2, //集合
		"value": []string{"A", "0", "1", "6"}, // 0,A-NEW, 1-PARTIALLY_FILLED, 6-PENDING_CANCEL
	}
	fieldConditions = append(fieldConditions, ordStatusCondition)
	input["field_conditions"] = fieldConditions
	input["sort_fields"] = []string{"f_ord_create_time"}
	input["sort_type"] = 0
	input["limit"] = 1000000000
	input["offset"] = 0 // 0-升序，1-降序
	inputJson, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	result, err := router.ForwardQueryTradeOrdersService[*options.Order](tradeOrdersServiceUrl, inputJson)
	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, errors.New(result.Message)
	}

	return result.Data, nil
}