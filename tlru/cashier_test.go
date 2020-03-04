package tlru_test

import (
	"cashier/internal/basecache"
	"cashier/tlru"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const MaxInt = int(^uint(0) >> 1)

var cacheItemTests = map[string]basecache.Item{
	"a": {
		Object:     1,
		Expiration: 0,
	},
	"b": {
		Object:     2.2,
		Expiration: 0,
	},
	"c": {
		Object:     "bar",
		Expiration: 0,
	},
	"d": {
		Object:     1,
		Expiration: 0,
	},
	"e": {
		Object:     2.2,
		Expiration: 0,
	},
	"f": {
		Object:     "bar",
		Expiration: 0,
	},
}

type TestStruct struct {
	Num      int
	Children []*TestStruct
}

func TestNew(t *testing.T) {
	tc := tlru.New(MaxInt, basecache.DefaultExpiration, 0)

	if assert.NotNil(t, tc) {
		assert.Equal(t, tc.GetMap(), map[string]basecache.Item{})
		//assert.Nil(t, tc.GetList())

		a, found := tc.Get("a")
		assert.False(t, found)
		assert.Nil(t, a)

		b, found := tc.Get("b")
		assert.False(t, found)
		assert.Nil(t, b)

		c, found := tc.Get("c")
		assert.False(t, found)
		assert.Nil(t, c)
	}
}

func TestCacheTimes(t *testing.T) {
	var found bool

	tc := tlru.New(MaxInt, 50*time.Millisecond, 1*time.Millisecond)
	tc.Set("a", 1, basecache.DefaultExpiration)
	tc.Set("b", 2, basecache.NoExpiration)
	tc.Set("c", 3, 20*time.Millisecond)
	tc.Set("d", 4, 70*time.Millisecond)

	<-time.After(25 * time.Millisecond)
	_, found = tc.Get("c")
	if found {
		t.Error("Found c when it should have been automatically deleted")
	}

	<-time.After(30 * time.Millisecond)
	_, found = tc.Get("a")
	if found {
		t.Error("Found a when it should have been automatically deleted")
	}

	_, found = tc.Get("b")
	if !found {
		t.Error("Did not find b even though it was set to never expire")
	}

	_, found = tc.Get("d")
	if !found {
		t.Error("Did not find d even though it was set to expire later than the default")
	}

	<-time.After(20 * time.Millisecond)
	_, found = tc.Get("d")
	if found {
		t.Error("Found d when it should have been automatically deleted (later than the default)")
	}
}

//func TestNewFrom(t *testing.T) {
//	m := map[string]basecache.Item{
//		"a": {
//			Object:     1,
//			Expiration: 0,
//		},
//		"b": {
//			Object:     2,
//			Expiration: 0,
//		},
//	}
//	tc := cashier.NewFrom(basecache.NoItemLimit, basecache.DefaultExpiration, 0, m)
//	a, found := tc.Get("a")
//	if !found {
//		t.Fatal("Did not find a")
//	}
//	if a.(int) != 1 {
//		t.Fatal("a is not 1")
//	}
//	b, found := tc.Get("b")
//	if !found {
//		t.Fatal("Did not find b")
//	}
//	if b.(int) != 2 {
//		t.Fatal("b is not 2")
//	}
//}

func TestStorePointerToStruct(t *testing.T) {
	tc := tlru.New(MaxInt, basecache.DefaultExpiration, 0)
	tc.Set("foo", &TestStruct{Num: 1}, basecache.DefaultExpiration)
	x, found := tc.Get("foo")
	if !found {
		t.Fatal("*TestStruct was not found for foo")
	}
	foo := x.(*TestStruct)
	foo.Num++

	y, found := tc.Get("foo")
	if !found {
		t.Fatal("*TestStruct was not found for foo (second time)")
	}
	bar := y.(*TestStruct)
	if bar.Num != 2 {
		t.Fatal("TestStruct.Num is not 2")
	}
}

func TestOnEvicted(t *testing.T) {
	tc := tlru.New(MaxInt, basecache.DefaultExpiration, 0)
	tc.Set("foo", 3, basecache.DefaultExpiration)
	if tc.GetOnEvicted() != nil {
		t.Fatal("tc.onEvicted is not nil")
	}
	works := false
	tc.OnEvicted(func(k string, v interface{}) {
		if k == "foo" && v.(int) == 3 {
			works = true
		}
		tc.Set("bar", 4, basecache.DefaultExpiration)
	})
	tc.Delete("foo")
	x, _ := tc.Get("bar")
	if !works {
		t.Error("works bool not true")
	}
	if x.(int) != 4 {
		t.Error("bar was not 4")
	}
}
