package server

func (a *APIServer) Serve() {
	a.Server.Handler = a.Mux
	a.Server.ListenAndServe()
}

func (a *APIServer) Stop() {
	a.Server.Close()
}
