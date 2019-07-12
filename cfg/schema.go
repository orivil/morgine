// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
)

type Configs map[string]interface{}

func (cs Configs) Unmarshal(schema interface{}) error {
	data, err := yaml.Marshal(cs)
	if err != nil {
		return err
	} else {
		return yaml.Unmarshal(data, schema)
	}
}

func (cs Configs) Copy() Configs {
	dst := make(Configs, len(cs))
	for key, value := range cs {
		dst[key] = value
	}
	return dst
}

func (cs Configs) Set(name string, value interface{}) {
	cs[name] = value
}

func (cs Configs) GetStr(name string) string {
	return cs[name].(string)
}

func (cs Configs) GetInt(name string) int {
	return cs[name].(int)
}

func (cs Configs) GetFloat(name string) float64 {
	return cs[name].(float64)
}

func (cs Configs) GetBool(name string) bool {
	return cs[name].(bool)
}

func (cs Configs) GetSliceStr(name string) []string {
	return cs[name].([]string)
}

func (cs Configs) GetSliceInt(name string) []int {
	return cs[name].([]int)
}

func (cs Configs) GetSliceFloat(name string) []float64 {
	return cs[name].([]float64)
}

func (cs Configs) GetSliceBool(name string) []bool {
	return cs[name].([]bool)
}

func (cs Configs) LoadOSEnv() (err error) {
	for key, value := range cs {
		osValue := os.Getenv(key)
		if osValue != "" {
			switch value.(type) {
			case string:
				cs[key] = osValue
			case int:
				cs[key], err = strconv.Atoi(osValue)
				if err != nil {
					return errors.Wrapf(err, "config '%s' need int value, got '%s'", key, osValue)
				}
			case bool:
				switch osValue {
				case "y", "Y", "yes", "YES", "Yes", "1", "t", "T", "true", "TRUE", "True":
					cs[key] = true
				case "n", "N", "no", "NO", "No", "0", "f", "F", "false", "FALSE", "False":
					cs[key] = false
				default:
					return errors.Wrapf(err, "config '%s' need bool value, got '%s'", key, osValue)
				}
			case float64:
				cs[key], err = strconv.ParseFloat(osValue, 64)
				if err != nil {
					return errors.Wrapf(err, "config '%s' need float64 value, got '%s'", key, osValue)
				}
			default:
				return errors.New("os config value only support 'string', 'int', 'float64' or 'bool'")
			}
		}
	}
	return nil
}
