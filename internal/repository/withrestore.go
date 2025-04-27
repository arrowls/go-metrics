package repository

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/sirupsen/logrus"
)

type restorableRepository struct {
	Metric
	storeInterval time.Duration
	storagePath   string
	logger        *logrus.Logger
}

func WithRestore(repo *Repository, storeInterval time.Duration, storagePath string, needRestore bool, logger *logrus.Logger) *Repository {
	restoreRepo := &restorableRepository{
		repo.Metric,
		storeInterval,
		storagePath,
		logger,
	}

	if storeInterval != 0 {
		go restoreRepo.runAsyncAction()
	}

	if needRestore {
		restoreRepo.restore()
	}

	return &Repository{Metric: restoreRepo}
}

func (r *restorableRepository) isSync() bool {
	return r.storeInterval == 0
}

func (r *restorableRepository) runSyncAction() {
	file, err := os.OpenFile(r.storagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		r.logger.Errorf("error while opening storage file: %+v", err)
		return
	}

	defer func() {
		if errClose := file.Close(); err != nil {
			err = errors.Join(err, errClose)
		}
	}()

	data := r.GetAll(context.Background())

	fileContent, err := json.Marshal(data)
	if err != nil {
		r.logger.Errorf("error while marshaling json: %+v", err)
		return
	}

	_, err = file.Write(fileContent)
	if err != nil {
		r.logger.Errorf("error while writing to file: %+v", err)
	}
}

func (r *restorableRepository) runAsyncAction() {
	for {
		time.Sleep(r.storeInterval)
		r.runSyncAction()
	}
}

func (r *restorableRepository) restore() {
	file, err := os.OpenFile(r.storagePath, os.O_RDONLY, 0666)
	if err != nil {
		r.logger.Errorf("error while opening storage file: %+v", err)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		r.logger.Errorf("error while reading storage file: %+v", err)
		return
	}

	var restoredStorage memstorage.MemStorage
	if err = json.Unmarshal(fileBytes, &restoredStorage); err != nil {
		r.logger.Errorf("error while decoding file: %+v", err)
		return
	}

	ctx := context.Background()

	for name, value := range restoredStorage.Counter {
		r.Metric.AddCounterValue(ctx, name, value)
	}

	for name, value := range restoredStorage.Gauge {
		r.Metric.AddGaugeValue(ctx, name, value)
	}
}

func (r *restorableRepository) AddGaugeValue(ctx context.Context, name string, value float64) {
	if r.isSync() {
		r.runSyncAction()
	}

	r.Metric.AddGaugeValue(ctx, name, value)
}

func (r *restorableRepository) AddCounterValue(ctx context.Context, name string, value int64) {
	if r.isSync() {
		r.runSyncAction()
	}

	r.Metric.AddCounterValue(ctx, name, value)
}
