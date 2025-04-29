package middleware

import (
	"net/http"
)

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
