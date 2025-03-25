package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrowls/go-metrics/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMetricService struct {
	mock.Mock
}

func (s *MockMetricService) CreateByType(metricType string, name string, stringValue string) error {
	args := s.Called(metricType, name, stringValue)

	return args.Error(0)
}

func (s *MockMetricService) GetList() *map[string]interface{} {
	return &map[string]interface{}{}
}
func (s *MockMetricService) GetItem(metricType string, name string) (string, error) {
	if metricType != "" && name != "" {
		return "123", nil
	}
	return "", errors.New("not found")
}

func createContext(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()

	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestMetricController_HandleNew(t *testing.T) {
	mockService := &MockMetricService{}

	controller := NewMetricController(&service.Service{Metric: mockService})

	t.Run("HandleNew/should create a new metric", func(t *testing.T) {
		mockService.On("CreateByType", "gauge", "TestMetric", "1.23").Return(nil)
		r := httptest.NewRequest(http.MethodPost, "/update/gauge/TestMetric/1.23", nil)
		w := httptest.NewRecorder()

		r = createContext(r, map[string]string{
			"type":  "gauge",
			"name":  "TestMetric",
			"value": "1.23",
		})

		controller.HandleNew(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		mockService.AssertCalled(t, "CreateByType", "gauge", "TestMetric", "1.23")
	})

	t.Run("HandleNew/empty value", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/update/gauge/TestMetric/", nil)
		w := httptest.NewRecorder()

		r = createContext(r, map[string]string{
			"type":  "gauge",
			"name":  "TestMetric",
			"value": "",
		})

		controller.HandleNew(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.Equal(t, "No metric value specified\n", w.Body.String())
	})

	t.Run("HandleNew/empty name", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/update/gauge//1.23", nil)
		w := httptest.NewRecorder()

		r = createContext(r, map[string]string{
			"type":  "gauge",
			"name":  "",
			"value": "1.23",
		})

		controller.HandleNew(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)

		assert.Equal(t, "No metric name specified\n", w.Body.String())
	})

	t.Run("HandleItem/success case", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/value/gauge/TestMetric", nil)
		w := httptest.NewRecorder()

		r = createContext(r, map[string]string{
			"type": "gauge",
			"name": "TestMetric",
		})

		controller.HandleItem(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "123", w.Body.String())
	})

	t.Run("HandleItem/invalid type", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/value//TestMetric", nil)
		w := httptest.NewRecorder()

		r = createContext(r, map[string]string{
			"type": "",
			"name": "TestMetric",
		})

		controller.HandleItem(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "not found")
	})

	t.Run("HandleItem/invalid name", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/value/gauge/", nil)
		w := httptest.NewRecorder()

		r = createContext(r, map[string]string{
			"type": "gauge",
			"name": "",
		})

		controller.HandleItem(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "No metric name specified")
	})
}
