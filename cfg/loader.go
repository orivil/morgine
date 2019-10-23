// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type ReadWriter interface {
	Read(name string) (data []byte, err error)
	Write()
}

type Loader struct {
	rw io.ReadWriter
	mu sync.Mutex
}

const (
	jsonExt = ".json"
	ymlExt  = ".yml"
	yamlExt = ".yaml"
)

// configuration file directory, could change from outside
var Dir = "configs"

// Unmarshal parses the JSON/YAML data and stores the result
// in the value pointed to by schema. If v is nil or not a pointer,
// Unmarshal returns an error.
//
// the config value comes form 3 ways:
// 1. form config file if config file exist
// 2. form the content value if config file not exist
// 3. form os environment, and os environment value will cover the config value
//
// file is the file name under the Dir directory
// content is the default file value, if the file not exist, create one by the
// content value.
func Unmarshal(file, content string, schema interface{}) error {
	err := os.MkdirAll(Dir, os.ModePerm)
	if err != nil {
		return err
	}
	filename := filepath.Join(Dir, file)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			data = []byte(content)
			err = ioutil.WriteFile(filename, data, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "create config file failed")
			}
		} else {
			return errors.Wrap(err, "read config file failed")
		}
	}
	ext := filepath.Ext(file)

	// for collecting config keys
	cfgs := make(Configs, 5)
	switch ext {
	case jsonExt:
		err = json.Unmarshal(data, &cfgs)
		if err != nil {
			return errors.Wrap(err, "unmarshal config file failed")
		}
		// get config value form os environment
		err = cfgs.LoadOSEnv()
		if err != nil {
			return err
		}
		data, err = json.Marshal(cfgs)
		if err != nil {
			return err
		}

	case yamlExt, ymlExt:
		err = yaml.Unmarshal(data, &cfgs)
		if err != nil {
			return errors.Wrap(err, "unmarshal config file failed")
		}
		// get config value form os environment
		err = cfgs.LoadOSEnv()
		if err != nil {
			return err
		}
		data, err = yaml.Marshal(cfgs)
		if err != nil {
			return errors.WithStack(err)
		}
	default:
		return errors.Errorf("config file type [%s] not supported", ext)
	}

	switch ext {
	case jsonExt:
		err = json.Unmarshal(data, schema)
		if err != nil {
			return errors.Wrap(err, "unmarshal config file failed")
		}

	case yamlExt, ymlExt:
		err = yaml.Unmarshal(data, schema)
		if err != nil {
			return errors.Wrap(err, "unmarshal config file failed")
		}
	default:
		return errors.Errorf("config file type [%s] not supported", ext)
	}
	return nil
}

// UnmarshalMap is a helper function for loading the config data to map
func UnmarshalMap(file, content string) (schema Configs, err error) {
	schema = make(Configs, 5)
	err = Unmarshal(file, content, &schema)
	if err != nil {
		return nil, err
	} else {
		return schema, nil
	}
}
