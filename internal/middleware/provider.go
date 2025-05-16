package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/arrowls/go-metrics/internal/di"
	"github.com/arrowls/go-metrics/internal/logger"
)

const diKey = "logging_middleware"

func ProvideLoggingMiddleware(container di.ContainerInterface) func(next http.Handler) http.Handler {
	middlewareInst := container.Get(diKey)
	if middleware, ok := middlewareInst.(func(next http.Handler) http.Handler); ok {
		return middleware
	}

	loggerInst := logger.ProvideLogger(container)

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			rw := NewResponseRecorder(w)
			next.ServeHTTP(rw, r)

			elapsed := time.Since(now)

			loggerInst.Infof("%s %s | Status %d %d bytes, took %v\n", r.Method, r.URL.Path, rw.status, rw.responseSize, elapsed)
		})
	}

	if err := container.Add(diKey, middleware); err != nil {
		log.Fatal(err)
	}

	return middleware
}
