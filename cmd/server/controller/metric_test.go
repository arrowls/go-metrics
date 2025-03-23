package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrowls/go-metrics/cmd/server/service"
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

func createContext(r *http.Request, typeValue, nameValue, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", typeValue)
	rctx.URLParams.Add("name", nameValue)
	rctx.URLParams.Add("value", value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestMetricController_HandleNew(t *testing.T) {
	mockService := &MockMetricService{}

	controller := MetricController{service.Service{
		Metric: mockService,
	}}

	t.Run("should create a new metric", func(t *testing.T) {
		mockService.On("CreateByType", "gauge", "TestMetric", "1.23").Return(nil)
		r := httptest.NewRequest("POST", "/update/gauge/TestMetric/1.23", nil)
		w := httptest.NewRecorder()

		r = createContext(r, "gauge", "TestMetric", "1.23")

		controller.HandleNew(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		mockService.AssertCalled(t, "CreateByType", "gauge", "TestMetric", "1.23")
	})

	t.Run("empty value", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/update/gauge/TestMetric/", nil)
		w := httptest.NewRecorder()

		r = createContext(r, "gauge", "TestMetric", "")

		controller.HandleNew(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.Equal(t, "No metric value specified\n", w.Body.String())
	})

	t.Run("empty name", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/update/gauge//1.23", nil)
		w := httptest.NewRecorder()

		r = createContext(r, "gauge", "", "1.23")

		controller.HandleNew(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)

		assert.Equal(t, "No metric name specified\n", w.Body.String())
	})
}
