// Copyright 2018-21 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package nocache

import (
	"context"
	"time"

	"github.com/rocketlaunchr/remember-go"
)

// NoCache is used for testing purposes.
type NoCache struct{}

// NewNoCache creates a NoCache struct.
func NewNoCache() *NoCache { return &NoCache{} }

// Conn will provide a "pretend" connection.
func (nc *NoCache) Conn(ctx context.Context) (remember.Cacher, error) { return nc, nil }

// StorePointer sets whether a storage driver requires itemToStore to be
// stored as a pointer or as a concrete value.
func (nc *NoCache) StorePointer() bool { return true }

// Get returns a value from the cache if the key exists.
func (nc *NoCache) Get(key string) (_ interface{}, found bool, _ error) { return nil, false, nil }

// Set sets a item into the cache for a particular key.
func (nc *NoCache) Set(key string, expiration time.Duration, itemToStore interface{}) error {
	return nil
}

// Close returns the connection back to the pool for storage drivers that utilize a pool.
func (nc *NoCache) Close() { return }

// Forget clears the value from the cache for the particular key.
func (nc *NoCache) Forget(key string) error { return nil }

// ForgetAll clears all values from the cache.
func (nc *NoCache) ForgetAll() error { return nil }
