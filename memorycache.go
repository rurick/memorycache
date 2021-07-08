// Copyright 2019 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

// this package provides to save any value in memory cache
// Important!
// Be careful to use this cache module in different processes when creating group of microservices
// the cache table is only valid in one process

// Usage:
// cache := New(10*time.Minute, 1*time.Hour)
// cache.Set("simple_key", "value", 1*time.Minute)
// cache.Set("simple_key2", "value2")
// ...
// v := cache.Get("simple_key")
// ...
// cache.Delete("simple_key")
//

package memorycache

import (
	"fmt"
	"sync"
	"time"
)

type (
	// Cache - cache storage
	Cache struct {
		sync.RWMutex
		items             map[string]Item
		defaultExpiration time.Duration
		cleanupInterval   time.Duration
	}

	// Item - cache item
	Item struct {
		Value      interface{}
		Expiration int64
		Created    time.Time
	}
)

var (
	singleTone sync.Once
)

// New - initializing a new memory cache
// defaultExpiration - time.Duration for cache life time
// cleanupInterval - time.Duration time interval for running garbage collector
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	items := make(map[string]Item)
	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.startGC()
	}

	return &cache
}

// Set save value to cache with key
// opt:
// if there is opt with type time.Duration it used as cache life time,
// else used defaultExpiration from New function
func (c *Cache) Set(key string, value interface{}, opt ...interface{}) {
	var duration time.Duration
	for _, o := range opt {
		switch o.(type) {
		case time.Duration:
			duration = o.(time.Duration)
		}
	}
	var expiration int64
	if duration <= 0 {
		duration = c.defaultExpiration
	}
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

}

// Get getting cached value by key
func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if item.Expiration > 0 {
		// cache expired
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}

	return item.Value, true
}

// Delete cache by key
// Return error if key not found
func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return fmt.Errorf("Key %s not found", key)
	}

	delete(c.items, key)
	return nil
}

// StartGC start Garbage Collection
func (c *Cache) startGC() {
	singleTone.Do(func() {
		go c.gc()
	})
}

// gc Garbage Collection
func (c *Cache) gc() {
	for {
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.deleteItems(keys)
		}
	}
}

// expiredKeys return key list which are expired
func (c *Cache) expiredKeys() (keys []string) {
	c.RLock()
	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}
	return
}

// deleteItems removes all the items which key in keys.
func (c *Cache) deleteItems(keys []string) {
	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		if _, found := c.items[k]; found {
			delete(c.items, k)
		}
	}
}
