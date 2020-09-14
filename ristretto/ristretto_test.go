// Copyright 2019-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package ristretto_test

import (
	"context"
	"testing"
	"time"

	rist "github.com/dgraph-io/ristretto"
	"github.com/rocketlaunchr/remember-go"
	"github.com/rocketlaunchr/remember-go/ristretto"
)

var cfg = &rist.Config{
	NumCounters: 1e7,     // number of keys to track frequency of (10M).
	MaxCost:     1 << 30, // maximum cost of cache (1GB).
	BufferItems: 64,      // number of keys per Get buffer.
}

func TestKeyBasicOperation(t *testing.T) {
	ctx := context.Background()
	var ms = ristretto.NewRistrettoStore(cfg)

	key := "key"
	exp := 10 * time.Minute

	slowQuery := func(ctx context.Context) (interface{}, error) {
		return "val", nil
	}

	actual, _, _ := remember.Cache(ctx, ms, key, exp, slowQuery)

	expected := "val"

	if actual.(string) != expected {
		t.Errorf("wrong val: expected: %v actual: %v", expected, actual)
	}
}

func TestFetchFromCacheAndDisableCache(t *testing.T) {
	ctx := context.Background()
	var ms = ristretto.NewRistrettoStore(cfg)

	key := "key"
	exp := 10 * time.Minute

	slowQuery := func(ctx context.Context) (interface{}, error) {
		return "val", nil
	}

	// warm up cache
	remember.Cache(ctx, ms, key, exp, slowQuery)

	// This time fetch from cache
	actual, _, _ := remember.Cache(ctx, ms, key, exp, slowQuery)

	expected := "val"

	if actual.(string) != expected {
		t.Errorf("wrong val: expected: %v actual: %v", expected, actual)
	}

	// Actual is now "val", Let's change it to "val2" and disable cache usage.

	slowQuery = func(ctx context.Context) (interface{}, error) {
		return "val2", nil
	}

	actual, _, _ = remember.Cache(ctx, ms, key, exp, slowQuery, remember.Options{DisableCacheUsage: true})

	expected = "val2"

	if actual.(string) != expected {
		t.Errorf("wrong val: expected: %v actual: %v", expected, actual)
	}
}
