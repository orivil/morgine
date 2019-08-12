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
	"sync"
	"time"
)

var closeFunc []func()

// 启动 http 服务
func Run(addr string) {
	var server = &http.Server{
		Addr:         addr,
		Handler:      DefaultServeMux,
		ErrorLog:     log.Emergency,
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}

	// 监听关闭信号
	closed := grace.ListenSignal(func() error {
		return server.Shutdown(context.Background())
	})

	// start http server
	log.Init.Printf("pid:[%d] listen on: http://localhost%s\n", os.Getpid(), addr)
	err := server.ListenAndServe()
	// handle server error
	if err != http.ErrServerClosed {
		log.Emergency.Fatalf("server closed: %s\n", err)
	}
	shutdown()
	// wait until server shutdown
	<-closed
}

// 启动 https 服务
func RunTLS(addr, cert, certKey string) {
	var server = &http.Server{
		Addr:         addr,
		Handler:      DefaultServeMux,
		ErrorLog:     log.Emergency,
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}

	// 监听关闭信号
	closed := grace.ListenSignal(func() error {
		return server.Shutdown(context.Background())
	})
	log.Init.Printf("pid:[%d] listen on: https://localhost%s\n", os.Getpid(), addr)
	err := server.ListenAndServeTLS(cert, certKey)
	// handle server error
	if err != http.ErrServerClosed {
		log.Emergency.Fatalf("server closed: %s\n", err)
	}
	shutdown()
	// wait until server shutdown
	<-closed
}

// 服务器退出监听后所执行的回调函数, 先注册将被后执行.
func AfterShutdown(call func()) {
	closeFunc = append(closeFunc, call)
}

var once = &sync.Once{}

func shutdown() {
	once.Do(func() {
		// 先注册将被后执行
		for ln := len(closeFunc); ln > 0; ln-- {
			closeFunc[ln-1]()
		}
	})
}
