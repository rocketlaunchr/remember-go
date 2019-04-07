package memory

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rocketlaunchr/remember-go"
)

// MemoryStore is used to create a in-memory cache.
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

// NewMemoryStoreFrom creates an in-memory cache directly a *cache.Cache
// object.
func NewMemoryStoreFrom(cache *cache.Cache) *MemoryStore {
	return &MemoryStore{
		cache: cache,
	}
}

func (c *MemoryStore) Conn(ctx context.Context) (remember.Cacher, error) {
	return c, nil
}

func (c *MemoryStore) StorePointer() bool {
	return false
}

func (c *MemoryStore) Get(key string) (interface{}, bool, error) {
	item, found := c.cache.Get(key)
	return item, found, nil
}

func (c *MemoryStore) Set(key string, expiration time.Duration, itemToStore interface{}) error {
	c.cache.Set(key, itemToStore, expiration)
	return nil
}

func (c *MemoryStore) Close() {}

func (c *MemoryStore) Forget(key string) error {
	c.cache.Delete(key)
	return nil
}

func (c *MemoryStore) ForgetAll() error {
	c.cache.Flush()
	return nil
}
