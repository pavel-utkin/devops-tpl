package storage

import (
	"context"
	"database/sql"
	"devops-tpl/internal/server/config"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"os"
	"log"
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
	repository.InitTables()

	return repository, nil
}

func (repository DBRepo) InitTables() error {
	_, err := repository.db.Exec("CREATE TABLE IF NOT EXISTS counter (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value BIGINT NOT NULL)")
	if err != nil {
		return fmt.Errorf("failed to create counter table: %w", err)
	}

	_, err = repository.db.Exec("CREATE TABLE IF NOT EXISTS gauge (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value DOUBLE PRECISION NOT NULL)")
	if err != nil {
		return fmt.Errorf("failed to create gauge table: %w", err)
	}

	return nil
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
	_, err := repository.db.Exec("INSERT INTO gauge (name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE set value = $2", key, *newMetricValue.Value)
	return err
}

func (repository DBRepo) updateCounter(key string, newMetricValue MetricValue) error {
	_, err := repository.db.Exec("INSERT INTO counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = counter.value + $2", key, *newMetricValue.Delta)
	return err
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
	metricValue := MetricValue{
		MType: MeticTypeGauge,
	}

	err := repository.db.QueryRow("SELECT value FROM gauge WHERE name = $1", key).Scan(&metricValue.Value)
	if err != nil {
		return metricValue, fmt.Errorf("gauge select error : %w", err)
	}
	return metricValue, nil
}

func (repository DBRepo) readCounter(key string) (MetricValue, error) {
	metricValue := MetricValue{
		MType: MeticTypeCounter,
	}

	err := repository.db.QueryRow("SELECT value FROM counter WHERE name = $1", key).Scan(&metricValue.Delta)
	if err != nil {
		return metricValue, fmt.Errorf("counter select error : %w", err)
	}
	return metricValue, nil
}

func (repository DBRepo) InitStateValues(DBSchema map[string]MetricValue) {
	for metricKey, metricValue := range DBSchema {
		repository.Update(metricKey, metricValue)
	}
}

func (repository DBRepo) ReadAll() map[string]MetricMap {
	var err error
	AllValues := map[string]MetricMap{}

	AllValues[MeticTypeCounter], err = repository.readAllCounter()
	log.Println(AllValues[MeticTypeCounter])
	if err != nil {
		return AllValues
	}

	AllValues[MeticTypeGauge], err = repository.readAllGauge()
	if err != nil {
		return AllValues
	}

	return AllValues
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

func (repository DBRepo) InitFromFile() {
	file, err := os.OpenFile(repository.config.File, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	var metricsDump map[string]MetricMap
	json.NewDecoder(file).Decode(&metricsDump)

	for _, metricList := range metricsDump {
		repository.InitStateValues(metricList)
	}
}
func (repository DBRepo) readAllCounter() (map[string]MetricValue, error) {
	allValues := map[string]MetricValue{}

	rows, err := repository.db.Query("SELECT name, value from counter")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var vKey string
		v := MetricValue{
			MType: MeticTypeCounter,
		}

		err = rows.Scan(&vKey, &v.Delta)
		if err != nil {
			return nil, err
		}

		allValues[vKey] = v
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return allValues, nil
}

func (repository DBRepo) readAllGauge() (map[string]MetricValue, error) {
	allValues := map[string]MetricValue{}

	rows, err := repository.db.Query("SELECT name, value from gauge")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var vKey string
		v := MetricValue{
			MType: MeticTypeGauge,
		}

		err = rows.Scan(&vKey, &v.Value)
		if err != nil {
			return nil, err
		}

		allValues[vKey] = v
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return allValues, nil
}
