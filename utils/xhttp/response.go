// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xhttp

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	HttpHeader http.Header
	Body *bytes.Buffer
	Status int
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		HttpHeader: make(http.Header),
		Body:       new(bytes.Buffer),
		Status:     200,
	}
}

func (h *ResponseWriter) Header() http.Header {
	return h.HttpHeader
}

func (h *ResponseWriter) Write(data []byte) (int, error) {
	return h.Body.Write(data)
}

func (h *ResponseWriter) WriteHeader(statusCode int) {
	h.Status = statusCode
}