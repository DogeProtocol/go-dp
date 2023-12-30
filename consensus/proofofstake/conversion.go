package proofofstake

import (
	"errors"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/core/state"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/internal/ethapi"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/systemcontracts/conversion"
	"math"
	"math/big"
)

func (p *ProofOfStake) GetCoinsForEthereumAddress(ethAddress common.Address, state *state.StateDB, header *types.Header) (*big.Int, error) {
	method := conversion.GetContract_Method_getAmount()
	abiData, err := conversion.GetConversionContract_ABI()
	if err != nil {
		log.Error("GetCoinsForEthereumAddress abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(conversion.CONVERSION_CONTRACT)

	// call
	data, err := encodeCall(&abiData, method, ethAddress)
	if err != nil {
		log.Error("Unable to pack GetCoinsForEthereumAddress", "error", err)
		return nil, err
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
		return nil, err
	}

	result, err := p.blockchain.ExecuteNoGas(msg, state, header)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("GetCoinsForEthereumAddress result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Error("GetCoinsForEthereumAddress UnpackIntoInterface", "err", err, "ethAddress", ethAddress)
		return nil, err
	}

	log.Info("===============GetCoinsForEthereumAddress", "out", out)

	return out, nil
}

func (p *ProofOfStake) GetConversionStatus(ethAddress common.Address, state *state.StateDB, header *types.Header) (bool, error) {
	method := conversion.GetContract_Method_getConversionStatus()
	abiData, err := conversion.GetConversionContract_ABI()
	if err != nil {
		log.Error("GetConversionStatus abi error", "err", err)
		return false, err
	}
	contractAddress := common.HexToAddress(conversion.CONVERSION_CONTRACT)

	// call
	data, err := encodeCall(&abiData, method, ethAddress)
	if err != nil {
		log.Error("Unable to pack GetConversionStatus", "error", err)
		return false, err
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
		return false, err
	}

	result, err := p.blockchain.ExecuteNoGas(msg, state, header)
	if err != nil {
		return false, err
	}

	if len(result) == 0 {
		return false, errors.New("GetConversionStatus result is 0")
	}

	var out bool

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Error("GetConversionStatus UnpackIntoInterface", "err", err, "ethAddress", ethAddress)
		return false, err
	}

	log.Info("===============GetConversionStatus", "out", out)

	return out, nil
}

func (p *ProofOfStake) SetConverted(ethereumAddress common.Address, quantumAddress common.Address,
	state *state.StateDB, header *types.Header) (*big.Int, error) {
	method := conversion.GetContract_Method_setConverted()
	abiData, err := conversion.GetConversionContract_ABI()
	if err != nil {
		log.Error("SetConverted abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(conversion.CONVERSION_CONTRACT)

	// call
	data, err := encodeCall(&abiData, method, ethereumAddress, quantumAddress)
	if err != nil {
		log.Error("Unable to pack SetConverted", "error", err)
		return nil, err
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
		return nil, err
	}

	result, err := p.blockchain.ExecuteNoGas(msg, state, header)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("SetConverted result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Error("SetConverted UnpackIntoInterface", "err", err, "ethAddress", ethereumAddress)
		return nil, err
	}

	log.Info("===============SetConverted", "out", out)

	return out, nil
}

func (p *ProofOfStake) GetQuantumAddress(ethAddress common.Address, state *state.StateDB, header *types.Header) (common.Address, error) {
	method := conversion.GetContract_Method_getQuantumAddress()
	abiData, err := conversion.GetConversionContract_ABI()
	if err != nil {
		log.Error("GetQuantumAddress abi error", "err", err)
		return ZERO_ADDRESS, err
	}
	contractAddress := common.HexToAddress(conversion.CONVERSION_CONTRACT)

	// call
	data, err := encodeCall(&abiData, method, ethAddress)
	if err != nil {
		log.Error("Unable to pack GetQuantumAddress", "error", err)
		return ZERO_ADDRESS, err
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
		return ZERO_ADDRESS, err
	}

	result, err := p.blockchain.ExecuteNoGas(msg, state, header)
	if err != nil {
		return ZERO_ADDRESS, err
	}

	if len(result) == 0 {
		return ZERO_ADDRESS, errors.New("GetQuantumAddress result is 0")
	}

	var out common.Address

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Error("GetQuantumAddress UnpackIntoInterface", "err", err, "ethAddress", ethAddress)
		return ZERO_ADDRESS, err
	}

	log.Info("===============GetQuantumAddress", "out", out)

	return out, nil
}
