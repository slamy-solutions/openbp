package pkcs

import (
	"sync"
	"time"
)

// slotCache represents a TTL slotCache that can store key/value pairs
type slotCache struct {
	sync.Mutex
	items      map[string]uint
	expiration map[string]time.Time
}

// NewSlotCache creates a new Cache with the specified TTL
func newSlotCache(ttl time.Duration) *slotCache {
	c := &slotCache{
		items:      make(map[string]uint),
		expiration: make(map[string]time.Time),
	}

	go func() {
		for {
			c.Lock()
			for key, expTime := range c.expiration {
				if time.Now().After(expTime) {
					delete(c.items, key)
					delete(c.expiration, key)
				}
			}
			c.Unlock()
			time.Sleep(ttl)
		}
	}()

	return c
}

// Set adds a key/value pair to the cache with the specified TTL
func (c *slotCache) Set(key string, value uint, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()

	c.items[key] = value
	c.expiration[key] = time.Now().Add(ttl)
}

// Get retrieves the value associated with the specified key from the cache
func (c *slotCache) Get(key string) (uint, bool) {
	c.Lock()
	defer c.Unlock()

	value, ok := c.items[key]
	if !ok {
		return 0, false
	}

	expTime, ok := c.expiration[key]
	if !ok {
		return 0, false
	}

	if time.Now().After(expTime) {
		delete(c.items, key)
		delete(c.expiration, key)
		return 0, false
	}

	return value, true
}
