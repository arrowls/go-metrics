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

	tests := []struct {
		name         string
		method       string
		url          string
		params       map[string]string
		mockSetup    func()
		expectedCode int
		expectedBody string
	}{
		{
			name:         "should create a new metric",
			method:       http.MethodPost,
			url:          "/update/gauge/TestMetric/1.23",
			params:       map[string]string{"type": "gauge", "name": "TestMetric", "value": "1.23"},
			mockSetup:    func() { mockService.On("CreateByType", "gauge", "TestMetric", "1.23").Return(nil) },
			expectedCode: http.StatusOK,
		},
		{
			name:         "empty value",
			method:       http.MethodPost,
			url:          "/update/gauge/TestMetric/",
			params:       map[string]string{"type": "gauge", "name": "TestMetric", "value": ""},
			expectedCode: http.StatusBadRequest,
			expectedBody: "No metric value specified\n",
		},
		{
			name:         "empty name",
			method:       http.MethodPost,
			url:          "/update/gauge//1.23",
			params:       map[string]string{"type": "gauge", "name": "", "value": "1.23"},
			expectedCode: http.StatusNotFound,
			expectedBody: "No metric name specified\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			r := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			r = createContext(r, tt.params)

			if tt.method == http.MethodPost {
				controller.HandleNew(w, r)
			} else {
				controller.HandleItem(w, r)
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}

			if tt.mockSetup != nil {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestMetricController_HandleItem(t *testing.T) {
	mockService := &MockMetricService{}

	controller := NewMetricController(&service.Service{Metric: mockService})

	tests := []struct {
		name         string
		method       string
		url          string
		params       map[string]string
		mockSetup    func()
		expectedCode int
		expectedBody string
	}{
		{
			name:   "HandleItem/success case",
			method: http.MethodGet,
			url:    "/value/gauge/TestMetric",
			params: map[string]string{
				"type": "gauge",
				"name": "TestMetric",
			},
			expectedBody: "123",
			expectedCode: http.StatusOK,
		},
		{
			name:   "HandleItem/invalid type",
			method: http.MethodGet,
			url:    "/value//TestMetric",
			params: map[string]string{
				"type": "",
				"name": "TestMetric",
			},
			expectedCode: http.StatusNotFound,
			expectedBody: "not found\n",
		},
		{
			name:   "HandleItem/invalid name",
			method: http.MethodGet,
			url:    "/value/gauge/",
			params: map[string]string{
				"type": "gauge",
				"name": "",
			},
			expectedCode: http.StatusNotFound,
			expectedBody: "No metric name specified\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			r = createContext(r, tt.params)

			controller.HandleItem(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
