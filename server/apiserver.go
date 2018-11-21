package server

import (
	"net/http"
)

type APIServer struct {
	Mux    *http.ServeMux
	Server http.Server
}

func NewAPIServer(address string) *APIServer {
	return &APIServer{
		Mux: http.NewServeMux(),
		Server: http.Server{
			Addr: address,
		},
	}
}
