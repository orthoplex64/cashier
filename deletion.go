package gocache

import (
	"container/list"
	"time"
)

type kvp struct {
	key   string
	value interface{}
}

func (c *cache) Delete(k string) {
	c.mu.Lock()
	v, evicted := c.delete(k)
	c.mu.Unlock()

	if evicted {
		c.onEvicted(k, v)
	}
}

func (c *cache) DeleteOldest() {
	if c.ll == nil {
		return
	}

	b := c.ll.Back()
	if b == nil {
		return
	}

	k := b.Value.(string)
	c.mu.Lock()
	v, evicted := c.delete(k)
	c.mu.Unlock()

	if evicted {
		c.onEvicted(k, v)
	}
}

func (c *cache) DeleteExpired() {
	var evictedItems []kvp
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := c.delete(k)
			if evicted {
				evictedItems = append(evictedItems, kvp{k, ov})
			}
		}
	}
	c.mu.Unlock()

	for _, v := range evictedItems {
		c.onEvicted(v.key, v.value)
	}
}

func (c *cache) delete(k string) (interface{}, bool) {
	if c.ll != nil {
		ek := &list.Element{Value: k}
		c.ll.Remove(ek)
	}

	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}

	delete(c.items, k)
	return nil, false
}

func (c *cache) Flush() {
	c.mu.Lock()
	if c.onEvicted != nil {
		for k, v := range c.items {
			c.onEvicted(k, v.Object)
		}
	}

	c.items = map[string]Item{}
	c.mu.Unlock()
}
