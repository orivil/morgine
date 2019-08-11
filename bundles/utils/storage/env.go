// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package storage

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/orivil/morgine/cfg"
	"github.com/pkg/errors"
	"strings"
)

var Env = &environment{LocalServeHost: "localhost"}

type environment struct {
	UseOSS             string `yaml:"use_oss"`
	LocalServeHost     string `yaml:"local_serve_host"`
	OSSEndpoint        string `yaml:"oss_endpoint"`
	OSSAccessKeyID     string `yaml:"oss_access_key_id"`
	OSSAccessKeySecret string `yaml:"oss_access_key_secret"`
}

func (e *environment) UseALiYunOSS() bool {
	return e.UseOSS == "yes" || e.UseOSS == "true"
}

func (e *environment) ServeHost(bucketName string) string {
	return bucketName + "." + strings.TrimPrefix(e.OSSEndpoint, "http://")
}

var defaultConfig = `# yes 使用 OSS 存储, no 则使用本地硬盘存储
use_oss: "no"

# 本地服务域名
local_serve_host: ""

# Endpoint以杭州为例，其它Region请按实际情况填写
oss_endpoint: "http://oss-cn-hangzhou.aliyuncs.com"

# 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM
# 账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
oss_access_key_id: "<yourAccessKeyId>"

oss_access_key_secret: "<yourAccessKeySecret>"
`

var client *oss.Client

func init() {
	err := cfg.Unmarshal("storage.yml", defaultConfig, Env)
	if err != nil {
		panic(err)
	}
	if Env.UseALiYunOSS() {
		client, err = oss.New(Env.OSSEndpoint, Env.OSSAccessKeyID, Env.OSSAccessKeySecret)
		if err != nil {
			panic(err)
		}
	}
}

// get bucket or create a new one with cross site rules
func Bucket(name string, corsRules []oss.CORSRule) (b *oss.Bucket, err error) {
	if client == nil {
		return nil, errors.Errorf("未开启 OSS 服务, 需要在配置文件中或环境变量中设置开启")
	}
	ok, err := client.IsBucketExist(name)
	if err != nil {
		return nil, errors.Wrapf(err, "检查 bucket [%s] 是否存在时出错", name)
	}
	if !ok {
		err = client.CreateBucket(name)
		if err != nil {
			return nil, errors.Wrapf(err, "创建 bucket [%s] 时出错", name)
		}
		err = client.SetBucketCORS(name, corsRules)
		if err != nil {
			return nil, err
		}
	}
	return client.Bucket(name)
}
