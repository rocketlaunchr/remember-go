// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package memcached

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/rocketlaunchr/remember-go"
)

// MemcachedStore is used to create a memcached-backed cache.
type MemcachedStore struct {
	client *memcache.Client
}

// NewMemcachedStore creates a memcached-backed cache.
func NewMemcachedStore(server ...string) *MemcachedStore {
	return &MemcachedStore{
		client: memcache.New(server...),
	}
}

// NewMemachedStoreFromSelector creates a memcached-backed cache.
func NewMemachedStoreFromSelector(ss memcache.ServerSelector) *MemcachedStore {
	return &MemcachedStore{
		client: memcache.NewFromSelector(ss),
	}
}

// Conn does nothing for this storage driver.
func (c *MemcachedStore) Conn(ctx context.Context) (remember.Cacher, error) {
	return c, nil
}

// StorePointer sets whether a storage driver requires itemToStore to be
// stored as a pointer or as a concrete value.
func (c *MemcachedStore) StorePointer() bool {
	return false // Not sure if this should be true or false. Try with both?
}

// Get retrieves a value from the cache. The key must be at most 250 bytes in length.
func (c *MemcachedStore) Get(key string) (_ interface{}, found bool, _ error) {

	item, err := c.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
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
func (c *MemcachedStore) Set(key string, expiration time.Duration, itemToStore interface{}) error {

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

// Close returns the connection back to the pool for storage drivers that utilize a pool.
func (c *MemcachedStore) Close() {}

// Forget clears the value from the cache for the particular key.
func (c *MemcachedStore) Forget(key string) error {
	return c.client.Delete(key)
}

// ForgetAll clears all values from the cache.
func (c *MemcachedStore) ForgetAll() error {
	return c.client.DeleteAll()
}
