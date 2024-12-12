package common

import (
	"container/list"
	"sync"
	"time"
)

// CacheItem represents an item stored in the cache.
type CacheItem struct {
	Key        interface{}
	Value      interface{}
	Timestamp  time.Time
}

// Cache is a thread-safe LRU cache.
type Cache struct {
	capacity     int
	items        map[interface{}]*list.Element
	evictionList *list.List
	mutex        sync.Mutex
	ttl          time.Duration // Time-to-live for cache items
}

// NewCache initializes a new Cache.
func NewCache(capacity int, ttl time.Duration) *Cache {
	return &Cache{
		capacity:     capacity,
		items:        make(map[interface{}]*list.Element),
		evictionList: list.New(),
		ttl:          ttl,
	}
}

// Get retrieves a value from the cache.
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		item := element.Value.(*CacheItem)
		if c.isExpired(item) {
			c.removeElement(element)
			return nil, false
		}
		// Move to front
		c.evictionList.MoveToFront(element)
		return item.Value, true
	}
	return nil, false
}

// Set adds or updates a value in the cache.
func (c *Cache) Set(key, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		c.evictionList.MoveToFront(element)
		item := element.Value.(*CacheItem)
		item.Value = value
		item.Timestamp = time.Now()
	} else {
		if c.evictionList.Len() >= c.capacity {
			// Remove least recently used item
			c.evict()
		}
		item := &CacheItem{
			Key:       key,
			Value:     value,
			Timestamp: time.Now(),
		}
		element := c.evictionList.PushFront(item)
		c.items[key] = element
	}
}

// Delete removes a value from the cache.
func (c *Cache) Delete(key interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		c.removeElement(element)
	}
}

// Clear removes all items from the cache.
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[interface{}]*list.Element)
	c.evictionList.Init()
}

// isExpired checks if a cache item has expired.
func (c *Cache) isExpired(item *CacheItem) bool {
	if c.ttl == 0 {
		return false
	}
	return time.Since(item.Timestamp) > c.ttl
}

// evict removes the least recently used item from the cache.
func (c *Cache) evict() {
	element := c.evictionList.Back()
	if element != nil {
		c.removeElement(element)
	}
}

// removeElement removes an element from the cache.
func (c *Cache) removeElement(element *list.Element) {
	c.evictionList.Remove(element)
	item := element.Value.(*CacheItem)
	delete(c.items, item.Key)
}
