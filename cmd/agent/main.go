package main

import (
	"time"

	"github.com/arrowls/go-metrics/internal/collector"
	"github.com/arrowls/go-metrics/internal/config"
	"github.com/arrowls/go-metrics/internal/updater"
	"github.com/sirupsen/logrus"
)

func main() {
	agentConfig := config.NewAgentConfig()
	metricProvider := collector.New()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	metricUpdater := updater.New(metricProvider, agentConfig.ServerEndpoint, logger, agentConfig.Key)

	stopChan := make(chan struct{})

	RunCollectionAndUpdate(
		metricProvider,
		time.Duration(agentConfig.PollInterval)*time.Second,
		metricUpdater,
		time.Duration(agentConfig.ReportInterval)*time.Second,
	)

	<-stopChan
}

func RunCollectionAndUpdate(
	provider collector.MetricProvider,
	collectInterval time.Duration,
	consumer updater.MetricConsumer,
	updateInterval time.Duration,
) {
	go func() {
		for {
			provider.Collect()
			time.Sleep(collectInterval)
		}
	}()

	go func() {
		for {
			time.Sleep(updateInterval)
			consumer.Update()
		}
	}()
}
