// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package accounts

import (
	"encoding/base64"
	"errors"
	"github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/bundles/utils/crypto"
	"github.com/orivil/morgine/xx"
	"strconv"
	"strings"
	"sync"
	"time"
)

var AuthorizationMaxAge = int64(3600) * 24 * 7 // 授权过期时间，单位：秒

var (
	ErrBasicAuthorizeFailed              = errors.New("basic authorize failed")
	ErrBasicAuthorizationExpired         = errors.New("basic authorization expired")
	ErrBasicAuthorizationFormatIncorrect = errors.New("basic authorization format incorrect")
)

type AdminSessions struct {
	data map[string]*model.Admin
	mu   sync.RWMutex
}

func (as *AdminSessions) Set(key string, admin *model.Admin) {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.data[key] = admin
}

func (as *AdminSessions) Get(key string) (admin *model.Admin) {
	if key != "" {
		as.mu.RLock()
		defer as.mu.RUnlock()
		return as.data[key]
	}
	return nil
}

func (as *AdminSessions) Del(key string) {
	as.mu.Lock()
	defer as.mu.Unlock()
	delete(as.data, key)
}

var AdminSessionContainer = &AdminSessions{
	data: make(map[string]*model.Admin, 10),
}

var loginContextKey = "login-admin"

func GetAdminFromSession(ctx *xx.Context) (*model.Admin, error) {
	admin, ok := ctx.Get(loginContextKey).(*model.Admin)
	if ok {
		return admin, nil
	}
	auth := ctx.Request.Header.Get("Authorization")
	if auth != "" {
		admin = AdminSessionContainer.Get(auth)
		if admin == nil {
			username, password, err := DecodeBasicAuthorization(auth)
			if err == nil {
				return nil, err
			} else {
				admin, _ = model.SignIn(username, password)
				if admin != nil {
					AdminSessionContainer.Set(username, admin)
				}
			}
		}
		if admin != nil {
			ctx.Set(loginContextKey, admin)
			return admin, nil
		}
	}
	return nil, ErrBasicAuthorizeFailed
}

func MustGetAdmin(ctx *xx.Context) *model.Admin {
	admin, _ := GetAdminFromSession(ctx)
	return admin
}

func EncodeBasicAuthorization(username, password string) (auth string, err error) {
	auth = username + ":" + password + ":" + strconv.FormatInt(time.Now().Unix()+AuthorizationMaxAge, 10)
	auth, err = crypto.EncryptString(auth)
	if err != nil {
		return "", err
	}
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth)), nil
}

const prefix = "Basic "

func DecodeBasicAuthorization(auth string) (username, password string, err error) {
	if !strings.HasPrefix(auth, prefix) {
		return "", "", ErrBasicAuthorizationFormatIncorrect
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	auth, err = crypto.DecryptString(string(c))
	cs := strings.Split(auth, ":")
	if len(cs) == 3 {
		expireTime, _ := strconv.ParseInt(cs[2], 10, 64)
		if time.Now().Unix()-expireTime >= 0 {
			return "", "", ErrBasicAuthorizationExpired
		}
		return cs[0], cs[1], nil
	}
	return "", "", ErrBasicAuthorizeFailed
}

func DelAdminFromSession(username string) {
	AdminSessionContainer.Del(username)
}
