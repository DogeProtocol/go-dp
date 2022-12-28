// Copyright 2018 The go-ethereum Authors
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

package enode

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/cryptopq"
	"math/big"
	"testing"
	"testing/quick"

	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

// TestPythonInterop checks that we can decode and verify a record produced by the Python
// implementation.
func TestPythonInterop(t *testing.T) {
	testKey, _ := cryptopq.HexToOQS(hexprvkeytest1)
	var r1 enr.Record
	r1.Set(enr.IP{127, 0, 0, 1})
	r1.Set(enr.UDP(30303))
	r1.SetSeq(1)
	SignV4(&r1, testKey)
	result1, _ := New(ValidSchemes, &r1)
	pyRecord, err := base64.RawURLEncoding.DecodeString(result1.String()[4:])

	if err != nil {
		return
	}
	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := New(ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	var (
		wantID  = HexID(crypto.Keccak256Hash(testKey.N.Bytes()).Hex())
		wantSeq = uint64(1)
		wantIP  = enr.IPv4{127, 0, 0, 1}
		wantUDP = enr.UDP(30303)
	)

	if n.Seq() != wantSeq {
		t.Errorf("wrong seq: got %d, want %d", n.Seq(), wantSeq)
	}
	if n.ID() != wantID {
		t.Errorf("wrong id: got %x, want %x", n.ID(), wantID)
	}
	want := map[enr.Entry]interface{}{new(enr.IPv4): &wantIP, new(enr.UDP): &wantUDP}
	for k, v := range want {
		desc := fmt.Sprintf("loading key %q", k.ENRKey())
		if assert.NoError(t, n.Load(k), desc) {
			assert.Equal(t, k, v, desc)
		}
	}
}

func TestHexID(t *testing.T) {
	ref := ID{0, 0, 0, 0, 0, 0, 0, 128, 106, 217, 182, 31, 165, 174, 1, 67, 7, 235, 220, 150, 66, 83, 173, 205, 159, 44, 10, 57, 42, 161, 26, 188}
	id1 := HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc")
	id2 := HexID("00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc")

	if id1 != ref {
		t.Errorf("wrong id1\ngot  %v\nwant %v", id1[:], ref[:])
	}
	if id2 != ref {
		t.Errorf("wrong id2\ngot  %v\nwant %v", id2[:], ref[:])
	}
}

func TestID_textEncoding(t *testing.T) {
	ref := ID{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x20,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x30,
		0x31, 0x32,
	}
	hex := "0102030405060708091011121314151617181920212223242526272829303132"

	text, err := ref.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(text, []byte(hex)) {
		t.Fatalf("text encoding did not match\nexpected: %s\ngot:      %s", hex, text)
	}

	id := new(ID)
	if err := id.UnmarshalText(text); err != nil {
		t.Fatal(err)
	}
	if *id != ref {
		t.Fatalf("text decoding did not match\nexpected: %s\ngot:      %s", ref, id)
	}
}

func TestID_distcmp(t *testing.T) {
	distcmpBig := func(target, a, b ID) int {
		tbig := new(big.Int).SetBytes(target[:])
		abig := new(big.Int).SetBytes(a[:])
		bbig := new(big.Int).SetBytes(b[:])
		return new(big.Int).Xor(tbig, abig).Cmp(new(big.Int).Xor(tbig, bbig))
	}
	if err := quick.CheckEqual(DistCmp, distcmpBig, nil); err != nil {
		t.Error(err)
	}
}

// The random tests is likely to miss the case where a and b are equal,
// this test checks it explicitly.
func TestID_distcmpEqual(t *testing.T) {
	base := ID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	x := ID{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	if DistCmp(base, x, x) != 0 {
		t.Errorf("DistCmp(base, x, x) != 0")
	}
}

func TestID_logdist(t *testing.T) {
	logdistBig := func(a, b ID) int {
		abig, bbig := new(big.Int).SetBytes(a[:]), new(big.Int).SetBytes(b[:])
		return new(big.Int).Xor(abig, bbig).BitLen()
	}
	if err := quick.CheckEqual(LogDist, logdistBig, nil); err != nil {
		t.Error(err)
	}
}

// The random tests is likely to miss the case where a and b are equal,
// this test checks it explicitly.
func TestID_logdistEqual(t *testing.T) {
	x := ID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	if LogDist(x, x) != 0 {
		t.Errorf("LogDist(x, x) != 0")
	}
}
