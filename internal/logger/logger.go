package logger

import (
	log "github.com/sirupsen/logrus"
)

func NewLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetLevel(log.InfoLevel)

	return logger
}
