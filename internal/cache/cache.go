package cache

import (
	"sync"
	"time"
)

// Item 缓存项
type Item struct {
	Value     interface{}
	ExpireAt  time.Time
}

// Cache 分析结果缓存
type Cache struct {
	items map[string]*Item
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewCache 创建缓存
func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]*Item),
		ttl:   ttl,
	}
	go c.autoCleanup()
	return c
}

// Set 设置缓存
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &Item{
		Value:    value,
		ExpireAt: time.Now().Add(c.ttl),
	}
}

// Get 获取缓存
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(item.ExpireAt) {
		return nil, false
	}
	return item.Value, true
}

// Delete 删除缓存
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Cleanup 清理过期缓存
func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, item := range c.items {
		if now.After(item.ExpireAt) {
			delete(c.items, k)
		}
	}
}

func (c *Cache) autoCleanup() {
	for {
		time.Sleep(c.ttl)
		c.Cleanup()
	}
}
