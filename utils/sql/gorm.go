// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package sql

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"sync"
)

var defaultConfig = `# 开启日志
db_log: true

# mysql postgres
db_dialect: "postgres"

# 数据库地址, 线上项目应该从OS环境变量中获取
db_host: "localhost"

# 数据库监听端口, 线上项目应该从OS环境变量中获取
db_port: ""

# 用户名, 线上项目应该从OS环境变量中获取
db_user: ""

# 密码, 线上项目应该从OS环境变量中获取1
db_password: ""

# 数据库名
db_name: ""

# 表前缀
db_sql_table_prefix: ""

# 最大空闲连接, 支持热重载
db_max_idle_connects: 5

# 最大活动连接, 支持热重载
db_max_opened_connects: 10

`

type Env struct {
	DBLog               bool   `yaml:"db_log"`
	DBDialect           string `yaml:"db_dialect"`
	DBHost              string `yaml:"db_host"`
	DBPort              string `yaml:"db_port"`
	DBUser              string `yaml:"db_user"`
	DBPassword          string `yaml:"db_password"`
	DBName              string `yaml:"db_name"`
	DBSqlTablePrefix    string `yaml:"db_sql_table_prefix"`
	DBMaxIdleConnects   int    `yaml:"db_max_idle_connects"`
	DBMaxOpenedConnects int    `yaml:"db_max_opened_connects"`
}

var once = sync.Once{}

func (e *Env) Connect() (*gorm.DB, error) {
	var arg string
	switch e.DBDialect {
	case "", "mysql":
		arg = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", e.DBUser, e.DBPassword, e.DBHost, e.DBPort, e.DBName)
	case "postgres":
		arg = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", e.DBUser, e.DBPassword, e.DBHost, e.DBPort, e.DBName)
	default:
		return nil, errors.New("目前只支持 postgres 或 mysql 数据库")
	}

	db, err := gorm.Open(e.DBDialect, arg)
	if err != nil {
		return nil, err
	}
	db.LogMode(e.DBLog)
	db.DB().SetMaxIdleConns(e.DBMaxIdleConnects)
	db.DB().SetMaxOpenConns(e.DBMaxOpenedConnects)
	if err = db.DB().Ping(); err != nil {
		return nil, err
	}
	if e.DBSqlTablePrefix != "" {
		db.Set("table_prefix", e.DBSqlTablePrefix)
		once.Do(func() {
			gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
				if tf, ok := db.Get("table_prefix"); ok {
					return tf.(string) + defaultTableName
				}
				return defaultTableName
			}
		})
	}
	return db, nil
}