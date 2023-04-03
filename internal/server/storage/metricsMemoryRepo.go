package storage

import (
	"devops-tpl/internal/server/config"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

const syncUploadSymbol = time.Duration(0)

type MetricValue struct {
	MType string   `json:"type" valid:"required,in(counter|gauge)"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type Metric struct {
	ID    string   `json:"id" valid:"required"`
	MType string   `json:"type" valid:"required,in(counter|gauge)"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (metric MetricValue) GetStringValue() string {
	switch metric.MType {
	case MeticTypeGauge:
		return fmt.Sprintf("%v", *metric.Value)
	case MeticTypeCounter:
		return fmt.Sprintf("%v", *metric.Delta)
	default:
		return ""
	}
}

// MetricsMemoryRepo - репо для приходящей статистики
type MetricsMemoryRepo struct {
	uploadMutex    *sync.RWMutex
	gaugeStorage   *MemoryRepo
	counterStorage *MemoryRepo
	config         config.StoreConfig
}

func NewMetricsMemoryRepo(config config.StoreConfig) MetricsMemoryRepo {
	var mmr MetricsMemoryRepo
	var err error

	mmr.config = config
	mmr.uploadMutex = &sync.RWMutex{}
	mmr.gaugeStorage, err = NewMemoryRepo()
	if err != nil {
		panic("gaugeMemoryRepo init error")
	}
	mmr.counterStorage, err = NewMemoryRepo()
	if err != nil {
		panic("counterMemoryRepo init error")
	}

	if mmr.config.Interval != syncUploadSymbol {
		mmr.IterativeUploadToFile()
	}

	return mmr
}

func (mmr MetricsMemoryRepo) Update(key string, newMetricValue MetricValue) error {
	switch newMetricValue.MType {
	case MeticTypeGauge:
		if newMetricValue.Value == nil {
			return errors.New("Metric Value is empty")
		}
		newMetricValue.Delta = nil

		return mmr.updateGaugeValue(key, newMetricValue)
	case MeticTypeCounter:
		if newMetricValue.Delta == nil {
			return errors.New("Metric Delta is empty")
		}
		newMetricValue.Value = nil

		return mmr.updateCounterValue(key, newMetricValue)
	default:
		return errors.New("Metric type is not defined")
	}
}

func (mmr MetricsMemoryRepo) updateGaugeValue(key string, newMetricValue MetricValue) error {
	mmr.uploadMutex.Lock()
	err := mmr.gaugeStorage.Write(key, newMetricValue)
	mmr.uploadMutex.Unlock()

	if err != nil {
		return err
	}

	if mmr.config.Interval == syncUploadSymbol {
		return mmr.UploadToFile()
	}

	return nil
}

func (mmr MetricsMemoryRepo) updateCounterValue(key string, newMetricValue MetricValue) error {
	//Чтение старого значения
	oldMetricValue, err := mmr.Read(key, MeticTypeCounter)
	if err != nil {
		var delta int64 = 0
		oldMetricValue = MetricValue{
			Delta: &delta,
		}
	}

	newValue := *oldMetricValue.Delta + *newMetricValue.Delta
	newMetricValue.Delta = &newValue

	mmr.uploadMutex.Lock()
	mmr.counterStorage.Write(key, newMetricValue)
	mmr.uploadMutex.Unlock()

	if mmr.config.Interval == syncUploadSymbol {
		return mmr.UploadToFile()
	}

	return nil
}

func (mmr MetricsMemoryRepo) Read(key string, metricType string) (MetricValue, error) {
	switch metricType {
	case MeticTypeGauge:
		return mmr.gaugeStorage.Read(key)
	case MeticTypeCounter:
		return mmr.counterStorage.Read(key)
	default:
		return MetricValue{}, errors.New("metricType not found")
	}
}

func (mmr MetricsMemoryRepo) UploadToFile() error {
	mmr.uploadMutex.Lock()
	defer mmr.uploadMutex.Unlock()
	if mmr.config.File == "" {
		return nil
	}

	file, err := os.OpenFile(mmr.config.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	allStates := mmr.ReadAll()
	err = json.NewEncoder(file).Encode(allStates)
	if err != nil {
		return err
	}

	return nil
}

func (mmr MetricsMemoryRepo) IterativeUploadToFile() error {
	tickerUpload := time.NewTicker(mmr.config.Interval)

	go func() {
		for range tickerUpload.C {
			mmr.UploadToFile()
		}
	}()

	return nil
}

func (mmr MetricsMemoryRepo) InitFromFile() {
	file, err := os.OpenFile(mmr.config.File, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err.Error())
		//panic(err.Error())
	}
	defer file.Close()

	var metricsDump map[string]MetricMap
	err = json.NewDecoder(file).Decode(&metricsDump)
	if err != nil {
		fmt.Println(err.Error())
		//panic(err.Error())
	}

	for _, metricList := range metricsDump {
		mmr.InitStateValues(metricList)
	}
}

func (mmr MetricsMemoryRepo) InitStateValues(DBSchema map[string]MetricValue) {
	for metricKey, metricValue := range DBSchema {
		mmr.Update(metricKey, metricValue)
	}
}

func (mmr MetricsMemoryRepo) ReadAll() map[string]MetricMap {
	return map[string]MetricMap{
		MeticTypeGauge:   mmr.gaugeStorage.GetSchemaDump(),
		MeticTypeCounter: mmr.counterStorage.GetSchemaDump(),
	}
}

func (mmr MetricsMemoryRepo) Close() error {
	err := mmr.gaugeStorage.Close()
	if err != nil {
		return err
	}
	err = mmr.counterStorage.Close()

	return err
}
