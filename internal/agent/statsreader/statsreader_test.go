package statsreader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRefresh(t *testing.T) {
	var metricsDump MetricsDump
	metricsDump.Refresh()
	metricsDump.Refresh()
	metricsDump.Refresh()

	assert.Equal(t, 3, int(metricsDump.PollCount))
}
