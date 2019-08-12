// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package accounts

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/bundles/utils/crypto"
	"github.com/orivil/morgine/xx"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

const AuthScheme = "Bearer "

var AuthorizationMaxAge = int64(3600) * 24 * 7 // 授权过期时间，单位：秒

var (
	ErrBasicAuthorizeFailed              = errors.New("basic authorize failed")
	ErrBasicAuthorizationExpired         = errors.New("basic authorization expired")
	ErrBasicAuthorizationFormatIncorrect = errors.New("basic authorization format incorrect")
)

type AdminSessions struct {
	data map[int]*model.Admin
	mu   sync.RWMutex
}

func (as *AdminSessions) Set(admin *model.Admin) {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.data[admin.ID] = admin
}

func (as *AdminSessions) Get(id int) (admin *model.Admin) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return as.data[id]
}

func (as *AdminSessions) Del(id int) {
	as.mu.Lock()
	defer as.mu.Unlock()
	delete(as.data, id)
}

var AdminSessionContainer = &AdminSessions{
	data: make(map[int]*model.Admin, 10),
}

var adminIDContextKey = "loginAdminID"

func NewAdminAuthMiddleware(key []byte) *xx.Handler {
	type ps struct {
		Authorization string
	}
	doc := &xx.Doc{
		Title: "管理员登陆中间件",
		Params: xx.Params{
			{
				Type:   xx.Header,
				Schema: &ps{},
			},
		},
		Responses: xx.Responses{
			{
				Body: xx.MsgData(xx.MsgTypeWarning, "管理员未登录"),
			},
			{
				Body: xx.MsgData(xx.MsgTypeWarning, "管理员授权过期"),
			},
			{
				Body: xx.MsgData(xx.MsgTypeWarning, "管理员授权失败"),
			},
		},
	}
	return &xx.Handler{
		Doc: doc,
		HandleFunc: func(ctx *xx.Context) {
			p := &ps{}
			err := ctx.Unmarshal(p)
			if err != nil {
				ctx.Error(err)
			} else {
				if p.Authorization == "" {
					ctx.MsgWarning("管理员未登录")
				} else {
					auth := strings.TrimPrefix(p.Authorization, AuthScheme)
					token, err := jwt.ParseWithClaims(auth, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
						return key, nil
					})
					if err != nil {
						ctx.Error(err)
					} else {
						if !token.Valid {
							ctx.MsgWarning("管理员授权失败")
						} else {
							claims := token.Claims.(*jwt.StandardClaims)
							if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
								ctx.MsgWarning("管理员授权过期")
							} else {
								id, err := strconv.Atoi(claims.Id)
								if err != nil {
									ctx.Error(err)
								} else {
									ctx.Set(adminIDContextKey, id)
								}
							}
						}
					}
				}
			}
		},
	}
}

var adminContextKey = "adminContextKey"

// 从上下文中获取管理员信息，需要在中间件之后使用
func GetAdminFromContext(ctx *xx.Context) (*model.Admin, bool) {
	admin, ok := ctx.Get(adminContextKey).(*model.Admin)
	if ok {
		return admin, ok
	} else {
		id, ok := ctx.Get(adminIDContextKey).(int)
		if !ok {
			return nil, false
		}
		admin, _ = model.GetAdmin(id)
		if admin != nil {
			ctx.Set(adminContextKey, admin)
			return admin, ok
		}
	}
	return nil, false
}

// 从上下文中获取管理员ID，需要在中间件之后使用
func GetAdminIDFromContext(ctx *xx.Context) (int, bool) {
	id, ok := ctx.Get(adminIDContextKey).(int)
	return id, ok
}

func LoginHandler(ctx *xx.Context) {
	ff
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

func DecodeBasicAuthorization(auth string) (username, password string, expireTime int64, err error) {
	if !strings.HasPrefix(auth, prefix) {
		return "", "", 0, ErrBasicAuthorizationFormatIncorrect
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return "", "", 0, errors.Wrapf(err, "decode authorization token [%s] failed", auth)
	}
	auth, err = crypto.DecryptString(string(c))
	if err != nil {
		return "", "", 0, errors.Wrapf(err, "decode authorization token [%s] failed", auth)
	}
	cs := strings.Split(auth, ":")
	if len(cs) == 3 {
		expireTime, err := strconv.ParseInt(cs[2], 10, 64)
		if err != nil {
			return "", "", 0, errors.Wrap(err, "parse token expire time")
		} else {
			return cs[0], cs[1], expireTime, nil
		}
	} else {
		return "", "", 0, errors.Errorf("authorization token [%s] format incorrect", auth)
	}
}

func DelAdminFromSession(username string) {
	AdminSessionContainer.Del(username)
}
