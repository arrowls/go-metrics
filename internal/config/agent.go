package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v7"
)

type AgentConfig struct {
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ServerEndpoint string `env:"ADDRESS"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func NewAgentConfig() AgentConfig {
	flag.IntVar(&agentConfig.ReportInterval, "r", reportIntervalDefault, "report interval in seconds")
	flag.IntVar(&agentConfig.PollInterval, "p", pollIntervalDefault, "collection interval in seconds")
	flag.StringVar(&agentConfig.ServerEndpoint, "a", serverEndpointDefault, "server endpoint url")
	flag.StringVar(&serverConfig.Key, "k", "", "encoding key")
	flag.IntVar(&agentConfig.RateLimit, "l", 0, "max requests to server")

	flag.Parse()

	if err := env.Parse(&agentConfig); err != nil {
		fmt.Printf("Failed to parse env: %v\n", err)
		os.Exit(1)
	}

	return agentConfig
}
