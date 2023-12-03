// Copyright 2014 The go-ethereum Authors
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

package crypto

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"testing"
	"time"
)

var testAddrHex = "970e8128ab834e8eac17ab8e3812f010678cf791"
var testPrivHex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"

// These tests are sanity checks.
// They should ensure that we don't e.g. use Sha3-224 instead of Sha3-256
// and that the sha3 library uses keccak-f permutation.
func TestKeccak256Hash(t *testing.T) {
	msg := []byte("abc")
	exp, _ := hex.DecodeString("3a985da74fe225b2045c172d6bd390bd855f086e3e9d525b46bfe24511431532")
	checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := Keccak256Hash(in); return h[:] }, msg, exp)
}

func TestKeccak256Hasher(t *testing.T) {
	msg := []byte{79, 110, 67, 104, 97, 110, 103, 101, 86, 97, 108, 105, 100, 97, 116, 111, 114, 40, 97, 100, 100, 114, 101, 115, 115, 44, 97, 100, 100, 114, 101, 115, 115, 44, 97, 100, 100, 114, 101, 115, 115, 41}
	fmt.Println("len", len(msg))
	hash := Keccak256(msg)
	expectedHash := []byte{102, 215, 164, 222, 167, 72, 81, 162, 220, 192, 57, 244, 193, 125, 210, 134, 32, 129, 8, 60, 41, 218, 206, 177, 169, 52, 103, 131, 222, 145, 133, 206}
	fmt.Println(Keccak256(msg))
	if bytes.Compare(hash, expectedHash) != 0 {
		t.Fatalf("failed")
	}

}

func Elapsed(startTime time.Time) int64 {
	end := time.Now().UnixNano() / int64(time.Millisecond)
	start := startTime.UnixNano() / int64(time.Millisecond)
	diff := end - start
	return diff
}

func TestKeccak256HasherPerf(t *testing.T) {
	msg := []byte{79, 110, 73, 110, 105, 116, 105, 97, 116, 101, 87, 105, 116, 104, 100, 114, 97, 119, 97, 108, 40, 97, 100, 100, 114, 101, 115, 115, 41}
	fmt.Println("len", len(msg))
	startTime := time.Now()
	var count int64
	count = 1000000
	var i int64
	for i = 0; i < count; i++ {
		Keccak256(msg)
	}
	fmt.Println("elapsed", Elapsed(startTime))
}

func checkhash(t *testing.T, name string, f func([]byte) []byte, msg, exp []byte) {
	sum := f(msg)
	if !bytes.Equal(exp, sum) {
		t.Fatalf("hash %s mismatch: want: %x have: %x", name, exp, sum)
	}
}

func checkAddr(t *testing.T, addr0, addr1 common.Address) {
	if addr0 != addr1 {
		t.Fatalf("address mismatch: want: %x have: %x", addr0, addr1)
	}
}
