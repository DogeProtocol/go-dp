package proofofstake

import (
	"context"
	"errors"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/core/state"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/internal/ethapi"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/rpc"
	"github.com/DogeProtocol/dp/systemcontracts/consensuscontext"
	"math"
	"strconv"
)

func (p *ProofOfStake) SetConsensusContext(key string, context [32]byte, state *state.StateDB, header *types.Header) error {
	method := consensuscontext.SET_CONTEXT_METHOD
	abiData, err := consensuscontext.GetConsensusContract_ABI()
	if err != nil {
		log.Error("SetConsensusContext abi error", "err", err)
		return err
	}
	contractAddress := consensuscontext.CONSENSUS_CONTEXT_CONTRACT_ADDRESS

	data, err := encodeCall(&abiData, method, key, context)
	if err != nil {
		log.Error("Unable to pack SetConsensusContext", "error", err)
		return err
	}

	msgData := (hexutil.Bytes)(data)
	var from common.Address
	from.CopyFrom(ZERO_ADDRESS)
	args := ethapi.TransactionArgs{
		From: &from,
		To:   &contractAddress,
		Data: &msgData,
	}

	msg, err := args.ToMessage(math.MaxUint64)
	if err != nil {
		return err
	}

	_, err = p.blockchain.ExecuteNoGas(msg, state, header)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProofOfStake) DeleteConsensusContext(key string, state *state.StateDB, header *types.Header) error {
	method := consensuscontext.DELETE_CONTEXT_METHOD
	abiData, err := consensuscontext.GetConsensusContract_ABI()
	if err != nil {
		log.Error("DeleteConsensusContext abi error", "err", err)
		return err
	}
	contractAddress := consensuscontext.CONSENSUS_CONTEXT_CONTRACT_ADDRESS

	data, err := encodeCall(&abiData, method, key)
	if err != nil {
		log.Error("Unable to pack DeleteConsensusContext", "error", err)
		return err
	}

	msgData := (hexutil.Bytes)(data)
	var from common.Address
	from.CopyFrom(ZERO_ADDRESS)
	args := ethapi.TransactionArgs{
		From: &from,
		To:   &contractAddress,
		Data: &msgData,
	}

	msg, err := args.ToMessage(math.MaxUint64)
	if err != nil {
		return err
	}

	_, err = p.blockchain.ExecuteNoGas(msg, state, header)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProofOfStake) GetConsensusContext(key string, blockHash common.Hash) ([32]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	var out [32]byte

	method := consensuscontext.GET_CONTEXT_METHOD

	abiData, err := consensuscontext.GetConsensusContract_ABI()
	if err != nil {
		log.Error("GetConsensusContext abi error", "err", err)
		return out, err
	}
	contractAddress := consensuscontext.CONSENSUS_CONTEXT_CONTRACT_ADDRESS

	// call
	data, err := abiData.Pack(method, key)
	if err != nil {
		log.Error("Unable to pack tx for GetConsensusContext", "error", err)
		return out, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		log.Error("Call", "err", err)
		return out, err
	}
	if len(result) == 0 {
		return out, errors.New("GetConsensusContext result is 0")
	}

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Debug("UnpackIntoInterface", "err", err, "key", key)
		return out, err
	}

	return out, nil
}

func GetBlockConsensusContextKey(blockNumber uint64) (string, error) {
	var key string
	if blockNumber <= CONSENSUS_CONTEXT_START_BLOCK {
		return key, errors.New("GetBlockConsensusContext blockNumber below CONSENSUS_CONTEXT_START_BLOCK")
	}

	//bc for block context
	key = "bc-" + strconv.FormatUint(blockNumber, 10)

	return key, nil
}
