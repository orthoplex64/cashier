package tlru

import (
	"cashier/internal/basecache"
	"time"
)

type kvp struct {
	key   string
	value interface{}
}

func (c *Cache) Delete(k string) {
	c.baseCache.RWMutex.Lock()
	v, evicted := c.baseCache.Delete(k)
	c.baseCache.RWMutex.Unlock()

	if evicted {
		c.baseCache.OnEvicted(k, v)
	}
}

func (c *Cache) DeleteOldest() {
	if c.ll == nil {
		return
	}

	b := c.ll.Back()
	if b == nil {
		return
	}

	k := b.Value.(string)
	c.baseCache.RWMutex.Lock()
	v, evicted := c.baseCache.Delete(k)
	c.baseCache.RWMutex.Unlock()

	if evicted {
		c.baseCache.OnEvicted(k, v)
	}
}

func (c *Cache) DeleteExpired() {
	var evictedItems []kvp
	now := time.Now().UnixNano()
	c.baseCache.RWMutex.Lock()
	for k, v := range c.baseCache.Items {
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := c.baseCache.Delete(k)
			if evicted {
				evictedItems = append(evictedItems, kvp{k, ov})
			}
		}
	}
	c.baseCache.RWMutex.Unlock()

	for _, v := range evictedItems {
		c.baseCache.OnEvicted(v.key, v.value)
	}
}

func (c *Cache) Flush() {
	c.baseCache.RWMutex.Lock()

	var toEvict map[string]basecache.Item

	if c.baseCache.OnEvicted != nil {
		toEvict = map[string]basecache.Item{}
		for k, v := range c.baseCache.Items {
			toEvict[k] = v
		}
	}

	c.baseCache.Items = map[string]basecache.Item{}
	c.baseCache.RWMutex.Unlock()

	// TODO :: Look at this again
	if len(toEvict) != 0 {
		for k, v := range toEvict {
			c.baseCache.OnEvicted(k, v.Object)
		}
	}
}
