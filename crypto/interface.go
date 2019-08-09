// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package crypto

type Crypto interface {
	Encrypt(text []byte) ([]byte, error)
	Decrypt(text []byte) ([]byte, error)
}
