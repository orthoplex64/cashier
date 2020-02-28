package tlru

import (
	"cashier/internal/basecache"
	"container/list"
	"fmt"
	"time"
)

func (c *Cache) Set(k string, x interface{}, d time.Duration) {
	c.baseCache.RWMutex.Lock()
	c.baseCache.Set(k, x, d)
	if c.ll != nil {
		ek := &list.Element{Value: k}
		if _, found := c.baseCache.Items[k]; found {
			c.ll.MoveToFront(ek)
		} else {
			c.ll.PushFront(ek)
		}

		if c.maxItems != 0 && uint(c.ll.Len()) > c.maxItems {
			c.DeleteOldest()
		}
	}
	c.baseCache.RWMutex.Unlock()
}

func (c *Cache) SetDefault(k string, x interface{}) {
	c.Set(k, x, basecache.DefaultExpiration)
}

func (c *Cache) Add(k string, x interface{}, d time.Duration) error {
	c.baseCache.RWMutex.Lock()
	if _, found := c.baseCache.Get(k); found {
		c.baseCache.RWMutex.Unlock()
		return fmt.Errorf("Item %s already exists", k)
	}

	c.baseCache.Set(k, x, d)
	c.baseCache.RWMutex.Unlock()
	return nil
}

func (c *Cache) Replace(k string, x interface{}, d time.Duration) error {
	c.baseCache.RWMutex.Lock()
	if _, found := c.baseCache.Get(k); !found {
		c.baseCache.RWMutex.Unlock()
		return fmt.Errorf("Item %s doesn't exist", k)
	}

	c.baseCache.Set(k, x, d)
	c.baseCache.RWMutex.Unlock()
	return nil
}
