// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package service

import (
	"github.com/orivil/morgine/cfg"
	"github.com/orivil/morgine/service"
)

var Container = service.NewContainer(true)

var configsProvider = service.NewServiceProvider(func(c *service.Container) (value interface{}, err error) {
	return nil, nil
})

func SetConfigs(c *service.Container, configs cfg.Configs) {
	c.SetCache(configsProvider, configs)
}

func GetConfigs(c *service.Container) cfg.Configs {
	configs, _ := c.Get(configsProvider)
	return configs.(cfg.Configs)
}