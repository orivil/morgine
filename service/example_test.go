// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package service_test

import (
	"fmt"
	"github.com/orivil/morgine/service"
	"unsafe"
)

type A struct {
	Name string
}

type B struct {
	Name string

	// dependent on struct "A"
	Dependence *A
}

// 相当于 A 的 "new()" 。该方法的用于初始化 A 对象。
var providerA service.Provider = func(c *service.Container) (interface{}, error) {

	return &A{Name: "struct A"}, nil
}

// 相当于 B 的 "new()" 方法。该方法用于初始化 B 对象。
var providerB service.Provider = func(c *service.Container) (interface{}, error) {

	// 获得 A 的单例对象
	a, err := c.Get(&providerA)
	if err != nil {
		return nil, err
	}

	// 总是获得新的 A 对象
	// a := c.GetNew(&providerA).(*A)

	// 注入依赖的对象。
	return &B{Name: "struct B", Dependence: a.(*A)}, nil
}

func ExampleContainer() {

	container1 := service.NewContainer(false)
	// 获取服务
	b := container1.MustGet(&providerB).(*B)
	fmt.Println(b.Name)

	// 测试是否注入依赖
	fmt.Println(b.Dependence.Name)

	container2 := service.NewContainer(true)

	// 测试单例模式
	b1 := container2.MustGet(&providerB).(*B)
	b2 := container2.MustGet(&providerB).(*B)
	fmt.Println(unsafe.Pointer(b1) == unsafe.Pointer(b2))

	// 测试工厂模式
	b3 := container2.MustGetNew(&providerB).(*B)
	fmt.Println(unsafe.Pointer(b2) == unsafe.Pointer(b3))

	// 测试两个不同实例所依赖的对象是否时同一个实例（取决于 providerB 的定义方式）
	fmt.Println(unsafe.Pointer(b1.Dependence) == unsafe.Pointer(b3.Dependence))
	// Output:
	// struct B
	// struct A
	// true
	// false
	// true
}
