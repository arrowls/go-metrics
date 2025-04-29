package repository

import (
	"log"
	"time"

	"github.com/arrowls/go-metrics/internal/config"
	"github.com/arrowls/go-metrics/internal/database"
	"github.com/arrowls/go-metrics/internal/di"
	"github.com/arrowls/go-metrics/internal/logger"
	"github.com/arrowls/go-metrics/internal/memstorage"
)

const diKey = "repository"

func ProvideRepository(container di.ContainerInterface) *Repository {
	if repo, ok := container.Get(diKey).(*Repository); ok {
		return repo
	}

	serverConfig := config.ProvideServerConfig(container)
	loggerInst := logger.ProvideLogger(container)

	var repo *Repository

	if serverConfig.DatabaseDSN == "" {
		storage := memstorage.GetInstance()
		repo = WithRestore(
			NewRepository(storage),
			time.Duration(serverConfig.StoreInterval)*time.Second,
			serverConfig.StorageFilePath,
			serverConfig.Restore,
			loggerInst,
		)
	} else {
		pool, err := database.ProvideDatabasePool(container)
		if err != nil {
			log.Fatal(err)
		}

		repo = &Repository{
			Metric: NewPostgresRepository(pool, loggerInst),
		}
	}

	if err := container.Add(diKey, repo); err != nil {
		log.Fatal(err)
	}

	return repo
}
