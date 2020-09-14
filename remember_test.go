// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package remember_test

import (
	"context"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/rocketlaunchr/remember-go"
	"github.com/rocketlaunchr/remember-go/memory"
)

type aLogger struct{}

func (l aLogger) Log(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func TestKeyGen(t *testing.T) {

	actuals := []string{
		remember.CreateKey(false, "+", "", 1, 2, 3),
		remember.CreateKeyStruct(struct {
			Search      string
			notExported string
			Ignored     string `json:"-"`
			Omit        string `json:"xxx"`
			Page        int
		}{"search", "z", "y", "ppp", 1}),
		remember.Hash("crc32-hash"),
		remember.CreateKey(true, "", "", 1, 2, 3),
	}

	expected := []string{
		"1+2+3",
		`{"Page":1,"Search":"search","xxx":"ppp"}`,
		"40ffd476",
		`^github.com/rocketlaunchr/remember-go_test.TestKeyGen_.+?github.com/rocketlaunchr/remember-go/remember_test.go_\d+_1 2 3$`,
	}

	for i := range actuals {
		actual := actuals[i]

		if i == 3 {
			match, _ := regexp.MatchString(expected[i], actual)
			if !match {
				t.Errorf("wrong val: expected (regex): %v actual: %v", expected[i], actual)
			}
		} else {
			if actual != expected[i] {
				t.Errorf("wrong val: expected: %v actual: %v", expected[i], actual)
			}
		}
	}
}

func TestKeyBasicOperation(t *testing.T) {
	ctx := context.Background()
	var ms = memory.NewMemoryStore(10 * time.Minute)

	key := "key"
	exp := 10 * time.Minute

	slowQuery := func(ctx context.Context) (interface{}, error) {
		return "val", nil
	}

	actual, _, _ := remember.Cache(ctx, ms, key, exp, slowQuery, remember.Options{Logger: aLogger{}, GobRegister: true})

	expected := "val"

	if actual.(string) != expected {
		t.Errorf("wrong val: expected: %v actual: %v", expected, actual)
	}
}

func TestFetchFromCacheAndDisableCache(t *testing.T) {
	ctx := context.Background()
	var ms = memory.NewMemoryStore(10 * time.Minute)

	key := "key"
	exp := 10 * time.Minute

	slowQuery := func(ctx context.Context) (interface{}, error) {
		return "val", nil
	}

	// warm up cache
	remember.Cache(ctx, ms, key, exp, slowQuery, remember.Options{Logger: aLogger{}, GobRegister: true})

	// This time fetch from cache
	actual, _, _ := remember.Cache(ctx, ms, key, exp, slowQuery, remember.Options{Logger: aLogger{}, GobRegister: true})

	expected := "val"

	if actual.(string) != expected {
		t.Errorf("wrong val: expected: %v actual: %v", expected, actual)
	}

	// Actual is now "val", Let's change it to "val2" and disable cache usage.

	slowQuery = func(ctx context.Context) (interface{}, error) {
		return "val2", nil
	}

	actual, _, _ = remember.Cache(ctx, ms, key, exp, slowQuery, remember.Options{Logger: aLogger{}, DisableCacheUsage: true})

	expected = "val2"

	if actual.(string) != expected {
		t.Errorf("wrong val: expected: %v actual: %v", expected, actual)
	}
}
