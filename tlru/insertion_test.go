package tlru_test

import (
	"cashier/internal/basecache"
	"cashier/tlru"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_Set(t *testing.T) {
	tc := tlru.New(MaxUint, basecache.DefaultExpiration, 0)

	for k, v := range cacheItemTests {
		tc.Set(k, v.Object, time.Duration(v.Expiration))
	}

	assert.Equal(t, tc.GetMap(), cacheItemTests)
}

func TestCache_SetDefault(t *testing.T) {
	tc := tlru.New(MaxUint, basecache.DefaultExpiration, 0)

	for k, v := range cacheItemTests {
		tc.Set(k, v.Object, time.Duration(v.Expiration))
	}

	assert.Equal(t, tc.GetMap(), cacheItemTests)
}

func TestCache_Add(t *testing.T) {
	tc := tlru.New(MaxUint, basecache.DefaultExpiration, 0)

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
	tc := tlru.New(MaxUint, basecache.DefaultExpiration, 0)

	for k, v := range cacheItemTests {
		tc.Set(k, v.Object, time.Duration(v.Expiration))
	}

	replaceItemsTest := map[string]basecache.Item{
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

	err := tc.Replace("keyNotInCache", "thisShouldError", basecache.DefaultExpiration)
	assert.Error(t, err)
}
