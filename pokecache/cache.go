// cache package provide utility to cache resonse from api call to save resource and bandwith
package cache

import (
	"sync"
	"time"
)

var mu sync.RWMutex

type cacheEntery struct {
	createdAt time.Time
	val       []byte
}

type Cache map[string]cacheEntery

func NewCache(interval time.Duration) Cache {
	ticker := time.NewTicker(interval)
	cache := make(Cache, 0)
	go func() {
		for {

			<-ticker.C
			cache.reap(interval)
		}
	}()
	return cache
}
func (cache *Cache) Add(key string, data []byte) ([]byte, bool) {
	mu.Lock()
	defer mu.Unlock()
	value, ok := (*cache)[key]
	currentTime := time.Now()
	(*cache)[key] = cacheEntery{val: data, createdAt: currentTime}
	return value.val, ok

}
func (cache Cache) Get(key string) ([]byte, bool) {
	mu.RLock()
	defer mu.RUnlock()
	value, ok := cache[key]
	return value.val, ok

}
func (cache *Cache) reap(interval time.Duration) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	for key, entry := range *cache {
		if now.Sub(entry.createdAt) > interval {
			delete(*cache, key)
		}
	}
}
