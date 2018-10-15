package server

import (
	"net/http"
)

type ApiServer struct {
	Mux    *http.ServeMux
	server http.Server
}

func NewApiServer(address string) *ApiServer {
	return &ApiServer{
		Mux: http.NewServeMux(),
		server: http.Server{
			Addr: address,
		},
	}
}
