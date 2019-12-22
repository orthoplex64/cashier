package gocache

import (
	"container/list"
	"time"
)

func (c *cache) Get(k string) (interface{}, bool) {
	if c.ll == nil {
		c.mu.RLock()
	} else {
		c.mu.Lock()
	}

	item, found := c.get(k)
	if c.ll == nil {
		c.mu.RUnlock()
	}
	if !found {
		return nil, false
	}

	if c.ll != nil {
		ek := &list.Element{Value: k}
		c.ll.MoveToFront(ek)
		c.mu.Unlock()
	}

	return item, true
}

func (c *cache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	c.mu.RUnlock()
	if !found {
		return nil, time.Time{}, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, time.Time{}, false
		}

		return item.Object, time.Unix(0, item.Expiration), true
	}

	return item.Object, time.Time{}, true
}

func (c *cache) get(k string) (interface{}, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Object, true
}

func (c *cache) ItemCount() int {
	c.mu.RLock()
	n := len(c.items)
	c.mu.RUnlock()
	return n
}
