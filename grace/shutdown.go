// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package grace

import (
	"github.com/orivil/morgine/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var mu = sync.Mutex{}

var c chan struct{}

var shutdowns []func() error

func ListenSignal(onShutdown func() error) (closed <-chan struct{}) {
	mu.Lock()
	defer mu.Unlock()
	shutdowns = append(shutdowns, onShutdown)
	if c == nil {
		closeChan := make(chan struct{})
		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt, os.Kill, syscall.SIGTERM)
			<-sigint
			Cancel()
		}()
		c = closeChan
	}
	return c
}

func Cancel() {
	mu.Lock()
	defer mu.Unlock()
	log.Init.Println("Shutting down...")
	for _, shutdown := range shutdowns {
		if err := shutdown(); err != nil {
			log.Warning.Printf("Shutdown: %v\n", err)
		}
	}
	close(c)
}
