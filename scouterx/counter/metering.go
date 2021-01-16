package counter

import (
	"sync"
	"time"
)

func getTimestamp() int {
	return int(time.Now().Unix())
}

type Metering struct {
	sync.Mutex

	bucketSize int
	time int
	pos int
	table []interface{}

	clear func(interface{})
}

type Handler interface {
	process(interface{})
}

func NewMetering(create func() interface{}, clear func(interface{})) *Metering {
	bucketSize := 301
	metering := &Metering{}
	metering.bucketSize = bucketSize
	metering.time = getTimestamp()
	metering.pos = int(metering.time % bucketSize)
	metering.table = make([]interface{}, bucketSize)
	metering.clear = clear

	for i := 0; i < bucketSize; i++ {
		metering.table[i] = create()
	}
	return metering
}

func (m *Metering) GetCurrentBucket() interface{} {
	pos := m.getPosition();
	return m.table[pos];
}

func (m *Metering) getPosition() int{
	m.Lock()
	defer m.Unlock()
	curTime := getTimestamp();
	if curTime != m.time {
		for i := 0; i < curTime-m.time && i < m.bucketSize; i++ {
			if m.pos + 1 > m.bucketSize - 1 {
				m.pos = 0
			} else {
				m.pos = m.pos + 1
			}
			m.clear(m.table[m.pos])
		}
		m.time = curTime
		m.pos = m.time % m.bucketSize
	}
	return m.pos;
}

func (m *Metering) check(period int) int {
	if period >= m.bucketSize {
		period = m.bucketSize - 1
	}
	return period
}

func (m *Metering) stepBack(pos int) int {
	if pos == 0 {
		pos = m.bucketSize - 1
	} else {
		pos--
	}
	return pos
}

func (m *Metering) SearchOnHandler(period int, handler func(interface{})) int {
	period = m.check(period)
	pos := m.getPosition()

	for i := 0; i < period; i++ {
		handler(m.table[pos])
		pos=m.stepBack(pos)
	}
	return period
}

func (m *Metering) search(period int) interface{} {
	period = m.check(period)
	pos := m.getPosition()
	out := make([]interface{}, period)
	for i := 0; i < period; i++ {
		out[i] = m.table[pos]
		pos = m.stepBack(pos)
	}
	return out
}

