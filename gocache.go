package gocache

import (
	"container/list"
	"errors"
	"runtime"
	"sync"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
	NoItemLimit                     = -1
)

var (
	ErrFound    = errors.New("found")
	ErrNotFound = errors.New("not found")
	ErrNoLRU    = errors.New("no LRU")
)

type Cache struct {
	*cache
}

type cache struct {
	maxItems          int
	defaultExpiration time.Duration
	mu                sync.RWMutex
	items             map[string]Item
	ll                *list.List
	onEvicted         func(string, interface{})
	janitor           *janitor
}

type Item struct {
	Object     interface{}
	Expiration int64
}

func New(mi int, defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)
	return newCacheWithJanitor(mi, defaultExpiration, cleanupInterval, items)
}

func NewFrom(mi int, defaultExpiration, cleanupInterval time.Duration, items map[string]Item) *Cache {
	return newCacheWithJanitor(mi, defaultExpiration, cleanupInterval, items)
}

func newCacheWithJanitor(mi int, de time.Duration, ci time.Duration, m map[string]Item) *Cache {
	c := newCache(mi, de, m)
	C := &Cache{c}
	if ci > 0 {
		runJanitor(c, ci)
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}

func newCache(mi int, de time.Duration, m map[string]Item) *cache {
	if de == DefaultExpiration {
		de = NoExpiration
	}

	c := &cache{
		maxItems:          mi,
		defaultExpiration: de,
		items:             m,
	}

	if c.maxItems != NoItemLimit {
		c.ll = list.New()
	}

	return c
}

func (c *cache) OnEvicted(f func(string, interface{})) {
	c.mu.Lock()
	c.onEvicted = f
	c.mu.Unlock()
}
