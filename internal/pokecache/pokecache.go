package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte // cached raw data
}

type Cache struct {
	cacheMap map[string]cacheEntry
	mux      sync.Mutex
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	tickChannel := ticker.C
	defer ticker.Stop()

	for true {
		reapTime := <-tickChannel
		c.mux.Lock()
		for key, entry := range c.cacheMap {
			if reapTime.After(entry.createdAt.Add(interval)) {
				delete(c.cacheMap, key)
			}
		}
		c.mux.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{cacheMap: map[string]cacheEntry{}, mux: sync.Mutex{}}
	go cache.reapLoop(interval)
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	c.cacheMap[key] = cacheEntry{createdAt: time.Now(), val: val}
	c.mux.Unlock()
	return
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	entry, ok := c.cacheMap[key]
	if !ok {
		return []byte{}, false
	}
	return entry.val, true
}
