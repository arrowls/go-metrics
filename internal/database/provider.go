package database

import (
	"context"
	"log"
	"time"

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

	dbConfig, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
	if err != nil {
		loggerInst.Fatal("failed to parse database config: " + err.Error())
		return nil, err
	}

	dbConfig.MaxConns = 20
	dbConfig.MinConns = 5
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = 30 * time.Minute
	dbConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
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
