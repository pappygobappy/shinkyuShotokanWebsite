package cache

import (
	"sync"
	"time"
)

type cachedItem struct {
	value     interface{}
	expiresAt time.Time
}

// MemoryCache is a thread-safe in-memory cache with TTL support
type MemoryCache struct {
	data            sync.Map
	ttl             time.Duration
	cleanupInterval time.Duration
}

// New creates a new MemoryCache with specified TTL
func New(ttl time.Duration) *MemoryCache {
	c := &MemoryCache{
		ttl:             ttl,
		cleanupInterval: ttl / 2,
	}

	go c.startCleanup()

	return c
}

// Get retrieves a value from cache, returns (value, found)
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	item, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}

	cached := item.(*cachedItem)
	if time.Now().After(cached.expiresAt) {
		c.data.Delete(key)
		return nil, false
	}

	return cached.value, true
}

// Set stores a value in cache with TTL
func (c *MemoryCache) Set(key string, value interface{}) {
	item := &cachedItem{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}

	c.data.Store(key, item)
}

// Delete removes a key from cache
func (c *MemoryCache) Delete(key string) {
	c.data.Delete(key)
}

// startCleanup runs periodically to remove expired entries
func (c *MemoryCache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		c.data.Range(func(key, value interface{}) bool {
			item := value.(*cachedItem)
			if time.Now().After(item.expiresAt) {
				c.data.Delete(key)
			}
			return true
		})
	}
}

// Reset clears all cache entries
func (c *MemoryCache) Reset() {
	c.data.Range(func(key, value interface{}) bool {
		c.data.Delete(key)
		return true
	})
}
