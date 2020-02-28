package tlru

import (
	"cashier/internal/basecache"
	"container/list"
)

func (c *Cache) GetMap() map[string]basecache.Item {
	c.baseCache.RWMutex.RLock()
	m := c.baseCache.Items
	c.baseCache.RWMutex.RUnlock()
	return m
}

func (c *Cache) GetList() *list.List {
	c.baseCache.RWMutex.RLock()
	l := c.ll
	c.baseCache.RWMutex.RUnlock()
	return l
}

func (c *Cache) GetOnEvicted() func(string, interface{}) {
	return c.baseCache.OnEvicted
}
