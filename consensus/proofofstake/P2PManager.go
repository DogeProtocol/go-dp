package proofofstake

import (
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/eth/protocols/eth"
)

type P2PHandler interface {
	BroadcastConsensusData(packet *eth.ConsensusPacket) error
	RequestTransactions(txns []common.Hash) error
	RequestConsensusData(packet *eth.RequestConsensusDataPacket) error
}
