package pokecache

import (
	"time"
	"sync"
	"fmt"
)

const ReapCacheTime = time.Second * 5

type Cache struct {
	entries map[string]cacheEntry
	mu sync.RWMutex
	interval time.Duration
	done chan struct{}
	stopOnce sync.Once
}
type cacheEntry struct {
	createdAt time.Time
	val []byte
}


func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
		interval: interval,
		done: make(chan struct{}),
	}
	go c.reapLoop()
	return c
}


func (c *Cache) Add(key string, val []byte) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if len(val) == 0 {
		return fmt.Errorf("val cannot be empty")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	}

	return nil
}


func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, found := c.entries[key]
	return entry.val, found
}


func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.reap()
		case <-c.done:
			return
		}
	}
}


func (c *Cache) Stop() {
	c.stopOnce.Do(func() {
		close(c.done)
	})
}


func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.entries {
		if time.Since(c.entries[key].createdAt) > c.interval {
			delete(c.entries, key)
		}
	}
}
