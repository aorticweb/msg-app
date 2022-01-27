package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type API struct {
	db     *gorm.DB
	router *mux.Router
	logger *log.Logger
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func NewAPI(db *gorm.DB, logger *log.Logger) *API {
	a := &API{db, mux.NewRouter(), logger}
	a.routes()
	return a
}

func (a *API) routes() {
	a.router.HandleFunc("/health", a.handleHealth()).Methods("GET")
}
