package middleware

import (
	"net/http"
	"rhino-common/domain_error"

	"github.com/gin-gonic/gin"
)

var(
	unauthorizedResponse =[]byte(`{"code":"`+domain_error.API_UNAUTHORIZED_ERR_CODE+`","msg":"API鉴权失败"}`)
)

func SetCredentials(apiToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 直接从HeaderMap获取（避免GetHeader的额外处理）
		tokens := c.Request.Header["X-Api-Token"]
		if len(tokens) == 0 {
			abortUnauthorized(c)
			return
		}

		// 2. 字节级比较（比字符串转换更高效）
		if apiToken != tokens[0] {
			abortUnauthorized(c)
			return
		}

		// 3. 鉴权通过
		c.Next()
	}
}

// 独立的失败处理函数（减少重复代码）
func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatus(http.StatusUnauthorized)
	c.Writer.Header().Set("Content-Type", "application/json")
	_, _ = c.Writer.Write(unauthorizedResponse)
}