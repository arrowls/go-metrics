package main

import (
	"log"
	"net/http"

	"github.com/arrowls/go-metrics/internal/controller"
	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/arrowls/go-metrics/internal/repository"
	"github.com/arrowls/go-metrics/internal/service"
)

func main() {
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	services := service.NewService(repo)
	controllers := controller.NewController(services)

	router := controllers.InitRoutes()

	log.Fatal(http.ListenAndServe(":8080", router))
}
