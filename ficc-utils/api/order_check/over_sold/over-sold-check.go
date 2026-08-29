package over_sold

import (
	"bytes"
	"encoding/json"
	"errors"
	"ficc-utils/common/utils/wechat"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	CheckOversoldResultHead = "【债券互换空头持仓提示】\n"
)

type OverSoldResponse struct {
	Total 	int 				`json:"total"`
	Data 	[]OverSoldPosition 	`json:"data"`
	Code    int          		`json:"Code"`
	Message string          	`json:"Message"`
}

type OverSoldPosition struct {
	Account           int       `json:"account"`
	Counterparty      string    `json:"counterparty"`
	Symbol            string    `json:"symbol"`
	NetPositionT1     float64   `json:"netPositionT1"`
	InitNetPositionT1 float64   `json:"initNetPositionT1"`
}

func Task_CheckOversoldPosition(positionServiceUrl, dataQryUrl, WebhookUrl string) error {
	passed, checkInfo, err := checkOversoldPosition(positionServiceUrl, dataQryUrl)
	if err != nil {
		log.Printf("Run Task Task_CheckOversoldPosition error: %v", err)
		return err
	}

	if !passed {
		msg := CheckOversoldResultHead + checkInfo
		log.Printf("Run Task Task_CheckOversoldPosition not passed:\n%s", msg)
		err = wechat.SendToWeChat(WebhookUrl, msg)
		if err != nil {
			log.Printf("Run Task Task_CheckOversoldPosition SendToWeChat error: %v", err)
			return err
		}
		return nil
	}
	log.Println("Run Task Task_CheckOversoldPosition check passed!")
	return nil
}

func checkOversoldPosition(positionServiceUrl, dataQryUrl string) (passed bool, checkInfo string, err error) {
	oversoldPositions, err := queryOversoldPosition(positionServiceUrl, dataQryUrl)
	if err != nil {
		return false, "", err
	}
	if len(oversoldPositions) > 0 {
		checkInfo = ""
		for _, position := range oversoldPositions {
			checkInfo += fmt.Sprintf("%s|%s, 日初: %s万元, 当前: %s万元\n", position.Counterparty, position.Symbol, strconv.FormatFloat(position.InitNetPositionT1/10000, 'f', -1, 64),  strconv.FormatFloat(position.NetPositionT1/10000, 'f', -1, 64))
		}
		return false, checkInfo, nil
	}
	return true, "债券互换空头持仓均为0\n", nil
}

func queryOversoldPosition(positionServiceUrl, dataQryUrl string) (result []OverSoldPosition, err error){
	//ctptyListT0, err := queryCtptyListT0(dataQryUrl)
	//if err != nil {
	//	log.Printf("queryCtptyListT0 error: %v\n", err)
	//	return nil, err
	//}
	//if len(ctptyListT0) == 0 {
	//	return nil, errors.New("CtptyListT0 is empty")
	//}

	input := make(map[string]any)
	input["select_fields"] = []string{"account", "counterparty", "symbol", "netPositionT1", "initNetPositionT1"}
	fieldConditions := []map[string]any{}
	netPositionT1Field := map[string]any{
		"field": "netPositionT1",
		"field_type": 1, //浮点
		"value_type": 1, //区间
		"value": []float64{-1000000000000, 0},
	}
	//accountField := map[string]any{
	//	"field": "account",
	//	"field_type": 2, //整形
	//	"value_type": 2, //集合
	//	"value": ctptyListT0,
	//}
	fieldConditions = append(fieldConditions, netPositionT1Field)
	//fieldConditions = append(fieldConditions, accountField)
	input["field_conditions"] = fieldConditions
	input["sort_fields"] = []string{"account", "symbol"}
	input["limit"] = 100000000
	input["offset"] = 0

	inputJson, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	log.Printf("invokeOverSoldPositionService positionServiceUrl(%s) input:%s", positionServiceUrl, inputJson)
	resp, err := http.Post(positionServiceUrl, "application/json", bytes.NewReader(inputJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("invokeOverSoldPositionService body:%s\n", body)
	output := OverSoldResponse{}
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}
	if output.Code != 0 {
		err = errors.New(output.Message)
		return nil, err
	}

	result = []OverSoldPosition{}
	for _, item := range output.Data {
		if item.NetPositionT1 < item.InitNetPositionT1 {
			result = append(result, item)
		}
	}
	log.Printf("invokeOverSoldPositionService result:%+v\n", result)
	return result, nil
}

func queryCtptyListT0(dataQryUrl string) (ctptyListT0 []int, err error) {
	qryUrl, err := url.Parse(dataQryUrl)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("collection", "ApplicationConfig")
	params.Set("key", "CtptyListT0")
	qryUrl.RawQuery = params.Encode()

	log.Printf("queryCtptyListT0 url:%s\n", qryUrl.String())
	resp, err := http.Get(qryUrl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("queryCtptyListT0 body:%s\n", body)
	output := struct {
		Code string `json:"code"`
		Msg string `json:"msg"`
		Stack string `json:"stack"`
		Total int32 `json:"total"`
		Data []struct {
			ConfigItemName string `json:"ConfigItemName"`
			ConfigItemValue string `json:"ConfigItemValue"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(body, &output)
	if output.Code != "" {
		return nil, errors.New(output.Msg)
	}
	if len(output.Data) < 1 {
		return nil, errors.New("result is empty")
	}
	ctptyListT0 = []int{}
	ids := strings.Split(output.Data[0].ConfigItemValue, ",")
	for _, id := range ids {
		ctptyId, err := strconv.Atoi(strings.TrimSpace(id))
		if err != nil {
			return nil, err
		}
		ctptyListT0 = append(ctptyListT0, ctptyId)
	}
	log.Printf("queryCtptyListT0 result:%+v\n", ctptyListT0)
	return ctptyListT0, nil

}