// Package logger содержит компоненты для логирования приложения.
// Предоставляет единообразный интерфейс для записи логов в JSON формате.
package logger

import (
	"github.com/arrowls/go-metrics/internal/di"
	log "github.com/sirupsen/logrus"
)

// diKey - ключ для хранения логгера в контейнере зависимостей
const diKey = "logger"

// ProvideLogger предоставляет экземпляр логгера из контейнера зависимостей.
// Если логгер не существует в контейнере, создает новый с JSON форматированием.
// Настраивает уровень логирования Info и добавляет логгер в контейнер.
// Возвращает настроенный экземпляр logrus.Logger.
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
