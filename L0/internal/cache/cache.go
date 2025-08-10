package cache

import (
	"L0/internal/models"
	"sync"
)

type OrderCache struct {
	items map[string]models.Order
	mu    sync.Mutex
}

func New() *OrderCache {
	return &OrderCache{
		items: make(map[string]models.Order),
	}
}

func (c *OrderCache) Set(key string, value models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = value
}

func (c *OrderCache) Get(key string) (models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, found := c.items[key]
	return value, found
}

func (c *OrderCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *OrderCache) Pop(key string) (models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, found := c.items[key]
	if found {
		delete(c.items, key)
	}
	return value, found
}
