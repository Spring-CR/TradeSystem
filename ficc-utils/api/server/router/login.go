package router

import (
	"bytes"
	"encoding/json"
	"ficc-utils/api/api_const"
	"ficc-utils/api/options"
	"ficc-utils/common/domain_error"
	"ficc-utils/common/server/middleware"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

// curl -XPOST -H "Content-Type: application/json"  -d '{"phoneNum":"17300000514","password":"a123123"}' http://10.51.136.21:5094/api/v1/login
func setLoginRouter(e *gin.Engine, loginServiceUrl, loginClientID, loginClientKey string, credProvider *CredProvider) {
	r := e.Group(api_const.RouteLogin)
	{
		r.POST("", func(c *gin.Context) {

			loginOpt := &options.LoginOpt{}
			if !middleware.BindInputOption(c, loginOpt) {
				return
			}

			if loginOpt.PhoneNum == "" || loginOpt.Password == "" {
				de := domain_error.Build(domain_error.ILLEGAL_LOGIN_OPTION_ERR_CODE, nil)
				middleware.ProcessDomainError(de, c)
				return
			}

			ok, _, err := authenticate(loginServiceUrl, loginClientID, loginClientKey, loginOpt.PhoneNum, loginOpt.Password)
			if err != nil {
				de := domain_error.Build(domain_error.LOGIN_FAIL_ERR_CODE, err)
				middleware.ProcessDomainError(de, c)
				return
			}
			if !ok {
				de := domain_error.Build(domain_error.LOGIN_FAIL_ERR_CODE, nil)
				middleware.ProcessDomainError(de, c)
				return
			}

			apiToken, date, accountMap, _ := credProvider.getTokenKey(loginOpt.PhoneNum)
			if len(accountMap) == 0 {
				de := domain_error.Build(domain_error.ACCOUNT_NOT_CONFIG_ERR_CODE, nil)
				middleware.ProcessDomainError(de, c)
				return
			}

			var keyList []int
			for k := range accountMap {
				keyList = append(keyList, k)
			}
			sort.Slice(keyList, func(i, j int) bool { return keyList[i] < keyList[j] })

			loginResult := options.LoginResult{
				ApiToken: apiToken,
				EffectiveDate: date,
				Accounts: keyList,
			}

			middleware.ResponseJson(c, loginResult)
		})
	}
}

// 鉴权函数
func authenticate(loginServiceUrl, loginClientID, loginClientKey, phoneNum, password string) (bool, string, error) {
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
