package proofofstake

import (
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/eth/protocols/eth"
)

type P2PHandler interface {
	SendConsensusPacket(peerList []string, packet *eth.ConsensusPacket) error
	BroadcastConsensusData(packet *eth.ConsensusPacket) error
	RequestTransactions(txns []common.Hash) error
	RequestConsensusData(packet *eth.RequestConsensusDataPacket) error
	GetLocalPeerId() string
}
