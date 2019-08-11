// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package storage

import "net/http"

// Interface is the object storage
type Storage interface {

	// IsExist if object is exist returns true
	IsExist(name string) (bool, error)

	// Write for saving the object to storage
	Write(name string, data []byte) error

	// Read for reading the object data form storage
	Read(name string) (data []byte, err error)

	// Remove for deleting the object from storage
	Remove(name string) error

	// GetServeUrl for getting the object http static file url
	GetServeUrl(name string) (url string, err error)
}

type HttpHandler interface {
	RoutePattern() string

	// HandleRequest for handling the http request and response the static file
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}
