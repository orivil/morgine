// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

// Package storage implements the object storage rules
package storage

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	ServeHost     string
	Dir           string
	HeaderHandler func(header http.Header)
}

func NewLocalStorage(dir, serveHost string, corsHandler func(header http.Header)) (*LocalStorage, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &LocalStorage{Dir: dir, ServeHost: serveHost, HeaderHandler: corsHandler}, nil
}

func (l *LocalStorage) file(name string) string {
	return filepath.Join(l.Dir, name)
}

func (l *LocalStorage) IsExist(name string) (bool, error) {
	_, err := os.Stat(l.file(name))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (l *LocalStorage) Write(name string, data []byte) error {
	return ioutil.WriteFile(l.file(name), data, os.ModePerm)
}

func (l *LocalStorage) Read(name string) (data []byte, err error) {
	return ioutil.ReadFile(name)
}

func (l *LocalStorage) Remove(name string) error {
	err := os.Remove(l.file(name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (l *LocalStorage) GetServeUrl(name string) (url string, err error) {
	return "http://" + l.ServeHost + "/" + l.file(name), nil
}

func (l *LocalStorage) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if l.HeaderHandler != nil {
		l.HeaderHandler(writer.Header())
	}
	http.ServeFile(writer, request, strings.TrimPrefix(request.URL.Path, "/"))
}

func (l *LocalStorage) RoutePattern() string {
	route := l.Dir
	if !strings.HasPrefix(route, "/") {
		route = "/" + route
	}
	if !strings.HasSuffix(route, "/") {
		route = route + "/"
	}
	return route
}
