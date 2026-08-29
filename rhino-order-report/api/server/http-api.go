package server

import (
	"fmt"
	"log"
	"net/http"
	"rhino-core/domain_cfg"
	"rhino-order-report/api/server/handlers"
	"rhino-order-report/api/server/router"
)

type HttpApiServer struct {
	server *http.Server
}

func NewHttpApiServer(port int, applicationCfg *domain_cfg.ApplicationCfg) (server *HttpApiServer) {
	currentReportHandler := handlers.NewCurrentReportHandler(applicationCfg)
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.GetApiRouters(currentReportHandler, applicationCfg.GetApiToken()),
	}
	log.Printf("start rhino order-report on port:%d\n", port)
	server = &HttpApiServer{server: s}
	return
}

func (s *HttpApiServer) Start() (err error) {
	return s.server.ListenAndServe()
}
