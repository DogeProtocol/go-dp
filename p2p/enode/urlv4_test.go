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
	"encoding/base64"
	"errors"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/hashingalgorithm"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/DogeProtocol/dp/p2p/enr"
)

var (
	h                 = hashingalgorithm.NewHashState()
	keyTest1, _       = cryptobase.SigAlg.GenerateKey()
	hexprvkeytest1, _ = cryptobase.SigAlg.PrivateKeyToHex(keyTest1)
	hexpubkeytest1, _ = cryptobase.SigAlg.PublicKeyToHex(&keyTest1.PublicKey)
	signTest1, _      = cryptobase.SigAlg.Sign(h.Sum(nil), keyTest1)
	hexsigntest1      = base64.RawURLEncoding.EncodeToString(signTest1)
)

func init() {
	lookupIPFunc = func(name string) ([]net.IP, error) {
		if name == "node.example.org" {
			return []net.IP{{33, 44, 55, 66}}, nil
		}
		return nil, errors.New("no such host")
	}
}

var parseNodeTests = []struct {
	input      string
	wantError  string
	wantResult *Node
}{
	// Records
	{

		input: "dynamic1",
		wantResult: func() *Node {
			testKey, _ := cryptobase.SigAlg.HexToPrivateKey(hexprvkeytest1)
			var r enr.Record
			r.Set(enr.IP{127, 0, 0, 1})
			r.SetSeq(99)
			SignV4(&r, testKey)
			n, _ := New(ValidSchemes, &r)
			return n
		}(),
	},
	// Invalid Records
	{
		input:     "enr:",
		wantError: "EOF", // could be nicer
	},
	{
		input:     "enr:x",
		wantError: "illegal base64 data at input byte 0",
	},
	{
		input: "dynamic2",

		wantError: enr.ErrInvalidSig.Error(),
	},
	// Complete node URLs with IP address and ports
	{
		input:     "enode://" + hexpubkeytest1 + "@invalid.:3",
		wantError: `no such host`,
	},
	{
		input:     "enode://" + hexpubkeytest1 + "@127.0.0.1:foo",
		wantError: `invalid port`,
	},
	{
		input:     "enode://" + hexpubkeytest1 + "@127.0.0.1:3?discport=foo",
		wantError: `invalid discport in query`,
	},
	{
		input: "enode://" + hexpubkeytest1 + "@127.0.0.1:52150",
		wantResult: NewV4(
			hexPubkey(hexpubkeytest1),
			net.IP{127, 0, 0, 1},
			52150,
		),
	},
	{
		input: "enode://" + hexpubkeytest1 + "@[::]:52150",
		wantResult: NewV4(
			hexPubkey(hexpubkeytest1),
			net.ParseIP("::"),
			52150,
		),
	},
	{
		input: "enode://" + hexpubkeytest1 + "@[2001:db8:3c4d:15::abcd:ef12]:52150",
		wantResult: NewV4(
			hexPubkey(hexpubkeytest1),
			net.ParseIP("2001:db8:3c4d:15::abcd:ef12"),
			52150,
		),
	},
	{
		input: "enode://" + hexpubkeytest1 + "@127.0.0.1:52150?discport=22334",
		wantResult: NewV4(
			hexPubkey(hexpubkeytest1),
			net.IP{0x7f, 0x0, 0x0, 0x1},
			52150,
		),
	},
	// Incomplete node URLs with no address
	{
		input: "enode://" + hexpubkeytest1,
		wantResult: NewV4(
			hexPubkey(hexpubkeytest1),
			nil, 0,
		),
	},
	// Invalid URLs
	{
		input:     "",
		wantError: errMissingPrefix.Error(),
	},
	{
		input:     hexpubkeytest1,
		wantError: errMissingPrefix.Error(),
	},
	{
		input:     "01010101",
		wantError: errMissingPrefix.Error(),
	},
	{
		input:     "enode://01010101@123.124.125.126:3",
		wantError: `invalid public key (invalid public key length)`,
	},
	{
		input:     "enode://01010101",
		wantError: `invalid public key (invalid public key length)`,
	},
	{
		input:     "http://foobar",
		wantError: errMissingPrefix.Error(),
	},
	{
		input:     "://foo",
		wantError: errMissingPrefix.Error(),
	},
}

func hexPubkey(h string) *signaturealgorithm.PublicKey {
	k, err := parsePubkey(h)
	if err != nil {
		panic(err)
	}
	return k
}

func TestParseNode(t *testing.T) {
	testKey, _ := cryptobase.SigAlg.HexToPrivateKey(hexprvkeytest1)

	var r enr.Record
	r.Set(enr.IP{127, 0, 0, 1})
	r.SetSeq(99)
	SignV4(&r, testKey)
	result1, _ := New(ValidSchemes, &r)

	var r1 enr.Record
	r1.SetSeq(99)
	SignV4(&r1, testKey)
	result2, _ := New(ValidSchemes, &r1)

	for _, test := range parseNodeTests {
		if test.input == "dynamic1" {
			test.input = result1.String()
		}
		if test.input == "dynamic2" {
			test.input = result2.String()
		}
		_, err := Parse(ValidSchemes, test.input)

		if test.wantError != "" {
			if err == nil {
				t.Errorf("test %q:\n  got nil error, expected %#q", test.input, test.wantError)
				continue
			} else if !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("test %q:\n  got error %#q, expected %#q", test.input, err.Error(), test.wantError)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("test %q:\n  unexpected error: %v", test.input, err)
				continue
			}

		}
	}
}

func TestNodeString(t *testing.T) {
	for i, test := range parseNodeTests {
		if test.wantError == "" && strings.HasPrefix(test.input, "enode://") {
			n, _ := Parse(ValidSchemes, test.input)
			if !reflect.DeepEqual(n, test.wantResult) {
				t.Errorf("test %d:\n  result mismatch:\ngot:  %#v\nwant: %#v", i, n, test.wantResult)
			}
		}
	}
}
