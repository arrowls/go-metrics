package updater

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/arrowls/go-metrics/internal/collector"
	"github.com/arrowls/go-metrics/internal/dto"
)

type Updater struct {
	provider  collector.MetricProvider
	serverURL string
}

func New(provider collector.MetricProvider, serverURL string) MetricConsumer {
	if !strings.HasPrefix(serverURL, "http://") {
		serverURL = "http://" + serverURL
	}
	return &Updater{
		provider,
		serverURL,
	}
}

func (u *Updater) Update() {
	data := u.provider.AsMap()
	var updateDto dto.Metrics

	for metricType, metricValue := range *data {
		updateDto.ID = metricType

		switch v := metricValue.(type) {
		case float64:
			updateDto.MType = "gauge"
			updateDto.Value = &v

		case int64:
			updateDto.MType = "counter"
			updateDto.Delta = &v

		default:
			fmt.Printf("Unsupported metric type: %s\n", metricValue)
		}

		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		if err := json.NewEncoder(gz).Encode(updateDto); err != nil {
			fmt.Printf("Error marshaling postDto to JSON: %v\n", err)
			return
		}

		gz.Close()

		req, err := http.NewRequest("POST", u.serverURL+"/update", &buf)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}
		req.Header.Set("Content-Encoding", "gzip")

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			fmt.Printf("Error posting metric: %v\n", err)
			return
		}
		res.Body.Close()
	}
}
