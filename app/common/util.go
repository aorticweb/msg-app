package common

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetIDFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return 0, errors.New("id not found in request")
	}
	return strconv.ParseInt(id, 10, 64)
}

func GetUsernameFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		return "", errors.New("username not found in request")
	}
	return username, nil
}

func WrapError(context string, e error) error {
	return fmt.Errorf("%s: %s", context, e)
}
