// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/rocketlaunchr/remember-go"
	red "github.com/rocketlaunchr/remember-go/redis"
)

var ctx = context.Background()

func TestKeyBasicOperation(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	var rs = red.NewRedisStore(&redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", s.Addr())
		},
	})

	key := "key"
	exp := 10 * time.Minute

	slowQuery := func(ctx context.Context) (interface{}, error) {
		return "val", nil
	}

	actual, _, _ := remember.Cache(ctx, rs, key, exp, slowQuery)

	expected := "val"

	if actual.(string) != expected {
		t.Errorf("wrong val: expected: %v actual: %v", expected, actual)
	}
}

func TestFetchFromCacheAndDisableCache(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	var rs = red.NewRedisStore(&redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", s.Addr())
		},
	})

	key := "key"
	exp := 10 * time.Minute

	slowQuery := func(ctx context.Context) (interface{}, error) {
		return "val", nil
	}

	// warm up cache
	remember.Cache(ctx, rs, key, exp, slowQuery)

	// This time fetch from cache
	actual, _, _ := remember.Cache(ctx, rs, key, exp, slowQuery)

	expected := "val"

	if actual.(string) != expected {
		t.Errorf("wrong val: expected: %v actual: %v", expected, actual)
	}

	// Actual is now "val", Let's change it to "val2" and disable cache usage.

	slowQuery = func(ctx context.Context) (interface{}, error) {
		return "val2", nil
	}

	actual, _, _ = remember.Cache(ctx, rs, key, exp, slowQuery, remember.Options{DisableCacheUsage: true})

	expected = "val2"

	if actual.(string) != expected {
		t.Errorf("wrong val: expected: %v actual: %v", expected, actual)
	}
}
