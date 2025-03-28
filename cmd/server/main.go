package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/arrowls/go-metrics/internal/controller"
	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/arrowls/go-metrics/internal/repository"
	"github.com/arrowls/go-metrics/internal/service"
)

const serverEndpointDefault = "localhost:8080"

type Config struct {
	ServerEndpoint *string
}

var config Config

func InitConfig() {
	config.ServerEndpoint = flag.String("a", serverEndpointDefault, "server endpoint url")

	flag.Parse()
}

func main() {
	InitConfig()
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	services := service.NewService(repo)
	controllers := controller.NewController(services)

	router := controllers.InitRoutes()

	log.Fatal(http.ListenAndServe(*config.ServerEndpoint, router))
}
