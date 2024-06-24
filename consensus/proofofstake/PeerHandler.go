package proofofstake

import (
	"errors"
	"github.com/DogeProtocol/dp/accounts"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/eth/protocols/eth"
	"github.com/DogeProtocol/dp/handler"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/rlp"
	"runtime/debug"
	"sync"
)

const MinConsensusNetworkProtocolVersion = byte(5)
const ConsensusNetworkProtocolVersion = byte(5)

type GetLatestBlockNumberFn func() uint64

var isConsensusRelay = true

type PeerDetails struct {
	capabilityDetails *CapabilityDetails
	peerId            string
}

type PacketSyncDetails struct {
	incomingPeerMap map[string]bool //List of peers who sent this packet
	packet          *eth.ConsensusPacket
	sendPeerMap     map[string]bool //List of peers to which this packet was sent
}

type PeerHandler struct {
	peerMap                map[string]*PeerDetails //Superset of connected peers
	peerLock               sync.Mutex
	p2pHandler             *handler.P2PHandler
	signFn                 SignerFn
	account                accounts.Account
	isConsensusRelay       bool
	getLatestBlockNumberFn GetLatestBlockNumberFn
	localPeerId            string
	consensusRelayMap      map[string]bool                    //List of connected ConsensusRelays
	syncPeerMap            map[string]bool                    //List of peers who have requested for consensus sync (i.e. ConsensusRelaying consensus packets)
	packetSyncMap          map[common.Hash]*PacketSyncDetails //packet hash is the key

	parentHashLock     sync.Mutex
	currentParentHash  common.Hash
	currentBlockNumber uint64
}

// Sent by a ConsensusRelay to another node
type CapabilityDetails struct {
	IsConsensusRelay bool   `json:"IsConsensusRelay" gencodec:"required"` //should always be true
	PeerId           string `json:"PeerId" gencodec:"required"`           //PeerId of the original sender
}

// Send by a node to a ConsensusRelay, to request consensus packets
type RequestConsensusSyncDetails struct {
	IsConsensusRelay bool   `json:"IsConsensusRelay" gencodec:"required"` //Whether requester is also a ConsensusRelay
	PeerId           string `json:"PeerId" gencodec:"required"`           //PeerId of the original sender (requester)
}

func NewPeerHandler(isConsensusRelay bool, getLatestBlockNumberFn GetLatestBlockNumberFn) *PeerHandler {
	if isConsensusRelay {
		log.Trace("NewPeerHandler isConsensusRelay")
	}
	return &PeerHandler{
		isConsensusRelay:       isConsensusRelay,
		getLatestBlockNumberFn: getLatestBlockNumberFn,
		peerMap:                make(map[string]*PeerDetails),
		consensusRelayMap:      make(map[string]bool),
		syncPeerMap:            make(map[string]bool),
		packetSyncMap:          make(map[common.Hash]*PacketSyncDetails),
	}
}

func (p *PeerHandler) SetP2PHandler(handler *handler.P2PHandler, localPeerId string) {
	p.p2pHandler = handler
	p.localPeerId = localPeerId
}

func (p *PeerHandler) SetSignFn(signFn SignerFn, account accounts.Account) {
	p.signFn = signFn
	p.account = account
}

func (p *PeerHandler) OnPeerConnected(peerId string) error {
	log.Debug("OnPeerConnected start", "peerId", peerId)
	p.peerLock.Lock()
	defer p.peerLock.Unlock()

	p.peerMap[peerId] = &PeerDetails{
		peerId: peerId,
	}

	if p.isConsensusRelay {
		go p.SendCapabilityPacket([]string{peerId})
	}

	log.Debug("OnPeerConnected done", "peerId", peerId)
	return nil
}

func (p *PeerHandler) OnPeerDisconnected(peerId string) error {
	log.Debug("OnPeerDisconnected start", "peerId", peerId)
	p.peerLock.Lock()
	defer p.peerLock.Unlock()

	delete(p.peerMap, peerId)
	delete(p.consensusRelayMap, peerId)
	delete(p.syncPeerMap, peerId)

	if len(p.consensusRelayMap) == 0 {
		go p.ConnectAvailableConsensusRelay()
	}

	log.Debug("OnPeerDisconnected done", "peerId", peerId)
	return nil
}

func (p *PeerHandler) HandleConsensusPacket(packet *eth.ConsensusPacket, fromPeerId string) error {
	log.Trace("PeerHandler HandleConsensusPacket", "fromPeerId", fromPeerId)
	if packet == nil || packet.Signature == nil || packet.ConsensusData == nil || len(packet.Signature) == 0 || len(packet.ConsensusData) == 0 {
		log.Debug("HandleConsensusPacket nil", "fromPeerId", fromPeerId)
		return InvalidPacketErr
	}

	var startIndex int
	if packet.ConsensusData[0] >= MinConsensusNetworkProtocolVersion {
		startIndex = 2
	} else {
		startIndex = 1
	}

	packetType := ConsensusPacketType(packet.ConsensusData[startIndex-1])

	if packetType == CONSENSUS_PACKET_TYPE_CAPABILITY {
		capabilityDetails := CapabilityDetails{}

		err := rlp.DecodeBytes(packet.ConsensusData[startIndex:], &capabilityDetails)
		if err != nil {
			log.Debug("PeerHandler HandleConsensusPacket", "error", err)
			return err
		}

		go p.HandleCapabilityPacket(&capabilityDetails, fromPeerId)
	} else if packetType == CONSENSUS_PACKET_TYPE_SYNC {
		requestConsensusSyncDetails := RequestConsensusSyncDetails{}

		err := rlp.DecodeBytes(packet.ConsensusData[startIndex:], &requestConsensusSyncDetails)
		if err != nil {
			log.Debug("PeerHandler HandleConsensusPacket", "error", err)
			return err
		}

		go p.HandleRequestConsensusSync(&requestConsensusSyncDetails, fromPeerId)
	} else if packetType >= CONSENSUS_PACKET_TYPE_PROPOSE_BLOCK && packetType <= CONSENSUS_PACKET_TYPE_COMMIT_BLOCK {
		if p.isConsensusRelay {
			go p.BroadcastToSyncPeers(packet, fromPeerId)
		}
	} else {
		log.Debug("PeerHandler unhandled packet type", "packetType", packetType, "fromPeerId", fromPeerId)
	}

	return nil
}

func (p *PeerHandler) HandleCapabilityPacket(capabilityDetails *CapabilityDetails, fromPeerId string) {
	log.Trace("PeerHandler HandleCapabilityPacket", "fromPeerId", fromPeerId)
	debug.PrintStack()
	if capabilityDetails.IsConsensusRelay == false || fromPeerId != capabilityDetails.PeerId {
		log.Debug("PeerHandler HandleCapabilityPacket", "fromPeerId", fromPeerId, "capabilityDetails peeerId", capabilityDetails.PeerId)
		return
	}
	p.peerLock.Lock()
	defer p.peerLock.Unlock()

	log.Trace("PeerHandler HandleCapabilityPacket unlock", "fromPeerId", fromPeerId)

	p.peerMap[capabilityDetails.PeerId] = &PeerDetails{
		peerId:            capabilityDetails.PeerId,
		capabilityDetails: capabilityDetails,
	}

	if p.isConsensusRelay || len(p.consensusRelayMap) == 0 {
		go p.SendRequestConsensusSyncPacket(capabilityDetails.PeerId)
	}
}

func (p *PeerHandler) HandleRequestConsensusSync(requestConsensusSyncDetails *RequestConsensusSyncDetails, fromPeerId string) {
	log.Trace("PeerHandler HandleRequestConsensusSync", "fromPeerId", fromPeerId)
	debug.PrintStack()
	if fromPeerId != requestConsensusSyncDetails.PeerId {
		log.Debug("PeerHandler HandleRequestConsensusSync", "fromPeerId", fromPeerId, "requestConsensusSyncDetails", requestConsensusSyncDetails.PeerId)
		return
	}
	p.peerLock.Lock()
	defer p.peerLock.Unlock()

	p.syncPeerMap[requestConsensusSyncDetails.PeerId] = true
}

func (p *PeerHandler) HandleRequestConsensusDataPacket(packet *eth.RequestConsensusDataPacket) ([]*eth.ConsensusPacket, error) {
	return make([]*eth.ConsensusPacket, 0), nil
}

func (p *PeerHandler) CreateConsensusPacket(data []byte) (*eth.ConsensusPacket, error) {
	log.Debug("PeerHandler CreateConsensusPacket")

	if p.signFn == nil {
		return nil, errors.New("signFn is not set")
	}
	dataToSign := append(ZERO_HASH.Bytes(), data...)
	var signature []byte
	var err error

	signature, err = p.signFn(p.account, accounts.MimetypeProofOfStake, dataToSign)

	if err != nil {
		log.Trace("PeerHandler CreateConsensusPacket failed", "err", err)
		return nil, err
	}

	packet := &eth.ConsensusPacket{
		ParentHash: ZERO_HASH,
	}

	packet.ConsensusData = make([]byte, len(data))
	copy(packet.ConsensusData, data)

	packet.Signature = make([]byte, len(signature))
	copy(packet.Signature, signature)

	return packet, nil
}

func (p *PeerHandler) SendCapabilityPacket(peerList []string) error {
	log.Debug("PeerHandler SendCapabilityPacket", "peer count", len(peerList))
	if p.p2pHandler == nil || p.isConsensusRelay == false || p.getLatestBlockNumberFn() < PACKET_PROTOCOL_START_BLOCK {
		return nil
	}

	capabilityDetails := &CapabilityDetails{
		IsConsensusRelay: true,
		PeerId:           p.localPeerId,
	}

	data, err := rlp.EncodeToBytes(capabilityDetails)

	if err != nil {
		log.Debug("PeerHandler SendCapabilityPacket EncodeToBytes", "error")
		return err
	}

	var dataToSend []byte
	dataToSend = append([]byte{ConsensusNetworkProtocolVersion}, append([]byte{byte(CONSENSUS_PACKET_TYPE_CAPABILITY)}, data...)...)

	packet, err := p.CreateConsensusPacket(dataToSend)
	if err != nil {
		log.Debug("PeerHandler SendCapabilityPacket CreateConsensusPacket", "error", err)
		return err
	}

	err = p.p2pHandler.SendConsensusPacket(peerList, packet)
	if err != nil {
		log.Debug("PeerHandler SendCapabilityPacket SendConsensusPacket", "error", err)
		return err
	}

	log.Trace("PeerHandler SendCapabilityPacket", "peer count", len(peerList))
	return nil
}

func (p *PeerHandler) ConnectAvailableConsensusRelay() {
	log.Trace("PeerHandler ConnectConsensusRelay lock")
	p.peerLock.Lock()
	defer p.peerLock.Unlock()
	log.Trace("PeerHandler ConnectConsensusRelay Unlock")

	for k, v := range p.peerMap {
		if v.capabilityDetails.IsConsensusRelay {
			go p.SendRequestConsensusSyncPacket(k)
			break
		}
	}
}

func (p *PeerHandler) SendRequestConsensusSyncPacket(peerId string) error {
	log.Trace("PeerHandler SendRequestConsensusSyncPacket", "peerId", peerId)
	if p.p2pHandler == nil || p.getLatestBlockNumberFn() < PACKET_PROTOCOL_START_BLOCK {
		log.Debug("PeerHandler SendRequestConsensusSyncPacket return", "peerId", peerId)
		return nil
	}

	consensusSyncDetails := &RequestConsensusSyncDetails{
		IsConsensusRelay: p.isConsensusRelay,
		PeerId:           p.localPeerId,
	}

	data, err := rlp.EncodeToBytes(consensusSyncDetails)

	if err != nil {
		log.Debug("PeerHandler SendRequestConsensusSyncPacket EncodeToBytes", "error", err, "peer", peerId)
		return err
	}

	var dataToSend []byte
	dataToSend = append([]byte{ConsensusNetworkProtocolVersion}, append([]byte{byte(CONSENSUS_PACKET_TYPE_SYNC)}, data...)...)

	packet, err := p.CreateConsensusPacket(dataToSend)
	if err != nil {
		log.Debug("PeerHandler SendRequestConsensusSyncPacket CreateConsensusPacket", "error", err, "peer", peerId)
		return err
	}

	err = p.p2pHandler.SendConsensusPacket([]string{peerId}, packet)
	if err != nil {
		log.Debug("PeerHandler SendRequestConsensusSyncPacket SendConsensusPacket", "error", err, "peer", peerId)
		return err
	}

	p.peerLock.Lock()
	defer p.peerLock.Unlock()
	p.consensusRelayMap[peerId] = true
	log.Trace("PeerHandler SendRequestConsensusSyncPacket done", "peerId", peerId)
	return nil
}

func (p *PeerHandler) ShouldRebroadCast(packet *eth.ConsensusPacket, fromPeerId string) bool {
	return false
}

func (p *PeerHandler) BroadcastLocalPacket(packet *eth.ConsensusPacket) int {
	if p.isConsensusRelay == true {
		syncPeerCount := p.BroadcastToSyncPeers(packet, p.localPeerId)
		consensusRelayCount := p.BroadcastToConsensusRelays(packet, p.localPeerId)
		return syncPeerCount + consensusRelayCount
	} else {
		return p.BroadcastToConsensusRelays(packet, p.localPeerId)
	}
}

func (p *PeerHandler) BroadcastToConsensusRelays(packet *eth.ConsensusPacket, fromPeerId string) int {
	p.peerLock.Lock()
	defer p.peerLock.Unlock()

	sendList := make([]string, 0)

	var packetSyncDetails *PacketSyncDetails
	packetSyncDetails, ok := p.packetSyncMap[packet.Hash()]
	if ok == false {
		packetSyncDetails = &PacketSyncDetails{
			incomingPeerMap: make(map[string]bool),
			packet:          packet,
			sendPeerMap:     make(map[string]bool),
		}
		p.packetSyncMap[packet.Hash()] = packetSyncDetails
	}

	sendPeerMap := packetSyncDetails.sendPeerMap

	alreadySentCount := 0
	for k, _ := range p.consensusRelayMap {
		_, ok := sendPeerMap[k]
		if ok {
			alreadySentCount = alreadySentCount + 1
			continue
		}
		sendList = append(sendList, []string{k}...)
		sendPeerMap[k] = true
	}

	packetSyncDetails.sendPeerMap = sendPeerMap
	p.packetSyncMap[packet.Hash()] = packetSyncDetails

	log.Info("BroadcastToConsensusRelays", "relay count", len(p.consensusRelayMap), "send list count", len(sendList), "alreadySentCount", alreadySentCount, "packetHash", packet.Hash(), "parentHash", packet.ParentHash)
	go p.p2pHandler.SendConsensusPacket(sendList, packet)

	return len(sendList)
}

func (p *PeerHandler) BroadcastToSyncPeers(packet *eth.ConsensusPacket, fromPeerId string) int {
	log.Trace("BroadcastToSyncPeers", "fromPeerId", fromPeerId, "packetHash", packet.Hash(), "parentHash", packet.ParentHash)
	p.peerLock.Lock()
	defer p.peerLock.Unlock()
	log.Trace("BroadcastToSyncPeers unlock")

	if packet.ParentHash.IsEqualTo(p.GetCurrentParentHash()) == false {
		log.Trace("BroadcastToSyncPeers unlock parentHash not matched")
		return 0
	}

	var packetSyncDetails *PacketSyncDetails
	packetSyncDetails, ok := p.packetSyncMap[packet.Hash()]
	if ok == false {
		packetSyncDetails = &PacketSyncDetails{
			incomingPeerMap: make(map[string]bool),
			packet:          packet,
			sendPeerMap:     make(map[string]bool),
		}
		p.packetSyncMap[packet.Hash()] = packetSyncDetails
	}

	incomingPeerMap := packetSyncDetails.incomingPeerMap
	incomingPeerMap[fromPeerId] = true

	sendPeerMap := packetSyncDetails.sendPeerMap

	sendPeerList := make([]string, 0)

	alreadySentCount := 0
	for peerId, _ := range p.syncPeerMap {
		if peerId == fromPeerId {
			continue
		}
		_, ok := sendPeerMap[peerId]
		if ok {
			alreadySentCount = alreadySentCount + 1
			continue
		}
		sendPeerList = append(sendPeerList, []string{peerId}...)
		sendPeerMap[peerId] = true
	}

	packetSyncDetails.incomingPeerMap = incomingPeerMap
	packetSyncDetails.sendPeerMap = sendPeerMap
	p.packetSyncMap[packet.Hash()] = packetSyncDetails

	log.Info("BroadcastToSyncPeers", "sendPeerMap count", len(sendPeerList), "sendPeerList count", len(sendPeerList), "syncPeerMap count", len(p.syncPeerMap), "alreadySentCount", alreadySentCount,
		"packetHash", packet.Hash(), "parentHash", packet.ParentHash)

	go p.p2pHandler.SendConsensusPacket(sendPeerList, packet)

	return len(sendPeerList)
}

func (p *PeerHandler) GetCurrentParentHash() common.Hash {
	p.parentHashLock.Lock()
	defer p.parentHashLock.Unlock()
	return p.currentParentHash
}

func (p *PeerHandler) SetCurrentParentHash(parentHash common.Hash, currentBlockNumber uint64) {
	p.parentHashLock.Lock()
	defer p.parentHashLock.Unlock()

	p.peerLock.Lock()
	defer p.peerLock.Unlock()

	p.currentParentHash = parentHash
	p.currentBlockNumber = currentBlockNumber

	//Cleanup old packets
	for k, v := range p.packetSyncMap {
		if v.packet.ParentHash.IsEqualTo(p.currentParentHash) == true {
			continue
		}
		delete(p.packetSyncMap, k)
	}

	if currentBlockNumber == PACKET_PROTOCOL_START_BLOCK { //Special case, to trigger on-going connections
		if p.isConsensusRelay {
			go p.SendCapabilityToAllPeers()
		}
	}
}

func (p *PeerHandler) SendCapabilityToAllPeers() {
	p.peerLock.Lock()
	defer p.peerLock.Unlock()

	peerList := make([]string, 0)

	for peerId, _ := range p.peerMap {
		peerList = append(peerList, []string{peerId}...)
	}

	p.SendCapabilityPacket(peerList)
}
