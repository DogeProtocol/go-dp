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

//This file was added for go-dogep project (Doge Protocol Platform)

package cryptopq

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/cryptopq/oqs"
	"io"
	"io/ioutil"
	"math/big"
	"os"
)

const HASH_PUBKEY_BYTES_INDEX_START = 12

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

// ToOQS creates a private key with the given D value.
func ToOQS(d []byte) (*oqs.PrivateKey, error) {
	return toOQS(d, true)
}

// ToOQSUnsafe blindly converts a binary blob to a private key. It should almost
// never be used unless you are sure the input is valid and want to avoid hitting
// errors due to bad origin encoding (0 prefixes cut off).
func ToOQSUnsafe(d []byte) *oqs.PrivateKey {
	privy, _ := toOQS(d, false)
	return privy
}

// toOQS creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toOQS(d []byte, strict bool) (*oqs.PrivateKey, error) {
	privKey, err := oqs.ConvertBytesToPrivate(d)
	//Check valid private key (C library)
	//get publickey
	if err != nil {
		return nil, err
	}

	pub, err := oqs.RecoverPubkeyByPrivate(privKey.D.Bytes())
	if err != nil {
		return nil, err
	}

	pubkey, err := oqs.ConvertBytesToPublic(pub)
	if err != nil {
		return nil, err
	}

	privKey.PublicKey = *pubkey
	//get private key
	return privKey, err
}

// FromOQS exports a private key into a binary dump.
func FromOQS(priv *oqs.PrivateKey) ([]byte, error) {
	return oqs.ExportPrivateKey(priv)
}

// UnmarshalPubkey converts bytes to an oqs public key.
func UnmarshalPubkey(pub []byte) (*oqs.PublicKey, error) {
	pubKey, error := oqs.ConvertBytesToPublic(pub)
	return pubKey, error
}

func FromOQSPub(pub *oqs.PublicKey) ([]byte, error) {
	return oqs.ExportPublicKey(pub)
}

// HexToOQS parses an oqs private key.
func HexToOQS(hexkey string) (*oqs.PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, err
	}

	if byteErr, ok := err.(hex.InvalidByteError); ok {
		return nil, fmt.Errorf("invalid hex character %q in private key", byte(byteErr))
	} else if err != nil {
		return nil, errors.New("invalid hex data for private key")
	}
	return ToOQS(b)
}

// LoadOQS loads an oqs private key from the given file.
func LoadOQS(file string) (*oqs.PrivateKey, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	r := bufio.NewReader(fd)
	buf := make([]byte, oqs.PrivateKeyLen+oqs.PrivateKeyLen)
	n, err := readASCII(buf, r)
	if err != nil {
		return nil, err
	} else if n != len(buf) {
		return nil, fmt.Errorf("key file too short, want oqs hex character")
	}
	if err := checkKeyFileEnd(r); err != nil {
		return nil, err
	}
	return HexToOQS(string(buf))
}

// readASCII reads into 'buf', stopping when the buffer is full or
// when a non-printable control character is encountered.
func readASCII(buf []byte, r *bufio.Reader) (n int, err error) {
	for ; n < len(buf); n++ {
		buf[n], err = r.ReadByte()
		switch {
		case err == io.EOF || buf[n] < '!':
			return n, nil
		case err != nil:
			return n, err
		}
	}
	return n, nil
}

// checkKeyFileEnd skips over additional newlines at the end of a key file.
func checkKeyFileEnd(r *bufio.Reader) error {
	for i := 0; ; i++ {
		b, err := r.ReadByte()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case b != '\n' && b != '\r':
			return fmt.Errorf("invalid character %q at end of key file", b)
		case i >= 2:
			return errors.New("key file too long, want 64 hex characters")
		}
	}
}

// SaveOQS saves a private key to the given file with
// restrictive permissions. The key data is saved hex-encoded.
func SaveOQS(file string, key *oqs.PrivateKey) error {
	data, err := FromOQS(key)
	if err != nil {
		return err
	}
	k := hex.EncodeToString(data)
	return ioutil.WriteFile(file, []byte(k), 0600)
}

// GenerateKey generates a new private key.
func GenerateKey() (*oqs.PrivateKey, error) {
	return oqs.GenerateKey()
}

func PubkeyToAddress(p oqs.PublicKey) (common.Address, error) {
	pubBytes, err := FromOQSPub(&p)
	tempAddr := common.Address{}
	if err != nil {
		return tempAddr, err
	}
	return common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[HASH_PUBKEY_BYTES_INDEX_START:]), nil
}
