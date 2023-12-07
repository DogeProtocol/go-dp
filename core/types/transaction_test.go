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

package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/rlp"
)

// The values in those tests are from the Transaction Tests
// at github.com/ethereum/tests.

var defaultSigner = londonSigner{
	chainId: big.NewInt(DEFAULT_CHAIN_ID),
}
var (
	baseTx = NewTransaction(
		3,
		testAddr,
		big.NewInt(10),
		2000,
		big.NewInt(1),
		common.FromHex("5544"),
	)

	privtestkey, _ = cryptobase.SigAlg.GenerateKey()
	hextestkey, _  = cryptobase.SigAlg.PrivateKeyToHex(privtestkey)
	hash1, _       = defaultSigner.Hash(baseTx)
	sigtest, _     = cryptobase.SigAlg.Sign(hash1.Bytes(), privtestkey)
	hexsigtest     = hex.EncodeToString(sigtest)
	parentHash     = common.HexToHash("0xabcdbaea6a6c7c4c2dfeb977efac326af552d87")

	testAddr = common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")

	emptyTx = NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil,
	)

	rightvrsTx, _ = baseTx.WithSignature(
		NewLondonSignerDefaultChain(),
		common.Hex2Bytes(hexsigtest),
	)

	emptyEip2718Tx = NewTx(&DefaultFeeTx{
		ChainID:    big.NewInt(1),
		Nonce:      3,
		To:         &testAddr,
		Value:      big.NewInt(10),
		Gas:        25000,
		MaxGasTier: GAS_TIER_DEFAULT,
		Data:       common.FromHex("5544"),
	})

	eipSigner  = NewLondonSignerDefaultChain()
	hash2, err = eipSigner.Hash(emptyEip2718Tx)

	sigtest2, _ = cryptobase.SigAlg.Sign(hash2.Bytes(), privtestkey)
	hexsigtest2 = hex.EncodeToString(sigtest2)

	signedEip2718Tx, _ = emptyEip2718Tx.WithSignature(
		eipSigner,
		common.Hex2Bytes(hexsigtest2),
	)
)

func TestDecodeEmptyTypedTx(t *testing.T) {
	input := []byte{0x80}
	var tx Transaction
	err := rlp.DecodeBytes(input, &tx)
	if err != errEmptyTypedTx {
		t.Fatal("wrong error:", err)
	}
}

func TestTransactionSigHash(t *testing.T) {
	homestead := NewLondonSignerDefaultChain()
	hash, err := homestead.Hash(emptyTx)
	if err != nil {
		t.Fatalf("failed")
	}
	if hash != common.HexToHash("c775b99e7ad12f50d819fcd602390467e28141316969f4b57f0626f74fe3b386") {
		t.Errorf("empty transaction hash mismatch, got %x", emptyTx.Hash())
	}
	hash, err = homestead.Hash(rightvrsTx)
	if err != nil {
		t.Fatalf("failed")
	}
	if hash != common.HexToHash("fe7a79529ed5f7c3375d06b26b186a8644e0e16c373d7a12be41c62d6042b77a") {
		t.Errorf("RightVRS transaction hash mismatch, got %x", rightvrsTx.Hash())
	}
}

func TestTransactionEncode(t *testing.T) {
	txb, err := rlp.EncodeToBytes(rightvrsTx)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	should := common.FromHex(hex.EncodeToString(txb))
	if !bytes.Equal(txb, should) {
		t.Errorf("encoded RLP mismatch, got %x", txb)
	}
}

func decodeTx(data []byte) (*Transaction, error) {
	var tx Transaction
	t, err := &tx, rlp.Decode(bytes.NewReader(data), &tx)
	return t, err
}

func defaultTestKey() (*signaturealgorithm.PrivateKey, common.Address) {
	key, _ := cryptobase.SigAlg.HexToPrivateKey(hextestkey)
	addr := cryptobase.SigAlg.PublicKeyToAddressNoError(&key.PublicKey)
	return key, addr
}

func TestRecipientEmpty(t *testing.T) {
	_, addr := defaultTestKey()
	tx, err := decodeTx(common.Hex2Bytes("f8498080808080011ca09b16de9d5bdee2cf56c28d16275a4da68cd30273e2525f3959f5d62557489921a0372ebd8fb3345f7db7b5a86d42e24d36e983e259b0664ceb8c227ec9af572f3d"))
	if err != nil {
		t.Fatal(err)
	}

	from, err := Sender(NewLondonSignerDefaultChain(), tx)
	if err != nil {
		t.Fatal(err)
	}
	if addr != from {
		t.Fatal("derived address doesn't match")
	}
}

func TestRecipientNormal(t *testing.T) {
	_, addr := defaultTestKey()

	tx, err := decodeTx(common.Hex2Bytes("f85d80808094000000000000000000000000000000000000000080011ca0527c0d8f5c63f7b9f41324a7c8a563ee1190bcbf0dac8ab446291bdbf32f5c79a0552c4ef0a09a04395074dab9ed34d3fbfb843c2f2546cc30fe89ec143ca94ca6"))
	if err != nil {
		t.Fatal(err)
	}

	from, err := Sender(NewLondonSignerDefaultChain(), tx)
	if err != nil {
		t.Fatal(err)
	}
	if addr != from {
		t.Fatal("derived address doesn't match")
	}
}

// Tests that if multiple transactions have the same price, the ones seen earlier
// are prioritized to avoid network spam attacks aiming for a specific ordering.
func TestTransactionSort(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*signaturealgorithm.PrivateKey, 5)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = cryptobase.SigAlg.GenerateKey()
	}
	signer := NewLondonSignerDefaultChain()

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}
	overallCount := 0
	for start, key := range keys {
		addr := cryptobase.SigAlg.PublicKeyToAddressNoError(&key.PublicKey)

		for i := 0; i < 5; i++ {
			tx, _ := SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)
			tx.time = time.Unix(0, int64(len(keys)-start))
			overallCount = overallCount + 1
			groups[addr] = append(groups[addr], tx)
			fmt.Println("txhash", tx.Hash(), addr)
		}
	}
	// Sort the transactions and cross check the nonce ordering
	parentHash := common.BytesToHash([]byte("test parent hash"))
	txset := NewTransactionsByNonce(signer, groups, parentHash)

	count := 0
	ok := txset.NextCursor()
	for ok == true {
		txn := txset.PeekCursor()
		from, _ := Sender(signer, txn)
		fmt.Println("Cursor", txn.Hash(), from, txn.Nonce())
		ok = txset.NextCursor()
		count = count + 1
	}
	if count != overallCount {
		t.Errorf("test count failed")
	}
	fmt.Println("count", count)
}

func TestTransactionSortIncreasing(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*signaturealgorithm.PrivateKey, 4)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = cryptobase.SigAlg.GenerateKey()
	}
	signer := NewLondonSignerDefaultChain()

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}
	txnCount := 0
	overallCount := 0
	for start, key := range keys {
		addr := cryptobase.SigAlg.PublicKeyToAddressNoError(&key.PublicKey)
		txnCount = txnCount + 1
		for i := 0; i < txnCount; i++ {
			tx, _ := SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)
			tx.time = time.Unix(0, int64(len(keys)-start))
			overallCount = overallCount + 1
			groups[addr] = append(groups[addr], tx)
			fmt.Println("txhash", tx.Hash(), addr)
		}
	}
	// Sort the transactions and cross check the nonce ordering
	parentHash := common.BytesToHash([]byte("test parent hash"))
	txset := NewTransactionsByNonce(signer, groups, parentHash)

	count := 0
	ok := txset.NextCursor()
	for ok == true {
		txn := txset.PeekCursor()
		from, _ := Sender(signer, txn)
		fmt.Println("Cursor", txn.Hash(), from, txn.Nonce())
		ok = txset.NextCursor()
		count = count + 1
	}
	if count != overallCount {
		t.Errorf("test count failed")
	}
	fmt.Println("count", count)
}

func TestTransactionSortDecreasing(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*signaturealgorithm.PrivateKey, 4)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = cryptobase.SigAlg.GenerateKey()
	}
	signer := NewLondonSignerDefaultChain()

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}
	txnCount := 4
	overallCount := 0
	for start, key := range keys {
		addr := cryptobase.SigAlg.PublicKeyToAddressNoError(&key.PublicKey)
		txnCount = txnCount - 1
		for i := 0; i < txnCount; i++ {
			tx, _ := SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)
			tx.time = time.Unix(0, int64(len(keys)-start))
			overallCount = overallCount + 1
			groups[addr] = append(groups[addr], tx)
			fmt.Println("txhash", tx.Hash(), addr)
		}
	}
	// Sort the transactions and cross check the nonce ordering
	parentHash := common.BytesToHash([]byte("test parent hash"))
	txset := NewTransactionsByNonce(signer, groups, parentHash)

	count := 0
	ok := txset.NextCursor()
	for ok == true {
		txn := txset.PeekCursor()
		from, _ := Sender(signer, txn)
		fmt.Println("Cursor", txn.Hash(), from, txn.Nonce())
		ok = txset.NextCursor()
		count = count + 1
	}
	if count != overallCount {
		t.Errorf("test count failed")
	}
	fmt.Println("count", count)
}

func TestTransactionSortIncreaseDecrease(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*signaturealgorithm.PrivateKey, 6)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = cryptobase.SigAlg.GenerateKey()
	}
	signer := NewLondonSignerDefaultChain()

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}
	txnCount := 0
	overallCount := 0
	for start, key := range keys {
		addr := cryptobase.SigAlg.PublicKeyToAddressNoError(&key.PublicKey)
		if txnCount == 2 {
			txnCount = txnCount - 1
		} else {
			txnCount = txnCount + 1
		}
		for i := 0; i < txnCount; i++ {
			tx, _ := SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)
			tx.time = time.Unix(0, int64(len(keys)-start))
			overallCount = overallCount + 1
			groups[addr] = append(groups[addr], tx)
			fmt.Println("txhash", tx.Hash(), addr)
		}
	}
	// Sort the transactions and cross check the nonce ordering
	parentHash := common.BytesToHash([]byte("test parent hash"))
	txset := NewTransactionsByNonce(signer, groups, parentHash)

	count := 0
	ok := txset.NextCursor()
	for ok == true {
		txn := txset.PeekCursor()
		from, _ := Sender(signer, txn)
		fmt.Println("Cursor", txn.Hash(), from, txn.Nonce())
		ok = txset.NextCursor()
		count = count + 1
	}
	if count != overallCount {
		t.Errorf("test count failed")
	}
	fmt.Println("count", count)
}

func TestTransactionSortSingle(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*signaturealgorithm.PrivateKey, 1)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = cryptobase.SigAlg.GenerateKey()
	}
	signer := NewLondonSignerDefaultChain()

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}
	overallCount := 0
	for start, key := range keys {
		addr := cryptobase.SigAlg.PublicKeyToAddressNoError(&key.PublicKey)
		for i := 0; i < 1; i++ {
			tx, _ := SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)
			tx.time = time.Unix(0, int64(len(keys)-start))
			overallCount = overallCount + 1
			groups[addr] = append(groups[addr], tx)
			fmt.Println("txhash", tx.Hash(), addr)
		}
	}
	// Sort the transactions and cross check the nonce ordering
	parentHash := common.BytesToHash([]byte("test parent hash"))
	txset := NewTransactionsByNonce(signer, groups, parentHash)

	count := 0
	ok := txset.NextCursor()
	for ok == true {
		txn := txset.PeekCursor()
		from, _ := Sender(signer, txn)
		fmt.Println("Cursor", txn.Hash(), from, txn.Nonce())
		ok = txset.NextCursor()
		count = count + 1
	}
	if count != overallCount {
		t.Errorf("test count failed")
	}
	fmt.Println("count", count)
}

func TestTransactionSortSingleAccount(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*signaturealgorithm.PrivateKey, 1)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = cryptobase.SigAlg.GenerateKey()
	}
	signer := NewLondonSignerDefaultChain()

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}
	txnCount := 5
	overallCount := 0
	for start, key := range keys {
		addr := cryptobase.SigAlg.PublicKeyToAddressNoError(&key.PublicKey)
		for i := 0; i < txnCount; i++ {
			tx, _ := SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)
			tx.time = time.Unix(0, int64(len(keys)-start))
			overallCount = overallCount + 1
			groups[addr] = append(groups[addr], tx)
			fmt.Println("txhash", tx.Hash(), addr)
		}
	}
	// Sort the transactions and cross check the nonce ordering
	parentHash := common.BytesToHash([]byte("test parent hash"))
	txset := NewTransactionsByNonce(signer, groups, parentHash)

	count := 0
	ok := txset.NextCursor()
	for ok == true {
		txn := txset.PeekCursor()
		from, _ := Sender(signer, txn)
		fmt.Println("Cursor", txn.Hash(), from, txn.Nonce())
		ok = txset.NextCursor()
		count = count + 1
	}
	if count != overallCount {
		t.Errorf("test count failed")
	}
	fmt.Println("count", count)
}

func TestTransactionSortNoTxns(t *testing.T) {
	signer := NewLondonSignerDefaultChain()

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}

	// Sort the transactions and cross check the nonce ordering
	parentHash := common.BytesToHash([]byte("test parent hash"))
	txset := NewTransactionsByNonce(signer, groups, parentHash)

	count := 0
	overallCount := 0
	ok := txset.NextCursor()
	for ok == true {
		txn := txset.PeekCursor()
		from, _ := Sender(signer, txn)
		fmt.Println("Cursor", txn.Hash(), from, txn.Nonce())
		ok = txset.NextCursor()
		count = count + 1
	}
	if count != overallCount {
		t.Errorf("test count failed")
	}
	fmt.Println("count", count)
}

func testTransactionNonceOrder_byCount(txnCount int, t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*signaturealgorithm.PrivateKey, 1)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = cryptobase.SigAlg.GenerateKey()
	}
	signer := NewLondonSignerDefaultChain()

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}
	overallCount := 0
	for start, key := range keys {
		addr := cryptobase.SigAlg.PublicKeyToAddressNoError(&key.PublicKey)

		txnList := make([]*Transaction, 0)
		for i := 0; i < txnCount; i++ {
			tx, _ := SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)
			tx.time = time.Unix(0, int64(len(keys)-start))
			overallCount = overallCount + 1
			txnList = append(txnList, tx)
			//groups[addr] = append(groups[addr], tx)
			//fmt.Println("txhash", tx.Hash(), addr)
		}
		for j := len(txnList) - 1; j >= 0; j-- {
			groups[addr] = append(groups[addr], txnList[j])
		}
	}
	// Sort the transactions and cross check the nonce ordering
	parentHash := common.BytesToHash([]byte("test parent hash"))
	txset := NewTransactionsByNonce(signer, groups, parentHash)

	count := 0
	ok := txset.NextCursor()
	prevNonce := uint64(0)
	for ok == true {
		txn := txset.PeekCursor()
		if txn.Nonce() < prevNonce {
			fmt.Println("failed", txn.Hash(), txn.Nonce(), prevNonce)
			t.Errorf("failed")
			t.Fatalf("failed")
		}
		prevNonce = txn.Nonce()
		from, _ := Sender(signer, txn)
		fmt.Println("Cursor", txn.Hash(), from, txn.Nonce(), prevNonce)
		ok = txset.NextCursor()
		count = count + 1
	}
	if count != overallCount {
		t.Errorf("test count failed")
	}
	fmt.Println("count", count)
}

func TestTransactionNonceOrder(t *testing.T) {
	testTransactionNonceOrder_byCount(10, t)
	testTransactionNonceOrder_byCount(1, t)
	testTransactionNonceOrder_byCount(2, t)
}

func testTransactionNonceOrder_skip_byCount(txnCount int, skipMap map[int]bool, outputCount int, t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*signaturealgorithm.PrivateKey, 1)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = cryptobase.SigAlg.GenerateKey()
	}
	signer := NewLondonSignerDefaultChain()

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}
	overallCount := 0
	for start, key := range keys {
		addr := cryptobase.SigAlg.PublicKeyToAddressNoError(&key.PublicKey)

		txnList := make([]*Transaction, 0)
		for i := 0; i < txnCount; i++ {
			tx, _ := SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)
			tx.time = time.Unix(0, int64(len(keys)-start))
			overallCount = overallCount + 1
			txnList = append(txnList, tx)
			//groups[addr] = append(groups[addr], tx)
			//fmt.Println("txhash", tx.Hash(), addr)
		}
		for j := len(txnList) - 1; j >= 0; j = j - 1 {
			_, ok := skipMap[j]
			if ok {
				continue
			}
			groups[addr] = append(groups[addr], txnList[j])
		}
	}
	// Sort the transactions and cross check the nonce ordering
	parentHash := common.BytesToHash([]byte("test parent hash"))
	txset := NewTransactionsByNonce(signer, groups, parentHash)

	count := 0
	ok := txset.NextCursor()
	prevNonce := uint64(0)
	for ok == true {
		txn := txset.PeekCursor()
		if txn.Nonce() < prevNonce {
			fmt.Println("failed", txn.Hash(), txn.Nonce(), prevNonce)
			t.Errorf("failed")
			t.Fatalf("failed")
		}
		prevNonce = txn.Nonce()
		from, _ := Sender(signer, txn)
		fmt.Println("Cursor", txn.Hash(), from, txn.Nonce(), prevNonce)
		ok = txset.NextCursor()
		count = count + 1
	}
	if count != outputCount {
		fmt.Println("count", count, outputCount)
		t.Errorf("test count failed")
	}
}

func TestTransactionNonceOrderSkip(t *testing.T) {
	testTransactionNonceOrder_skip_byCount(10, map[int]bool{1: true}, 1, t)
	testTransactionNonceOrder_skip_byCount(10, map[int]bool{5: true}, 5, t)
	testTransactionNonceOrder_skip_byCount(10, map[int]bool{0: true}, 9, t)
	testTransactionNonceOrder_skip_byCount(10, map[int]bool{9: true}, 9, t)
	testTransactionNonceOrder_skip_byCount(0, map[int]bool{0: true}, 0, t)
	testTransactionNonceOrder_skip_byCount(1, map[int]bool{0: true}, 0, t)
}

// TestTransactionCoding tests serializing/de-serializing to/from rlp and JSON.
func TestTransactionCoding(t *testing.T) {
	key, err := cryptobase.SigAlg.GenerateKey()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}
	var (
		signer    = NewLondonSigner(common.Big1)
		addr      = common.HexToAddress("0x0000000000000000000000000000000000000001")
		recipient = common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
		accesses  = AccessList{{Address: addr, StorageKeys: []common.Hash{{0}}}}
	)
	for i := uint64(0); i < 500; i++ {
		var txdata TxData
		switch i % 5 {
		case 0:
			// Legacy tx.
			txdata = &DefaultFeeTx{
				ChainID:    big.NewInt(DEFAULT_CHAIN_ID),
				Nonce:      i,
				To:         &recipient,
				Gas:        1,
				MaxGasTier: GAS_TIER_DEFAULT,
				AccessList: accesses,
				Data:       []byte("abcdef"),
			}
		}
		tx, err := SignNewTx(key, signer, txdata)
		if err != nil {
			t.Fatalf("could not sign transaction: %v", err)
		}
		// RLP
		parsedTx, err := encodeDecodeBinary(tx)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(parsedTx, tx)

		// JSON
		parsedTx, err = encodeDecodeJSON(tx)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(parsedTx, tx)
	}
}

func encodeDecodeJSON(tx *Transaction) (*Transaction, error) {
	data, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("json encoding failed: %v", err)
	}
	var parsedTx = &Transaction{}
	if err := json.Unmarshal(data, &parsedTx); err != nil {
		return nil, fmt.Errorf("json decoding failed: %v", err)
	}
	return parsedTx, nil
}

func encodeDecodeBinary(tx *Transaction) (*Transaction, error) {
	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("rlp encoding failed: %v", err)
	}
	var parsedTx = &Transaction{}
	if err := parsedTx.UnmarshalBinary(data); err != nil {
		return nil, fmt.Errorf("rlp decoding failed: %v", err)
	}
	return parsedTx, nil
}

func assertEqual(orig *Transaction, cpy *Transaction) error {
	// compare nonce, price, gaslimit, recipient, amount, payload, V, R, S
	if want, got := orig.Hash(), cpy.Hash(); want != got {
		return fmt.Errorf("parsed tx differs from original tx, want %v, got %v", want, got)
	}
	if want, got := orig.ChainId(), cpy.ChainId(); want.Cmp(got) != 0 {
		return fmt.Errorf("invalid chain id, want %d, got %d", want, got)
	}
	if orig.AccessList() != nil {
		if !reflect.DeepEqual(orig.AccessList(), cpy.AccessList()) {
			return fmt.Errorf("access list wrong!")
		}
	}
	return nil
}
