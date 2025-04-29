package config

import (
	"log"

	"github.com/arrowls/go-metrics/internal/di"
)

const diKey = "server_config"

func ProvideServerConfig(container di.ContainerInterface) ServerConfig {
	config := container.Get(diKey)
	if cfg, ok := config.(ServerConfig); ok {
		return cfg
	}

	cfg := newServerConfig()

	err := container.Add(diKey, cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
