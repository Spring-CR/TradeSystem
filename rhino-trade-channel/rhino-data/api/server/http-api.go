package server

import (
	"fmt"
	"log"
	"net/http"
	"rhino-data/api/server/handlers"
	"rhino-data/api/server/router"
	"rhino-data/datamap"
)

type HttpApiServer struct {
	server *http.Server
}

func NewHttpApiServer(port int, apiToken string, autoSyncRepo *datamap.AutoSyncRepo) (server *HttpApiServer) {
	query := handlers.NewDataQueryHandler(autoSyncRepo)
	syncLog := handlers.NewSyncLogHandler(autoSyncRepo)
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.GetApiRouters(query, syncLog, apiToken),
	}
	log.Printf("start rhino data service on port:%d\n", port)
	server = &HttpApiServer{server: s}
	return
}

func (s *HttpApiServer) Start() (err error) {
	return s.server.ListenAndServe()
}