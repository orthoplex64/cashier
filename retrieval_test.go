package gocache_test

import (
	"gocache"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_Get(t *testing.T) {
	tc := gocache.New(gocache.NoItemLimit, gocache.DefaultExpiration, 0)

	for k, v := range cacheItemTests {
		tc.Set(k, v.Object, time.Duration(v.Expiration))
	}

	for k, v := range cacheItemTests {
		item, found := tc.Get(k)
		assert.True(t, found)
		assert.Equal(t, item, v.Object)
	}

	item, found := tc.Get("keyNotInCache")
	assert.False(t, found)
	assert.Nil(t, item)
}

func TestCache_GetWithExpiration(t *testing.T) {
	tc := gocache.New(gocache.NoItemLimit, gocache.DefaultExpiration, 0)

	for k, v := range cacheItemTests {
		tc.Set(k, v.Object, time.Duration(v.Expiration))
	}

	for k, v := range cacheItemTests {
		item, et, found := tc.GetWithExpiration(k)
		assert.True(t, found)
		assert.Equal(t, item, v.Object)
		if v.Expiration > 0 {
			assert.Equal(t, et, time.Unix(0, v.Expiration))
		} else {
			assert.Equal(t, et, time.Time{})
		}
	}

	item, et, found := tc.GetWithExpiration("keyNotInCache")
	assert.False(t, found)
	assert.Nil(t, item)
	assert.Equal(t, et, time.Time{})
}
