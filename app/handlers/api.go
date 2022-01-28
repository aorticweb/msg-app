package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type APIResponse struct {
	code    int
	data    interface{}
	message string
	err     error
}

func newBadResponse(code int, message string, err error) *APIResponse {
	return &APIResponse{code, nil, message, err}
}

func newGoodResponse(code int, data interface{}) *APIResponse {
	return &APIResponse{code, data, "", nil}
}

type HandlerFunc func(http.ResponseWriter, *http.Request) *APIResponse

type API struct {
	db       *gorm.DB
	router   *mux.Router
	logger   *log.Logger
	validate *validator.Validate
}

func NewAPI(db *gorm.DB, logger *log.Logger) *API {
	a := &API{db, mux.NewRouter(), logger, validator.New()}
	a.routes()
	return a
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *API) middleware(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		resp := next(w, r)
		defer func() {
			a.logger.Printf("[%s] %s response [%d]: %s", r.Method, r.URL.Path, resp.code, time.Now().Sub(startTime))
		}()
		if http.StatusOK <= resp.code && resp.code < 300 {
			a.okResponse(w, resp.code, resp.data)
			return
		}
		if resp.err != nil {
			a.logger.Println(resp.err)
		}
		if resp.message != "" {
			http.Error(w, resp.message, resp.code)
			return
		}
		w.WriteHeader(resp.code)
	}
}

func (a *API) okResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func wrapError(context string, e error) error {
	return fmt.Errorf("%s: %s", context, e)
}
