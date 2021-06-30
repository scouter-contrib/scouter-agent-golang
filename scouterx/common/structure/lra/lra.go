package lra

import (
	"container/list"
	"sync"
)

var lock sync.RWMutex
var itemExists = struct{}{}

type Cache struct {
	MaxEntries int
	OnEvicted func(key Key, value interface{})

	lst   *list.List
	table map[interface{}]*list.Element
}

type Key interface{}

type entry struct {
	key   Key
	value interface{}
}

func New(maxEntries int) *Cache {
	return &Cache{
		MaxEntries: maxEntries,
		lst:        list.New(),
		table:      make(map[interface{}]*list.Element),
	}
}

func (c *Cache) AddKey(key Key) {
	c.Add(key, itemExists)
}

func (c *Cache) Add(key Key, value interface{}) {
	lock.Lock()
	defer lock.Unlock()
	if ee, ok := c.table[key]; ok {
		ee.Value.(*entry).value = value
		return
	}
	ele := c.lst.PushFront(&entry{key, value})
	c.table[key] = ele
	if c.MaxEntries != 0 && c.lst.Len() > c.MaxEntries {
		c.removeOldest()
	}
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key Key) interface{} {
	lock.RLock()
	defer lock.RUnlock()
	if ele, contains := c.table[key]; contains {
		return ele.Value.(*entry).value
	}
	return nil
}

func (c *Cache) Contains(key interface{}) bool {
	lock.RLock()
	defer lock.RUnlock()
	if _, contains := c.table[key]; !contains {
		return false
	}
	return true
}

func (c *Cache) Empty() bool {
	lock.RLock()
	defer lock.RUnlock()
	return c.Size() == 0
}


func (c *Cache) Remove(key Key) {
	lock.Lock()
	defer lock.Unlock()
	if ele, hit := c.table[key]; hit {
		c.removeElement(ele)
	}
}

func (c *Cache) removeOldest() {
	ele := c.lst.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.lst.Remove(e)
	kv := e.Value.(*entry)
	delete(c.table, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

func (c *Cache) Size() int {
	lock.RLock()
	defer lock.RUnlock()
	return c.lst.Len()
}

func (c *Cache) Clear() {
	lock.Lock()
	defer lock.Unlock()
	if c.OnEvicted != nil {
		for _, e := range c.table {
			kv := e.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.lst = nil
	c.table = nil
}

func (c *Cache) GetValues() []interface{} {
	l := make([]interface{}, 0, int(float32(c.lst.Len()) * 1.2))
	lock.RLock()
	defer lock.RUnlock()
	for _, e := range c.table {
		lock.RUnlock()
		if e != nil {
			kv := e.Value.(*entry)
			if kv != nil && kv.value != nil {
				l = append(l, kv.value)
			}
		}
		lock.RLock()
	}
	return l
}
