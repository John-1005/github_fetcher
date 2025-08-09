package githubcache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAddGet(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	defer cache.Stop()

	cache.Add("key1", []byte("value1"))
	val, ok := cache.Get("key1")
	require.True(t, ok, "expected to find key1")
	require.Equal(t, "value1", string(val))
}

func TestExpirationInGet(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)
	defer cache.Stop()

	cache.Add("key2", []byte("value2"))
	time.Sleep(100 * time.Millisecond)

	val, ok := cache.Get("key2")
	require.False(t, ok, "expected key2 to be expired")
	require.Nil(t, val)
}

func TestReapLoop(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)
	defer cache.Stop()

	cache.Add("key3", []byte("value3"))
	time.Sleep(100 * time.Millisecond)

	cache.mu.Lock()
	_, exists := cache.data["key3"]
	cache.mu.Unlock()

	require.False(t, exists, "expected key3 to be removed by reapLoop")
}

func TestCorruptedCacheFile(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "cache.json")

	// Write invalid JSON to the file
	err := os.WriteFile(cacheFile, []byte("not-json"), 0644)
	require.NoError(t, err)

	cache := &Cache{
		data:     make(map[string]cacheEntry),
		interval: 1 * time.Hour,
		stopCh:   make(chan struct{}),
		filePath: cacheFile,
	}

	file, err := os.ReadFile(cacheFile)
	require.NoError(t, err)

	err = json.Unmarshal(file, &cache.data)
	require.Error(t, err, "expected unmarshal to fail for corrupted file")
}

func TestOverwriteKey(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	defer cache.Stop()

	cache.Add("key5", []byte("oldvalue"))
	cache.Add("key5", []byte("newvalue"))

	val, ok := cache.Get("key5")
	require.True(t, ok, "expected to find key5")
	require.Equal(t, "newvalue", string(val))
}

func TestMultipleKeys(t *testing.T) {
	cache := NewCache(1 * time.Hour)
	defer cache.Stop()

	cache.Add("key6", []byte("value6"))
	cache.Add("key7", []byte("value7"))

	val1, ok1 := cache.Get("key6")
	val2, ok2 := cache.Get("key7")

	require.True(t, ok1, "expected to find key6")
	require.True(t, ok2, "expected to find key7")
	require.Equal(t, "value6", string(val1))
	require.Equal(t, "value7", string(val2))
}

