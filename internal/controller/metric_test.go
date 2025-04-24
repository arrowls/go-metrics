package controller

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMetricService struct {
	mock.Mock
}

func (s *MockMetricService) Create(dto *dto.CreateMetric) error {
	args := s.Called(dto.Type, dto.Name, dto.Value)

	return args.Error(0)
}

func (s *MockMetricService) GetList() *map[string]interface{} {
	return &map[string]interface{}{}
}
func (s *MockMetricService) GetItem(dto *dto.GetMetric) (string, error) {
	if dto.Type != "" && dto.Name != "" {
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
	mockLogger := logrus.New()
	mockLogger.SetOutput(io.Discard)
	errorHandler := apperrors.NewHTTPErrorHandler(mockLogger)

	controller := NewMetricController(&service.Service{Metric: mockService}, errorHandler)

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
			mockSetup:    func() { mockService.On("Create", "gauge", "TestMetric", "1.23").Return(nil) },
			expectedCode: http.StatusOK,
		},
		{
			name:         "empty value",
			method:       http.MethodPost,
			url:          "/update/gauge/TestMetric/",
			params:       map[string]string{"type": "gauge", "name": "TestMetric", "value": ""},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Failed to read request: metric value not specified"}` + "\n",
		},
		{
			name:         "empty name",
			method:       http.MethodPost,
			url:          "/update/gauge//1.23",
			params:       map[string]string{"type": "gauge", "name": "", "value": "1.23"},
			expectedCode: http.StatusNotFound,
			expectedBody: `{"message":"Failed to read request: metric name not specified"}` + "\n",
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
	mockLogger := logrus.New()
	mockLogger.SetOutput(io.Discard)
	errorHandler := apperrors.NewHTTPErrorHandler(mockLogger)

	controller := NewMetricController(&service.Service{Metric: mockService}, errorHandler)

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
			expectedBody: `{"message":"Failed to read request: unknown metric type:"}` + "\n",
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
			expectedBody: `{"message":"Failed to read request: metric name is not specified"}` + "\n",
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

func TestMetricController_HandleNewFromBody(t *testing.T) {
	mockService := &MockMetricService{}
	mockLogger := logrus.New()
	mockLogger.SetOutput(io.Discard)
	errorHandler := apperrors.NewHTTPErrorHandler(mockLogger)

	controller := NewMetricController(&service.Service{Metric: mockService}, errorHandler)

	tests := []struct {
		name               string
		expectedStatusCode int
		body               []byte
		mockSetup          func()
	}{
		{
			name:               "success case",
			expectedStatusCode: 200,
			body: []byte(`{
				"type":"gauge",
				"id":"TestMetric",
				"value":1
			}`),
			mockSetup: func() {
				mockService.On("Create", "gauge", "TestMetric", "1").Return(nil)
			},
		},
		{
			name:               "error in mapper",
			body:               nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "error in service",
			body: []byte(`{
				"type":"gauge",
				"id":"TestMetric",
				"value":2
			}`),
			expectedStatusCode: http.StatusBadRequest,
			mockSetup: func() {
				mockService.On("Create", "gauge", "TestMetric", "2").Return(apperrors.ErrBadRequest)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			body := bytes.NewReader(tt.body)
			r := httptest.NewRequest("POST", "/update", body)
			w := httptest.NewRecorder()

			controller.HandleNewFromBody(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			bodyBytes, err := io.ReadAll(w.Body)
			assert.Nil(t, err)
			if tt.expectedStatusCode == http.StatusOK {
				assert.Contains(t, string(bodyBytes), `"value":123`)
			}
		})
	}
}

func TestMetricController_HandleGetItemFromBody(t *testing.T) {
	mockService := &MockMetricService{}
	mockLogger := logrus.New()
	mockLogger.SetOutput(io.Discard)
	errorHandler := apperrors.NewHTTPErrorHandler(mockLogger)

	controller := NewMetricController(&service.Service{Metric: mockService}, errorHandler)

	tests := []struct {
		name         string
		mockSetup    func()
		expectedCode int
		expectedBody string
		body         []byte
	}{
		{
			name:         "success case",
			expectedBody: `{"id":"TestMetric","type":"gauge","value":123}`,
			expectedCode: http.StatusOK,
			body: []byte(`{
				"type":"gauge",
				"id":"TestMetric"
			}`),
		},
		{
			name:         "error in mapper",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Error reading request: could not read the request body"}` + "\n",
			body: []byte(`{
				"type":"invalid_type",
				"id":"TestMetric",
			}`),
		},
		{
			name:         "error in service",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Error reading request: could not read the request body"}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewReader(tt.body)
			r := httptest.NewRequest(http.MethodPost, "/value", body)
			w := httptest.NewRecorder()

			controller.HandleGetItemFromBody(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
