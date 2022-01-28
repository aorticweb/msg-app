package api

import (
	"encoding/json"
	"net/http"

	c "github.com/aorticweb/msg-app/app/common"
	"github.com/aorticweb/msg-app/app/crud"
	m "github.com/aorticweb/msg-app/app/model"
)

func (a *API) handleMessagePost() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *c.APIResponse {
		var messageInput m.ComposedMessage
		err := json.NewDecoder(r.Body).Decode(&messageInput)
		if err != nil {
			return c.NewBadResponse(http.StatusBadRequest, "invalid request", c.WrapError("JSON decoding error", err))
		}
		if err = a.validate.Struct(messageInput); err != nil {
			return c.NewBadResponse(http.StatusBadRequest, "invalid request", nil)
		}
		message, badResp := messageInput.Validate(a.db)
		if badResp != nil {
			return badResp
		}
		dbMessage, err := crud.CreateMessage(a.db, message)
		if err != nil {
			return c.NewBadResponse(http.StatusNotFound, "", c.WrapError("failed to create message", err))
		}
		respMessage := m.ResponseMessageFromDBMessage(dbMessage)
		return c.NewGoodResponse(http.StatusCreated, respMessage)
	}
}

func (a *API) handleMessageReplyPost() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *c.APIResponse {
		var messageInput m.ReplyMessage
		err := json.NewDecoder(r.Body).Decode(&messageInput)
		if err != nil {
			return c.NewBadResponse(http.StatusBadRequest, "invalid request", c.WrapError("JSON decoding error", err))
		}
		if err = a.validate.Struct(messageInput); err != nil {
			return c.NewBadResponse(http.StatusBadRequest, "invalid request", nil)
		}
		reID, err := c.GetIDFromRequest(r)
		if err != nil {
			return c.NewBadResponse(http.StatusBadRequest, "invalid request", nil)
		}
		message, badResp := messageInput.Validate(a.db, reID)
		if badResp != nil {
			return badResp
		}
		dbMessage, err := crud.CreateMessage(a.db, message)
		if err != nil {
			return c.NewBadResponse(http.StatusNotFound, "", c.WrapError("failed to create message", err))
		}
		respMessage := m.ResponseMessageFromDBMessage(dbMessage)
		return c.NewGoodResponse(http.StatusCreated, respMessage)
	}
}
