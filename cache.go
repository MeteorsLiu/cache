package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

type Cache[T any] struct {
	sync.Mutex
	setFunc func() (T, error)
	lastSet atomic.Pointer[time.Time]

	err    error
	value  T
	expire time.Duration
}

func NewCache[T any](f func() (T, error), expire time.Duration) *Cache[T] {
	c := &Cache[T]{
		setFunc: f,
		expire:  expire,
	}
	c.set()
	return c
}

func (c *Cache[T]) set() {
	c.value, c.err = c.setFunc()
	// we still update the timestamp.
	now := time.Now().Add(c.expire)
	c.lastSet.Store(&now)
}

func (c *Cache[T]) Get() (T, error) {
	last := c.lastSet.Load()
	// time is not updated yet, wait the lock.
	// until is optimized in Unix-like system(Linux),
	// which is much faster than time.Now().Sub()
	if time.Until(*last) <= 0 {
		c.Lock()
		if last == c.lastSet.Load() {
			c.set()
		}
		c.Unlock()
	}
	return c.value, c.err
}
