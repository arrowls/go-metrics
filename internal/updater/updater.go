package updater

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/arrowls/go-metrics/internal/collector"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/mappers"
)

type Updater struct {
	provider  collector.MetricProvider
	serverURL string
	wg        *sync.WaitGroup
}

func New(provider collector.MetricProvider, serverURL string) MetricConsumer {
	if !strings.HasPrefix(serverURL, "http://") {
		serverURL = "http://" + serverURL
	}
	return &Updater{
		provider,
		serverURL,
		&sync.WaitGroup{},
	}
}

func (u *Updater) Update() {
	data := u.provider.AsMap()

	for metricType, metricValue := range *data {
		updateDto, err := mappers.MetricToDTO(metricType, metricValue)
		if err != nil {
			continue
		}

		u.wg.Add(1)
		go u.updateFromDto(updateDto)
	}

	u.wg.Wait()
}

func (u *Updater) updateFromDto(updateDto *dto.Metrics) {
	defer u.wg.Done()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if err := json.NewEncoder(gz).Encode(updateDto); err != nil {
		fmt.Printf("Error marshaling postDto to JSON: %v\n", err)
		return
	}

	if errClose := gz.Close(); errClose != nil {
		fmt.Printf("Error closing gzip writer: %+v", errClose)
		return
	}
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

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			fmt.Printf("Error closing Body: %+v", errClose)
		}
	}()
}
