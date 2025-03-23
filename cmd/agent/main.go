package main

import (
	"time"

	"github.com/arrowls/go-metrics/internal/collector"
	"github.com/arrowls/go-metrics/internal/updater"
)

// env или argv
const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	metricProvider := collector.New()
	// env или argv
	metricUpdater := updater.New(metricProvider, "http://localhost:8080")

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
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			time.Sleep(reportInterval)
			consumer.Update()
		}
	}()
}
