package scanblock

import (
	"sync"
	"sync/atomic"
	"time"
)

const minWaitBetweenCacheCleans = 10 * time.Minute

// Cache holds entries about remote IPs and their status.
type Cache struct {
	lock        sync.RWMutex
	entries     map[string]*CacheEntry
	lastCleaned atomic.Int64
}

// CacheEntry holds status information about a remote IP.
type CacheEntry struct {
	TotalRequests atomic.Uint64
	ScanRequests  atomic.Uint64
	FirstSeen     atomic.Int64
	LastSeen      atomic.Int64
	Blocking      atomic.Bool
}

// NewCache creates and returns a new cache.
func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]*CacheEntry, 10000),
	}
}

// GetEntry returns an entry from the cache.
// Read-locks the cache.
func (c *Cache) GetEntry(key string) *CacheEntry {
	c.lock.RLock()
	defer c.lock.RUnlock()

	// Get entry and return.
	return c.entries[key]
}

// CreateEntry creates a new entry in the cache.
// Write-locks the cache.
func (c *Cache) CreateEntry(key string) *CacheEntry {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if there already is an entry.
	entry, ok := c.entries[key]
	if ok {
		return entry
	}

	// Create and return.
	entry = &CacheEntry{}
	c.entries[key] = entry
	return entry
}

// CleanEntries removes all cache entries that weren't touched for at least maxAge.
// Only executes if cache was not recently cleaned.
// Write-locks the cache.
func (c *Cache) CleanEntries(maxAge time.Duration) (removedEntries int) {
	// Check if cache was recently cleaned already.
	lastCleaned := c.lastCleaned.Load()
	if lastCleaned > time.Now().Add(-minWaitBetweenCacheCleans).Unix() {
		// Cache was recently cleaned, skip this time.
		return 0
	}
	// Check if we get to clean it.
	nowUnix := time.Now().Unix()
	if !c.lastCleaned.CompareAndSwap(lastCleaned, nowUnix) {
		// Someone else just started cleaning, abort!
		return 0
	}

	// Calculate clean threshold.
	cleanOlderThan := time.Now().Add(-maxAge).Unix()

	c.lock.Lock()
	defer c.lock.Unlock()

	// Remove all entries that are too old from the cache.
	for key, entry := range c.entries {
		if entry.LastSeen.Load() < cleanOlderThan {
			delete(c.entries, key)
			removedEntries++
		}
	}
	return removedEntries
}
