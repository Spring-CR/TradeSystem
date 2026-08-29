package server

import (
	"log"
	"net/http"
	"rhino-api/api/server/router"
)

var(
	server *http.Server
)

func StartApiServer(addr, positionServiceUrl string) error {
	server = &http.Server{
		Addr: addr,
		Handler: router.ApiRouter(positionServiceUrl),
	}
	log.Printf("start rhino api-server on address:%s\n", addr)
	return server.ListenAndServe()
}