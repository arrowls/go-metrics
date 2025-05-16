package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/config"
	"github.com/arrowls/go-metrics/internal/database"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/arrowls/go-metrics/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type PostgresRepository struct {
	db     *pgxpool.Pool
	logger *logrus.Logger
}

func NewPostgresRepository(db *pgxpool.Pool, logger *logrus.Logger) Metric {
	return &PostgresRepository{
		db,
		logger,
	}
}
func (m *PostgresRepository) AddGaugeValue(ctx context.Context, name string, value float64) error {
	err := utils.WithRetry(func() (bool, error) {
		_, dbErr := m.db.Exec(ctx, `
		  INSERT INTO gauges (name, value) 
		  VALUES ($1, $2)
		  ON CONFLICT (name) DO UPDATE
		  SET value = $2
		  `,
			name, strconv.FormatFloat(value, 'f', -1, 64),
		)

		if dbErr == nil {
			return false, nil
		}

		if database.IsConnectionException(dbErr.Error()) {
			return true, dbErr
		}
		return false, dbErr
	})

	if err != nil {
		m.logger.Error(err)
		return fmt.Errorf("error while inserting gauge value")
	}
	return nil
}

func (m *PostgresRepository) AddCounterValue(ctx context.Context, name string, value int64) error {
	err := utils.WithRetry(func() (bool, error) {
		_, dbErr := m.db.Exec(ctx, `
		  INSERT INTO counters (name, value) 
		  VALUES ($1, $2)
		  ON CONFLICT (name) DO UPDATE
		  SET value = counters.value + $2`,
			name, value,
		)

		if dbErr == nil {
			return false, nil
		}

		if database.IsConnectionException(dbErr.Error()) {
			return true, dbErr
		}
		return false, dbErr
	})

	if err != nil {
		m.logger.Error(err)
		return fmt.Errorf("error while inserting counter value")
	}
	return nil
}

func (m *PostgresRepository) GetAll(ctx context.Context) (memstorage.MemStorage, error) {
	storage := memstorage.GetInstance()

	err := utils.WithRetry(func() (bool, error) {
		gaugeRows, dbErr := m.db.Query(ctx, "SELECT name, value FROM gauges")
		if dbErr != nil {
			m.logger.Error(dbErr)

			if database.IsConnectionException(dbErr.Error()) {
				return true, dbErr
			}

			return false, fmt.Errorf("error while fetching gauges from database")
		}

		defer gaugeRows.Close()

		for gaugeRows.Next() {
			var name string
			var value float64
			err := gaugeRows.Scan(&name, &value)
			if err != nil {
				m.logger.Errorf("error while fetching gauge from database: %v", err)
				continue
			}

			storage.Gauge[name] = value
		}

		if err := gaugeRows.Err(); err != nil {
			m.logger.Errorf("error while fetching gauge rows from database: %v", err)
			return false, err
		}

		return false, nil
	})

	if err != nil {
		return *storage, err
	}

	err = utils.WithRetry(func() (bool, error) {
		counterRows, dbErr := m.db.Query(ctx, "SELECT name, value FROM counters")
		if dbErr != nil {
			m.logger.Error(dbErr)

			if database.IsConnectionException(dbErr.Error()) {
				return true, dbErr
			}

			return false, fmt.Errorf("error while fetching gauges from database")
		}

		defer counterRows.Close()

		for counterRows.Next() {
			var name string
			var value int64
			err = counterRows.Scan(&name, &value)
			if err != nil {
				m.logger.Errorf("error while fetching counter from database: %v", err)
				continue
			}

			storage.Counter[name] = value
		}

		if err := counterRows.Err(); err != nil {
			m.logger.Errorf("error while fetching counter rows from database: %v", err)
			return false, err
		}

		return false, nil
	})

	if err != nil {
		return *storage, err
	}

	return *storage, nil
}

func (m *PostgresRepository) GetGaugeItem(ctx context.Context, name string) (float64, error) {
	var value float64

	row := m.db.QueryRow(ctx, "SELECT value FROM gauges WHERE name = $1", name)
	err := row.Scan(&value)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errors.Join(apperrors.ErrNotFound, fmt.Errorf("gauge %s not found", name))
	}

	if err != nil {
		m.logger.Error(err)
		return 0, fmt.Errorf("error fetching value from database")
	}

	return value, nil
}

func (m *PostgresRepository) GetCounterItem(ctx context.Context, name string) (int64, error) {
	var value int64

	row := m.db.QueryRow(ctx, "SELECT value FROM counters WHERE name = $1", name)
	err := row.Scan(&value)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errors.Join(apperrors.ErrNotFound, fmt.Errorf("counter %s not found", name))
	}

	if err != nil {
		m.logger.Error(err)
		return 0, fmt.Errorf("error fetching value from database")
	}

	return value, nil
}

func (m *PostgresRepository) CheckConnection(ctx context.Context) bool {
	if err := m.db.Ping(ctx); err != nil {
		return false
	}
	return true
}

func (m *PostgresRepository) CreateBatch(ctx context.Context, batch []dto.CreateMetric) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting tx")
	}

	defer func() {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			m.logger.Error(errRollback)
		}
	}()

	for _, metric := range batch {
		if metric.Type == config.GaugeType {
			_, err = tx.Exec(ctx, `
			  INSERT INTO gauges (name, value) 
			  VALUES ($1, $2)
			  ON CONFLICT (name) DO UPDATE
			  SET value = $2
			`, metric.Name, metric.Value)

			if err != nil {
				return fmt.Errorf("error creating a batch on %s", metric.Name)
			}
		}

		if metric.Type == config.CounterType {
			_, err = tx.Exec(ctx, `
				  INSERT INTO counters (name, value) 
				  VALUES ($1, $2)
				  ON CONFLICT (name) DO UPDATE
				  SET value = counters.value + $2`,
				metric.Name, metric.Value,
			)

			if err != nil {
				return fmt.Errorf("error creating a batch on %s", metric.Name)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error creating a batch")
	}

	return nil
}
