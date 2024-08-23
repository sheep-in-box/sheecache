package sheecache

import (
	"github.com/sheep-in-box/sheecache/lru"
	"sync"
)

// cache is a thread-safe cache that implements a least-recently-used (LRU) cache.
type cache struct {
	mu         sync.Mutex // protects lru
	lru        *lru.Cache // lru cache
	cacheBytes int64      // max memory bytes allowed
}

// add adds a value to the cache.
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil { // lazy initialization
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

// get looks up a key's value from the cache.
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
