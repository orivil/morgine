// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"fmt"
)

func main() {
	fmt.Println(2147483647 >> 10 )
	fmt.Println(2147483647 >> 20 )
	fmt.Println(2147483647 >> 21 )
	fmt.Println(2147483647 >> 22 )
	fmt.Println(2147483647 >> 30 )
	fmt.Println(1 << 10 )
	fmt.Println(1 << 20 )
	fmt.Println(1 << 21 )
}
