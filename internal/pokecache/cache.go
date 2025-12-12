package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries      map[string]cacheEntry
	reapInterval time.Duration
	mutex        *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	entries := make(map[string]cacheEntry)
	cache := Cache{
		entries:      entries,
		reapInterval: interval,
		mutex:        &sync.Mutex{},
	}

	go cache.reapLoop()

	return cache
}

func (c *Cache) Add(key string, value []byte) {
	entry := cacheEntry{createdAt: time.Now(), val: value}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = entry
}

func (c *Cache) Get(key *string) ([]byte, bool) {
	if key == nil {
		return make([]byte, 0), false
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, ok := c.entries[*key]
	if ok {
		return entry.val, true
	}

	return make([]byte, 0), false
}

func (c *Cache) reapLoop() {
	for {
		c.mutex.Lock()
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) >= c.reapInterval {
				delete(c.entries, key)
			}
		}
		c.mutex.Unlock()

		time.Sleep(c.reapInterval)
	}
}
