package server

func (a *ApiServer) Serve() {
	a.server.Handler = a.Mux
	a.server.ListenAndServe()
}

func (a *ApiServer) Stop() {
	a.server.Close()
}
