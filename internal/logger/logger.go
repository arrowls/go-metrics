package logger

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func NewLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetLevel(log.DebugLevel)

	today := time.Now().Format("2006-01-02")

	logFile, err := os.OpenFile(fmt.Sprintf("./logs/%s.log", today), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger.SetOutput(logFile)
	return logger
}
