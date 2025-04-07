package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v7"
)

const (
	reportIntervalDefault = 10
	pollIntervalDefault   = 2
	serverEndpointDefault = "localhost:8080"
)

type ServerConfig struct {
	ServerEndpoint string `env:"ADDRESS"`
}

type AgentConfig struct {
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ServerEndpoint string `env:"ADDRESS"`
}

var serverConfig ServerConfig
var agentConfig AgentConfig

func NewServerConfig() ServerConfig {
	flag.StringVar(&serverConfig.ServerEndpoint, "a", serverEndpointDefault, "server endpoint url")

	flag.Parse()

	if err := env.Parse(&serverConfig); err != nil {
		fmt.Printf("Failed to parse env: %v\n", err)
		os.Exit(1)
	}

	return serverConfig
}

func NewAgentConfig() AgentConfig {
	flag.IntVar(&agentConfig.ReportInterval, "r", reportIntervalDefault, "report interval in seconds")
	flag.IntVar(&agentConfig.PollInterval, "p", pollIntervalDefault, "collection interval in seconds")
	flag.StringVar(&agentConfig.ServerEndpoint, "a", serverEndpointDefault, "server endpoint url")

	flag.Parse()

	if err := env.Parse(&agentConfig); err != nil {
		fmt.Printf("Failed to parse env: %v\n", err)
		os.Exit(1)
	}

	return agentConfig
}
