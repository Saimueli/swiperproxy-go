package proxy

import (
	"net/http"
	"sync"
	"time"
)

type CacheEntry struct {
	Body      []byte
	Headers   http.Header
	timestamp time.Time
}

type Cache struct {
	mu         sync.RWMutex
	entries    map[string]CacheEntry
	ttl        time.Duration
	maxEntries int
}

func NewCache(ttl time.Duration, maxEntries int) *Cache {
	return &Cache{
		entries:    make(map[string]CacheEntry),
		ttl:        ttl,
		maxEntries: maxEntries,
	}
}

func (c *Cache) Get(key string) (CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok {
		return CacheEntry{}, false
	}
	if time.Since(entry.timestamp) > c.ttl {
		return CacheEntry{}, false
	}
	return entry, true
}

func (c *Cache) Set(key string, entry CacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) >= c.maxEntries {
		for k := range c.entries {
			delete(c.entries, k)
			break
		}
	}
	entry.timestamp = time.Now()
	c.entries[key] = entry
}