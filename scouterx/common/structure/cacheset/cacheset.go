package cacheset

import (
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"strconv"
	"sync"
)

var lock sync.RWMutex

type CacheSet struct {
	table    map[interface{}]struct{}
	ordering *singlylinkedlist.List
	maxSize  int
}

var itemExists = struct{}{}

func New(maxSize int) *CacheSet {
	set := &CacheSet{
		table:    make(map[interface{}]struct{}),
		ordering: singlylinkedlist.New(),
		maxSize:  maxSize,
	}
	return set
}

func (set *CacheSet) Add(item interface{}) {
	lock.Lock()
	defer lock.Unlock()
	var contains bool
	if _, contains = set.table[item]; !contains {
		set.removeExceeded()
		set.table[item] = itemExists
		set.ordering.Append(item)
	}
}

func (set *CacheSet) removeExceeded() {
	for set.ordering.Size() >= set.maxSize {
		item, exist := set.ordering.Get(0)
		if exist {
			set.ordering.Remove(0)
			delete(set.table, item)
		}
	}
}

func (set *CacheSet) Contains(item interface{}) bool {
	lock.RLock()
	defer lock.RUnlock()
	if _, contains := set.table[item]; !contains {
		return false
	}
	return true
}

func (set *CacheSet) Empty() bool {
	lock.Lock()
	defer lock.Unlock()
	return set.Size() == 0
}

func (set *CacheSet) Size() int {
	return set.ordering.Size()
}

func (set *CacheSet) Clear() {
	lock.Lock()
	defer lock.Unlock()
	set.table = make(map[interface{}]struct{})
	set.ordering.Clear()
}

func (set *CacheSet) String() string {
	return "CacheSet[" + strconv.Itoa(set.ordering.Size()) + "]"
}

