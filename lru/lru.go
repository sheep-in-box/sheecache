package lru

import "container/list"

// Cache is an LRU cache, and is not thread-safe.
type Cache struct {
	maxBytes int64      // maximum memory bytes allowed
	nBytes   int64      // memory bytes currently used
	ll       *list.List // doubly linked list to store the cache items
	cache    map[string]*list.Element
	// callback function when an entry is purged, can be nil.
	OnEvicted func(key string, value Value) //某条记录被移除时的回调函数，可以为nil
}

// entry is a key-value pair.
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes.
type Value interface {
	Len() int
}

// New constructs a new Cache with the specified maximum bytes and eviction callback.
func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get retrieves the value for the specified key, returning nil if the key is not found.
// The boolean 'ok' is true if the key is found, and false otherwise.
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item.
func (c *Cache) RemoveOldest() {
	element := c.ll.Back()
	if element != nil {
		c.ll.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add inserts a key-value pair into the Cache, updating the existing pair if the key is already present.
func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		element := c.ll.PushFront(&entry{key, value})
		c.cache[key] = element
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// Len is the number of Cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
