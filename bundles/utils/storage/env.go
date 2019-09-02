// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package storage

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

var env =
`# yes 使用 OSS 存储, no 则使用本地硬盘存储
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

type Env struct {
	UseOSS             string `yaml:"use_oss"`
	LocalServeHost     string `yaml:"local_serve_host"`
	OSSEndpoint        string `yaml:"oss_endpoint"`
	OSSAccessKeyID     string `yaml:"oss_access_key_id"`
	OSSAccessKeySecret string `yaml:"oss_access_key_secret"`
	client *oss.Client
	Storage
}

func (e *Env) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(e)
	if err != nil {
		return err
	} else {
		if e.UseALiYunOSS() {
			e.client, err = oss.New(e.OSSEndpoint, e.OSSAccessKeyID, e.OSSAccessKeySecret)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Env) InitLocalStorage(dir string, corsHandler func(header http.Header)) (*LocalStorage, error) {
	return NewLocalStorage(dir, e.LocalServeHost, corsHandler)
}

func (e *Env) InitOssStorage(bucketName, cdnHost string, urlMaxAge int64, corsRules []oss.CORSRule) (*OssStorage, error) {
	return NewOssStorage(bucketName, corsRules, urlMaxAge, cdnHost, e.OssServeHost(bucketName))
}

func (e *Env) UseALiYunOSS() bool {
	return e.UseOSS == "yes" || e.UseOSS == "true"
}

func (e *Env) OssServeHost(bucketName string) string {
	return bucketName + "." + strings.TrimPrefix(e.OSSEndpoint, "http://")
}

// get bucket or create a new one with cross site rules
func (e *Env) Bucket(name string, corsRules []oss.CORSRule) (b *oss.Bucket, err error) {
	ok, err := e.client.IsBucketExist(name)
	if err != nil {
		return nil, errors.Wrapf(err, "检查 bucket [%s] 是否存在时出错", name)
	}
	if !ok {
		err = e.client.CreateBucket(name)
		if err != nil {
			return nil, errors.Wrapf(err, "创建 bucket [%s] 时出错", name)
		}
		err = e.client.SetBucketCORS(name, corsRules)
		if err != nil {
			return nil, err
		}
	}
	return e.client.Bucket(name)
}
