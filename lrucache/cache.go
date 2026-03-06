package lrucache

import (
	"container/list"
	"sync"
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

// Cache is a concurrency-safe fixed-capacity LRU cache.
// Most recently used items are kept at the front.
type Cache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	order    *list.List
	items    map[K]*list.Element
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		panic("lrucache: capacity must be > 0")
	}

	return &Cache[K, V]{
		capacity: capacity,
		order:    list.New(),
		items:    make(map[K]*list.Element, capacity),
	}
}

func (c *Cache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		elem.Value.(*entry[K, V]).value = value
		c.order.MoveToFront(elem)
		return
	}

	elem := c.order.PushFront(&entry[K, V]{
		key:   key,
		value: value,
	})
	c.items[key] = elem

	if c.order.Len() <= c.capacity {
		return
	}

	c.evictOldestLocked()
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}

	c.order.MoveToFront(elem)
	return elem.Value.(*entry[K, V]).value, true
}

func (c *Cache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return false
	}

	c.order.Remove(elem)
	delete(c.items, key)
	return true
}

func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}

func (c *Cache[K, V]) Capacity() int {
	return c.capacity
}

func (c *Cache[K, V]) evictOldestLocked() {
	oldest := c.order.Back()
	if oldest == nil {
		return
	}

	c.order.Remove(oldest)
	delete(c.items, oldest.Value.(*entry[K, V]).key)
}

// Usage
//
// cache := lrucache.New[string, int](2)
// cache.Put("a", 1)
// cache.Put("b", 2)
// cache.Get("a") // "a" becomes most recently used
// cache.Put("c", 3) // "b" is evicted
