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
	exp, _ := hex.DecodeString("4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45")
	checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := Keccak256Hash(in); return h[:] }, msg, exp)
}

func TestKeccak256Hasher(t *testing.T) {
	msg1 := []byte{79, 110, 82, 101, 119, 97, 114, 100, 40, 97, 100, 100, 114, 101, 115, 115, 44, 117, 105, 110, 116, 50, 53, 54, 41}
	fmt.Println("len", len(msg1))
	fmt.Println(Keccak256(msg1))

	//msg := []byte("abc")
	//exp, _ := hex.DecodeString("4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45")
	//checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := HashDataToBytes(in); return h[:] }, msg, exp)
}

func Elapsed(startTime time.Time) int64 {
	end := time.Now().UnixNano() / int64(time.Millisecond)
	start := startTime.UnixNano() / int64(time.Millisecond)
	diff := end - start
	return diff
}

func TestKeccak256HasherPerf(t *testing.T) {
	msg1 := []byte{79, 110, 73, 110, 105, 116, 105, 97, 116, 101, 87, 105, 116, 104, 100, 114, 97, 119, 97, 108, 40, 97, 100, 100, 114, 101, 115, 115, 41}
	fmt.Println("len", len(msg1))
	startTime := time.Now()
	var count int64
	count = 1000000
	var i int64
	for i = 0; i < count; i++ {
		Keccak256(msg1)
	}
	fmt.Println("elapsed", Elapsed(startTime))

	//msg := []byte("abc")
	//exp, _ := hex.DecodeString("4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45")
	//checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := HashDataToBytes(in); return [:] }, msg, exp)
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
