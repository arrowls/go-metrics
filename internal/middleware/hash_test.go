package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrowls/go-metrics/internal/config"
	"github.com/sirupsen/logrus"
)

func TestNewHashingMiddleware_WithoutKey(t *testing.T) {
	serverConfig := config.ServerConfig{Key: ""}
	logger := logrus.New()

	middleware := NewHashingMiddleware(serverConfig, logger)

	req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte("test body")))
	recorder := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.Write([]byte("test response"))
	})

	middleware(handler).ServeHTTP(recorder, req)

	if !handlerCalled {
		t.Error("Handler was not called when no key provided")
	}
}

func TestNewHashingMiddleware_WithValidHash(t *testing.T) {
	serverConfig := config.ServerConfig{Key: "test-key"}
	logger := logrus.New()

	middleware := NewHashingMiddleware(serverConfig, logger)

	body := []byte("test body")
	hasher := hmac.New(sha256.New, []byte(serverConfig.Key))
	hasher.Write(body)
	hash := hex.EncodeToString(hasher.Sum(nil))

	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set(config.HashHeaderName, hash)
	recorder := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.Write([]byte("test response"))
	})

	middleware(handler).ServeHTTP(recorder, req)

	if !handlerCalled {
		t.Error("Handler was not called with valid hash")
	}

	// Проверяем, что в ответе установлен хеш
	responseHash := recorder.Header().Get(config.HashHeaderName)
	if responseHash == "" {
		t.Error("Response hash not set")
	}
}

func TestNewHashingMiddleware_WithInvalidHash(t *testing.T) {
	serverConfig := config.ServerConfig{Key: "test-key"}
	logger := logrus.New()

	middleware := NewHashingMiddleware(serverConfig, logger)

	body := []byte("test body")
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set(config.HashHeaderName, "invalid-hash")
	recorder := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.Write([]byte("test response"))
	})

	middleware(handler).ServeHTTP(recorder, req)

	if handlerCalled {
		t.Error("Handler was called with invalid hash")
	}

	// Проверяем, что установлен правильный Content-Type
	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type to be application/json")
	}

	// Проверяем, что тело ответа содержит сообщение об ошибке
	if len(recorder.Body.Bytes()) == 0 {
		t.Error("Expected error response body")
	}
}

func TestNewHashingMiddleware_WithoutHashHeader(t *testing.T) {
	serverConfig := config.ServerConfig{Key: "test-key"}
	logger := logrus.New()

	middleware := NewHashingMiddleware(serverConfig, logger)

	body := []byte("test body")
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.Write([]byte("test response"))
	})

	middleware(handler).ServeHTTP(recorder, req)

	if !handlerCalled {
		t.Error("Handler was not called without hash header")
	}
}

func TestNewHashingMiddleware_EmptyBody(t *testing.T) {
	serverConfig := config.ServerConfig{Key: "test-key"}
	logger := logrus.New()

	middleware := NewHashingMiddleware(serverConfig, logger)

	req := httptest.NewRequest("GET", "/test", nil)
	recorder := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.Write([]byte("test response"))
	})

	middleware(handler).ServeHTTP(recorder, req)

	if !handlerCalled {
		t.Error("Handler was not called with empty body")
	}
}

func TestResponseWriter_Write(t *testing.T) {
	recorder := httptest.NewRecorder()
	writer := &responseWriter{
		ResponseWriter: recorder,
		body:           []byte{},
	}

	testData := []byte("test data")
	n, err := writer.Write(testData)

	if err != nil {
		t.Errorf("Write returned error: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}

	if len(writer.body) != len(testData) {
		t.Errorf("Expected body length %d, got %d", len(testData), len(writer.body))
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	writer := &responseWriter{
		ResponseWriter: recorder,
	}

	writer.WriteHeader(http.StatusNotFound)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected recorder status %d, got %d", http.StatusNotFound, recorder.Code)
	}
}

func TestNewResponseRecorder(t *testing.T) {
	recorder := httptest.NewRecorder()
	responseRecorder := NewResponseRecorder(recorder)

	if responseRecorder == nil {
		t.Fatal("NewResponseRecorder returned nil")
	}

	if responseRecorder.ResponseWriter != recorder {
		t.Error("ResponseWriter not set correctly")
	}

	if responseRecorder.status != http.StatusOK {
		t.Errorf("Expected initial status %d, got %d", http.StatusOK, responseRecorder.status)
	}

	if responseRecorder.responseSize != 0 {
		t.Errorf("Expected initial response size 0, got %d", responseRecorder.responseSize)
	}
}
