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

package proofofstake

import (
	"fmt"
	"github.com/DogeProtocol/dp/accounts/abi"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core"
	"github.com/DogeProtocol/dp/core/rawdb"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/core/vm"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"github.com/DogeProtocol/dp/params"
	"github.com/DogeProtocol/dp/systemcontracts/staking"
	"math/big"
	"testing"
)

// This test case is a repro of an annoying bug that took us forever to catch.
// In Clique PoA networks (Rinkeby, Görli, etc), consecutive blocks might have
// the same state root (no block subsidy, empty block). If a node crashes, the
// chain ends up losing the recent state and needs to regenerate it from blocks
// already in the database. The bug was that processing the block *prior* to an
// empty one **also completes** the empty one, ending up in a known-block error.
func TestReimportMirroredState(t *testing.T) {
	// Initialize a Clique chain with a single signer
	privtestkey, _ := cryptobase.SigAlg.GenerateKey()
	hextestkey, _ := cryptobase.SigAlg.PrivateKeyToHex(privtestkey)

	var (
		db     = rawdb.NewMemoryDatabase()
		key, _ = cryptobase.SigAlg.HexToPrivateKey(hextestkey)
		addr   = cryptobase.SigAlg.PublicKeyToAddressNoError(&key.PublicKey)
		engine = New(params.AllProofOfStakeProtocolChanges, db, nil, common.Hash{})
	)
	genspec := &core.Genesis{
		ExtraData: make([]byte, extraVanity+common.AddressLength+extraSeal),
		Alloc: map[common.Address]core.GenesisAccount{
			addr: {Balance: big.NewInt(10000000000000000)},
		},
	}
	copy(genspec.ExtraData[extraVanity:], addr[:])
	genesis := genspec.MustCommit(db)

	// Generate a batch of blocks, each properly signed
	chain, _ := core.NewBlockChain(db, nil, params.AllProofOfStakeProtocolChanges, engine, vm.Config{}, nil, nil)
	defer chain.Stop()
	signer := types.NewLondonSigner(chain.Config().ChainID)

	blocks, _ := core.GenerateChain(params.AllProofOfStakeProtocolChanges, genesis, engine, db, 3, func(i int, block *core.BlockGen) {
		// The chain maker doesn't have access to a chain, so the difficulty will be
		// lets unset (nil). Set it here to the correct value.
		block.SetDifficulty(diffInTurn)

		// We want to simulate an empty middle block, having the same state as the
		// first one. The last is needs a state change again to force a reorg.
		if i != 1 {
			tx, err := types.SignTx(types.NewTransaction(block.TxNonce(addr), common.Address{0x00}, new(big.Int), params.TxGas, nil, nil), signer, key)
			if err != nil {
				panic(err)
			}
			block.AddTxWithChain(chain, tx)
		}
	})
	for i, block := range blocks {
		header := block.Header()
		if i > 0 {
			header.ParentHash = blocks[i-1].Hash()
		}
		header.Extra = make([]byte, extraVanity+extraSeal)
		header.Difficulty = diffInTurn

		sig, _ := cryptobase.SigAlg.Sign(SealHash(header).Bytes(), key)

		copy(header.Extra[len(header.Extra)-extraSeal:], sig)

		blocks[i] = block.WithSeal(header)
	}
	// Insert the first two blocks and make sure the chain is valid
	db = rawdb.NewMemoryDatabase()
	genspec.MustCommit(db)

	chain, _ = core.NewBlockChain(db, nil, params.AllProofOfStakeProtocolChanges, engine, vm.Config{}, nil, nil)
	defer chain.Stop()

	if _, err := chain.InsertChain(blocks[:2]); err != nil {
		t.Fatalf("failed to insert initial blocks: %v", err)
	}
	if head := chain.CurrentBlock().NumberU64(); head != 2 {
		t.Fatalf("chain head mismatch: have %d, want %d", head, 2)
	}

	// Simulate a crash by creating a new chain on top of the database, without
	// flushing the dirty states out. Insert the last block, triggering a sidechain
	// reimport.
	chain, _ = core.NewBlockChain(db, nil, params.AllProofOfStakeProtocolChanges, engine, vm.Config{}, nil, nil)
	defer chain.Stop()

	if _, err := chain.InsertChain(blocks[2:]); err != nil {
		t.Fatalf("failed to insert final block: %v", err)
	}
	if head := chain.CurrentBlock().NumberU64(); head != 3 {
		t.Fatalf("chain head mismatch: have %d, want %d", head, 3)
	}
}

func TestSealHash(t *testing.T) {
	have := SealHash(&types.Header{
		Difficulty: new(big.Int),
		Number:     new(big.Int),
		Extra:      make([]byte, 32+cryptobase.SigAlg.SignatureWithPublicKeyLength()),
	})
	want := common.HexToHash("0xbd3d1fa43fbc4c5bfcc91b179ec92e2861df3654de60468beb908ff805359e8f") //sha3
	//want := common.HexToHash("0xe28be2bd8ff4897d07cd4fbb59b291de87746ac2cf264f57b3b696c3ddf9f99b") //sha3sha256
	if have != want {
		t.Errorf("have %x, want %x", have, want)
	}
}

func TestFlattenTxnMap(t *testing.T) {
	txnList, txnAddressMap := flattenTxnMap(nil)
	if txnList != nil && txnAddressMap != nil {
		t.Fatalf("failed")
	}

	// Generate a batch of accounts to start with
	keys := make([]*signaturealgorithm.PrivateKey, 4)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = cryptobase.SigAlg.GenerateKey()
	}
	signer := types.NewLondonSignerDefaultChain()

	groups := map[common.Address]types.Transactions{}
	txnCount := 0
	overallCount := 0
	for _, key := range keys {
		addr := cryptobase.SigAlg.PublicKeyToAddressNoError(&key.PublicKey)
		txnCount = txnCount + 1
		for i := 0; i < txnCount; i++ {
			tx, _ := types.SignTx(types.NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)
			overallCount = overallCount + 1
			groups[addr] = append(groups[addr], tx)
			fmt.Println("txhash", tx.Hash(), addr)
		}
	}

	txnList, txnAddressMap = flattenTxnMap(groups)
	if txnList == nil && txnAddressMap == nil {
		t.Fatalf("failed")
	}

	if len(txnList) != overallCount {
		t.Fatalf("failed")
	}

	if len(txnAddressMap) != overallCount {
		t.Fatalf("failed")
	}

	for addr, txns := range groups {
		for _, txn := range txns {
			addrResult, ok := txnAddressMap[txn.Hash()]
			if ok == false {
				t.Fatalf("failed")
			}
			if addr.IsEqualTo(addrResult) == false {
				t.Fatalf("failed")
			}
		}
	}

	for txnhash, addr := range txnAddressMap {
		addrResult, ok := groups[addr]
		if ok == false {
			t.Fatalf("failed")
		}
		found := false
		for _, t := range addrResult {
			hash := t.Hash()
			if hash.IsEqualTo(txnhash) {
				found = true
				break
			}
		}
		if found == false {
			t.Fatalf("failed")
		}
	}

	resultMap, err := recreateTxnMap(txnList, txnAddressMap, groups)
	if err != nil {
		t.Fatalf("failed")
	}

	for k, v := range groups {
		txns, ok := resultMap[k]
		if ok == false {
			t.Fatalf("failed")
		}

		for _, t1 := range v {
			found := false
			for _, t2 := range txns {
				t2hash := t2.Hash()
				if t2hash.IsEqualTo(t1.Hash()) {
					found = true
					break
				}
			}
			if found == false {
				t.Fatalf("failed")
			}
		}
	}

	for k, v := range resultMap {
		txns, ok := groups[k]
		if ok == false {
			t.Fatalf("failed")
		}

		for _, t1 := range v {
			found := false
			for _, t2 := range txns {
				t2hash := t2.Hash()
				if t2hash.IsEqualTo(t1.Hash()) {
					found = true
					break
				}
			}
			if found == false {
				t.Fatalf("failed")
			}
		}
	}

}

func encCall(abi *abi.ABI, method string, args ...interface{}) ([]byte, error) {
	return abi.Pack(method, args...)
}

func encCallOuter(abi *abi.ABI, method string, args ...interface{}) ([]byte, error) {
	return encCall(abi, method, args...)
}

func TestPack(t *testing.T) {
	method := staking.GetContract_Method_AddDepositorSlashing()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		fmt.Println("AddDepositorSlashing abi error", err)
		t.Fatalf("failed")
	}

	// call
	slashedAmount := big.NewInt(10)
	_, err = encCallOuter(&abiData, method, ZERO_ADDRESS, slashedAmount)
	//data, err := abiData.Pack(method, depositor, slashedAmount)
	if err != nil {
		fmt.Println("Unable to pack AddDepositorSlashing", "error", err)
		t.Fatalf("failed")
	}
}

func TestRewardYearly(t *testing.T) {
	for i := 1; i <= 350; i++ {
		blockNumber := rewardStartBlock.Int64() + (blockYearly.Int64() * int64(i))
		startBlockNumber := big.NewInt(blockNumber - blockYearly.Int64())
		startReward := new(big.Int).Set(getReward(startBlockNumber))

		endBlockNumber := big.NewInt(blockNumber - 1)
		endReward := new(big.Int).Set(getReward(endBlockNumber))

		fmt.Println("Year : ", i,
			" Block Range : ", startBlockNumber, " - ", endBlockNumber,
			" Block reward range : ", startReward, " - ", endReward)
	}
}

func TestRewardBlocks(t *testing.T) {

	startBlockNumber := big.NewInt(22338000 - 1000)
	endBlockNumber := big.NewInt(22338000)
	incrementBlock := big.NewInt(1)

	for startBlockNumber.Int64() <= endBlockNumber.Int64() {
		reward := new(big.Int).Set(getReward(startBlockNumber))
		fmt.Println("Block Number : ", startBlockNumber, " reward : ", reward)
		startBlockNumber = common.SafeAddBigInt(startBlockNumber, incrementBlock)
	}

}
