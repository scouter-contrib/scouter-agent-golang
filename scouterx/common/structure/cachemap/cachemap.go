package cachemap

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"strconv"
	"sync"
)

var lock sync.RWMutex

type CacheMap struct {
	table           *linkedhashmap.Map
	maxSize         int
}

func (m *CacheMap) GetValues() []interface{} {
	l := make([]interface{}, 0, int(float32(m.table.Size()) * 1.2))
	lock.RLock()
	defer lock.RUnlock()
	for _, v := range m.table.Values() {
		lock.RUnlock()
		if v != nil {
			l = append(l, v)
		}
		lock.RLock()
	}
	return l
}

func New(maxSize int) *CacheMap {
	m := &CacheMap{
		table:           linkedhashmap.New(),
		maxSize:         maxSize,
	}
	return m
}

func (m *CacheMap) Add(key interface{}, item interface{}) {
	lock.Lock()
	defer lock.Unlock()
	if _, contains := m.table.Get(key); !contains {
		m.removeExceeded()
		m.table.Put(key, item)
	}
}

func (m *CacheMap) Remove(key interface{}) {
	lock.Lock()
	defer lock.Unlock()
	m.table.Remove(key)
}

func (m *CacheMap) removeExceeded() {
	removalCount := m.table.Size() - m.maxSize
	if removalCount < 0 {
		return
	}
	var removals []interface{}
	iterator := m.table.Iterator()

	for i := removalCount; i >= 0; i-- {
		iterator.Next()
		removals = append(removals, iterator.Key())
	}
	for _, removal := range removals {
		m.table.Remove(removal)
	}
}

func (m *CacheMap) Contains(key interface{}) bool {
	lock.RLock()
	defer lock.RUnlock()
	if _, contains := m.table.Get(key); !contains {
		return false
	}
	return true
}

func (m *CacheMap) Get(key interface{}) interface{} {
	lock.RLock()
	defer lock.RUnlock()
	value, found := m.table.Get(key)
	if !found {
		return nil
	} else {
		return value
	}
}

func (m *CacheMap) Empty() bool {
	lock.Lock()
	defer lock.Unlock()
	return m.Size() == 0
}

func (m *CacheMap) Size() int {
	return m.table.Size()
}

func (m *CacheMap) Clear() {
	lock.Lock()
	defer lock.Unlock()
	m.table = linkedhashmap.New()
}

func (m *CacheMap) String() string {
	return "CacheMap[" + strconv.Itoa(m.table.Size()) + "]"
}
