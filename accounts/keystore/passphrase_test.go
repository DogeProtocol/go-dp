// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package keystore

import (
	"github.com/DogeProtocol/dp/accounts"
	"os"
	"strconv"
	"testing"
)

const (
	veryLightScryptN = 2
	veryLightScryptP = 1
)

// Tests that a json key file can be decrypted and encrypted in multiple rounds.
func TestKeyEncryptDecrypt(t *testing.T) {

	// Do a few rounds of decryption and encryption
	for i := 0; i < 3; i++ {
		dir, ks := tmpKeyStore(t, true)
		defer os.RemoveAll(dir)

		pass := strconv.Itoa(i)
		a1, err := ks.NewAccount(pass)
		if err != nil {
			t.Fatal(err)
		}
		if err := ks.Unlock(a1, pass); err != nil {
			t.Fatal(err)
		}
		if _, err := ks.SignHash(accounts.Account{Address: a1.Address}, testSigData); err != nil {
			t.Fatal(err)
		}
	}

}
