package cache

import (
	"log/slog"
	"sync"
	"time"
)

type Cache[T any] interface {
	Get(key string) (T, bool)
	Set(key string, value T)
	Delete(key string)
}

type entry[T any] struct {
	value    T
	expireAt time.Time
}

type cache[T any] struct {
	data            map[string]entry[T]
	mutex           sync.RWMutex
	ttlDuration     time.Duration
	cleanupInterval time.Duration
}

func NewCache[T any](ttlDuration time.Duration) *cache[T] {
	c := &cache[T]{
		data:            make(map[string]entry[T]),
		ttlDuration:     ttlDuration,
		cleanupInterval: ttlDuration / 2,
	}
	c.startCleanupTimer()
	return c
}

func (c *cache[T]) Get(key string) (T, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.data[key]
	if !exists || time.Now().After(entry.expireAt) {
		var zero T
		return zero, false
	}
	return entry.value, true
}

func (c *cache[T]) Set(key string, value T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = entry[T]{
		value:    value,
		expireAt: time.Now().Add(c.ttlDuration),
	}
}

func (c *cache[T]) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

func (c *cache[T]) startCleanupTimer() {
    go func() {
        ticker := time.NewTicker(c.cleanupInterval)
        for range ticker.C {
			slog.Debug("Cleaning up expired entries in cache")
            c.cleanExpiredEntries()
        }
    }()
}

func (c *cache[T]) cleanExpiredEntries() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, entry := range c.data {
		if now.After(entry.expireAt) {
			delete(c.data, key)
		}
	}
}
