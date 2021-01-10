package conf

import "sync"

var observerLock sync.RWMutex
var observers = make(map[string]Runnable)

type Runnable interface {
	Run()
}

func AddToConfObserver(name string, r Runnable) {
	observerLock.Lock()
	defer observerLock.Unlock()

	observers[name] = r
}

func confChangeNotify() {
	observerLock.RLock()
	defer observerLock.RUnlock()

	for _, r := range observers {
		r.Run()
	}
}
