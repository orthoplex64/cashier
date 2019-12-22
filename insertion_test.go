package gocache_test

import (
	"gocache"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_Set(t *testing.T) {
	tc := gocache.New(gocache.NoItemLimit, gocache.DefaultExpiration, 0)

	for k, v := range cacheItemTests {
		tc.Set(k, v.Object, time.Duration(v.Expiration))
	}

	assert.Equal(t, tc.GetMap(), cacheItemTests)
}

func TestCache_SetDefault(t *testing.T) {
	tc := gocache.New(gocache.NoItemLimit, gocache.DefaultExpiration, 0)

	for k, v := range cacheItemTests {
		tc.Set(k, v.Object, time.Duration(v.Expiration))
	}

	assert.Equal(t, tc.GetMap(), cacheItemTests)
}

func TestCache_Add(t *testing.T) {
	tc := gocache.New(gocache.NoItemLimit, gocache.DefaultExpiration, 0)

	for k, v := range cacheItemTests {
		err := tc.Add(k, v.Object, time.Duration(v.Expiration))
		assert.NoError(t, err)
	}

	assert.Equal(t, tc.GetMap(), cacheItemTests)

	testErrKey := "c"
	testErrItem := cacheItemTests[testErrKey]
	err := tc.Add(testErrKey, testErrItem.Object, time.Duration(testErrItem.Expiration))
	assert.Error(t, err)
}

func TestCache_Replace(t *testing.T) {
	tc := gocache.New(gocache.NoItemLimit, gocache.DefaultExpiration, 0)

	for k, v := range cacheItemTests {
		tc.Set(k, v.Object, time.Duration(v.Expiration))
	}

	replaceItemsTest := map[string]gocache.Item{
		"a": {
			Object:     2,
			Expiration: 0,
		},
		"b": {
			Object:     3.5,
			Expiration: 0,
		},
		"c": {
			Object:     "bar2",
			Expiration: 0,
		},
		"d": {
			Object:     2,
			Expiration: 0,
		},
		"e": {
			Object:     3.5,
			Expiration: 0,
		},
		"f": {
			Object:     "bar2",
			Expiration: 0,
		},
	}

	for k, v := range replaceItemsTest {
		err := tc.Replace(k, v.Object, time.Duration(v.Expiration))
		assert.NoError(t, err)
	}

	assert.Equal(t, tc.GetMap(), replaceItemsTest)

	err := tc.Replace("keyNotInCache", "thisShouldError", gocache.DefaultExpiration)
	assert.Error(t, err)
}
