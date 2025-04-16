package updater

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
	caughtRequestsMap := map[string]bool{
		"/update/counter/MetricName1/1":      false,
		"/update/gauge/MetricName2/2.220000": false,
		"/update/counter/MetricName3/3":      false,
		"/update/gauge/MetricName4/4.400000": false,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, found := caughtRequestsMap[r.URL.Path]

		require.Truef(t, found, "%s not found in caught requests", r.URL.Path)
		require.Falsef(t, value, "Metric update called more than once: %s", r.URL.Path)

		caughtRequestsMap[r.URL.Path] = true
	}))

	defer server.Close()

	mockProvider := MockProvider{}
	updater := New(mockProvider, server.URL)

	updater.Update()

	for url, caught := range caughtRequestsMap {
		assert.Truef(t, caught, "expected to catch a request to %s", url)
	}
}
