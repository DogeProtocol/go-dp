package backupmanager

import (
	"crypto/rand"
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/trie"
	"math/big"
	"testing"
)

// Returns a random hash
func randHash() common.Hash {
	var h common.Hash
	rand.Read(h[:])
	return h
}

func randAddress() common.Address {
	var a common.Address
	rand.Read(a[:])
	return a
}

// largeNumber returns a very large big.Int.
func largeNumber(megabytes int) *big.Int {
	buf := make([]byte, megabytes*1024*1024)
	rand.Read(buf)
	bigint := new(big.Int)
	bigint.SetBytes(buf)
	return bigint
}

func TestBackup(t *testing.T) {
	tmpdir := t.TempDir()

	bm, err := NewBackupManager(tmpdir)
	if err != nil {
		fmt.Println("err", err)
		t.Fatalf("failed NewBackupManager")
	}

	header := &types.Header{
		MixDigest:   randHash(),
		ReceiptHash: randHash(),
		TxHash:      randHash(),
		Nonce:       types.BlockNonce{},
		Extra:       []byte{},
		Bloom:       types.Bloom{},
		GasUsed:     0,
		Coinbase:    common.Address{},
		GasLimit:    0,
		Time:        1337,
		ParentHash:  randHash(),
		Root:        randHash(),
		Number:      largeNumber(2),
		Difficulty:  largeNumber(2),
	}

	var receipts []*types.Receipt
	var txs [10000]*types.Transaction

	for i := 0; i < 10000; i++ {
		to := randAddress()
		var data [16000]byte
		baseTx := types.NewDefaultFeeTransactionSimple(0, &to, big.NewInt(100), 21000, data[:])
		rawTx := types.NewTx(baseTx)
		txs[i] = rawTx
	}

	block := types.NewBlock(header, txs[:], receipts, trie.NewStackTrie(nil))
	err = bm.BackupBlock(block)
	if err != nil {
		fmt.Println("err", err)
		t.Fatalf("failed backup block 1")
	}

	err = bm.BlockExists(block.Hash())
	if err != nil {
		fmt.Println("err", err)
		t.Fatalf("failed block exists")
	}

	blockRet, err := bm.GetBlock(block.Hash())
	if err != nil {
		fmt.Println("err", err)
		t.Fatalf("failed GetBlock")
	}

	if block.IsInternalDataEqualTo(blockRet) == false {
		t.Fatalf("block comparison failed")
	}

	blockHash, err := bm.GetBlockHash(block.NumberU64())
	if err != nil {
		t.Fatalf("GetBlockHash failed")
	}
	if blockHash.IsEqualTo(block.Hash()) == false {
		t.Fatalf("block hash comparison failed")
	}

	err = bm.BackupBlock(block) //overwrite check
	if err != nil {
		fmt.Println("err", err)
		t.Fatalf("failed BackupBlock 2")
	}

	blockHash, err = bm.GetBlockHash(block.NumberU64())
	if err != nil {
		t.Fatalf("GetBlockHash failed")
	}
	if blockHash.IsEqualTo(block.Hash()) == false {
		t.Fatalf("block hash comparison failed")
	}

	for i := 0; i < 100; i++ {
		to := randAddress()
		var data [16000]byte
		baseTx := types.NewDefaultFeeTransactionSimple(0, &to, big.NewInt(100), 21000, data[:])
		rawTx := types.NewTx(baseTx)
		err = bm.BackupTransaction(rawTx)
		if err != nil {
			fmt.Println("err", err)
			t.Fatalf("failed 3")
		}

		err = bm.BackupTransaction(rawTx) //overwrite check
		if err != nil {
			fmt.Println("err", err)
			t.Fatalf("failed 3")
		}

		err = bm.TrsansactionExists(rawTx.Hash())
		if err != nil {
			fmt.Println("err", err)
			t.Fatalf("failed transaction exists")
		}
	}

	err = bm.Close()
	if err != nil {
		fmt.Println("err", err)
		t.Fatalf("failed 4")
	}
}
