package cachemap

import (
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"strconv"
	"sync"
)

var lock sync.RWMutex

type CacheMap struct {
	table    map[interface{}]interface{}
	ordering *singlylinkedlist.List
	orderPos map[interface{}]int
	maxSize  int
}

func New(maxSize int) *CacheMap {
	m := &CacheMap{
		table:    make(map[interface{}]interface{}),
		ordering: singlylinkedlist.New(),
		orderPos: make(map[interface{}]int),
		maxSize:  maxSize,
	}
	return m
}

func (m *CacheMap) Add(key interface{}, item interface{}) {
	lock.Lock()
	defer lock.Unlock()
	if _, contains := m.table[key]; !contains {
		m.removeExceeded()
		m.table[key] = item
		m.ordering.Append(key)
		m.orderPos[key] = m.ordering.Size() - 1
	}
}

func (m *CacheMap) Remove(key interface{}) {
	lock.Lock()
	defer lock.Unlock()
	delete(m.table, key)
	index, contains := m.orderPos[key]
	if contains {
		m.ordering.Remove(index)
		delete(m.orderPos, key)
	}
}

func (m *CacheMap) removeExceeded() {
	for m.ordering.Size() >= m.maxSize {
		key, exist := m.ordering.Get(0)
		if exist {
			m.ordering.Remove(0)
			delete(m.orderPos, key)
			delete(m.table, key)
		}
	}
}

func (m *CacheMap) Contains(key interface{}) bool {
	lock.RLock()
	defer lock.RUnlock()
	if _, contains := m.table[key]; !contains {
		return false
	}
	return true
}

func (m *CacheMap) Get(key interface{}) interface{} {
	lock.RLock()
	defer lock.RUnlock()
	return m.table[key]
}

func (m *CacheMap) Empty() bool {
	lock.Lock()
	defer lock.Unlock()
	return m.Size() == 0
}

func (m *CacheMap) Size() int {
	return m.ordering.Size()
}

func (m *CacheMap) Clear() {
	lock.Lock()
	defer lock.Unlock()
	m.table = make(map[interface{}]interface{})
	m.ordering.Clear()
	m.orderPos = make(map[interface{}]int)
}

func (m *CacheMap) String() string {
	return "CacheMap[" + strconv.Itoa(m.ordering.Size()) + "]"
}
