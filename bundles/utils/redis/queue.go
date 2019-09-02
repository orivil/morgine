// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package morgine_redis

import (
	"github.com/go-redis/redis"
	"strconv"
)

// 队列存储器
type QueueStorage struct {
	client *redis.Client
}

func NewQueueStorage(client *redis.Client) *QueueStorage {
	return &QueueStorage{
		client: client,
	}
}

func (rs *QueueStorage) Del(key string, members ...string) error {
	args := make([]interface{}, len(members))
	for key, member := range members {
		args[key] = member
	}
	return rs.client.ZRem(key, args...).Err()
}

// 添加或修改队列成员
func (rs *QueueStorage) Set(key, member string, activeAt int64) error {
	return rs.client.ZAdd(key, &redis.Z{
		Score:  float64(activeAt),
		Member: member,
	}).Err()
}

type QueueDeleter func(key string, members ...string) error

type QueueWalker func(key string, members []string, deleter QueueDeleter) (isContinue bool)

// 获取队列过期成员并删除过期成员
func (rs *QueueStorage) Range(expireAt, limit int64, walk QueueWalker) (err error) {
	var cursor uint64
	var keys []string
	for {
		keys, cursor, err = rs.client.Scan(cursor, "", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			for _, key := range keys {
				var rangeBy = &redis.ZRangeBy{
					Min:    "0",
					Max:    strconv.FormatInt(expireAt, 10),
					Offset: 0,
					Count:  limit,
				}
				for {
					members, err := rs.client.ZRangeByScore(key, rangeBy).Result()
					if err != nil {
						if err == redis.Nil {
							break
						} else {
							return err
						}
					} else {
						if len(members) > 0 {
							if !walk(key, members, rs.Del) {
								break
							} else {
								rangeBy.Offset += int64(len(members))
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
