package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	api "github.com/aorticweb/msg-app/app/handlers"

	"github.com/aorticweb/msg-app/app/crud"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var groupname string = "Griffondor"
var usernames = []string{"Harry", "Hermione", "Ron"}

func groupSuccessPayload(t *testing.T) io.Reader {
	jsonStr := []byte(fmt.Sprintf(`{"groupname": "%s", "usernames": ["Harry", "Hermione", "Ron"]}`, groupname))
	return bytes.NewBuffer(jsonStr)
}

func groupUserNotRegisterPayload(t *testing.T) io.Reader {
	jsonStr := []byte(fmt.Sprintf(`{"groupname": "%s", "usernames": ["Harry", "Hermione", "Ron", "Malfoy"]}`, groupname))
	return bytes.NewBuffer(jsonStr)
}
func createUsers(t *testing.T, db *gorm.DB) {
	users := []crud.User{}
	for _, name := range usernames {
		users = append(users, crud.User{Username: name})
	}
	err := db.Create(&users).Error
	require.NoError(t, err)
}

func TestGroupPostSuccess(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)

	createUsers(t, db)
	resp, err := http.Post(url(srv.URL, "/groups"), "application/json", groupSuccessPayload(t))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data api.GroupPost
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
	var data api.GroupPost
	err = json.NewDecoder(resp.Body).Decode(&data)

	require.NoError(t, err)
	require.Equal(t, groupname, data.Groupname)
	require.Equal(t, usernames, data.Usernames)

	resp, err = http.Post(url(srv.URL, "/groups"), "application/json", groupSuccessPayload(t))
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
}
