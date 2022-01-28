package model

import (
	"net/http"
	"time"

	c "github.com/aorticweb/msg-app/app/common"
	"github.com/aorticweb/msg-app/app/crud"
	"gorm.io/gorm"
)

type ReplyMessage struct {
	Sender  string `json:"sender" validate:"required"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

func (rm *ReplyMessage) ValidateSender(db *gorm.DB) (*crud.User, *c.APIResponse) {
	sender, exist, err := crud.FindUser(db, rm.Sender)
	if err != nil {
		return nil, c.NewBadResponse(http.StatusInternalServerError, "", c.WrapError("failed to query users", err))
	}
	if !exist {
		return nil, c.NewBadResponse(http.StatusNotFound, "user with given username does not exist", nil)
	}
	return sender, nil
}

func (rm *ReplyMessage) Validate(db *gorm.DB, reID int64) (*crud.Message, *c.APIResponse) {
	msg := crud.Message{
		Subject: rm.Subject,
		Body:    rm.Body,
		SentAt:  time.Now().UTC(),
	}
	sender, badResp := rm.ValidateSender(db)
	if badResp != nil {
		return nil, badResp
	}
	msg.Sender = sender
	reMessage, exist, err := crud.GetMessage(db, reID)
	if err != nil {
		return nil, c.NewBadResponse(http.StatusInternalServerError, "", c.WrapError("failed to query message", err))
	}
	if !exist {
		return nil, c.NewBadResponse(http.StatusNotFound, "message with given id does not exist", nil)
	}
	msg.REID = &reMessage.ID
	return &msg, nil
}

type ComposedMessage struct {
	ReplyMessage
	Recipient map[string]string `json:"recipient" validate:"required"` // Either crud.User or crud.Group
}

func (m *ComposedMessage) Validate(db *gorm.DB) (*crud.Message, *c.APIResponse) {
	msg := crud.Message{
		Subject: m.Subject,
		Body:    m.Body,
		SentAt:  time.Now().UTC(),
	}
	username, usernameFound := m.Recipient["username"]
	groupname, groupnameFound := m.Recipient["groupname"]
	if usernameFound && groupnameFound {
		return nil, c.NewBadResponse(http.StatusBadRequest, "invalid request", nil)
	}
	sender, badResp := m.ValidateSender(db)
	if badResp != nil {
		return nil, badResp
	}
	msg.Sender = sender
	if usernameFound {
		user, exist, err := crud.FindUser(db, username)
		if err != nil {
			return nil, c.NewBadResponse(http.StatusInternalServerError, "", c.WrapError("failed to query user", err))
		}
		if !exist {
			return nil, c.NewBadResponse(http.StatusNotFound, "recipient user with given username does not exist", nil)
		}
		msg.Recipient = user
		return &msg, nil
	}
	if groupnameFound {
		group, exist, err := crud.FindGroup(db, groupname)
		if err != nil {
			return nil, c.NewBadResponse(http.StatusInternalServerError, "", c.WrapError("failed to query group", err))
		}
		if !exist {
			return nil, c.NewBadResponse(http.StatusNotFound, "recipient group with given groupname does not exist", nil)
		}
		msg.Group = group
		return &msg, nil
	}
	return nil, c.NewBadResponse(http.StatusBadRequest, "invalid request", nil)
}

type Message struct {
	ComposedMessage
	ID     int64     `json:"id" validate:"required"`
	RE     *int64    `json:"re"`
	SentAt time.Time `json:"sent_at" validate:"required"`
}

func ResponseMessageFromDBMessage(m *crud.Message) *Message {
	msg := Message{
		ID: m.ID,
		ComposedMessage: ComposedMessage{
			ReplyMessage: ReplyMessage{
				Subject: m.Subject,
				Body:    m.Body,
				Sender:  m.Sender.Username,
			},
			Recipient: make(map[string]string),
		},
		SentAt: m.SentAt,
	}
	// Purposefully not raising an error here if
	// both user and group are missing because of db constraint
	if m.Group != nil {
		msg.Recipient["groupname"] = m.Group.Groupname
	} else if m.Recipient != nil {
		msg.Recipient["username"] = m.Recipient.Username
	}
	if m.REID != nil {
		msg.RE = m.REID
	}
	return &msg
}
