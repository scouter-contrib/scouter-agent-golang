package counter

import (
	"sync"
)

var (
	onceServiceMetering sync.Once
	serviceMetering     *ServiceMetering
)

type ServiceMetering struct {
	sync.Mutex
	metering *Metering
}

type ServiceCounter struct {
	Tps       float32
	Elapsed   int
	ErrorRate float32
}

type ServiceBucket struct {
	count int
	elapsed int
	error int
}

func GetServiceMeter() *ServiceMetering {
	onceServiceMetering.Do(func() {
		serviceMetering = &ServiceMetering{
			metering: NewMetering(
				func() interface{} {
					return &ServiceBucket{}
				},
				func(b interface{}) {
					sb := b.(*ServiceBucket)
					sb.count = 0
					sb.elapsed = 0
					sb.error = 0
				},
			),
		}
	})
	return serviceMetering
}


func (g *ServiceMetering) Add(elapsed int, err bool) {
	g.Lock()
	defer g.Unlock()
	if elapsed < 0 {elapsed = 0}
	b := g.metering.GetCurrentBucket().(*ServiceBucket)
	b.count++
	b.elapsed += elapsed
	if (err) {
		b.error++
	}
}

func (g *ServiceMetering) GetAllCounter(period int) *ServiceCounter {
	var countSum int
	var elapsedSum int
	var errorSum int

	period = g.metering.SearchOnHandler(period, func(b interface{}) {
		sb := b.(*ServiceBucket)
		countSum += sb.count
		elapsedSum += sb.elapsed
		errorSum += sb.error
	})

	var tps = float32(countSum) / float32(period)
	var elapsed int
	var errorRate float32

	if countSum != 0 {
		elapsed = elapsedSum / countSum
		errorRate = float32(errorSum) / float32(countSum) * 100.0
	}

	return &ServiceCounter{
		Tps:       tps,
		Elapsed:   elapsed,
		ErrorRate: errorRate,
	}
}
