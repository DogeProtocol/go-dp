// Copyright 2017 The go-ethereum Authors
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

//go:build !nacl && !js && cgo && !gofuzz
// +build !nacl,!js,cgo,!gofuzz

package cryptopq

import (
	"fmt"
	"github.com/ethereum/go-ethereum/cryptopq/oqs"
)

// RecoverPublicKey returns the uncompressed public key that created the given signature.
func RecoverPublicKey(hash, sig []byte) ([]byte, error) {
	return oqs.RecoverPubkey(hash, sig)
}

// SigToPub returns the public key that created the given signature.
func SigToPub(hash, sig []byte) (*oqs.PublicKey, error) {
	s, err := RecoverPublicKey(hash, sig)
	if err != nil {
		return nil, err
	}
	pub, error := oqs.ConvertBytesToPublic(s)
	if error != nil {
		return nil, error
	}
	return pub, nil
}

// Sign calculates an OQS signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given digest cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func Sign(digestHash []byte, prv *oqs.PrivateKey) (sig []byte, err error) {
	seckey, err := oqs.ExportPrivateKey(prv)
	if err != nil {
		return nil, err
	}
	return oqs.Sign(digestHash, seckey)
}

// VerifySignature checks that the given public key created signature over digest.
// The public key should be in compressed (33 bytes) or uncompressed (65 bytes) format.
// The signature should have the 64 byte [R || S] format.
func VerifySignature(pubkey, digestHash, signature []byte) bool {
	return oqs.VerifySignature(pubkey, digestHash, signature)
}

// DecompressPubkey parses a public key in the 33-byte compressed format.
func DecompressPubkey(pubkey []byte) (*oqs.PublicKey, error) {
	n, err := oqs.DecompressPubkey(pubkey)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &oqs.PublicKey{N: n}, nil
}

// CompressPubkey encodes a public key to the 33-byte compressed format.
func CompressPubkey(pubkey *oqs.PublicKey) []byte {
	return oqs.CompressPubkey(pubkey.N)
}