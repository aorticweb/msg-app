package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aorticweb/msg-app/app/crud"
	"github.com/aorticweb/msg-app/app/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createGroup(t *testing.T, db *gorm.DB, users []crud.User) *crud.Group {
	group, err := crud.CreateGroup(db, groupname, users)
	require.NoError(t, err)
	return group
}

func messageReplySuccess(t *testing.T, sender *crud.User) model.ComposedMessage {
	return model.ComposedMessage{
		ReplyMessage: model.ReplyMessage{
			Sender:  sender.Username,
			Subject: "Greetings",
			Body:    "You are Hired",
		},
	}
}

func messageUserSuccess(t *testing.T, sender *crud.User, recipient *crud.User) model.ComposedMessage {
	return model.ComposedMessage{
		Recipient: map[string]string{"username": recipient.Username},
		ReplyMessage: model.ReplyMessage{
			Sender:  sender.Username,
			Subject: "Greetings",
			Body:    "You are Hired",
		},
	}
}

func messageGroupSuccess(t *testing.T, sender *crud.User, recipient *crud.Group) model.ComposedMessage {
	return model.ComposedMessage{
		Recipient: map[string]string{"groupname": recipient.Groupname},
		ReplyMessage: model.ReplyMessage{
			Sender:  sender.Username,
			Subject: "New Hire",
			Body:    "Hey team, please welcome our new hire",
		},
	}
}

func assertDBMessage(t *testing.T, dbMsg *crud.Message, msg *model.Message) {
	require.Equal(t, msg.Body, dbMsg.Body)
	require.Equal(t, msg.Subject, dbMsg.Subject)
	require.Equal(t, msg.Sender, dbMsg.Sender.Username)

	groupname, ok1 := msg.Recipient["groupname"]
	username, ok2 := msg.Recipient["username"]
	require.True(t, ok1 || ok2)
	if ok1 {
		require.Equal(t, groupname, dbMsg.Group.Groupname)
	} else {
		require.Equal(t, username, dbMsg.Recipient.Username)
	}
	if msg.RE != nil {
		require.Equal(t, msg.RE, dbMsg.REID)
	} else {
		require.True(t, dbMsg.REID == nil)
	}
}

func TestSendMessageToUser(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)

	msg := messageUserSuccess(t, &users[0], &users[1])
	resp, err := http.Post(url(srv.URL, "/messages"), "application/json", toPayload(t, msg))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data model.Message
	err = json.NewDecoder(resp.Body).Decode(&data)

	require.NoError(t, err)
	require.Equal(t, msg.Body, data.Body)
	require.Equal(t, msg.Subject, data.Subject)
	require.Equal(t, msg.Sender, data.Sender)
	require.Equal(t, msg.Recipient, data.Recipient)
	require.True(t, data.RE == nil)

	sentAtDiff := (time.Now().UTC().Sub(data.SentAt))
	require.True(t, sentAtDiff < time.Second)

	dbMsg, exist, err := crud.GetMessage(db, data.ID)
	require.NoError(t, err)
	require.True(t, exist)
	assertDBMessage(t, dbMsg, &data)
}

func TestSendMessageToUserRecipientUserNotFound(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)

	msg := messageUserSuccess(t, &users[0], &users[1])
	msg.Recipient["username"] = "Fake User"
	resp, err := http.Post(url(srv.URL, "/messages"), "application/json", toPayload(t, msg))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	var data model.Message
	err = json.NewDecoder(resp.Body).Decode(&data)

	err = db.First(&crud.Message{}).Error
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestSendMessageToUserSenderNotFound(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)

	msg := messageUserSuccess(t, &users[0], &users[1])
	msg.Sender = "Fake User"
	resp, err := http.Post(url(srv.URL, "/messages"), "application/json", toPayload(t, msg))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	var data model.Message
	err = json.NewDecoder(resp.Body).Decode(&data)

	err = db.First(&crud.Message{}).Error
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestMessageToGroup(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)
	group := createGroup(t, db, users)

	msg := messageGroupSuccess(t, &users[0], group)
	resp, err := http.Post(url(srv.URL, "/messages"), "application/json", toPayload(t, msg))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data model.Message
	err = json.NewDecoder(resp.Body).Decode(&data)

	require.NoError(t, err)
	require.Equal(t, msg.Body, data.Body)
	require.Equal(t, msg.Subject, data.Subject)
	require.Equal(t, msg.Sender, data.Sender)
	require.Equal(t, msg.Recipient, data.Recipient)
	require.True(t, data.RE == nil)

	sentAtDiff := (time.Now().UTC().Sub(data.SentAt))
	require.True(t, sentAtDiff < time.Second)

	dbMsg, exist, err := crud.GetMessage(db, data.ID)
	require.NoError(t, err)
	require.True(t, exist)
	assertDBMessage(t, dbMsg, &data)
}

func TestSendMessageToGroupRecipientNotFound(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)
	group := createGroup(t, db, users)

	msg := messageGroupSuccess(t, &users[0], group)
	msg.Recipient["groupname"] = "Fake Group"
	resp, err := http.Post(url(srv.URL, "/messages"), "application/json", toPayload(t, msg))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	var data model.Message
	err = json.NewDecoder(resp.Body).Decode(&data)

	err = db.First(&crud.Message{}).Error
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestSendReplyMessageToUser(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)
	existingMsg := crud.Message{
		Sender:    &users[0],
		Recipient: &users[1],
		Subject:   "Waiting for reply",
		Body:      "Please Respond",
		SentAt:    time.Now().UTC(),
	}
	_, err := crud.CreateMessage(db, &existingMsg)
	require.NoError(t, err)

	msg := messageReplySuccess(t, &users[1])
	resp, err := http.Post(url(srv.URL, fmt.Sprintf("/messages/%d/replies", existingMsg.ID)), "application/json", toPayload(t, msg))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data model.Message
	err = json.NewDecoder(resp.Body).Decode(&data)

	require.NoError(t, err)
	require.Equal(t, msg.Body, data.Body)
	require.Equal(t, msg.Subject, data.Subject)
	require.Equal(t, msg.Sender, data.Sender)

	require.True(t, *data.RE == existingMsg.ID)

	sentAtDiff := (time.Now().UTC().Sub(data.SentAt))
	require.True(t, sentAtDiff < time.Second)

	dbMsg, exist, err := crud.GetMessage(db, data.ID)
	require.NoError(t, err)
	require.True(t, exist)
	assertDBMessage(t, dbMsg, &data)
}

func TestSendReplyMessageToGroup(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)
	group := createGroup(t, db, users)
	existingMsg := crud.Message{
		Sender:  &users[0],
		Group:   group,
		Subject: "Waiting for reply",
		Body:    "Please Respond",
		SentAt:  time.Now().UTC(),
	}
	_, err := crud.CreateMessage(db, &existingMsg)
	require.NoError(t, err)

	msg := messageReplySuccess(t, &users[1])
	resp, err := http.Post(url(srv.URL, fmt.Sprintf("/messages/%d/replies", existingMsg.ID)), "application/json", toPayload(t, msg))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data model.Message
	err = json.NewDecoder(resp.Body).Decode(&data)

	require.NoError(t, err)
	require.Equal(t, msg.Body, data.Body)
	require.Equal(t, msg.Subject, data.Subject)
	require.Equal(t, msg.Sender, data.Sender)

	require.True(t, *data.RE == existingMsg.ID)

	sentAtDiff := (time.Now().UTC().Sub(data.SentAt))
	require.True(t, sentAtDiff < time.Second)

	dbMsg, exist, err := crud.GetMessage(db, data.ID)
	require.NoError(t, err)
	require.True(t, exist)
	assertDBMessage(t, dbMsg, &data)
}

func TestSendReplyMessageToUserMessageNotFound(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)

	msg := messageReplySuccess(t, &users[1])
	resp, err := http.Post(url(srv.URL, "/messages/1/replies"), "application/json", toPayload(t, msg))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
