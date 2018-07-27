package server

import (
	"net/http"

	"github.com/tlmiller/disttrust/server/healthz"
)

func (a *ApiServer) Serve() {
	mux := http.NewServeMux()
	healthz.InstallHandler(mux, a.healthzChecks...)
	a.server.Handler = mux
	a.server.ListenAndServe()
}
