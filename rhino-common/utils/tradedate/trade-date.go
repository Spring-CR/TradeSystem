package tradedate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GetTradeDay 获取交易日
func GetTradeDay(apiURL, tradeDate string, days int, exchange, appId, appSecret string) (string, error) {
	// 构建请求URL
	params := url.Values{}
	params.Add("tradeDate", tradeDate)
	params.Add("days", fmt.Sprintf("%d", days))
	params.Add("exchange", exchange)

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	// 创建请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("AppId", appId)
	req.Header.Set("AppSecret", appSecret)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var result struct {
		ErrCode struct {
			Code int `json:"code"`
		} `json:"errCode"`
		Data struct {
			CalDate string `json:"calDate"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v, body:%s", err, body)
	}

	// 检查错误码
	if result.ErrCode.Code != 200 {
		return "", fmt.Errorf("接口返回错误，状态码: %d", result.ErrCode.Code)
	}

	return result.Data.CalDate, nil
}