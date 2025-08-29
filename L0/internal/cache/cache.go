package cache

import (
	"L0/internal/models"
	"container/list"
	"sync"
)

type OrderCache struct {
	items      map[string]*list.Element
	list       *list.List
	mu         sync.Mutex
	maxEntries int
}

type entry struct {
	key   string
	value models.Order
}

func New(maxEntries int) *OrderCache {
	if maxEntries <= 0 {
		maxEntries = 1000
	}
	return &OrderCache{
		items:      make(map[string]*list.Element),
		list:       list.New(),
		maxEntries: maxEntries,
	}
}

func (c *OrderCache) Set(key string, value models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.list.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return
	}

	elem := c.list.PushFront(&entry{key, value})
	c.items[key] = elem

	if c.list.Len() >= c.maxEntries {
		c.removeOldest()
	}
}

func (c *OrderCache) Get(key string) (models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.list.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return models.Order{}, false
}

func (c *OrderCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.removeElement(elem)
	}
}

func (c *OrderCache) Pop(key string) (models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		value := elem.Value.(*entry).value
		c.removeElement(elem)
		return value, true
	}
	return models.Order{}, false
}

func (c *OrderCache) removeOldest() {
	elem := c.list.Back()
	if elem != nil {
		c.removeElement(elem)
	}
}

func (c *OrderCache) removeElement(e *list.Element) {
	c.list.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
}
