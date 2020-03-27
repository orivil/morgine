// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package limiter

import (
	"sync"
	"time"
)

type user struct {
	session string
	failed int
	waitAt *time.Time
}

type TimeProvider interface {
	// 获得失败 failed 次数之后的等待时间
	GetWaitTime(failed int) time.Duration

	// 获得当前时间
	GetNowTime() time.Time
}

// OperationContainer 用于操作失败检测，例如用户登录失败记录，验证吗检测失败记录等
type OperationContainer struct {
	users map[string]*user
	timeProvider TimeProvider
	mu sync.Mutex
}

func NewOperationContainer(timeProvider TimeProvider) *OperationContainer {
	return &OperationContainer {
		users:    make(map[string]*user, 50),
		timeProvider: timeProvider,
		mu:       sync.Mutex{},
	}
}

// 操作成功，删除失败记录
func (uc *OperationContainer) Success(session string) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	delete(uc.users, session)
}

// 添加失败记录，获得等待时间
func (uc *OperationContainer) Failed(session string) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	usr := uc.users[session]
	if usr == nil {
		usr = &user {
			session:session,
			failed: 1,
		}
		uc.users[session] = usr
	} else {
		usr.failed++
	}
	waitDuration := uc.timeProvider.GetWaitTime(usr.failed)
	if waitDuration > 0 {
		wa := uc.timeProvider.GetNowTime().Add(waitDuration)
		usr.waitAt = &wa
	} else {
		usr.waitAt = nil
	}
}

// 检测当前时间是否允许操作
func (uc *OperationContainer) Allow(session string) bool {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	usr := uc.users[session]
	if usr != nil && usr.waitAt != nil {
		return usr.waitAt.Before(uc.timeProvider.GetNowTime())
	} else {
		return true
	}
}

// 获得用户等待时间
func (uc *OperationContainer) GetWaitTime(session string) (waitAt *time.Time) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	usr := uc.users[session]
	if usr == nil {
		return nil
	} else {
		return usr.waitAt
	}
}
