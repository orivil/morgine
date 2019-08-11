// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package crypto

import (
	"github.com/orivil/morgine/cfg"
	"github.com/orivil/morgine/utils/crypto"
	"sync"
)

var aesCrypt = &crypt{}

type crypt struct {
	cry crypto.Interface
	mu  sync.RWMutex
}

func (c *crypt) Get() crypto.Interface {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cry
}

func (c *crypt) Set(key string) error {
	cry, err := crypto.NewCFBCrypto(key)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cry = cry
	return nil
}

func EncryptString(str string) (encrypted string, err error) {
	data, err := aesCrypt.Get().Encrypt([]byte(str))
	if err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

func DecryptString(str string) (decrypted string, err error) {
	data, err := aesCrypt.Get().Decrypt([]byte(str))
	if err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

var defaultFile = `# aes 加密密钥
aes_key: "change this pass"`

func init() {
	configs, err := cfg.UnmarshalMap("aes-crypt.yml", defaultFile)
	if err != nil {
		panic(err)
	}
	err = aesCrypt.Set(configs.GetStr("aes_key"))
	if err != nil {
		panic(err)
	}
}
