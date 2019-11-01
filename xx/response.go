// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"bufio"
	"net"
	"net/http"
)

type response struct {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	statusCode int
}

func (r *response) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}

func (r *response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

func (r *response) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}