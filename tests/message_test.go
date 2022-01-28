package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/aorticweb/msg-app/app/crud"
	"github.com/aorticweb/msg-app/app/model"
	"github.com/stretchr/testify/require"
	"github.com/thanhpk/randstr"
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
	require.Equal(t, existingMsg.Sender.Username, data.Recipient["username"])
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

func TestGetMessageSentToUser(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)
	existingMsg := crud.Message{
		Sender:    &users[0],
		Recipient: &users[1],
		Subject:   "Super Message",
		Body:      "Super Mario",
		SentAt:    time.Now().UTC(),
	}
	_, err := crud.CreateMessage(db, &existingMsg)
	require.NoError(t, err)

	resp, err := http.Get(url(srv.URL, fmt.Sprintf("/messages/%d", existingMsg.ID)))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var data model.Message
	err = json.NewDecoder(resp.Body).Decode(&data)

	sentAtDiff := (time.Now().UTC().Sub(data.SentAt))
	require.True(t, sentAtDiff < time.Second)

	dbMsg, exist, err := crud.GetMessage(db, data.ID)
	require.NoError(t, err)
	require.True(t, exist)
	assertDBMessage(t, dbMsg, &data)
}

func TestGetMessageSentToUserWithRe(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)
	baseMsg := crud.Message{
		Sender:    &users[0],
		Recipient: &users[1],
		Subject:   "Super Message",
		Body:      "Super Mario",
		SentAt:    time.Now().UTC(),
	}
	_, err := crud.CreateMessage(db, &baseMsg)
	require.NoError(t, err)
	existingMsg := crud.Message{
		Sender:    &users[1],
		Recipient: &users[1], // If sent as POST recipient is the same as base message
		REID:      &baseMsg.ID,
		Subject:   "Reply Super Message",
		Body:      "Super Luigi",
		SentAt:    time.Now().UTC(),
	}
	_, err = crud.CreateMessage(db, &existingMsg)
	require.NoError(t, err)

	resp, err := http.Get(url(srv.URL, fmt.Sprintf("/messages/%d", existingMsg.ID)))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var data model.Message
	err = json.NewDecoder(resp.Body).Decode(&data)

	sentAtDiff := (time.Now().UTC().Sub(data.SentAt))
	require.True(t, sentAtDiff < time.Second)

	dbMsg, exist, err := crud.GetMessage(db, data.ID)
	require.NoError(t, err)
	require.True(t, exist)
	assertDBMessage(t, dbMsg, &data)
}

func TestGetMessageSentToUserMessageNotFound(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)
	existingMsg := crud.Message{
		Sender:    &users[0],
		Recipient: &users[1],
		Subject:   "Super Message",
		Body:      "Super Mario",
		SentAt:    time.Now().UTC(),
	}
	_, err := crud.CreateMessage(db, &existingMsg)
	require.NoError(t, err)

	resp, err := http.Get(url(srv.URL, fmt.Sprintf("/messages/%d", 150)))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetMessageReplies(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)

	baseMsg := crud.Message{
		Sender:    &users[0],
		Recipient: &users[1],
		Subject:   "Super Message",
		Body:      "Super Mario",
		SentAt:    time.Now().UTC(),
	}
	err := db.Create(&baseMsg).Error
	require.NoError(t, err)
	messages := []crud.Message{
		{
			Sender:    &users[1],
			Recipient: &users[1],
			Subject:   "Super Message 2",
			Body:      "Super Mario 2",
			SentAt:    time.Now().UTC(),
			REID:      &baseMsg.ID,
		},
		{
			Sender:    &users[0],
			Recipient: &users[1],
			Subject:   "Super Message 3",
			Body:      "Super Mario 3",
			SentAt:    time.Now().UTC(),
			REID:      &baseMsg.ID,
		},
		{
			Sender:    &users[2],
			Recipient: &users[1],
			Subject:   "Super Message 4",
			Body:      "Super Mario 4",
			SentAt:    time.Now().UTC(),
		},
	}
	err = db.Create(&messages).Error
	require.NoError(t, err)

	resp, err := http.Get(url(srv.URL, fmt.Sprintf("/messages/%d/replies", baseMsg.ID)))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var data []model.Message
	err = json.NewDecoder(resp.Body).Decode(&data)
	require.True(t, len(data) == 2)

	// TODO:
	// add tests on messages returned
}

func TestGetMessageRepliesMessageNotFound(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)

	baseMsg := crud.Message{
		Sender:    &users[0],
		Recipient: &users[1],
		Subject:   "Super Message",
		Body:      "Super Mario",
		SentAt:    time.Now().UTC(),
	}
	err := db.Create(&baseMsg).Error
	require.NoError(t, err)

	messages := []crud.Message{
		{
			Sender:    &users[1],
			Recipient: &users[1],
			Subject:   "Super Message 2",
			Body:      "Super Mario 2",
			SentAt:    time.Now().UTC(),
		},
		{
			Sender:    &users[0],
			Recipient: &users[1],
			Subject:   "Super Message 3",
			Body:      "Super Mario 3",
			SentAt:    time.Now().UTC(),
		},
		{
			Sender:    &users[2],
			Recipient: &users[1],
			Subject:   "Super Message 4",
			Body:      "Super Mario 4",
			SentAt:    time.Now().UTC(),
		},
	}
	err = db.Create(&messages).Error
	require.NoError(t, err)

	resp, err := http.Get(url(srv.URL, fmt.Sprintf("/messages/%d/replies", 150)))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetUserMailbox(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)
	users := createUsers(t, db)
	group := createGroup(t, db, users)
	group_not_a_member := createGroup(t, db, users[1:])

	baseMsgs := []crud.Message{
		{ // part of the inbox
			Sender:    &users[1],
			Recipient: &users[0],
			Subject:   randstr.String(50),
			Body:      randstr.String(100),
			SentAt:    time.Now().UTC(),
		},
		{
			// part of the inbox
			Sender:  &users[2],
			Group:   group,
			Subject: randstr.String(150),
			Body:    randstr.String(340),
			SentAt:  time.Now().UTC(),
		},
		{
			Sender:  &users[0],
			Group:   group_not_a_member,
			Subject: randstr.String(12),
			Body:    randstr.String(600),
			SentAt:  time.Now().UTC(),
		},
		{
			Sender:  &users[1],
			Group:   group_not_a_member,
			Subject: randstr.String(12),
			Body:    randstr.String(600),
			SentAt:  time.Now().UTC(),
		},
	}
	err := db.Create(&baseMsgs).Error
	require.NoError(t, err)

	childMsgs := []crud.Message{
		{
			// User Reply to received message
			Sender:    &users[0],
			Recipient: &users[1],
			Subject:   randstr.String(50),
			Body:      randstr.String(100),
			SentAt:    time.Now().UTC(),
			REID:      &baseMsgs[0].ID,
		},
		{
			// part of the inbox
			// Other user replies to group
			Sender:  &users[1],
			Group:   group,
			Subject: randstr.String(150),
			Body:    randstr.String(340),
			SentAt:  time.Now().UTC(),
			REID:    &baseMsgs[1].ID,
		},
		{
			// Other user replies message sent by user[0] in a
			//  group user[0] is not a member of
			Sender:  &users[1],
			Group:   group_not_a_member,
			Subject: randstr.String(12),
			Body:    randstr.String(600),
			SentAt:  time.Now().UTC(),
			REID:    &baseMsgs[2].ID,
		},
		{
			// Other user replies to group user[0] is not a member of
			Sender:  &users[2],
			Group:   group_not_a_member,
			Subject: randstr.String(12),
			Body:    randstr.String(600),
			SentAt:  time.Now().UTC(),
		},
	}
	err = db.Create(&childMsgs).Error
	require.NoError(t, err)

	validMsgs := []crud.Message{baseMsgs[0], baseMsgs[1], childMsgs[1]}
	resp, err := http.Get(url(srv.URL, fmt.Sprintf("/users/%s/mailbox", users[0].Username)))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var data []model.Message
	err = json.NewDecoder(resp.Body).Decode(&data)
	log.Println(len(data))
	require.True(t, len(data) == len(validMsgs))

	dataMap := map[int64]model.Message{}
	for _, msg := range data {
		dataMap[msg.ID] = msg
	}
	var found bool
	for _, validMsg := range validMsgs {
		_, found = dataMap[validMsg.ID]
		require.True(t, found)
	}
	// TODO:
	// add tests on messages returned
}
