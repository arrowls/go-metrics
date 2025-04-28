package middleware

import (
	"net/http"

	"github.com/arrowls/go-metrics/internal/logger"
)

func NewProvideLoggerMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = logger.Provide(r)
			next.ServeHTTP(w, r)
		})
	}
}
