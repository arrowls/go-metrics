package updater

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/stretchr/testify/assert"
)

type MockProvider struct{}

func (p MockProvider) AsMap() *map[string]interface{} {
	return &map[string]interface{}{
		"MetricName1": int64(1),
		"MetricName2": float64(2.22),
		"MetricName3": int64(3),
		"MetricName4": float64(4.4),
	}
}

func (p MockProvider) Collect() {}

func TestUpdater_Update(t *testing.T) {
	data := *MockProvider{}.AsMap()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reader, err := gzip.NewReader(r.Body)
		assert.Nil(t, err)

		body, err := io.ReadAll(reader)
		assert.Nil(t, err)

		var metric dto.Metrics

		err = json.Unmarshal(body, &metric)
		assert.Nil(t, err)

		value, found := data[metric.ID]

		assert.True(t, found)

		switch v := value.(type) {
		case float64:
			assert.Equal(t, v, *metric.Value)
		case int64:
			assert.Equal(t, v, *metric.Delta)
		default:
			t.Errorf("unknown metric type: %T", value)

		}
	}))

	defer server.Close()

	mockProvider := MockProvider{}
	updater := New(mockProvider, server.URL)

	updater.Update()
}
