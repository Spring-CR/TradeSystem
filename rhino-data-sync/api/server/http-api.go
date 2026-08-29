package server

import (
	"fmt"
	"log"
	"net/http"
	"rhino-core/domain_cfg"
	"rhino-data-sync/api/server/handlers"
	"rhino-data-sync/api/server/router"
)

type HttpApiServer struct {
	server *http.Server
}

func NewHttpApiServer(port int, apiToken string, dataSyncCfg *domain_cfg.DataSyncConfig) (server *HttpApiServer) {
	dataSyncNotify := handlers.NewDataSyncNotifyHandler(dataSyncCfg)
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.GetApiRouters(dataSyncNotify, apiToken),
	}
	log.Printf("start rhino data-sync on port:%d\n", port)
	server = &HttpApiServer{server: s}
	return
}

func (s *HttpApiServer) Start() (err error) {
	return s.server.ListenAndServe()
}