package proofofstake

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/accounts"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/eth/protocols/eth"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/rlp"
	"math/big"
	"math/rand"
	"runtime/debug"
	"sort"
	"sync"
	"time"
)

type GetValidatorsFn func(blockHash common.Hash) (map[common.Address]*big.Int, error)
type DoesFinalizedTransactionExistFn func(txnHash common.Hash) (bool, error)

type OutOfOrderPacket struct {
	ReceivedTime time.Time
	Packet       *eth.ConsensusPacket
}

type ConsensusHandler struct {
	account                         accounts.Account
	signFn                          SignerFn
	p2pHandler                      P2PHandler
	blockStateDetailsMap            map[common.Hash]*BlockStateDetails
	outOfOrderPacketsMap            map[common.Hash][]*OutOfOrderPacket
	outerPacketLock                 sync.Mutex
	innerPacketLock                 sync.Mutex
	getValidatorsFn                 GetValidatorsFn
	doesFinalizedTransactionExistFn DoesFinalizedTransactionExistFn
	currentParentHash               common.Hash

	timeStatMap map[string]int

	nilVoteBlocks            uint64
	okVoteBlocks             uint64
	totalTransactions        uint64
	maxTransactionsInBlock   uint64
	maxTransactionsBlockTime int64
	initTime                 time.Time
	initialized              bool
}

type BlockConsensusData struct {
	BlockProposer          common.Address   `json:"blockProposer" gencodec:"required"`
	VoteType               VoteType         `json:"voteType" gencodec:"required"`
	ProposalHash           common.Hash      `json:"proposalHash" gencodec:"required"`
	PrecommitHash          common.Hash      `json:"precommitHash" gencodec:"required"`
	NilvotedBlockProposers []common.Address `json:"nilvotedBlockProposers" gencodec:"required"`
	Round                  byte
}

type BlockAdditionalConsensusData struct {
	ConsensusPackets []eth.ConsensusPacket `json:"consensusPackets" gencodec:"required"`
	InitTime         uint64                `json:"initTime" gencodec:"required"`
}

//todo: use mono clock

const BLOCK_TIMEOUT_MS = 9000 //relative to start of block locally
const REQUEST_CONSENSUS_DATA_PERCENT = 20
const BLOCK_CLEANUP_TIME_MS = 900
const MAX_ROUND_WITH_TXNS = 4

var STARTUP_DELAY_MS = int64(120000)

type BlockRoundState byte
type VoteType byte
type ConsensusPacketType byte
type RequestConsensusDataType byte

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
	MIN_BLOCK_VOTE_PERCENTAGE int64 = 70
	MIN_VALIDATORS            int   = 3
	MAX_VALIDATORS            int   = 128
)

var (
	MIN_VALIDATOR_DEPOSIT                               *big.Int       = big.NewInt(100000)
	MIN_BLOCK_DEPOSIT                                   *big.Int       = big.NewInt(1000000)
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

	selfCommited     bool
	selfCommitPacket *eth.ConsensusPacket

	proposer common.Address
}

type BlockStateDetails struct {
	filteredValidatorsDepositMap      map[common.Address]*big.Int
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
}

type ProposalDetails struct {
	Txns  []common.Hash `json:"Txns" gencodec:"required"`
	Round byte          `json:"Round" gencodec:"required"`
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

func getBlockProposer(parentHash common.Hash, filteredValidatorDepositMap *map[common.Address]*big.Int, round byte) (common.Address, error) {

	var proposer common.Address

	if len(*filteredValidatorDepositMap) < MIN_VALIDATORS {
		return proposer, errors.New("min validators not found")
	}

	validators := make([]common.Address, len(*filteredValidatorDepositMap))
	i := 0
	for k, _ := range *filteredValidatorDepositMap {
		validators[i].CopyFrom(k)
		//fmt.Println("getBlockProposer validator", k, "copied", validators[i])
		i = i + 1
	}

	sort.Slice(validators, func(i, j int) bool {
		vi := crypto.Keccak256Hash(parentHash.Bytes(), validators[i].Bytes(), []byte{round}).Bytes()
		vj := crypto.Keccak256Hash(parentHash.Bytes(), validators[j].Bytes(), []byte{round}).Bytes()
		return bytes.Compare(vi, vj) == -1
	})

	//fmt.Println("block proposer are 0", validators[0], "1", validators[1], "2", validators[2])
	proposer = validators[0]
	//fmt.Println("proposer", proposer, "round", round)

	return proposer, nil
}

func filterValidators(parentHash common.Hash, valDepMap *map[common.Address]*big.Int) (filteredValidators map[common.Address]bool, filteredDepositValue *big.Int, blockMinWeightedProposalsRequired *big.Int, err error) {
	validatorsDepositMap := *valDepMap

	totalDepositValue := big.NewInt(0)
	valCount := 0
	for val, depositValue := range validatorsDepositMap {
		if depositValue.Cmp(MIN_VALIDATOR_DEPOSIT) == -1 {
			fmt.Println("Skipping validator with low balance", val, depositValue)
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

func (cph *ConsensusHandler) initializeBlockStateIfRequired(parentHash common.Hash) error {
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
	}
	blockStateDetails := cph.blockStateDetailsMap[parentHash]

	validators, err := cph.getValidatorsFn(parentHash)
	if err != nil {
		fmt.Println("getValidatorsFn", err)
		delete(cph.blockStateDetailsMap, parentHash)
		return err
	}

	var filteredValidators map[common.Address]bool
	filteredValidators, blockStateDetails.totalBlockDepositValue, blockStateDetails.blockMinWeightedProposalsRequired, err = filterValidators(parentHash, &validators)
	if err != nil {
		delete(cph.blockStateDetailsMap, parentHash)
		return err
	}

	if blockStateDetails.totalBlockDepositValue.Cmp(MIN_BLOCK_DEPOSIT) == -1 {
		delete(cph.blockStateDetailsMap, parentHash)
		return errors.New("min block deposit not met")
	}

	//blockStateDetails.totalBlockDepositValue = big.NewInt(0)
	for addr, _ := range filteredValidators {
		depositValue := validators[addr]
		blockStateDetails.filteredValidatorsDepositMap[addr] = depositValue
		//blockStateDetails.totalBlockDepositValue = common.SafeAddBigInt(blockStateDetails.totalBlockDepositValue, depositValue)
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[cph.account.Address]
	if ok == false {
		fmt.Println("Not a validator in this block")
	}

	//blockStateDetails.blockMinWeightedProposalsRequired = common.SafeRelativePercentageBigInt(blockStateDetails.totalBlockDepositValue, MIN_BLOCK_TRANSACTION_WEIGHTED_PROPOSALS_PERCENTAGE)
	//fmt.Println("blockStateDetails.totalBlockDepositValue", blockStateDetails.totalBlockDepositValue)
	//fmt.Println("blockStateDetails.blockMinWeightedProposalsRequired", blockStateDetails.blockMinWeightedProposalsRequired)
	//fmt.Println("blockMinWeightedProposalsRequired", blockStateDetails.blockMinWeightedProposalsRequired)
	//fmt.Println("totalBlockDepositValue", blockStateDetails.totalBlockDepositValue,
	//"blockMinWeightedProposalsRequired", blockStateDetails.blockMinWeightedProposalsRequired)

	cph.blockStateDetailsMap[parentHash] = blockStateDetails
	cph.currentParentHash = parentHash

	err = cph.initializeNewBlockRound()
	if err != nil {
		delete(cph.blockStateDetailsMap, parentHash)
		return errors.New("min block deposit not met")
	}

	return nil
}

func (cph *ConsensusHandler) initializeNewBlockRound() error {
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
	}

	if blockRoundDetails.Round > 1 {
		//fmt.Println("initializeNewBlockRound", blockStateDetails.currentRound, cph.account.Address)
	}

	proposer, err := getBlockProposer(cph.currentParentHash, &blockStateDetails.filteredValidatorsDepositMap, blockRoundDetails.Round)
	if err != nil {
		return err
	}

	blockRoundDetails.proposer = proposer
	blockStateDetails.blockRoundMap[blockRoundDetails.Round] = blockRoundDetails
	blockStateDetails.currentRound = blockRoundDetails.Round
	cph.blockStateDetailsMap[cph.currentParentHash] = blockStateDetails

	return nil
}

func (cph *ConsensusHandler) isBlockProposer(parentHash common.Hash, filteredValidatorDepositMap *map[common.Address]*big.Int, round byte) (bool, error) {
	blockProposer, err := getBlockProposer(parentHash, filteredValidatorDepositMap, round)
	if err != nil {
		fmt.Println("isBlockProposer", err)
		return false, err
	}
	return blockProposer.IsEqualTo(cph.account.Address), nil
}

func (cph *ConsensusHandler) HandleConsensusPacket(packet *eth.ConsensusPacket) error {
	//fmt.Println("HandleConsensusPacket1", packet.ParentHash)
	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	if cph.signFn == nil {
		return nil
	}

	if cph.initialized == false || HasExceededTimeThreshold(cph.initTime, STARTUP_DELAY_MS) == false {
		return errors.New("startup delay")
	}

	if packet == nil {
		debug.PrintStack()
		panic("packet is nil")
	}

	err := cph.processPacket(packet)
	if err == OutOfOrderPackerErr {
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

	return err
}

func (cph *ConsensusHandler) processPacket(packet *eth.ConsensusPacket) error {
	if packet == nil {
		debug.PrintStack()
	}
	dataToVerify := append(packet.ParentHash.Bytes(), packet.ConsensusData...)
	digestHash := crypto.Keccak256(dataToVerify)
	pubKey, err := cryptobase.SigAlg.PublicKeyFromSignature(digestHash, packet.Signature)
	if err != nil {
		return InvalidPacketErr
	}
	if cryptobase.SigAlg.Verify(pubKey.PubData, digestHash, packet.Signature) == false {
		return InvalidPacketErr
	}

	validator, err := cryptobase.SigAlg.PublicKeyToAddress(pubKey)
	if err != nil {
		return InvalidPacketErr
	}

	packetType := ConsensusPacketType(packet.ConsensusData[0])
	if packetType == CONSENSUS_PACKET_TYPE_PROPOSE_BLOCK {
		return cph.handleProposeBlockPacket(validator, packet, false)
	} else if packetType == CONSENSUS_PACKET_TYPE_ACK_BLOCK_PROPOSAL {
		return cph.handleAckBlockProposalPacket(validator, packet)
	} else if packetType == CONSENSUS_PACKET_TYPE_PRECOMMIT_BLOCK {
		return cph.handlePrecommitPacket(validator, packet, false)
	} else if packetType == CONSENSUS_PACKET_TYPE_COMMIT_BLOCK {
		return cph.handleCommitPacket(validator, packet, false)
	}

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
					//fmt.Println("processOutOfOrderPackets 2")
				}
			} else {
				//fmt.Println("processOutOfOrderPackets 3")
			}
		}

		if len(unprocessedPackets) == 0 {
			//fmt.Println("processOutOfOrderPackets 4")
			delete(cph.outOfOrderPacketsMap, key)
		} else {
			//fmt.Println("processOutOfOrderPackets 5")
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

	blockConsensusData = &BlockConsensusData{
		VoteType:               blockRoundDetails.blockVoteType,
		NilvotedBlockProposers: make([]common.Address, 0),
		Round:                  blockStateDetails.currentRound,
	}
	if blockConsensusData.VoteType == VOTE_TYPE_OK {
		blockConsensusData.BlockProposer.CopyFrom(blockRoundDetails.proposer)
		blockConsensusData.ProposalHash.CopyFrom(blockRoundDetails.proposalHash)
	} else {
		blockConsensusData.BlockProposer.CopyFrom(ZERO_ADDRESS)
		blockConsensusData.ProposalHash.CopyFrom(getNilVoteProposalHash(parentHash, blockStateDetails.currentRound))
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

		roundProposer, err := getBlockProposer(parentHash, &blockStateDetails.filteredValidatorsDepositMap, r)
		if err != nil {
			return nil, nil, err
		}

		roundPktCount = roundPktCount + len(blockRoundDetails.proposalAckPackets) + len(blockRoundDetails.precommitPackets) + len(blockRoundDetails.commitPackets)
		if roundPktCount == 0 {
			fmt.Println("consensusdata", cph.account.Address, blockStateDetails.currentRound, r)
			return nil, nil, errors.New("no consensus packets for round")
		}

		if blockConsensusData.VoteType == VOTE_TYPE_NIL {
			blockConsensusData.NilvotedBlockProposers = append(blockConsensusData.NilvotedBlockProposers, roundProposer)
		} else {
			//if VoteType is VOTE_TUPE_OK, it means that all proposers less than currentRound will be NIL VOTED (except if only one round)
			if blockStateDetails.currentRound != byte(1) && r < blockStateDetails.currentRound {
				blockConsensusData.NilvotedBlockProposers = append(blockConsensusData.NilvotedBlockProposers, roundProposer)
			}
		}

		fmt.Println("===================>consensusdata", cph.account.Address, blockStateDetails.currentRound, r, roundPktCount)
	}

	blockAdditionalConsensusData.ConsensusPackets = make([]eth.ConsensusPacket, len(consensusPackets))
	for i, packet := range consensusPackets {
		blockAdditionalConsensusData.ConsensusPackets[i] = eth.NewConsensusPacket(&packet)
	}

	//packetRoundMap, err := ParseConsensusPackets(parentHash, &blockAdditionalConsensusData.ConsensusPackets, blockStateDetails.filteredValidatorsDepositMap)
	//if err != nil {
	//	return nil, nil, err
	//}

	//for r := byte(1); r <= blockRoundDetails.Round; r++ {
	//	_, ok := packetRoundMap[r]
	//	if ok == false {
	//		ParseConsensusPackets(parentHash, &blockAdditionalConsensusData.ConsensusPackets, blockStateDetails.filteredValidatorsDepositMap)
	//		return nil, nil, errors.New("packet not found")
	//	}
	//}

	if blockConsensusData.VoteType == VOTE_TYPE_NIL {
		err = ValidateBlockConsensusDataInner(nil, parentHash, blockConsensusData, blockAdditionalConsensusData, &blockStateDetails.filteredValidatorsDepositMap)
	} else {
		err = ValidateBlockConsensusDataInner(blockRoundDetails.proposalTxns, parentHash, blockConsensusData, blockAdditionalConsensusData, &blockStateDetails.filteredValidatorsDepositMap)
	}

	if err != nil {
		return nil, nil, err
	}

	return blockConsensusData, blockAdditionalConsensusData, nil
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

	return crypto.Keccak256Hash(data, parentHash.Bytes(), []byte{round})
}

func (cph *ConsensusHandler) handleProposeBlockPacket(validator common.Address, packet *eth.ConsensusPacket, self bool) error {
	cph.innerPacketLock.Lock()
	defer cph.innerPacketLock.Unlock()

	//fmt.Println("validator proposal", validator, "self", cph.account.Address)
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
		//fmt.Println("handleProposeTransactionsPacket8", err)
		return err
	}

	if proposalDetails.Round != blockRoundDetails.Round {
		return OutOfOrderPackerErr
	}

	if blockRoundDetails.state >= BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS {
		//fmt.Println("handleProposeBlockPacket BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS")
		return OutOfOrderPackerErr
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[validator]
	if ok == false {
		//fmt.Println("handleProposeTransactionsPacket6")
		return errors.New("invalid validator")
	}

	if blockRoundDetails.proposer.IsEqualTo(validator) == false {
		return errors.New("invalid proposer")
	}

	if validator.IsEqualTo(cph.account.Address) == true && self == false {
		return errors.New("self packet from elsewhere")
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
			fmt.Println("doesFinalizedTransactionExistFn", err)
			return err
		}
		if exists {
			fmt.Println("doesFinalizedTransactionExistFn true", proposalDetails.Txns[i].Hex())
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
				//fmt.Println("------------------------->unknown txn", proposalDetails.Txns[i], "validator", cph.account.Address)
			} else {
				//fmt.Println("known txn", proposalDetails.Txns[i], "validator", cph.account.Address)
			}
		}
		if len(unknownTxns) > 0 {
			err = cph.p2pHandler.RequestTransactions(unknownTxns)
			if err != nil {
				//fmt.Println("RequestTransactions error", err)
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

	if Elapsed(blockStateDetails.initTime) > BLOCK_TIMEOUT_MS*2 {
		//fmt.Println("handleAckBlockProposalPacket", "me", cph.account.Address, "validator", validator)
	}

	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	_, ok = blockRoundDetails.validatorProposalAcks[validator]
	if ok == true {
		//return errors.New("already received proposal")
		//todo: compare
	} else {
		if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS {
			if Elapsed(blockStateDetails.initTime) > BLOCK_TIMEOUT_MS*3 {
				//fmt.Println("proposalAckDetails new")
			}
		}
	}

	proposalAckDetails := &ProposalAckDetails{}

	err := rlp.DecodeBytes(packet.ConsensusData[1:], proposalAckDetails)
	if err != nil {
		//fmt.Println("handleAckBlockProposalPacket err5", err)
		return err
	}

	if proposalAckDetails.Round != blockStateDetails.currentRound {
		if proposalAckDetails.Round > blockStateDetails.currentRound {
			blockStateDetails.highestProposalRoundSeen = proposalAckDetails.Round
			cph.blockStateDetailsMap[packet.ParentHash] = blockStateDetails
		}
		return OutOfOrderPackerErr
	}

	if proposalAckDetails.ProposalAckVoteType != VOTE_TYPE_OK && proposalAckDetails.ProposalAckVoteType != VOTE_TYPE_NIL {
		fmt.Println("proposalAckDetails.ProposalAckVoteType", proposalAckDetails.ProposalAckVoteType)
		return errors.New("invalid vote type")
	}

	if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS || blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PRECOMMITS {
		//fmt.Println("handleAckBlockProposalPacket waiting", cph.account.Address)
	} else if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL {
		//if proposalAckDetails.ProposalAckVoteType != VOTE_TYPE_NIL {
		//return errors.New("invalid state")
		//}
		//return errors.New("invalid state")
	} else if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_COMMITS {
		//return errors.New("invalid state")
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
		fmt.Println("invalid 1", err)
		return 0, ZERO_ADDRESS, err
	}
	if cryptobase.SigAlg.Verify(pubKey.PubData, digestHash, packet.Signature) == false {
		fmt.Println("invalid 2")
		return 0, ZERO_ADDRESS, InvalidPacketErr
	}

	validator, err := cryptobase.SigAlg.PublicKeyToAddress(pubKey)
	if err != nil {
		fmt.Println("invalid 3", err)
		return 0, ZERO_ADDRESS, err
	}

	packetType := ConsensusPacketType(packet.ConsensusData[0])
	if packetType == CONSENSUS_PACKET_TYPE_PROPOSE_BLOCK {
		details := ProposalDetails{}

		err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
		if err != nil {
			fmt.Println("invalid 4", err)
			return 0, ZERO_ADDRESS, err
		}

		return details.Round, validator, nil
	} else if packetType == CONSENSUS_PACKET_TYPE_ACK_BLOCK_PROPOSAL {
		details := ProposalAckDetails{}

		err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
		if err != nil {
			fmt.Println("invalid 5", err)
			return 0, ZERO_ADDRESS, err
		}

		return details.Round, validator, nil
	} else if packetType == CONSENSUS_PACKET_TYPE_PRECOMMIT_BLOCK {
		details := PreCommitDetails{}

		err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
		if err != nil {
			fmt.Println("invalid 6", err)
			return 0, ZERO_ADDRESS, err
		}

		return details.Round, validator, nil
	} else if packetType == CONSENSUS_PACKET_TYPE_COMMIT_BLOCK {
		details := CommitDetails{}

		err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
		if err != nil {
			fmt.Println("invalid 7", err)
			return 0, ZERO_ADDRESS, err
		}

		return details.Round, validator, nil
	}

	fmt.Println("invalid 8", err, "packetType", packetType)

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

func (cph *ConsensusHandler) shouldMoveToNextRound(parentHash common.Hash) (bool, error) {
	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	//Find validators in greater rounds
	valMap := make(map[common.Address]bool)
	for _, pktList := range cph.outOfOrderPacketsMap {
		for _, pkt := range pktList {
			if pkt.Packet.ParentHash.IsEqualTo(parentHash) {
				round, validator, err := parsePacket(pkt.Packet)
				if err != nil {
					fmt.Println("parsePacket", "err", err)
					return false, err
				}
				if round <= blockStateDetails.currentRound {
					continue
				}
				valMap[validator] = true
				//fmt.Println("shouldMoveToNextRound", "valInGreaterRound", validator)
			}
		}
	}

	totalGreaterRoundDepositCount := big.NewInt(0)
	currentRoundDepositSoFar := big.NewInt(0)
	for val, depositAmount := range blockStateDetails.filteredValidatorsDepositMap {
		_, ok := valMap[val]
		if ok == false {
			_, ok1 := blockRoundDetails.validatorPrecommits[val]
			if ok1 == true {
				currentRoundDepositSoFar = common.SafeAddBigInt(depositAmount, currentRoundDepositSoFar)
				//fmt.Println("currentRoundDepositSoFar", "val", val, "depositAmount", depositAmount, "currentRoundDepositSoFar", currentRoundDepositSoFar)
			} else {
				//fmt.Println("currentRoundDepositSoFar val no", val)
			}
		} else {
			totalGreaterRoundDepositCount = common.SafeAddBigInt(depositAmount, totalGreaterRoundDepositCount)
			//fmt.Println("totalGreaterRoundDepositCount", "val", val, "depositAmount", depositAmount, "totalGreaterRoundDepositCount", totalGreaterRoundDepositCount)
		}
	}

	if currentRoundDepositSoFar.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 {
		return false, nil
	}

	//If there are votes in greater rounds, so currentRound can never get that much deposit
	balanceDepositVotesRequiredCurrentRound := common.SafeSubBigInt(blockStateDetails.blockMinWeightedProposalsRequired, currentRoundDepositSoFar)
	/*fmt.Println("shouldMoveNextRound",
	"blockMinWeightedProposalsRequired", blockStateDetails.blockMinWeightedProposalsRequired,
	"balanceDepositVotesRequiredCurrentRound", balanceDepositVotesRequiredCurrentRound,
	"currentRoundDepositSoFar", currentRoundDepositSoFar, "totalGreaterRoundDepositCount", totalGreaterRoundDepositCount,
	"precommit count", len(blockRoundDetails.validatorPrecommits),
	"self precomitted", blockRoundDetails.selfPrecommited,
	"val", cph.account.Address)*/
	if totalGreaterRoundDepositCount.Cmp(balanceDepositVotesRequiredCurrentRound) >= 0 {
		return true, nil
	}

	return false, nil
}

func (cph *ConsensusHandler) handlePrecommitPacket(validator common.Address, packet *eth.ConsensusPacket, self bool) error {
	cph.innerPacketLock.Lock()
	defer cph.innerPacketLock.Unlock()

	if self == false {
		//fmt.Println("precommit other")
	}

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
		//fmt.Println("handlePrecommitPacket BLOCK_STATE_WAITING_FOR_PRECOMMITS")
		return OutOfOrderPackerErr
	}

	_, ok = blockStateDetails.filteredValidatorsDepositMap[validator]
	if ok == false {
		//fmt.Println("handleProposeTransactionsPacket6")
		return errors.New("invalid validator")
	}

	if validator.IsEqualTo(cph.account.Address) == true && self == false {
		return errors.New("self packet from elsewhere")
	}

	_, ok = blockRoundDetails.validatorPrecommits[validator]
	if ok == true {
		//return errors.New("already received precommit")
		//todo: check
	} else {
		//fmt.Println("precommit")
	}

	_, ok = blockRoundDetails.validatorProposalAcks[validator]
	if ok == false {
		//fmt.Println("did not receive proposal ack from validator")
		//return OutOfOrderPackerErr
	}

	precommitDetails := &PreCommitDetails{}

	err := rlp.DecodeBytes(packet.ConsensusData[1:], precommitDetails)
	if err != nil {
		//fmt.Println("handlePrecommitPacket err5", err)
		return err
	}

	if precommitDetails.Round != blockStateDetails.currentRound {
		return OutOfOrderPackerErr
	}

	if precommitDetails.PrecommitHash.IsEqualTo(blockRoundDetails.precommitHash) == false {
		fmt.Println("precommit error", "incoming", precommitDetails.PrecommitHash, "expected", blockRoundDetails.precommitHash, "me", cph.account.Address, "validator", validator)
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
		}

		//totalVotesPercentage := common.SafePercentageOfBigInt(totalVotesDepositCount, blockStateDetails.totalBlockDepositValue)
		if totalVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 {
			blockStateDetails.precommitTime = Elapsed(blockStateDetails.initTime)
			blockRoundDetails.state = BLOCK_STATE_WAITING_FOR_COMMITS
		}
	}

	//fmt.Println("precommit", "totalVotesPercentage", totalVotesPercentage)

	pkt := eth.NewConsensusPacket(packet)
	blockRoundDetails.precommitPackets[validator] = &pkt

	blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
	cph.blockStateDetailsMap[packet.ParentHash] = blockStateDetails
	//fmt.Println("handlePrecommitPacket done", packet.ParentHash)

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
		//fmt.Println("handleProposeTransactionsPacket6")
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
		//fmt.Println("handlePrecommitPacket err5", err)
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
		}

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

			fmt.Println("BlockStats", "maxTxnsInBlock", cph.maxTransactionsInBlock, "totalTxns", cph.totalTransactions, "okBlocks", cph.okVoteBlocks, "nilBlocks", cph.nilVoteBlocks)
			for statKey, statVal := range cph.timeStatMap {
				if statVal > 0 {
					fmt.Println("     TimeStatsBlockCount", "stat", statKey, "blocks", statVal)
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

func (cph *ConsensusHandler) proposeBlock(parentHash common.Hash, txns []common.Hash) error {
	var packet *eth.ConsensusPacket
	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if blockRoundDetails.selfProposed == true {
		return cph.broadCast(blockRoundDetails.selfProposalPacket)
	}

	proposalDetails := &ProposalDetails{}

	proposalDetails.Round = blockStateDetails.currentRound
	if blockStateDetails.currentRound < MAX_ROUND_WITH_TXNS { //No transactions after this round, to reduce chance of FLP
		proposalDetails.Txns = make([]common.Hash, len(txns))
		for i := 0; i < len(proposalDetails.Txns); i++ {
			proposalDetails.Txns[i].CopyFrom(txns[i])
		}
	} else {
		proposalDetails.Txns = make([]common.Hash, 0)
	}
	//fmt.Println("ProposeBlock with txns", len(proposalDetails.Txns))

	data, err := rlp.EncodeToBytes(proposalDetails)

	if err != nil {
		return err
	}

	dataToSend := append([]byte{byte(CONSENSUS_PACKET_TYPE_PROPOSE_BLOCK)}, data...)

	packet, err = cph.createConsensusPacket(parentHash, dataToSend)
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
		packet, err := cph.createConsensusPacket(parentHash, dataToSend)
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
	}

	okVotes := 0
	nilVotes := 0
	mismatchedVotes := 0
	okVotesDepositCount := big.NewInt(0)
	nilVotesDepositCount := big.NewInt(0)
	totalVotesDepositCount := big.NewInt(0)

	for val, ack := range blockRoundDetails.validatorProposalAcks {
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

	//fmt.Println("timeout", "okVotesPercentage", okVotesPercentage, "nilVotesPercentage", nilVotesPercentage, "totalVotesPercentage", totalVotesPercentage)

	if okVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 && blockRoundDetails.selfAckProposalVoteType == VOTE_TYPE_OK {
		blockRoundDetails.state = BLOCK_STATE_WAITING_FOR_PRECOMMITS
		blockRoundDetails.precommitHash.CopyFrom(getOkVotePreCommitHash(parentHash, blockRoundDetails.proposalHash, blockStateDetails.currentRound))
	} else if nilVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 { //handle timeout differently?
		//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>NilVote 2")
		blockRoundDetails.state = BLOCK_STATE_WAITING_FOR_PRECOMMITS
		blockRoundDetails.precommitHash.CopyFrom(getNilVotePreCommitHash(parentHash, blockStateDetails.currentRound))
	} else {
		if totalVotesDepositCount.Cmp(blockStateDetails.totalBlockDepositValue) >= 0 ||
			totalVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 && HasExceededTimeThreshold(blockRoundDetails.initTime,
				int64(BLOCK_TIMEOUT_MS*int(blockRoundDetails.Round)*2)) {
			blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
			cph.blockStateDetailsMap[parentHash] = blockStateDetails
			err := cph.initializeNewBlockRound()
			if err != nil {
				return err
			}
			return nil
		}
	}

	blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
	cph.blockStateDetailsMap[parentHash] = blockStateDetails

	return cph.broadCast(blockRoundDetails.selfAckPacket)
}

func (cph *ConsensusHandler) ackBlockProposal(parentHash common.Hash) error {
	//fmt.Println("ackBlockProposal", cph.account.Address)
	blockStateDetails := cph.blockStateDetailsMap[parentHash]
	blockRoundDetails := blockStateDetails.blockRoundMap[blockStateDetails.currentRound]

	if blockRoundDetails.selfAckd == true {
		shouldPropose, err := cph.isBlockProposer(parentHash, &blockStateDetails.filteredValidatorsDepositMap, blockStateDetails.currentRound)
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

		if blockStateDetails.currentRound >= MAX_ROUND_WITH_TXNS && len(blockRoundDetails.blockProposalDetails.Txns) > 0 {
			return errors.New("unexpected transaction count")
		}

		//Find if any new transactions we don't know yet
		unknownTxns := make([]common.Hash, 0)
		for i := 0; i < len(blockRoundDetails.blockProposalDetails.Txns); i++ {
			_, txnExists := blockRoundDetails.selfKnownTransactions[blockRoundDetails.blockProposalDetails.Txns[i]]
			if txnExists == false {
				//fmt.Println("===============ackBlockProposal unknown txns", blockRoundDetails.blockProposalDetails.Txns[i])
				unknownTxns = append(unknownTxns, blockRoundDetails.blockProposalDetails.Txns[i])
			}
		}
		if len(unknownTxns) > 0 {
			err := cph.p2pHandler.RequestTransactions(unknownTxns)
			if err != nil {
				return errors.New("unknown transactions")
			}
		}

		proposalAckDetails := &ProposalAckDetails{
			ProposalHash:        blockRoundDetails.proposalHash,
			ProposalAckVoteType: VOTE_TYPE_OK,
			Round:               blockStateDetails.currentRound,
		}

		data, err := rlp.EncodeToBytes(proposalAckDetails)

		if err != nil {
			return err
		}

		dataToSend := append([]byte{byte(CONSENSUS_PACKET_TYPE_ACK_BLOCK_PROPOSAL)}, data...)
		packet, err := cph.createConsensusPacket(parentHash, dataToSend)
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
			fmt.Println("unexpected")
			return errors.New("unexpected")
		}
	}

	//fmt.Println("ack", "totalBlockDepositValue", blockStateDetails.totalBlockDepositValue, "okVotesDepositCount",
	//	okVotesDepositCount, "okVotesPercentage", okVotesPercentage, "nilVotesPercentage", nilVotesPercentage, "totalVotesPercentage", totalVotesPercentage)

	if okVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 && blockRoundDetails.selfAckProposalVoteType == VOTE_TYPE_OK { //For ok votes, vote type should match
		blockRoundDetails.state = BLOCK_STATE_WAITING_FOR_PRECOMMITS
		blockRoundDetails.precommitHash.CopyFrom(getOkVotePreCommitHash(parentHash, blockRoundDetails.proposalHash, blockStateDetails.currentRound))
		blockRoundDetails.blockVoteType = VOTE_TYPE_OK
		blockStateDetails.ackProposalTime = Elapsed(blockStateDetails.initTime)
	} else if nilVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 { //handle timeout differently? for nil votes, it is ok to accept NIL vote even if self vote is OK
		blockRoundDetails.state = BLOCK_STATE_WAITING_FOR_PRECOMMITS
		blockRoundDetails.precommitHash.CopyFrom(getNilVotePreCommitHash(parentHash, blockStateDetails.currentRound))
		blockRoundDetails.blockVoteType = VOTE_TYPE_NIL
	} else {
		if totalVotesDepositCount.Cmp(blockStateDetails.totalBlockDepositValue) >= 0 ||
			totalVotesDepositCount.Cmp(blockStateDetails.blockMinWeightedProposalsRequired) >= 0 && HasExceededTimeThreshold(blockRoundDetails.initTime, int64(BLOCK_TIMEOUT_MS*int(blockRoundDetails.Round)*2)) {
			blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
			cph.blockStateDetailsMap[parentHash] = blockStateDetails
			err := cph.initializeNewBlockRound()
			if err != nil {
				return err
			}
			return nil
		}
	}

	blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
	cph.blockStateDetailsMap[parentHash] = blockStateDetails

	return cph.broadCast(blockRoundDetails.selfAckPacket)
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
	packet, err := cph.createConsensusPacket(parentHash, dataToSend)
	if err != nil {
		return err
	}

	err = cph.handlePrecommitPacket(cph.account.Address, packet, true)

	if err != nil {
		//fmt.Println("precommitBlock handlePrecommitPacket error", err)
		return err
	}

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
	packet, err := cph.createConsensusPacket(parentHash, dataToSend)
	if err != nil {
		return err
	}

	err = cph.handleCommitPacket(cph.account.Address, packet, true)

	if err != nil {
		//fmt.Println("commitBlock handleCommitPacket error", err)
		return err
	}

	return cph.broadCast(packet)
}

func (cph *ConsensusHandler) HandleTransactions(parentHash common.Hash, txns []common.Hash) error {
	cph.outerPacketLock.Lock()
	defer cph.outerPacketLock.Unlock()

	if cph.initialized == false {
		cph.initTime = time.Now()
		cph.initialized = true
		log.Info("Starting up...")
		return errors.New("starting up")
	}

	if HasExceededTimeThreshold(cph.initTime, STARTUP_DELAY_MS) == false {
		log.Info("Waiting to startup...", "elapsed ms", Elapsed(cph.initTime))
		fmt.Println("Waiting to startup...", "elapsed ms", Elapsed(cph.initTime))
		return errors.New("starting up")
	}

	err := cph.initializeBlockStateIfRequired(parentHash)
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

	shouldPropose, err := cph.isBlockProposer(parentHash, &blockStateDetails.filteredValidatorsDepositMap, blockStateDetails.currentRound)
	if err != nil {
		return err
	}

	cph.processOutOfOrderPackets(parentHash)

	err = errors.New("not ready yet")
	if shouldPropose {
		fmt.Println("parentHash", parentHash, "round", blockStateDetails.currentRound, "state", blockRoundDetails.state, "vote", blockRoundDetails.blockVoteType,
			"shouldPropose", shouldPropose, "currTxns", len(txns), "okBlocks", cph.okVoteBlocks, "nilBlocks", cph.nilVoteBlocks,
			"totalTxs", cph.totalTransactions, "maxTxns", cph.maxTransactionsInBlock, "blockTIme", cph.maxTransactionsBlockTime, "txns", len(txns))
	}

	if blockRoundDetails.state == BLOCK_STATE_WAITING_FOR_PROPOSAL {
		for _, txn := range txns {
			blockRoundDetails.selfKnownTransactions[txn] = true
		}
		blockStateDetails.blockRoundMap[blockStateDetails.currentRound] = blockRoundDetails
		cph.blockStateDetailsMap[parentHash] = blockStateDetails

		if shouldPropose {
			cph.proposeBlock(parentHash, txns)
		} else {
			if HasExceededTimeThreshold(blockRoundDetails.initTime, int64(BLOCK_TIMEOUT_MS*int(blockRoundDetails.Round))) {
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
		shouldMove, err := cph.shouldMoveToNextRound(parentHash)
		if err == nil && shouldMove {
			return cph.initializeNewBlockRound()
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

func (cph *ConsensusHandler) createConsensusPacket(parentHash common.Hash, data []byte) (*eth.ConsensusPacket, error) {
	if cph.signFn == nil {
		return nil, errors.New("signFn is not set")
	}
	dataToSign := append(parentHash.Bytes(), data...)
	signature, err := cph.signFn(cph.account, accounts.MimetypeProofOfStake, dataToSign)
	if err != nil {
		////fmt.Println("signAndSend failed", err)
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

func (cph *ConsensusHandler) broadCast(packet *eth.ConsensusPacket) error {
	if packet == nil {
		debug.PrintStack()
		panic("packet is nil")
	}
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
	elapsed := Elapsed(blockStateDetails.initTime)
	if elapsed < BLOCK_TIMEOUT_MS {
		return nil
	}

	r := rand.Intn(100)
	//fmt.Println("r", r, "elapsed", elapsed, "timeout", BLOCK_TIMEOUT_MS, "state", blockStateDetails.state)
	if r > REQUEST_CONSENSUS_DATA_PERCENT {
		return nil
	}

	//fmt.Println("requestConsensusData 1")
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

	if cph.initialized == false || HasExceededTimeThreshold(cph.initTime, STARTUP_DELAY_MS) == false {
		return nil, errors.New("startup delay")
	}

	requestDetails := RequestConsensusPacketDetails{}

	err := rlp.DecodeBytes(packet.RequestData, &requestDetails)
	if err != nil {
		//fmt.Println("handleProposeTransactionsPacket8", err)
		return nil, err
	}

	blockStateDetails, ok := cph.blockStateDetailsMap[packet.ParentHash]
	if ok == false {
		return nil, errors.New("unknown parentHash")
	}

	var packets []*eth.ConsensusPacket
	packets = make([]*eth.ConsensusPacket, 0)

	var r byte
	for r = 1; r <= blockStateDetails.currentRound; r++ {
		blockRoundDetails := blockStateDetails.blockRoundMap[r]

		if blockRoundDetails.state >= BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS && blockRoundDetails.proposalPacket != nil {
			packets = append(packets, blockRoundDetails.proposalPacket)
		}

		if blockRoundDetails.state >= BLOCK_STATE_WAITING_FOR_PROPOSAL_ACKS {
			for _, pkt := range blockRoundDetails.proposalAckPackets {
				packets = append(packets, pkt)
			}
		}

		if blockRoundDetails.state >= BLOCK_STATE_WAITING_FOR_PRECOMMITS {
			for _, pkt := range blockRoundDetails.precommitPackets {
				packets = append(packets, pkt)
			}
		}

		if blockRoundDetails.state >= BLOCK_STATE_WAITING_FOR_COMMITS {
			for _, pkt := range blockRoundDetails.commitPackets {
				packets = append(packets, pkt)
			}
		}
	}

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
