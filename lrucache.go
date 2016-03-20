// Package lrucache provides a byte-size-limited implementation of
// httpcache.Cache that stores data in memory.
package lrucache

import (
	"container/list"
	"sync"
)

// LruCache is a thread-safe, in-memory httpcache.Cache that evicts the
// least recently used entries from memory when the MaxSize (in bytes) limit
// would be exceeded. Use the New constructor to create one.
type LruCache struct {
	MaxSize int64

	mu    sync.Mutex
	cache map[string]*list.Element
	lru   *list.List // Front is least-recent
	size  int64
}

// New creates an LruCache that will restrict itself to maxSize bytes of memory.
func New(maxSize int64) *LruCache {
	c := &LruCache{
		MaxSize: maxSize,
		lru:     list.New(),
		cache:   make(map[string]*list.Element),
	}

	return c
}

// Get returns the []byte representation of a cached response and a bool
// set to true if the key was found.
func (c *LruCache) Get(key string) (value []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if le, ok := c.cache[key]; ok {
		c.lru.MoveToBack(le)
		return le.Value.(*entry).value, true
	}
	return nil, false
}

// Set stores the []byte representation of a response for a given key.
func (c *LruCache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if le, ok := c.cache[key]; ok {
		c.lru.MoveToBack(le)
		e := le.Value.(*entry)
		c.size = c.size + int64(len(value)) - int64(len(e.value))
		e.value = value
	} else {
		e := &entry{key: key, value: value}
		c.cache[key] = c.lru.PushBack(e)
		c.size = c.size + e.size()
	}

	for c.size > c.MaxSize {
		le := c.lru.Front()
		if le == nil {
			panic("LruCache: non-zero size but empty lru")
		}
		c.deleteElement(le)
	}
}

// Delete removes the value associated with a key.
func (c *LruCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if le, ok := c.cache[key]; ok {
		c.deleteElement(le)
	}
}

func (c *LruCache) deleteElement(le *list.Element) {
	c.lru.Remove(le)
	e := le.Value.(*entry)
	delete(c.cache, e.key)
	c.size = c.size - e.size()
}

// Rough estimate of map + entry object + string + byte slice overheads in bytes.
const entryOverhead = 128

type entry struct {
	key   string
	value []byte
}

func (e *entry) size() int64 {
	return entryOverhead + int64(len(e.key)) + int64(len(e.value))
}
