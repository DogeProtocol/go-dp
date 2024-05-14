package proofofstake

import (
	"bytes"
	"errors"
	"github.com/DogeProtocol/dp/accounts"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/hybrideds"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"github.com/DogeProtocol/dp/eth/protocols/eth"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/node"
	"github.com/DogeProtocol/dp/params"
	"github.com/DogeProtocol/dp/rlp"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"sync"
	"time"
)

type GetValidatorsFn func(blockHash common.Hash) (map[common.Address]*big.Int, error)
type DoesFinalizedTransactionExistFn func(txnHash common.Hash) (bool, error)
type ListValidatorsAsMapFn func(blockHash common.Hash) (map[common.Address]*ValidatorDetailsV2, error)

type OutOfOrderPacket struct {
	ReceivedTime time.Time
	Packet       *eth.ConsensusPacket
}

type ConsensusHandler struct {
	account                         accounts.Account
	signFn                          SignerFn
	signFnWithContext               SignerFnWithContext
	p2pHandler                      P2PHandler
	blockStateDetailsMap            map[common.Hash]*BlockStateDetails
	outOfOrderPacketsMap            map[common.Hash][]*OutOfOrderPacket
	outerPacketLock                 sync.Mutex
	innerPacketLock                 sync.Mutex
	p2pLock                         sync.Mutex
	getValidatorsFn                 GetValidatorsFn
	listValidatorsFn                ListValidatorsAsMapFn
	doesFinalizedTransactionExistFn DoesFinalizedTransactionExistFn
	currentParentHash               common.Hash

	timeStatMap map[string]int

	nilVoteBlocks                uint64
	okVoteBlocks                 uint64
	totalTransactions            uint64
	maxTransactionsInBlock       uint64
	maxTransactionsBlockTime     int64
	initTime                     time.Time
	initialized                  bool
	packetHashLastSentMap        map[common.Hash]time.Time
	lastRequestConsensusDataTime time.Time

	lastBlockNumber           uint64
	lastBlockNumberChangeTime time.Time

	packetStats PacketStats
}

type PacketStats struct {
	TotalIncomingPacketCount uint64
}

type BlockConsensusData struct {
	BlockProposer         common.Address   `json:"blockProposer" gencodec:"required"`
	VoteType              VoteType         `json:"voteType" gencodec:"required"`
	ProposalHash          common.Hash      `json:"proposalHash" gencodec:"required"`
	PrecommitHash         common.Hash      `json:"precommitHash" gencodec:"required"`
	SlashedBlockProposers []common.Address `json:"nilvotedBlockProposers" gencodec:"required"`
	Round                 byte
	SelectedTransactions  []common.Hash `json:"selectedTransactions" gencodec:"required"` //this will be a super-set of transactions that actually got executed
	BlockTime             uint64        `json:"blockTime" gencodec:"required"`
}

type BlockAdditionalConsensusData struct {
	ConsensusPackets []eth.ConsensusPacket `json:"consensusPackets" gencodec:"required"`
	InitTime         uint64                `json:"initTime" gencodec:"required"`
}

// todo: use mono clock
var BLOCK_TIMEOUT_MS = int64(60000)
var FULL_BLOCK_TIMEOUT_MS = int64(90000)
var ACK_BLOCK_TIMEOUT_MS = 300000 //relative to start of block locally
var BLOCK_CLEANUP_TIME_MS = int64(900000)
var MAX_ROUND = byte(2)
var BROADCAST_RESEND_DELAY = int64(10000)
var BROADCAST_CLEANUP_DELAY = int64(1800000)
var CONSENSUS_DATA_REQUEST_RESEND_DELAY = int64(30000)
var STARTUP_DELAY_MS = int64(120000)
var BLOCK_PERIOD_TIME_CHANGE = uint64(64) //propose timeChanges every N blocks
var ALLOWED_TIME_SKEW_MINUTES = 3.0
var SKIP_HASH_CHECK = false
var STALE_BLOCK_WARN_TIME = int64(1800 * 1000)
var BLOCK_PROPOSER_OFFLINE_NIL_BLOCK_MULTIPLIER = uint64(2)
var BLOCK_PROPOSER_OFFLINE_MAX_DELAY_BLOCK_COUNT = uint64(1024)

type BlockRoundState byte
type VoteType byte
type ConsensusPacketType byte
type RequestConsensusDataType byte
type NewRoundReason byte

var InvalidPacketErr = errors.New("invalid packet")
var OutOfOrderPackerErr = errors.New("packet received out of order")

const (
	BLOCK_STATE_UNKNOWN                   BlockRoundState = 0
	BLOCK_STATE_WAITING_FOR_PROPOSAL      BlockRoundState = 1
	BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS BlockRoundState = 2
	BLOCK_STATE_WAITING_FOR_PRECOMMITS    BlockRoundState = 3
	BLOCK_STATE_WAITING_FOR_COMMITS       BlockRoundState = 4
	BLOCK_STATE_RECEIVED_COMMITS          BlockRoundState = 5
)

const (
	CONSENSUS_PACKET_TYPE_PROPOSE_BLOCK      ConsensusPacketType = 0
	CONSENSUS_PACKET_TYPE_ACK_BLOCK_PROPOSAL ConsensusPacketType = 1
	CONSENSUS_PACKET_TYPE_PRECOMMIT_BLOCK    ConsensusPacketType = 2
	CONSENSUS_PACKET_TYPE_COMMIT_BLOCK       ConsensusPacketType = 3

	PROPOSAL_KEY_PREFIX     = "proposal"
	ACK_PROPOSAL_KEY_PREFIX = "ackProposal"
	PRECOMMIT_KEY_PREFIX    = "precommit"
	COMMIT_KEY_PREFIX       = "commit"
)

const (
	VOTE_TYPE_OK  VoteType = 1
	VOTE_TYPE_NIL VoteType = 2
)

const (
	MIN_VALIDATORS int = 3
	MAX_VALIDATORS int = 128
)

const (
	NEW_ROUND_REASON_START                                NewRoundReason = 1
	NEW_ROUND_REASON_WAIT_ACK_BLOCK_PROPOSAL_TIMEOUT      NewRoundReason = 2
	NEW_ROUND_REASON_WAIT_ACK_BLOCK_PROPOSAL_HIGHER_ROUND NewRoundReason = 3
	NEW_ROUND_REASON_WAIT_PRECOMMIT_TIMEOUT               NewRoundReason = 4
)

const (
	GENESIS_BLOCK_HASH = "0x2c8127f13d50434052128a88c9c9f79a27d44a1145e51f6fd250b6e247369e99"
)

var (
	//Use genesis block as context, so that cryptographic state-proof using full-signature-mode can be verified using the genesis file itself (if signatures from proposal blocks of genesis validators are available as well).
	//Eventually when 70% (staked coins) of genesis proposers have full-signed proposal blocks that also contain genesis hash as part of the message, it means that the first state-proof is achieved for the chain.
	FULL_SIGN_CONTEXT = append(common.Hex2Bytes(GENESIS_BLOCK_HASH), []byte{crypto.DILITHIUM_ED25519_SPHINCS_FULL_ID}...)
)

var (
	MIN_VALIDATOR_DEPOSIT                               *big.Int       = params.EtherToWei(big.NewInt(5000000))
	MIN_BLOCK_DEPOSIT                                   *big.Int       = params.EtherToWei(big.NewInt(500000000000))
	MIN_BLOCK_TRANSACTION_WEIGHTED_PROPOSALS_PERCENTAGE *big.Int       = big.NewInt(70)
	ZERO_HASH                                           common.Hash    = common.BytesToHash([]byte{0})
	ZERO_ADDRESS                                        common.Address = common.BytesToAddress([]byte{0})
)

type BlockRoundDetails struct {
	Round byte

	proposalPacket     *eth.ConsensusPacket
	proposalAckPackets map[common.Address]*eth.ConsensusPacket
	precommitPackets   map[common.Address]*eth.ConsensusPacket
	commitPackets      map[common.Address]*eth.ConsensusPacket

	state                BlockRoundState
	precommitHash        common.Hash
	blockVoteType        VoteType
	blockProposalDetails *ProposalDetails
	proposalHash         common.Hash
	proposalTxns         []common.Hash

	validatorProposalAcks map[common.Address]*ProposalAckDetails
	validatorPrecommits   map[common.Address]*PreCommitDetails
	validatorCommits      map[common.Address]*CommitDetails
	validatorsDepositMap  map[common.Address]*big.Int
	initTime              time.Time

	selfKnownTransactions map[common.Hash]bool

	selfProposed       bool
	selfProposalPacket *eth.ConsensusPacket
	selfProposedTime   time.Time

	selfAckd                bool
	selfAckPacket           *eth.ConsensusPacket
	selfAckProposalVoteType VoteType

	selfPrecommited     bool
	selfPrecommitPacket *eth.ConsensusPacket
	precommitInitTime   time.Time

	selfCommited     bool
	selfCommitPacket *eth.ConsensusPacket

	proposer common.Address

	newRoundReason NewRoundReason
}

type BlockStateDetails struct {
	filteredValidatorsDepositMap      map[common.Address]*big.Int
	validatorDetailsMap               *map[common.Address]*ValidatorDetailsV2
	totalBlockDepositValue            *big.Int
	blockMinWeightedProposalsRequired *big.Int
	initTime                          time.Time
	blockRoundMap                     map[byte]*BlockRoundDetails
	currentRound                      byte
	parentHash                        common.Hash
	highestProposalRoundSeen          byte

	//stats
	proposalTime    int64
	ackProposalTime int64
	precommitTime   int64
	commitTime      int64
	blockNumber     uint64
}

type ProposalDetails struct {
	Txns      []common.Hash `json:"Txns" gencodec:"required"`
	Round     byte          `json:"Round" gencodec:"required"`
	BlockTime uint64        `json:"BlockTime" gencodec:"required"` //Is only valid for blocks divisible by 256. Only hour and minute should be set, rest should be zero.
}

type ProposalAckDetails struct {
	ProposalHash        common.Hash `json:"PrecommitHash" gencodec:"required"`
	ProposalAckVoteType VoteType    `json:"VoteType" gencodec:"required"`
	Round               byte        `json:"Round" gencodec:"required"`
}

type PreCommitDetails struct {
	PrecommitHash common.Hash `json:"PrecommitHash" gencodec:"required"` //Hash of txns + ProposalAckVoteType
	Round         byte        `json:"Round" gencodec:"required"`
}

type CommitDetails struct {
	CommitHash common.Hash `json:"CommitHash" gencodec:"required"` //Hash of txns + ProposalAckVoteType
	Round      byte        `json:"Round" gencodec:"required"`
}

type RequestConsensusPacketDetails struct {
	RequestProposal       bool             `json:"RequestProposal" gencodec:"required"`
	ValidatorProposalAcks []common.Address `json:"ValidatorProposalAcks" gencodec:"required"`
	ValidatorPrecommits   []common.Address `json:"ValidatorPrecommits" gencodec:"required"`
	ValidatorCommits      []common.Address `json:"ValidatorCommits" gencodec:"required"`
	Round                 byte             `json:"Round" gencodec:"required"`
}

func GetTimeStateBucket(state string, ms int64) string {
	key := state + "-"
	if ms < 1000 {
		return key + "-0s-to-1s"
	} else if ms > 1000 && ms < 10000 {
		return key + "-1s-to-10s"
	} else if ms > 10000 && ms < 30000 {
		return key + "-10s-to-30s"
	} else {
		return key + "-30s+"
	}

	return key
}

func NewConsensusPacketHandler() *ConsensusHandler {
	timeStatMap := make(map[string]int)

	timeStatMap[PROPOSAL_KEY_PREFIX+"-0s-to-1s"] = 0
	timeStatMap[PROPOSAL_KEY_PREFIX+"-1s-to-10s"] = 0
	timeStatMap[PROPOSAL_KEY_PREFIX+"-10s-to-30s"] = 0
	timeStatMap[PROPOSAL_KEY_PREFIX+"-30s+"] = 0

	timeStatMap[ACK_PROPOSAL_KEY_PREFIX+"-0s-to-1s"] = 0
	timeStatMap[ACK_PROPOSAL_KEY_PREFIX+"-1s-to-10s"] = 0
	timeStatMap[ACK_PROPOSAL_KEY_PREFIX+"-10s-to-30s"] = 0
	timeStatMap[ACK_PROPOSAL_KEY_PREFIX+"-30s+"] = 0

	timeStatMap[PRECOMMIT_KEY_PREFIX+"-0s-to-1s"] = 0
	timeStatMap[PRECOMMIT_KEY_PREFIX+"-1s-to-10s"] = 0
	timeStatMap[PRECOMMIT_KEY_PREFIX+"-10s-to-30s"] = 0
	timeStatMap[PRECOMMIT_KEY_PREFIX+"-30s+"] = 0

	timeStatMap[COMMIT_KEY_PREFIX+"-0s-to-1s"] = 0
	timeStatMap[COMMIT_KEY_PREFIX+"-1s-to-10s"] = 0
	timeStatMap[COMMIT_KEY_PREFIX+"-10s-to-30s"] = 0
	timeStatMap[COMMIT_KEY_PREFIX+"-30s+"] = 0

	return &ConsensusHandler{
		blockStateDetailsMap: make(map[common.Hash]*BlockStateDetails),
		outOfOrderPacketsMap: make(map[common.Hash][]*OutOfOrderPacket),
		timeStatMap:          timeStatMap,
	}
}

func (cph *ConsensusHandler) SetValidatorsFunction(getValidatorsFn GetValidatorsFn) {
	cph.getValidatorsFn = getValidatorsFn
}

func (cph *ConsensusHandler) SetListValidatorsFunction(listValidatorsFn ListValidatorsAsMapFn) {
	cph.listValidatorsFn = listValidatorsFn
}

func (cph *ConsensusHandler) SetTransactionsFunction(doesFinalizedTransactionExistFn DoesFinalizedTransactionExistFn) {
	cph.doesFinalizedTransactionExistFn = doesFinalizedTransactionExistFn
}

func (cph *ConsensusHandler) isValidator(parentHash common.Hash) (bool, error) {
	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	blockStateDetails, ok := cph.blockStateDetailsMap[parentHash]
	if ok == false {
		return false, errors.New("block hash not found")
	}

	_, found := blockStateDetails.filteredValidatorsDepositMap[cph.account.Address]
	return found, nil
}

func getBlockProposer(parentHash common.Hash, filteredValidatorDepositMap *map[common.Address]*big.Int, round byte, validatorDetailsMap *map[common.Address]*ValidatorDetailsV2, blockNumber uint64) (common.Address, error) {
	if blockNumber >= BLOCK_PROPOSER_NIL_BLOCK_START_BLOCK {
		return getBlockProposerV2(parentHash, validatorDetailsMap, round, blockNumber)
	}
	var proposer common.Address

	if len(*filteredValidatorDepositMap) < MIN_VALIDATORS {
		return proposer, errors.New("min validators not found")
	}

	validators := make([]common.Address, len(*filteredValidatorDepositMap))
	i := 0
	for k, _ := range *filteredValidatorDepositMap {
		validators[i].CopyFrom(k)
		log.Trace("getBlockProposer validator", "v", validators[i], "i", i)
		i = i + 1
	}

	sort.Slice(validators, func(i, j int) bool {
		vi := crypto.Keccak256Hash(parentHash.Bytes(), validators[i].Bytes(), []byte{round}).Bytes()
		vj := crypto.Keccak256Hash(parentHash.Bytes(), validators[j].Bytes(), []byte{round}).Bytes()
		return bytes.Compare(vi, vj) == -1
	})

	proposer = validators[0]
	log.Trace("getBlockProposer", "proposer", proposer, "round", round)

	return proposer, nil
}

func canPropose(valDetails *ValidatorDetailsV2, currentBlockNumber uint64) bool {
	if valDetails.LastNiLBlock.Cmp(new(big.Int)) == 0 {
		return true
	}

	slotsMissed := float64(valDetails.NilBlockCount.Uint64() / BLOCK_PROPOSER_OFFLINE_NIL_BLOCK_MULTIPLIER)
	blockDelay := uint64(math.Pow(2.0, slotsMissed))
	if blockDelay > BLOCK_PROPOSER_OFFLINE_MAX_DELAY_BLOCK_COUNT {
		blockDelay = BLOCK_PROPOSER_OFFLINE_MAX_DELAY_BLOCK_COUNT
	}

	nextProposalBlock := valDetails.LastNiLBlock.Uint64() + blockDelay
	result := currentBlockNumber >= nextProposalBlock
	log.Debug("canPropose", "LastNiLBlock", valDetails.LastNiLBlock, "NilBlockCount", valDetails.NilBlockCount,
		"slotsMissed", slotsMissed, "blockDelay", blockDelay, "nextProposalBlock", nextProposalBlock,
		"currentBlockNumber", currentBlockNumber, "canPropose", result, "validator", valDetails.Validator)
	return result
}

func getBlockProposerV2(parentHash common.Hash, validatorMap *map[common.Address]*ValidatorDetailsV2, round byte, blockNumber uint64) (common.Address, error) {
	var proposer common.Address

	if len(*validatorMap) < MIN_VALIDATORS {
		return proposer, errors.New("getBlockProposerV2 min validators not found")
	}

	selectedValMap := make(map[common.Address]*ValidatorDetailsV2)
	for valAddr, valDetails := range *validatorMap {
		if canPropose(valDetails, blockNumber) == false {
			continue
		}
		selectedValMap[valAddr] = valDetails
	}

	//If fewer proposers than MIN_VALIDATORS, then select everyone, something is wrong
	if len(selectedValMap) < MIN_VALIDATORS {
		for valAddr, valDetails := range *validatorMap {
			selectedValMap[valAddr] = valDetails
		}
	}

	validators := make([]common.Address, len(selectedValMap))
	j := 0
	for valAddr, _ := range selectedValMap {
		validators[j] = valAddr
		j = j + 1
	}

	sort.Slice(validators, func(i, j int) bool {
		vi := crypto.Keccak256Hash(parentHash.Bytes(), validators[i].Bytes(), []byte{round}).Bytes()
		vj := crypto.Keccak256Hash(parentHash.Bytes(), validators[j].Bytes(), []byte{round}).Bytes()
		return bytes.Compare(vi, vj) == -1
	})

	proposer = validators[0]
	log.Trace("getBlockProposerV2", "proposer", proposer, "round", round)

	return proposer, nil
}

func filterValidators(parentHash common.Hash, valDepMap *map[common.Address]*big.Int) (filteredValidators map[common.Address]bool, filteredDepositValue *big.Int, blockMinWeightedProposalsRequired *big.Int, err error) {
	validatorsDepositMap := *valDepMap

	totalDepositValue := big.NewInt(0)
	valCount := 0
	for val, depositValue := range validatorsDepositMap { //todo: this should be based on netBalance
		if depositValue.Cmp(MIN_VALIDATOR_DEPOSIT) == -1 {
			log.Trace("Skipping validator with low balance", "val", val, "depositValue", depositValue)
			delete(validatorsDepositMap, val)
			continue
		}
		totalDepositValue = common.SafeAddBigInt(totalDepositValue, depositValue)
		valCount = valCount + 1
	}

	if valCount < MIN_VALIDATORS {
		return nil, nil, nil, errors.New("number of validators less than minimum")
	}

	if totalDepositValue.Cmp(MIN_BLOCK_DEPOSIT) == -1 {
		return nil, nil, nil, errors.New("min block deposit not met")
	}

	filteredValidators = make(map[common.Address]bool)

	if len(validatorsDepositMap) <= MAX_VALIDATORS {
		for validator := range validatorsDepositMap {
			filteredValidators[validator] = true
		}
	} else {
		rng, err := cryptobase.DRNG.InitializeWithSeed(parentHash)
		if err != nil {
			return nil, nil, nil, err
		}

		zero := big.NewInt(0)
		byteMax := big.NewInt(255)
		depositValueSoFar := big.NewInt(0)

		validatorList := make([]common.Address, len(validatorsDepositMap))
		ctr := 0
		for validator, _ := range validatorsDepositMap {
			validatorList[ctr] = validator
			ctr = ctr + 1
		}

		sort.Slice(validatorList, func(i, j int) bool {
			vi := crypto.Keccak256Hash(parentHash.Bytes(), validatorList[i].Bytes()).Bytes()
			vj := crypto.Keccak256Hash(parentHash.Bytes(), validatorList[j].Bytes()).Bytes()
			return bytes.Compare(vi, vj) == -1
		})

		for _, validator := range validatorList {
			depositValue := validatorsDepositMap[validator]
			randByte := big.NewInt(int64(rng.NextByte()))

			//normalize depositValue to byte-max value since random generator only returns bytes
			normalizedDepositValue := common.SafeDivBigInt(common.SafeMulBigInt(byteMax, depositValue), totalDepositValue)
			if normalizedDepositValue.Cmp(zero) < 0 || normalizedDepositValue.Cmp(byteMax) > 0 {
				return nil, nil, nil, errors.New("invalid normalizedDepositValue")
			}

			if normalizedDepositValue.Cmp(randByte) >= 0 {
				filteredValidators[validator] = true
				depositValueSoFar = common.SafeAddBigInt(depositValueSoFar, depositValue)
			}
		}

		if len(filteredValidators) < MAX_VALIDATORS || MIN_BLOCK_DEPOSIT.Cmp(depositValueSoFar) > 0 {
			for _, validator := range validatorList {
				_, ok := filteredValidators[validator]
				if ok == false {
					//this needs optimization, since validators first in the list get the benefit
					filteredValidators[validator] = true
					depositValue := validatorsDepositMap[validator]
					depositValueSoFar = common.SafeAddBigInt(depositValueSoFar, depositValue)
					if len(filteredValidators) == MAX_VALIDATORS && MIN_BLOCK_DEPOSIT.Cmp(depositValueSoFar) <= 0 {
						break
					}
				}
			}
		}
	}

	filteredDepositValue = big.NewInt(0)
	for val, _ := range filteredValidators {
		depositValue := validatorsDepositMap[val]
		filteredDepositValue = common.SafeAddBigInt(filteredDepositValue, depositValue)
	}

	if filteredDepositValue.Cmp(MIN_BLOCK_DEPOSIT) == -1 {
		return nil, nil, nil, errors.New("min block deposit not met for filteredDepositValue")
	}

	blockMinWeightedProposalsRequired = common.SafeRelativePercentageBigInt(filteredDepositValue, MIN_BLOCK_TRANSACTION_WEIGHTED_PROPOSALS_PERCENTAGE)

	return filteredValidators, filteredDepositValue, blockMinWeightedProposalsRequired, nil
}

func (cph *ConsensusHandler) initializeBlockStateIfRequired(parentHash common.Hash, blockNumber uint64) error {
	_, ok := cph.blockStateDetailsMap[parentHash]

	if ok == true {
		return nil
	}

	cph.blockStateDetailsMap[parentHash] = &BlockStateDetails{
		blockRoundMap:                make(map[byte]*BlockRoundDetails),
		filteredValidatorsDepositMap: make(map[common.Address]*big.Int),
		initTime:                     time.Now(),
		parentHash:                   parentHash,
		highestProposalRoundSeen:     0,
		blockNumber:                  blockNumber,
	}
	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	cph.lastRequestConsensusDataTime = time.Now()

	validators, err := cph.getValidatorsFn(parentHash)
	if err != nil {
		log.Error("getValidatorsFn", "err", err)
		delete(cph.blockStateDetailsMap, parentHash)
		return err
	}

	var filteredValidators map[common.Address]bool
	filteredValidators, blockStateDetails.totalBlockDepositValue, blockStateDetails.blockMinWeightedProposalsRequired, err = filterValidators(parentHash, &validators)
	if err != nil {
		delete(cph.blockStateDetailsMap, parentHash)
		return err
	}

	if blockNumber >= BLOCK_PROPOSER_NIL_BLOCK_START_BLOCK {
		validatorDetailsMap, err := cph.listValidatorsFn(parentHash)
		if err != nil {
			log.Error("listValidatorsFn", "err", err)
			return err
		}
		for valAddr, valDetails := range validatorDetailsMap {
			if valDetails.IsValidationPaused {
				delete(validatorDetailsMap, valAddr)
				continue
			}
			_, ok := filteredValidators[valAddr]
			if ok == false {
				delete(validatorDetailsMap, valAddr)
			}
		}
		blockStateDetails.validatorDetailsMap = &validatorDetailsMap
	}

	if blockStateDetails.totalBlockDepositValue.Cmp(MIN_BLOCK_DEPOSIT) == -1 {
		delete(cph.blockStateDetailsMap, parentHash)
		return errors.New("min block deposit not met")
	}

	for addr, _ := range filteredValidators {
		depositValue := validators[addr]
		blockStateDetails.filteredValidatorsDepositMap[addr] = depositValue
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[cph.account.Address]
	if ok == false {
		log.Error("Not a validator in this block")
	}

	log.Debug("blockStateDetails", "totalBlockDepositValue", blockStateDetails.totalBlockDepositValue,
		"blockMinWeightedProposalsRequired", blockStateDetails.blockMinWeightedProposalsRequired)

	cph.blockStateDetailsMap[parentHash] = blockStateDetails
	cph.currentParentHash = parentHash

	err = cph.initializeNewBlockRound(NEW_ROUND_REASON_START)
	if err != nil {
		delete(cph.blockStateDetailsMap, parentHash)
		return errors.New("min block deposit not met")
	}

	return cph.SaveHash(parentHash)
}

func (cph *ConsensusHandler) initializeNewBlockRound(newRoundReason NewRoundReason) error {
	blockStateDetails := cph.blockStateDetailsMap[cph.currentParentHash]

	blockRoundDetails := &BlockRoundDetails{
		Round:                 blockStateDetails.currentRound + 1,
		state:                 BLOCK_STATE_WAITING_FOR_PROPOSAL,
		validatorProposalAcks: make(map[common.Address]*ProposalAckDetails),
		validatorPrecommits:   make(map[common.Address]*PreCommitDetails),
		validatorCommits:      make(map[common.Address]*CommitDetails),
		selfProposed:          false,
		selfAckd:              false,
		selfPrecommited:       false,
		initTime:              time.Now(),
		proposalAckPackets:    make(map[common.Address]*eth.ConsensusPacket),
		precommitPackets:      make(map[common.Address]*eth.ConsensusPacket),
		commitPackets:         make(map[common.Address]*eth.ConsensusPacket),
		selfKnownTransactions: make(map[common.Hash]bool),
		newRoundReason:        newRoundReason,
	}

	if blockRoundDetails.Round > 1 {
		log.Trace("initializeNewBlockRound", "currentRound", blockStateDetails.currentRound, "Address", cph.account.Address)
	}

	proposer, err := getBlockProposer(cph.currentParentHash, &blockStateDetails.filteredValidatorsDepositMap, blockRoundDetails.Round, blockStateDetails.validatorDetailsMap, blockStateDetails.blockNumber)
	if err != nil {
		return err
	}

	blockRoundDetails.proposer = proposer
	blockStateDetails.blockRoundMap[blockRoundDetails.Round] = blockRoundDetails
	blockStateDetails.currentRound = blockRoundDetails.Round

	cph.blockStateDetailsMap[cph.currentParentHash] = blockStateDetails

	return nil
}

func (cph *ConsensusHandler) isBlockProposer(parentHash common.Hash, filteredValidatorDepositMap *map[common.Address]*big.Int, round byte, blockStateDetails *BlockStateDetails) (bool, error) {
	blockProposer, err := getBlockProposer(parentHash, filteredValidatorDepositMap, round, blockStateDetails.validatorDetailsMap, blockStateDetails.blockNumber)

	if err != nil {
		log.Trace("isBlockProposer", "err", err)
		return false, err
	}
	return blockProposer.IsEqualTo(cph.account.Address), nil
}

func (cph *ConsensusHandler) HandleConsensusPacket(packet *eth.ConsensusPacket) error {
	log.Trace("HandleConsensusPacket", "ParentHash", packet.ParentHash)
	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	if packet == nil || packet.Signature == nil || packet.ConsensusData == nil || len(packet.Signature) == 0 || len(packet.ConsensusData) == 0 {
		log.Debug("HandleConsensusPacket nil")
		return errors.New("invalid packet, nil data")
	}

	if cph.signFn == nil {
		return nil
	}

	if cph.initialized == false || HasExceededTimeThreshold(cph.initTime, STARTUP_DELAY_MS) == false {
		log.Trace("received consensus packet, but consensus is not ready yet")
		return nil
	}

	cph.LogIncomingPacketStats()
	err := cph.processPacket(packet)
	if errors.Is(err, OutOfOrderPackerErr) {
		pkt := eth.NewConsensusPacket(packet)
		_, ok := cph.outOfOrderPacketsMap[packet.ParentHash]
		if ok == false {
			cph.outOfOrderPacketsMap[packet.ParentHash] = make([]*OutOfOrderPacket, 0)
		}
		oooPacket := &OutOfOrderPacket{
			ReceivedTime: time.Now(),
			Packet:       &pkt,
		}
		cph.outOfOrderPacketsMap[packet.ParentHash] = append(cph.outOfOrderPacketsMap[packet.ParentHash], oooPacket)
	}

	log.Debug("HandleConsensusPacket error", "err", err)
	return err
}

func shouldSignFull(blockNumber uint64) bool {
	if blockNumber >= FULL_SIGN_PROPOSAL_CUTOFF_BLOCK && blockNumber%FULL_SIGN_PROPOSAL_FREQUENCY_BLOCKS == 0 {
		return true
	}
	return false
}

func (cph *ConsensusHandler) processPacket(packet *eth.ConsensusPacket) error {
	if packet == nil || packet.ConsensusData == nil || len(packet.ConsensusData) < 1 || packet.Signature == nil || len(packet.Signature) < hybrideds.CRYPTO_SIGNATURE_BYTES {
		log.Debug("processPacket nil")
		return errors.New("nil packet")
	}
	packetType := ConsensusPacketType(packet.ConsensusData[0])

	dataToVerify := append(packet.ParentHash.Bytes(), packet.ConsensusData...)
	digestHash := crypto.Keccak256(dataToVerify)
	var pubKey *signaturealgorithm.PublicKey
	var err error

	if packetType == CONSENSUS_PACKET_TYPE_PROPOSE_BLOCK && len(packet.Signature) != cryptobase.SigAlg.SignatureWithPublicKeyLength() { //for verify, it is ok not to check the blockNumber for full
		pubKey, err = cryptobase.SigAlg.PublicKeyFromSignatureWithContext(digestHash, packet.Signature, FULL_SIGN_CONTEXT)
		if err != nil {
			log.Debug("processPacket invalid 1")
			return InvalidPacketErr
		}

		if cryptobase.SigAlg.VerifyWithContext(pubKey.PubData, digestHash, packet.Signature, FULL_SIGN_CONTEXT) == false {
			return InvalidPacketErr
		}
	} else {
		pubKey, err = cryptobase.SigAlg.PublicKeyFromSignature(digestHash, packet.Signature)
		if err != nil {
			log.Debug("processPacket invalid 2")
			return InvalidPacketErr
		}

		if cryptobase.SigAlg.Verify(pubKey.PubData, digestHash, packet.Signature) == false {
			log.Debug("processPacket invalid 3")
			return InvalidPacketErr
		}
	}

	validator, err := cryptobase.SigAlg.PublicKeyToAddress(pubKey)
	if err != nil {
		log.Debug("processPacket invalid 4")
		return InvalidPacketErr
	}

	log.Trace("processPacket", "validator", validator, "packetType", packetType)
	if packetType == CONSENSUS_PACKET_TYPE_PROPOSE_BLOCK {
		return cph.handleProposeBlockPacket(validator, packet, false)
	} else if packetType == CONSENSUS_PACKET_TYPE_ACK_BLOCK_PROPOSAL {
		return cph.handleAckBlockProposalPacket(validator, packet)
	} else if packetType == CONSENSUS_PACKET_TYPE_PRECOMMIT_BLOCK {
		return cph.handlePrecommitPacket(validator, packet, false)
	} else if packetType == CONSENSUS_PACKET_TYPE_COMMIT_BLOCK {
		return cph.handleCommitPacket(validator, packet, false)
	}

	log.Debug("processPacket unknown packet type")
	return errors.New("unknown packet type")
}

func (cph *ConsensusHandler) processOutOfOrderPackets(parentHash common.Hash) error {
	unprocessedPackets := make([]*OutOfOrderPacket, 0)

	for key, pktList := range cph.outOfOrderPacketsMap {
		for _, pkt := range pktList {
			if pkt.Packet.ParentHash.IsEqualTo(parentHash) {
				err := cph.processPacket(pkt.Packet)
				if err != nil {
					unprocessedPackets = append(unprocessedPackets, &OutOfOrderPacket{
						Packet:       pkt.Packet,
						ReceivedTime: pkt.ReceivedTime,
					})
				} else {
					log.Trace("processOutOfOrderPackets A", "err", err)
				}
			} else {
				log.Trace("processOutOfOrderPackets mismatch", "packet parentHash", pkt.Packet.ParentHash, "current parentHash", parentHash)
			}
		}

		if len(unprocessedPackets) == 0 {
			log.Trace("processOutOfOrderPackets none")
			delete(cph.outOfOrderPacketsMap, key)
		} else {
			log.Trace("processOutOfOrderPackets", "count", len(unprocessedPackets))
			cph.outOfOrderPacketsMap[parentHash] = unprocessedPackets
		}
	}

	return nil
}

func (cph *ConsensusHandler) getBlockRoundState(parentHash common.Hash, round byte) (blockRoundState BlockRoundState, voteType VoteType, voteCount int, err error) {
	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	blockStateDetails, ok := cph.blockStateDetailsMap[parentHash]
	if ok == false {
		return BLOCK_STATE_UNKNOWN, 0, 0, nil
	}

	if round > blockStateDetails.currentRound || round == 0 {
		return BLOCK_STATE_UNKNOWN, 0, 0, nil
	}

	roundDetails := blockStateDetails.blockRoundMap[round]

	return roundDetails.state, roundDetails.blockVoteType, len(roundDetails.validatorProposalAcks), nil
}

func (cph *ConsensusHandler) getBlockConsensusData(parentHash common.Hash) (blockConsensusData *BlockConsensusData, blockAdditionalConsensusData *BlockAdditionalConsensusData, err error) {
	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	blockStateDetails, ok := cph.blockStateDetailsMap[parentHash]
	if ok == false {
		return nil, nil, errors.New("block doesn't exist")
	}

	if blockStateDetails.currentRound == 0 {
		return nil, nil, errors.New("invalid block round")
	}

	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if blockRoundDetails.state != BLOCK_STATE_RECEIVED_COMMITS {
		return nil, nil, errors.New("block state not done commit yet")
	}

	if blockRoundDetails.blockVoteType != VOTE_TYPE_OK && blockRoundDetails.blockVoteType != VOTE_TYPE_NIL {
		log.Warn("getBlockConsensusData", "voteType", blockRoundDetails.blockVoteType)
		return nil, nil, errors.New("getBlockConsensusData invalid vote type e")
	}

	blockConsensusData = &BlockConsensusData{
		VoteType:              blockRoundDetails.blockVoteType,
		SlashedBlockProposers: make([]common.Address, 0),
		Round:                 blockStateDetails.currentRound,
		SelectedTransactions:  make([]common.Hash, 0),
	}
	if blockConsensusData.VoteType == VOTE_TYPE_OK {
		blockConsensusData.BlockProposer.CopyFrom(blockRoundDetails.proposer)
		blockConsensusData.ProposalHash.CopyFrom(blockRoundDetails.proposalHash)
		blockConsensusData.BlockTime = blockRoundDetails.blockProposalDetails.BlockTime

		if blockRoundDetails.proposalTxns != nil {
			blockConsensusData.SelectedTransactions = make([]common.Hash, len(blockRoundDetails.proposalTxns))
			for i := 0; i < len(blockRoundDetails.proposalTxns); i++ {
				blockConsensusData.SelectedTransactions[i].CopyFrom(blockRoundDetails.proposalTxns[i])
			}
		}
	} else {
		blockConsensusData.BlockProposer.CopyFrom(ZERO_ADDRESS)
		blockConsensusData.ProposalHash.CopyFrom(getNilVoteProposalHash(parentHash, blockStateDetails.currentRound))
		blockConsensusData.BlockTime = 0
	}

	blockConsensusData.PrecommitHash.CopyFrom(blockRoundDetails.precommitHash)

	blockAdditionalConsensusData = &BlockAdditionalConsensusData{
		InitTime: uint64(blockStateDetails.initTime.UnixNano() / int64(time.Millisecond)),
	}

	consensusPackets := make([]eth.ConsensusPacket, 0)

	for r := byte(1); r <= blockStateDetails.currentRound; r++ {
		roundPktCount := 0

		blockRoundDetails := blockStateDetails.blockRoundMap[r]
		if blockRoundDetails.proposalPacket != nil {
			consensusPackets = append(consensusPackets, eth.NewConsensusPacket(blockRoundDetails.proposalPacket))
			roundPktCount = roundPktCount + 1
		}

		for _, pkt := range blockRoundDetails.proposalAckPackets {
			consensusPackets = append(consensusPackets, eth.NewConsensusPacket(pkt))
		}
		for _, pkt := range blockRoundDetails.precommitPackets {
			consensusPackets = append(consensusPackets, eth.NewConsensusPacket(pkt))
		}
		for _, pkt := range blockRoundDetails.commitPackets {
			consensusPackets = append(consensusPackets, eth.NewConsensusPacket(pkt))
		}

		roundProposer, err := getBlockProposer(parentHash, &blockStateDetails.filteredValidatorsDepositMap, r,
			blockStateDetails.validatorDetailsMap, blockStateDetails.blockNumber)
		if err != nil {
			return nil, nil, err
		}

		roundPktCount = roundPktCount + len(blockRoundDetails.proposalAckPackets) + len(blockRoundDetails.precommitPackets) + len(blockRoundDetails.commitPackets)
		if roundPktCount == 0 {
			log.Trace("consensusdata", "Address", cph.account.Address, "currentRound", blockStateDetails.currentRound, "r", r)
			return nil, nil, errors.New("no consensus packets for round")
		}

		if blockConsensusData.VoteType == VOTE_TYPE_NIL {
			if r < MAX_ROUND { //since MAX_ROUND is by default NIL vote
				blockConsensusData.SlashedBlockProposers = append(blockConsensusData.SlashedBlockProposers, roundProposer)
			}
		} else {
			//if VoteType is VOTE_TYPE_OK, it means that all proposers less than currentRound will be NIL VOTED (except if only one round)
			if blockStateDetails.currentRound != byte(1) && r < blockStateDetails.currentRound {
				blockConsensusData.SlashedBlockProposers = append(blockConsensusData.SlashedBlockProposers, roundProposer)
			}
		}

		log.Trace("consensusdata", "Address", cph.account.Address, "currentRound", blockStateDetails.currentRound, "r", r, "roundPktCount", roundPktCount)
	}

	blockAdditionalConsensusData.ConsensusPackets = make([]eth.ConsensusPacket, len(consensusPackets))
	for i, packet := range consensusPackets {
		blockAdditionalConsensusData.ConsensusPackets[i] = eth.NewConsensusPacket(&packet)
	}

	if blockConsensusData.VoteType == VOTE_TYPE_NIL {
		err = ValidateBlockConsensusDataInner(nil, parentHash, blockConsensusData, blockAdditionalConsensusData, &blockStateDetails.filteredValidatorsDepositMap, blockStateDetails.blockNumber, blockStateDetails.validatorDetailsMap)
	} else {
		err = ValidateBlockConsensusDataInner(blockRoundDetails.proposalTxns, parentHash, blockConsensusData, blockAdditionalConsensusData, &blockStateDetails.filteredValidatorsDepositMap, blockStateDetails.blockNumber, blockStateDetails.validatorDetailsMap)
	}

	if err != nil {
		return nil, nil, err
	}

	return blockConsensusData, blockAdditionalConsensusData, nil
}

func (cph *ConsensusHandler) getBlockRound(parentHash common.Hash, round byte) (blockRoundState *BlockRoundDetails, err error) {
	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	blockStateDetails, ok := cph.blockStateDetailsMap[parentHash]
	if ok == false {
		return nil, nil
	}

	return blockStateDetails.blockRoundMap[round], nil
}

func (cph *ConsensusHandler) getBlockState(parentHash common.Hash) (blockRoundState BlockRoundState, round byte, err error) {
	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	blockStateDetails, ok := cph.blockStateDetailsMap[parentHash]
	if ok == false {
		return BLOCK_STATE_UNKNOWN, 0, nil
	}

	return blockStateDetails.blockRoundMap[blockStateDetails.currentRound].state, blockStateDetails.currentRound, nil
}

func (cph *ConsensusHandler) getBlockSelectedTransactions(parentHash common.Hash) (txns []common.Hash, err error) {
	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	blockStateDetails, ok := cph.blockStateDetailsMap[parentHash]
	if ok == false {
		return nil, errors.New("block doesn't exist")
	}

	if blockStateDetails.currentRound == 0 {
		return nil, errors.New("invalid block round")
	}

	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if blockRoundDetails.state != BLOCK_STATE_RECEIVED_COMMITS {
		return nil, errors.New("block state not done commit yet")
	}

	if blockRoundDetails.blockVoteType == VOTE_TYPE_NIL {
		return nil, nil
	} else {
		txns = make([]common.Hash, len(blockRoundDetails.proposalTxns))
		for i := 0; i < len(blockRoundDetails.proposalTxns); i++ {
			txns[i].CopyFrom(blockRoundDetails.proposalTxns[i])
		}

		return txns, nil
	}
}

func (cph *ConsensusHandler) getBlockVote(parentHash common.Hash) (VoteType, error) {
	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	blockStateDetails, ok := cph.blockStateDetailsMap[parentHash]
	if ok == false {
		return VOTE_TYPE_NIL, nil
	}

	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	return blockRoundDetails.selfAckProposalVoteType, nil
}

func GetCombinedTxnHash(parentHash common.Hash, round byte, txns []common.Hash) common.Hash {
	var txnList []common.Hash
	txnList = make([]common.Hash, len(txns))
	for i := 0; i < len(txns); i++ {
		txnList[i].CopyFrom(txns[i])
	}

	sort.Slice(txnList, func(i, j int) bool {
		return bytes.Compare(txnList[i].Bytes(), txnList[j].Bytes()) == -1
	})

	var data []byte
	data = make([]byte, 0)
	for _, txn := range txnList {
		data = append(data, txn.Bytes()...)
	}

	hash := crypto.Keccak256Hash(data, parentHash.Bytes(), []byte{round})
	log.Trace("GetCombinedTxnHash", "parentHash", parentHash, "round", round, "txn count", len(txns), "hash", hash)
	return hash
}

func (cph *ConsensusHandler) handleProposeBlockPacket(validator common.Address, packet *eth.ConsensusPacket, self bool) error {
	cph.innerPacketLock.Lock()
	defer cph.innerPacketLock.Unlock()

	log.Trace("validator proposal", "validator", validator, "self", cph.account.Address, "hash", packet.ParentHash)
	blockStateDetails, ok := cph.blockStateDetailsMap[packet.ParentHash]
	if ok == false {
		return errors.New("unknown parentHash")
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[cph.account.Address]
	if ok == false {
		return errors.New("not a validator in this block")
	}

	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	proposalDetails := ProposalDetails{}

	err := rlp.DecodeBytes(packet.ConsensusData[1:], &proposalDetails)
	if err != nil {
		log.Trace("handleProposeTransactionsPacket8", err)
		return err
	}

	if proposalDetails.Round != blockRoundDetails.Round {
		return OutOfOrderPackerErr
	}

	if blockRoundDetails.state >= BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS {
		log.Trace("handleProposeBlockPacket BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS")
		return OutOfOrderPackerErr
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[validator]
	if ok == false {
		log.Trace("handleProposeTransactionsPacket6")
		return errors.New("invalid validator")
	}

	if blockRoundDetails.proposer.IsEqualTo(validator) == false {
		return errors.New("invalid proposer")
	}

	if ValidateBlockProposalTimeConsensus(blockStateDetails.blockNumber, proposalDetails.BlockTime) == false {
		return errors.New("block time validation failed, skipping packet")
	}

	if validator.IsEqualTo(cph.account.Address) == true && self == false {
		return errors.New("self packet from elsewhere")
	}

	if blockStateDetails.currentRound >= MAX_ROUND && len(proposalDetails.Txns) > 0 {
		return errors.New("unexpected transaction count when handling blockProposal")
	}

	proposalHash := GetCombinedTxnHash(packet.ParentHash, proposalDetails.Round, proposalDetails.Txns)
	if blockRoundDetails.proposalPacket != nil && proposalHash.IsEqualTo(blockRoundDetails.proposalHash) == false {
		return errors.New("invalid proposal hash")
	}

	blockRoundDetails.blockProposalDetails = &proposalDetails
	blockRoundDetails.proposalHash.CopyFrom(proposalHash)
	blockRoundDetails.proposalTxns = make([]common.Hash, len(proposalDetails.Txns))
	for i := 0; i < len(proposalDetails.Txns); i++ {
		exists, err := cph.doesFinalizedTransactionExistFn(proposalDetails.Txns[i])
		if err != nil {
			log.Trace("doesFinalizedTransactionExistFn", "err", err)
			return err
		}
		if exists {
			log.Trace("doesFinalizedTransactionExistFn true", proposalDetails.Txns[i].Hex())
			return errors.New("transaction already finalized")
		}
		blockRoundDetails.proposalTxns[i].CopyFrom(proposalDetails.Txns[i])
	}

	if self == false {
		//Find if any new transactions we don't know yet
		unknownTxns := make([]common.Hash, 0)
		for i := 0; i < len(proposalDetails.Txns); i++ {
			_, txnExists := blockRoundDetails.selfKnownTransactions[proposalDetails.Txns[i]]
			if txnExists == false {
				unknownTxns = append(unknownTxns, proposalDetails.Txns[i])
				log.Trace("handleProposeBlockPacket unknown", "txn", proposalDetails.Txns[i], "validator")
			} else {
				log.Trace("known txn", "txn", proposalDetails.Txns[i], "validator")
			}
		}
		if len(unknownTxns) > 0 {
			err = cph.p2pHandler.RequestTransactions(unknownTxns)
			if err != nil {
				log.Trace("RequestTransactions error", "err", err)
			}
		} else {
			blockRoundDetails.state = BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS
			blockStateDetails.proposalTime = Elapsed(blockStateDetails.initTime)
		}
	} else {
		blockStateDetails.proposalTime = Elapsed(blockStateDetails.initTime)
		blockRoundDetails.state = BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS
		blockRoundDetails.selfProposed = true
		blockRoundDetails.selfProposalPacket = packet
		blockRoundDetails.selfProposedTime = time.Now()
	}

	pkt := eth.NewConsensusPacket(packet)
	blockRoundDetails.proposalPacket = &pkt
	blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails

	cph.blockStateDetailsMap[packet.ParentHash] = blockStateDetails

	return nil
}

func (cph *ConsensusHandler) handleAckBlockProposalPacket(validator common.Address, packet *eth.ConsensusPacket) error {
	cph.innerPacketLock.Lock()
	defer cph.innerPacketLock.Unlock()

	blockStateDetails, ok := cph.blockStateDetailsMap[packet.ParentHash]
	if ok == false {
		return errors.New("unknown parentHash")
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[cph.account.Address]
	if ok == false {
		return errors.New("not a validator in this block")
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[validator]
	if ok == false {
		return errors.New("invalid validator")
	}

	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	_, ok = blockRoundDetails.validatorProposalAcks[validator]
	if ok == true {

		//todo: compare
	} else {
		if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS {

		}
	}

	proposalAckDetails := &ProposalAckDetails{}

	err := rlp.DecodeBytes(packet.ConsensusData[1:], proposalAckDetails)
	if err != nil {
		log.Trace("handleAckBlockProposalPacket err5", "err", err)
		return err
	}

	if proposalAckDetails.Round != blockStateDetails.currentRound {
		log.Trace("handleAckBlockProposalPacket", "Round", proposalAckDetails.Round, "currentRound", blockStateDetails.currentRound)
		if proposalAckDetails.Round > blockStateDetails.currentRound {
			blockStateDetails.highestProposalRoundSeen = proposalAckDetails.Round
			cph.blockStateDetailsMap[packet.ParentHash] = blockStateDetails
		}
		return OutOfOrderPackerErr
	}

	if proposalAckDetails.ProposalAckVoteType != VOTE_TYPE_OK && proposalAckDetails.ProposalAckVoteType != VOTE_TYPE_NIL {
		log.Trace("proposalAckDetails.ProposalAckVoteType", "ProposalAckVoteType", proposalAckDetails.ProposalAckVoteType)
		return errors.New("invalid vote type c")
	}

	if proposalAckDetails.Round >= MAX_ROUND && proposalAckDetails.ProposalAckVoteType != VOTE_TYPE_NIL {
		log.Trace("invalid vote type d", "validator", validator)
		return errors.New("invalid vote type, expected nil vote")
	}

	log.Trace("handleAckBlockProposalPacket blockRoundDetails", "state", blockRoundDetails.state)
	if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS || blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PRECOMMITS {
		log.Trace("handleAckBlockProposalPacket waiting")
	} else if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL {
	} else if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_COMMITS {

	} else {
		return errors.New("invalid state")
	}

	blockRoundDetails.validatorProposalAcks[validator] = proposalAckDetails

	pkt := eth.NewConsensusPacket(packet)
	blockRoundDetails.proposalAckPackets[validator] = &pkt
	blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails

	cph.blockStateDetailsMap[packet.ParentHash] = blockStateDetails

	return nil
}

func parsePacket(packet *eth.ConsensusPacket) (byte, common.Address, error) {
	dataToVerify := append(packet.ParentHash.Bytes(), packet.ConsensusData...)
	digestHash := crypto.Keccak256(dataToVerify)
	pubKey, err := cryptobase.SigAlg.PublicKeyFromSignature(digestHash, packet.Signature)
	if err != nil {
		log.Trace("invalid 1", "err", err)
		return 0, ZERO_ADDRESS, err
	}
	if cryptobase.SigAlg.Verify(pubKey.PubData, digestHash, packet.Signature) == false {
		log.Trace("invalid 2")
		return 0, ZERO_ADDRESS, InvalidPacketErr
	}

	validator, err := cryptobase.SigAlg.PublicKeyToAddress(pubKey)
	if err != nil {
		log.Trace("invalid 3", "err", err)
		return 0, ZERO_ADDRESS, err
	}

	packetType := ConsensusPacketType(packet.ConsensusData[0])
	if packetType == CONSENSUS_PACKET_TYPE_PROPOSE_BLOCK {
		details := ProposalDetails{}

		err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
		if err != nil {
			log.Trace("invalid 4", "err", err)
			return 0, ZERO_ADDRESS, err
		}

		return details.Round, validator, nil
	} else if packetType == CONSENSUS_PACKET_TYPE_ACK_BLOCK_PROPOSAL {
		details := ProposalAckDetails{}

		err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
		if err != nil {
			log.Trace("invalid 5", "err", err)
			return 0, ZERO_ADDRESS, err
		}

		return details.Round, validator, nil
	} else if packetType == CONSENSUS_PACKET_TYPE_PRECOMMIT_BLOCK {
		details := PreCommitDetails{}

		err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
		if err != nil {
			log.Trace("invalid 6", "err", err)
			return 0, ZERO_ADDRESS, err
		}

		return details.Round, validator, nil
	} else if packetType == CONSENSUS_PACKET_TYPE_COMMIT_BLOCK {
		details := CommitDetails{}

		err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
		if err != nil {
			log.Trace("invalid 7", "err", err)
			return 0, ZERO_ADDRESS, err
		}

		return details.Round, validator, nil
	}

	log.Trace("invalid 8", "err", err, "packetType", packetType)

	return 0, ZERO_ADDRESS, InvalidPacketErr
}

func (cph *ConsensusHandler) findTotalDepositsInGreaterRound(parentHash common.Hash) *big.Int {
	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	//Find deposit in greater rounds
	valMap := make(map[common.Address]bool)
	for _, pktList := range cph.outOfOrderPacketsMap {
		for _, pkt := range pktList {
			if pkt.Packet.ParentHash.IsEqualTo(parentHash) {
				round, validator, err := parsePacket(pkt.Packet)
				if err != nil {
					continue
				}
				if round <= blockStateDetails.currentRound {
					continue
				}
				_, ok := blockRoundDetails.validatorPrecommits[validator]
				if ok { //if precommit from this validator, skip counting it
					continue
				}
				valMap[validator] = true
			}
		}
	}

	totalGreaterRoundDepositCount := big.NewInt(0)
	for val, depositAmount := range blockStateDetails.filteredValidatorsDepositMap {
		exists, ok := valMap[val]
		if ok == false || exists == false {
			continue
		}
		totalGreaterRoundDepositCount = common.SafeAddBigInt(depositAmount, totalGreaterRoundDepositCount)
	}

	return totalGreaterRoundDepositCount
}

func (cph *ConsensusHandler) shouldMoveToNextRoundProposalAcks(parentHash common.Hash) (bool, error) {
	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	//Find validators in greater rounds
	valMap := make(map[common.Address]bool)
	for _, pktList := range cph.outOfOrderPacketsMap {
		for _, pkt := range pktList {
			if pkt.Packet.ParentHash.IsEqualTo(parentHash) {
				round, validator, err := parsePacket(pkt.Packet)
				if err != nil {
					log.Trace("parsePacket", "err", err)
					continue
				}
				if round <= blockStateDetails.currentRound {
					continue
				}
				valMap[validator] = true
				log.Trace("shouldMoveToNextRoundProposalAcks", "valInGreaterRound", validator)
			}
		}
	}

	totalGreaterRoundDepositCount := big.NewInt(0)
	currentRoundDepositSoFar := big.NewInt(0)
	for val, depositAmount := range blockStateDetails.filteredValidatorsDepositMap {
		_, ok := valMap[val]
		if ok == false {
			_, ok1 := blockRoundDetails.validatorProposalAcks[val]
			if ok1 == true {
				currentRoundDepositSoFar = common.SafeAddBigInt(depositAmount, currentRoundDepositSoFar)
				log.Trace("currentRoundDepositSoFar received packet", "val", val, "depositAmount", depositAmount, "currentRoundDepositSoFar", currentRoundDepositSoFar)
			} else {
				log.Trace("currentRoundDepositSoFar did not receive packet from validator", "val", val)
			}
		} else {
			totalGreaterRoundDepositCount = common.SafeAddBigInt(depositAmount, totalGreaterRoundDepositCount)
			log.Trace("totalGreaterRoundDepositCount", "val", val, "depositAmount", depositAmount, "totalGreaterRoundDepositCount", totalGreaterRoundDepositCount)
		}
	}

	if currentRoundDepositSoFar.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 {
		return false, nil
	}

	//If there are votes in greater rounds
	balanceDepositVotesRequiredCurrentRound := common.SafeSubBigInt(blockStateDetails.blockMinWeightedProposalsRequired, currentRoundDepositSoFar)
	log.Trace("shouldMoveToNextRoundProposalAcks",
		"blockMinWeightedProposalsRequired", blockStateDetails.blockMinWeightedProposalsRequired,
		"balanceDepositVotesRequiredCurrentRound", balanceDepositVotesRequiredCurrentRound,
		"currentRoundDepositSoFar", currentRoundDepositSoFar,
		"totalGreaterRoundDepositCount", totalGreaterRoundDepositCount,
		"validatorProposalAcks count", len(blockRoundDetails.validatorProposalAcks),
		"self selfProposed", blockRoundDetails.selfProposed,
		"val", cph.account.Address)
	if totalGreaterRoundDepositCount.Cmp(balanceDepositVotesRequiredCurrentRound) >= 0 {
		return true, nil
	}

	return false, nil
}

// No nil checks, call only if nil has been checked already
func getPacketType(packet *eth.ConsensusPacket) ConsensusPacketType {
	packetType := ConsensusPacketType(packet.ConsensusData[0])
	return packetType
}

func (cph *ConsensusHandler) shouldMoveToNextRoundPrecommit(parentHash common.Hash) (bool, error) {
	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if HasExceededTimeThreshold(blockRoundDetails.precommitInitTime, int64(ACK_BLOCK_TIMEOUT_MS*int(blockRoundDetails.Round))) == false {
		log.Trace("shouldMoveToNextRoundPrecommit time not met", "blockRoundDetails.precommitInitTime", blockRoundDetails.precommitInitTime)
		return false, nil
	}

	//Find validators in greater rounds
	valMap := make(map[common.Address]bool)
	valCommitMap := make(map[common.Address]bool)
	for _, pktList := range cph.outOfOrderPacketsMap {
		for _, pkt := range pktList {
			if pkt.Packet.ParentHash.IsEqualTo(parentHash) {
				round, validator, err := parsePacket(pkt.Packet)
				if err != nil {
					log.Trace("parsePacket", "err", err)
					continue
				}
				packetType := getPacketType(pkt.Packet)
				if round == blockStateDetails.currentRound && packetType == CONSENSUS_PACKET_TYPE_COMMIT_BLOCK { //todo: check commitHash
					return false, nil //todo: verify percentage
				}
				if round <= blockStateDetails.currentRound {
					continue
				}
				valMap[validator] = true
				log.Trace("shouldMoveToNextRound", "valInGreaterRound", validator)
			}
		}
	}

	totalGreaterRoundDepositCount := big.NewInt(0)
	currentRoundDepositSoFar := big.NewInt(0)
	currentRoundCommitsDepositSoFar := big.NewInt(0)
	for val, depositAmount := range blockStateDetails.filteredValidatorsDepositMap {
		_, okCommit := valCommitMap[val]
		if okCommit {
			currentRoundCommitsDepositSoFar = common.SafeAddBigInt(depositAmount, currentRoundCommitsDepositSoFar)
			log.Trace("currentRoundCommitsDepositSoFar", "val", val, "depositAmount", depositAmount, "currentRoundCommitsDepositSoFar", currentRoundCommitsDepositSoFar)
			continue
		}

		_, ok := valMap[val]
		if ok == false {
			_, ok1 := blockRoundDetails.validatorPrecommits[val]
			if ok1 == true {
				currentRoundDepositSoFar = common.SafeAddBigInt(depositAmount, currentRoundDepositSoFar)
				log.Trace("currentRoundDepositSoFar", "val", val, "depositAmount", depositAmount, "currentRoundDepositSoFar", currentRoundDepositSoFar)
			} else {
				log.Trace("currentRoundDepositSoFar val no", "val", val)
			}
		} else {
			totalGreaterRoundDepositCount = common.SafeAddBigInt(depositAmount, totalGreaterRoundDepositCount)
			log.Trace("totalGreaterRoundDepositCount", "val", val, "depositAmount", depositAmount, "totalGreaterRoundDepositCount", totalGreaterRoundDepositCount)
		}
	}

	currentRoundDepositSoFar = common.SafeAddBigInt(currentRoundCommitsDepositSoFar, currentRoundDepositSoFar)
	if currentRoundDepositSoFar.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 {
		return false, nil
	}

	balanceDepositVotesRequiredCurrentRound := common.SafeSubBigInt(blockStateDetails.blockMinWeightedProposalsRequired, currentRoundDepositSoFar)
	log.Trace("shouldMoveNextRound",
		"blockMinWeightedProposalsRequired", blockStateDetails.blockMinWeightedProposalsRequired,
		"balanceDepositVotesRequiredCurrentRound", balanceDepositVotesRequiredCurrentRound,
		"currentRoundDepositSoFar", currentRoundDepositSoFar,
		"totalGreaterRoundDepositCount", totalGreaterRoundDepositCount,
		"precommit count", len(blockRoundDetails.validatorPrecommits),
		"self precomitted", blockRoundDetails.selfPrecommited,
		"val", cph.account.Address)
	if totalGreaterRoundDepositCount.Cmp(balanceDepositVotesRequiredCurrentRound) >= 0 {
		return true, nil
	}

	return false, nil
}

func (cph *ConsensusHandler) handlePrecommitPacket(validator common.Address, packet *eth.ConsensusPacket, self bool) error {
	cph.innerPacketLock.Lock()
	defer cph.innerPacketLock.Unlock()

	blockStateDetails, ok := cph.blockStateDetailsMap[packet.ParentHash]
	if ok == false {
		return errors.New("unknown parentHash")
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[cph.account.Address]
	if ok == false {
		return errors.New("not a validator in this block")
	}

	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if blockRoundDetails.state != BLOCK_STATE_WAITING_FOR_PRECOMMITS {
		return OutOfOrderPackerErr
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[validator]
	if ok == false {
		log.Trace("handleProposeTransactionsPacket6")
		return errors.New("invalid validator")
	}

	if validator.IsEqualTo(cph.account.Address) == true && self == false {
		return errors.New("self packet from elsewhere")
	}

	_, ok = blockRoundDetails.validatorPrecommits[validator]
	if ok == true {
		//todo: check
	} else {

	}

	_, ok = blockRoundDetails.validatorProposalAcks[validator]
	if ok == false {
		log.Trace("did not receive proposal ack from validator")
	}

	precommitDetails := &PreCommitDetails{}

	err := rlp.DecodeBytes(packet.ConsensusData[1:], precommitDetails)
	if err != nil {
		log.Trace("handlePrecommitPacket err5", err)
		return err
	}

	if precommitDetails.Round != blockStateDetails.currentRound {
		log.Trace("handlePrecommitPacket OutOfOrderPackerErr", "round", precommitDetails.Round, "currentRound", blockStateDetails.currentRound)
		return OutOfOrderPackerErr
	}

	if precommitDetails.PrecommitHash.IsEqualTo(blockRoundDetails.precommitHash) == false {
		log.Trace("precommit error", "incoming", precommitDetails.PrecommitHash, "expected", blockRoundDetails.precommitHash, "me", cph.account.Address, "validator", validator)
		return errors.New("invalid Precommit Hash")
	}

	blockRoundDetails.validatorPrecommits[validator] = precommitDetails
	if self {
		blockRoundDetails.selfPrecommited = true
		blockRoundDetails.selfPrecommitPacket = packet
		log.Trace("self precomitted")
	}

	if blockRoundDetails.selfPrecommited {
		totalVotesDepositCount := big.NewInt(0)
		for val, _ := range blockRoundDetails.validatorPrecommits {
			totalVotesDepositCount = common.SafeAddBigInt(totalVotesDepositCount, blockStateDetails.filteredValidatorsDepositMap[val])
			log.Trace("Precommits", "validator", val, "deposit", blockStateDetails.filteredValidatorsDepositMap[val])
		}

		log.Debug("handlePrecommitPacket", "totalVotesDepositCount", totalVotesDepositCount, "blockMinWeightedProposalsRequired", blockStateDetails.blockMinWeightedProposalsRequired)
		if totalVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 {
			blockStateDetails.precommitTime = Elapsed(blockStateDetails.initTime)
			blockRoundDetails.state = BLOCK_STATE_WAITING_FOR_COMMITS
		}
	}

	pkt := eth.NewConsensusPacket(packet)
	blockRoundDetails.precommitPackets[validator] = &pkt

	blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
	cph.blockStateDetailsMap[packet.ParentHash] = blockStateDetails
	log.Trace("handlePrecommitPacket done", "ParentHash", packet.ParentHash)

	return nil
}

func (cph *ConsensusHandler) handleCommitPacket(validator common.Address, packet *eth.ConsensusPacket, self bool) error {
	cph.innerPacketLock.Lock()
	defer cph.innerPacketLock.Unlock()

	blockStateDetails, ok := cph.blockStateDetailsMap[packet.ParentHash]
	if ok == false {
		return errors.New("unknown parentHash")
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[cph.account.Address]
	if ok == false {
		return errors.New("not a validator in this block")
	}

	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]
	if blockRoundDetails.state != BLOCK_STATE_WAITING_FOR_COMMITS {
		return OutOfOrderPackerErr
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[validator]
	if ok == false {
		log.Trace("handleProposeTransactionsPacket6")
		return errors.New("invalid validator")
	}

	if validator.IsEqualTo(cph.account.Address) == true && self == false {
		return errors.New("self packet from elsewhere")
	}

	_, ok = blockRoundDetails.validatorCommits[validator]
	if ok == true {
		//todo: check
	} else {
	}

	_, ok = blockRoundDetails.validatorProposalAcks[validator]
	if ok == false {

	}

	commitDetails := &CommitDetails{}

	err := rlp.DecodeBytes(packet.ConsensusData[1:], commitDetails)
	if err != nil {
		log.Trace("handlePrecommitPacket err5", "err", err)
		return err
	}

	if commitDetails.Round != blockStateDetails.currentRound {
		return OutOfOrderPackerErr
	}

	var commitHash common.Hash
	commitHash.CopyFrom(crypto.Keccak256Hash(blockRoundDetails.precommitHash.Bytes()))
	if commitDetails.CommitHash.IsEqualTo(commitHash) == false { //PrecommitHash and commitHash should be the same
		return errors.New("invalid commit Hash")
	}

	blockRoundDetails.validatorCommits[validator] = commitDetails
	if self {
		blockRoundDetails.selfCommited = true
		blockRoundDetails.selfCommitPacket = packet
	}

	if blockRoundDetails.selfCommited {
		totalVotesDepositCount := big.NewInt(0)
		for val, _ := range blockRoundDetails.validatorCommits {
			totalVotesDepositCount = common.SafeAddBigInt(totalVotesDepositCount, blockStateDetails.filteredValidatorsDepositMap[val])
			log.Trace("Commits", "validator", val, "deposit", blockStateDetails.filteredValidatorsDepositMap[val])
		}

		log.Debug("handleCommitPacket", "totalVotesDepositCount", totalVotesDepositCount, "blockMinWeightedProposalsRequired", blockStateDetails.blockMinWeightedProposalsRequired)

		if totalVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 {
			if blockRoundDetails.blockVoteType == VOTE_TYPE_NIL {
				cph.nilVoteBlocks = cph.nilVoteBlocks + 1
			} else if blockRoundDetails.blockVoteType == VOTE_TYPE_OK {
				cph.okVoteBlocks = cph.okVoteBlocks + 1
				txnCountInBlock := uint64(len(blockRoundDetails.proposalTxns))
				cph.totalTransactions = cph.totalTransactions + txnCountInBlock
				if txnCountInBlock > cph.maxTransactionsInBlock {
					cph.maxTransactionsInBlock = txnCountInBlock
					cph.maxTransactionsBlockTime = Elapsed(blockStateDetails.initTime)
				}
			}
			blockStateDetails.commitTime = Elapsed(blockStateDetails.initTime)

			//stats
			cph.timeStatMap[GetTimeStateBucket(PROPOSAL_KEY_PREFIX, blockStateDetails.proposalTime)]++
			cph.timeStatMap[GetTimeStateBucket(ACK_PROPOSAL_KEY_PREFIX, blockStateDetails.ackProposalTime-blockStateDetails.proposalTime)]++
			cph.timeStatMap[GetTimeStateBucket(PRECOMMIT_KEY_PREFIX, blockStateDetails.precommitTime-blockStateDetails.ackProposalTime)]++
			cph.timeStatMap[GetTimeStateBucket(COMMIT_KEY_PREFIX, blockStateDetails.commitTime-blockStateDetails.precommitTime)]++

			log.Trace("BlockStats", "maxTxnsInBlock", cph.maxTransactionsInBlock, "totalTxns", cph.totalTransactions, "okBlocks", cph.okVoteBlocks, "nilBlocks", cph.nilVoteBlocks)
			for statKey, statVal := range cph.timeStatMap {
				if statVal > 0 {
					log.Trace("TimeStatsBlockCount", "stat", statKey, "blocks", statVal)
				}
			}

			blockRoundDetails.state = BLOCK_STATE_RECEIVED_COMMITS
		}
	}

	pkt := eth.NewConsensusPacket(packet)
	blockRoundDetails.commitPackets[validator] = &pkt

	blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
	cph.blockStateDetailsMap[packet.ParentHash] = blockStateDetails

	return nil
}

func HasExceededTimeThreshold(startTime time.Time, thresholdMs int64) bool {
	end := time.Now().UnixNano() / int64(time.Millisecond)
	start := startTime.UnixNano() / int64(time.Millisecond)
	diff := end - start
	if diff >= thresholdMs {
		return true
	} else {
		return false
	}
}

func Elapsed(startTime time.Time) int64 {
	end := time.Now().UnixNano() / int64(time.Millisecond)
	start := startTime.UnixNano() / int64(time.Millisecond)
	diff := end - start
	return diff
}

func GetProposalTime(blockNumber uint64) uint64 {
	if blockNumber == 1 || blockNumber%BLOCK_PERIOD_TIME_CHANGE == 0 {
		blockTime := uint64(time.Now().UTC().Unix())
		if blockTime%60 != 0 {
			blockTime = blockTime - (blockTime % 60)
		}

		return blockTime
	} else {
		return 0
	}
}

func ValidateBlockProposalTimeConsensus(blockNumber uint64, proposedTime uint64) bool {
	if blockNumber == 1 || blockNumber%BLOCK_PERIOD_TIME_CHANGE == 0 {
		if proposedTime == 0 {
			return false
		}

		tm := time.Unix(int64(proposedTime), 0)
		if tm.Second() != 0 || tm.Nanosecond() != 0 { //No granularity at anything other than minute level allowed, to reduce ability to manipulate blockHash
			return false
		}
		currTimeVal := time.Now().UTC().Unix() //Note that packet may have arrived late. So, these comparisions are approximate.
		if currTimeVal%60 != 0 {
			currTimeVal = currTimeVal - (currTimeVal % 60)
		}
		currTime := time.Unix(currTimeVal, 0)

		if currTime.Before(tm) {
			difference := tm.Sub(currTime)
			if difference.Minutes() > ALLOWED_TIME_SKEW_MINUTES {
				return false
			}
		} else if currTime.After(tm) {
			difference := currTime.Sub(tm)
			if difference.Minutes() > ALLOWED_TIME_SKEW_MINUTES {
				return false
			}
		}
	} else {
		if proposedTime != 0 {
			return false
		}
	}

	return true
}

func (cph *ConsensusHandler) proposeBlock(parentHash common.Hash, txns []common.Hash, blockNumber uint64) error {
	var packet *eth.ConsensusPacket
	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if blockRoundDetails.selfProposed == true {
		return cph.broadCast(blockRoundDetails.selfProposalPacket)
	}

	proposalDetails := &ProposalDetails{}

	proposalDetails.Round = blockStateDetails.currentRound
	if blockStateDetails.currentRound < MAX_ROUND { //No transactions after this round, to reduce chance of FLP
		proposalDetails.Txns = make([]common.Hash, len(txns))
		for i := 0; i < len(proposalDetails.Txns); i++ {
			proposalDetails.Txns[i].CopyFrom(txns[i])
		}
	} else {
		proposalDetails.Txns = make([]common.Hash, 0)
	}
	proposalDetails.BlockTime = GetProposalTime(blockNumber)

	log.Trace("ProposeBlock with txns", "count", len(proposalDetails.Txns))

	data, err := rlp.EncodeToBytes(proposalDetails)

	if err != nil {
		return err
	}

	dataToSend := append([]byte{byte(CONSENSUS_PACKET_TYPE_PROPOSE_BLOCK)}, data...)

	fullSignNeeded := shouldSignFull(blockNumber)
	packet, err = cph.createConsensusPacket(parentHash, dataToSend, fullSignNeeded)
	if err != nil {
		return err
	}

	proposalDetails1 := ProposalDetails{}

	err = rlp.DecodeBytes(packet.ConsensusData[1:], &proposalDetails1)
	if err != nil {
		return err
	}

	err = cph.handleProposeBlockPacket(cph.account.Address, packet, true)

	if err != nil {
		return err
	}

	return cph.broadCast(packet)
}

func (cph *ConsensusHandler) ackBlockProposalTimeout(parentHash common.Hash) error {
	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if blockRoundDetails.selfAckd == true {
		log.Trace("ackBlockProposalTimeout selfAckd", "parentHash", parentHash)
	} else {
		if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL {
		} else {
			if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS {
			} else {
				return errors.New("unexpected state")
			}
		}

		proposalAckDetails := &ProposalAckDetails{
			ProposalAckVoteType: VOTE_TYPE_NIL,
			Round:               blockStateDetails.currentRound,
		}

		proposalAckDetails.ProposalHash.CopyFrom(getNilVoteProposalHash(parentHash, blockStateDetails.currentRound))

		data, err := rlp.EncodeToBytes(&proposalAckDetails)

		if err != nil {
			return err
		}

		dataToSend := append([]byte{byte(CONSENSUS_PACKET_TYPE_ACK_BLOCK_PROPOSAL)}, data...)
		packet, err := cph.createConsensusPacket(parentHash, dataToSend, false)
		if err != nil {
			return err
		}

		pkt := eth.NewConsensusPacket(packet)
		blockRoundDetails.proposalAckPackets[cph.account.Address] = &pkt
		blockRoundDetails.validatorProposalAcks[cph.account.Address] = proposalAckDetails
		blockRoundDetails.selfAckd = true
		blockRoundDetails.selfAckPacket = packet
		blockRoundDetails.selfAckProposalVoteType = proposalAckDetails.ProposalAckVoteType
		blockRoundDetails.blockVoteType = VOTE_TYPE_NIL
		log.Trace("blockVoteType a3", "parentHash", parentHash)
	}

	okVotes := 0
	nilVotes := 0
	mismatchedVotes := 0
	okVotesDepositCount := big.NewInt(0)
	nilVotesDepositCount := big.NewInt(0)
	totalVotesDepositCount := big.NewInt(0)

	for val, ack := range blockRoundDetails.validatorProposalAcks {
		log.Trace("validatorProposalAcks", "validator", val, "deposit", blockStateDetails.filteredValidatorsDepositMap[val], "voteType", ack.ProposalAckVoteType)
		if ack.ProposalAckVoteType == VOTE_TYPE_OK {
			if ack.ProposalHash.IsEqualTo(blockRoundDetails.proposalHash) {
				okVotes = okVotes + 1
				okVotesDepositCount = common.SafeAddBigInt(okVotesDepositCount, blockStateDetails.filteredValidatorsDepositMap[val])
			} else {
				mismatchedVotes = mismatchedVotes + 1
			}
			totalVotesDepositCount = common.SafeAddBigInt(totalVotesDepositCount, blockStateDetails.filteredValidatorsDepositMap[val])
		} else if ack.ProposalAckVoteType == VOTE_TYPE_NIL {
			nilVotes = nilVotes + 1
			nilVotesDepositCount = common.SafeAddBigInt(nilVotesDepositCount, blockStateDetails.filteredValidatorsDepositMap[val])
			totalVotesDepositCount = common.SafeAddBigInt(totalVotesDepositCount, blockStateDetails.filteredValidatorsDepositMap[val])
		} else {
			return errors.New("unexpected")
		}
	}

	log.Debug("ackBlockProposalTimeout", "totalBlockDepositValue", blockStateDetails.totalBlockDepositValue, "okVotesDepositCount",
		okVotesDepositCount, "okVotesDepositCount", okVotesDepositCount, "nilVotesDepositCount", nilVotesDepositCount,
		"totalVotesDepositCount", totalVotesDepositCount, "okVotes", okVotes, "nilVotes", nilVotes, "mismatchedVotes", mismatchedVotes,
		"blockMinWeightedProposalsRequired", blockStateDetails.blockMinWeightedProposalsRequired)

	if okVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 {
		//do nothing
	} else if nilVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 { //handle timeout differently?
		blockRoundDetails.state = BLOCK_STATE_WAITING_FOR_PRECOMMITS
		blockRoundDetails.precommitInitTime = time.Now()
		blockRoundDetails.blockVoteType = VOTE_TYPE_NIL
		blockRoundDetails.precommitHash.CopyFrom(getNilVotePreCommitHash(parentHash, blockStateDetails.currentRound))
	} else {
		if HasExceededTimeThreshold(blockRoundDetails.initTime, int64(ACK_BLOCK_TIMEOUT_MS*int(blockRoundDetails.Round))) {
			if totalVotesDepositCount.Cmp(blockStateDetails.totalBlockDepositValue) >= 0 ||
				totalVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 {
				blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
				cph.blockStateDetailsMap[parentHash] = blockStateDetails
				err := cph.initializeNewBlockRound(NEW_ROUND_REASON_WAIT_ACK_BLOCK_PROPOSAL_TIMEOUT)
				if err != nil {
					return err
				}
				return nil
			} else {
				ok, err := cph.shouldMoveToNextRoundProposalAcks(parentHash)
				if err != nil {
					return err
				}
				if ok == true {
					blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
					cph.blockStateDetailsMap[parentHash] = blockStateDetails
					err := cph.initializeNewBlockRound(NEW_ROUND_REASON_WAIT_ACK_BLOCK_PROPOSAL_HIGHER_ROUND)
					if err != nil {
						return err
					}
					return nil
				}
			}
		}
	}

	blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
	cph.blockStateDetailsMap[parentHash] = blockStateDetails

	err := cph.broadCast(blockRoundDetails.selfAckPacket)
	if err != nil {
		return err
	}

	return cph.broadcastPreviousRoundPackets(parentHash)
}

func (cph *ConsensusHandler) broadcastPreviousRoundPackets(parentHash common.Hash) error {
	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if blockRoundDetails.Round > 1 {
		for i := byte(1); i < blockRoundDetails.Round; i = i + 1 {
			prevBlockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]
			if prevBlockRoundDetails.selfAckd {
				log.Trace("Broadcasting selfAckPacket", "parentHash", parentHash, "round", i)
				err := cph.broadCast(prevBlockRoundDetails.selfAckPacket)
				if err != nil {
					log.Error("broadcastPreviousRoundPackets", "err", err)
					return err
				}
			}
		}
	}

	return nil
}

func (cph *ConsensusHandler) ackBlockProposal(parentHash common.Hash) error {
	log.Trace("ackBlockProposal")
	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if blockRoundDetails.selfAckd == true {
		shouldPropose, err := cph.isBlockProposer(parentHash, &blockStateDetails.filteredValidatorsDepositMap, blockStateDetails.currentRound, blockStateDetails)
		if err != nil {
			return err
		}
		if shouldPropose {
			cph.broadCast(blockRoundDetails.selfProposalPacket)
		}
	} else {
		if blockRoundDetails.state != BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS {
			return errors.New("unexpected state")
		}

		if blockStateDetails.currentRound >= MAX_ROUND && len(blockRoundDetails.blockProposalDetails.Txns) > 0 {
			return errors.New("unexpected transaction count")
		} else {
			//Find if any new transactions we don't know yet
			unknownTxns := make([]common.Hash, 0)
			for i := 0; i < len(blockRoundDetails.blockProposalDetails.Txns); i++ {
				_, txnExists := blockRoundDetails.selfKnownTransactions[blockRoundDetails.blockProposalDetails.Txns[i]]
				if txnExists == false {
					log.Trace("===============ackBlockProposal unknown txns", "hash", blockRoundDetails.blockProposalDetails.Txns[i])
					unknownTxns = append(unknownTxns, blockRoundDetails.blockProposalDetails.Txns[i])
				}
			}
			if len(unknownTxns) > 0 {
				err := cph.p2pHandler.RequestTransactions(unknownTxns)
				if err != nil {
					return errors.New("unknown transactions")
				}
			}
		}

		var voteType VoteType
		if blockStateDetails.currentRound >= MAX_ROUND {
			voteType = VOTE_TYPE_NIL
		} else {
			voteType = VOTE_TYPE_OK
		}

		proposalAckDetails := &ProposalAckDetails{
			ProposalAckVoteType: voteType,
			Round:               blockStateDetails.currentRound,
		}

		if blockStateDetails.currentRound >= MAX_ROUND {
			proposalAckDetails.ProposalHash.CopyFrom(getNilVoteProposalHash(parentHash, blockStateDetails.currentRound))
		} else {
			proposalAckDetails.ProposalHash.CopyFrom(blockRoundDetails.proposalHash)
		}

		data, err := rlp.EncodeToBytes(proposalAckDetails)

		if err != nil {
			return err
		}

		dataToSend := append([]byte{byte(CONSENSUS_PACKET_TYPE_ACK_BLOCK_PROPOSAL)}, data...)
		packet, err := cph.createConsensusPacket(parentHash, dataToSend, false)
		if err != nil {
			return err
		}

		pkt := eth.NewConsensusPacket(packet)
		blockRoundDetails.proposalAckPackets[cph.account.Address] = &pkt
		blockRoundDetails.validatorProposalAcks[cph.account.Address] = proposalAckDetails
		blockRoundDetails.selfAckd = true
		blockRoundDetails.selfAckPacket = packet
		blockRoundDetails.selfAckProposalVoteType = proposalAckDetails.ProposalAckVoteType
	}

	okVotesCount := 0
	nilVotesCount := 0
	mismatchedVotesCount := 0
	okVotesDepositCount := big.NewInt(0)
	nilVotesDepositCount := big.NewInt(0)
	totalVotesDepositCount := big.NewInt(0)

	for val, ack := range blockRoundDetails.validatorProposalAcks {
		log.Trace("validatorProposalAcks", "validator", val, "deposit", blockStateDetails.filteredValidatorsDepositMap[val], "voteType", ack.ProposalAckVoteType)
		if ack.ProposalAckVoteType == VOTE_TYPE_OK {
			if ack.ProposalHash.IsEqualTo(blockRoundDetails.proposalHash) {
				okVotesCount = okVotesCount + 1
				okVotesDepositCount = common.SafeAddBigInt(okVotesDepositCount, blockStateDetails.filteredValidatorsDepositMap[val])
			} else {
				mismatchedVotesCount = mismatchedVotesCount + 1
			}
			totalVotesDepositCount = common.SafeAddBigInt(totalVotesDepositCount, blockStateDetails.filteredValidatorsDepositMap[val])
		} else if ack.ProposalAckVoteType == VOTE_TYPE_NIL {
			nilVotesCount = nilVotesCount + 1
			nilVotesDepositCount = common.SafeAddBigInt(nilVotesDepositCount, blockStateDetails.filteredValidatorsDepositMap[val])
			totalVotesDepositCount = common.SafeAddBigInt(totalVotesDepositCount, blockStateDetails.filteredValidatorsDepositMap[val])
		} else {
			log.Trace("unexpected")
			return errors.New("unexpected")
		}
	}

	log.Debug("ackBlockProposal", "totalBlockDepositValue", blockStateDetails.totalBlockDepositValue, "okVotesDepositCount",
		okVotesDepositCount, "okVotesDepositCount", okVotesDepositCount, "nilVotesDepositCount", nilVotesDepositCount, "totalVotesDepositCount", totalVotesDepositCount,
		"okVotes", okVotesCount, "nilVotes", nilVotesCount, "blockMinWeightedProposalsRequired", blockStateDetails.blockMinWeightedProposalsRequired)

	if okVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 && blockRoundDetails.selfAckProposalVoteType == VOTE_TYPE_OK { //For ok votes, vote type should match
		blockRoundDetails.state = BLOCK_STATE_WAITING_FOR_PRECOMMITS
		blockRoundDetails.precommitInitTime = time.Now()
		blockRoundDetails.precommitHash.CopyFrom(getOkVotePreCommitHash(parentHash, blockRoundDetails.proposalHash, blockStateDetails.currentRound))
		blockRoundDetails.blockVoteType = VOTE_TYPE_OK
		log.Trace("blockVoteType a1", "parentHash", parentHash)
		blockStateDetails.ackProposalTime = Elapsed(blockStateDetails.initTime)
	} else if nilVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 { //handle timeout differently? for nil votes, it is ok to accept NIL vote even if self vote is OK
		blockRoundDetails.state = BLOCK_STATE_WAITING_FOR_PRECOMMITS
		blockRoundDetails.precommitInitTime = time.Now()
		blockRoundDetails.precommitHash.CopyFrom(getNilVotePreCommitHash(parentHash, blockStateDetails.currentRound))
		log.Trace("blockVoteType a2", "parentHash", parentHash)
		blockRoundDetails.blockVoteType = VOTE_TYPE_NIL
	} else {
		if totalVotesDepositCount.Cmp(blockStateDetails.totalBlockDepositValue) >= 0 ||
			totalVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 && HasExceededTimeThreshold(blockRoundDetails.initTime, int64(ACK_BLOCK_TIMEOUT_MS*int(blockRoundDetails.Round))) {
			blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
			cph.blockStateDetailsMap[parentHash] = blockStateDetails
			err := cph.initializeNewBlockRound(NEW_ROUND_REASON_WAIT_ACK_BLOCK_PROPOSAL_TIMEOUT)
			if err != nil {
				return err
			}
			log.Trace("blockVoteType a3", "parentHash", parentHash)
			return nil
		} else {
			ok, err := cph.shouldMoveToNextRoundProposalAcks(parentHash)
			if err != nil {
				return err
			}
			if ok == true {
				blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
				cph.blockStateDetailsMap[parentHash] = blockStateDetails
				err := cph.initializeNewBlockRound(NEW_ROUND_REASON_WAIT_ACK_BLOCK_PROPOSAL_HIGHER_ROUND)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
	cph.blockStateDetailsMap[parentHash] = blockStateDetails

	err := cph.broadCast(blockRoundDetails.selfAckPacket)
	if err != nil {
		return err
	}

	return cph.broadcastPreviousRoundPackets(parentHash)
}

func (cph *ConsensusHandler) precommitBlock(parentHash common.Hash) error {
	blockStateDetails, ok := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if ok == false || blockRoundDetails.selfAckPacket == nil {
		return errors.New("invalid state 1")
	}

	if blockRoundDetails.state != BLOCK_STATE_WAITING_FOR_PRECOMMITS {
		return errors.New("invalid state 2")
	}

	if blockRoundDetails.selfPrecommited == true {
		log.Trace("precommitBlock broadcast")
		cph.broadCast(blockRoundDetails.selfAckPacket)
		cph.broadCast(blockRoundDetails.selfPrecommitPacket)
		return cph.handlePrecommitPacket(cph.account.Address, blockRoundDetails.selfPrecommitPacket, true)
	}

	precommit := &PreCommitDetails{
		PrecommitHash: blockRoundDetails.precommitHash,
		Round:         blockRoundDetails.Round,
	}

	data, err := rlp.EncodeToBytes(precommit)

	if err != nil {
		return err
	}

	dataToSend := append([]byte{byte(CONSENSUS_PACKET_TYPE_PRECOMMIT_BLOCK)}, data...)
	packet, err := cph.createConsensusPacket(parentHash, dataToSend, false)
	if err != nil {
		return err
	}

	err = cph.handlePrecommitPacket(cph.account.Address, packet, true)

	if err != nil {
		log.Trace("precommitBlock handlePrecommitPacket error", err)
		return err
	}

	cph.broadCast(blockRoundDetails.selfAckPacket)
	return cph.broadCast(packet)
}

func (cph *ConsensusHandler) commitBlock(parentHash common.Hash) error {
	blockStateDetails, ok := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if ok == false || blockRoundDetails.selfAckPacket == nil || blockRoundDetails.selfPrecommitPacket == nil {
		return errors.New("invalid state 1")
	}

	if blockRoundDetails.state != BLOCK_STATE_WAITING_FOR_COMMITS {
		return errors.New("invalid state 2")
	}

	if blockRoundDetails.selfCommited == true {
		cph.broadCast(blockRoundDetails.selfAckPacket)
		cph.broadCast(blockRoundDetails.selfPrecommitPacket)
		cph.broadCast(blockRoundDetails.selfCommitPacket)
		return cph.handleCommitPacket(cph.account.Address, blockRoundDetails.selfCommitPacket, true)
	}

	commitDetails := &CommitDetails{
		Round: blockRoundDetails.Round,
	}
	commitDetails.CommitHash.CopyFrom(getCommitHash(blockRoundDetails.precommitHash))

	data, err := rlp.EncodeToBytes(commitDetails)

	if err != nil {
		return err
	}

	dataToSend := append([]byte{byte(CONSENSUS_PACKET_TYPE_COMMIT_BLOCK)}, data...)
	packet, err := cph.createConsensusPacket(parentHash, dataToSend, false)
	if err != nil {
		return err
	}

	err = cph.handleCommitPacket(cph.account.Address, packet, true)

	if err != nil {
		log.Trace("commitBlock handleCommitPacket error", "err", err)
		return err
	}

	return cph.broadCast(packet)
}

func (cph *ConsensusHandler) DoesPreviousHashMatch(parentHash common.Hash) (bool, error) {
	skipCheck := os.Getenv("SKIP_CONSENSUS_STARTUP_HASH_CHECK")
	if SKIP_HASH_CHECK || (len(skipCheck) > 0 && skipCheck == "1") {
		log.Warn("SKIP_CONSENSUS_STARTUP_HASH_CHECK is set, skipping hash check")
		return false, nil
	}

	datadir := node.DefaultDataDir()
	hashFilePath := filepath.Join(datadir, "previoushash.txt")
	log.Trace("DoesPreviousHashMatch", "path", hashFilePath, "parentHash", parentHash)

	if _, err := os.Stat(hashFilePath); errors.Is(err, os.ErrNotExist) {
		log.Trace("DoesPreviousHashMatch previous hash not found")
		return false, nil
	}

	b, err := ioutil.ReadFile(hashFilePath)
	if err != nil {
		log.Warn("DoesPreviousHashMatch", "err", err, "hashFilePath", hashFilePath)
		return false, err
	}
	hash := common.HexToHash(string(b))

	if hash.IsEqualTo(parentHash) {
		return true, nil
	}

	log.Trace("Previous doesn't match current parentHash, is ok to proceed with consensus", "previous", hash.Hex(), "current parentHash", parentHash.Hex())
	return false, nil
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

func (cph *ConsensusHandler) SaveHash(parentHash common.Hash) error {
	datadir := node.DefaultDataDir()
	hashFilePath := filepath.Join(datadir, "previoushash.txt")

	if err := ensureDir(datadir); err != nil {
		return err
	}

	f, err := os.Create(hashFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(parentHash.Hex())
	if err != nil {
		return err
	}

	return nil
}

func (cph *ConsensusHandler) HandleConsensus(parentHash common.Hash, txns []common.Hash, blockNumber uint64) error {
	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	if cph.initialized == false {

		matched, err := cph.DoesPreviousHashMatch(parentHash)
		if err != nil {
			log.Warn("DoesPreviousHashMatch on parent hash failed")
			return err
		}

		if matched {
			log.Warn("Previous block hash before restart matches current parentHash. Will wait for one block to get mined before starting.", "parentHash", parentHash)
			return errors.New("Waiting for previous block to mine")
		}

		cph.initTime = time.Now()
		cph.initialized = true
		cph.packetHashLastSentMap = make(map[common.Hash]time.Time)
		cph.packetStats = PacketStats{}

		log.Info("Starting up...")
		return errors.New("starting up")
	}

	if cph.lastBlockNumber == blockNumber {
		if Elapsed(cph.lastBlockNumberChangeTime) >= STALE_BLOCK_WARN_TIME {
			log.Warn("Stale Block. Please check your connection.", "blockNumber", blockNumber, "lastBlockChangeTime", cph.lastBlockNumberChangeTime)
		}
	} else {
		cph.lastBlockNumber = blockNumber
		cph.lastBlockNumberChangeTime = time.Now()
	}

	if HasExceededTimeThreshold(cph.initTime, STARTUP_DELAY_MS) == false {
		log.Info("Waiting to startup...", "elapsed ms", Elapsed(cph.initTime), "pending txn count", len(txns), "STARTUP_DELAY_MS", STARTUP_DELAY_MS)
		return errors.New("starting up")
	}

	err := cph.initializeBlockStateIfRequired(parentHash, blockNumber)
	if err != nil {
		return err
	}

	cph.cleanupBlockState()

	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	_, ok := blockStateDetails.filteredValidatorsDepositMap[cph.account.Address]
	if ok == false {
		return errors.New("not a validator in this block")
	}

	shouldPropose, err := cph.isBlockProposer(parentHash, &blockStateDetails.filteredValidatorsDepositMap, blockStateDetails.currentRound, blockStateDetails)
	if err != nil {
		return err
	}

	cph.processOutOfOrderPackets(parentHash)

	err = errors.New("not ready yet")
	log.Info("HandleConsensus", "parentHash", parentHash, "blockNumber", blockNumber, "currentRound", blockStateDetails.currentRound, "state", blockRoundDetails.state, "blockVoteType", blockRoundDetails.blockVoteType,
		"selfAckProposalVoteType", blockRoundDetails.selfAckProposalVoteType,
		"shouldPropose", shouldPropose, "currTxns", len(txns), "okVoteBlocks", cph.okVoteBlocks, "nilVoteBlocks", cph.nilVoteBlocks,
		"totalTransactions", cph.totalTransactions, "maxTransactionsInBlock", cph.maxTransactionsInBlock, "maxTransactionsBlockTime", cph.maxTransactionsBlockTime,
		"pending txns", len(txns), "TotalIncomingPackets", cph.packetStats.TotalIncomingPacketCount, "newRoundReason", blockRoundDetails.newRoundReason)

	if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL {
		for _, txn := range txns {
			blockRoundDetails.selfKnownTransactions[txn] = true
		}
		blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
		cph.blockStateDetailsMap[parentHash] = blockStateDetails

		if shouldPropose {
			cph.proposeBlock(parentHash, txns, blockNumber)
		} else {
			var timeoutMs int64
			if shouldSignFull(blockNumber) {
				timeoutMs = FULL_BLOCK_TIMEOUT_MS
			} else {
				timeoutMs = BLOCK_TIMEOUT_MS
			}
			if HasExceededTimeThreshold(blockRoundDetails.initTime, timeoutMs*int64(blockRoundDetails.Round)) {
				cph.ackBlockProposalTimeout(parentHash)
			} else {
				cph.requestConsensusData(blockStateDetails)
			}
		}
	} else if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS {
		for _, txn := range txns {
			blockRoundDetails.selfKnownTransactions[txn] = true
		}
		blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
		cph.blockStateDetailsMap[parentHash] = blockStateDetails

		cph.requestConsensusData(blockStateDetails)
		if blockRoundDetails.selfAckProposalVoteType == VOTE_TYPE_NIL {
			cph.ackBlockProposalTimeout(parentHash)
		} else {
			cph.ackBlockProposal(parentHash)
		}
	} else if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PRECOMMITS {
		shouldMove, err := cph.shouldMoveToNextRoundPrecommit(parentHash)
		if err == nil && shouldMove {
			return cph.initializeNewBlockRound(NEW_ROUND_REASON_WAIT_PRECOMMIT_TIMEOUT)
		} else {
			cph.requestConsensusData(blockStateDetails)
			cph.precommitBlock(parentHash)
		}
	} else if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_COMMITS {
		cph.requestConsensusData(blockStateDetails)
		cph.commitBlock(parentHash)
	} else if blockRoundDetails.state == BLOCK_STATE_RECEIVED_COMMITS {
		cph.broadCast(blockRoundDetails.selfCommitPacket)
		err = nil
	}

	return err
}

func (cph *ConsensusHandler) createConsensusPacket(parentHash common.Hash, data []byte, fullSign bool) (*eth.ConsensusPacket, error) {
	if cph.signFn == nil {
		return nil, errors.New("signFn is not set")
	}
	dataToSign := append(parentHash.Bytes(), data...)
	var signature []byte
	var err error
	if fullSign {
		log.Debug("createConsensusPacket", "parentHash", parentHash, "fullSign", fullSign)
		signature, err = cph.signFnWithContext(cph.account, accounts.MimetypeProofOfStake, dataToSign, FULL_SIGN_CONTEXT)
	} else {
		log.Trace("createConsensusPacket", "parentHash", parentHash, "fullSign", fullSign)
		signature, err = cph.signFn(cph.account, accounts.MimetypeProofOfStake, dataToSign)
	}
	if err != nil {
		log.Trace("createConsensusPacket signAndSend failed", "err", err)
		return nil, err
	}

	packet := &eth.ConsensusPacket{
		ParentHash: parentHash,
	}

	packet.ConsensusData = make([]byte, len(data))
	copy(packet.ConsensusData, data)

	packet.Signature = make([]byte, len(signature))
	copy(packet.Signature, signature)

	return packet, nil
}

func (cph *ConsensusHandler) cleanupBroadcast() {
	for k, v := range cph.packetHashLastSentMap {
		elapsed := Elapsed(v)
		if elapsed >= BROADCAST_CLEANUP_DELAY {
			delete(cph.packetHashLastSentMap, k)
		}
	}
}

func (cph *ConsensusHandler) broadCast(packet *eth.ConsensusPacket) error {
	cph.p2pLock.Lock()
	defer cph.p2pLock.Unlock()
	if packet == nil {
		debug.PrintStack()
		return errors.New("packet is nil")
	}

	dataToHash := append(packet.ParentHash.Bytes(), packet.ConsensusData...)
	digestHash := crypto.Keccak256(dataToHash)
	var hash common.Hash
	hash.SetBytes(digestHash)

	packetType := ConsensusPacketType(packet.ConsensusData[0])
	lastSent, ok := cph.packetHashLastSentMap[hash]
	if ok == false {
		cph.packetHashLastSentMap[hash] = time.Now()
		log.Trace("Broadcasting packet", "hash", hash, "packetType", packetType)
	} else {
		elapsed := Elapsed(lastSent)
		if elapsed > BROADCAST_RESEND_DELAY {
			cph.packetHashLastSentMap[hash] = time.Now()
			log.Trace("Rebroadcasting packet", "hash", hash, "packetType", packetType)
		} else {
			log.Trace("Skipping broadcasting packet", "hash", hash, "packetType", packetType)
			return nil
		}
	}

	cph.cleanupBroadcast()
	go cph.p2pHandler.BroadcastConsensusData(packet)
	return nil
}

func (cph *ConsensusHandler) getRequestConsensusDataPacket(blockStateDetails *BlockStateDetails) (*RequestConsensusPacketDetails, error) {
	requestPacketDetails := RequestConsensusPacketDetails{
		RequestProposal:       false,
		ValidatorProposalAcks: make([]common.Address, 0),
		ValidatorPrecommits:   make([]common.Address, 0),
	}
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL {
		requestPacketDetails.RequestProposal = true
	} else if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS {
		for val, _ := range blockStateDetails.filteredValidatorsDepositMap {
			_, ok := blockRoundDetails.proposalAckPackets[val]
			if ok == false {
				requestPacketDetails.ValidatorProposalAcks = append(requestPacketDetails.ValidatorProposalAcks, val)
			}
		}
	} else if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PRECOMMITS {
		for val, _ := range blockStateDetails.filteredValidatorsDepositMap {
			_, ok := blockRoundDetails.proposalAckPackets[val]
			if ok == false {
				requestPacketDetails.ValidatorProposalAcks = append(requestPacketDetails.ValidatorProposalAcks, val)
			}
		}

		for val, _ := range blockStateDetails.filteredValidatorsDepositMap {
			_, ok := blockRoundDetails.precommitPackets[val]
			if ok == false {
				requestPacketDetails.ValidatorPrecommits = append(requestPacketDetails.ValidatorPrecommits, val)
			}
		}
	} else if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_COMMITS {
		for val, _ := range blockStateDetails.filteredValidatorsDepositMap {
			_, ok := blockRoundDetails.proposalAckPackets[val]
			if ok == false {
				requestPacketDetails.ValidatorProposalAcks = append(requestPacketDetails.ValidatorProposalAcks, val)
			}
		}

		for val, _ := range blockStateDetails.filteredValidatorsDepositMap {
			_, ok := blockRoundDetails.precommitPackets[val]
			if ok == false {
				requestPacketDetails.ValidatorPrecommits = append(requestPacketDetails.ValidatorPrecommits, val)
			}
		}

		for val, _ := range blockStateDetails.filteredValidatorsDepositMap {
			_, ok := blockRoundDetails.commitPackets[val]
			if ok == false {
				requestPacketDetails.ValidatorCommits = append(requestPacketDetails.ValidatorCommits, val)
			}
		}
	} else {
		return nil, errors.New("unknown state")
	}

	return &requestPacketDetails, nil
}

func (cph *ConsensusHandler) requestConsensusData(blockStateDetails *BlockStateDetails) error {
	cph.p2pLock.Lock()
	defer cph.p2pLock.Unlock()

	dataToHash := append(blockStateDetails.parentHash.Bytes())
	digestHash := crypto.Keccak256(dataToHash)
	var hash common.Hash
	hash.SetBytes(digestHash)

	lastSent, ok := cph.packetHashLastSentMap[hash]
	if ok == false {
		cph.packetHashLastSentMap[hash] = time.Now()
		log.Trace("requestConsensusData packet", "hash", hash)
	} else {
		elapsed := Elapsed(lastSent)
		if elapsed > BROADCAST_RESEND_DELAY*3 {
			cph.packetHashLastSentMap[hash] = time.Now()
			log.Trace("requestConsensusData packet", "hash", hash)
		} else {
			log.Trace("Skipping requestConsensusData packet", "hash", hash)
			return nil
		}
	}

	elapsed := Elapsed(blockStateDetails.initTime)
	if elapsed < BLOCK_TIMEOUT_MS {
		return nil
	}

	elapsed = Elapsed(cph.lastRequestConsensusDataTime)
	if elapsed < CONSENSUS_DATA_REQUEST_RESEND_DELAY {
		return nil
	}
	cph.lastRequestConsensusDataTime = time.Now()

	log.Trace("requestConsensusData 1")
	requestPacketDetails, err := cph.getRequestConsensusDataPacket(blockStateDetails)
	if err != nil {
		return err
	}

	data, err := rlp.EncodeToBytes(requestPacketDetails)
	if err != nil {
		return err
	}

	packet := eth.RequestConsensusDataPacket{}

	packet.RequestData = make([]byte, len(data))
	copy(packet.RequestData, data)
	packet.ParentHash = blockStateDetails.parentHash

	go cph.p2pHandler.RequestConsensusData(&packet)

	return nil
}

func (cph *ConsensusHandler) cleanupBlockState() {
	for key, blockStateDetails := range cph.blockStateDetailsMap {
		if blockStateDetails.parentHash.IsEqualTo(cph.currentParentHash) {
			continue
		}

		if Elapsed(blockStateDetails.initTime) >= BLOCK_CLEANUP_TIME_MS {
			delete(cph.blockStateDetailsMap, key)
		}
	}
}

func (cph *ConsensusHandler) HandleRequestConsensusDataPacket(packet *eth.RequestConsensusDataPacket) ([]*eth.ConsensusPacket, error) {

	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	cph.innerPacketLock.Lock()
	defer cph.innerPacketLock.Unlock()

	if packet == nil || packet.RequestData == nil || len(packet.RequestData) == 0 {
		return nil, errors.New("invalid request consensus data packet")
	}

	if cph.initialized == false || HasExceededTimeThreshold(cph.initTime, STARTUP_DELAY_MS) == false {
		return nil, errors.New("received request for consensus packet, but consensus is not ready yet")
	}

	requestDetails := RequestConsensusPacketDetails{}

	err := rlp.DecodeBytes(packet.RequestData, &requestDetails)
	if err != nil {
		log.Trace("handleProposeTransactionsPacket8", "err", err)
		return nil, err
	}

	blockStateDetails, ok := cph.blockStateDetailsMap[packet.ParentHash]
	if ok == false {
		return nil, errors.New("unknown parentHash")
	}

	var packets []*eth.ConsensusPacket
	packets = make([]*eth.ConsensusPacket, 0)

	var r byte
	proposalCount := 0
	ackCount := 0
	precommitCount := 0
	commitCount := 0
	for r = 1; r <= blockStateDetails.currentRound; r++ {
		blockRoundDetails := blockStateDetails.blockRoundMap[r]

		if blockRoundDetails.state >= BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS && blockRoundDetails.proposalPacket != nil {
			packets = append(packets, blockRoundDetails.proposalPacket)
			proposalCount = proposalCount + 1
		}

		if blockRoundDetails.state >= BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS {
			for _, pkt := range blockRoundDetails.proposalAckPackets {
				packets = append(packets, pkt)
				ackCount = ackCount + 1
			}
		}

		if blockRoundDetails.state >= BLOCK_STATE_WAITING_FOR_PRECOMMITS {
			for _, pkt := range blockRoundDetails.precommitPackets {
				packets = append(packets, pkt)
				precommitCount = precommitCount + 1
			}
		}

		if blockRoundDetails.state >= BLOCK_STATE_WAITING_FOR_COMMITS {
			for _, pkt := range blockRoundDetails.commitPackets {
				packets = append(packets, pkt)
				commitCount = commitCount + 1
			}
		}
	}

	log.Trace("HandleRequestConsensusDataPacket", "ParentHash", packet.ParentHash, "count", len(packets), "proposalCount", proposalCount, "ackCount", ackCount, "precommitCount", precommitCount, "commitCount", commitCount)

	return packets, nil
}

func getNilVoteProposalHash(parentHash common.Hash, round byte) common.Hash {
	return crypto.Keccak256Hash(parentHash.Bytes(), []byte("proposal"), ZERO_HASH.Bytes(), []byte{round}, []byte{byte(VOTE_TYPE_NIL)})
}

func getCommitHash(precommitHash common.Hash) common.Hash {
	return crypto.Keccak256Hash(precommitHash.Bytes())
}

func getOkVotePreCommitHash(parentHash common.Hash, proposalHash common.Hash, round byte) common.Hash {
	return crypto.Keccak256Hash(parentHash.Bytes(), proposalHash.Bytes(), []byte{round}, []byte{byte(VOTE_TYPE_OK)})
}

func getNilVotePreCommitHash(parentHash common.Hash, round byte) common.Hash {
	return crypto.Keccak256Hash(parentHash.Bytes(), []byte("precommit"), ZERO_HASH.Bytes(), []byte{round}, []byte{byte(VOTE_TYPE_NIL)})
}

func (cph *ConsensusHandler) LogIncomingPacketStats() {
	cph.packetStats.TotalIncomingPacketCount = cph.packetStats.TotalIncomingPacketCount + 1
	log.Trace("LogIncomingPacketStats", "TotalIncomingPacketCount", cph.packetStats.TotalIncomingPacketCount)
}
