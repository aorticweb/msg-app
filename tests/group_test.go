package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/aorticweb/msg-app/app/model"

	"github.com/aorticweb/msg-app/app/crud"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var groupname string = "Griffondor"
var usernames = []string{"Harry", "Hermione", "Ron"}

func groupSuccessPayload(t *testing.T) io.Reader {
	group := model.GroupPost{
		Groupname: groupname,
		Usernames: usernames,
	}
	return toPayload(t, group)
}

func groupUserNotRegisterPayload(t *testing.T) io.Reader {
	group := model.GroupPost{
		Groupname: groupname,
		Usernames: append(usernames, "Malfoy"),
	}
	return toPayload(t, group)
}

func createUsers(t *testing.T, db *gorm.DB) []crud.User {
	users := []crud.User{}
	for _, name := range usernames {
		users = append(users, crud.User{Username: name})
	}
	err := db.Create(&users).Error
	require.NoError(t, err)
	return users
}

func TestGroupPostSuccess(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)

	createUsers(t, db)
	resp, err := http.Post(url(srv.URL, "/groups"), "application/json", groupSuccessPayload(t))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data model.GroupPost
	err = json.NewDecoder(resp.Body).Decode(&data)

	require.NoError(t, err)
	require.Equal(t, groupname, data.Groupname)
	require.Equal(t, usernames, data.Usernames)
}

func TestGroupPostUserNotRegistered(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)

	createUsers(t, db)
	resp, err := http.Post(url(srv.URL, "/groups"), "application/json", groupUserNotRegisterPayload(t))
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestGroupPostGroupAlreadyExist(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)

	createUsers(t, db)
	resp, err := http.Post(url(srv.URL, "/groups"), "application/json", groupSuccessPayload(t))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data model.GroupPost
	err = json.NewDecoder(resp.Body).Decode(&data)

	require.NoError(t, err)
	require.Equal(t, groupname, data.Groupname)
	require.Equal(t, usernames, data.Usernames)

	resp, err = http.Post(url(srv.URL, "/groups"), "application/json", groupSuccessPayload(t))
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
}
