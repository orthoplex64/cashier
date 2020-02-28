package tlru

import (
	"container/list"
	"time"
)

func (c *Cache) Get(k string) (interface{}, bool) {
	c.baseCache.RWMutex.Lock()

	item, found := c.baseCache.Get(k)
	if !found {
		c.baseCache.RWMutex.Unlock()
		return nil, false
	}

	ek := &list.Element{Value: k}
	c.ll.MoveToFront(ek)

	c.baseCache.RWMutex.Unlock()
	return item, true
}

func (c *Cache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	c.baseCache.RWMutex.Lock()
	item, found := c.baseCache.Items[k]

	if !found {
		c.baseCache.RWMutex.Unlock()
		return nil, time.Time{}, false
	}

	c.baseCache.RWMutex.Unlock()
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, time.Time{}, false
		}

		return item.Object, time.Unix(0, item.Expiration), true
	}

	return item.Object, time.Time{}, true
}

func (c *Cache) ItemCount() int {
	c.baseCache.RWMutex.RLock()
	n := len(c.baseCache.Items)
	c.baseCache.RWMutex.RUnlock()
	return n
}
