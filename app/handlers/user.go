package api

import (
	"encoding/json"
	"net/http"

	"github.com/aorticweb/msg-app/app/crud"
)

func (a *API) handleUserPost() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *APIResponse {
		var userInput crud.User
		err := json.NewDecoder(r.Body).Decode(&userInput)
		if err != nil {
			return newBadResponse(http.StatusBadRequest, "invalid request", wrapError("JSON decoding error", err))
		}
		if err = a.validate.Struct(userInput); err != nil {
			return newBadResponse(http.StatusBadRequest, "invalid request", nil)
		}
		exist, err := crud.UserExist(a.db, userInput.Username)
		if err != nil {
			return newBadResponse(http.StatusInternalServerError, "", wrapError("failed to query users", err))
		}
		if exist {
			return newBadResponse(http.StatusConflict, "user with the same username already registered", nil)
		}
		err = crud.CreateUser(a.db, userInput)
		if err != nil {
			return newBadResponse(http.StatusInternalServerError, "", wrapError("failed to create user", err))
		}
		return newGoodResponse(http.StatusCreated, userInput)
	}
}
