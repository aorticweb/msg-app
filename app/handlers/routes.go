package api

func (a *API) routes() {
	a.router.HandleFunc("/groups", a.middleware(a.handleGroupPost())).Methods("POST")

	a.router.HandleFunc("/health", a.middleware(a.handleHealth())).Methods("GET")

	a.router.HandleFunc("/messages/{id}", a.middleware(a.handleMessageGet())).Methods("GET")
	a.router.HandleFunc("/messages", a.middleware(a.handleMessagePost())).Methods("POST")

	a.router.HandleFunc("/messages/{id}/replies", a.middleware(a.handleMessageRepliesGet())).Methods("GET")
	a.router.HandleFunc("/messages/{id}/replies", a.middleware(a.handleMessageReplyPost())).Methods("POST")

	a.router.HandleFunc("/users", a.middleware(a.handleUserPost())).Methods("POST")

	a.router.HandleFunc("/users/{username}/mailbox", a.middleware(a.handleInboxGet())).Methods("GET")
}
