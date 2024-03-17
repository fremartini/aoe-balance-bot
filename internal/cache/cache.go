package cache

import (
	"time"
)

type Cache[T, K comparable] struct {
	entries map[T]entry[K]
	expiry  float64
}

type entry[T any] struct {
	timestamp time.Time
	value     T
}

func New[T, K comparable](expiry uint) *Cache[T, K] {
	return &Cache[T, K]{
		entries: map[T]entry[K]{},
		expiry:  float64(expiry),
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

	if c.isExpired(time.Now(), value.timestamp) {
		delete(c.entries, key)
		return nil, false
	}

	return &value.value, true
}

func (c *Cache[T, K]) isExpired(now, timeStamp time.Time) bool {
	return now.Sub(timeStamp).Hours() > c.expiry
}
