package server

import (
	"net/http"

	"github.com/tlmiller/disttrust/server/healthz"
)

type ApiServer struct {
	healthzChecks []healthz.Checker
	server        http.Server
	stopCh        chan struct{}
}

func NewApiServer(address string) *ApiServer {
	return &ApiServer{
		server: http.Server{
			Addr: address,
		},
		stopCh: make(chan struct{}),
	}
}
