package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/arrowls/go-metrics/internal/collector"
	"github.com/arrowls/go-metrics/internal/updater"
	"github.com/caarlos0/env/v7"
)

const (
	reportIntervalDefault = 10
	pollIntervalDefault   = 2
	serverEndpointDefault = "localhost:8080"
)

type Config struct {
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ServerEndpoint string `env:"ADDRESS"`
}

var config Config

func InitConfig() {
	flag.IntVar(&config.ReportInterval, "r", reportIntervalDefault, "report interval in seconds")
	flag.IntVar(&config.PollInterval, "p", pollIntervalDefault, "collection interval in seconds")
	flag.StringVar(&config.ServerEndpoint, "a", serverEndpointDefault, "server endpoint url")

	flag.Parse()

	if err := env.Parse(&config); err != nil {
		fmt.Printf("Failed to parse env: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	InitConfig()
	metricProvider := collector.New()
	metricUpdater := updater.New(metricProvider, config.ServerEndpoint)

	stopChan := make(chan struct{})

	RunCollectionAndUpdate(
		metricProvider,
		metricUpdater,
	)

	<-stopChan
}

func RunCollectionAndUpdate(
	provider collector.MetricProvider,
	consumer updater.MetricConsumer,
) {
	go func() {
		for {
			provider.Collect()
			time.Sleep(time.Duration(config.PollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(config.ReportInterval) * time.Second)
			consumer.Update()
		}
	}()
}
