package githubcache

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type cacheEntry struct {
	Val       []byte
	CreatedAt time.Time
}

type Cache struct {
	data     map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
	stopCh   chan struct{}
	filePath string
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]cacheEntry)

	data, err := json.MarshalIndent(c.data, "", " ")
	if err == nil {
		os.WriteFile(c.filePath, data, 0644)
	}
}

func NewCache(interval time.Duration) *Cache {

	c := &Cache{
		data:     make(map[string]cacheEntry),
		interval: interval,
		stopCh:   make(chan struct{}),
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		fmt.Printf("Error getting user cache directory: %v\n", err)
	}

	appCacheDir := filepath.Join(cacheDir, "github_fetcher")

	err = os.MkdirAll(appCacheDir, 0755)
	if err != nil {
		fmt.Printf("Error creating application cache directory: %v\n", err)
	}

	fullPath := filepath.Join(appCacheDir, "cache.json")
	c.filePath = fullPath

	_, err = os.Stat(fullPath)
	if err == nil {
		file, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Printf("Failed to read file: %v\n", err)
		}

		err = json.Unmarshal(file, &c.data)
		if err != nil {
			fmt.Printf("Failed to unmarshal json: %v\n", err)
			c.data = make(map[string]cacheEntry)
		}

	}

	go c.reapLoop()

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheEntry{
		Val:       val,
		CreatedAt: time.Now(),
	}

	fmt.Printf("[verbose] Writing cache to %s\n", c.filePath)

	data, err := json.MarshalIndent(c.data, "", " ")
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err)
		return
	}

	err = os.WriteFile(c.filePath, data, 0644)
	if err != nil {
		log.Printf("Failed to write file: %v\n", err)
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, found := c.data[key]

	if !found {
		return nil, false
	}

	if time.Since(entry.CreatedAt) > c.interval {
		delete(c.data, key)
		return nil, false
	}

	return entry.Val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, entry := range c.data {
				if time.Since(entry.CreatedAt) > c.interval {
					delete(c.data, key)
				}
			}

			c.mu.Unlock()
		case <-c.stopCh:
			return
		}
	}
}

func (c *Cache) Stop() {
	if c.stopCh != nil {
		close(c.stopCh)
		c.stopCh = nil
	}

}

func (c *Cache) CacheFilePath() string {
	return c.filePath
}
