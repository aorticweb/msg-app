package api

import (
	"net/http"
)

func (a *API) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c int
		a.db.Raw("SELECT COUNT(table_name) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'user';").Scan(&c)
		if c != 1 {
			a.logger.Error("Health check failed, could not reach database")
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}
