// Copyright 2018-21 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package remember

import (
	"context"
	"time"
)

// Options is used to change caching behavior.
type Options struct {

	// DisableCacheUsage disables the cache.
	// It can be useful during debugging.
	DisableCacheUsage bool

	// UseFreshData will ignore content in the cache and always pull fresh data.
	// The pulled data will subsequently be saved in the cache.
	UseFreshData bool

	// Logger, when set, will log error and debug messages.
	Logger Logger

	// OnlyLogErrors, when set, will only log errors (but not debug messages).
	// For production, this should be set to true.
	OnlyLogErrors bool

	// GobRegister registers with the gob encoder the data type returned by the
	// SlowRetrieve function.
	// Some storage drivers may require this to be set.
	// Setting this to true will slightly impact concurrency performance.
	// It is usually better to set this to false, but register all structs
	// inside an init(). Otherwise you will encounter complaints from the gob package
	// if a Logger is provided.
	// See: https://golang.org/pkg/encoding/gob/#Register
	GobRegister bool
}

// SlowRetrieve obtains a result when the key is not found in the cache.
// It is usually (but not limited to) a query to a database with some additional
// processing of the returned data. The function must return a value that is compatible
// with the gob package for some storage drivers.
type SlowRetrieve func(ctx context.Context) (interface{}, error)

// Conner allows a storage driver to provide a connection from the pool
// in order to communicate with it.
type Conner interface {
	Conn(ctx context.Context) (Cacher, error)
}

// Cacher is the interface that all storage drivers must implement.
type Cacher interface {
	// StorePointer sets whether a storage driver requires itemToStore to be
	// stored as a pointer or as a concrete value.
	StorePointer() bool

	// Get returns a value from the cache if the key exists.
	Get(key string) (item interface{}, found bool, err error)

	// Set sets a item into the cache for a particular key.
	Set(key string, expiration time.Duration, itemToStore interface{}) error

	// Close returns the connection back to the pool for storage drivers that utilize a pool.
	Close()

	// Forget clears the value from the cache for the particular key.
	Forget(key string) error

	// ForgetAll clears all values from the cache.
	ForgetAll() error
}
