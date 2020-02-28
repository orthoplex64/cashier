package basecache

import (
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

type TimeAwareCache interface {
	DeleteExpired()
}

type Item struct {
	Object     interface{}
	Expiration int64
}

type BaseCache struct {
	defaultExpiration time.Duration
	RWMutex           sync.RWMutex
	Items             map[string]Item
	OnEvicted         func(string, interface{})
	Janitor           *janitor
}

func (c *BaseCache) SetJanitor() *janitor {
	return c.Janitor
}

func (c *BaseCache) GetJanitor(j *janitor) {
	c.Janitor = j
}

func NewCache(d time.Duration, ci time.Duration, m map[string]Item, f func(string, interface{})) *BaseCache {
	var j *janitor
	if ci > 0 {
		j = &janitor{
			Interval: ci,
			Stop:     make(chan bool),
		}
	}

	return &BaseCache{
		defaultExpiration: d,
		Items:             m,
		OnEvicted:         f,
		Janitor:           j,
	}
}

func InitJanitor(c TimeAwareCache, bc *BaseCache, f func(TimeAwareCache)) TimeAwareCache {
	if bc.Janitor == nil {
		return c
	}

	go bc.Janitor.Run(c)
	runtime.SetFinalizer(c, f)

	// TODO :: Do we need to return this or will passing the reference suffice?
	return c
}

func (c *BaseCache) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	c.Items[k] = Item{
		Object:     x,
		Expiration: e,
	}
}

func (c *BaseCache) Get(k string) (interface{}, bool) {
	item, found := c.Items[k]
	if !found {
		return nil, false
	}

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Object, true
}

func (c *BaseCache) Delete(k string) (interface{}, bool) {
	if c.OnEvicted != nil {
		if v, found := c.Items[k]; found {
			delete(c.Items, k)
			return v.Object, true
		}
	}

	delete(c.Items, k)
	return nil, false
}
