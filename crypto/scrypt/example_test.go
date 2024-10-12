// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scrypt_test

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/ethclient/crypto/scrypt"
)

func Example() {
	// DO NOT use this salt value; generate your own random salt. 8 bytes is
	// a good length.
	salt := []byte{0xc8, 0x28, 0xf2, 0x58, 0xa7, 0x6a, 0xad, 0x7b}

	dk, err := scrypt.Key([]byte("some password"), salt, 1<<15, 8, 1, 32, uint(0x30))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(dk))
	// Output: f/vqXLfTHi5L9w8Jy09AzozKOZoaHimO6gMHjDqXLm4=
}
