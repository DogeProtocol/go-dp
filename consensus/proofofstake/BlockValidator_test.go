package proofofstake

import (
	"crypto/rand"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/rlp"
	"github.com/DogeProtocol/dp/trie"
	"math/big"
	"runtime/debug"
	"strings"
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

func TestBlock_NilNegative(t *testing.T) {

	//Case 1
	BlockNilTest(nil, nil, t, "ValidateBlockConsensusData nil")

	//Case 2
	blockConsensusData := &BlockConsensusData{
		VoteType:              VOTE_TYPE_OK,
		SlashedBlockProposers: make([]common.Address, 0),
		Round:                 1,
		SelectedTransactions:  make([]common.Hash, 0),
	}

	BlockNilTest(blockConsensusData, nil, t, "ValidateBlockConsensusData nil")
}

func BlockNilTest(blockConsensusData *BlockConsensusData, blockAdditionalConsensusData *BlockAdditionalConsensusData, t *testing.T, expectedError string) {
	header := &types.Header{
		MixDigest:             randHash(),
		ReceiptHash:           randHash(),
		TxHash:                randHash(),
		Nonce:                 types.BlockNonce{},
		Extra:                 []byte{},
		Bloom:                 types.Bloom{},
		GasUsed:               0,
		Coinbase:              common.Address{},
		GasLimit:              0,
		Time:                  1337,
		ParentHash:            randHash(),
		Root:                  randHash(),
		Number:                largeNumber(2),
		Difficulty:            largeNumber(2),
		ConsensusData:         nil,
		UnhashedConsensusData: nil,
	}

	if blockConsensusData != nil {
		data, err := rlp.EncodeToBytes(blockConsensusData)
		if err != nil {
			t.Fatalf("EncodeToBytes failed 1")
		}
		header.ConsensusData = make([]byte, len(data))
		copy(header.ConsensusData, data)
	}

	if blockAdditionalConsensusData != nil {
		data, err := rlp.EncodeToBytes(blockAdditionalConsensusData)
		if err != nil {
			t.Fatalf("EncodeToBytes failed 2")
		}
		header.UnhashedConsensusData = make([]byte, len(data))
		copy(header.UnhashedConsensusData, data)
	}

	var receipts []*types.Receipt
	var txs [10]*types.Transaction

	for i := 0; i < 10; i++ {
		to := randAddress()
		var data [16000]byte
		baseTx := types.NewDefaultFeeTransactionSimple(0, &to, big.NewInt(100), 21000, data[:])
		rawTx := types.NewTx(baseTx)
		txs[i] = rawTx
	}

	block := types.NewBlock(header, txs[:], receipts, trie.NewStackTrie(nil))
	valMap := make(map[common.Address]*big.Int)
	err := ValidateBlockConsensusData(block, &valMap, nil)
	if err == nil || strings.Compare(err.Error(), expectedError) != 0 {
		debug.PrintStack()
		t.Fatalf("BlockNilTest failed")
	}
}
