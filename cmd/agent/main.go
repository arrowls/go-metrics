package main

import (
	"flag"
	"time"

	"github.com/arrowls/go-metrics/internal/collector"
	"github.com/arrowls/go-metrics/internal/updater"
)

const reportIntervalDefault = 2
const pollIntervalDefault = 10
const serverEndpointDefault = "localhost:8080"

type Config struct {
	ReportInterval *int
	PollInterval   *int
	ServerEndpoint *string
}

var config Config

func InitConfig() {
	config.ReportInterval = flag.Int("r", reportIntervalDefault, "report interval in seconds")
	config.PollInterval = flag.Int("p", pollIntervalDefault, "collection interval in seconds")
	config.ServerEndpoint = flag.String("s", serverEndpointDefault, "server endpoint url")

	flag.Parse()
}

func main() {
	InitConfig()
	metricProvider := collector.New()
	metricUpdater := updater.New(metricProvider, *config.ServerEndpoint)

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
			time.Sleep(time.Duration(*config.PollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(*config.ReportInterval) * time.Second)
			consumer.Update()
		}
	}()
}
