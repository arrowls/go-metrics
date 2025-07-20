// Package apperrors содержит определения ошибок приложения.
// Включает стандартные ошибки и функции для формирования JSON ответов с ошибками.
package apperrors

import (
	"encoding/json"
	"errors"
	"strings"
)

// Стандартные ошибки приложения
var (
	// ErrNotFound - ошибка когда запрашиваемый ресурс не найден
	ErrNotFound = errors.New("not found")
	// ErrBadRequest - ошибка когда запрос содержит некорректные данные
	ErrBadRequest = errors.New("bad request")
	// ErrUnknown - неизвестная ошибка сервера
	ErrUnknown = errors.New("unknown error")
)

// ErrorResponseWithMessage представляет структуру JSON ответа с ошибкой.
// Используется для возврата ошибок в HTTP API.
type ErrorResponseWithMessage struct {
	// Message - текстовое сообщение об ошибке
	Message string `json:"message"`
}

// ErrorResponse формирует JSON ответ с сообщением об ошибке.
// Очищает сообщение от переносов строк и приводит первую букву к верхнему регистру.
// Возвращает сериализованный JSON в виде байтового слайса.
func ErrorResponse(message string) []byte {
	message = strings.Replace(message, "\n", "", -1)
	message = strings.Trim(message, " \n")

	if len(message) > 0 {
		messageSplit := strings.Split(message, "")
		message = strings.ToUpper(messageSplit[0]) + strings.Join(messageSplit[1:], "")
	}

	response := ErrorResponseWithMessage{
		message,
	}
	responseBytes, _ := json.Marshal(response)
	return responseBytes
}
