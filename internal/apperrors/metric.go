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

func ErrorResponse(message string) []byte {
	message = strings.Replace(message, "\n", "", -1)
	message = strings.Trim(message, " \n")

	messageSplit := strings.Split(message, "")
	message = strings.ToUpper(messageSplit[0]) + strings.Join(messageSplit[1:], "")

	response := ErrorResponseWithMessage{
		message,
	}
	responseBytes, _ := json.Marshal(response)
	return responseBytes
}
