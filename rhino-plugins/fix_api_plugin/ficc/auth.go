package ficc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 鉴权函数
func Authenticate(loginServiceUrl, loginClientID, loginClientKey, phoneNum, password string) (bool, string, error) {
	// 创建请求体
	reqBody := map[string]string{
		"user_id":    phoneNum,
		"password":   password,
		"client_id":  loginClientID,
		"client_key": loginClientKey,
	}

	// 将请求体转换为JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return false, "", fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 发送HTTP请求
	resp, err := http.Post(loginServiceUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, "", fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("读取响应体失败: %v", err)
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("解析响应JSON失败: %v", err)
	}

	// 检查鉴权是否成功
	if errCode, ok := result["errCode"].(map[string]interface{}); ok {
		if code, ok := errCode["code"].(float64); ok && int(code) == 200 {
			// 提取JWT
			if data, ok := result["data"].(map[string]interface{}); ok {
				if jwt, ok := data["jwt"].(string); ok {
					return true, jwt, nil
				}
			}
			return true, "", nil // 成功但没有jwt的情况
		}
	}

	return false, "", nil
}