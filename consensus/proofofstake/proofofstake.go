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

// Package proofofstake implements the proof-of-authority consensus engine.
package proofofstake

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/conversionutil"
	"github.com/DogeProtocol/dp/core"
	"github.com/DogeProtocol/dp/core/state"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/handler"
	"github.com/DogeProtocol/dp/internal/ethapi"
	"github.com/DogeProtocol/dp/systemcontracts/consensuscontext"
	"github.com/DogeProtocol/dp/systemcontracts/conversion"
	"github.com/DogeProtocol/dp/systemcontracts/staking"
	"github.com/DogeProtocol/dp/trie"
	"io"
	"math/big"
	"sync"
	"time"

	"github.com/DogeProtocol/dp/accounts"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/consensus"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/ethdb"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/params"
	"github.com/DogeProtocol/dp/rlp"
	"github.com/DogeProtocol/dp/rpc"
	"github.com/DogeProtocol/dp/systemcontracts/staking/stakingv2"
	lru "github.com/hashicorp/golang-lru"
)

const (
	inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	inmemorySignatures = 4096 // Number of recent block signatures to keep in memory
)

// ProofOfStake proof-of-authority protocol constants.
var (
	maxSystemBalance = new(big.Int).Mul(big.NewInt(100), big.NewInt(params.Ether))

	epochLength = uint64(30000) // Default number of blocks after which to checkpoint and reset the pending votes

	extraVanity = 32                                               // Fixed number of extra-data prefix bytes reserved for validator vanity
	extraSeal   = cryptobase.SigAlg.SignatureWithPublicKeyLength() // Fixed number of extra-data suffix bytes reserved for validator seal

	nonceAuthVote = hexutil.MustDecode("0xffffffffffffffff") // Magic nonce number to vote on adding a new validator
	nonceDropVote = hexutil.MustDecode("0x0000000000000000") // Magic nonce number to vote on removing a validator.

	diffInTurn = big.NewInt(2) // Block difficulty for in-turn signatures
	diffNoTurn = big.NewInt(1) // Block difficulty for out-of-turn signatures

	slashAmount = params.EtherToWei(big.NewInt(10))

	rewardStartBlockNumber = uint64(277204)
	slashStartBlockNumber  = uint64(1497600)

	FULL_SIGN_PROPOSAL_CUTOFF_BLOCK     = uint64(421888)
	FULL_SIGN_PROPOSAL_FREQUENCY_BLOCKS = uint64(4096)

	STAKING_CONTRACT_V2_CUTOFF_BLOCK  = FULL_SIGN_PROPOSAL_CUTOFF_BLOCK
	CONSENSUS_CONTEXT_START_BLOCK     = FULL_SIGN_PROPOSAL_CUTOFF_BLOCK
	CONSENSUS_CONTEXT_MAX_BLOCK_COUNT = uint64(512000)

	VALIDATOR_NIL_BLOCK_START_BLOCK      = STAKING_CONTRACT_V2_CUTOFF_BLOCK + 1
	BLOCK_PROPOSER_NIL_BLOCK_START_BLOCK = VALIDATOR_NIL_BLOCK_START_BLOCK + 16
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")

	// errInvalidCheckpointBeneficiary is returned if a checkpoint/epoch transition
	// block has a beneficiary set to non-zeroes.
	errInvalidCheckpointBeneficiary = errors.New("beneficiary in checkpoint block non-zero")

	// errInvalidVote is returned if a nonce value is something else that the two
	// allowed constants of 0x00..0 or 0xff..f.
	errInvalidVote = errors.New("vote nonce not 0x00..0 or 0xff..f")

	// errInvalidCheckpointVote is returned if a checkpoint/epoch transition block
	// has a vote nonce set to non-zeroes.
	errInvalidCheckpointVote = errors.New("vote nonce in checkpoint block non-zero")

	// errMissingVanity is returned if a block's extra-data section is shorter than
	// 32 bytes, which is required to store the signer vanity.
	errMissingVanity = errors.New("extra-data 32 byte vanity prefix missing")

	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")

	// errMismatchingEpochValidators is returned if a sprint block contains a
	// list of filteredValidatorsDepositMap different than the one the local node calculated.
	errMismatchingEpochValidators = errors.New("mismatching validator list on epoch block")

	// errExtraSigners is returned if non-checkpoint block contain signer data in
	// their extra-data fields.
	errExtraSigners = errors.New("non-checkpoint block contains extra signer list")

	// errInvalidCheckpointSigners is returned if a checkpoint block contains an
	// invalid list of signers (i.e. non divisible by 20 bytes).
	errInvalidCheckpointSigners = errors.New("invalid signer list on checkpoint block")

	// errMismatchingCheckpointSigners is returned if a checkpoint block contains a
	// list of signers different than the one the local node calculated.
	errMismatchingCheckpointSigners = errors.New("mismatching signer list on checkpoint block")

	// errInvalidMixDigest is returned if a block's mix digest is non-zero.
	errInvalidMixDigest = errors.New("non-zero mix digest")

	// errInvalidDifficulty is returned if the difficulty of a block neither 1 or 2.
	errInvalidDifficulty = errors.New("invalid difficulty")

	// errWrongDifficulty is returned if the difficulty of a block doesn't match the
	// turn of the signer.
	errWrongDifficulty = errors.New("wrong difficulty")

	// errInvalidTimestamp is returned if the timestamp of a block is lower than
	// the previous block's timestamp + the minimum block period.
	errInvalidTimestamp = errors.New("invalid timestamp")

	// errInvalidVotingChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errInvalidVotingChain = errors.New("invalid voting chain")

	// errUnauthorizedSigner is returned if a header is signed by a non-authorized entity.
	errUnauthorizedSigner = errors.New("unauthorized signer")

	// errRecentlySigned is returned if a header is signed by an authorized entity
	// that already signed a header recently, thus is temporarily not allowed to.
	errRecentlySigned = errors.New("recently signed")

	// errCoinBaseMisMatch is returned if a header's coinbase do not match with signature
	errCoinBaseMisMatch = errors.New("coinbase do not match with signature")
)

// SignerFn hashes and signs the data to be signed by a backing account.
type SignerFn func(signer accounts.Account, mimeType string, message []byte) ([]byte, error)
type SignerFnWithContext func(signer accounts.Account, mimeType string, message []byte, context []byte) ([]byte, error)
type SignerTxFn func(accounts.Account, *types.Transaction, *big.Int) (*types.Transaction, error)

// ProofOfStake is the proof-of-authority consensus engine proposed to support the
// Ethereum testnet following the Ropsten attacks.
type ProofOfStake struct {
	chainConfig *params.ChainConfig        // Chain config
	config      *params.ProofOfStakeConfig // Consensus engine configuration parameters
	genesisHash common.Hash
	db          ethdb.Database // Database to store and retrieve snapshot checkpoints

	recents    *lru.ARCCache // Snapshots for recent block to speed up reorgs
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	proposals map[common.Address]bool // Current list of proposals we are pushing

	signer            types.Signer
	validator         common.Address
	signFn            SignerFn // Signer function to authorize hashes with
	signFnWithContext SignerFnWithContext
	signTxFn          SignerTxFn

	ethAPI *ethapi.PublicBlockChainAPI

	lock sync.RWMutex // Protects the validator fields

	// The fields below are for testing only
	fakeDiff bool // Skip difficulty verifications

	consensusHandler *ConsensusHandler

	account    *accounts.Account
	blockchain *core.BlockChain
}

// New creates a ProofOfStake proof-of-authority consensus engine with the initial
// signers set to the ones provided by the user.
func New(chainConfig *params.ChainConfig, db ethdb.Database,
	ethAPI *ethapi.PublicBlockChainAPI, genesisHash common.Hash) *ProofOfStake {
	// Set any missing consensus parameters to their defaults
	conf := *chainConfig

	if conf.ProofOfStake.Epoch == 0 {
		conf.ProofOfStake.Epoch = epochLength
	}
	// Allocate the snapshot caches and c.ProofOfStakereate the engine
	recents, _ := lru.NewARC(inmemorySnapshots)
	signatures, _ := lru.NewARC(inmemorySignatures)

	packetHandler := NewConsensusPacketHandler()

	proofofstake := &ProofOfStake{
		chainConfig:      chainConfig,
		config:           conf.ProofOfStake,
		genesisHash:      genesisHash,
		db:               db,
		ethAPI:           ethAPI,
		recents:          recents,
		signatures:       signatures,
		proposals:        make(map[common.Address]bool),
		signer:           types.NewLondonSigner(chainConfig.ChainID),
		consensusHandler: packetHandler,
	}

	proofofstake.consensusHandler.getValidatorsFn = proofofstake.GetValidators
	proofofstake.consensusHandler.listValidatorsFn = proofofstake.ListValidatorsAsMap
	proofofstake.consensusHandler.doesFinalizedTransactionExistFn = proofofstake.DoesFinalizedTransactionExistFn

	return proofofstake
}

func (c *ProofOfStake) SetP2PHandler(handler *handler.P2PHandler) {
	c.consensusHandler.p2pHandler = handler
}

func (c *ProofOfStake) SetBlockchain(blockchain *core.BlockChain) {
	c.blockchain = blockchain
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (c *ProofOfStake) Author(header *types.Header) (common.Address, error) {
	return ZERO_ADDRESS, nil
}

func (c *ProofOfStake) DoesFinalizedTransactionExistFn(txnHash common.Hash) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	return c.ethAPI.DoesFinalizedTransactionExist(ctx, txnHash)
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (c *ProofOfStake) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	return c.verifyHeader(chain, header, nil)
}

func flattenTxnMap(txnMap map[common.Address]types.Transactions) ([]common.Hash, map[common.Hash]common.Address) {
	if txnMap == nil {
		return nil, nil
	}

	count := 0
	for _, v := range txnMap {
		count = count + v.Len()
	}

	txnList := make([]common.Hash, count)
	txnAddressMap := make(map[common.Hash]common.Address)
	i := 0
	for k, v := range txnMap {
		for _, txn := range v {
			log.Trace("flattenTxnMap", "Hash", txn.Hash())
			txnList[i].CopyFrom(txn.Hash())
			txnAddressMap[txnList[i]] = k
			i = i + 1
		}
	}

	return txnList, txnAddressMap
}

func recreateTxnMap(selectedTxns []common.Hash, txnAddressMap map[common.Hash]common.Address, txnMap map[common.Address]types.Transactions) (map[common.Address]types.Transactions, error) {
	if selectedTxns == nil {
		return nil, nil
	}

	resultMap := make(map[common.Address]types.Transactions)
	for _, txnHash := range selectedTxns {
		addr, ok := txnAddressMap[txnHash]
		if ok == false {
			log.Debug("recreateTxnMap not fouud", "tx", txnHash)
			for k, v := range txnAddressMap {
				log.Trace("recreateTxnMap txnAddressMap", "k", k, "v", v)
			}
			return nil, errors.New("unknown transaction")
		}
		txnList, ok := txnMap[addr]
		if ok == false {
			return nil, errors.New("unknown address")
		}
		for _, txnInner := range txnList {
			hash := txnInner.Hash()
			if hash.IsEqualTo(txnHash) {
				_, ok := resultMap[addr]
				if ok == false {
					resultMap[addr] = make([]*types.Transaction, 0)
				}
				resultMap[addr] = append(resultMap[addr], txnInner)
			}
		}
	}

	return resultMap, nil
}

func (c *ProofOfStake) IsBlockReadyToSeal(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB) bool {
	blockState, _, err := c.consensusHandler.getBlockState(header.ParentHash)
	if err != nil {
		log.Trace("IsBlockReadyToSeal", "blockState", blockState, "err", err)
		return false
	}
	if blockState != BLOCK_STATE_RECEIVED_COMMITS {
		return false
	}

	return true
}

// HandleTransactions selects the transactions for including in the block according to the consensus rules.
func (c *ProofOfStake) HandleTransactions(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txnMap map[common.Address]types.Transactions) (map[common.Address]types.Transactions, error) {
	if c.signFn == nil {
		return nil, errors.New("not a miner")
	}
	txns, txnAddressMap := flattenTxnMap(txnMap)

	err := c.consensusHandler.HandleConsensus(header.ParentHash, txns, header.Number.Uint64())
	if err != nil {
		return nil, err
	}

	blockState, round, err := c.consensusHandler.getBlockState(header.ParentHash)
	if err != nil {
		log.Debug("getBlockState", "err", err)
		return nil, err
	}
	if blockState != BLOCK_STATE_RECEIVED_COMMITS {
		return nil, errors.New("not ready yet")
	}
	vote, err := c.consensusHandler.getBlockVote(header.ParentHash)
	if err != nil {
		log.Debug("getBlockVote", "err", err)
		return nil, err
	}

	selectedTxns, err := c.consensusHandler.getBlockSelectedTransactions(header.ParentHash)
	if err != nil {
		log.Debug("getBlockSelectedTransactions", "err", err)
		return nil, err
	}
	if selectedTxns == nil {
		log.Debug("getBlockSelectedTransactions nil")
		return nil, nil
	}

	log.Debug("HandleTransactions", "in txn count", len(txns), "out txn count", len(selectedTxns), "round", round, "vote", vote)
	for _, t := range txns {
		log.Trace("HandleTransactions intxns", "txn", t)
	}
	for _, t := range selectedTxns {
		log.Trace("HandleTransactions outtxns", "txn", t)
	}

	resultMap, err := recreateTxnMap(selectedTxns, txnAddressMap, txnMap)
	if err != nil {
		log.Debug("recreateTxnMap", "err", err)
		return nil, err
	}

	return resultMap, nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (c *ProofOfStake) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for i, header := range headers {
			err := c.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (c *ProofOfStake) verifyHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}

	number := header.Number.Uint64()

	// Don't waste time checking blocks from the future
	if header.Time > uint64(time.Now().Unix()) {
		return consensus.ErrFutureBlock
	}
	// Checkpoint blocks need to enforce zero beneficiary

	// Check that the extra-data contains both the vanity and signature
	if len(header.Extra) < extraVanity {
		return errMissingVanity
	}

	// Ensure that the extra-data contains a signer list on checkpoint, but none otherwise

	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != (common.Hash{}) {
		return errInvalidMixDigest
	}

	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if number > 0 {
		if header.Difficulty == nil || header.Difficulty.Uint64() != number {
			return errInvalidDifficulty
		}
	}
	// Verify that the gas limit is <= 2^63-1
	cap := uint64(0x7fffffffffffffff)
	if header.GasLimit > cap {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, cap)
	} // If all checks passed, validate any special fields for hard forks

	// All basic checks passed, verify cascading fields
	return c.verifyCascadingFields(chain, header, parents)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (c *ProofOfStake) verifyCascadingFields(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	// Ensure that the block's timestamp isn't too close to its parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}

	return c.verifySeal(chain, header, parents)
}

// verifySeal checks whether the signature contained in the header satisfies the
// consensus protocol requirements. The method accepts an optional list of parent
// headers that aren't yet part of the local blockchain to generate the snapshots
// from.
func (c *ProofOfStake) verifySeal(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	if header.ConsensusData == nil || header.UnhashedConsensusData == nil {
		log.Trace("ValidateBlockConsensusData nil")
		return errors.New("nil consensusdata")
	}

	blockConsensusData := &BlockConsensusData{}
	err := rlp.DecodeBytes(header.ConsensusData, &blockConsensusData)
	if err != nil {
		return err
	}

	blockAdditionalConsensusData := &BlockAdditionalConsensusData{}
	err = rlp.DecodeBytes(header.UnhashedConsensusData, &blockAdditionalConsensusData)
	if err != nil {
		return err
	}

	if blockConsensusData.Round < 1 {
		return errors.New("verifySeal round")
	}

	if blockConsensusData.PrecommitHash.IsEqualTo(ZERO_HASH) {
		return errors.New("ValidateBlockConsensusData PrecommitHash ProposalHash zero_hash")
	}

	if blockConsensusData.Round > 1 {
		if len(blockConsensusData.SlashedBlockProposers) < int(blockConsensusData.Round-1) {
			return errors.New("ValidateBlockConsensusData SlashedBlockProposers length")
		}
	}

	if blockConsensusData.VoteType == VOTE_TYPE_NIL {
		if blockConsensusData.BlockProposer.IsEqualTo(ZERO_ADDRESS) == false {
			return errors.New("ValidateBlockConsensusData BlockProposer false")
		}

		//todo: deep validate block proposers
	} else if blockConsensusData.VoteType == VOTE_TYPE_OK {
		if blockConsensusData.BlockProposer.IsEqualTo(ZERO_ADDRESS) {
			return errors.New("ValidateBlockConsensusData BlockProposer true")
		}
	} else {
		return errors.New("unknown VoteType")
	}
	return nil
}

func (c *ProofOfStake) PostPare(chain consensus.ChainHeaderReader, header *types.Header) error {
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header.Time = parent.Time + c.config.Period

	return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (c *ProofOfStake) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	header.Coinbase = common.Address{}
	header.Nonce = types.BlockNonce{}
	number := header.Number.Uint64()
	header.Difficulty = header.Number

	if len(header.Extra) < extraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, extraVanity-len(header.Extra))...)
	}
	header.Extra = header.Extra[:extraVanity]
	header.Extra = append(header.Extra, make([]byte, extraSeal)...)

	header.MixDigest = common.Hash{}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header.Time = parent.Time + c.config.Period

	return nil
}

func (c *ProofOfStake) VerifyBlock(chain consensus.ChainHeaderReader, block *types.Block) error {
	header := block.Header()
	number := header.Number.Uint64()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	currentNumber := uint64(c.ethAPI.BlockNumber())
	currentHeader, err := c.ethAPI.GetHeaderByNumberInner(ctx, rpc.BlockNumber(currentNumber))
	if err != nil {
		log.Trace("VerifyBlock 1", "err", err)
		return err
	}

	if number != currentNumber+1 || header.ParentHash.IsEqualTo(currentHeader.Hash()) == false {
		log.Trace("VerifyBlock 2", "err", err)
		return err
	}

	validatorDepositMap, err := c.GetValidators(header.ParentHash)
	if err != nil {
		log.Trace("VerifyBlock 3", "err", err)
		return err
	}

	var valDetailsMap map[common.Address]*ValidatorDetailsV2
	if number >= BLOCK_PROPOSER_NIL_BLOCK_START_BLOCK {
		valDetailsMap, err = c.ListValidatorsAsMap(header.ParentHash)
		if err != nil {
			return err
		}
	}

	err = ValidateBlockConsensusData(block, &validatorDepositMap, &valDetailsMap)
	if err != nil {
		log.Trace("ValidateBlockConsensusData", "err", err)
	}

	return err
}

func (c *ProofOfStake) Convert(header *types.Header, state *state.StateDB, txn *types.Transaction) error {
	msg, err := txn.AsMessage(c.signer)
	if err != nil {
		return err
	}

	log.Info("Conversion txn", "txn", txn.Hash(), "from", msg.From())

	eAddress, err := conversionutil.VerifyDataAndGetEthereumAddress(msg.From(), txn.Data())
	if err != nil {
		return err
	}

	//skip txn and proceed
	ethAddress := common.HexToAddress(eAddress) //representative address
	coins, err := c.GetCoinsForEthereumAddress(ethAddress, state, header)
	if err != nil {
		return err
	}

	if coins.Cmp(big.NewInt(0)) <= 0 {
		return nil
	}

	converted, err := c.GetConversionStatus(ethAddress, state, header)
	if err != nil {
		return err
	}

	if converted == true {
		log.Info("Conversion txn already converted, skipping", "txn", txn.Hash(), "from", msg.From())
		return nil
	}

	retCoins, err := c.SetConverted(ethAddress, msg.From(), state, header)
	if err != nil {
		return err
	}

	retQuantAddress, err := c.GetQuantumAddress(ethAddress, state, header)
	if err != nil {
		return err
	}

	log.Info("=================Conversion successful", "ethAddress", eAddress, "quantumAddress", msg.From(), "coins", coins, "retQuantAddress", retQuantAddress, "retCoins", retCoins)

	return nil
}

// Finalize implements consensus.Engine
func (c *ProofOfStake) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction) error {
	if txs == nil {
		txs = make([]*types.Transaction, 0)
	} else {
		for _, tx := range txs {
			signerHash, err := c.signer.Hash(tx)
			if err != nil {
				return err
			}
			if !tx.Verify(signerHash.Bytes()) {
				log.Trace("Txn Verify failed", "Hash", tx.Hash())
				return errors.New("Transaction verify failed")
			} else {
				log.Trace("Txn Verify ok", "Hash", tx.Hash())
			}
		}
	}

	// should not happen. Once happen, stop the node is better than broadcast the block
	if header.GasLimit < header.GasUsed {
		return errors.New("gas consumption of system txs exceed the gas limit")
	}

	blockConsensusData := &BlockConsensusData{}
	err := rlp.DecodeBytes(header.ConsensusData, &blockConsensusData)
	if err != nil {
		return err
	}

	//Conversions
	if blockConsensusData.VoteType == VOTE_TYPE_OK && txs != nil {
		for _, txn := range txs {
			if txn.To().IsEqualTo(conversion.CONVERSION_CONTRACT_ADDRESS) {
				err = c.Convert(header, state, txn)
				if err != nil {
					log.Info("Convert error", "err", err)
					return err
				}
			}
		}
	}

	//Block Slashing
	//If Round = 1, then it means PROPOSER was likely offline, as opposed to Round = 2 which means validators were not able to get consensus on time
	if blockConsensusData.Round == 1 && blockConsensusData.SlashedBlockProposers != nil && len(blockConsensusData.SlashedBlockProposers) > 0 && header.Number.Uint64() >= slashStartBlockNumber {
		for _, val := range blockConsensusData.SlashedBlockProposers {
			depositor, err := c.GetDepositorOfValidator(val, header.ParentHash)
			if err != nil {
				return err
			}
			log.Trace("depositor slashing", "depositor", depositor)
			slashTotal, err := c.AddDepositorSlashing(header.ParentHash, depositor, slashAmount, state, header)
			if err != nil {
				log.Trace("AddDepositorSlashing err", "err", err)
				return err
			}
			log.Trace("slashed amount", "slashTotal", slashTotal, "slashAmount", slashAmount, "depositor", depositor)

			if c.signFn != nil && val.IsEqualTo(c.validator) {
				log.Warn("You account got a slashing!", "parentHash", header.ParentHash)
			}
		}
	}

	//Validator nil block
	//If Round = 1, then it means PROPOSER was likely offline, as opposed to Round = 2 which means validators were not able to get consensus on time
	if blockConsensusData.VoteType == VOTE_TYPE_NIL && blockConsensusData.Round == 1 && blockConsensusData.SlashedBlockProposers != nil &&
		len(blockConsensusData.SlashedBlockProposers) > 0 && header.Number.Uint64() >= VALIDATOR_NIL_BLOCK_START_BLOCK {
		for _, val := range blockConsensusData.SlashedBlockProposers {
			err = c.SetNilBlock(val, state, header)
			if err != nil {
				log.Error("SetNilBlock err", "err", err)
				return err
			}
		}
	}

	//Block Rewards
	if blockConsensusData.VoteType == VOTE_TYPE_OK && header.Number.Uint64() >= rewardStartBlockNumber {
		blockProposerRewardAmount := GetReward(header.Number)

		//Add same amount of reward to Staking Contract, so that it is available for withdrawal later on
		err := c.accumulateBalance(state, blockProposerRewardAmount, common.HexToAddress(staking.GetStakingContract_Address_String()))
		if err != nil {
			log.Error("accumulateBalance staking contract err", "err", err)
			return err
		}

		//Get depositor of validator
		depositor, err := c.GetDepositorOfValidator(blockConsensusData.BlockProposer, header.ParentHash)
		if err != nil {
			log.Error("GetDepositorOfValidator", "err", err)
			return err
		}

		//Update staking contract with reward details
		blockProposerRewardAmountTotal, err := c.AddDepositorReward(header.ParentHash, depositor, blockProposerRewardAmount, state, header)
		if err != nil {
			log.Error("AddDepositorReward err", "err", err)
			return err
		}
		log.Trace("Reward amount", "BlockNumber", header.Number, "blockProposerRewardAmountTotal", blockProposerRewardAmountTotal, "blockProposerRewardAmount", blockProposerRewardAmount, "BlockProposer", blockConsensusData.BlockProposer)

		if blockConsensusData.VoteType == VOTE_TYPE_OK && c.signFn != nil && blockConsensusData.BlockProposer.IsEqualTo(c.validator) {
			log.Info("You potentially proposed and mined a new block!", "BlockNumber", header.Number, "parentHash", header.ParentHash)
		}

		//Validator nil block reset
		if header.Number.Uint64() > VALIDATOR_NIL_BLOCK_START_BLOCK {
			err = c.ResetNilBlock(blockConsensusData.BlockProposer, state, header)
			if err != nil {
				log.Error("ResetNilBlock err", "err", err)
				return err
			}
		}
	}

	//Staking V2
	if header.Number.Uint64() == STAKING_CONTRACT_V2_CUTOFF_BLOCK {
		log.Info("Setting stakingv2 contract code", "blockNumber", STAKING_CONTRACT_V2_CUTOFF_BLOCK)
		stakingContractCode := common.FromHex(stakingv2.STAKING_RUNTIME_BIN)
		state.SetCode(staking.STAKING_CONTRACT_ADDRESS, stakingContractCode)
	}

	//Consensus Context
	if header.Number.Uint64() == CONSENSUS_CONTEXT_START_BLOCK {
		log.Info("Setting consensus context contract code", "blockNumber", CONSENSUS_CONTEXT_START_BLOCK)
		consensuscontextContractCode := common.FromHex(consensuscontext.CONSENSUS_CONTEXT_RUNTIME_BIN)
		state.SetCode(consensuscontext.CONSENSUS_CONTEXT_CONTRACT_ADDRESS, consensuscontextContractCode)
	}

	if header.Number.Uint64() > CONSENSUS_CONTEXT_START_BLOCK {
		key, err := GetBlockConsensusContextKey(header.Number.Uint64())
		if err != nil {
			log.Error("GetBlockConsensusContext err", "err", err)
			return err
		}
		var consensuscontext [32]byte
		copy(consensuscontext[:], header.ParentHash.Bytes())
		err = c.SetConsensusContext(key, consensuscontext, state, header)
		if err != nil {
			log.Error("SetConsensusContext err", "err", err)
			return err
		}

		//Remove the oldest key
		if header.Number.Uint64() > (CONSENSUS_CONTEXT_START_BLOCK + CONSENSUS_CONTEXT_MAX_BLOCK_COUNT) {
			oldKey, err := GetBlockConsensusContextKey(header.Number.Uint64() - CONSENSUS_CONTEXT_MAX_BLOCK_COUNT)
			if err != nil {
				log.Error("GetBlockConsensusContextKey oldKey err", "err", err)
				return err
			}

			err = c.DeleteConsensusContext(oldKey, state, header)
			if err != nil {
				log.Error("DeleteConsensusContext oldKey err", "err", err)
				return err
			}
		}
	}

	//Fix blocktime
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if (header.Number.Uint64() == 1 || header.Number.Uint64()%BLOCK_PERIOD_TIME_CHANGE == 0) && blockConsensusData.VoteType == VOTE_TYPE_OK && parent.Time < blockConsensusData.BlockTime {
		header.Time = blockConsensusData.BlockTime
	} else {
		header.Time = parent.Time + c.config.Period
	}

	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))

	log.Info("Finalize Block", "root", header.Root, "hash", header.Hash(), "number", header.Number, "txn count", len(txs))

	return nil
}

func (c *ProofOfStake) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
	err := c.Finalize(chain, header, state, txs)
	if err != nil {
		return nil, err
	}

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, receipts, trie.NewStackTrie(nil)), nil
}

func (c *ProofOfStake) FinalizeAndAssembleWithConsensus(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
	// Sealing the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return nil, errUnknownBlock
	}

	blockState, round, err := c.consensusHandler.getBlockState(header.ParentHash)
	if err != nil {
		log.Trace("getBlockState", "err", err)
		return nil, err
	}

	if blockState != BLOCK_STATE_RECEIVED_COMMITS {
		log.Trace("FinalizeAndAssembleWithConsensus BLOCK_STATE_WAITING_FOR_COMMITS", round)
		return nil, errors.New("Block state not yet BLOCK_STATE_WAITING_FOR_COMMITS")
	}

	blockConsensusData, blockAdditionalConsensusData, err := c.consensusHandler.getBlockConsensusData(header.ParentHash)
	if err != nil {
		log.Trace("getBlockConsensusData", "err", err)
		return nil, err
	}
	data, err := rlp.EncodeToBytes(blockConsensusData)
	if err != nil {
		log.Trace("EncodeToBytes blockConsensusData", "err", err)
		return nil, err
	}
	header.ConsensusData = make([]byte, len(data))
	copy(header.ConsensusData, data)

	data, err = rlp.EncodeToBytes(blockAdditionalConsensusData)
	if err != nil {
		log.Trace("EncodeToBytes blockAdditionalConsensusData", "err", err)
		return nil, err
	}
	header.UnhashedConsensusData = make([]byte, len(data))
	copy(header.UnhashedConsensusData, data)

	err = c.Finalize(chain, header, state, txs)
	if err != nil {
		return nil, err
	}

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, receipts, trie.NewStackTrie(nil)), nil
}

// Authorize injects a private key into the consensus engine to mint new blocks
// with.
func (c *ProofOfStake) Authorize(validator common.Address, signFn SignerFn, signFnWithContext SignerFnWithContext, signTxFn SignerTxFn, account accounts.Account) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.validator = validator
	c.signFn = signFn
	c.signTxFn = signTxFn
	c.signFnWithContext = signFnWithContext

	c.consensusHandler.signFn = signFn
	c.consensusHandler.signFnWithContext = signFnWithContext
	c.consensusHandler.account = account
}

// Seal implements consensus.Engine, attempting to create a sealed block using
// the local signing credentials.
func (c *ProofOfStake) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	header := block.Header()
	log.Info("Seal Block", "Hash", block.ParentHash().String(), "Number", header.Number)

	delay := time.Second * 1
	go func() {
		select {
		case <-stop:
			return
		case <-time.After(delay):
		}

		select {
		case results <- block.WithSeal(header):
		default:
			log.Warn("Sealing result is not read by miner", "sealhash", SealHash(header))
		}
	}()
	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have:
// * DIFF_NOTURN(2) if BLOCK_NUMBER % SIGNER_COUNT != SIGNER_INDEX
// * DIFF_INTURN(1) if BLOCK_NUMBER % SIGNER_COUNT == SIGNER_INDEX
func (c *ProofOfStake) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(parent.Number.Int64() + 1)
}

// SealHash returns the hash of a block prior to it being sealed.
func (c *ProofOfStake) SealHash(header *types.Header) common.Hash {
	return SealHash(header)
}

// Close implements consensus.Engine. It's a noop for proofofstake as there are no background threads.
func (c *ProofOfStake) Close() error {
	return nil
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the validator voting.
func (c *ProofOfStake) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{{
		Namespace: "proofofstake",
		Version:   "1.0",
		Service:   &API{chain: chain, proofofstake: c},
		Public:    true,
	}}
}

func (c *ProofOfStake) GetConsensusPacketHandler() *ConsensusHandler {
	return c.consensusHandler
}

// SealHash returns the hash of a block prior to it being sealed.
func SealHash(header *types.Header) (hash common.Hash) {
	buff := new(bytes.Buffer)
	encodeSigHeader(buff, header)
	hash.SetBytes(crypto.Keccak256(buff.Bytes()))
	return hash
}

func encodeSigHeader(w io.Writer, header *types.Header) {
	enc := []interface{}{
		header.ParentHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-cryptobase.SigAlg.SignatureWithPublicKeyLength()], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
	}

	if err := rlp.Encode(w, enc); err != nil {
		panic("can't encode: " + err.Error())
	}
}
