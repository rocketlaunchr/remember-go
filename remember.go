package remember

import (
	"context"
	"encoding/gob"
	"fmt"
	"time"
)

// Cache is used to return a value from the cache if available. If unavailable, it will obtain the value by
// calling fn and then subsequently saving it into the cache.
// The behavior can be modified by providing an optional Options struct.
func Cache(ctx context.Context, c Conner, key string, expiration time.Duration, fn SlowRetrieve, options ...Options) (_ interface{}, found bool, _ error) {

	var (
		disableCache bool
		fresh        bool
		logger       Logger
		gobRegister  bool
	)

	if options != nil {
		disableCache = options[0].DisableCacheUsage
		fresh = options[0].UseFreshData
		logger = options[0].Logger
		gobRegister = options[0].GobRegister
	}

	// Check if cache has been disabled
	if disableCache {
		if logger != nil {
			logger.Log(logPatternBlue, "[cache disabled] Grabbing from SlowRetrieve key: "+key)
		}

		out, err := fn(ctx)
		if err != nil {
			if logger != nil {
				logger.Log(logPatternBlue, "[cache disabled] Grabbing (cache disabled) from SlowRetrieve key: "+key+" error: "+err.Error())
			}
			return nil, false, err
		}

		return out, false, nil
	}

	// Obtain cache connection
	cache, err := c.Conn(ctx)
	if err != nil {
		if logger != nil {
			logger.Log(logPatternRed, "could not obtain connection for cache")
		}
		return nil, false, err
	}
	defer cache.Close()

	var item interface{}

	if fresh {
		if logger != nil {
			logger.Log(logPatternBlue, "Grabbing (fresh) from SlowRetrieve key: "+key)
		}
		goto fresh
	}

	// Check if item exists
	item, found, err = cache.Get(key)
	if err != nil {
		// Error when attempting to fetch from cache
		if logger != nil {
			logger.Log(logPatternRed, "could not fetch from cache key: "+key+" error: "+err.Error())
		}
	}

	if found && err == nil {
		// Item exists in cache
		if logger != nil {
			logger.Log(logPatternBlue, "Found in Cache key: "+key)
		}

		return item, true, nil
	}

	if logger != nil {
		logger.Log(logPatternBlue, "Grabbing from SlowRetrieve key: "+key)
	}

fresh:
	// Item does not exist in cache so grab it from the fn
	itemToStore, err := fn(ctx)
	if err != nil {
		return nil, false, err
	}

	if gobRegister {
		func(itemToStore interface{}) {
			defer func() {
				err := recover()
				if err != nil {
					if logger != nil {
						logger.Log(logPatternRed, fmt.Sprintf("gob register: %v", err))
					}
				}
			}()
			gob.Register(itemToStore)
		}(itemToStore)
	}

	// Store item in Cache
	if cache.StorePointer() {
		err = cache.Set(key, expiration, &itemToStore)
	} else {
		err = cache.Set(key, expiration, itemToStore)
	}
	if err != nil {
		// Storage failed
		if logger != nil {
			logger.Log(logPatternRed, "Could not store item to memcache key: "+key+" "+err.Error()+" "+fmt.Sprintf("%+v", itemToStore))
		}
	}

	return itemToStore, false, nil
}
