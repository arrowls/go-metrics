package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

type CompressWriter struct {
	responseWriter http.ResponseWriter
	gzipWriter     *gzip.Writer
	pooled         bool
}

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
