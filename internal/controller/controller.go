// Package controller содержит HTTP контроллеры для обработки запросов.
// Предоставляет интерфейсы и реализации для работы с метриками и публичными эндпоинтами.
package controller

import (
	"net/http"
)

// Metric определяет интерфейс для HTTP обработчиков метрик.
// Включает методы для создания, получения и пакетной обработки метрик.
type Metric interface {
	// HandleNew обрабатывает POST запросы для создания метрик через URL параметры.
	// Поддерживает формат /update/{type}/{name}/{value}
	HandleNew(rw http.ResponseWriter, r *http.Request)

	// HandleItem обрабатывает GET запросы для получения метрик через URL параметры.
	// Поддерживает формат /value/{type}/{name}
	HandleItem(rw http.ResponseWriter, r *http.Request)

	// HandleNewFromBody обрабатывает POST запросы с JSON телом для создания метрик.
	// Принимает JSON объект с полями type, name, value.
	HandleNewFromBody(rw http.ResponseWriter, r *http.Request)

	// HandleGetItemFromBody обрабатывает POST запросы с JSON телом для получения метрик.
	// Принимает JSON объект с полями type, name и возвращает значение в JSON.
	HandleGetItemFromBody(rw http.ResponseWriter, r *http.Request)

	// HandleCreateBatch обрабатывает POST запросы для пакетного создания метрик.
	// Принимает массив JSON объектов с метриками.
	HandleCreateBatch(rw http.ResponseWriter, r *http.Request)
}

// Public определяет интерфейс для публичных HTTP обработчиков.
// Включает методы для отображения информации и проверки состояния.
type Public interface {
	// HandlePublic обрабатывает запросы к публичным эндпоинтам.
	// Возвращает общую информацию о сервисе.
	HandlePublic(rw http.ResponseWriter, r *http.Request)

	// HandleIndex обрабатывает запросы к корневому эндпоинту.
	// Возвращает HTML страницу с информацией о доступных метриках.
	HandleIndex(rw http.ResponseWriter, r *http.Request)

	// Ping обрабатывает запросы для проверки доступности сервиса.
	// Возвращает статус соединения с хранилищем метрик.
	Ping(rw http.ResponseWriter, r *http.Request)
}

// ErrorHandler определяет интерфейс для обработки ошибок HTTP запросов.
// Предоставляет единообразный способ возврата ошибок клиенту.
type ErrorHandler interface {
	// Handle обрабатывает ошибки и возвращает соответствующий HTTP ответ.
	// Принимает ResponseWriter и ошибку для формирования ответа.
	Handle(w http.ResponseWriter, err error)
}

// Controller содержит все HTTP контроллеры приложения.
// Предоставляет централизованный доступ к обработчикам запросов.
type Controller struct {
	// Metric - контроллер для работы с метриками
	Metric Metric
	// Public - контроллер для публичных эндпоинтов
	Public Public
}
