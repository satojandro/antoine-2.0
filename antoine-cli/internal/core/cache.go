package core

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type CacheManager struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

func NewCacheManager() *CacheManager {
	cm := &CacheManager{
		items: make(map[string]CacheItem),
	}

	// Iniciar limpieza periódica
	go cm.cleanupExpired()

	return cm
}

func (c *CacheManager) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(duration).UnixNano()
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *CacheManager) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Value, true
}

func (c *CacheManager) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *CacheManager) Health() error {
	return nil
}

func (c *CacheManager) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]CacheItem)
	return nil
}

func (c *CacheManager) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now().UnixNano()
		for key, item := range c.items {
			if now > item.Expiration {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
