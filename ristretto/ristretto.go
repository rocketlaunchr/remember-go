// Copyright 2018-21 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package ristretto

import (
	"context"
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/rocketlaunchr/remember-go"
)

// NoExpiration is used to indicate that data should not expire from the cache.
const NoExpiration time.Duration = 0

// ErrItemDropped signifies that the item to store was not inserted into the cache.
//
// See: https://godoc.org/github.com/dgraph-io/ristretto#Cache.Set
var ErrItemDropped = errors.New("item dropped")

// RistrettoStore is used to create an in-memory ristretto cache.
//
// See: https://godoc.org/github.com/dgraph-io/ristretto
type RistrettoStore struct {
	Cache       *ristretto.Cache
	DefaultCost *int64
}

// NewRistrettoStore creates an in-memory ristretto cache.
//
// See: https://godoc.org/github.com/dgraph-io/ristretto#Config
func NewRistrettoStore(config *ristretto.Config, defaultCost ...int64) *RistrettoStore {
	cache, err := ristretto.NewCache(config)
	if err != nil {
		panic(err)
	}

	var dc *int64
	if len(defaultCost) > 0 {
		dc = &defaultCost[0]
	}

	return &RistrettoStore{
		Cache:       cache,
		DefaultCost: dc,
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
// The cost is always set to 1, unless over-ridden at creation.
//
// See: https://godoc.org/github.com/dgraph-io/ristretto#Cache.SetWithTTL
func (r *RistrettoStore) Set(key string, expiration time.Duration, itemToStore interface{}) error {
	var (
		stored bool
		cost   int64 = 1
	)
	if r.DefaultCost != nil {
		cost = *r.DefaultCost
	}
	if expiration == NoExpiration {
		stored = r.Cache.Set(key, itemToStore, cost)
	} else {
		stored = r.Cache.SetWithTTL(key, itemToStore, cost, expiration)
	}

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
