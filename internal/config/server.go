package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v7"
)

const (
	reportIntervalDefault  = 10
	pollIntervalDefault    = 2
	serverEndpointDefault  = "localhost:8080"
	storeIntervalDefault   = 300
	storageFilePathDefault = "metrics.json"
	restoreDefault         = false
)

type ServerConfig struct {
	ServerEndpoint  string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	StorageFilePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

var serverConfig ServerConfig
var agentConfig AgentConfig

func NewServerConfig() ServerConfig {
	flag.StringVar(&serverConfig.ServerEndpoint, "a", serverEndpointDefault, "server endpoint url")
	flag.IntVar(&serverConfig.StoreInterval, "i", storeIntervalDefault, "interval to write metrics to file")
	flag.StringVar(&serverConfig.StorageFilePath, "f", storageFilePathDefault, "file to write metrics backup")
	flag.BoolVar(&serverConfig.Restore, "r", restoreDefault, "restore on startup")

	flag.Parse()

	if err := env.Parse(&serverConfig); err != nil {
		fmt.Printf("Failed to parse env: %v\n", err)
		os.Exit(1)
	}

	return serverConfig
}
