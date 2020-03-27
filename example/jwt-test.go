// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

// +build ignore

package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"net/http"
	"time"
)

const (
	k = "1885df74d00dbbe19274c6d955feeb5b"
)

func main() {
	//生成token
	//提供三种加密方式SigningMethodHS256（sha256）SigningMethodHS384,SigningMethodHS512
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(1 * time.Second).Unix(),
		Id:        "1",
	})
	fmt.Println(token)
	t, _ := token.SignedString([]byte(k))
	time.Sleep(1 * time.Second)
	//创建req验证token
	req, _ := http.NewRequest("GET", "test", nil)
	req.Header.Add("token", t)
	fmt.Println(req.Header)
	token2, err := request.ParseFromRequest(req, request.HeaderExtractor{"token"}, func(token *jwt.Token) (interface{}, error) {
		return []byte(k), nil
	}, request.WithClaims(&jwt.StandardClaims{}))

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(token2.Valid)
	sc := token2.Claims.(*jwt.StandardClaims)
	fmt.Println(sc)
}
