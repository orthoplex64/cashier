package cashier

import (
	"container/list"
	"fmt"
	"time"
)

func (c *cache) Set(k string, x interface{}, d time.Duration) {
	c.mu.Lock()
	c.set(k, x, d)
	c.mu.Unlock()
}

func (c *cache) SetDefault(k string, x interface{}) {
	c.Set(k, x, DefaultExpiration)
}

func (c *cache) Add(k string, x interface{}, d time.Duration) error {
	if c.ll != nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		if _, found := c.get(k); found {
			return fmt.Errorf("Item %s already exists", k)
		}

		c.set(k, x, d)
		return nil
	}

	c.mu.RLock()
	_, found := c.get(k)
	c.mu.RUnlock()
	if found {
		return fmt.Errorf("Item %s already exists", k)
	}

	c.mu.Lock()
	c.set(k, x, d)
	c.mu.Unlock()
	return nil
}

func (c *cache) Replace(k string, x interface{}, d time.Duration) error {
	c.mu.RLock()
	_, found := c.get(k)
	c.mu.RUnlock()
	if !found {
		return fmt.Errorf("Item %s doesn't exist", k)
	}

	c.mu.Lock()
	c.set(k, x, d)
	c.mu.Unlock()

	return nil
}

func (c *cache) set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	if c.ll != nil {
		ek := &list.Element{Value: k}
		if _, found := c.items[k]; found {
			c.ll.MoveToFront(ek)

			c.items[k] = Item{
				Object:     x,
				Expiration: e,
			}
			return
		}
		c.ll.PushFront(ek)

		if c.maxItems != 0 && c.ll.Len() > c.maxItems {
			c.DeleteOldest()
		}
	}

	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
}
