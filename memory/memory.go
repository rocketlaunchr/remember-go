package memory

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rocketlaunchr/remember-go"
)

// MemoryStore is used to create an in-memory cache.
type MemoryStore struct {
	cache *cache.Cache
}

// NewMemoryStore creates an in-memory cache where the expired items
// are deleted based on the cleanupInterval duration.
func NewMemoryStore(cleanupInterval time.Duration) *MemoryStore {
	return &MemoryStore{
		cache: cache.New(cache.NoExpiration, cleanupInterval),
	}
}

// NewMemoryStoreFrom creates an in-memory cache directly from a *cache.Cache object.
func NewMemoryStoreFrom(cache *cache.Cache) *MemoryStore {
	return &MemoryStore{
		cache: cache,
	}
}

// Conn does nothing for this storage driver.
func (c *MemoryStore) Conn(ctx context.Context) (remember.Cacher, error) {
	return c, nil
}

// StorePointer sets whether a storage driver requires itemToStore to be
// stored as a pointer or as a concrete value.
func (c *MemoryStore) StorePointer() bool {
	return false
}

// Get returns a value from the cache if the key exists.
func (c *MemoryStore) Get(key string) (_ interface{}, found bool, _ error) {
	item, found := c.cache.Get(key)
	return item, found, nil
}

// Set sets a item into the cache for a particular key.
func (c *MemoryStore) Set(key string, expiration time.Duration, itemToStore interface{}) error {
	c.cache.Set(key, itemToStore, expiration)
	return nil
}

// Close returns the connection back to the pool for storage drivers that utilize a pool.
// For this driver, it does nothing.
func (c *MemoryStore) Close() {}

// Forget clears the value from the cache for the particular key.
func (c *MemoryStore) Forget(key string) error {
	c.cache.Delete(key)
	return nil
}

// ForgetAll clears all values from the cache.
func (c *MemoryStore) ForgetAll() error {
	c.cache.Flush()
	return nil
}
