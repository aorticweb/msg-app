package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	c "github.com/aorticweb/msg-app/app/common"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) *c.APIResponse

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
			a.logger.Printf("[%s] %s response [%d]: %s", r.Method, r.URL.Path, resp.Code, time.Now().Sub(startTime))
		}()
		if http.StatusOK <= resp.Code && resp.Code < 300 {
			a.okResponse(w, resp.Code, resp.Data)
			return
		}
		if resp.Err != nil {
			a.logger.Println(resp.Err)
		}
		if resp.Message != "" {
			http.Error(w, resp.Message, resp.Code)
			return
		}
		w.WriteHeader(resp.Code)
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
