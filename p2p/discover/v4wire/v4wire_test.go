// Copyright 2019 The go-ethereum Authors
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

package v4wire

import (
	"encoding/hex"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/rlp"
	"github.com/davecgh/go-spew/spew"
	"net"
	"reflect"
	"testing"
)

const (
	WirePublicKeyLen = 241
)

var (
	key1, _ = cryptobase.SigAlg.GenerateKey()
	key2, _ = cryptobase.SigAlg.GenerateKey()
	key3, _ = cryptobase.SigAlg.GenerateKey()
	key4, _ = cryptobase.SigAlg.GenerateKey()

	hexpubkey  = hex.EncodeToString(key1.PriData)
	hexpubkey1 = hex.EncodeToString(key1.PriData)
	hexpubkey2 = hex.EncodeToString(key2.PriData)
	hexpubkey3 = hex.EncodeToString(key3.PriData)
	hexpubkey4 = hex.EncodeToString(key4.PriData)

	hexPrivatekey, _ = cryptobase.SigAlg.PrivateKeyToHex(key1)
)

// EIP-8 test vectors.
var testPackets = []struct {
	input      string
	wantPacket Packet ///interface{}
}{
	{

		input: "",
		wantPacket: &Ping{
			Version:    4,
			From:       Endpoint{net.ParseIP("127.0.0.1").To4(), 3322, 5544},
			To:         Endpoint{net.ParseIP("::1"), 2222, 3333},
			Expiration: 1136239445,
		},
	},
	{

		input: "",
		wantPacket: &Ping{
			Version:    4,
			From:       Endpoint{net.ParseIP("127.0.0.1").To4(), 3322, 5544},
			To:         Endpoint{net.ParseIP("::1"), 2222, 3333},
			Expiration: 1136239445,
			ENRSeq:     1,
			Rest:       []rlp.RawValue{{0x02}},
		},
	},
	{

		input: "",
		wantPacket: &Findnode{
			Target:     hexPubkey(hexpubkey),
			Expiration: 1136239445,
			Rest:       []rlp.RawValue{{0x82, 0x99, 0x99}, {0x83, 0x99, 0x99, 0x99}},
		},
	},
	{

		input: "",
		wantPacket: &Neighbors{
			Nodes: []Node{
				{
					ID:  hexPubkey(hexpubkey1),
					IP:  net.ParseIP("99.33.22.55").To4(),
					UDP: 4444,
					TCP: 4445,
				},
				{
					ID:  hexPubkey(hexpubkey2),
					IP:  net.ParseIP("1.2.3.4").To4(),
					UDP: 1,
					TCP: 1,
				},
				{
					ID:  hexPubkey(hexpubkey3),
					IP:  net.ParseIP("2001:db8:3c4d:15::abcd:ef12"),
					UDP: 3333,
					TCP: 3333,
				},
				{
					ID:  hexPubkey(hexpubkey4),
					IP:  net.ParseIP("2001:db8:85a3:8d3:1319:8a2e:370:7348"),
					UDP: 999,
					TCP: 1000,
				},
			},
			Expiration: 1136239445,
			Rest:       []rlp.RawValue{{0x01}, {0x02}, {0x03}},
		},
	},
}

// This test checks that the decoder accepts packets according to EIP-8.
func TestForwardCompatibility(t *testing.T) {
	testkey, _ := cryptobase.SigAlg.HexToPrivateKey(hexPrivatekey)
	wantNodeKey := EncodePubkey(&testkey.PublicKey)

	for _, test := range testPackets {

		req := test.wantPacket
		packet1, _, _ := Encode(testkey, req)
		test.input = hex.EncodeToString(packet1)

		input, err := hex.DecodeString(test.input)
		if err != nil {
			t.Fatalf("invalid hex: %s", test.input)
		}
		packet, nodekey, _, err := Decode(input)
		if err != nil {
			t.Errorf("did not accept packet %s\n%v", test.input, err)
			continue
		}
		if !reflect.DeepEqual(packet, test.wantPacket) {
			t.Errorf("got %s\nwant %s", spew.Sdump(packet), spew.Sdump(test.wantPacket))
		}
		if WirePubKeyEquals(nodekey, wantNodeKey) == false {
			t.Errorf("got id %v\nwant id %v", nodekey, wantNodeKey)
		}
	}

}

func hexPubkey(h string) (ret Pubkey) {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	return CreateWirePubKey(b)
}
