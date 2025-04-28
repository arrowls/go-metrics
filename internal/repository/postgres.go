package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/logger"
	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Metric {
	return &PostgresRepository{
		db,
	}
}
func (m *PostgresRepository) AddGaugeValue(ctx context.Context, name string, value float64) error {
	loggerInst := logger.Inject(ctx)

	_, err := m.db.Exec(ctx, `
	  INSERT INTO gauges (name, value) 
	  VALUES ($1, $2)
	  ON CONFLICT (name) DO UPDATE
	  SET value = $2`,
		name, strconv.FormatFloat(value, 'f', -1, 64),
	)

	if err != nil {
		loggerInst.Error(err)
		return fmt.Errorf("error while inserting gauge value")
	}
	return nil
}

func (m *PostgresRepository) AddCounterValue(ctx context.Context, name string, value int64) error {
	loggerInst := logger.Inject(ctx)

	_, err := m.db.Exec(ctx, `
	  INSERT INTO counters (name, value) 
	  VALUES ($1, $2)
	  ON CONFLICT (name) DO UPDATE
	  SET value = counters.value + $2`,
		name, value,
	)

	if err != nil {
		loggerInst.Error(err)
		return fmt.Errorf("error while inserting counter value")
	}
	return nil
}

func (m *PostgresRepository) GetAll(ctx context.Context) (memstorage.MemStorage, error) {
	storage := memstorage.GetInstance()
	loggerInst := logger.Inject(ctx)

	gaugeRows, err := m.db.Query(ctx, "SELECT name, value FROM gauges")
	if err != nil {
		loggerInst.Error(err)
		return *storage, fmt.Errorf("error while fetching gauges from database")
	}

	for gaugeRows.Next() {
		var name string
		var value float64
		err = gaugeRows.Scan(&name, &value)
		if err == nil {
			storage.Gauge[name] = value
		}
	}

	counterRows, err := m.db.Query(ctx, "SELECT name, value FROM counters")
	if err != nil {
		loggerInst.Error(err)
		return *storage, fmt.Errorf("error while fetching counters from database: %w", err)
	}

	for counterRows.Next() {
		var name string
		var value int64
		err = gaugeRows.Scan(&name, &value)
		if err == nil {
			storage.Counter[name] = value
		}
	}

	return *storage, nil
}

func (m *PostgresRepository) GetGaugeItem(ctx context.Context, name string) (float64, error) {
	loggerInst := logger.Inject(ctx)
	var value float64

	row := m.db.QueryRow(ctx, "SELECT value FROM gauges WHERE name = $1", name)
	err := row.Scan(&value)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errors.Join(apperrors.ErrNotFound, fmt.Errorf("gauge %s not found", name))
	}

	if err != nil {
		loggerInst.Error(err)
		return 0, fmt.Errorf("error fetching value from database")
	}

	return value, nil
}

func (m *PostgresRepository) GetCounterItem(ctx context.Context, name string) (int64, error) {
	loggerInst := logger.Inject(ctx)
	var value int64

	row := m.db.QueryRow(ctx, "SELECT value FROM counters WHERE name = $1", name)
	err := row.Scan(&value)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errors.Join(apperrors.ErrNotFound, fmt.Errorf("counter %s not found", name))
	}

	if err != nil {
		loggerInst.Error(err)
		return 0, fmt.Errorf("error fetching value from database")
	}

	return value, nil
}

func (m *PostgresRepository) CheckConnection(ctx context.Context) bool {
	if err := m.db.Ping(ctx); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
