package cashier_test

import (
	"cashier"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_Delete(t *testing.T) {
	tc := cashier.New(cashier.NoItemLimit, cashier.DefaultExpiration, 0)

	for k, v := range cacheItemTests {
		tc.Set(k, v.Object, time.Duration(v.Expiration))
	}

	var deletionTestsMap = map[string]cashier.Item{}
	for k, v := range cacheItemTests {
		deletionTestsMap[k] = v
	}

	for k := range cacheItemTests {
		tc.Delete(k)
		delete(deletionTestsMap, k)
		assert.Equal(t, tc.GetMap(), deletionTestsMap)
	}

	tc.Delete("keyNotInCache")
	assert.Equal(t, tc.GetMap(), map[string]cashier.Item{})
}
