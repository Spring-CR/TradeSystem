package router

import (
	"log"
	"net/http/httputil"
	"net/url"
	"rhino-api/api/api_const"
	"time"

	"github.com/gin-gonic/gin"
)

func setPositionRouter(e*gin.Engine, positionServiceUrl string) {

	url, err := url.Parse(positionServiceUrl)
    if err != nil {
		time.Sleep(3*time.Second)
        log.Fatalf("Failed to parse target URL %s: %v", positionServiceUrl, err)
    }

	proxy := httputil.NewSingleHostReverseProxy(url)

	r := e.Group(api_const.RoutePositions)
	{
		r.Any("/*proxyPath", func(c *gin.Context) {
			// 设定代理的目标路径
			c.Request.URL.Path = "/query/posFundRedis"
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}
}