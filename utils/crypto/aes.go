// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var ErrIncorrectCipherText = errors.New("incorrect cipher text")

type CFBCrypto struct {
	block cipher.Block
}

// 16, 24, 32 位 key 分别对应：AES-128, AES-192, or AES-256 加解密方式
func NewCFBCrypto(key string) (ac *CFBCrypto, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	} else {
		return &CFBCrypto{block: block}, nil
	}
}

func (ac *CFBCrypto) Encrypt(text []byte) ([]byte, error) {
	cipherText := make([]byte, aes.BlockSize+len(text))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(ac.block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], []byte(text))
	return []byte(base64.URLEncoding.EncodeToString(cipherText)), nil
}

func (ac *CFBCrypto) Decrypt(text []byte) ([]byte, error) {
	data, err := base64.URLEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, ErrIncorrectCipherText
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(ac.block, iv)
	cfb.XORKeyStream(data, data)
	return data, nil
}
