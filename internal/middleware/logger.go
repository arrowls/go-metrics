package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLoggingMiddleware(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			rw := NewResponseRecorder(w)
			next.ServeHTTP(rw, r)

			elapsed := time.Since(now)

			logger.Infof("%s %s | Status %d %d bytes, took %v\n", r.Method, r.URL.Path, rw.status, rw.responseSize, elapsed)
		})
	}
}

type ResponseRecorder struct {
	http.ResponseWriter
	status       int
	responseSize int
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{ResponseWriter: w, status: http.StatusOK}
}

func (rw *ResponseRecorder) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.responseSize += size
	return size, err
}

func (rw *ResponseRecorder) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
