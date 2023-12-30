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

package types

import (
	"fmt"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"math/big"
	"testing"

	"github.com/DogeProtocol/dp/common"
)

func TestChainId(t *testing.T) {
	key, _ := defaultTestKey()

	tx := NewTransaction(0, common.Address{}, new(big.Int), 0, new(big.Int), nil)

	var err error
	tx, err = SignTx(tx, NewLondonSigner(big.NewInt(DEFAULT_CHAIN_ID)), key)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Sender(NewLondonSigner(big.NewInt(2)), tx)
	if err != ErrInvalidChainId {
		t.Error("expected error:", ErrInvalidChainId)
	}

	_, err = Sender(NewLondonSigner(big.NewInt(DEFAULT_CHAIN_ID)), tx)
	if err != nil {
		t.Error("expected no error")
	}
}

func TestHash(t *testing.T) {
	to := common.BytesToAddress([]byte{1})
	accesses := AccessList{{Address: to, StorageKeys: []common.Hash{{0}}}}
	accesses2 := AccessList{{Address: to, StorageKeys: []common.Hash{{2}}}}
	s := big.NewInt(5)
	r := big.NewInt(6)
	v := big.NewInt(7)

	chainId := big.NewInt(DEFAULT_CHAIN_ID)
	innerTx := DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      1,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       []byte{1, 2, 3},
		Gas:        10,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    []byte{2},
		AccessList: accesses,
		V:          v,
		R:          r,
		S:          s,
	}

	tx := NewTx(&innerTx)

	signer := LatestSignerForChainID(chainId)
	keypair, err := cryptobase.SigAlg.GenerateKey()
	if err != nil {
		t.Fatalf("failed")
	}

	_, err = SignTx(tx, signer, keypair)
	if err != nil {
		t.Fatalf("failed")
	}
	origHash, err := signer.Hash(tx)
	if err != nil {
		t.Fatalf("failed")
	}

	//Chain ID change
	innerTx1 := DefaultFeeTx{
		ChainID:    big.NewInt(2),
		Nonce:      1,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       []byte{1, 2, 3},
		Gas:        10,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    []byte{2},
		AccessList: accesses,
		V:          v,
		R:          r,
		S:          s,
	}

	tx1 := NewTx(&innerTx1)
	gotHash, err := signer.Hash(tx1)
	if err == nil {
		t.Fatalf("failed")
	}

	//Nonce change
	innerTx1 = DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      2,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       []byte{1, 2, 3},
		Gas:        10,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    []byte{2},
		AccessList: accesses,
		V:          v,
		R:          r,
		S:          s,
	}

	tx1 = NewTx(&innerTx1)
	gotHash, err = signer.Hash(tx1)
	if err != nil {
		t.Fatalf("failed")
	}

	if gotHash.IsEqualTo(origHash) {
		fmt.Println("gotHash", gotHash, "origHash", origHash)
		t.Fatalf("failed")
	}

	//To address change
	to1 := common.BytesToAddress([]byte{20})
	innerTx1 = DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      1,
		To:         &to1,
		Value:      big.NewInt(100),
		Data:       []byte{1, 2, 3},
		Gas:        10,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    []byte{2},
		AccessList: accesses,
		V:          v,
		R:          r,
		S:          s,
	}

	tx1 = NewTx(&innerTx1)
	gotHash, err = signer.Hash(tx1)
	if err != nil {
		t.Fatalf("failed")
	}

	if gotHash.IsEqualTo(origHash) {
		fmt.Println("gotHash", gotHash, "origHash", origHash)
		t.Fatalf("failed")
	}

	//Value change
	innerTx1 = DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      1,
		To:         &to,
		Value:      big.NewInt(500),
		Data:       []byte{1, 2, 3},
		Gas:        10,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    []byte{2},
		AccessList: accesses,
		V:          v,
		R:          r,
		S:          s,
	}

	tx1 = NewTx(&innerTx1)
	gotHash, err = signer.Hash(tx1)
	if err != nil {
		t.Fatalf("failed")
	}

	if gotHash.IsEqualTo(origHash) {
		fmt.Println("gotHash", gotHash, "origHash", origHash)
		t.Fatalf("failed")
	}

	//Data change
	innerTx1 = DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      1,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       []byte{1, 2, 3, 4, 5},
		Gas:        10,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    []byte{2},
		AccessList: accesses,
		V:          v,
		R:          r,
		S:          s,
	}

	tx1 = NewTx(&innerTx1)
	gotHash, err = signer.Hash(tx1)
	if err != nil {
		t.Fatalf("failed")
	}

	if gotHash.IsEqualTo(origHash) {
		fmt.Println("gotHash", gotHash, "origHash", origHash)
		t.Fatalf("failed")
	}

	//Data nil
	innerTx1 = DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      1,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       nil,
		Gas:        10,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    []byte{2},
		AccessList: accesses,
		V:          v,
		R:          r,
		S:          s,
	}

	tx1 = NewTx(&innerTx1)
	gotHash, err = signer.Hash(tx1)
	if err != nil {
		t.Fatalf("failed")
	}

	if gotHash.IsEqualTo(origHash) {
		fmt.Println("gotHash", gotHash, "origHash", origHash)
		t.Fatalf("failed")
	}

	//Gas change
	innerTx1 = DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      1,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       []byte{1, 2, 3},
		Gas:        20,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    []byte{2},
		AccessList: accesses,
		V:          v,
		R:          r,
		S:          s,
	}

	tx1 = NewTx(&innerTx1)
	gotHash, err = signer.Hash(tx1)
	if err != nil {
		t.Fatalf("failed")
	}

	if gotHash.IsEqualTo(origHash) {
		fmt.Println("gotHash", gotHash, "origHash", origHash)
		t.Fatalf("failed")
	}

	//Gas tier change
	innerTx1 = DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      1,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       []byte{1, 2, 3},
		Gas:        10,
		MaxGasTier: GAS_TIER_2X,
		Remarks:    []byte{2},
		AccessList: accesses,
		V:          v,
		R:          r,
		S:          s,
	}

	tx1 = NewTx(&innerTx1)
	gotHash, err = signer.Hash(tx1)
	if err != nil {
		t.Fatalf("failed")
	}

	if gotHash.IsEqualTo(origHash) {
		fmt.Println("gotHash", gotHash, "origHash", origHash)
		t.Fatalf("failed")
	}

	//Remarks change
	innerTx1 = DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      1,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       []byte{1, 2, 3},
		Gas:        10,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    []byte{2, 3},
		AccessList: accesses,
		V:          v,
		R:          r,
		S:          s,
	}

	tx1 = NewTx(&innerTx1)
	gotHash, err = signer.Hash(tx1)
	if err != nil {
		t.Fatalf("failed")
	}

	if gotHash.IsEqualTo(origHash) {
		fmt.Println("gotHash", gotHash, "origHash", origHash)
		t.Fatalf("failed")
	}

	//Remarks nil
	innerTx1 = DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      1,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       []byte{1, 2, 3},
		Gas:        10,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    nil,
		AccessList: accesses,
		V:          v,
		R:          r,
		S:          s,
	}

	tx1 = NewTx(&innerTx1)
	gotHash, err = signer.Hash(tx1)
	if err != nil {
		t.Fatalf("failed")
	}

	if gotHash.IsEqualTo(origHash) {
		fmt.Println("gotHash", gotHash, "origHash", origHash)
		t.Fatalf("failed")
	}

	//Access list change
	innerTx1 = DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      1,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       []byte{1, 2, 3},
		Gas:        10,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    []byte{2},
		AccessList: accesses2,
		V:          v,
		R:          r,
		S:          s,
	}

	tx1 = NewTx(&innerTx1)
	gotHash, err = signer.Hash(tx1)
	if err != nil {
		t.Fatalf("failed")
	}

	if gotHash.IsEqualTo(origHash) {
		fmt.Println("gotHash", gotHash, "origHash", origHash)
		t.Fatalf("failed")
	}

	//V,R,S change
	innerTx1 = DefaultFeeTx{
		ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
		Nonce:      1,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       []byte{1, 2, 3},
		Gas:        10,
		MaxGasTier: GAS_TIER_DEFAULT,
		Remarks:    []byte{2},
		AccessList: accesses,
	}

	tx1 = NewTx(&innerTx1)
	gotHash, err = signer.Hash(tx1)
	if err != nil {
		t.Fatalf("failed")
	}

	if gotHash.IsEqualTo(origHash) == false {
		fmt.Println("gotHash", gotHash, "origHash", origHash)
		t.Fatalf("failed")
	}
}
