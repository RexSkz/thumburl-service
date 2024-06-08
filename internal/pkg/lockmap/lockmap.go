package lockmap

import (
	"fmt"
	"sync"
)

var globalMutex sync.RWMutex
var lockMap = make(map[string]*sync.Mutex)

func Lock(key string) {
	globalMutex.RLock()

	if lock, ok := lockMap[key]; ok {
		globalMutex.RUnlock()
		lock.Lock()
	} else {
		globalMutex.RUnlock()
		globalMutex.Lock()
		lockMap[key] = &sync.Mutex{}
		lockMap[key].Lock()
		globalMutex.Unlock()
	}

	fmt.Printf("locked %s\n", key)
}

func Unlock(key string) {
	globalMutex.RLock()

	if _, ok := lockMap[key]; ok {
		lockMap[key].Unlock()
	}

	globalMutex.RUnlock()
	fmt.Printf("unlocked %s\n", key)
}
