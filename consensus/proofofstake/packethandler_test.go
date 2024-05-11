package proofofstake

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/accounts"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"github.com/DogeProtocol/dp/eth/protocols/eth"
	"github.com/DogeProtocol/dp/params"
	"github.com/DogeProtocol/dp/rlp"
	"math/big"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const TEST_ITERATIONS int = 1

var waitLock sync.Mutex
var waitMap map[common.Address]bool
var packetDropCount int32
var packetSentCount int32
var TEST_CONSENSUS_BLOCK_NUMBER = uint64(1)

type ValidatorDetailsTest struct {
	balance *big.Int
	key     *signaturealgorithm.PrivateKey
}

type ValidatorManager struct {
	valMap map[common.Address]*ValidatorDetailsTest
}

type MockP2PManager struct {
	packetMutex                     sync.Mutex
	txnMapMutex                     sync.Mutex
	mockP2pHandlers                 map[common.Address]*MockP2PHandler
	blockPacketValidatorMap         map[common.Address]bool
	blockPacketsBetweenValidatorMap map[common.Hash]bool
	txnFinalize                     map[common.Hash]bool
}

type MockNetworkDetails struct {
	packetLoss int
	latencyMs  int
}

type MockP2PHandler struct {
	mockLock         sync.Mutex
	mockP2pManager   *MockP2PManager
	validator        common.Address
	validatorTxns    []common.Hash
	consensusHandler *ConsensusHandler
	validatorDetails *ValidatorDetailsTest
	networkDetails   MockNetworkDetails
}

func (m *MockP2PManager) DoesFinalizedTransactionExistFn(txnHash common.Hash) (bool, error) {
	m.txnMapMutex.Lock()
	defer m.txnMapMutex.Unlock()

	exists, ok := m.txnFinalize[txnHash]
	return exists && ok, nil
}

func (m *MockP2PManager) SetTransactionFinalizeState(txnHash common.Hash, state bool) {
	m.txnMapMutex.Lock()
	defer m.txnMapMutex.Unlock()

	m.txnFinalize[txnHash] = state
}

func (m *MockP2PManager) BlockPacketsBetweenValidators(val1 common.Address, val2 common.Address) {
	m.packetMutex.Lock()
	defer m.packetMutex.Unlock()
	hash1 := crypto.Keccak256Hash(val1.Bytes(), val2.Bytes())
	hash2 := crypto.Keccak256Hash(val2.Bytes(), val1.Bytes())
	m.blockPacketsBetweenValidatorMap[hash1] = true
	m.blockPacketsBetweenValidatorMap[hash2] = true
}

func (m *MockP2PManager) UnblockPacketsBetweenValidators(val1 common.Address, val2 common.Address) {
	m.packetMutex.Lock()
	defer m.packetMutex.Unlock()
	hash1 := crypto.Keccak256Hash(val1.Bytes(), val2.Bytes())
	hash2 := crypto.Keccak256Hash(val2.Bytes(), val1.Bytes())
	m.blockPacketsBetweenValidatorMap[hash1] = false
	m.blockPacketsBetweenValidatorMap[hash2] = false
}

func (m *MockP2PManager) ArePacketsBetweenValidatorsBlocked(val1 common.Address, val2 common.Address) bool {
	m.packetMutex.Lock()
	defer m.packetMutex.Unlock()
	hash1 := crypto.Keccak256Hash(val1.Bytes(), val2.Bytes())
	hash2 := crypto.Keccak256Hash(val2.Bytes(), val1.Bytes())
	ok1 := m.blockPacketsBetweenValidatorMap[hash1]
	ok2 := m.blockPacketsBetweenValidatorMap[hash2]

	if ok1 || ok2 {
		return true
	}

	return false
}

func (m *MockP2PManager) DeleteAllPacketBlocks() bool {
	m.packetMutex.Lock()
	defer m.packetMutex.Unlock()

	for k, _ := range m.blockPacketsBetweenValidatorMap {
		delete(m.blockPacketsBetweenValidatorMap, k)
	}

	for k, _ := range m.blockPacketValidatorMap {
		delete(m.blockPacketValidatorMap, k)
	}

	return false
}

func (m *MockP2PManager) BlockValidatorPackets(val common.Address) {
	m.packetMutex.Lock()
	defer m.packetMutex.Unlock()
	m.blockPacketValidatorMap[val] = true
}

func (m *MockP2PManager) UnblockValidatorPackets(val common.Address) {
	m.packetMutex.Lock()
	defer m.packetMutex.Unlock()
	m.blockPacketValidatorMap[val] = false
}

func (m *MockP2PManager) IsValidatorPacketsBlocked(val common.Address) bool {
	m.packetMutex.Lock()
	defer m.packetMutex.Unlock()
	return m.blockPacketValidatorMap[val]
}

func getSigner(packet *eth.ConsensusPacket) (common.Address, error) {
	dataToVerify := append(packet.ParentHash.Bytes(), packet.ConsensusData...)
	digestHash := crypto.Keccak256(dataToVerify)

	packetType := ConsensusPacketType(packet.ConsensusData[0])
	if shouldSignFull(TEST_CONSENSUS_BLOCK_NUMBER) && packetType == CONSENSUS_PACKET_TYPE_PROPOSE_BLOCK && packet.ParentHash.IsEqualTo(getTestParentHash(TEST_CONSENSUS_BLOCK_NUMBER)) {
		pubKey, err := cryptobase.SigAlg.PublicKeyFromSignatureWithContext(digestHash, packet.Signature, FULL_SIGN_CONTEXT)
		if err != nil {
			return ZERO_ADDRESS, err
		}
		if cryptobase.SigAlg.VerifyWithContext(pubKey.PubData, digestHash, packet.Signature, FULL_SIGN_CONTEXT) == false {
			return ZERO_ADDRESS, InvalidPacketErr
		}

		validator, err := cryptobase.SigAlg.PublicKeyToAddress(pubKey)
		if err != nil {
			return ZERO_ADDRESS, err
		}

		return validator, nil
	} else {
		pubKey, err := cryptobase.SigAlg.PublicKeyFromSignature(digestHash, packet.Signature)
		if err != nil {
			return ZERO_ADDRESS, err
		}
		if cryptobase.SigAlg.Verify(pubKey.PubData, digestHash, packet.Signature) == false {
			return ZERO_ADDRESS, InvalidPacketErr
		}

		validator, err := cryptobase.SigAlg.PublicKeyToAddress(pubKey)
		if err != nil {
			return ZERO_ADDRESS, err
		}

		return validator, nil
	}
}

func (p *MockP2PHandler) BroadcastConsensusData(packet *eth.ConsensusPacket) error {
	for _, val := range p.mockP2pManager.mockP2pHandlers {
		handler := val.consensusHandler
		if bytes.Compare(handler.account.Address.Bytes(), p.validator.Bytes()) != 0 {
			if p.networkDetails.packetLoss > 0 {
				r := rand.Intn(100)
				if r <= p.networkDetails.packetLoss {
					atomic.AddInt32(&packetDropCount, 1)
					continue
				}
				atomic.AddInt32(&packetSentCount, 1)
			}
			handler := val.consensusHandler
			if p.mockP2pManager.IsValidatorPacketsBlocked(handler.account.Address) {
				continue
			}
			signer, err := getSigner(packet)
			if err != nil {
				panic("unexpected")
			}
			if p.mockP2pManager.IsValidatorPacketsBlocked(signer) {
				continue
			}
			if p.mockP2pManager.ArePacketsBetweenValidatorsBlocked(handler.account.Address, p.validator) {
				continue
			}
			if p.mockP2pManager.ArePacketsBetweenValidatorsBlocked(handler.account.Address, signer) {
				continue
			}
			if p.mockP2pManager.ArePacketsBetweenValidatorsBlocked(p.validator, signer) {
				continue
			}
			err = handler.HandleConsensusPacket(packet)
			if err != nil {
				continue
			}
		}
	}

	//fmt.Println("packets dropped", packetDropCount, "packets sent", packetSentCount)
	return nil
}

func (p *MockP2PHandler) SetValidatorTransactions(txns []common.Hash) {
	p.mockLock.Lock()
	defer p.mockLock.Unlock()
	p.validatorTxns = make([]common.Hash, len(txns))
	for i := 0; i < len(txns); i++ {
		p.validatorTxns[i].CopyFrom(txns[i])
	}
}

func (p *MockP2PHandler) AppendValidatorTransactions(txns []common.Hash) {
	p.mockLock.Lock()
	defer p.mockLock.Unlock()
	p.validatorTxns = append(p.validatorTxns, txns...)
}

func (p *MockP2PHandler) GetValidatorTransactions() []common.Hash {
	p.mockLock.Lock()
	defer p.mockLock.Unlock()
	txns := make([]common.Hash, len(p.validatorTxns))
	for i := 0; i < len(txns); i++ {
		txns[i].CopyFrom(p.validatorTxns[i])
	}
	return txns
}

func (p *MockP2PHandler) RequestTransactions(txns []common.Hash) error {
	p.mockLock.Lock()
	defer p.mockLock.Unlock()

	for _, val := range p.mockP2pManager.mockP2pHandlers {
		handler := val.consensusHandler

		if handler.account.Address.IsEqualTo(p.validator) == false {
			if p.networkDetails.packetLoss > 0 {
				r := rand.Intn(100)
				if r <= p.networkDetails.packetLoss {
					continue
				}
			}
			val.AppendValidatorTransactions(txns)
		}
	}
	return nil
}

func (p *MockP2PHandler) RequestConsensusData(packet *eth.RequestConsensusDataPacket) error {
	//p.mockLock.Lock()
	//defer p.mockLock.Unlock()

	for _, val := range p.mockP2pManager.mockP2pHandlers {
		handler := val.consensusHandler

		if handler.account.Address.IsEqualTo(p.validator) == false {
			if p.networkDetails.packetLoss > 0 {
				r := rand.Intn(100)
				if r <= p.networkDetails.packetLoss {
					continue
				}
			}
			if p.mockP2pManager.IsValidatorPacketsBlocked(val.validator) {
				continue
			}
			if p.mockP2pManager.ArePacketsBetweenValidatorsBlocked(val.validator, p.validator) {
				continue
			}
			consensusPackets, err := handler.HandleRequestConsensusDataPacket(packet)
			if err != nil {
				return err
			}

			for _, pkt := range consensusPackets {
				p.consensusHandler.HandleConsensusPacket(pkt)
			}
		}
	}

	return nil
}

func NewValidatorManager(numKeys int) *ValidatorManager {
	valManager := &ValidatorManager{}
	valManager.valMap = make(map[common.Address]*ValidatorDetailsTest)

	for i := 0; i < numKeys; i++ {
		valKey, _ := cryptobase.SigAlg.GenerateKey()
		valAddress, _ := cryptobase.SigAlg.PublicKeyToAddress(&valKey.PublicKey)
		valManager.valMap[valAddress] = &ValidatorDetailsTest{
			key:     valKey,
			balance: params.EtherToWei(big.NewInt(500000000000)),
		}
	}

	return valManager
}

func (vm *ValidatorManager) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	val, ok := vm.valMap[account.Address]
	// If the key exists
	if ok == false {
		return nil, errors.New("validator does not exist " + account.Address.String())
	}

	hash := crypto.Keccak256(data)
	return cryptobase.SigAlg.Sign(hash, val.key)
}

func (vm *ValidatorManager) SignDataWithContext(account accounts.Account, mimeType string, data []byte, context []byte) ([]byte, error) {
	val, ok := vm.valMap[account.Address]
	// If the key exists
	if ok == false {
		return nil, errors.New("validator does not exist " + account.Address.String())
	}

	hash := crypto.Keccak256(data)
	return cryptobase.SigAlg.SignWithContext(hash, val.key, context)
}

func (vm *ValidatorManager) GetValidatorsFn(blockHash common.Hash) (map[common.Address]*big.Int, error) {
	valBalanceMap := make(map[common.Address]*big.Int)
	for addr, val := range vm.valMap {
		valBalanceMap[addr] = val.balance
	}

	return valBalanceMap, nil
}

func Initialize(numKeys int) (vm *ValidatorManager, mockp2pManager *MockP2PManager, validatorMap *map[common.Address]*big.Int) {
	STARTUP_DELAY_MS = int64(2000)
	BLOCK_TIMEOUT_MS = int64(6000)
	ACK_BLOCK_TIMEOUT_MS = 18000 //relative to start of block locally
	BLOCK_CLEANUP_TIME_MS = int64(60000)
	MAX_ROUND = byte(2)
	BROADCAST_RESEND_DELAY = int64(100)
	BROADCAST_CLEANUP_DELAY = int64(1800000)
	CONSENSUS_DATA_REQUEST_RESEND_DELAY = int64(60000)
	SKIP_HASH_CHECK = true

	waitMap = make(map[common.Address]bool)
	vm = NewValidatorManager(numKeys)

	valMap, _ := vm.GetValidatorsFn(common.Hash{})

	mockp2pManager = &MockP2PManager{
		mockP2pHandlers:                 make(map[common.Address]*MockP2PHandler),
		blockPacketValidatorMap:         make(map[common.Address]bool),
		blockPacketsBetweenValidatorMap: make(map[common.Hash]bool),
		txnFinalize:                     make(map[common.Hash]bool),
	}

	for addr, _ := range valMap {
		consensusHandler := NewConsensusPacketHandler()
		consensusHandler.getValidatorsFn = vm.GetValidatorsFn
		consensusHandler.doesFinalizedTransactionExistFn = mockp2pManager.DoesFinalizedTransactionExistFn
		account := accounts.Account{
			Address: addr,
		}
		consensusHandler.signFn = vm.SignData
		consensusHandler.signFnWithContext = vm.SignDataWithContext
		consensusHandler.account = account
		p2pHandler := &MockP2PHandler{
			mockP2pManager:   mockp2pManager,
			validator:        addr,
			consensusHandler: consensusHandler,
			validatorDetails: vm.valMap[addr],
			networkDetails:   MockNetworkDetails{},
		}
		consensusHandler.p2pHandler = p2pHandler
		p2pHandler.consensusHandler = consensusHandler
		mockp2pManager.mockP2pHandlers[addr] = p2pHandler
	}

	validatorMap = &valMap

	return
}

func WaitBlockCommit(parentHash common.Hash, mockp2pHandler *MockP2PHandler, t *testing.T) {
	waitLock.Lock()

	valAddress := mockp2pHandler.validator
	_, ok := waitMap[valAddress]
	if ok == true {
		t.Fatalf("validator wait already exists")
	}
	waitLock.Unlock()

	for {
		txns := mockp2pHandler.GetValidatorTransactions()
		err := mockp2pHandler.consensusHandler.HandleConsensus(parentHash, txns, TEST_CONSENSUS_BLOCK_NUMBER)
		if err != nil {
			//fmt.Println("HandleTransactions err", err)
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func ValidateBlockConsensusDataTest(parentHash common.Hash, p2p *MockP2PManager, validatorMap *map[common.Address]*big.Int, t *testing.T) {
	for _, handler := range p2p.mockP2pHandlers {
		blockState, _, err := handler.consensusHandler.getBlockState(parentHash)
		if err != nil {
			fmt.Println("ValidateBlockConsensusData getBlockState", err)
			t.Fatalf("failed")
		}
		if blockState != BLOCK_STATE_RECEIVED_COMMITS {
			continue
		}
		txns, err := handler.consensusHandler.getBlockSelectedTransactions(parentHash)
		if err != nil {
			t.Fatalf("failed")
		}

		blockConsensusData, blockAdditionalConsensusData, err := handler.consensusHandler.getBlockConsensusData(parentHash)
		if err != nil {
			fmt.Println("ValidateBlockConsensusData getBlockConsensusData", "err", err, "val", handler.validator)
			t.Fatalf("failed")
		} else {
			fmt.Println("ValidateBlockConsensusData getBlockConsensusData ok", handler.validator)
		}

		if blockConsensusData == nil || blockAdditionalConsensusData == nil {
			fmt.Println("ValidateBlockConsensusData nil")
			t.Fatalf("failed")
		}

		data, err := rlp.EncodeToBytes(blockAdditionalConsensusData)
		if err != nil {
			fmt.Println("EncodeToBytes", err)
			t.Fatalf("failed")
		}
		//fmt.Println("data len", len(data))
		blockAdditionalConsensusDataDecoded := BlockAdditionalConsensusData{}

		err = rlp.DecodeBytes(data, &blockAdditionalConsensusDataDecoded)
		if err != nil {
			t.Fatalf("failed")
		}
		//fmt.Println("len(blockAdditionalConsensusData.ConsensusPackets)",
		//	len(blockAdditionalConsensusData.ConsensusPackets), len(blockAdditionalConsensusDataDecoded.ConsensusPackets), blockAdditionalConsensusData.InitTime, blockAdditionalConsensusDataDecoded.InitTime)

		if len(blockAdditionalConsensusData.ConsensusPackets) != len(blockAdditionalConsensusDataDecoded.ConsensusPackets) {
			t.Fatalf("failed")
		}

		if blockAdditionalConsensusData.InitTime != blockAdditionalConsensusDataDecoded.InitTime {
			t.Fatalf("failed")
		}

		data2, err := rlp.EncodeToBytes(blockAdditionalConsensusDataDecoded)
		if err != nil {
			t.Fatalf("failed")
		}

		if bytes.Compare(data, data2) != 0 {
			t.Fatalf("failed")
		}

		err = ValidateBlockConsensusDataInner(txns, parentHash, blockConsensusData, blockAdditionalConsensusData, validatorMap, TEST_CONSENSUS_BLOCK_NUMBER, nil)
		if err != nil {
			fmt.Println("ValidateBlockConsensusDataInner", err, handler.validator)
			t.Fatalf("ValidateBlockConsensusDataInner failed")
		}
	}
}

func ValidateTest(validatorMap *map[common.Address]*big.Int, startTime int64, parentHash common.Hash, p2p *MockP2PManager, minPass int, maxWaitCount int,
	expectedVoteMap map[VoteType]bool, expectedState BlockRoundState, t *testing.T) bool {
	i := 0
	j := 0
	count := 0
	for {
		for _, handler := range p2p.mockP2pHandlers {
			blockState, _, err := handler.consensusHandler.getBlockState(parentHash)
			if err != nil {
				return false
			}
			blockVote, err := handler.consensusHandler.getBlockVote(parentHash)
			if err != nil {
				return false
			}
			if blockState == expectedState {
				if expectedVoteMap[blockVote] {
					i = i + 1
				} else {
					j = j + 1
				}
			}
		}
		if i >= minPass {
			PrintState(parentHash, p2p.mockP2pHandlers, startTime)
			ValidateBlockConsensusDataTest(parentHash, p2p, validatorMap, t)
			return true
		}
		if j >= minPass {
			PrintState(parentHash, p2p.mockP2pHandlers, startTime)
			ValidateBlockConsensusDataTest(parentHash, p2p, validatorMap, t)
			return false
		}
		if count == maxWaitCount {
			PrintState(parentHash, p2p.mockP2pHandlers, startTime)
			ValidateBlockConsensusDataTest(parentHash, p2p, validatorMap, t)
			return false
		} else {
			PrintState(parentHash, p2p.mockP2pHandlers, startTime)
			time.Sleep(time.Second * 1)
			i = 0
			count = count + 1
		}
	}
}

func getTestParentHash(blockNumber uint64) common.Hash {
	return common.BytesToHash([]byte(strconv.FormatUint(TEST_CONSENSUS_BLOCK_NUMBER, 10)))
}

func testPacketHandler_basic(numKeys int, t *testing.T) {
	_, p2p, valMap := Initialize(numKeys)

	parentHash := getTestParentHash(TEST_CONSENSUS_BLOCK_NUMBER)

	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		go WaitBlockCommit(parentHash, h, t)
	}

	if ValidateTest(valMap, startTime, parentHash, p2p, numKeys, 10, map[VoteType]bool{VOTE_TYPE_OK: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		txnList, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if err != nil || txnList == nil || len(txnList) != 0 {
			t.Fatalf("failed")
		}
	}
}

func TestPacketHandler_basic(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_basic(4, t)
	}
}

func TestPacketHandler_basic_max_validators(t *testing.T) {
	for i := 1; i <= 1; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_basic(8, t)
	}
}

func testPacketHandler_min_basic(t *testing.T) {
	numKeys := 4
	_, p2p, valMap := Initialize(numKeys)

	parentHash := common.BytesToHash([]byte{1})

	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	proposer, _ := getBlockProposer(parentHash, valMap, 1, nil, TEST_CONSENSUS_BLOCK_NUMBER)

	skipped := false
	c := 0
	skipList := make(map[common.Address]bool)
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if h.validator.IsEqualTo(proposer) == false && skipped == false {
			skipped = true
			skipList[h.validator] = true
			continue
		}
		go WaitBlockCommit(parentHash, h, t)
		c = c + 1
	}

	fmt.Println("c", c)

	if ValidateTest(valMap, startTime, parentHash, p2p, 3, 10, map[VoteType]bool{VOTE_TYPE_OK: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		txnList, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if skipList[h.validator] {
			if err == nil {
				t.Fatalf("failed")
			}
		} else {
			if err != nil || txnList == nil || len(txnList) != 0 {
				t.Fatalf("failed")
			}
		}
	}
}

func TestPacketHandler_min_basic(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_min_basic(t)
	}
}

func testPacketHandler_extended_failure(t *testing.T, numKeys int, minPass int, responsiveValidators int) {
	_, p2p, valMap := Initialize(numKeys)

	parentHash := common.BytesToHash([]byte{1})

	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	proposer, _ := getBlockProposer(parentHash, valMap, 1, nil, TEST_CONSENSUS_BLOCK_NUMBER)

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if h.validator.IsEqualTo(proposer) == true {
			go WaitBlockCommit(parentHash, h, t)
			break
		}
	}

	skipList := make(map[common.Address]bool)
	c := 1
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if h.validator.IsEqualTo(proposer) == true {
			continue
		}
		if c >= responsiveValidators {
			skipList[h.validator] = true
			c = c + 1
			continue
		}
		go WaitBlockCommit(parentHash, h, t)
		c = c + 1
	}

	fmt.Println("c", c)

	if ValidateTest(valMap, startTime, parentHash, p2p, minPass, 10, map[VoteType]bool{VOTE_TYPE_OK: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		txnList, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if skipList[h.validator] {
			if err == nil {
				t.Fatalf("failed")
			}
		} else {
			if err != nil || txnList == nil || len(txnList) != 0 {
				t.Fatalf("failed")
			}
		}
	}
}

func TestPacketHandler_extended_failure(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_extended_failure(t, 4, 3, 3)
		//testPacketHandler_extended_failure(t, 32, 23, 23)
	}
}

func testPacketHandler_block_proposer_timedout(t *testing.T) {
	numKeys := 4
	_, p2p, valMap := Initialize(numKeys)

	parentHash := common.BytesToHash([]byte{1})
	c := 1
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	proposer, _ := getBlockProposer(parentHash, valMap, 1, nil, TEST_CONSENSUS_BLOCK_NUMBER)

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if h.validator.IsEqualTo(proposer) {
			continue //proposer timeout simulation
		}
		go WaitBlockCommit(parentHash, h, t)
		c = c + 1
	}

	if ValidateTest(valMap, startTime, parentHash, p2p, 3, 20, map[VoteType]bool{VOTE_TYPE_NIL: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		txnList, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if h.validator.IsEqualTo(proposer) {
			if err == nil {
				t.Fatalf("failed")
			}
		} else {
			if err != nil || txnList != nil {
				t.Fatalf("failed")
			}
		}
	}
}

func TestPacketHandler_block_proposer_timedout(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_block_proposer_timedout(t)
	}
}

func testPacketHandler_min_negative(t *testing.T, numKeys int, minPass int, unresponsiveValCount int) {
	_, p2p, valMap := Initialize(numKeys)

	parentHash := common.BytesToHash([]byte{1})
	c := 1
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	proposer, _ := getBlockProposer(parentHash, valMap, 1, nil, TEST_CONSENSUS_BLOCK_NUMBER)
	skipList := make(map[common.Address]bool)

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if h.validator.IsEqualTo(proposer) == false && len(skipList) < unresponsiveValCount {
			skipList[h.validator] = true
			continue
		}
		go WaitBlockCommit(parentHash, h, t)
		c = c + 1
	}

	if ValidateTest(valMap, startTime, parentHash, p2p, minPass, 20, map[VoteType]bool{VOTE_TYPE_OK: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == true {
		t.Fatalf("failed")
	}

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		_, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if err == nil {
			t.Fatalf("failed")
		}
	}
}

func TestPacketHandler_min_negative(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_min_negative(t, 4, 3, 2)
	}
}

func testPacketHandler_no_round2_then_round2(t *testing.T, numKeys int, minPass int) {
	_, p2p, valMap := Initialize(numKeys)

	parentHash := common.BytesToHash([]byte{1})
	c := 1
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	proposer, _ := getBlockProposer(parentHash, valMap, 1, nil, TEST_CONSENSUS_BLOCK_NUMBER)
	skipCount := 0
	unresponsiveValCount := 2
	var valSkipList []common.Address
	valSkipList = make([]common.Address, 0)

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if h.validator.IsEqualTo(proposer) == false && skipCount < unresponsiveValCount {
			_, round, err := h.consensusHandler.getBlockState(parentHash)
			if err != nil || round <= 1 {
				skipCount = skipCount + 1
				valSkipList = append(valSkipList, h.validator)
				continue
			}
		}

		numTxns := 1
		txns := make([]common.Hash, numTxns)
		for i := 0; i < numTxns; i++ {
			txns[i] = common.BytesToHash([]byte{byte(rand.Intn(255))})
			h.SetValidatorTransactions(txns)
		}
		h.SetValidatorTransactions(txns)

		go WaitBlockCommit(parentHash, h, t)
		c = c + 1
	}

	breakLoop := false
	checkTime := time.Now()

	fmt.Println("Stage 2")
	for {
		commitCount := 0
		for _, handler := range p2p.mockP2pHandlers {
			h := handler
			state, round, _ := h.consensusHandler.getBlockState(parentHash)
			if round > byte(1) {
				t.Fatalf("failed")
			}
			if state == BLOCK_STATE_RECEIVED_COMMITS {
				commitCount = commitCount + 1
				if commitCount >= minPass {
					breakLoop = true
					break
				}
			}
			if HasExceededTimeThreshold(checkTime, int64(BLOCK_TIMEOUT_MS*2)) {
				for _, v := range valSkipList {
					vh := p2p.mockP2pHandlers[v]

					numTxns := 1
					txns := make([]common.Hash, numTxns)
					for i := 0; i < numTxns; i++ {
						txns[i] = common.BytesToHash([]byte{byte(rand.Intn(255))})
					}
					vh.SetValidatorTransactions(txns)
					go WaitBlockCommit(parentHash, vh, t)
				}
				breakLoop = true
				break
			}
		}
		if breakLoop {
			break
		}
		PrintState(parentHash, p2p.mockP2pHandlers, startTime)
		time.Sleep(time.Second * 1)
	}

	fmt.Println("ValidateTest start")
	if ValidateTest(valMap, startTime, parentHash, p2p, 3, 60, map[VoteType]bool{VOTE_TYPE_NIL: true, VOTE_TYPE_OK: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	hasRound2 := false
	commitCount := 0
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		blockRoundState, round, _ := h.consensusHandler.getBlockState(parentHash)
		if round > byte(1) {
			hasRound2 = true
		}
		if blockRoundState == BLOCK_STATE_RECEIVED_COMMITS {
			commitCount = commitCount + 1
		}
	}
	fmt.Println("commitCount", commitCount)
	if commitCount >= 3 {
		return
	}
	if hasRound2 == false {
		t.Fatalf("failed")
	}
}

func TestPacketHandler_no_round2_then_round2(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_no_round2_then_round2(t, 4, 3)
	}
}

func testPacketHandler_bifurcated(t *testing.T) {
	_, p2p, valMap := Initialize(4)
	parentHash := common.BytesToHash([]byte{1})
	proposer, _ := getBlockProposer(parentHash, valMap, 1, nil, TEST_CONSENSUS_BLOCK_NUMBER)
	c := 0
	startTime := time.Now().UnixNano() / int64(time.Millisecond)

	valList := make([]common.Address, 4)
	j := 0
	var p2pManager *MockP2PManager
	proposerIndex := 0
	for _, handler := range p2p.mockP2pHandlers {
		valList[j].CopyFrom(handler.validator)
		if handler.validator.IsEqualTo(proposer) {
			proposerIndex = j
		}
		j = j + 1
		p2pManager = handler.mockP2pManager
	}
	if proposerIndex != 0 {
		var tmp common.Address
		tmp.CopyFrom(valList[0])
		valList[0].CopyFrom(proposer)
		valList[proposerIndex].CopyFrom(tmp)
	}
	p2pManager.BlockPacketsBetweenValidators(valList[0], valList[2])
	p2pManager.BlockPacketsBetweenValidators(valList[0], valList[3])

	var prev *MockP2PHandler
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if c%2 == 1 {
			h.mockP2pManager.BlockPacketsBetweenValidators(h.validator, prev.validator)
		}

		numTxns := 1
		txns := make([]common.Hash, numTxns)
		for i := 0; i < numTxns; i++ {
			txns[i] = common.BytesToHash([]byte{byte(rand.Intn(255))})
		}
		h.SetValidatorTransactions(txns)
		prev = h
		go WaitBlockCommit(parentHash, h, t)
		c = c + 1
	}

	loopCount := 0
	for {
		commitCount := 0
		for _, handler := range p2p.mockP2pHandlers {
			h := handler
			state, _, _ := h.consensusHandler.getBlockState(parentHash)
			if state == BLOCK_STATE_RECEIVED_COMMITS {
				if loopCount <= 60 {
					t.Fatalf("failed")
				} else {
					commitCount = commitCount + 1
				}
			}
		}
		PrintState(parentHash, p2p.mockP2pHandlers, startTime)
		if commitCount >= 3 {
			return
		}
		loopCount = loopCount + 1
		if loopCount == 60 {
			p2pManager.DeleteAllPacketBlocks()
		} else if loopCount >= 120 {
			t.Fatalf("failed")
		}
		fmt.Println("loopCount", loopCount)
		time.Sleep(time.Second * 1)
	}
}

func TestPacketHandler_bifurcated(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_bifurcated(t)
	}
}

func testPacketHandler_round2(t *testing.T, numKeys int, minPass int) {
	_, p2p, valMap := Initialize(numKeys)
	parentHash := common.BytesToHash([]byte{1})
	c := 0
	startTime := time.Now().UnixNano() / int64(time.Millisecond)

	var prev *MockP2PHandler
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if c%2 == 1 {
			h.mockP2pManager.BlockPacketsBetweenValidators(h.validator, prev.validator)
		}

		numTxns := 1
		txns := make([]common.Hash, numTxns)
		for i := 0; i < numTxns; i++ {
			txns[i] = common.BytesToHash([]byte{byte(rand.Intn(255))})
		}
		h.SetValidatorTransactions(txns)
		prev = h
		go WaitBlockCommit(parentHash, h, t)
		c = c + 1
	}

	breakLoop := false
	loopCount := 0
	for {
		for _, handler := range p2p.mockP2pHandlers {
			h := handler
			_, round, _ := h.consensusHandler.getBlockState(parentHash)
			if round == 2 {
				h.mockP2pManager.DeleteAllPacketBlocks()
				breakLoop = true
				break
			}
		}
		if breakLoop {
			break
		}
		PrintState(parentHash, p2p.mockP2pHandlers, startTime)
		loopCount = loopCount + 1
		if loopCount >= 60 {
			t.Fatalf("failed")
		}
		time.Sleep(time.Second * 1)
	}

	fmt.Println("===============Round 2 start")
	if ValidateTest(valMap, startTime, parentHash, p2p, minPass, 60, map[VoteType]bool{VOTE_TYPE_OK: true, VOTE_TYPE_NIL: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	txnCountOk := 0
	var combinedTxnListHash common.Hash
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		txnList, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if err == nil {
			txnCountOk = txnCountOk + 1
			vote, err := h.consensusHandler.getBlockVote(parentHash)
			if err != nil {
				t.Fatalf("failed")
			}
			if vote == VOTE_TYPE_OK {
				//if txnList == nil || len(txnList) < 1 {
				//	t.Fatalf("failed")
				//}

				_, round, _ := h.consensusHandler.getBlockState(parentHash)
				if err != nil {
					t.Fatalf("failed")
				}

				txnHash := GetCombinedTxnHash(parentHash, round, txnList)
				if combinedTxnListHash.IsEqualTo(ZERO_HASH) {
					combinedTxnListHash.CopyFrom(txnHash)
				} else {
					if combinedTxnListHash.IsEqualTo(txnHash) != true {
						t.Fatalf("failed")
					}
				}
			} else {
				if txnList != nil || len(txnList) > 0 {
					t.Fatalf("failed")
				}
			}
		}
	}
	if txnCountOk < minPass {
		t.Fatalf("failed")
	}
}

func TestPacketHandler_round2(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_round2(t, 4, 3)
	}
}

func testPacketHandler_basic_txns(t *testing.T) {
	numKeys := 4
	_, p2p, valMap := Initialize(numKeys)

	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	parentHash := common.BytesToHash([]byte{1})
	numTxns := 5
	txns := make([]common.Hash, numTxns)
	for i := 0; i < numTxns; i++ {
		txns[i] = common.BytesToHash([]byte{byte(i)})
	}

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		h.SetValidatorTransactions(txns)
		go WaitBlockCommit(parentHash, h, t)
	}

	if ValidateTest(valMap, startTime, parentHash, p2p, numKeys, 10, map[VoteType]bool{VOTE_TYPE_OK: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	txnCountOk := 0
	var combinedTxnListHash common.Hash
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		txnList, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if err == nil {
			txnCountOk = txnCountOk + 1
			if txnList == nil || len(txnList) < 5 {
				t.Fatalf("failed")
			}

			_, round, _ := h.consensusHandler.getBlockState(parentHash)
			if err != nil {
				t.Fatalf("failed")
			}

			txnHash := GetCombinedTxnHash(parentHash, round, txnList)
			if combinedTxnListHash.IsEqualTo(ZERO_HASH) {
				combinedTxnListHash.CopyFrom(txnHash)
			} else {
				if combinedTxnListHash.IsEqualTo(txnHash) != true {
					t.Fatalf("failed")
				}
			}
		}
	}
	if txnCountOk < numKeys {
		t.Fatalf("failed")
	}
}

func TestPacketHandler_basic_txns(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_basic_txns(t)
	}
}

func testPacketHandler_basic_txns_finalize_fail(t *testing.T) {
	numKeys := 4
	_, p2p, valMap := Initialize(numKeys)

	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	parentHash := common.BytesToHash([]byte{1})
	numTxns := 5
	txns := make([]common.Hash, numTxns)
	for i := 0; i < numTxns; i++ {
		txns[i] = common.BytesToHash([]byte{byte(i)})
		if i > 3 {
			p2p.SetTransactionFinalizeState(txns[i], true)
		}
	}

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		h.SetValidatorTransactions(txns)
		go WaitBlockCommit(parentHash, h, t)
	}

	if ValidateTest(valMap, startTime, parentHash, p2p, numKeys-1, 60, map[VoteType]bool{VOTE_TYPE_NIL: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}
}

func TestPacketHandler_basic_txns_finalize_fail(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_basic_txns_finalize_fail(t)
	}
}

func testPacketHandler_split_txns(t *testing.T) {
	numKeys := 6
	_, p2p, valMap := Initialize(numKeys)
	rand.Seed(time.Now().UnixNano())

	parentHash := common.BytesToHash([]byte{1})

	j := 0
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	for _, handler := range p2p.mockP2pHandlers {
		numTxns := 1
		txns := make([]common.Hash, numTxns)
		for i := 0; i < numTxns; i++ {
			txns[i] = common.BytesToHash([]byte{byte(rand.Intn(255))})
			j = j + 1
		}
		handler.SetValidatorTransactions(txns)
	}

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		go WaitBlockCommit(parentHash, h, t)

	}

	if ValidateTest(valMap, startTime, parentHash, p2p, numKeys, 15, map[VoteType]bool{VOTE_TYPE_OK: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	txnCountOk := 0
	var combinedTxnListHash common.Hash
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		txnList, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if err == nil {
			txnCountOk = txnCountOk + 1
			if txnList == nil || len(txnList) < 1 {
				t.Fatalf("failed")
			}

			_, round, _ := h.consensusHandler.getBlockState(parentHash)
			if err != nil {
				t.Fatalf("failed")
			}

			txnHash := GetCombinedTxnHash(parentHash, round, txnList)
			if combinedTxnListHash.IsEqualTo(ZERO_HASH) {
				combinedTxnListHash.CopyFrom(txnHash)
			} else {
				if combinedTxnListHash.IsEqualTo(txnHash) != true {
					t.Fatalf("failed")
				}
			}
		}
	}
	if txnCountOk < numKeys {
		t.Fatalf("failed")
	}
}

func TestPacketHandler_split_txns(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_split_txns(t)
	}
}

func testPacketHandler_split_increasing_txns(t *testing.T) {
	numKeys := 8
	_, p2p, valMap := Initialize(numKeys)

	parentHash := common.BytesToHash([]byte{1})

	j := 0
	numTxns := 0
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	for _, handler := range p2p.mockP2pHandlers {
		numTxns = numTxns + 1
		txns := make([]common.Hash, numTxns)
		for i := 0; i < numTxns; i++ {
			txns[i] = common.BytesToHash([]byte{byte(j)})
			j = j + 1
		}
		h := handler
		handler.SetValidatorTransactions(txns)
		go WaitBlockCommit(parentHash, h, t)
	}

	if ValidateTest(valMap, startTime, parentHash, p2p, 7, 30, map[VoteType]bool{VOTE_TYPE_OK: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	txnCountOk := 0
	var combinedTxnListHash common.Hash
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		txnList, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if err == nil {
			txnCountOk = txnCountOk + 1
			if txnList == nil || len(txnList) < 1 {
				t.Fatalf("failed")
			}

			_, round, _ := h.consensusHandler.getBlockState(parentHash)
			if err != nil {
				t.Fatalf("failed")
			}

			txnHash := GetCombinedTxnHash(parentHash, round, txnList)
			if combinedTxnListHash.IsEqualTo(ZERO_HASH) {
				combinedTxnListHash.CopyFrom(txnHash)
			} else {
				if combinedTxnListHash.IsEqualTo(txnHash) != true {
					t.Fatalf("failed")
				}
			}
		}
	}
	if txnCountOk < 7 {
		t.Fatalf("failed")
	}
}

func TestPacketHandler_split_increasing_txns(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_split_increasing_txns(t)
	}
}

func testPacketHandler_split_increasing_txns_some_unresponsive(t *testing.T, numVal int, minPass int, responsiveValidators int) {
	numKeys := numVal
	_, p2p, valMap := Initialize(numKeys)

	parentHash := common.BytesToHash([]byte{1})
	proposer, _ := getBlockProposer(parentHash, valMap, 1, nil, TEST_CONSENSUS_BLOCK_NUMBER)

	j := 0
	numTxns := 0
	startTime := time.Now().UnixNano() / int64(time.Millisecond)

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if h.validator.IsEqualTo(proposer) == true {
			numTxns = numTxns + 1
			txns := make([]common.Hash, numTxns)
			for i := 0; i < numTxns; i++ {
				txns[i] = common.BytesToHash([]byte{byte(j)})
				j = j + 1
			}
			h.SetValidatorTransactions(txns)
			go WaitBlockCommit(parentHash, h, t)
			break
		}
	}

	c := 1
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if h.validator.IsEqualTo(proposer) == true {
			continue
		}
		if c >= responsiveValidators {
			break
		}
		numTxns = numTxns + 1
		txns := make([]common.Hash, numTxns)
		for i := 0; i < numTxns; i++ {
			txns[i] = common.BytesToHash([]byte{byte(j)})
			j = j + 1
		}
		h.SetValidatorTransactions(txns)
		go WaitBlockCommit(parentHash, h, t)
		c = c + 1
	}

	if ValidateTest(valMap, startTime, parentHash, p2p, minPass, 145, map[VoteType]bool{VOTE_TYPE_OK: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	txnCountOk := 0
	var combinedTxnListHash common.Hash
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		txnList, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if err == nil {
			txnCountOk = txnCountOk + 1
			if txnList == nil || len(txnList) < 1 {
				t.Fatalf("failed")
			}

			_, round, _ := h.consensusHandler.getBlockState(parentHash)
			if err != nil {
				t.Fatalf("failed")
			}

			txnHash := GetCombinedTxnHash(parentHash, round, txnList)
			if combinedTxnListHash.IsEqualTo(ZERO_HASH) {
				combinedTxnListHash.CopyFrom(txnHash)
			} else {
				if combinedTxnListHash.IsEqualTo(txnHash) != true {
					t.Fatalf("failed")
				}
			}
		}
	}
	if txnCountOk < minPass {
		t.Fatalf("failed")
	}
}

func TestPacketHandler_split_increasing_txns_some_unresponsive(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_split_increasing_txns_some_unresponsive(t, 4, 3, 3)
	}
}

func testPacketHandler_packet_loss_txns(t *testing.T) {
	numKeys := 10
	minPass := 7
	_, p2p, valMap := Initialize(numKeys)

	parentHash := common.BytesToHash([]byte{1})

	j := 0
	numTxns := 0
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	for _, handler := range p2p.mockP2pHandlers {
		numTxns = numTxns + 1
		txns := make([]common.Hash, numTxns)
		for i := 0; i < numTxns; i++ {
			txns[i] = common.BytesToHash([]byte{byte(j)})
			j = j + 1
		}
		h := handler
		handler.SetValidatorTransactions(txns)
		h.networkDetails.packetLoss = 50
		go WaitBlockCommit(parentHash, h, t)
	}

	if ValidateTest(valMap, startTime, parentHash, p2p, minPass, 600, map[VoteType]bool{VOTE_TYPE_OK: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	txnCountOk := 0
	var combinedTxnListHash common.Hash
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		txnList, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if err == nil {
			txnCountOk = txnCountOk + 1
			if txnList == nil || len(txnList) < 1 {
				continue
			}

			_, round, _ := h.consensusHandler.getBlockState(parentHash)
			if err != nil {
				t.Fatalf("failed")
			}

			txnHash := GetCombinedTxnHash(parentHash, round, txnList)
			if combinedTxnListHash.IsEqualTo(ZERO_HASH) {
				combinedTxnListHash.CopyFrom(txnHash)
			} else {
				if combinedTxnListHash.IsEqualTo(txnHash) != true {
					t.Fatalf("failed")
				}
			}
		}
	}
	if txnCountOk < minPass {
		t.Fatalf("failed")
	}
}

func TestPacketHandler_packet_loss_txns(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_packet_loss_txns(t)
	}
}

func testPacketHandler_packet_loss_txns_some_unresponsive(t *testing.T, numVal int, minPass int, responsiveValidators int) {
	numKeys := numVal
	_, p2p, valMap := Initialize(numKeys)

	parentHash := common.BytesToHash([]byte{1})
	proposer, _ := getBlockProposer(parentHash, valMap, 1, nil, TEST_CONSENSUS_BLOCK_NUMBER)

	j := 0
	numTxns := 0
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	//fmt.Println("precommit", crypto.Keccak256Hash(ZERO_HASH.Bytes(), []byte{byte(VOTE_TYPE_NIL)}))

	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if h.validator.IsEqualTo(proposer) == true {
			numTxns = numTxns + 1
			txns := make([]common.Hash, numTxns)
			for i := 0; i < numTxns; i++ {
				txns[i] = common.BytesToHash([]byte{byte(j)})
				j = j + 1
			}
			h.SetValidatorTransactions(txns)
			h.networkDetails.packetLoss = 50
			go WaitBlockCommit(parentHash, h, t)
			break
		}
	}

	c := 1
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		if h.validator.IsEqualTo(proposer) == true {
			continue
		}
		if c >= responsiveValidators {
			break
		}
		numTxns = numTxns + 1
		txns := make([]common.Hash, numTxns)
		for i := 0; i < numTxns; i++ {
			txns[i] = common.BytesToHash([]byte{byte(j)})
			j = j + 1
		}
		handler.SetValidatorTransactions(txns)
		go WaitBlockCommit(parentHash, h, t)
		c = c + 1
	}

	fmt.Println("c", c)

	if ValidateTest(valMap, startTime, parentHash, p2p, minPass, 300, map[VoteType]bool{VOTE_TYPE_OK: true, VOTE_TYPE_NIL: true}, BLOCK_STATE_RECEIVED_COMMITS, t) == false {
		t.Fatalf("failed")
	}

	txnCountOk := 0
	nilVoteCount := 0
	var combinedTxnListHash common.Hash
	for _, handler := range p2p.mockP2pHandlers {
		h := handler
		state, round, err := h.consensusHandler.getBlockState(parentHash)
		if err != nil {
			t.Fatalf("failed")
		}
		if state != BLOCK_STATE_RECEIVED_COMMITS {
			continue
		}

		vote, err := h.consensusHandler.getBlockVote(parentHash)
		if err != nil {
			t.Fatalf("failed")
		}

		if vote == VOTE_TYPE_NIL {
			nilVoteCount = nilVoteCount + 1
			continue
		}

		txnList, err := h.consensusHandler.getBlockSelectedTransactions(parentHash)
		if err == nil {
			txnCountOk = txnCountOk + 1
			if txnList == nil {
				t.Fatalf("failed")
			}
			if len(txnList) < 1 {
				return
			}

			txnHash := GetCombinedTxnHash(parentHash, round, txnList)
			if combinedTxnListHash.IsEqualTo(ZERO_HASH) {
				combinedTxnListHash.CopyFrom(txnHash)
			} else {
				if combinedTxnListHash.IsEqualTo(txnHash) != true {
					t.Fatalf("failed")
				}
			}
		}
	}
	if txnCountOk < minPass && nilVoteCount < minPass {
		t.Fatalf("failed")
	}
}

func TestPacketHandler_packet_loss_txns_some_unresponsive(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_packet_loss_txns_some_unresponsive(t, 10, 7, 7)
	}
}

func PrintState(parentHash common.Hash, mockP2pHandlers map[common.Address]*MockP2PHandler, startTime int64) {
	var roundMap map[byte]int32
	var roundStateMap map[byte]map[BlockRoundState]int32
	var roundVoteMap map[byte]map[VoteType]int32

	roundStateMap = make(map[byte]map[BlockRoundState]int32)
	roundVoteMap = make(map[byte]map[VoteType]int32)

	roundMap = make(map[byte]int32)

	for _, handler := range mockP2pHandlers {
		state, round, _ := handler.consensusHandler.getBlockState(parentHash)
		vote, _ := handler.consensusHandler.getBlockVote(parentHash)
		fmt.Println("   validator", handler.validator, "vote", vote, "round", round, "state", state)
		var i byte
		for i = 1; i <= round; i++ {
			roundState, voteType, voteCount, _ := handler.consensusHandler.getBlockRoundState(parentHash, i)
			fmt.Println("       ", "round", i, "roundState", roundState, "VoteType", voteType, "voteCount", voteCount)
		}

		roundCount, ok := roundMap[round]
		if ok == false {
			roundMap[round] = 1
			roundStateMap[round] = make(map[BlockRoundState]int32)
			roundVoteMap[round] = make(map[VoteType]int32)

			roundStateMap[round][state] = 1
			roundVoteMap[round][vote] = 1
		} else {
			roundMap[round] = roundCount + 1
			roundStateMap[round][state] = roundStateMap[round][state] + 1
			roundVoteMap[round][vote] = roundVoteMap[round][vote] + 1
		}
	}

	for r, i := range roundMap {
		fmt.Println("round", r, "count", i)

		for s, j := range roundStateMap[r] {
			fmt.Println("     state", s, "count", j)
		}

		for v, k := range roundVoteMap[r] {
			fmt.Println("     VoteType", v, "count", k)
		}

		if r > 1 {
			for _, handler := range mockP2pHandlers {
				roundDetails, _ := handler.consensusHandler.getBlockRound(parentHash, r)
				if roundDetails != nil {
					fmt.Println("     round reason", "val", handler.consensusHandler.account.Address, "reason", roundDetails.newRoundReason)
				}
			}
		}
	}

	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	timeTaken := endTime - startTime
	fmt.Println("timeTaken ms", timeTaken)
}

func TestPacketHandler_packet_loss_txns_some_unresponsive_extended(t *testing.T) {
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_packet_loss_txns_some_unresponsive(t, 10, 7, 7)
	}
}

func testFilterValidatorsTest(t *testing.T, parentHash common.Hash, validatorsDepositMap map[common.Address]*big.Int, shouldPass bool) *big.Int {
	resultMap, filteredDepositValue, _, err := filterValidators(parentHash, &validatorsDepositMap)
	if err == nil {
		if shouldPass == false {
			t.Fatalf("failed")
		}
	} else {
		fmt.Println("filterValidators error", err)
		if shouldPass == true {
			t.Fatalf("filterValidators failed")
		}
		return nil
	}

	if MIN_BLOCK_DEPOSIT.Cmp(filteredDepositValue) > 0 {
		t.Fatalf("failed")
	}

	fmt.Println("selected validator count", len(resultMap), "total validators", len(validatorsDepositMap))
	if len(resultMap) < MIN_VALIDATORS {
		t.Fatalf("failed")
	}

	if len(resultMap) > MAX_VALIDATORS {
		t.Fatalf("failed")
	}

	if len(validatorsDepositMap) <= MAX_VALIDATORS && len(resultMap) != len(validatorsDepositMap) {
		t.Fatalf("failed")
	}

	if len(validatorsDepositMap) > MAX_VALIDATORS && len(resultMap) != MAX_VALIDATORS {
		t.Fatalf("failed")
	}

	totalDeposit := big.NewInt(0)
	for val, include := range resultMap {
		if include == false {
			t.Fatalf("failed")
		}
		depositValue, ok := validatorsDepositMap[val]
		if ok == false {
			t.Fatalf("unexpected validator")
		}

		totalDeposit = common.SafeAddBigInt(totalDeposit, depositValue)
		fmt.Println("Selected", "validator", val, "deposit", depositValue)
	}
	fmt.Println("filteredDepositValue", filteredDepositValue, "totalDeposit", totalDeposit)

	if totalDeposit.Cmp(filteredDepositValue) > 0 {
		t.Fatalf("failed")
	}

	if MIN_BLOCK_DEPOSIT.Cmp(totalDeposit) > 0 {
		t.Fatalf("failed")
	}

	return filteredDepositValue
}

func TestFilterValidators_negative(t *testing.T) {
	parentHash := common.BytesToHash([]byte{100})
	validatorsDepositMap := make(map[common.Address]*big.Int)
	testFilterValidatorsTest(t, parentHash, validatorsDepositMap, false)

	val1 := common.BytesToAddress([]byte{1})
	val2 := common.BytesToAddress([]byte{2})
	val3 := common.BytesToAddress([]byte{3})

	validatorsDepositMap[val1] = big.NewInt(1000000)
	validatorsDepositMap[val2] = big.NewInt(2000000)
	testFilterValidatorsTest(t, parentHash, validatorsDepositMap, false)

	validatorsDepositMap[val1] = big.NewInt(10000)
	validatorsDepositMap[val2] = big.NewInt(20000)
	validatorsDepositMap[val3] = big.NewInt(30000)
	testFilterValidatorsTest(t, parentHash, validatorsDepositMap, false)

	b := byte(0)
	for i := 0; i < MAX_VALIDATORS*2; i++ {
		val := common.BytesToAddress([]byte{b})
		validatorsDepositMap[val] = big.NewInt(1000)
		b = b + 1
	}
	testFilterValidatorsTest(t, parentHash, validatorsDepositMap, false)
}

func TestFilterValidators_positive(t *testing.T) {
	parentHash := common.BytesToHash([]byte{100})
	validatorsDepositMap := make(map[common.Address]*big.Int)

	val1 := common.BytesToAddress([]byte{1})
	val2 := common.BytesToAddress([]byte{2})
	val3 := common.BytesToAddress([]byte{3})

	validatorsDepositMap[val1] = params.EtherToWei(big.NewInt(100000000000))
	validatorsDepositMap[val2] = params.EtherToWei(big.NewInt(200000000000))
	validatorsDepositMap[val3] = params.EtherToWei(big.NewInt(400000000000))
	fmt.Println("Test1")
	testFilterValidatorsTest(t, parentHash, validatorsDepositMap, true)

	b := byte(0)
	for i := 0; i < MAX_VALIDATORS/2; i++ {
		val := common.BytesToAddress([]byte{b})
		validatorsDepositMap[val] = params.EtherToWei(big.NewInt(10000000000))
		b = b + 1
	}
	fmt.Println("Test2")
	testFilterValidatorsTest(t, parentHash, validatorsDepositMap, true)

	b = byte(0)
	for i := 0; i < MAX_VALIDATORS; i++ {
		val := common.BytesToAddress([]byte{b})
		validatorsDepositMap[val] = params.EtherToWei(big.NewInt(5000000000))
		b = b + 1
	}
	fmt.Println("Test3")
	testFilterValidatorsTest(t, parentHash, validatorsDepositMap, true)

	b = byte(0)
	for i := 0; i < MAX_VALIDATORS+1; i++ {
		val := common.BytesToAddress([]byte{b})
		validatorsDepositMap[val] = params.EtherToWei(big.NewInt(5000000000))
		b = b + 1
	}
	fmt.Println("Test4")
	testFilterValidatorsTest(t, parentHash, validatorsDepositMap, true)
}

func TestFilterValidators_positive_Extended(t *testing.T) {
	parentHash := common.BytesToHash([]byte{100})
	validatorsDepositMap := make(map[common.Address]*big.Int)

	b := byte(0)
	for i := 0; i < MAX_VALIDATORS+1; i++ {
		val := common.BytesToAddress([]byte{b})
		validatorsDepositMap[val] = params.EtherToWei(big.NewInt(5000000000))
		b = b + 1
	}
	testFilterValidatorsTest(t, parentHash, validatorsDepositMap, true)
}

func TestFilterValidators_positive_Tough(t *testing.T) {
	for test := 0; test < 2; test++ {
		validatorsDepositMap := make(map[common.Address]*big.Int)

		b := byte(0)
		for i := 1; i < 255; i++ {
			val := common.BytesToAddress([]byte{b})
			validatorsDepositMap[val] = params.EtherToWei(big.NewInt(1000000000))
			b = b + 1
		}

		for i := 1; i < 255; i++ {
			val := common.BytesToAddress([]byte{b})
			validatorsDepositMap[val] = params.EtherToWei(big.NewInt(20000000000))
			b = b + 1
		}

		parentHash1 := common.BytesToHash([]byte{100})
		totalDeposit := testFilterValidatorsTest(t, parentHash1, validatorsDepositMap, true)
		expected := params.EtherToWei(big.NewInt(2522000000000))
		if totalDeposit.Cmp(expected) != 0 {
			fmt.Println("dep", params.WeiToEther(totalDeposit))
			t.Fatalf("failed a")
		}

		parentHash2 := common.BytesToHash([]byte{200})
		totalDeposit = testFilterValidatorsTest(t, parentHash2, validatorsDepositMap, true)
		if totalDeposit.Cmp(params.EtherToWei(big.NewInt(2522000000000))) != 0 {
			fmt.Println("dep", params.WeiToEther(totalDeposit))
			t.Fatalf("failed b")
		}

		parentHash3 := common.BytesToHash([]byte{255})
		totalDeposit = testFilterValidatorsTest(t, parentHash3, validatorsDepositMap, true)
		if totalDeposit.Cmp(params.EtherToWei(big.NewInt(2522000000000))) != 0 {
			fmt.Println("dep", params.WeiToEther(totalDeposit))
			t.Fatalf("failed c")
		}
	}
}

func TestFilterValidators_positive_low_balance(t *testing.T) {
	for test := 0; test < 2; test++ {
		validatorsDepositMap := make(map[common.Address]*big.Int)

		val1 := common.BytesToAddress([]byte{1})
		validatorsDepositMap[val1] = params.EtherToWei(big.NewInt(1000))

		val2 := common.BytesToAddress([]byte{2})
		validatorsDepositMap[val2] = params.EtherToWei(big.NewInt(900000000000))

		val3 := common.BytesToAddress([]byte{3})
		validatorsDepositMap[val3] = params.EtherToWei(big.NewInt(10000000))

		val4 := common.BytesToAddress([]byte{4})
		validatorsDepositMap[val4] = params.EtherToWei(big.NewInt(5000000))

		parentHash1 := common.BytesToHash([]byte{100})
		totalDeposit := testFilterValidatorsTest(t, parentHash1, validatorsDepositMap, true)
		if totalDeposit.Cmp(params.EtherToWei(big.NewInt(900015000000))) != 0 {
			fmt.Println("dep", params.WeiToEther(totalDeposit))
			t.Fatalf("failed")
		}
	}
}

func TestFilterValidators_positive_low_balance_negative_total(t *testing.T) {
	for test := 0; test < 2; test++ {
		validatorsDepositMap := make(map[common.Address]*big.Int)

		val1 := common.BytesToAddress([]byte{1})
		validatorsDepositMap[val1] = big.NewInt(1000)

		val2 := common.BytesToAddress([]byte{2})
		validatorsDepositMap[val2] = big.NewInt(100000)

		val3 := common.BytesToAddress([]byte{3})
		validatorsDepositMap[val3] = big.NewInt(200000)

		val4 := common.BytesToAddress([]byte{4})
		validatorsDepositMap[val4] = big.NewInt(300000)

		parentHash1 := common.BytesToHash([]byte{100})
		testFilterValidatorsTest(t, parentHash1, validatorsDepositMap, false)
	}
}

func TestFilterValidators_positive_low_balance_negative(t *testing.T) {
	for test := 0; test < 2; test++ {
		validatorsDepositMap := make(map[common.Address]*big.Int)

		b := byte(0)
		for i := 1; i < 255; i++ {
			val := common.BytesToAddress([]byte{b})
			validatorsDepositMap[val] = big.NewInt(1000)
			b = b + 1
		}

		val2 := common.BytesToAddress([]byte{1, 2})
		validatorsDepositMap[val2] = big.NewInt(100000)

		val3 := common.BytesToAddress([]byte{1, 3})
		validatorsDepositMap[val3] = big.NewInt(1000000)

		parentHash1 := common.BytesToHash([]byte{100})
		testFilterValidatorsTest(t, parentHash1, validatorsDepositMap, false)
	}
}

func TestBlockProposalTime(t *testing.T) {
	for i := uint64(0); i < 1000000000; i += 256 {
		if GetProposalTime(i) == 0 {
			fmt.Println(i)
			t.Fatalf("failed 1")
		}
	}

	t1 := GetProposalTime(256)
	tm := time.Unix(int64(t1), 0)
	fmt.Println(tm)

	if tm.Second() != 0 || tm.Nanosecond() != 0 {
		t.Fatalf("failed 2")
	}

	val := tm.Unix()
	if val%60 != 0 {
		t.Fatalf("failed 3")
	}

	if GetProposalTime(1) == 0 {
		t.Fatalf("failed 4")
	}
}

func TestValidateBlockProposalTime(t *testing.T) {
	if ValidateBlockProposalTime(1, GetProposalTime(1)) == false {
		t.Fatalf("failed 1")
	}

	if ValidateBlockProposalTime(256, GetProposalTime(256)) == false {
		t.Fatalf("failed 2")
	}

	if ValidateBlockProposalTime(2, GetProposalTime(2)) == false {
		t.Fatalf("failed 3")
	}

	if ValidateBlockProposalTime(1, GetProposalTime(1)+1) == true {
		t.Fatalf("failed 4")
	}
}

func TestValidateBlockProposalTimeConsensus(t *testing.T) {
	if ValidateBlockProposalTimeConsensus(1, GetProposalTime(1)) == false {
		t.Fatalf("failed 1")
	}

	if ValidateBlockProposalTimeConsensus(256, GetProposalTime(256)) == false {
		t.Fatalf("failed 2")
	}

	if ValidateBlockProposalTimeConsensus(2, GetProposalTime(2)) == false {
		t.Fatalf("failed 3")
	}

	if ValidateBlockProposalTimeConsensus(1, GetProposalTime(2)) == true {
		t.Fatalf("failed 4")
	}

	if ValidateBlockProposalTimeConsensus(1, GetProposalTime(1)+1) == true {
		t.Fatalf("failed 5")
	}

	tm := time.Now().UTC().Add(time.Minute * time.Duration(1)).Unix()
	if tm%60 != 0 {
		tm = tm - (tm % 60)
	}
	if ValidateBlockProposalTimeConsensus(1, uint64(tm)) == false {
		t.Fatalf("failed 6")
	}

	tm = time.Now().UTC().Add(time.Minute * time.Duration(2)).Unix()
	if tm%60 != 0 {
		tm = tm - (tm % 60)
	}
	if ValidateBlockProposalTimeConsensus(1, uint64(tm)) == false {
		t.Fatalf("failed 7")
	}

	tm = time.Now().UTC().Add(time.Minute * time.Duration(3)).Unix()
	if tm%60 != 0 {
		tm = tm - (tm % 60)
	}
	if ValidateBlockProposalTimeConsensus(1, uint64(tm)) == false {
		t.Fatalf("failed 8")
	}

	tm = time.Now().UTC().Add(time.Minute * time.Duration(4)).Unix()
	if tm%60 != 0 {
		tm = tm - (tm % 60)
	}
	if ValidateBlockProposalTimeConsensus(1, uint64(tm)) == true {
		t.Fatalf("failed 9")
	}

	tm = time.Now().UTC().Add(time.Minute * time.Duration(-1)).Unix()
	if tm%60 != 0 {
		tm = tm - (tm % 60)
	}
	if ValidateBlockProposalTimeConsensus(1, uint64(tm)) == false {
		t.Fatalf("failed 10")
	}

	tm = time.Now().UTC().Add(time.Minute * time.Duration(-2)).Unix()
	if tm%60 != 0 {
		tm = tm - (tm % 60)
	}
	if ValidateBlockProposalTimeConsensus(1, uint64(tm)) == false {
		t.Fatalf("failed 11")
	}

	tm = time.Now().UTC().Add(time.Minute * time.Duration(-3)).Unix()
	if tm%60 != 0 {
		tm = tm - (tm % 60)
	}
	if ValidateBlockProposalTimeConsensus(1, uint64(tm)) == false {
		t.Fatalf("failed 12")
	}

	tm = time.Now().UTC().Add(time.Minute * time.Duration(-4)).Unix()
	if tm%60 != 0 {
		tm = tm - (tm % 60)
	}
	if ValidateBlockProposalTimeConsensus(1, uint64(tm)) == true {
		t.Fatalf("failed 13")
	}
}

func Test_consensuspacket_negative(t *testing.T) {
	numKeys := 4
	_, p2p, _ := Initialize(numKeys)

	for _, handler := range p2p.mockP2pHandlers {
		h := handler

		var p0 eth.ConsensusPacket
		err := h.consensusHandler.HandleConsensusPacket(&p0)
		if err == nil {
			t.Fatalf("failed1")
		}

		p1 := eth.ConsensusPacket{}
		err = h.consensusHandler.HandleConsensusPacket(&p1)

		if err == nil {
			t.Fatalf("failed2")
		}

		p2 := eth.ConsensusPacket{
			Signature: make([]byte, 10),
		}

		err = h.consensusHandler.HandleConsensusPacket(&p2)

		if err == nil {
			t.Fatalf("failed3")
		}

		p3 := eth.ConsensusPacket{
			ConsensusData: make([]byte, 10),
		}

		err = h.consensusHandler.HandleConsensusPacket(&p3)

		if err == nil {
			t.Fatalf("failed4")
		}

		p4 := eth.ConsensusPacket{
			Signature:     make([]byte, 10),
			ConsensusData: make([]byte, 0),
		}

		err = h.consensusHandler.HandleConsensusPacket(&p4)

		if err == nil {
			t.Fatalf("failed5")
		}

		p5 := eth.ConsensusPacket{
			Signature:     make([]byte, 0),
			ConsensusData: make([]byte, 10),
		}

		err = h.consensusHandler.HandleConsensusPacket(&p5)

		if err == nil {
			t.Fatalf("failed6")
		}
	}

}

func Test_requestconsensuspacket_negative(t *testing.T) {
	numKeys := 4
	_, p2p, _ := Initialize(numKeys)

	for _, handler := range p2p.mockP2pHandlers {
		h := handler

		var p0 eth.RequestConsensusDataPacket
		_, err := h.consensusHandler.HandleRequestConsensusDataPacket(&p0)
		if err == nil {
			t.Fatalf("failed1")
		}

		p1 := eth.RequestConsensusDataPacket{}
		_, err = h.consensusHandler.HandleRequestConsensusDataPacket(&p1)

		if err == nil {
			t.Fatalf("failed2")
		}

		p2 := eth.RequestConsensusDataPacket{
			RequestData: make([]byte, 0),
		}

		_, err = h.consensusHandler.HandleRequestConsensusDataPacket(&p2)

		if err == nil {
			t.Fatalf("failed3")
		}
	}
}

func Test_shouldSignFull(t *testing.T) {
	for i := uint64(0); i < FULL_SIGN_PROPOSAL_CUTOFF_BLOCK; i++ {
		if shouldSignFull(uint64(i)) == true {
			t.Fatalf("failed 1")
		}
	}

	for i := FULL_SIGN_PROPOSAL_CUTOFF_BLOCK; i < FULL_SIGN_PROPOSAL_CUTOFF_BLOCK*100; i += FULL_SIGN_PROPOSAL_FREQUENCY_BLOCKS {
		if shouldSignFull(uint64(i)) == false {
			t.Fatalf("failed 2")
		}
	}

	for i := FULL_SIGN_PROPOSAL_CUTOFF_BLOCK + 1; i < FULL_SIGN_PROPOSAL_CUTOFF_BLOCK+FULL_SIGN_PROPOSAL_FREQUENCY_BLOCKS-1; i++ {
		if shouldSignFull(uint64(i)) == true {
			t.Fatalf("failed 3")
		}
	}
}

func TestPacketHandler_basic_fullsign(t *testing.T) {
	fmt.Println("TestPacketHandler_basic_fullsign starting")
	TEST_CONSENSUS_BLOCK_NUMBER = FULL_SIGN_PROPOSAL_CUTOFF_BLOCK
	for i := 1; i <= TEST_ITERATIONS; i++ {
		fmt.Println("iteration", i)
		testPacketHandler_basic(4, t)
	}
	TEST_CONSENSUS_BLOCK_NUMBER = uint64(1)
	fmt.Println("TestPacketHandler_basic_fullsign done")
}

func canProposeTest(lastNilBlock int64, nilBlockCount int64, currentBlock uint64, expected bool) bool {
	valDetails := &ValidatorDetailsV2{
		LastNiLBlock:  big.NewInt(lastNilBlock),
		NilBlockCount: big.NewInt(nilBlockCount),
	}

	result := canPropose(valDetails, currentBlock)
	if result != expected {
		return false
	}

	return true
}

func TestPacketHandler_canPropose(t *testing.T) {
	if canProposeTest(0, 0, 100, true) == false {
		t.Fatalf("failed")
	}
	if canProposeTest(0, 10, 100, true) == false {
		t.Fatalf("failed")
	}
	if canProposeTest(1, 1, 2, true) == false {
		t.Fatalf("failed")
	}
	if canProposeTest(1, 1, 3, true) == false {
		t.Fatalf("failed")
	}

	if canProposeTest(50, 1, 51, true) == false {
		t.Fatalf("failed")
	}

	for i := uint64(1); i < 16; i++ {
		if canProposeTest(50, int64(i*BLOCK_PROPOSER_OFFLINE_NIL_BLOCK_MULTIPLIER), 51, true) == true {
			t.Fatalf("failed")
		}
	}
}

func testGetBlockProposerV2(validatorMap *map[common.Address]*ValidatorDetailsV2, expected common.Address, blockNumber uint64) bool {
	parentHash := common.BytesToHash([]byte(strconv.FormatInt(int64(blockNumber), 10)))
	proposer, err := getBlockProposerV2(parentHash, validatorMap, 1, blockNumber)
	if err != nil {
		fmt.Println("err", err)
		return false
	}

	fmt.Println("proposer", proposer, "expected", expected)

	return proposer.IsEqualTo(expected)
}

func TestPacketHandler_getBlockProposerV2(t *testing.T) {
	validatorMap := make(map[common.Address]*ValidatorDetailsV2)

	for i := 0; i < 100; i++ {
		v := &ValidatorDetailsV2{
			Validator:     common.BytesToAddress([]byte(string(rune(i)))),
			LastNiLBlock:  new(big.Int),
			NilBlockCount: new(big.Int),
		}
		validatorMap[v.Validator] = v
	}

	for i := 101; i < 128; i++ {
		v := &ValidatorDetailsV2{
			Validator:     common.BytesToAddress([]byte(string(rune(i)))),
			LastNiLBlock:  big.NewInt(50),
			NilBlockCount: big.NewInt(10),
		}
		validatorMap[v.Validator] = v
	}

	if testGetBlockProposerV2(&validatorMap, common.HexToAddress("0x0000000000000000000000000000000000000000000000000000000000000059"), 81) == false {
		t.Fatalf("failed")
	}

	if testGetBlockProposerV2(&validatorMap, common.HexToAddress("0x0000000000000000000000000000000000000000000000000000000000000056"), 85) == false {
		t.Fatalf("failed")
	}

	if testGetBlockProposerV2(&validatorMap, common.HexToAddress("0x000000000000000000000000000000000000000000000000000000000000005a"), 50) == false {
		t.Fatalf("failed")
	}

	validatorMap = make(map[common.Address]*ValidatorDetailsV2)
	for i := 0; i < MIN_VALIDATORS; i++ {
		if i == 0 {
			v := &ValidatorDetailsV2{
				Validator:     common.BytesToAddress([]byte(string(rune(i)))),
				LastNiLBlock:  big.NewInt(20),
				NilBlockCount: big.NewInt(100),
			}
			validatorMap[v.Validator] = v
		} else {
			v := &ValidatorDetailsV2{
				Validator:     common.BytesToAddress([]byte(string(rune(i)))),
				LastNiLBlock:  new(big.Int),
				NilBlockCount: new(big.Int),
			}
			validatorMap[v.Validator] = v
		}
	}

	if testGetBlockProposerV2(&validatorMap, common.HexToAddress("0x0000000000000000000000000000000000000000000000000000000000000001"), 50) == false {
		t.Fatalf("failed")
	}

	validatorMap = make(map[common.Address]*ValidatorDetailsV2)
	if testGetBlockProposerV2(&validatorMap, common.HexToAddress("0x0000000000000000000000000000000000000000000000000000000000000001"), 50) == true {
		t.Fatalf("failed")
	}
}
