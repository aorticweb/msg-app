package api

import (
	"errors"
	"net/http"
)

func (a *API) handleHealth() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *APIResponse {
		var c int
		a.db.Raw("SELECT COUNT(table_name) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'user';").Scan(&c)
		if c != 1 {
			err := errors.New("Health check failed, could not reach database")
			return newBadResponse(http.StatusServiceUnavailable, "Unavailable Ressource", err)
		}
		return &APIResponse{http.StatusOK, nil, "", nil}
	}
}
