package main

import (
	"log"
	"net/http"

	"github.com/arrowls/go-metrics/cmd/server/controller"
	"github.com/arrowls/go-metrics/cmd/server/repository"
	"github.com/arrowls/go-metrics/cmd/server/service"
)

func main() {
	repo := repository.NewRepository()
	services := service.NewService(repo)
	controllers := controller.NewController(services)

	mux := controllers.InitRoutes()

	log.Fatal(http.ListenAndServe(":8080", mux))
}
