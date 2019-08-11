// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package morgine_redis

import (
	"github.com/go-redis/redis"
)

type ListWalker func(key string, members []string) (isContinue bool)

type ListStorage struct {
	client *redis.Client
}

func NewListStorage(client *redis.Client) *ListStorage {
	return &ListStorage{client: client}
}

func (ls *ListStorage) Push(key string, members ...interface{}) error {
	return ls.client.LPush(key, members...).Err()
}

func (ls *ListStorage) Range(limit int64, walk ListWalker) (err error) {
	var cursor uint64
	var keys []string
	for {
		keys, cursor, err = ls.client.Scan(cursor, "", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			for _, key := range keys {
				for {
					var start int64 = 0
					stop := limit - 1
					members, err := ls.client.LRange(key, start, stop).Result()
					if err != nil {
						if err == redis.Nil {
							break
						} else {
							return err
						}
					} else {
						if ln := len(members); ln > 0 {
							if !walk(key, members) {
								break
							} else {
								// 直接删除列表, 并不能保证列表已被正常处理
								err = ls.client.LTrim(key, limit, -1).Err()
								if err != nil {
									return err
								}
							}
						} else {
							break
						}
					}
				}
			}
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}
