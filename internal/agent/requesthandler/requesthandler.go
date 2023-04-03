package requesthandler

import (
	"devops-tpl/internal/agent/config"
	"devops-tpl/internal/agent/statsreader"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func oneStatUpload(httpClient *resty.Client, statType, statName, statValue string) error {
	resp, err := httpClient.R().
		SetPathParams(map[string]string{
			"host":  config.ServerHost,
			"port":  strconv.Itoa(config.ServerPort),
			"type":  statType,
			"name":  statName,
			"value": statValue,
		}).
		SetHeader("Content-Type", "text/plain").
		Post("http://{host}:{port}/update/{type}/{name}/{value}")

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New("HTTP Status != 200")
	}

	return nil
}

func oneStatUploadJSON(httpClient *resty.Client, statType string, statName string, statValue string) error {
	OneMetrics := struct {
		ID    string  `json:"id"`
		MType string  `json:"type"`
		Delta int64   `json:"delta"`
		Value float64 `json:"value"`
	}{
		ID:    statName,
		MType: statType,
	}

	var err error
	switch OneMetrics.MType {
	case "counter":
		var metricValue int64
		metricValue, err = strconv.ParseInt(statValue, 10, 64)
		OneMetrics.Delta = metricValue
	case "gauge":
		var metricValue float64
		metricValue, err = strconv.ParseFloat(statValue, 64)
		OneMetrics.Value = metricValue
	default:
		return errors.New("unknow statType")
	}
	if err != nil {
		return errors.New("invalid statValue")
	}

	statJSON, err := json.Marshal(OneMetrics)
	if err != nil {
		return err
	}

	resp, err := httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(statJSON)).
		SetPathParams(map[string]string{
			"host": config.ServerHost,
			"port": strconv.Itoa(config.ServerPort),
		}).
		Post("http://{host}:{port}/update/")

	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP Status: %v (not 200)", resp.StatusCode())
	}

	return nil
}

func MemoryStatsUpload(httpClient *resty.Client, memoryStats statsreader.MemoryStatsDump) error {
	reflectMemoryStats := reflect.ValueOf(memoryStats)
	typeOfMemoryStats := reflectMemoryStats.Type()
	errorGroup := new(errgroup.Group)

	for i := 0; i < reflectMemoryStats.NumField(); i++ {
		statName := typeOfMemoryStats.Field(i).Name
		statValue := fmt.Sprintf("%v", reflectMemoryStats.Field(i).Interface())
		statType := strings.Split(typeOfMemoryStats.Field(i).Type.String(), ".")[1]

		errorGroup.Go(func() error {
			return oneStatUpload(httpClient, statType, statName, statValue)
		})
	}

	return errorGroup.Wait()
}
