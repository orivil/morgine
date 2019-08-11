// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package morgine_redis

import (
	"github.com/go-redis/redis"
	"time"
)

// 集合存储器
type SetStorage struct {
	client *redis.Client
}

func NewSetStorage(client *redis.Client) *SetStorage {
	return &SetStorage{client: client}
}

// 设置元素到集合中
func (ls *SetStorage) Set(key string, value string, maxAge time.Duration) error {
	err := ls.client.SAdd(key, value).Err()
	if err != nil {
		return err
	}
	if maxAge > 0 {
		return ls.client.Expire(key, maxAge).Err()
	}
	return nil
}

// 删除集合中的元素
func (ls *SetStorage) Del(key string, value ...string) error {
	var vs = make([]interface{}, len(value))
	for key, v := range value {
		vs[key] = v
	}
	return ls.client.SRem(key, vs...).Err()
}

// 检查元素在集合中是否存在
func (ls *SetStorage) IsExist(key string, value string) (bool, error) {
	return ls.client.SIsMember(key, value).Result()
}

// 统计集合元素数量
func (ls *SetStorage) Count(key string) (total int64, err error) {
	return ls.client.SCard(key).Result()
}

// 获得集合元素成员
func (ls *SetStorage) Get(key string, limit int64, cursor uint64) (values []string, nextCursor uint64, err error) {
	return ls.client.SScan(key, uint64(cursor), "", limit).Result()
}
