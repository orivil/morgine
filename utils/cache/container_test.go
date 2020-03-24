// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cache_test

import (
	"github.com/orivil/morgine/utils/cache"
	"testing"
)

// 1045 ns/op
func BenchmarkContainer_Set(b *testing.B) {
	c := cache.NewContainer(3)
	for i:=0; i<b.N; i++ {
		c.Set(i, i, 0)
	}
	b.Log(c.Len())
}

// 1682 ns/op
func BenchmarkContainer_SetExpire(b *testing.B) {
	c := cache.NewContainer(3)
	for i:=0; i<b.N; i++ {
		c.Set(i, i, 10)
	}
	b.Log(c.Len())
}

var container = func() *cache.Container {
	c := cache.NewContainer(3)
	for i:=0; i< 4000000; i++ {
		c.Set(i, i, 10)
	}
	return c
}()

// 318 ns/op
func BenchmarkContainer_Get(b *testing.B) {
	for i:=0; i<b.N; i++ {
		container.Get(i)
	}
	b.Log(container.Len())
}