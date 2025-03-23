package updater

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/arrowls/go-metrics/internal/collector"
)

type Updater struct {
	provider  collector.MetricProvider
	serverURL string
}

func New(provider collector.MetricProvider, serverURL string) MetricConsumer {
	return &Updater{
		provider,
		serverURL,
	}
}

func (u *Updater) Update() {
	data := u.provider.AsMap()

	for metricType, metricValue := range *data {
		switch reflect.TypeOf(metricValue).Kind() {
		case reflect.Float64:
			u.postGauge(metricType, metricValue.(float64))

		case reflect.Int64:
			u.postCounter(metricType, metricValue.(int64))

		default:
			fmt.Printf("Unsupported metric type: %s\n", metricValue)
		}

	}
}

func (u *Updater) postGauge(metricType string, metricValue float64) {
	url := fmt.Sprintf("%s/update/gauge/%s/%f", u.serverURL, metricType, metricValue)

	resp, err := http.Post(url, "text/plain", nil)

	if err != nil {
		fmt.Printf("Error posting metric to server: %v\n", err)
	}

	resp.Body.Close()
}
func (u *Updater) postCounter(metricType string, metricValue int64) {
	url := fmt.Sprintf("%s/update/counter/%s/%d", u.serverURL, metricType, metricValue)
	resp, err := http.Post(url, "text/plain", nil)

	if err != nil {
		fmt.Printf("Error posting metric to server: %v\n", err)
	}
	resp.Body.Close()
}
