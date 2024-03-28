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
	maxSize uint
}

type entry[T any] struct {
	value        T
	lastAccessed time.Time
}

func New[T, K comparable](expiry uint, maxSize uint, logger *logger.Logger) *Cache[T, K] {
	return &Cache[T, K]{
		entries: map[T]entry[K]{},
		expiry:  float64(expiry),
		logger:  logger,
		maxSize: maxSize,
	}
}

func (c *Cache[T, K]) Insert(key T, value K) {
	if len(c.entries) >= int(c.maxSize) {
		c.logger.Info("Cache is full. Removing expired entries ...")

		c.removeOldEntries()
	}

	now := time.Now()
	v := entry[K]{
		value:        value,
		lastAccessed: now,
	}

	c.logger.Infof("Inserted %v (%v) into cache", key, prettyPrint(value))

	c.entries[key] = v
}

func (c *Cache[T, K]) Contains(key T) (*K, bool) {
	value, exists := c.entries[key]

	if !exists {
		return nil, false
	}

	now := time.Now()

	if c.isStale(now, value.lastAccessed) {
		c.logger.Infof("Stale data. Removing %v (%v) from cache", key, prettyPrint(value.value))
		delete(c.entries, key)
		return nil, false
	}

	value.lastAccessed = now

	c.logger.Infof("Found %v (%v) in cache", key, prettyPrint(value.value))

	return &value.value, true
}

func (c *Cache[T, K]) removeOldEntries() {
	for key, value := range c.entries {
		if !c.isStale(time.Now(), value.lastAccessed) {
			continue
		}

		c.logger.Infof("Deleted %v (%v) from cache", key, prettyPrint(value.value))
		delete(c.entries, key)
	}
}

func (c *Cache[T, K]) isStale(now, timeStamp time.Time) bool {
	return now.Sub(timeStamp).Hours() > c.expiry
}

func prettyPrint(a any) string {
	b, _ := json.Marshal(a)

	return string(b)
}
