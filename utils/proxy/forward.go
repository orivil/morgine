// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type RequestHandler func(req *http.Request)
type ResponseHandler func(server *http.Response, client http.ResponseWriter)

type Forward struct {
	RequestHandler  RequestHandler
	ResponseHandler ResponseHandler
	RoundTripper    http.RoundTripper
}

func NewForwardProxy() *Forward {
	return &Forward{
		RoundTripper: http.DefaultTransport,
		RequestHandler: func(req *http.Request) {
			if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
				if prior, ok := req.Header["X-Forwarded-For"]; ok {
					clientIP = strings.Join(prior, ", ") + ", " + clientIP
				}
				req.Header.Set("X-Forwarded-For", clientIP)
			}
		},
	}
}

func (fwd *Forward) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL)
	if fwd.RequestHandler != nil {
		fwd.RequestHandler(req)
	}
	res, err := fwd.RoundTripper.RoundTrip(req)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), http.StatusBadGateway)
		//rw.WriteHeader(http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	if fwd.ResponseHandler != nil {
		fwd.ResponseHandler(res, rw)
	}
	rw.WriteHeader(res.StatusCode)
	_, err = io.Copy(rw, res.Body)
	if err != nil {
		panic(err)
	}
}
