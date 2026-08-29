package titans

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-core/schema"
	"time"
)

var (
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
)

// 定义与JSON结构对应的嵌套结构体
type response struct {
	Data *data `json:"data"` // 使用指针便于判空
}

type data struct {
	QueryResults []result `json:"queryResults"` // 数组可能为空
}

type result struct {
	AvailableQty float64 `json:"availableQty"` // 直接映射目标字段
}

func getAvailableQty(jsonText []byte) (availableQty float64, ok bool) {
	// 1. 解析JSON到结构体
	var resp response
	if err := json.Unmarshal(jsonText, &resp); err != nil {
		return 0, false // JSON语法错误
	}

	// 2. 逐层检查路径是否存在
	if resp.Data == nil {
		return 0, false // data或queryResults为空
	}

	if len(resp.Data.QueryResults) == 0 {
		return 0, true // data或queryResults为空
	}

	// 3. 提取首个元素的availableQty（根据需求调整索引）
	firstResult := resp.Data.QueryResults[0]
	return firstResult.AvailableQty, true
}

func (a *TitansCrossBorderAPIAdapter) checkPosition(order *schema.TradeOrder) {

	if order.MsgSeq > 0 || !a.softCheckPosition || order.Side != string(enum.Side_Sell) {
		return
	}

	val, ok, err := attrutil.GetAttrValue(order.ExtendAttrMap, "keyInstrumentId", enum.AttrValueType_INT)
	if !ok || err != nil {
		return
	}

	log.Printf("checkPosition, keyInstrumentId:%v\n", val)

	// 请求URL
	url := fmt.Sprintf("%s/api/titans/trading/1.0.0/hk/entrust/position?pageNum=1&pageSize=20", a.titansApiBase)

	// 请求体数据
	requestBody := map[string]interface{}{
		"keyUnderlyingInstIdList": []int{val.(int)},
	}

	// 转换为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("JSON序列化失败: %v\n", err)
		return
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建请求失败: %v\n", err)
		return
	}

	// 设置请求头
	req.Header.Set("AppId", "olts")
	req.Header.Set("AppSecret", a.titansApiSecret)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 打印响应状态
	log.Printf("响应状态: %s\n", resp.Status)

	// 打印响应体
	log.Println("响应体:")
	body, _ := io.ReadAll(resp.Body)

	availableQty, ok := getAvailableQty(body)
	if !ok {
		return
	}

	// 持仓不足
	if availableQty < order.OrderQty {
		order.ExtendAttrMap["insufficientPosition"] = true
	} else {
		order.ExtendAttrMap["insufficientPosition"] = false
	}

	jsText, err := json.MarshalToString(order.ExtendAttrMap)
	if err != nil {
		log.Printf("fail to MarshalToString, order.ExtendAttrMap:%v, error:%v\n", order.ExtendAttrMap, err)
	} else {
		order.ExtendAttr = jsText
	}
}
