package tlru

import (
	"cashier/internal/basecache"
	"container/list"
	"time"
)

type Cache struct {
	maxItems  int
	baseCache *basecache.BaseCache
	ll        *list.List
}

func New(maxItems int, defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]basecache.Item)
	c := &Cache{
		maxItems:  maxItems,
		baseCache: basecache.NewCache(defaultExpiration, cleanupInterval, items, nil),
		ll:        list.New(),
	}

	c = basecache.InitJanitor(c, c.baseCache, stopJanitor).(*Cache)
	return c
}

func stopJanitor(tc basecache.TimeAwareCache) {
	c := tc.(*Cache)
	c.baseCache.Janitor.Stop <- struct{}{}
}

func (c *Cache) OnEvicted(f func(string, interface{})) {
	c.baseCache.RWMutex.Lock()
	c.baseCache.OnEvicted = f
	c.baseCache.RWMutex.Unlock()
}
