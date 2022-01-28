package api

import (
	"encoding/json"
	"net/http"

	"github.com/aorticweb/msg-app/app/crud"
)

type GroupPost struct {
	Groupname string   `json:"groupname" validate:"required"`
	Usernames []string `json:"usernames" validate:"required"`
}

func (a *API) handleGroupPost() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *APIResponse {
		var groupInput GroupPost
		err := json.NewDecoder(r.Body).Decode(&groupInput)
		if err != nil {
			return newBadResponse(http.StatusBadRequest, "invalid request", wrapError("JSON decoding error", err))
		}
		if err = a.validate.Struct(groupInput); err != nil {
			return newBadResponse(http.StatusBadRequest, "invalid request", nil)
		}

		users, err := crud.FindUsers(a.db, groupInput.Usernames)
		if err != nil {
			return newBadResponse(http.StatusInternalServerError, "", wrapError("failed to query users", err))
		}
		if len(users) != len(groupInput.Usernames) {
			return newBadResponse(http.StatusConflict, "one or more group member username does not exist", nil)
		}

		exist, err := crud.GroupExists(a.db, groupInput.Groupname)
		if err != nil {
			return newBadResponse(http.StatusInternalServerError, "", wrapError("failed to query groups", err))
		}
		if exist {
			return newBadResponse(http.StatusConflict, "group with the same Groupname already registered", nil)
		}

		err = crud.CreateGroup(a.db, groupInput.Groupname, users)
		if err != nil {
			return newBadResponse(http.StatusInternalServerError, "", wrapError("failed to create group", err))
		}
		return newGoodResponse(http.StatusCreated, groupInput)
	}
}
