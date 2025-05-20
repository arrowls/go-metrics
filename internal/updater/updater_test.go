package updater

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reader, err := gzip.NewReader(r.Body)
		assert.Nil(t, err)

		body, err := io.ReadAll(reader)
		assert.Nil(t, err)

		expectedBody := `[{"id":"MetricName2","type":"gauge","value":2.22},{"id":"MetricName3","type":"counter","delta":3},{"id":"MetricName4","type":"gauge","value":4.4},{"id":"MetricName1","type":"counter","delta":1}]`

		var actual, expected []dto.Metrics

		require.NoError(t, json.Unmarshal(body, &actual), "Ошибка декодирования body")
		require.NoError(t, json.Unmarshal([]byte(expectedBody), &expected), "Ошибка декодирования expectedBody")

		expectedMap := make(map[string]dto.Metrics)
		for _, m := range expected {
			key := m.ID + "|" + m.MType
			expectedMap[key] = m
		}

		for _, m := range actual {
			key := m.ID + "|" + m.MType
			expectedMetric, ok := expectedMap[key]
			require.True(t, ok, "Метрика с ID=%s и типом=%s не найдена в ожидаемых данных", m.ID, m.MType)

			if expectedMetric.Delta == nil {
				require.Nil(t, m.Delta, "Delta для метрики ID=%s должно быть nil", m.ID)
			} else {
				require.NotNil(t, m.Delta, "Delta для метрики ID=%s не должно быть nil", m.ID)
				require.Equal(t, *expectedMetric.Delta, *m.Delta, "Delta для метрики ID=%s не совпадает", m.ID)
			}

			if expectedMetric.Value == nil {
				require.Nil(t, m.Value, "Value для метрики ID=%s должно быть nil", m.ID)
			} else {
				require.NotNil(t, m.Value, "Value для метрики ID=%s не должно быть nil", m.ID)
				require.Equal(t, *expectedMetric.Value, *m.Value, "Value для метрики ID=%s не совпадает", m.ID)
			}
		}
	}))

	defer server.Close()

	mockProvider := MockProvider{}
	logger := logrus.New()
	logger.SetOutput(io.Discard)

	ch := make(chan *map[string]interface{}, 1)
	updater := New(mockProvider, server.URL, logger, "", ch)

	ch <- mockProvider.AsMap()
	updater.Update()
}
