package server

import (
	"ficc-utils/api/server/router"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HttpApiServer struct {
	server *http.Server
}
func NewHttpApiServer(port int, dataQryServiceUrl, loginServiceUrl, loginClientID, loginClientKey, capitalServiceUrl, positionServiceUrl, tradeOrdersServiceUrl, hisTradeOrdersServiceUrl, execReportServiceUrl, pricingHubYtmConvUrl, pricingHubAuth, dataQryUrl, appId, appSecret string, queryCycle time.Duration, local *time.Location) (server *HttpApiServer) {

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.GetApiRouters(dataQryServiceUrl, loginServiceUrl, loginClientID, loginClientKey, capitalServiceUrl, positionServiceUrl, tradeOrdersServiceUrl, hisTradeOrdersServiceUrl, execReportServiceUrl, pricingHubYtmConvUrl, pricingHubAuth, dataQryUrl, appId, appSecret, queryCycle, local),
	}
	log.Printf("start ficc-utils on port:%d\n", port)
	server = &HttpApiServer{server: s}
	return
}

func (s *HttpApiServer) Start() (err error) {
	return s.server.ListenAndServe()
}
