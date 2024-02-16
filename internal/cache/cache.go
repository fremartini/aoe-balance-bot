package cache

import (
	"time"
)

type Cache[T, K comparable] struct {
	entries map[T]entry[K]
}

type entry[T any] struct {
	timestamp time.Time
	value     T
}

func New[T, K comparable]() *Cache[T, K] {
	return &Cache[T, K]{
		entries: map[T]entry[K]{},
	}
}

func (c *Cache[T, K]) Insert(key T, value K) {
	now := time.Now()

	v := entry[K]{
		timestamp: now,
		value:     value,
	}

	c.entries[key] = v
}

func (c *Cache[T, K]) Contains(key T) (*K, bool) {
	value, exists := c.entries[key]

	if !exists {
		return nil, false
	}

	if isExpired(time.Now(), value.timestamp) {
		delete(c.entries, key)
		return nil, false
	}

	return &value.value, true
}

func isExpired(now, timeStamp time.Time) bool {
	return now.Sub(timeStamp).Hours() > 24
}
