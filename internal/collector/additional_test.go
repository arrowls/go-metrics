package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAdditionalCollector(t *testing.T) {
	collector := NewAdditionalCollector()
	assert.NotNil(t, collector, "NewAdditionalCollector() returned nil")
}

func TestAdditionalCollector_Collect(t *testing.T) {
	collector := NewAdditionalCollector().(*AdditionalCollector)

	collector.Collect()

	// Проверяем, что метрики памяти заполнены
	assert.Greater(t, collector.TotalMemory, float64(0), "Expected TotalMemory to be collected")
	assert.GreaterOrEqual(t, collector.FreeMemory, float64(0), "Expected FreeMemory to be collected")

	// Проверяем, что FreeMemory не больше TotalMemory
	assert.LessOrEqual(t, collector.FreeMemory, collector.TotalMemory,
		"FreeMemory should not be greater than TotalMemory")

	// Проверяем CPU utilization
	assert.Greater(t, collector.CPUutilization1, float64(0), "Expected CPUutilization1 to be positive")
}

func TestAdditionalCollector_AsMap(t *testing.T) {
	collector := NewAdditionalCollector().(*AdditionalCollector)

	// Собираем метрики
	collector.Collect()

	// Получаем карту метрик
	metricsMap := collector.AsMap()

	assert.NotNil(t, metricsMap, "AsMap() returned nil")

	// Проверяем наличие всех ключей
	expectedKeys := []string{"TotalMemory", "FreeMemory", "CPUutilization1"}
	for _, key := range expectedKeys {
		_, exists := (*metricsMap)[key]
		assert.True(t, exists, "Expected key %s in metrics map", key)
	}

	// Проверяем типы значений
	totalMemory, ok := (*metricsMap)["TotalMemory"].(float64)
	assert.True(t, ok, "Expected TotalMemory to be float64")
	assert.Greater(t, totalMemory, float64(0), "Expected TotalMemory to be positive")

	freeMemory, ok := (*metricsMap)["FreeMemory"].(float64)
	assert.True(t, ok, "Expected FreeMemory to be float64")
	assert.GreaterOrEqual(t, freeMemory, float64(0), "Expected FreeMemory to be non-negative")

	cpuUtil, ok := (*metricsMap)["CPUutilization1"].(float64)
	assert.True(t, ok, "Expected CPUutilization1 to be float64")
	assert.Greater(t, cpuUtil, float64(0), "Expected CPUutilization1 to be positive")
}

func TestAdditionalCollector_Concurrency(t *testing.T) {
	collector := NewAdditionalCollector().(*AdditionalCollector)

	// Запускаем несколько горутин для тестирования конкурентности
	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func() {
			collector.Collect()
			collector.AsMap()
			done <- true
		}()
	}

	// Ждем завершения всех горутин
	for i := 0; i < 5; i++ {
		<-done
	}

	// Проверяем, что метрики все еще валидны
	assert.Greater(t, collector.TotalMemory, float64(0), "Expected TotalMemory to remain valid after concurrent access")
	assert.GreaterOrEqual(t, collector.FreeMemory, float64(0), "Expected FreeMemory to remain valid after concurrent access")
	assert.Greater(t, collector.CPUutilization1, float64(0), "Expected CPUutilization1 to remain valid after concurrent access")
}
