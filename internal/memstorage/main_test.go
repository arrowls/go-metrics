package memstorage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInstance(t *testing.T) {
	// Сбрасываем синглтон для чистого теста
	instance = nil
	once = sync.Once{}

	storage1 := GetInstance()
	assert.NotNil(t, storage1, "GetInstance returned nil")

	storage2 := GetInstance()
	assert.NotNil(t, storage2, "GetInstance returned nil on second call")

	// Проверяем, что это тот же экземпляр
	assert.Equal(t, storage1, storage2, "GetInstance returned different instances")

	// Проверяем инициализацию полей
	assert.NotNil(t, storage1.Gauge, "Gauge map not initialized")
	assert.NotNil(t, storage1.Counter, "Counter map not initialized")
	assert.NotNil(t, storage1.RWMutex, "RWMutex not initialized")
}

func TestMemStorage_Concurrency(t *testing.T) {
	// Сбрасываем синглтон для чистого теста
	instance = nil
	once = sync.Once{}

	storage := GetInstance()

	// Запускаем несколько горутин для тестирования конкурентности
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Записываем в gauge
			storage.Lock()
			storage.Gauge["test_gauge"] = float64(id)
			storage.Unlock()

			// Читаем из gauge
			storage.RLock()
			_ = storage.Gauge["test_gauge"]
			storage.RUnlock()

			// Записываем в counter
			storage.Lock()
			storage.Counter["test_counter"] = int64(id)
			storage.Unlock()

			// Читаем из counter
			storage.RLock()
			_ = storage.Counter["test_counter"]
			storage.RUnlock()
		}(i)
	}

	wg.Wait()

	// Проверяем, что данные записались
	storage.RLock()
	defer storage.RUnlock()

	assert.NotEmpty(t, storage.Gauge, "Expected gauge data to be written")
	assert.NotEmpty(t, storage.Counter, "Expected counter data to be written")
}

func TestMemStorage_DataIntegrity(t *testing.T) {
	// Сбрасываем синглтон для чистого теста
	instance = nil
	once = sync.Once{}

	storage := GetInstance()

	// Тестируем запись и чтение gauge
	storage.Lock()
	storage.Gauge["test_gauge"] = 123.45
	storage.Unlock()

	storage.RLock()
	value, exists := storage.Gauge["test_gauge"]
	storage.RUnlock()

	assert.True(t, exists, "Expected gauge value to exist")
	assert.Equal(t, 123.45, value, "Expected gauge value 123.45")

	// Тестируем запись и чтение counter
	storage.Lock()
	storage.Counter["test_counter"] = 100
	storage.Unlock()

	storage.RLock()
	counterValue, exists := storage.Counter["test_counter"]
	storage.RUnlock()

	assert.True(t, exists, "Expected counter value to exist")
	assert.Equal(t, int64(100), counterValue, "Expected counter value 100")
}

func TestMemStorage_MultipleInstances(t *testing.T) {
	// Сбрасываем синглтон для чистого теста
	instance = nil
	once = sync.Once{}

	// Получаем несколько экземпляров
	storage1 := GetInstance()
	storage2 := GetInstance()
	storage3 := GetInstance()

	// Проверяем, что все ссылки указывают на один объект
	assert.Equal(t, storage1, storage2, "All instances should be the same")
	assert.Equal(t, storage2, storage3, "All instances should be the same")

	// Записываем данные в один экземпляр
	storage1.Lock()
	storage1.Gauge["shared_gauge"] = 999.99
	storage1.Counter["shared_counter"] = 999
	storage1.Unlock()

	// Проверяем, что данные доступны через другой экземпляр
	storage2.RLock()
	gaugeValue, exists := storage2.Gauge["shared_gauge"]
	storage2.RUnlock()

	assert.True(t, exists, "Expected shared gauge value to exist")
	assert.Equal(t, 999.99, gaugeValue, "Expected shared gauge value 999.99")

	storage3.RLock()
	counterValue, exists := storage3.Counter["shared_counter"]
	storage3.RUnlock()

	assert.True(t, exists, "Expected shared counter value to exist")
	assert.Equal(t, int64(999), counterValue, "Expected shared counter value 999")
}
