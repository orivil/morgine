// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package storage

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/orivil/morgine/bundles/utils/ticker"
	"io/ioutil"
	"strings"
	"sync"
)

type OssStorage struct {
	*oss.Bucket
	urlMaxAge int64
	urls      map[string]*expiredUrl
	cdnHost   string
	ossHost   string
	mu        sync.RWMutex
}

func NewOssStorage(name string, corsRules []oss.CORSRule, urlMaxAge int64, cdnHost, ossHost string) (*OssStorage, error) {
	bucket, err := Bucket(name, corsRules)
	if err != nil {
		return nil, err
	}
	return &OssStorage{
		Bucket:    bucket,
		urls:      make(map[string]*expiredUrl, 100),
		urlMaxAge: urlMaxAge,
		cdnHost:   cdnHost,
		ossHost:   ossHost,
	}, nil
}

type expiredUrl struct {
	url       string
	expiredAt int64
}

func (os *OssStorage) IsExist(name string) (bool, error) {
	return os.Bucket.IsObjectExist(name)
}

func (os *OssStorage) Write(name string, data []byte) error {
	return os.Bucket.PutObject(name, bytes.NewBuffer(data))
}

func (os *OssStorage) Read(name string) (data []byte, err error) {
	buf, err := os.Bucket.GetObject(name)
	if err != nil {
		return nil, err
	}
	defer buf.Close()
	return ioutil.ReadAll(buf)
}

func (os *OssStorage) Remove(name string) error {
	return os.Bucket.DeleteObject(name)
}

func (os *OssStorage) GetServeUrl(name string) (url string, err error) {
	now := ticker.Now().Unix()
	os.mu.Lock()
	defer os.mu.Unlock()
	eUrl, ok := os.urls[name]
	if ok {
		if eUrl.expiredAt > now {
			return eUrl.url, nil
		}
	}
	url, err = os.Bucket.SignURL(name, oss.HTTPGet, os.urlMaxAge)
	if err != nil {
		return "", err
	} else {
		if os.cdnHost != "" && os.ossHost != "" {
			url = strings.Replace(url, os.ossHost, os.cdnHost, 1)
		}
		eUrl = &expiredUrl{url: url, expiredAt: now + os.urlMaxAge}
		os.urls[name] = eUrl
		return eUrl.url, nil
	}
}
