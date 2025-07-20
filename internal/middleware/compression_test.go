package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCompressWriter(t *testing.T) {
	recorder := httptest.NewRecorder()
	writer := NewCompressWriter(recorder)

	if writer == nil {
		t.Fatal("NewCompressWriter returned nil")
	}

	if writer.responseWriter != recorder {
		t.Error("responseWriter not set correctly")
	}

	if writer.gzipWriter == nil {
		t.Error("gzipWriter not initialized")
	}
}

func TestCompressWriter_Write(t *testing.T) {
	recorder := httptest.NewRecorder()
	writer := NewCompressWriter(recorder)

	testData := []byte("test data")
	n, err := writer.Write(testData)

	if err != nil {
		t.Errorf("Write returned error: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}

	// Закрываем writer для завершения сжатия
	writer.Close()

	// Проверяем, что данные были сжаты
	if len(recorder.Body.Bytes()) == 0 {
		t.Error("No data written to response")
	}
}

func TestCompressWriter_Header(t *testing.T) {
	recorder := httptest.NewRecorder()
	writer := NewCompressWriter(recorder)

	header := writer.Header()
	if header == nil {
		t.Error("Header returned nil")
	}

	// Проверяем, что можем установить заголовки
	header.Set("Content-Type", "application/json")
	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Error("Header not set correctly")
	}
}

func TestCompressWriter_WriteHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	writer := NewCompressWriter(recorder)

	writer.WriteHeader(http.StatusOK)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestNewCompressReader(t *testing.T) {
	// Создаем сжатые данные
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write([]byte("test data"))
	gzipWriter.Close()

	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	compressReader, err := NewCompressReader(reader)

	if err != nil {
		t.Errorf("NewCompressReader returned error: %v", err)
	}

	if compressReader == nil {
		t.Fatal("NewCompressReader returned nil")
	}

	// Читаем данные
	data, err := io.ReadAll(compressReader)
	if err != nil {
		t.Errorf("ReadAll returned error: %v", err)
	}

	if string(data) != "test data" {
		t.Errorf("Expected 'test data', got '%s'", string(data))
	}

	compressReader.Close()
}

func TestCompressReader_Close(t *testing.T) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write([]byte("test data"))
	gzipWriter.Close()

	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	compressReader, err := NewCompressReader(reader)

	if err != nil {
		t.Fatalf("NewCompressReader returned error: %v", err)
	}

	err = compressReader.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestNewCompressionMiddleware(t *testing.T) {
	// Тест с Accept-Encoding: gzip
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	recorder := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.Write([]byte("test response"))
	})

	middleware := NewCompressionMiddleware(handler)
	middleware.ServeHTTP(recorder, req)

	if !handlerCalled {
		t.Error("Handler was not called")
	}

	// Проверяем, что установлен заголовок Content-Encoding
	if recorder.Header().Get("Content-Encoding") != "gzip" {
		t.Error("Content-Encoding header not set to gzip")
	}

	// Проверяем, что данные сжаты
	if len(recorder.Body.Bytes()) == 0 {
		t.Error("No response body")
	}
}

func TestNewCompressionMiddleware_WithoutGzip(t *testing.T) {
	// Тест без Accept-Encoding
	req := httptest.NewRequest("GET", "/test", nil)

	recorder := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.Write([]byte("test response"))
	})

	middleware := NewCompressionMiddleware(handler)
	middleware.ServeHTTP(recorder, req)

	if !handlerCalled {
		t.Error("Handler was not called")
	}

	// Проверяем, что заголовок Content-Encoding не установлен
	if recorder.Header().Get("Content-Encoding") != "" {
		t.Error("Content-Encoding header should not be set")
	}

	// Проверяем, что данные не сжаты
	if recorder.Body.String() != "test response" {
		t.Errorf("Expected 'test response', got '%s'", recorder.Body.String())
	}
}

func TestNewCompressionMiddleware_WithGzipBody(t *testing.T) {
	// Создаем сжатые данные для тела запроса
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write([]byte("compressed request body"))
	gzipWriter.Close()

	req := httptest.NewRequest("POST", "/test", bytes.NewReader(buf.Bytes()))
	req.Header.Set("Content-Encoding", "gzip")

	recorder := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read body: %v", err)
			return
		}
		if string(body) != "compressed request body" {
			t.Errorf("Expected 'compressed request body', got '%s'", string(body))
		}
		w.Write([]byte("test response"))
	})

	middleware := NewCompressionMiddleware(handler)
	middleware.ServeHTTP(recorder, req)

	if !handlerCalled {
		t.Error("Handler was not called")
	}
}
