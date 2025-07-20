package apperrors

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		expectedPrefix string
	}{
		{
			name:           "simple message",
			message:        "test error",
			expectedPrefix: "Test error",
		},
		{
			name:           "message with newlines",
			message:        "test\nerror\nmessage",
			expectedPrefix: "Testerrormessage",
		},
		{
			name:           "message with spaces",
			message:        "  test  error  message  ",
			expectedPrefix: "Test  error  message",
		},
		{
			name:           "empty message",
			message:        "",
			expectedPrefix: "",
		},
		{
			name:           "single character",
			message:        "a",
			expectedPrefix: "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ErrorResponse(tt.message)

			assert.NotEmpty(t, result, "Expected non-empty response")

			// Парсим JSON
			var response ErrorResponseWithMessage
			err := json.Unmarshal(result, &response)
			assert.NoError(t, err, "Failed to unmarshal JSON")

			// Проверяем, что сообщение начинается с заглавной буквы
			assert.True(t, strings.HasPrefix(response.Message, tt.expectedPrefix),
				"Expected message to start with '%s', got '%s'", tt.expectedPrefix, response.Message)

			// Проверяем, что нет переносов строк
			assert.False(t, strings.Contains(response.Message, "\n"), "Message should not contain newlines")
		})
	}
}

func TestErrorResponseWithMessage_JSON(t *testing.T) {
	message := "test error message"
	response := ErrorResponseWithMessage{
		Message: message,
	}

	// Сериализуем в JSON
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err, "Failed to marshal JSON")

	// Десериализуем обратно
	var parsedResponse ErrorResponseWithMessage
	err = json.Unmarshal(jsonData, &parsedResponse)
	assert.NoError(t, err, "Failed to unmarshal JSON")

	assert.Equal(t, message, parsedResponse.Message, "Expected message '%s', got '%s'", message, parsedResponse.Message)
}

func TestErrorResponse_Formatting(t *testing.T) {
	// Тест с сообщением, содержащим различные символы
	message := "Test\nError\nMessage\nWith\nNewlines"
	result := ErrorResponse(message)

	var response ErrorResponseWithMessage
	err := json.Unmarshal(result, &response)
	assert.NoError(t, err, "Failed to unmarshal JSON")

	// Проверяем, что переносы строк удалены
	assert.False(t, strings.Contains(response.Message, "\n"), "Newlines should be removed from message")

	// Проверяем, что сообщение начинается с заглавной буквы
	assert.True(t, strings.HasPrefix(response.Message, "TestErrorMessageWithNewlines"),
		"Expected message to start with 'TestErrorMessageWithNewlines', got '%s'", response.Message)
}

func TestErrorResponse_EdgeCases(t *testing.T) {
	// Тест с пустой строкой
	result := ErrorResponse("")

	var response ErrorResponseWithMessage
	err := json.Unmarshal(result, &response)
	assert.NoError(t, err, "Failed to unmarshal JSON")

	assert.Empty(t, response.Message, "Expected empty message")

	// Тест с одним символом
	result = ErrorResponse("a")

	err = json.Unmarshal(result, &response)
	assert.NoError(t, err, "Failed to unmarshal JSON")

	assert.Equal(t, "A", response.Message, "Expected 'A'")

	// Тест с пробелами
	result = ErrorResponse("  test  ")

	err = json.Unmarshal(result, &response)
	assert.NoError(t, err, "Failed to unmarshal JSON")

	assert.Equal(t, "Test", response.Message, "Expected 'Test'")
}

func TestErrorConstants(t *testing.T) {
	// Проверяем, что константы ошибок определены
	assert.NotNil(t, ErrNotFound, "ErrNotFound should be defined")
	assert.NotNil(t, ErrBadRequest, "ErrBadRequest should be defined")
	assert.NotNil(t, ErrUnknown, "ErrUnknown should be defined")

	// Проверяем, что это разные ошибки
	assert.NotEqual(t, ErrNotFound, ErrBadRequest, "ErrNotFound and ErrBadRequest should be different")
	assert.NotEqual(t, ErrBadRequest, ErrUnknown, "ErrBadRequest and ErrUnknown should be different")
	assert.NotEqual(t, ErrNotFound, ErrUnknown, "ErrNotFound and ErrUnknown should be different")
}
