package common

import "net/http"

type APIResponse struct {
	Code    int
	Data    interface{}
	Message string
	Err     error
}

func NewBadResponse(code int, message string, err error) *APIResponse {
	return &APIResponse{code, nil, message, err}
}

func NewGoodResponse(code int, data interface{}) *APIResponse {
	return &APIResponse{code, data, "", nil}
}

var InvalidRequestResponse APIResponse = *NewBadResponse(http.StatusBadRequest, "invalid request", nil)
