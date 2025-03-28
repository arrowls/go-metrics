package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arrowls/go-metrics/internal/controller"
	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/arrowls/go-metrics/internal/repository"
	"github.com/arrowls/go-metrics/internal/service"
	"github.com/caarlos0/env/v7"
)

const (
	serverEndpointDefault = "localhost:8080"
)

type Config struct {
	ServerEndpoint string `env:"ADDRESS"`
}

var config Config

func InitConfig() {
	flag.StringVar(&config.ServerEndpoint, "a", serverEndpointDefault, "server endpoint url")

	flag.Parse()

	if err := env.Parse(&config); err != nil {
		fmt.Printf("Failed to parse env: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	InitConfig()
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	services := service.NewService(repo)
	controllers := controller.NewController(services)

	router := controllers.InitRoutes()

	log.Fatal(http.ListenAndServe(config.ServerEndpoint, router))
}
