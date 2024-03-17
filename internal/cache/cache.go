package cache

import (
	"aoe-bot/internal/logger"
	"encoding/json"
	"time"
)

type Cache[T, K comparable] struct {
	entries map[T]entry[K]
	expiry  float64
	logger  *logger.Logger
}

type entry[T any] struct {
	timestamp time.Time
	value     T
}

func New[T, K comparable](expiry uint, logger *logger.Logger) *Cache[T, K] {
	return &Cache[T, K]{
		entries: map[T]entry[K]{},
		expiry:  float64(expiry),
		logger:  logger,
	}
}

func (c *Cache[T, K]) Insert(key T, value K) {
	now := time.Now()

	v := entry[K]{
		timestamp: now,
		value:     value,
	}

	c.logger.Infof("Inserted %v (%v) into cache", key, prettyPrint(value))

	c.entries[key] = v
}

func (c *Cache[T, K]) Contains(key T) (*K, bool) {
	value, exists := c.entries[key]

	if !exists {
		return nil, false
	}

	if c.isExpired(time.Now(), value.timestamp) {
		c.logger.Infof("Deleted %v (%v) from cache", key, prettyPrint(value.value))
		delete(c.entries, key)
		return nil, false
	}

	c.logger.Infof("Found %v (%v) in cache", key, prettyPrint(value.value))

	return &value.value, true
}

func (c *Cache[T, K]) isExpired(now, timeStamp time.Time) bool {
	return now.Sub(timeStamp).Hours() > c.expiry
}

func prettyPrint(a any) string {
	b, _ := json.Marshal(a)

	return string(b)
}
