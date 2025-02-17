package cache

import (
	"encoding/json"
	"log"
	"time"
)

type Cache[T, K comparable] struct {
	entries map[T]entry[K]
	expiry  float64
	maxSize int
}

type entry[T any] struct {
	value        T
	lastAccessed time.Time
}

func New[T, K comparable](expiry uint, maxSize uint) *Cache[T, K] {
	return &Cache[T, K]{
		entries: map[T]entry[K]{},
		expiry:  float64(expiry),
		maxSize: int(maxSize),
	}
}

func (c *Cache[T, K]) Insert(key T, value K) {
	if len(c.entries) >= c.maxSize {
		log.Print("Cache is full. Removing expired entries ...")

		c.removeOldEntries()
	}

	now := time.Now()
	v := entry[K]{
		value:        value,
		lastAccessed: now,
	}

	log.Printf("Inserted %v (%v) into cache", key, prettyPrint(value))

	c.entries[key] = v
}

func (c *Cache[T, K]) Contains(key T) (*K, bool) {
	value, exists := c.entries[key]

	if !exists {
		return nil, false
	}

	now := time.Now()

	if c.isStale(now, value.lastAccessed) {
		log.Printf("Stale data. Removing %v (%v) from cache", key, prettyPrint(value.value))
		delete(c.entries, key)
		return nil, false
	}

	value.lastAccessed = now

	log.Printf("Found %v (%v) in cache", key, prettyPrint(value.value))

	return &value.value, true
}

func (c *Cache[T, K]) removeOldEntries() {
	now := time.Now()

	for key, value := range c.entries {
		if !c.isStale(now, value.lastAccessed) {
			continue
		}

		log.Printf("Deleted %v (%v) from cache", key, prettyPrint(value.value))
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
