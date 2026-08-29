package http

import (
	"fmt"
	"log"
	"net/http"
	"rhino-core/order_domain"
	"rhino-engine/api/api_adapter"
	"rhino-engine/api/http/server/handlers"
	"rhino-engine/api/http/server/router"
)

type HttpApiServer struct {
	server *http.Server
}

func NewHttpApiServer(port int, apiAdapter api_adapter.APIAdapter, engine *order_domain.OrderEngine, workerCount int) (server *HttpApiServer) {

	orderHandler := handlers.NewOrderHandler(apiAdapter, engine, workerCount)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.GetApiRouters(orderHandler, engine.GetApiToken()),
	}
	log.Printf("start rhino api-server on port:%d\n", port)
	server = &HttpApiServer{server: s}
	return
}

func (s *HttpApiServer) Start() (err error) {
	return s.server.ListenAndServe()
}
