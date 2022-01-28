package api

import (
	"errors"
	"net/http"

	c "github.com/aorticweb/msg-app/app/common"
)

func (a *API) handleHealth() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *c.APIResponse {
		var i int
		a.db.Raw("SELECT COUNT(table_name) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'user';").Scan(&i)
		if i != 1 {
			err := errors.New("Health check failed, could not reach database")
			return c.NewBadResponse(http.StatusServiceUnavailable, "Unavailable Ressource", err)
		}
		return &c.APIResponse{http.StatusOK, nil, "", nil}
	}
}
