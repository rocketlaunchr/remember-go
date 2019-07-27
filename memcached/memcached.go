package memcached

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemachedStore is used to create a memcached-backed cache.
type MemachedStore struct {
	client *memcache.Client
}

func NewMemachedStore(server ...string) *MemachedStore {
	return &MemachedStore{
		client: memcache.New(server...),
	}
}

func NewMemachedStoreFromSelector(ss memcache.ServerSelector) *MemachedStore {
	return &MemachedStore{
		client: memcache.NewFromSelector(ss),
	}
}

func (c *MemachedStore) Conn(ctx context.Context) (remember.Cacher, error) {
	return c, nil
}

func (c *MemachedStore) StorePointer() bool {
	return false // Not sure if this should be true or false. Try with both.
}

// Get retrieves a value from the cache. The key must be at most 250 bytes in length.
func (c *MemachedStore) Get(key string) (interface{}, bool, error) {

	item, err := c.client.Get(key)
	if err != nil {
		if err == memcached.ErrCacheMiss {
			return nil, false, nil
		}
		return nil, false, err
	}

	var output interface{}

	err = gob.NewDecoder(bytes.NewBuffer(item.Value)).Decode(&output)
	if err != nil {
		return nil, true, err // Could not decode cached data
	}

	return output, true, nil
}

// Set stores a value in the cache. The key must be at most 250 bytes in length.
func (c *MemachedStore) Set(key string, expiration time.Duration, itemToStore interface{}) error {

	var exp int32

	if expiration != 0 {
		exp = int32(time.Now().Add(expiration).Unix())
	}

	// Convert item to bytes
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(itemToStore)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:        key,
		Expiration: exp,
		Value:      b.Bytes(),
	}

	return c.client.Set(item)
}

func (c *MemachedStore) Close() {}

func (c *MemachedStore) Forget(key string) error {
	return c.client.Delete(key)
}

func (c *MemachedStore) ForgetAll() error {
	return c.client.DeleteAll()
}
