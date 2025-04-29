package apperrors

import (
	"log"

	"github.com/arrowls/go-metrics/internal/di"
	"github.com/arrowls/go-metrics/internal/logger"
)

const diKey = "http_error_handler"

func ProvideHTTPErrorHandler(container di.ContainerInterface) *HTTPErrorHandler {
	handlerInst := container.Get(diKey)
	if handler, ok := handlerInst.(*HTTPErrorHandler); ok {
		return handler
	}

	loggerInst := logger.ProvideLogger(container)
	handler := NewHTTPErrorHandler(loggerInst)

	if err := container.Add(diKey, handler); err != nil {
		log.Fatal(err)
	}

	return handler
}
