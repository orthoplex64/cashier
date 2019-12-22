package gocache

import (
	"container/list"
)

func (c *cache) GetMap() map[string]Item {
	c.mu.RLock()
	m := c.items
	c.mu.RUnlock()
	return m
}

func (c *cache) GetList() *list.List {
	c.mu.RLock()
	l := c.ll
	c.mu.RUnlock()
	return l
}

func (c *cache) GetOnEvicted() func(string, interface{}) {
	return c.onEvicted
}
