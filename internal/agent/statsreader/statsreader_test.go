package statsreader

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefresh(t *testing.T) {
	metricsDump, err := NewMetricsDump()
	assert.NoError(t, err)
	metricsDump.Refresh()
	metricsDump.Refresh()
	metricsDump.Refresh()

	assert.Equal(t, 3, int(metricsDump.MetricsCounter["PollCount"]))
}

func ExampleMetricsDump() {
	metricsDump, err := NewMetricsDump()
	if err != nil {
		log.Fatal(err)
	}

	metricsDump.Refresh()

	fmt.Println(metricsDump.MetricsGauge)
	fmt.Println(metricsDump.MetricsCounter)
}

func TestRefreshExtra(t *testing.T) {
	metricsDump, err := NewMetricsDump()
	assert.NoError(t, err)

	err = metricsDump.RefreshExtra()
	assert.NoError(t, err)

	t.Log(metricsDump.MetricsCounter)
	_, ok := metricsDump.MetricsGauge["TotalMemory"]
	assert.True(t, ok)

	_, ok = metricsDump.MetricsGauge["FreeMemory"]
	assert.True(t, ok)
}
