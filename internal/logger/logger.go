package logger

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type loggerKey struct{}

func NewLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetLevel(log.InfoLevel)

	return logger
}

func Provide(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), loggerKey{}, NewLogger())
	return r.WithContext(ctx)
}

func Inject(ctx context.Context) *log.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*log.Logger); ok {
		return logger
	}
	return NewLogger()
}
