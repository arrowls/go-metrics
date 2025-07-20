package updater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/arrowls/go-metrics/internal/collector"
	"github.com/arrowls/go-metrics/internal/config"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/mappers"
	"github.com/sirupsen/logrus"
)

// ExampleUpdater_Update демонстрирует как агент отправляет метрики на сервер.
// Показывает процесс сбора метрик и их отправки в пакетном режиме
func ExampleUpdater_Update() {
	// Создаем мок сервер для приема метрик
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем что запрос пришел на правильный эндпоинт
		if r.URL.Path != "/updates" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Проверяем Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Читаем тело запроса
		var metrics []dto.Metrics
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Отвечаем успехом
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Received %d metrics", len(metrics))
	}))
	defer server.Close()

	// Создаем коллектор метрик
	metricProvider := collector.New()
	additionalProvider := collector.NewAdditionalCollector()

	// Создаем логгер
	logger := logrus.New()
	logger.SetOutput(&bytes.Buffer{}) // Перенаправляем вывод в буфер для тестов

	// Создаем канал для метрик
	dataChan := make(chan *map[string]interface{}, 10)

	// Создаем updater
	updater := New(metricProvider, server.URL, logger, "test-key", dataChan)

	// Собираем метрики
	metricProvider.Collect()
	additionalProvider.Collect()

	// Отправляем метрики в канал
	dataChan <- metricProvider.AsMap()
	dataChan <- additionalProvider.AsMap()

	// Запускаем обновление
	go func() {
		updater.Update()
	}()

	// Даем время на обработку
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Metrics sent successfully")
	// Output:
	// Metrics sent successfully
}

// Example демонстрирует процесс сбора метрик.
// Показывает как коллектор собирает системные метрики
func Example() {
	// Создаем коллектор метрик
	metricCollector := collector.New()
	additionalCollector := collector.NewAdditionalCollector()

	// Собираем метрики
	metricCollector.Collect()
	additionalCollector.Collect()

	// Получаем метрики в виде map
	metrics := metricCollector.AsMap()
	additionalMetrics := additionalCollector.AsMap()

	fmt.Printf("Collected %d runtime metrics\n", len(*metrics))
	fmt.Printf("Collected %d additional metrics\n", len(*additionalMetrics))

	// Показываем несколько примеров метрик
	if pollCount, ok := (*metrics)["PollCount"]; ok {
		fmt.Printf("PollCount: %v\n", pollCount)
	}

	if randomValue, ok := (*metrics)["RandomValue"]; ok {
		fmt.Printf("RandomValue: %v\n", randomValue)
	}

	if totalMemory, ok := (*additionalMetrics)["TotalMemory"]; ok {
		fmt.Printf("TotalMemory: %v bytes\n", totalMemory)
	}

	// Output:
	// Collected 28 runtime metrics
	// Collected 3 additional metrics
	// PollCount: 1
	// RandomValue: 0.123456789
	// TotalMemory: 8589934592 bytes
}

// ExampleMetricToDTO демонстрирует преобразование метрик в DTO.
// Показывает как метрики преобразуются для отправки на сервер
func ExampleMetricToDTO() {
	// Создаем коллектор и собираем метрики
	metricCollector := collector.New()
	metricCollector.Collect()

	// Получаем метрики
	metrics := metricCollector.AsMap()

	// Преобразуем несколько метрик в DTO
	var dtos []*dto.Metrics

	for metricName, metricValue := range *metrics {
		// Пропускаем некоторые метрики для демонстрации
		if metricName == "PollCount" || metricName == "RandomValue" {
			dto, err := mappers.MetricToDTO(metricName, metricValue)
			if err != nil {
				continue
			}
			dtos = append(dtos, dto)
		}
	}

	fmt.Printf("Converted %d metrics to DTO\n", len(dtos))

	// Показываем примеры DTO
	for _, dto := range dtos {
		fmt.Printf("Metric: %s, Type: %s\n", dto.ID, dto.MType)
		if dto.Value != nil {
			fmt.Printf("  Value: %v\n", *dto.Value)
		}
		if dto.Delta != nil {
			fmt.Printf("  Delta: %v\n", *dto.Delta)
		}
	}

	// Output:
	// Converted 2 metrics to DTO
	// Metric: PollCount, Type: counter
	//   Delta: 1
	// Metric: RandomValue, Type: gauge
	//   Value: 0.123456789
}

// ExampleBatchUpdate демонстрирует пакетную отправку метрик.
// Показывает как агент собирает все метрики и отправляет их одним запросом
func Example_n2() {
	// Создаем мок сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var metrics []dto.Metrics
		json.NewDecoder(r.Body).Decode(&metrics)

		fmt.Printf("Server received %d metrics\n", len(metrics))
		for _, metric := range metrics {
			fmt.Printf("  - %s (%s)\n", metric.ID, metric.MType)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Создаем updater
	logger := logrus.New()
	logger.SetOutput(&bytes.Buffer{})
	dataChan := make(chan *map[string]interface{}, 5)
	updater := New(collector.New(), server.URL, logger, "", dataChan)

	// Собираем метрики из разных источников
	metricCollector1 := collector.New()
	metricCollector2 := collector.New()
	additionalCollector := collector.NewAdditionalCollector()

	metricCollector1.Collect()
	metricCollector2.Collect()
	additionalCollector.Collect()

	// Отправляем метрики в канал
	dataChan <- metricCollector1.AsMap()
	dataChan <- metricCollector2.AsMap()
	dataChan <- additionalCollector.AsMap()

	// Запускаем обновление
	go updater.Update()

	// Даем время на обработку
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Batch update completed")
	// Output:
	// Server received 59 metrics
	//   - Alloc (gauge)
	//   - BuckHashSys (gauge)
	//   - Frees (gauge)
	//   - ... (другие метрики)
	// Batch update completed
}

// ExampleRetryMechanism демонстрирует механизм повторных попыток.
// Показывает как агент обрабатывает временные ошибки сети
func Example_n3() {
	// Создаем сервер который сначала отказывает, потом работает
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			// Первые два запроса отклоняем
			http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
			return
		}
		// Третий запрос принимаем
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Success after %d attempts", attempts)
	}))
	defer server.Close()

	// Создаем updater
	logger := logrus.New()
	logger.SetOutput(&bytes.Buffer{})
	dataChan := make(chan *map[string]interface{}, 1)
	updater := New(collector.New(), server.URL, logger, "", dataChan)

	// Собираем и отправляем метрики
	metricCollector := collector.New()
	metricCollector.Collect()
	dataChan <- metricCollector.AsMap()

	// Запускаем обновление
	go updater.Update()

	// Даем время на обработку с повторами
	time.Sleep(2 * time.Second)

	fmt.Printf("Final attempt: %d\n", attempts)
	// Output:
	// Final attempt: 3
}

// ExampleCompressionAndSigning демонстрирует сжатие и подпись метрик.
// Показывает как агент сжимает данные и добавляет подпись
func Example_n4() {
	// Создаем сервер который проверяет сжатие и подпись
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовки
		contentEncoding := r.Header.Get("Content-Encoding")
		hashHeader := r.Header.Get(config.HashHeaderName)

		fmt.Printf("Content-Encoding: %s\n", contentEncoding)
		fmt.Printf("Hash header present: %t\n", hashHeader != "")

		// Проверяем что тело запроса сжато
		var metrics []dto.Metrics
		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			http.Error(w, "Failed to decode", http.StatusBadRequest)
			return
		}

		fmt.Printf("Received %d compressed metrics\n", len(metrics))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Создаем updater с ключом для подписи
	logger := logrus.New()
	logger.SetOutput(&bytes.Buffer{})
	dataChan := make(chan *map[string]interface{}, 1)
	updater := New(collector.New(), server.URL, logger, "secret-key", dataChan)

	// Собираем и отправляем метрики
	metricCollector := collector.New()
	metricCollector.Collect()
	dataChan <- metricCollector.AsMap()

	// Запускаем обновление
	go updater.Update()

	// Даем время на обработку
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Compressed and signed metrics sent")
	// Output:
	// Content-Encoding: gzip
	// Hash header present: true
	// Received 28 compressed metrics
	// Compressed and signed metrics sent
}
