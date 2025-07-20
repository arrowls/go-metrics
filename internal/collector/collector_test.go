package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	collector := New()
	assert.NotNil(t, collector, "New() returned nil")
}

func TestCollector_Collect(t *testing.T) {
	collector := New().(*Collector)

	// Сохраняем начальные значения
	initialPollCount := collector.PollCount

	collector.Collect()

	// Проверяем, что PollCount увеличился
	assert.Equal(t, initialPollCount+1, collector.PollCount, "Expected PollCount to be incremented")

	// Проверяем, что метрики заполнены
	assert.True(t, collector.Alloc > 0 || collector.HeapAlloc > 0, "Expected memory metrics to be collected")

	// Проверяем, что RandomValue в диапазоне [0, 1)
	assert.GreaterOrEqual(t, collector.RandomValue, 0.0, "Expected RandomValue to be >= 0")
	assert.Less(t, collector.RandomValue, 1.0, "Expected RandomValue to be < 1")
}

func TestCollector_AsMap(t *testing.T) {
	collector := New().(*Collector)

	// Собираем метрики
	collector.Collect()

	// Получаем карту метрик
	metricsMap := collector.AsMap()

	assert.NotNil(t, metricsMap, "AsMap() returned nil")

	// Проверяем наличие ключевых метрик
	expectedKeys := []string{"Alloc", "HeapAlloc", "PollCount", "RandomValue"}
	for _, key := range expectedKeys {
		_, exists := (*metricsMap)[key]
		assert.True(t, exists, "Expected key %s in metrics map", key)
	}

	// Проверяем типы значений
	pollCount, ok := (*metricsMap)["PollCount"].(int64)
	assert.True(t, ok, "Expected PollCount to be int64")
	assert.GreaterOrEqual(t, pollCount, int64(0), "Expected PollCount to be non-negative")

	randomValue, ok := (*metricsMap)["RandomValue"].(float64)
	assert.True(t, ok, "Expected RandomValue to be float64")
	assert.GreaterOrEqual(t, randomValue, 0.0, "Expected RandomValue to be >= 0")
	assert.Less(t, randomValue, 1.0, "Expected RandomValue to be < 1")
}

func TestCollector_Concurrency(t *testing.T) {
	collector := New().(*Collector)

	// Запускаем несколько горутин для тестирования конкурентности
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			collector.Collect()
			collector.AsMap()
			done <- true
		}()
	}

	// Ждем завершения всех горутин
	for i := 0; i < 10; i++ {
		<-done
	}

	// Проверяем, что PollCount увеличился на 10
	assert.Equal(t, int64(10), collector.PollCount, "Expected PollCount to be 10")
}
