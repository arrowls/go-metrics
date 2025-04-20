package main

import (
	"log"
	"net/http"
	"strings"
	"time"

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
	repo := repository.WithRestore(
		repository.NewRepository(storage),
		time.Duration(serverConfig.StoreInterval)*time.Second,
		serverConfig.StorageFilePath,
		serverConfig.Restore,
		loggerInst,
	)
	services := service.NewService(repo)
	controllers := controller.NewController(services, errorHandler)

	router := controllers.InitRoutes(loggerInst)

	log.Fatal(http.ListenAndServe(serverConfig.ServerEndpoint, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "//")
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		router.ServeHTTP(w, r)
	})))
}
