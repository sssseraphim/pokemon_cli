package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entrys map[string]cacheEntry
	mux    *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		entrys: make(map[string]cacheEntry),
		mux:    &sync.Mutex{},
	}

	go c.readLoop(interval)

	return c
}

func (c *Cache) AddCache(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.entrys[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		val:       val,
	}
}

func (c *Cache) GetCache(key string) (val []byte, exists bool) {
	entry, ok := c.entrys[key]
	return entry.val, ok
}

func (c *Cache) readLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now().UTC(), interval)
	}
}

func (c *Cache) reap(timeNow time.Time, interval time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	for k, v := range c.entrys {
		if v.createdAt.Before(timeNow.Add(-interval)) {
			delete(c.entrys, k)
		}
	}
}
