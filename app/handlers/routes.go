package api

func (a *API) routes() {
	a.router.HandleFunc("/health", a.middleware(a.handleHealth())).Methods("GET")
	a.router.HandleFunc("/users", a.middleware(a.handleUserPost())).Methods("POST")
	a.router.HandleFunc("/groups", a.middleware(a.handleGroupPost())).Methods("POST")
}
