package result

import (
	"time"

	"github.com/GrapefruitCat030/gfc_dcache/tools/benchmark/client"
)

// TimeStat is a struct to store the time statistics
type TimeStat struct {
	Count int           // Count is the number of operations
	Time  time.Duration // Time is the total time of the operations
}

type BenchmarkRes struct {
	GetCount         int
	MissCount        int
	SetCount         int
	StatisticBuckets []TimeStat // [idx, TimeStat]: idx set by the time in milliseconds
}

func (r *BenchmarkRes) timeStatistic(bucket int, stat TimeStat) {
	if bucket > len(r.StatisticBuckets)-1 {
		newBuckets := make([]TimeStat, bucket+1)
		copy(newBuckets, r.StatisticBuckets)
		r.StatisticBuckets = newBuckets
	}
	r.StatisticBuckets[bucket].Count += stat.Count
	r.StatisticBuckets[bucket].Time += stat.Time
}

func (r *BenchmarkRes) CollectTimeCost(dur time.Duration, opType string) {
	bucketIdx := int(dur / time.Millisecond)
	r.timeStatistic(bucketIdx, TimeStat{1, dur})
	switch opType {
	case client.OperationTypeGet:
		r.GetCount++
	case client.OperationTypeSet:
		r.SetCount++
	default:
		r.MissCount++
	}
}

func (r *BenchmarkRes) AddResult(src *BenchmarkRes) {
	r.GetCount += src.GetCount
	r.MissCount += src.MissCount
	r.SetCount += src.SetCount
	for idx, stat := range src.StatisticBuckets {
		r.timeStatistic(idx, stat)
	}
}
