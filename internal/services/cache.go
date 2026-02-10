package services

import (
	// Standard library imports for string manipulation, concurrency control, and time management
	"strings" // Used for prefix matching when deleting cache entries by prefix
	"sync"    // Provides RWMutex for thread-safe concurrent access to cache storage
	"time"    // Used for managing cache entry expiration times and cleanup intervals
)

// cacheItem represents a single cached value with its expiration timestamp.
// Each item stores arbitrary data and tracks when it should be considered stale.
type cacheItem struct {
	value     interface{} // The cached data (can be any type)
	expiresAt time.Time   // Absolute time when this cache entry becomes invalid
}

// Cache provides a thread-safe in-memory key-value store with automatic expiration.
// It supports storing arbitrary values with configurable TTL (time-to-live) and
// includes automatic cleanup of expired entries. The cache is safe for concurrent
// use by multiple goroutines through read-write mutex synchronization.
//
// The cache runs a background goroutine that periodically removes expired entries
// to prevent unbounded memory growth.
type Cache struct {
	mu    sync.RWMutex           // Read-write mutex for thread-safe concurrent access
	items map[string]cacheItem   // Internal storage mapping keys to cached items
}

// NewCache creates and initializes a new Cache instance with automatic cleanup.
// The function starts a background goroutine that runs the cleanup loop to
// periodically remove expired entries every 5 minutes.
//
// The cleanup goroutine runs for the lifetime of the application and helps
// prevent memory leaks by removing stale cache entries.
//
// Returns:
//   - *Cache: Initialized cache ready for concurrent read/write operations
func NewCache() *Cache {
	// Initialize the cache with an empty items map
	c := &Cache{
		items: make(map[string]cacheItem),
	}

	// Start the background cleanup goroutine to remove expired entries periodically.
	// This goroutine will continue running for the lifetime of the cache instance.
	go c.cleanupLoop()

	return c
}

// Get retrieves a value from the cache by its key, checking for expiration.
// This method is thread-safe and uses a read lock to allow concurrent reads.
//
// If the key doesn't exist or the cached item has expired, the function returns
// nil and false. Expired items are not automatically removed during Get calls;
// they're simply treated as cache misses and will be cleaned up during the
// periodic cleanup cycle.
//
// Parameters:
//   - key: The cache key to look up
//
// Returns:
//   - interface{}: The cached value if found and not expired, nil otherwise
//   - bool: true if the value was found and is still valid, false otherwise
func (c *Cache) Get(key string) (interface{}, bool) {
	// Use read lock to allow multiple concurrent readers
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Attempt to retrieve the cached item
	item, found := c.items[key]

	// Return cache miss if the key doesn't exist or the item has expired.
	// We check expiration on every Get to ensure stale data is never returned,
	// even if cleanup hasn't run yet.
	if !found || time.Now().After(item.expiresAt) {
		return nil, false
	}

	// Return the cached value for valid, non-expired entries
	return item.value, true
}

// Set stores a value in the cache with the specified time-to-live (TTL).
// This method is thread-safe and uses a write lock to ensure exclusive access
// during the update operation.
//
// If a value already exists for the given key, it will be replaced with the new
// value and expiration time. The expiration time is calculated as an absolute
// timestamp (current time + TTL) rather than a relative duration.
//
// Parameters:
//   - key: The cache key to store the value under
//   - value: The value to cache (can be any type)
//   - ttlSeconds: Time-to-live in seconds before the cached value expires
func (c *Cache) Set(key string, value interface{}, ttlSeconds int) {
	// Use write lock to ensure exclusive access during modification
	c.mu.Lock()
	defer c.mu.Unlock()

	// Store the cache item with an absolute expiration timestamp.
	// We use absolute time rather than duration to avoid recalculating
	// expiration on every Get call.
	c.items[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(time.Duration(ttlSeconds) * time.Second),
	}
}

// Delete removes a single cache entry by its exact key.
// This method is thread-safe and uses a write lock for exclusive access.
//
// If the key doesn't exist, the operation is a no-op (no error is returned).
// This is useful for cache invalidation when specific resources are updated
// or deleted.
//
// Parameters:
//   - key: The exact cache key to remove
func (c *Cache) Delete(key string) {
	// Use write lock to ensure exclusive access during deletion
	c.mu.Lock()
	defer c.mu.Unlock()

	// Delete the cache entry (no-op if key doesn't exist)
	delete(c.items, key)
}

// DeleteByPrefix removes all cache entries whose keys start with the given prefix.
// This method is thread-safe and uses a write lock for exclusive access.
//
// This is particularly useful for invalidating groups of related cache entries,
// such as all entries related to a specific resource type (e.g., "product:*" to
// clear all product-related cache entries when products are updated).
//
// The function iterates through all cache entries, so performance may degrade
// with very large caches. Consider using specific key patterns for better
// performance in production environments.
//
// Parameters:
//   - prefix: The string prefix to match against cache keys
func (c *Cache) DeleteByPrefix(prefix string) {
	// Use write lock to ensure exclusive access during bulk deletion
	c.mu.Lock()
	defer c.mu.Unlock()

	// Iterate through all cache entries and delete those matching the prefix.
	// This approach is simple but may be slow for large caches with many keys.
	for k := range c.items {
		if strings.HasPrefix(k, prefix) {
			delete(c.items, k)
		}
	}
}

// cleanupLoop runs as a background goroutine that periodically removes expired
// cache entries. This prevents memory leaks by ensuring stale data doesn't
// accumulate indefinitely in the cache.
//
// The cleanup process runs every 5 minutes and acquires a write lock to safely
// iterate through and remove expired items. The 5-minute interval balances
// memory efficiency with lock contention - running too frequently could impact
// performance in high-concurrency scenarios.
//
// This function is designed to run for the lifetime of the Cache instance and
// should only be called once during cache initialization.
func (c *Cache) cleanupLoop() {
	// Create a ticker that fires every 5 minutes for periodic cleanup
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop() // Ensure ticker resources are released if the loop exits

	// Run cleanup on each tick
	for range ticker.C {
		// Acquire write lock for safe iteration and deletion.
		// This blocks all reads and writes during cleanup, so we keep
		// the critical section as short as possible.
		c.mu.Lock()

		// Capture current time once to avoid multiple system calls
		now := time.Now()

		// Iterate through all cache entries and remove expired ones.
		// We check expiration by comparing the current time against
		// the stored expiresAt timestamp.
		for k, item := range c.items {
			if now.After(item.expiresAt) {
				delete(c.items, k)
			}
		}

		// Release the lock to allow normal cache operations to resume
		c.mu.Unlock()
	}
}
