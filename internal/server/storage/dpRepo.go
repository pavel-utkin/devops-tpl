package storage

import (
	"context"
	"database/sql"
	"devops-tpl/internal/server/config"
	"errors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

type DBRepo struct {
	config config.StoreConfig
	db     *sql.DB
}

func NewDBRepo(config config.StoreConfig) (DBRepo, error) {
	var repository DBRepo
	repository.config = config

	db, err := sql.Open("pgx",
		repository.config.DatabaseDSN)
	if err != nil {
		return DBRepo{}, err
	}
	repository.db = db

	return repository, nil
}

func (repository DBRepo) Update(key string, newMetricValue MetricValue) error {
	switch newMetricValue.MType {
	case MeticTypeGauge:
		if newMetricValue.Value == nil {
			return errors.New("Metric Value is empty")
		}
		newMetricValue.Delta = nil

		return repository.updateGauge(key, newMetricValue)
	case MeticTypeCounter:
		if newMetricValue.Delta == nil {
			return errors.New("Metric Delta is empty")
		}
		newMetricValue.Value = nil

		return repository.updateCounter(key, newMetricValue)
	default:
		return errors.New("Metric type is not defined")
	}
}

func (repository DBRepo) updateGauge(key string, newMetricValue MetricValue) error {
	return nil
}

func (repository DBRepo) updateCounter(key string, newMetricValue MetricValue) error {
	return nil
}

func (repository DBRepo) Read(key string, metricType string) (MetricValue, error) {
	switch metricType {
	case MeticTypeGauge:
		return repository.readGauge(key)
	case MeticTypeCounter:
		return repository.readCounter(key)
	default:
		return MetricValue{}, errors.New("metricType not found")
	}
}

func (repository DBRepo) readGauge(key string) (MetricValue, error) {
	return MetricValue{}, errors.New("not found")
}

func (repository DBRepo) readCounter(key string) (MetricValue, error) {
	return MetricValue{}, errors.New("not found")
}

func (repository DBRepo) InitStateValues(DBSchema map[string]MetricValue) {
	for metricKey, metricValue := range DBSchema {
		repository.Update(metricKey, metricValue)
	}
}

func (repository DBRepo) ReadAll() map[string]MetricMap {
	return map[string]MetricMap{}
}

func (repository DBRepo) Close() error {
	return repository.db.Close()
}

func (repository DBRepo) Ping() error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := repository.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
