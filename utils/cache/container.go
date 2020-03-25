// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cache

import (
	"sync"
	"time"
)

type Container struct {
	values expireValues
	locker sync.Mutex
}

type expireValues map[interface{}]*expireValue

type expireValue struct {
	expireAt *time.Time
	value interface{}
}

func NewContainer() *Container {
	return &Container{values: make(expireValues)}
}

func (c *Container) Flash(key interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()
	delete(c.values, key)
}

func (c *Container) Len() int {
	c.locker.Lock()
	defer c.locker.Unlock()
	return len(c.values)
}

func (c *Container) CheckAndDelExpires(checkNum int, now time.Time) (deleted int) {
	c.locker.Lock()
	defer c.locker.Unlock()
	for key, value := range c.values {
		if checkNum == 0 {
			break
		}
		if value.expireAt != nil && value.expireAt.After(now) {
			delete(c.values, key)
			deleted ++
		}
		deleted--
	}
	return deleted
}

// expireAt 为 nil 则不过期
func (c *Container) Set(key, value interface{}, expireAt *time.Time) {
	c.locker.Lock()
	defer c.locker.Unlock()
	vue := &expireValue {
		expireAt: nil,
		value:    value,
	}
	if expireAt != nil {
		vue.expireAt = expireAt
	}
	c.values[key] = vue
}

func (c *Container) Get(key interface{}) (value interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()
	vue := c.values[key]
	if vue != nil {
		if vue.expireAt != nil {
			if vue.expireAt.After(time.Now()) {
				return vue.value
			} else {
				delete(c.values, key)
				return nil
			}
		} else {
			return vue.value
		}
	} else {
		return nil
	}
}