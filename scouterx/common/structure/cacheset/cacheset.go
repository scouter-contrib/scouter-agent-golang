package cacheset

import (
	"github.com/emirpasic/gods/sets/linkedhashset"
	"strconv"
	"sync"
)

var lock sync.RWMutex

type CacheSet struct {
	table    *linkedhashset.Set
	maxSize  int
}

var itemExists = struct{}{}

func New(maxSize int) *CacheSet {
	s := &CacheSet{
		table:           linkedhashset.New(),
		maxSize:         maxSize,
	}
	return s
}

func (set *CacheSet) Add(item interface{}) {
	lock.Lock()
	defer lock.Unlock()
	if !set.table.Contains(item) {
		set.removeExceeded()
		set.table.Add(item)
	}
}

func (set *CacheSet) removeExceeded() {
	removalCount := set.table.Size() - set.maxSize
	if removalCount < 0 {
		return
	}
	var removals []interface{}
	iterator := set.table.Iterator()

	for i := removalCount; i >= 0; i-- {
		iterator.Next()
		removals = append(removals, iterator.Value())
	}
	for _, removal := range removals {
		set.table.Remove(removal)
	}
}

func (set *CacheSet) Contains(item interface{}) bool {
	lock.RLock()
	defer lock.RUnlock()
	return set.table.Contains(item)
}

func (set *CacheSet) Empty() bool {
	lock.Lock()
	defer lock.Unlock()
	return set.table. Size() == 0
}

func (set *CacheSet) Size() int {
	return set.table.Size()
}

func (set *CacheSet) Clear() {
	lock.Lock()
	defer lock.Unlock()
	set.table = linkedhashset.New()
}

func (set *CacheSet) String() string {
	return "CacheSet[" + strconv.Itoa(set.Size()) + "]"
}

