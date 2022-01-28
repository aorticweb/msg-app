package api

import (
	"encoding/json"
	"net/http"

	c "github.com/aorticweb/msg-app/app/common"
	"github.com/aorticweb/msg-app/app/crud"
	"github.com/aorticweb/msg-app/app/model"
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

func (a *API) handleMessageGet() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *c.APIResponse {
		messageID, err := c.GetIDFromRequest(r)
		if err != nil {
			return c.NewBadResponse(http.StatusBadRequest, "invalid request", nil)
		}
		dbMessage, exist, err := crud.GetMessage(a.db, messageID)
		if err != nil {
			return c.NewBadResponse(http.StatusInternalServerError, "", err)
		}
		if !exist {
			return c.NewBadResponse(http.StatusNotFound, "message not found", nil)
		}
		respMessage := m.ResponseMessageFromDBMessage(dbMessage)
		return c.NewGoodResponse(http.StatusOK, respMessage)
	}
}

func (a *API) handleMessageRepliesGet() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *c.APIResponse {
		messageID, err := c.GetIDFromRequest(r)
		if err != nil {
			return c.NewBadResponse(http.StatusBadRequest, "invalid request", nil)
		}
		_, exist, err := crud.GetMessage(a.db, messageID)
		if err != nil {
			return c.NewBadResponse(http.StatusInternalServerError, "", err)
		}
		if !exist {
			return c.NewBadResponse(http.StatusNotFound, "message not found", nil)
		}
		dbMessages, err := crud.GetMessageReplies(a.db, messageID)
		var data []model.Message

		for _, msg := range dbMessages {
			data = append(data, *m.ResponseMessageFromDBMessage(&msg))
		}
		return c.NewGoodResponse(http.StatusOK, data)
	}
}

func (a *API) handleInboxGet() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *c.APIResponse {
		username, err := c.GetUsernameFromRequest(r)
		if err != nil {
			return c.NewBadResponse(http.StatusBadRequest, "invalid request", nil)
		}
		user, exist, err := crud.FindUser(a.db, username)
		if err != nil {
			return c.NewBadResponse(http.StatusInternalServerError, "", c.WrapError("failed to query user", err))
		}
		if !exist {
			return c.NewBadResponse(http.StatusNotFound, "user with given username does not exist", nil)
		}
		dbMessages, err := crud.GetUserMailbox(a.db, user.ID)
		var data []model.Message

		for _, msg := range dbMessages {
			data = append(data, *m.ResponseMessageFromDBMessage(&msg))
		}
		return c.NewGoodResponse(http.StatusOK, data)
	}
}
