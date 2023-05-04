package statsreader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRefresh(t *testing.T) {
	metricsDump, err := NewMetricsDump()
	assert.NoError(t, err)
	metricsDump.Refresh()
	metricsDump.Refresh()
	metricsDump.Refresh()

	assert.Equal(t, 3, int(metricsDump.MetricsCounter["PollCount"]))
}
