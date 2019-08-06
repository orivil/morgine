// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"github.com/orivil/morgine/cfg"
)

var serverConfig = `# 是否使用 https 加密协议
use_ssl: false

# 监听 http 端口, 如果开启 SSL, 则当前端口将作为跳转端口, 将所有请求都跳转到 https_port
http_port: ":9090"

# 监听 https 端口
https_port: ":443"

ssl_certificate: ""

ssl_certificate_key: ""

# 是否开启 debug
debug: true

# 开启请求日志, 支持热加载, 只对设置过 middles.Logger 中间件的请求有效
open_log: true`

var Env = &env{
	Debug:     true,
	OpenLog:   true,
	HttpPort:  ":9090",
	HttpsPort: ":443",
}

func init() {
	err := cfg.Unmarshal("server.yml", serverConfig, Env)
	if err != nil {
		panic(err)
	}
}

type env struct {
	UseSSL            bool   `yaml:"use_ssl" json:"use_ssl"`
	SSLCertificate    string `yaml:"ssl_certificate" json:"ssl_certificate"`
	SSLCertificateKey string `yaml:"ssl_certificate_key" json:"ssl_certificate_key"`
	HttpPort          string `yaml:"http_port" json:"http_port"`
	HttpsPort         string `yaml:"https_port" json:"https_port"`
	Debug             bool   `yaml:"debug" json:"debug"`
	OpenLog           bool   `yaml:"open_log" json:"open_log"`
}
