package statsreader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRefresh(t *testing.T) {
	var memStatistics MemoryStatsDump
	memStatistics.Refresh()
	memStatistics.Refresh()
	memStatistics.Refresh()

	assert.Equal(t, 3, int(memStatistics.PollCount))
}
