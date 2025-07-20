package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkNewCompressWriter(b *testing.B) {
	recorder := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewCompressWriter(recorder)
	}
}

func BenchmarkCompressWriter_Write(b *testing.B) {
	recorder := httptest.NewRecorder()
	writer := NewCompressWriter(recorder)
	testData := []byte("test data for compression benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer.Write(testData)
	}
	writer.Close()
}

func BenchmarkNewCompressReader(b *testing.B) {
	// Создаем сжатые данные
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write([]byte("test data"))
	gzipWriter.Close()

	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewCompressReader(reader)
	}
}

func BenchmarkCompressReader_Read(b *testing.B) {
	// Создаем сжатые данные
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write([]byte("test data for reading benchmark"))
	gzipWriter.Close()

	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	compressReader, _ := NewCompressReader(reader)

	buffer := make([]byte, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compressReader.Read(buffer)
	}
	compressReader.Close()
}

func BenchmarkNewCompressionMiddleware(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test response"))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewCompressionMiddleware(handler)
	}
}

func BenchmarkCompressionMiddleware_ServeHTTP_WithGzip(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test response for compression middleware benchmark"))
	})

	middleware := NewCompressionMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	recorder := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		middleware.ServeHTTP(recorder, req)
	}
}

func BenchmarkCompressionMiddleware_ServeHTTP_WithoutGzip(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test response for compression middleware benchmark"))
	})

	middleware := NewCompressionMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	recorder := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		middleware.ServeHTTP(recorder, req)
	}
}

func BenchmarkCompressionMiddleware_ServeHTTP_WithGzipBody(b *testing.B) {
	// Создаем сжатые данные для тела запроса
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write([]byte("compressed request body for benchmark"))
	gzipWriter.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, r.Body)
	})

	middleware := NewCompressionMiddleware(handler)

	req := httptest.NewRequest("POST", "/test", bytes.NewReader(buf.Bytes()))
	req.Header.Set("Content-Encoding", "gzip")
	recorder := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		middleware.ServeHTTP(recorder, req)
	}
}

func BenchmarkGzipCompression(b *testing.B) {
	testData := []byte("test data for gzip compression benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)
		gzipWriter.Write(testData)
		gzipWriter.Close()
	}
}

func BenchmarkGzipDecompression(b *testing.B) {
	// Создаем сжатые данные
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write([]byte("test data for gzip decompression benchmark"))
	gzipWriter.Close()

	compressedData := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(compressedData)
		gzipReader, _ := gzip.NewReader(reader)
		io.Copy(io.Discard, gzipReader)
		gzipReader.Close()
	}
}
