// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"github.com/orivil/morgine/proxy"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var passwd = "wen123456"

func main() {
	rand.Seed(time.Now().Unix())
	pxy := proxy.NewForwardProxy()
	//pxy.RequestHandler = func(client, server *http.Request) {
	//	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	//	server.Header.Set("X-Forwarded-For", ip)
	//}
	pxy.RequestHandler = nil
	//var handler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
	//	fmt.Println(request.URL)
	//	//p := request.Header.Get("passwd")
	//	//if p != passwd {
	//	//	_, err := writer.Write([]byte("wrong password"))
	//	//	if err != nil {
	//	//		panic(err)
	//	//	}
	//	//} else {
	//	pxy.ServeHTTP(writer, request)
	//	//}
	//}
	log.Fatal(http.ListenAndServe(":8081", pxy))
}
