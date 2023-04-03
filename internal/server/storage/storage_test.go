package storage

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemoryRepoRW(t *testing.T) {
	memoryRepo, err := NewMemoryRepo()
	require.NoError(t, err)

	var counterValueExpect int64 = 50
	memoryRepo.Write("PollCount", MetricValue{
		MType: "counter",
		Delta: &counterValueExpect,
	})
	counterValueReal, err := memoryRepo.Read("PollCount")

	require.NoError(t, err)
	require.Equal(t, counterValueExpect, *counterValueReal.Delta)
}

func TestMemoryRepoReadEmpty(t *testing.T) {
	memoryRepo, err := NewMemoryRepo()
	require.NoError(t, err)
	_, err = memoryRepo.Read("username")
	require.Error(t, err)
}

func TestUpdateCounterValue(t *testing.T) {
	memStatsStorage := NewMemStatsMemoryRepo()

	var startValue int64 = 7
	var incrementValue int64 = 22
	err := memStatsStorage.Update("PollCount", MetricValue{
		MType: "counter",
		Delta: &startValue,
	})
	require.NoError(t, err)
	err = memStatsStorage.Update("PollCount", MetricValue{
		MType: "counter",
		Delta: &incrementValue,
	})
	require.NoError(t, err)
	PollCount, err := memStatsStorage.ReadValue("PollCount", "counter")
	require.NoError(t, err)

	require.Equal(t, int64(29), *PollCount.Delta)
}
