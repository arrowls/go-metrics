package main

import (
	"log"
	"net/http"

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

	serverConfig := config.NewServerConfig()
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	services := service.NewService(repo)
	controllers := controller.NewController(services)

	router := controllers.InitRoutes(loggerInst)

	log.Fatal(http.ListenAndServe(serverConfig.ServerEndpoint, router))
}
