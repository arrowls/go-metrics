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

	generatorChan := make(chan *map[string]interface{}, agentConfig.RateLimit)
	metricUpdater := updater.New(metricProvider, agentConfig.ServerEndpoint, logger, agentConfig.Key, generatorChan)

	collectionTicker := time.NewTicker(time.Duration(agentConfig.PollInterval) * time.Second)
	updateTicker := time.NewTicker(time.Duration(agentConfig.ReportInterval) * time.Second)

	for {
		select {
		case <-collectionTicker.C:
			go func() {
				metricProvider.Collect()
				generatorChan <- metricProvider.AsMap()
			}()
		case <-updateTicker.C:
			go metricUpdater.Update()
		}
	}
}
