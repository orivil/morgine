// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package grace_test

import (
	"context"
	"database/sql"
	"github.com/orivil/morgine/utils/grace"
	"log"
	"net/http"
)

func ExampleListenSignal() {
	var db sql.DB
	server := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}

	// listen "interrupt" and "kill" signal
	closed := grace.ListenSignal(func() error {
		return server.Shutdown(context.Background())
	})

	// register callback before server shutdown
	server.RegisterOnShutdown(func() {
		db.Close()
	})

	// err will never be nil
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// wait until handled shutdown error
	<-closed
}
