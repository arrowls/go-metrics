package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/config"
	"github.com/arrowls/go-metrics/internal/di"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/logger"
	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/arrowls/go-metrics/internal/repository"
	"github.com/arrowls/go-metrics/internal/service"
)

// ExampleMetricController_HandleNew демонстрирует создание метрики через URL параметры.
// Показывает как отправить POST запрос на /update/{type}/{name}/{value}
func ExampleMetricController_HandleNew() {
	// Создаем тестовый сервер
	container := di.NewContainer()
	loggerInst := logger.ProvideLogger(container)
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	metricService := service.NewMetricService(repo)
	errorHandler := apperrors.NewHTTPErrorHandler(loggerInst)
	controller := NewMetricController(&service.Service{Metric: metricService}, errorHandler)

	// Создаем HTTP запрос для создания gauge метрики
	req := httptest.NewRequest("POST", "/update/gauge/testMetric/123.45", nil)
	w := httptest.NewRecorder()

	// Выполняем запрос
	controller.HandleNew(w, req)

	fmt.Printf("Status Code: %d\n", w.Code)
	fmt.Printf("Response: %s\n", w.Body.String())
	// Output:
	// Status Code: 200
	// Response:
}

// ExampleMetricController_HandleItem демонстрирует получение метрики через URL параметры.
// Показывает как отправить GET запрос на /value/{type}/{name}
func ExampleMetricController_HandleItem() {
	// Создаем тестовый сервер
	container := di.NewContainer()
	loggerInst := logger.ProvideLogger(container)
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	metricService := service.NewMetricService(repo)
	errorHandler := apperrors.NewHTTPErrorHandler(loggerInst)
	controller := NewMetricController(&service.Service{Metric: metricService}, errorHandler)

	// Сначала создаем метрику
	createReq := httptest.NewRequest("POST", "/update/gauge/testMetric/123.45", nil)
	createW := httptest.NewRecorder()
	controller.HandleNew(createW, createReq)

	// Теперь получаем значение метрики
	req := httptest.NewRequest("GET", "/value/gauge/testMetric", nil)
	w := httptest.NewRecorder()

	controller.HandleItem(w, req)

	fmt.Printf("Status Code: %d\n", w.Code)
	fmt.Printf("Metric Value: %s\n", w.Body.String())
	// Output:
	// Status Code: 200
	// Metric Value: 123.45
}

// ExampleMetricController_HandleNewFromBody демонстрирует создание метрики через JSON.
// Показывает как отправить POST запрос на /update с JSON телом
func ExampleMetricController_HandleNewFromBody() {
	// Создаем тестовый сервер
	container := di.NewContainer()
	loggerInst := logger.ProvideLogger(container)
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	metricService := service.NewMetricService(repo)
	errorHandler := apperrors.NewHTTPErrorHandler(loggerInst)
	controller := NewMetricController(&service.Service{Metric: metricService}, errorHandler)

	// Создаем JSON данные для gauge метрики
	metricData := dto.Metrics{
		ID:    "jsonGaugeMetric",
		MType: config.GaugeType,
		Value: func() *float64 { v := 456.78; return &v }(),
	}

	jsonData, _ := json.Marshal(metricData)
	req := httptest.NewRequest("POST", "/update", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	controller.HandleNewFromBody(w, req)

	fmt.Printf("Status Code: %d\n", w.Code)
	fmt.Printf("Response: %s\n", w.Body.String())
	// Output:
	// Status Code: 200
	// Response: {"id":"jsonGaugeMetric","type":"gauge","value":456.78}
}

// ExampleMetricController_HandleGetItemFromBody демонстрирует получение метрики через JSON.
// Показывает как отправить POST запрос на /value с JSON телом
func ExampleMetricController_HandleGetItemFromBody() {
	// Создаем тестовый сервер
	container := di.NewContainer()
	loggerInst := logger.ProvideLogger(container)
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	metricService := service.NewMetricService(repo)
	errorHandler := apperrors.NewHTTPErrorHandler(loggerInst)
	controller := NewMetricController(&service.Service{Metric: metricService}, errorHandler)

	// Сначала создаем метрику
	createData := dto.Metrics{
		ID:    "counterMetric",
		MType: config.CounterType,
		Delta: func() *int64 { v := int64(100); return &v }(),
	}
	createJSON, _ := json.Marshal(createData)
	createReq := httptest.NewRequest("POST", "/update", bytes.NewReader(createJSON))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	controller.HandleNewFromBody(createW, createReq)

	// Теперь получаем значение метрики через JSON
	getData := dto.Metrics{
		ID:    "counterMetric",
		MType: config.CounterType,
	}
	getJSON, _ := json.Marshal(getData)
	req := httptest.NewRequest("POST", "/value", bytes.NewReader(getJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	controller.HandleGetItemFromBody(w, req)

	fmt.Printf("Status Code: %d\n", w.Code)
	fmt.Printf("Response: %s\n", w.Body.String())
	// Output:
	// Status Code: 200
	// Response: {"id":"counterMetric","type":"counter","delta":100}
}

// ExampleMetricController_HandleCreateBatch демонстрирует пакетное создание метрик.
// Показывает как отправить POST запрос на /updates с массивом JSON метрик
func ExampleMetricController_HandleCreateBatch() {
	// Создаем тестовый сервер
	container := di.NewContainer()
	loggerInst := logger.ProvideLogger(container)
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	metricService := service.NewMetricService(repo)
	errorHandler := apperrors.NewHTTPErrorHandler(loggerInst)
	controller := NewMetricController(&service.Service{Metric: metricService}, errorHandler)

	// Создаем массив метрик для пакетного создания
	metrics := []dto.Metrics{
		{
			ID:    "batchGauge1",
			MType: config.GaugeType,
			Value: func() *float64 { v := 111.11; return &v }(),
		},
		{
			ID:    "batchCounter1",
			MType: config.CounterType,
			Delta: func() *int64 { v := int64(222); return &v }(),
		},
		{
			ID:    "batchGauge2",
			MType: config.GaugeType,
			Value: func() *float64 { v := 333.33; return &v }(),
		},
	}

	jsonData, _ := json.Marshal(metrics)
	req := httptest.NewRequest("POST", "/updates", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	controller.HandleCreateBatch(w, req)

	fmt.Printf("Status Code: %d\n", w.Code)
	fmt.Printf("Batch created successfully\n")
	// Output:
	// Status Code: 200
	// Batch created successfully
}

// ExamplePublicController_Ping демонстрирует проверку состояния сервиса.
// Показывает как отправить GET запрос на /ping
func ExamplePublicController_Ping() {
	// Создаем тестовый сервер
	container := di.NewContainer()
	_ = logger.ProvideLogger(container)
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	metricService := service.NewMetricService(repo)
	controller := NewPublicController(&service.Service{Metric: metricService})

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	controller.Ping(w, req)

	fmt.Printf("Status Code: %d\n", w.Code)
	if w.Code == http.StatusOK {
		fmt.Println("Service is healthy")
	} else {
		fmt.Println("Service is not responding")
	}
	// Output:
	// Status Code: 200
	// Service is healthy
}

// ExampleHTTPClient_CompleteWorkflow демонстрирует полный рабочий процесс с HTTP клиентом.
// Показывает как использовать все основные эндпоинты в реальном приложении
func Example_n1() {
	// В реальном приложении вы бы использовали http.Client
	// Здесь показан пример с httptest для демонстрации

	// Создаем тестовый сервер
	container := di.NewContainer()
	loggerInst := logger.ProvideLogger(container)
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	metricService := service.NewMetricService(repo)
	errorHandler := apperrors.NewHTTPErrorHandler(loggerInst)
	metricController := NewMetricController(&service.Service{Metric: metricService}, errorHandler)
	publicController := NewPublicController(&service.Service{Metric: metricService})

	// 1. Проверяем состояние сервиса
	pingReq := httptest.NewRequest("GET", "/ping", nil)
	pingW := httptest.NewRecorder()
	publicController.Ping(pingW, pingReq)
	fmt.Printf("Ping Status: %d\n", pingW.Code)

	// 2. Создаем gauge метрику через URL
	gaugeReq := httptest.NewRequest("POST", "/update/gauge/cpu_usage/75.5", nil)
	gaugeW := httptest.NewRecorder()
	metricController.HandleNew(gaugeW, gaugeReq)
	fmt.Printf("Create Gauge Status: %d\n", gaugeW.Code)

	// 3. Создаем counter метрику через JSON
	counterData := dto.Metrics{
		ID:    "request_count",
		MType: config.CounterType,
		Delta: func() *int64 { v := int64(1); return &v }(),
	}
	counterJSON, _ := json.Marshal(counterData)
	counterReq := httptest.NewRequest("POST", "/update", bytes.NewReader(counterJSON))
	counterReq.Header.Set("Content-Type", "application/json")
	counterW := httptest.NewRecorder()
	metricController.HandleNewFromBody(counterW, counterReq)
	fmt.Printf("Create Counter Status: %d\n", counterW.Code)

	// 4. Получаем gauge метрику через URL
	getGaugeReq := httptest.NewRequest("GET", "/value/gauge/cpu_usage", nil)
	getGaugeW := httptest.NewRecorder()
	metricController.HandleItem(getGaugeW, getGaugeReq)
	fmt.Printf("Get Gauge Value: %s\n", strings.TrimSpace(getGaugeW.Body.String()))

	// 5. Получаем counter метрику через JSON
	getCounterData := dto.Metrics{
		ID:    "request_count",
		MType: config.CounterType,
	}
	getCounterJSON, _ := json.Marshal(getCounterData)
	getCounterReq := httptest.NewRequest("POST", "/value", bytes.NewReader(getCounterJSON))
	getCounterReq.Header.Set("Content-Type", "application/json")
	getCounterW := httptest.NewRecorder()
	metricController.HandleGetItemFromBody(getCounterW, getCounterReq)
	fmt.Printf("Get Counter Response: %s\n", strings.TrimSpace(getCounterW.Body.String()))

	// 6. Создаем пакет метрик
	batchData := []dto.Metrics{
		{
			ID:    "memory_usage",
			MType: config.GaugeType,
			Value: func() *float64 { v := 512.0; return &v }(),
		},
		{
			ID:    "error_count",
			MType: config.CounterType,
			Delta: func() *int64 { v := int64(5); return &v }(),
		},
	}
	batchJSON, _ := json.Marshal(batchData)
	batchReq := httptest.NewRequest("POST", "/updates", bytes.NewReader(batchJSON))
	batchReq.Header.Set("Content-Type", "application/json")
	batchW := httptest.NewRecorder()
	metricController.HandleCreateBatch(batchW, batchReq)
	fmt.Printf("Batch Create Status: %d\n", batchW.Code)

	// Output:
	// Ping Status: 200
	// Create Gauge Status: 200
	// Create Counter Status: 200
	// Get Gauge Value: 75.5
	// Get Counter Response: {"id":"request_count","type":"counter","delta":1}
	// Batch Create Status: 200
}

// ExampleErrorHandling демонстрирует обработку ошибок при работе с API.
// Показывает различные сценарии ошибок и их обработку
func Example_n2() {
	// Создаем тестовый сервер
	container := di.NewContainer()
	loggerInst := logger.ProvideLogger(container)
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	metricService := service.NewMetricService(repo)
	errorHandler := apperrors.NewHTTPErrorHandler(loggerInst)
	controller := NewMetricController(&service.Service{Metric: metricService}, errorHandler)

	// 1. Попытка получить несуществующую метрику
	req1 := httptest.NewRequest("GET", "/value/gauge/nonexistent", nil)
	w1 := httptest.NewRecorder()
	controller.HandleItem(w1, req1)
	fmt.Printf("Non-existent metric - Status: %d, Response: %s\n", w1.Code, strings.TrimSpace(w1.Body.String()))

	// 2. Попытка создать метрику с некорректным типом
	req2 := httptest.NewRequest("POST", "/update/invalid/test/123", nil)
	w2 := httptest.NewRecorder()
	controller.HandleNew(w2, req2)
	fmt.Printf("Invalid type - Status: %d, Response: %s\n", w2.Code, strings.TrimSpace(w2.Body.String()))

	// 3. Попытка отправить некорректный JSON
	invalidJSON := []byte(`{"invalid": "json"`)
	req3 := httptest.NewRequest("POST", "/update", bytes.NewReader(invalidJSON))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	controller.HandleNewFromBody(w3, req3)
	fmt.Printf("Invalid JSON - Status: %d, Response: %s\n", w3.Code, strings.TrimSpace(w3.Body.String()))

	// Output:
	// Non-existent metric - Status: 404, Response: {"message":"Failed to get metric: not foundgauge nonexistent not found"}
	// Invalid type - Status: 404, Response: {"message":"Failed to create metric: not foundunknown metric type: invalid"}
	// Invalid JSON - Status: 400, Response: {"message":"Failed to read request: bad requestfailed to read the body of the request: unexpected end of JSON input"}
}
