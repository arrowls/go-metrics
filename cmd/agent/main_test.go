package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/arrowls/go-metrics/internal/config"
	"github.com/stretchr/testify/assert"
)

type MockMetricsProvider struct {
	CollectCalledTimes int
}

func (m *MockMetricsProvider) Collect() {
	m.CollectCalledTimes++
}

func (m *MockMetricsProvider) AsMap() *map[string]interface{} {
	return &map[string]interface{}{}
}

type MockMetricsConsumer struct {
	UpdateCalledTimes int
}

func (m *MockMetricsConsumer) Update() {
	m.UpdateCalledTimes++
}

func TestRunCollectionAndUpdate(t *testing.T) {
	agentConfig := config.NewAgentConfig()
	var pollInterval = time.Duration(agentConfig.PollInterval) * time.Second
	var reportInterval = time.Duration(agentConfig.ReportInterval) * time.Second

	fmt.Println(pollInterval, reportInterval)

	fmt.Println("Testing RunCollectionAndUpdate() started")

	provider := &MockMetricsProvider{}
	consumer := &MockMetricsConsumer{}

	RunCollectionAndUpdate(provider, pollInterval, consumer, reportInterval)

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, provider.CollectCalledTimes)

	time.Sleep(pollInterval)

	assert.Equal(t, 2, provider.CollectCalledTimes)

	time.Sleep(reportInterval - pollInterval)

	assert.Equal(t, 1, consumer.UpdateCalledTimes)

	time.Sleep(reportInterval)

	assert.Equal(t, 2, consumer.UpdateCalledTimes)
}
