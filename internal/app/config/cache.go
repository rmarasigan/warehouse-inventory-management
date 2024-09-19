package config

import "sync"

var cachedconfig *cache

type cache struct {
	data map[string]any
	sync.RWMutex
}

// NewCache initializes the cache only once.
func NewCache() {
	if cachedconfig == nil {
		cachedconfig = &cache{
			data: make(map[string]any),
		}
	}
}

// GetCache retrieves the value for the provided key from the cache.
func GetCache(key string) any {
	// Read lock for concurrent access.
	cachedconfig.RLock()
	defer cachedconfig.RUnlock()

	return cachedconfig.data[key]
}

// SetCache sets a value in the cache for the provided key.
func SetCache(key string, value any) {
	// Write lock for concurrent access.
	cachedconfig.Lock()
	defer cachedconfig.Unlock()

	cachedconfig.data[key] = value
}

// DeleteCache removes the value associated with the provided key from
// the cache.
func DeleteCache(key string) {
	// Write lock for concurrent access.
	cachedconfig.Lock()
	defer cachedconfig.Unlock()

	// delete is a no-op if the key is missing
	delete(cachedconfig.data, key)
}
