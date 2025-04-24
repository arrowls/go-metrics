package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type CompressWriter struct {
	responseWriter http.ResponseWriter
	gzipWriter     *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		responseWriter: w,
		gzipWriter:     gzip.NewWriter(w),
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
	return c.gzipWriter.Close()
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
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			writer.Header().Set("Content-Encoding", "gzip")
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			reader, err := NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			defer func() {
				if errClose := reader.Close(); errClose != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			r.Body = reader
		}

		next.ServeHTTP(writer, r)
	})
}
