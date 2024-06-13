package lockmap

import (
	"context"
	"sync"

	"thumburl-service/internal/pkg/logger"
)

var globalMutex sync.RWMutex
var lockMap = make(map[string]*sync.Mutex)

func Lock(ctx context.Context, key string) {
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

	logger.Infow(
		ctx,
		"devtool action locked",
		"key", key,
	)
}

func Unlock(ctx context.Context, key string) {
	globalMutex.RLock()

	if _, ok := lockMap[key]; ok {
		lockMap[key].Unlock()
	}

	globalMutex.RUnlock()

	logger.Infow(
		ctx,
		"devtool action unlocked",
		"key", key,
	)
}
