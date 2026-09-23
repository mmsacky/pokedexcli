package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cachedData map[string]cacheEntry
	interval   time.Duration
	mu         sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	data      []byte
}

func NewCache(interval time.Duration) *Cache {

	newCache := &Cache{
		cachedData: make(map[string]cacheEntry),
		interval:   interval,
	}

	go newCache.reapLoop()

	return newCache

}

func (c *Cache) Add(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cachedData[key] = cacheEntry{
		createdAt: time.Now(),
		data:      data,
	}

}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cachedItem, ok := c.cachedData[key]

	if ok {
		return cachedItem.data, true
	}
	return nil, false

}

func (c *Cache) reapLoop() {

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {

		c.mu.Lock()
		for key, value := range c.cachedData {
			if time.Since(value.createdAt) > c.interval {
				// fmt.Println("a cache record was deleted")
				delete(c.cachedData, key)
			} else {
			}
		}
		c.mu.Unlock()
	}
}
