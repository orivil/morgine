// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package morgine_redis

import (
	"github.com/go-redis/redis"
	"github.com/orivil/morgine/cfg"
)

var defaultConfig = `# 数据库地址
redis_addr: "localhost:6379"

# 数据库密码
redis_password: ""
`

type Env struct {
	RedisAddr     string `yaml:"redis_addr"`
	RedisPassword string `yaml:"redis_password"`
}

func (e *Env) Connect(db int) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     e.RedisAddr,
		Password: e.RedisPassword,
		DB:       db,
	})
	err = client.Ping().Err()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func InitConfig(fileName string, db int, call func(client *redis.Client)) {
	env := &Env{}
	err := cfg.Unmarshal(fileName, defaultConfig, env)
	if err != nil {
		panic(err)
	}
	client, err := env.Connect(db)
	if err != nil {
		panic(err)
	} else {
		call(client)
	}
}
