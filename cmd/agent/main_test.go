package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/arrowls/go-metrics/internal/config"
	"github.com/stretchr/testify/assert"
)

type MockMetricsProvider struct {
	CollectCalledTimes int
	*sync.Mutex
}

func (m *MockMetricsProvider) Get() int {
	m.Lock()
	defer m.Unlock()
	return m.CollectCalledTimes
}

func (m *MockMetricsProvider) Collect() {
	m.Lock()
	m.CollectCalledTimes++
	m.Unlock()
}

func (m *MockMetricsProvider) AsMap() *map[string]interface{} {
	return &map[string]interface{}{}
}

type MockMetricsConsumer struct {
	UpdateCalledTimes int
	*sync.Mutex
}

func (m *MockMetricsConsumer) Get() int {
	m.Lock()
	defer m.Unlock()
	return m.UpdateCalledTimes
}
func (m *MockMetricsConsumer) Update() {
	m.Lock()
	m.UpdateCalledTimes++
	m.Unlock()
}

func TestRunCollectionAndUpdate(t *testing.T) {
	agentConfig := config.NewAgentConfig()
	var pollInterval = time.Duration(agentConfig.PollInterval) * time.Second
	var reportInterval = time.Duration(agentConfig.ReportInterval) * time.Second

	fmt.Println("Testing RunCollectionAndUpdate() started")

	provider := &MockMetricsProvider{
		0,
		&sync.Mutex{},
	}
	consumer := &MockMetricsConsumer{
		0,
		&sync.Mutex{},
	}

	RunCollectionAndUpdate(provider, pollInterval, consumer, reportInterval)

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, provider.Get())

	time.Sleep(pollInterval)

	assert.Equal(t, 2, provider.Get())

	time.Sleep(reportInterval - pollInterval)

	assert.Equal(t, 1, consumer.Get())

	time.Sleep(reportInterval)

	assert.Equal(t, 2, consumer.Get())
}
