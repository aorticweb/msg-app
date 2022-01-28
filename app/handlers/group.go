package api

import (
	"encoding/json"
	"net/http"

	c "github.com/aorticweb/msg-app/app/common"
	"github.com/aorticweb/msg-app/app/crud"
	m "github.com/aorticweb/msg-app/app/model"
)

func (a *API) handleGroupPost() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *c.APIResponse {
		var groupInput m.GroupPost
		err := json.NewDecoder(r.Body).Decode(&groupInput)
		if err != nil {
			return c.NewBadResponse(http.StatusBadRequest, "invalid request", c.WrapError("JSON decoding error", err))
		}
		if err = a.validate.Struct(groupInput); err != nil {
			return c.NewBadResponse(http.StatusBadRequest, "invalid request", nil)
		}

		users, err := crud.FindUsers(a.db, groupInput.Usernames)
		if err != nil {
			return c.NewBadResponse(http.StatusInternalServerError, "", c.WrapError("failed to query users", err))
		}
		if len(users) != len(groupInput.Usernames) {
			return c.NewBadResponse(http.StatusConflict, "one or more group member username does not exist", nil)
		}

		exist, err := crud.GroupExists(a.db, groupInput.Groupname)
		if err != nil {
			return c.NewBadResponse(http.StatusInternalServerError, "", c.WrapError("failed to query groups", err))
		}
		if exist {
			return c.NewBadResponse(http.StatusConflict, "group with the same Groupname already registered", nil)
		}

		err = crud.CreateGroup(a.db, groupInput.Groupname, users)
		if err != nil {
			return c.NewBadResponse(http.StatusInternalServerError, "", c.WrapError("failed to create group", err))
		}
		return c.NewGoodResponse(http.StatusCreated, groupInput)
	}
}
