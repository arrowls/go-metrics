package updater

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/arrowls/go-metrics/internal/collector"
	"github.com/arrowls/go-metrics/internal/config"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/mappers"
	"github.com/arrowls/go-metrics/internal/utils"
	"github.com/sirupsen/logrus"
)

type Updater struct {
	provider  collector.MetricProvider
	serverURL string
	logger    *logrus.Logger
	encodeKey string
	dataChan  chan *map[string]interface{}
}

func New(provider collector.MetricProvider, serverURL string, logger *logrus.Logger, encodeKey string, dataChan chan *map[string]interface{}) MetricConsumer {
	if !strings.HasPrefix(serverURL, "http://") {
		serverURL = "http://" + serverURL
	}
	return &Updater{
		provider,
		serverURL,
		logger,
		encodeKey,
		dataChan,
	}
}

func (u *Updater) Update() {
	var batch []*map[string]interface{}

loop:
	for {
		select {
		case data := <-u.dataChan:
			batch = append(batch, data)
		default:
			break loop
		}
	}

	for _, data := range batch {
		go func(data *map[string]interface{}) {
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

				u.logger.Error(err)

				return false, nil
			})
		}(data)
	}
}

func (u *Updater) updateFromDto(updateDto []*dto.Metrics) error {
	jsonBody, err := json.Marshal(updateDto)
	if err != nil {
		u.logger.Errorf("error converting data to JSON: %v", err)
		return err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(jsonBody); err != nil {
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

	if u.encodeKey != "" {
		hasher := hmac.New(sha256.New, []byte(u.encodeKey))
		hasher.Write(jsonBody)
		sum := hex.EncodeToString(hasher.Sum(nil))
		req.Header.Set(config.HashHeaderName, sum)
	}

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
