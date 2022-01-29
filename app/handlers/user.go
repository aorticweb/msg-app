package api

import (
	"encoding/json"
	"net/http"

	c "github.com/aorticweb/msg-app/app/common"
	"github.com/aorticweb/msg-app/app/crud"
)

func (a *API) handleUserPost() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *c.APIResponse {
		var userInput crud.User
		err := json.NewDecoder(r.Body).Decode(&userInput)
		if err != nil {
			return c.NewBadResponse(http.StatusBadRequest, "invalid request", c.WrapError("JSON decoding error", err))
		}
		if err = a.validate.Struct(userInput); err != nil {
			return &c.InvalidRequestResponse
		}
		exist, err := crud.UserExist(a.db, userInput.Username)
		if err != nil {
			return c.NewBadResponse(http.StatusInternalServerError, "", c.WrapError("failed to query users", err))
		}
		if exist {
			return c.NewBadResponse(http.StatusConflict, "user with the same username already registered", nil)
		}
		err = crud.CreateUser(a.db, userInput)
		if err != nil {
			return c.NewBadResponse(http.StatusInternalServerError, "", c.WrapError("failed to create user", err))
		}
		return c.NewGoodResponse(http.StatusCreated, userInput)
	}
}
