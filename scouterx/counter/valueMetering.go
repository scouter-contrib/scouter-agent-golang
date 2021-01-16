package counter

import (
	"sync"
)

type ValueMetering struct {
	sync.Mutex
	metering *Metering
}

type ValueMetric struct {
	Sum float64
	Avg float64
}

type ValueBucket struct {
	value float64
	count int32
}

func NewValueMeter() *ValueMetering {
	return &ValueMetering{
		metering: NewMetering(
			func() interface{} {
				return &ValueBucket{}
			},
			func(b interface{}) {
				sb := b.(*ValueBucket)
				sb.value = 0
				sb.count = 0
			},
		),
	}
}

func (g *ValueMetering) Add(value float64) {
	g.Lock()
	defer g.Unlock()
	b := g.metering.GetCurrentBucket().(*ValueBucket)
	b.count++
	b.value += value
}

func (g *ValueMetering) GetAllCounter(period int) *ValueMetric {
	var sum float64
	var count int32
	var avg float64

	period = g.metering.SearchOnHandler(period, func(b interface{}) {
		vb := b.(*ValueBucket)
		sum += vb.value
		count += vb.count
	})
	if count == 0 {
		avg = 0
	} else {
		avg = sum / float64(count)
	}

	return &ValueMetric{
		Sum: sum,
		Avg: avg,
	}
}
