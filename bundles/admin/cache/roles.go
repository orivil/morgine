// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	admin_model "github.com/orivil/morgine/bundles/admin/model"
	"strconv"
	"time"
)

var Roles *roleCache

func InitRoleCache(keyPrefix string, client *redis.Client) {
	Roles = &roleCache{client:client, key:keyPrefix}
}

type roleCache struct {
	client *redis.Client
	key string
}

func (rc *roleCache) getKey(userID int) string {
	return rc.key + strconv.Itoa(userID)
}

func (rc *roleCache) Get(userID int) ([]*admin_model.Role, error) {
	v, err := rc.client.Get(rc.getKey(userID)).Result()
	if err != nil {
		return nil, err
	} else {
		var rs []*admin_model.Role
		err = json.Unmarshal([]byte(v), &rs)
		if err != nil {
			return nil, err
		} else {
			return rs, nil
		}
	}
}

func (rc *roleCache) Set(userID int, rs []*admin_model.Role) error {
	data, err := json.Marshal(rs)
	if err != nil {
		return err
	}
	return rc.client.Set(rc.getKey(userID), string(data), 24 * time.Hour).Err()
}