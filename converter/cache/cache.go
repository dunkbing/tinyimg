package cache

import (
	"sync"
)

// Cache is a simple in-memory cache with generic key and value types.
type Cache[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewCache creates a new Cache with generic key and value types.
func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		data: make(map[K]V),
	}
}

// Set sets a value in the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// Get gets a value from the cache.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

// Delete deletes a value from the cache.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
