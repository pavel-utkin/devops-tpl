package requestorhandler

import (
	"devops-tpl/internal/agent/config"
	"devops-tpl/internal/agent/metricsreader"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"reflect"
	"strconv"
	"strings"
)

func oneStatUpload(httpClient *resty.Client, errorChan chan error, statType string, statName string, statValue string) {
	resp, err := httpClient.R().
		SetPathParams(map[string]string{
			"host":  config.ConfigServerHost,
			"port":  strconv.Itoa(config.ConfigServerPort),
			"type":  statType,
			"name":  statName,
			"value": statValue,
		}).
		SetHeader("Content-Type", "text/plain").
		Post("http://{host}:{port}/update/{type}/{name}/{value}")

	if err != nil {
		fmt.Println(err)
		errorChan <- err
	}
	if resp.StatusCode() != 200 {
		errorChan <- errors.New("HTTP Status != 200")
	}

	errorChan <- nil
}

func MemoryStatsUpload(httpClient *resty.Client, memoryStats metricsreader.MemoryStatsDump) error {
	reflectMemoryStats := reflect.ValueOf(memoryStats)
	typeOfMemoryStats := reflectMemoryStats.Type()
	errorChan := make(chan error)
	defer close(errorChan)

	for i := 0; i < reflectMemoryStats.NumField(); i++ {
		statName := typeOfMemoryStats.Field(i).Name
		statValue := fmt.Sprintf("%v", reflectMemoryStats.Field(i).Interface())
		statType := strings.Split(typeOfMemoryStats.Field(i).Type.String(), ".")[1]

		go oneStatUpload(httpClient, errorChan, statType, statName, statValue)
	}

	for i := 0; i < reflectMemoryStats.NumField(); i++ {
		error := <-errorChan

		if error != nil {
			return error
		}
	}

	return nil
}
