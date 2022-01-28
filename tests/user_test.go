package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/aorticweb/msg-app/app/crud"
	"github.com/stretchr/testify/require"
)

var user crud.User = crud.User{Username: "Bobby"}

func userSuccessPayload(t *testing.T) io.Reader {
	jsonStr := []byte(fmt.Sprintf(`{"username": "%s"}`, user.Username))
	return bytes.NewBuffer(jsonStr)
}

func userMissingUsernamePayload(t *testing.T) io.Reader {
	jsonStr := []byte(fmt.Sprintf(`{"invalid": "%s"}`, user.Username))
	return bytes.NewBuffer(jsonStr)
}

func userInvalidPayload(t *testing.T) io.Reader {
	jsonStr := []byte(fmt.Sprintf(`invalid{"username": "%s"}`, user.Username))
	return bytes.NewBuffer(jsonStr)
}

func TestUserPostSuccess(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)

	resp, err := http.Post(url(srv.URL, "/users"), "application/json", userSuccessPayload(t))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data crud.User
	err = json.NewDecoder(resp.Body).Decode(&data)

	require.NoError(t, err)
	require.Equal(t, user.Username, data.Username)
}

func TestUserPostFailsForDuplicate(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)

	resp, err := http.Post(url(srv.URL, "/users"), "application/json", userSuccessPayload(t))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data crud.User
	err = json.NewDecoder(resp.Body).Decode(&data)

	require.NoError(t, err)
	require.Equal(t, user.Username, data.Username)

	resp, err = http.Post(url(srv.URL, "/users"), "application/json", userSuccessPayload(t))
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestUserPostFailsMissingUsername(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)

	resp, err := http.Post(url(srv.URL, "/users"), "application/json", userMissingUsernamePayload(t))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUserPostInvalidPayload(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)

	resp, err := http.Post(url(srv.URL, "/users"), "application/json", userInvalidPayload(t))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
