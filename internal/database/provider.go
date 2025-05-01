package database

import (
	"context"
	"log"

	"github.com/arrowls/go-metrics/internal/config"
	"github.com/arrowls/go-metrics/internal/di"
	"github.com/arrowls/go-metrics/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const diKey = "database_pool"

func ProvideDatabasePool(container di.ContainerInterface) (*pgxpool.Pool, error) {
	cfg := config.ProvideServerConfig(container)
	loggerInst := logger.ProvideLogger(container)

	poolInst := container.Get(diKey)
	if pool, ok := poolInst.(*pgxpool.Pool); ok {
		return pool, nil
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		loggerInst.Fatal("failed to connect to database: " + err.Error())
		return nil, err
	}

	TryMigrations(cfg.DatabaseDSN, loggerInst)

	if err = container.Add(diKey, pool); err != nil {
		log.Fatal(err)
	}

	return pool, nil
}
