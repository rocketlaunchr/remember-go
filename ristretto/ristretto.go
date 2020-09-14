// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package ristretto

import (
	"context"
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/rocketlaunchr/remember-go"
)

// ErrItemDropped signifies that the item to store was not inserted into the cache.
//
// See: https://godoc.org/github.com/dgraph-io/ristretto#Cache.Set
var ErrItemDropped = errors.New("item dropped")

// RistrettoStore is used to create an in-memory ristretto cache.
//
// See: https://godoc.org/github.com/dgraph-io/ristretto
type RistrettoStore struct {
	Cache *ristretto.Cache
}

// NewRistrettoStore creates an in-memory ristretto cache.
//
// See: https://godoc.org/github.com/dgraph-io/ristretto#Config
func NewRistrettoStore(config *ristretto.Config) *RistrettoStore {
	cache, err := ristretto.NewCache(config)
	if err != nil {
		panic(err)
	}

	return &RistrettoStore{
		Cache: cache,
	}
}

// Conn does nothing for this storage driver.
func (r *RistrettoStore) Conn(ctx context.Context) (remember.Cacher, error) {
	return r, nil
}

// StorePointer sets whether a storage driver requires itemToStore to be
// stored as a pointer or as a concrete value.
func (r *RistrettoStore) StorePointer() bool {
	return false
}

// Get returns a value from the cache if the key exists.
// It is possible for nil to be returned while found is also true.
//
// See: https://godoc.org/github.com/dgraph-io/ristretto#Cache.Get
func (r *RistrettoStore) Get(key string) (_ interface{}, found bool, _ error) {
	item, found := r.Cache.Get(key)
	return item, found, nil
}

// Set sets a item into the cache for a particular key.
//
// See: https://godoc.org/github.com/dgraph-io/ristretto#Cache.SetWithTTL
func (r *RistrettoStore) Set(key string, expiration time.Duration, itemToStore interface{}) error {
	stored := r.Cache.SetWithTTL(key, itemToStore, 1, expiration)
	if stored {
		return nil
	}
	return ErrItemDropped
}

// Close returns the connection back to the pool for storage drivers that utilize a pool.
// For this driver, it does nothing.
func (r *RistrettoStore) Close() {}

// Forget clears the value from the cache for the particular key.
//
// See: https://godoc.org/github.com/dgraph-io/ristretto#Cache.Del
func (r *RistrettoStore) Forget(key string) error {
	r.Cache.Del(key)
	return nil
}

// ForgetAll clears all values from the cache.
// Note that this is not an atomic operation.
//
// See: https://godoc.org/github.com/dgraph-io/ristretto#Cache.Clear
func (r *RistrettoStore) ForgetAll() error {
	r.Cache.Clear()
	return nil
}
