package updater

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/arrowls/go-metrics/internal/collector"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/mappers"
	"github.com/arrowls/go-metrics/internal/utils"
	"github.com/sirupsen/logrus"
)

type Updater struct {
	provider  collector.MetricProvider
	serverURL string
	logger    *logrus.Logger
}

func New(provider collector.MetricProvider, serverURL string, logger *logrus.Logger) MetricConsumer {
	if !strings.HasPrefix(serverURL, "http://") {
		serverURL = "http://" + serverURL
	}
	return &Updater{
		provider,
		serverURL,
		logger,
	}
}

func (u *Updater) Update() {
	data := u.provider.AsMap()
	var metrics []*dto.Metrics

	for metricType, metricValue := range *data {
		updateDto, err := mappers.MetricToDTO(metricType, metricValue)
		if err != nil {
			continue
		}

		metrics = append(metrics, updateDto)
	}

	_ = utils.WithRetry(func() (bool, error) {
		err := u.updateFromDto(metrics)
		if err == nil {
			return false, nil
		}

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return true, netErr
		}

		if strings.Contains(err.Error(), "connection refused") {
			return true, err
		}

		if errors.Is(err, io.EOF) {
			return true, err
		}

		return false, nil
	})
}

func (u *Updater) updateFromDto(updateDto []*dto.Metrics) error {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if err := json.NewEncoder(gz).Encode(updateDto); err != nil {
		u.logger.Errorf("Error marshaling postDto to JSON: %v\n", err)
		return err
	}

	if errClose := gz.Close(); errClose != nil {
		u.logger.Errorf("Error closing gzip writer: %+v", errClose)
		return errClose
	}
	req, err := http.NewRequest("POST", u.serverURL+"/updates", &buf)
	if err != nil {
		u.logger.Errorf("Error creating request: %v\n", err)
		return err
	}
	req.Header.Set("Content-Encoding", "gzip")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		u.logger.Errorf("Error posting metric: %v\n", err)
		return err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			u.logger.Errorf("Error closing Body: %+v", errClose)
		}
	}()

	return nil
}
