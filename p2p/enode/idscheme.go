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
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/cryptopq"
	"github.com/ethereum/go-ethereum/cryptopq/oqs"
	"io"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// List of known secure identity schemes.
var ValidSchemes = enr.SchemeMap{
	"v4": V4ID{},
}

var ValidSchemesForTesting = enr.SchemeMap{
	"v4":   V4ID{},
	"null": NullID{},
}

// v4ID is the "v4" identity scheme.
type V4ID struct{}

// SignV4 signs a record using the v4 scheme.
func SignV4(r *enr.Record, privkey *oqs.PrivateKey) error {
	// Copy r to avoid modifying it if signing fails.
	cpy := *r
	cpy.Set(enr.ID("v4"))
	cpy.Set(PqPubKey(privkey.PublicKey))

	h := sha3.NewLegacyKeccak256()
	rlp.Encode(h, cpy.AppendElements(nil))
	sig, err := cryptopq.Sign(h.Sum(nil), privkey)
	if err != nil {
		return err
	}
	////sig = sig[:len(sig)-1] // remove v
	if err = cpy.SetSig(V4ID{}, sig); err == nil {
		*r = cpy
	}
	return err
}

func (V4ID) Verify(r *enr.Record, sig []byte) error {
	var entry s256raw

	if err := r.Load(&entry); err != nil {
		return err
	} else if len(entry) != oqs.PublicKeyLen {
		return fmt.Errorf("invalid public key")
	}
	h := sha3.NewLegacyKeccak256()
	rlp.Encode(h, r.AppendElements(nil))
	if !cryptopq.VerifySignature(entry, h.Sum(nil), sig) {
		return enr.ErrInvalidSig
	}
	return nil
}

func (V4ID) NodeAddr(r *enr.Record) []byte {
	var pubkey PqPubKey
	err := r.Load(&pubkey)
	if err != nil {
		return nil
	}
	buf := make([]byte, oqs.PublicKeyLen)
	math.ReadBits(pubkey.N, buf)
	return crypto.Keccak256(buf)
}

// PqPubKey is the "secp256k1" key, which holds a public key.
type PqPubKey oqs.PublicKey

func (v PqPubKey) ENRKey() string { return "secp256k1" }


// EncodeRLP implements rlp.Encoder.
func (v PqPubKey) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, cryptopq.CompressPubkey((*oqs.PublicKey)(&v)))
}

// DecodeRLP implements rlp.Decoder.
func (v *PqPubKey) DecodeRLP(s *rlp.Stream) error {
	buf, err := s.Bytes()
	if err != nil {
		return err
	}
	pk, err := cryptopq.DecompressPubkey(buf)
	if err != nil {
		return err
	}
	*v = (PqPubKey)(*pk)
	return nil
}

// s256raw is an unparsed secp256k1 public key entry.
type s256raw []byte

func (s256raw) ENRKey() string { return "secp256k1" }

// v4CompatID is a weaker and insecure version of the "v4" scheme which only checks for the
// presence of a secp256k1 public key, but doesn't verify the signature.
type v4CompatID struct {
	V4ID
}

func (v4CompatID) Verify(r *enr.Record, sig []byte) error {
	var pubkey PqPubKey
	return r.Load(&pubkey)
}

func signV4Compat(r *enr.Record, pubkey *oqs.PublicKey) {
	r.Set((*PqPubKey)(pubkey))
	if err := r.SetSig(v4CompatID{}, []byte{}); err != nil {
		panic(err)
	}
}

// NullID is the "null" ENR identity scheme. This scheme stores the node
// ID in the record without any signature.
type NullID struct{}

func (NullID) Verify(r *enr.Record, sig []byte) error {
	return nil
}

func (NullID) NodeAddr(r *enr.Record) []byte {
	var id ID
	r.Load(enr.WithEntry("nulladdr", &id))
	return id[:]
}

func SignNull(r *enr.Record, id ID) *Node {
	r.Set(enr.ID("null"))
	r.Set(enr.WithEntry("nulladdr", id))
	if err := r.SetSig(NullID{}, []byte{}); err != nil {
		panic(err)
	}
	return &Node{r: *r, id: id}
}
