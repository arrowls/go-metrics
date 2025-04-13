package main

import (
	"log"
	"net/http"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/config"
	"github.com/arrowls/go-metrics/internal/controller"
	"github.com/arrowls/go-metrics/internal/logger"
	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/arrowls/go-metrics/internal/repository"
	"github.com/arrowls/go-metrics/internal/service"
)

func main() {
	loggerInst := logger.NewLogger()
	loggerInst.Info("Starting application")

	errorHandler := apperrors.NewHTTPErrorHandler(loggerInst)

	serverConfig := config.NewServerConfig()
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	services := service.NewService(repo)
	controllers := controller.NewController(services, errorHandler)

	router := controllers.InitRoutes(loggerInst)

	log.Fatal(http.ListenAndServe(serverConfig.ServerEndpoint, router))
}
