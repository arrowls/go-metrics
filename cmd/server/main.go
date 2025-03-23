package main

import (
	"log"
	"net/http"

	"github.com/arrowls/go-metrics/cmd/server/controller"
	"github.com/arrowls/go-metrics/cmd/server/repository"
	"github.com/arrowls/go-metrics/cmd/server/service"
	"github.com/arrowls/go-metrics/internal/memstorage"
)

func main() {
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	services := service.NewService(repo)
	controllers := controller.NewController(services)

	router := controllers.InitRoutes()

	log.Fatal(http.ListenAndServe(":8080", router))
}
