/**
 * Copyright 2019 orivil.com. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found at https://mit-license.org.
 */

package xx

import (
	"context"
	"github.com/orivil/morgine/log"
	"github.com/orivil/morgine/utils/grace"
	"net/http"
	"os"
	"time"
)

var Server = &http.Server{
	Handler:      DefaultServeMux,
	ErrorLog:     log.Emergency,
	WriteTimeout: 20 * time.Second,
	ReadTimeout:  20 * time.Second,
}

var closeFunc []func()

var beforeFunc []func()

// 开始监听端口
func Run() {
	for _, call := range beforeFunc {
		call()
	}

	// 监听关闭信号
	closed := grace.ListenSignal(func() error {
		return Server.Shutdown(context.Background())
	})
	var err error
	if Env.UseSSL {
		// start https server
		log.Init.Printf("pid:[%d] listen on: https://localhost%s\n", os.Getpid(), Env.HttpsPort)
		Server.Addr = Env.HttpsPort
		err = Server.ListenAndServeTLS(Env.SSLCertificate, Env.SSLCertificateKey)
	} else {
		// start http server
		log.Init.Printf("pid:[%d] listen on: http://localhost%s\n", os.Getpid(), Env.HttpPort)
		Server.Addr = Env.HttpPort
		err = Server.ListenAndServe()
	}
	// handle server error
	if err != http.ErrServerClosed {
		log.Emergency.Fatalf("Server: %s\n", err)
	}
	// 先注册将被后执行
	for ln := len(closeFunc); ln > 0; ln-- {
		closeFunc[ln-1]()
	}
	// wait until server shutdown error
	<-closed
}

// OnShutdown 是 http.Server.RegisterOnShutdown() 的别名, 经过注册的函数将在 Server 关闭的时候调用,
// 需要注意的是: 当 Server 关闭时, 通过该方法注册的函数会被同时执行, 没有先后顺序
func OnShutdown(call func()) {
	Server.RegisterOnShutdown(call)
}

// 服务器退出监听后所执行的回调函数, 先注册将被后执行.
func AfterShutdown(call func()) {
	closeFunc = append(closeFunc, call)
}

func BeforeRun(call func()) {
	beforeFunc = append(beforeFunc, call)
}
