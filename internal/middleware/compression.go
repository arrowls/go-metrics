// Package middleware содержит HTTP middleware компоненты.
// Включает сжатие, логирование, хеширование и другие промежуточные обработчики.
package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// gzipWriterPool - пул gzip writer'ов для переиспользования объектов.
// Уменьшает нагрузку на GC при частом использовании сжатия.
var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

// CompressWriter оборачивает http.ResponseWriter для сжатия ответов.
// Автоматически сжимает данные в gzip формат при отправке клиенту.
type CompressWriter struct {
	// responseWriter - оригинальный ResponseWriter
	responseWriter http.ResponseWriter
	// gzipWriter - gzip writer для сжатия данных
	gzipWriter *gzip.Writer
	// pooled - флаг использования объекта из пула
	pooled bool
}

// NewCompressWriter создает новый CompressWriter.
// Получает gzip writer из пула и настраивает его для записи в responseWriter.
func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	gw := gzipWriterPool.Get().(*gzip.Writer)
	gw.Reset(w)

	return &CompressWriter{
		responseWriter: w,
		gzipWriter:     gw,
		pooled:         true,
	}
}

func (c *CompressWriter) Header() http.Header {
	return c.responseWriter.Header()
}

func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.gzipWriter.Write(p)
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	c.responseWriter.WriteHeader(statusCode)
}

func (c *CompressWriter) Close() error {
	err := c.gzipWriter.Close()
	if c.pooled {
		gzipWriterPool.Put(c.gzipWriter)
	}
	return err
}

type CompressReader struct {
	reader     io.ReadCloser
	gzipReader *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		reader:     r,
		gzipReader: zr,
	}, nil
}

func (c *CompressReader) Read(p []byte) (n int, err error) {
	return c.gzipReader.Read(p)
}

func (c *CompressReader) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.gzipReader.Close()
}

func NewCompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := w

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			compressWriter := NewCompressWriter(w)
			writer = compressWriter
			defer func() {
				if errClose := compressWriter.Close(); errClose != nil {
					fmt.Printf("%v1", errClose)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			writer.Header().Set("Content-Encoding", "gzip")
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			reader, err := NewCompressReader(r.Body)
			if err != nil {
				fmt.Printf("%v2", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			defer func() {
				if errClose := reader.Close(); errClose != nil {
					fmt.Printf("%v3", errClose)

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			r.Body = reader
		}

		next.ServeHTTP(writer, r)
	})
}
