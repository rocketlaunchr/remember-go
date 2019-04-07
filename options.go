package remember

import (
	"context"
	"time"
)

// CacheOptions is used to change caching behavior.
type Options struct {

	// DisableCacheUsage disables the cache.
	// It can be useful during debugging.
	DisableCacheUsage bool

	// UseFreshData will ignore content in the cache and always pull fresh data.
	// The pulled data will subsequently be saved in the cache.
	UseFreshData bool

	// Logger, when set, will turn on excessive logging.
	Logger Logger

	// GobRegister registers the struct returned by SlowRetrieve function with the gob encoder.
	// Some storage drivers may require this to be set.
	// Setting this to true will slightly impact concurrency performance.
	// It is usually better to set this to false, but register all structs
	// inside an init(). Otherwise you will encounter complaints from the gob package
	// if a Logger is provided.
	// See: https://golang.org/pkg/encoding/gob/#Register
	GobRegister bool
}

// SlowRetrieve obtains a value when the key is not found in the cache.
type SlowRetrieve func(ctx context.Context) (interface{}, error)

type Conner interface {
	Conn(ctx context.Context) (Cacher, error)
}

type Cacher interface {
	StorePointer() bool

	Get(key string) (item interface{}, found bool, err error)

	Set(key string, expiration time.Duration, itemToStore interface{}) error

	Close()
}
