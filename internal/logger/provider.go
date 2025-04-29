package logger

import (
	"github.com/arrowls/go-metrics/internal/di"
	log "github.com/sirupsen/logrus"
)

const diKey = "logger"

func ProvideLogger(container di.ContainerInterface) *log.Logger {
	loggerInst := container.Get(diKey)

	if logger, ok := loggerInst.(*log.Logger); ok {
		return logger
	}

	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetLevel(log.InfoLevel)

	err := container.Add(diKey, logger)
	if err != nil {
		log.Fatal(err)
	}

	return logger
}
