// Copyright 2018-21 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rocketlaunchr/remember-go"
)

// RedisStore is used to create a redis-backed cache.
type RedisStore struct {
	Pool *redis.Pool
}

// NewRedisStore creates a redis-backed cache directly from a redis
// pool object.
func NewRedisStore(redisPool *redis.Pool) *RedisStore {
	return &RedisStore{
		Pool: redisPool,
	}
}

// Conn will provide a single redis connection from the redis connection pool.
func (c *RedisStore) Conn(ctx context.Context) (remember.Cacher, error) {

	conn, err := c.Pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	return &RedisConn{
		conn: conn,
	}, nil
}

// RedisConn represents a single connection to the redis pool.
type RedisConn struct {
	conn redis.Conn
}

// StorePointer sets whether a storage driver requires itemToStore to be
// stored as a pointer or as a concrete value.
func (c *RedisConn) StorePointer() bool {
	return true
}

// Get returns a value from the cache if the key exists.
func (c *RedisConn) Get(key string) (_ interface{}, found bool, _ error) {

	val, err := redis.Bytes(c.conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			// Key not found
			return nil, false, nil
		}
		return nil, false, err
	}

	var output interface{}

	err = gob.NewDecoder(bytes.NewBuffer(val)).Decode(&output)
	if err != nil {
		return nil, true, err // Could not decode cached data
	}

	return output, true, nil
}

// Set sets a item into the cache for a particular key.
func (c *RedisConn) Set(key string, expiration time.Duration, itemToStore interface{}) error {

	// Convert item to bytes
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(itemToStore)
	if err != nil {
		return err
	}

	_, err = c.conn.Do("SET", key, b.Bytes(), "EX", int(expiration.Seconds()))
	return err
}

// Close returns the connection back to the pool for storage drivers that utilize a pool.
func (c *RedisConn) Close() {
	c.conn.Close()
}

// Forget clears the value from the cache for the particular key.
func (c *RedisConn) Forget(key string) error {
	_, err := c.conn.Do("DEL", key)
	return err
}

// ForgetAll clears all values from the cache.
func (c *RedisConn) ForgetAll() error {
	_, err := c.conn.Do("FLUSHDB")
	return err
}
