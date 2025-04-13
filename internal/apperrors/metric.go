package apperrors

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrNotFound   = errors.New("")
	ErrBadRequest = errors.New("")
	ErrUnknown    = errors.New("")
)

type ErrorResponseWithMessage struct {
	Message string `json:"message"`
}

func ErrorResponse(message string) string {
	response := ErrorResponseWithMessage{
		strings.Trim(strings.Replace(message, "\n", "", -1), " \n"),
	}
	responseBytes, _ := json.Marshal(response)
	return string(responseBytes)
}
