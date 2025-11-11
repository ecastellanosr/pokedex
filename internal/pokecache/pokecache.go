package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct { //cache with entries and mutex for concurrency safety
	CacheEntries map[string]*cacheEntry
	mu           *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type NoCacheError struct {
	Url string
}

func (d NoCacheError) Error() string {
	return fmt.Sprintf("No cache with url: %v", d.Url)
}

func NewCache(interval time.Duration) *Cache { //create a new cache
	cache := Cache{
		CacheEntries: make(map[string]*cacheEntry),
		mu:           &sync.Mutex{},
	}
	go cache.reapLoop(interval) //concurrent cache cleaning
	return &cache
}

func (c *Cache) Add(key string, val []byte) { // add a cacheEntry to cache
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CacheEntries[key] = &cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	return
}

func (c *Cache) Get(key string) ([]byte, bool) { //Get a cacheEntry
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.CacheEntries[key]
	if !ok {
		return nil, ok
	}
	fmt.Println("its getting used")
	return entry.val, ok
}
func (c *Cache) reapLoop(interval time.Duration) { //clean cache that were created above interval time
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	deleteEntries := make(chan bool)
	go func() {
		time.Sleep(interval)
		deleteEntries <- true
	}()
	for {
		select {
		case <-deleteEntries:
			if len(c.CacheEntries) == 0 {
				continue
			}
			for key, entry := range c.CacheEntries {
				time := time.Now().Sub(entry.createdAt) //time between creation and now
				if time < interval {
					continue
				}
				c.deleteEntry(key)
			}
		case <-ticker.C:
			continue
		}
	}
}

func (c *Cache) deleteEntry(key string) { //delete an especific entry
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.CacheEntries, key)
	return
}
