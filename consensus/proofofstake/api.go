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

package proofofstake

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/consensus"
	"github.com/DogeProtocol/dp/internal/ethapi"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/rlp"
	"github.com/DogeProtocol/dp/rpc"
	"github.com/DogeProtocol/dp/systemcontracts/conversion"
	"math/big"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	chain        consensus.ChainHeaderReader
	proofofstake *ProofOfStake
}

type blockNumberOrHashOrRLP struct {
	*rpc.BlockNumberOrHash
	RLP hexutil.Bytes `json:"rlp,omitempty"`
}

func (sb *blockNumberOrHashOrRLP) UnmarshalJSON(data []byte) error {
	bnOrHash := new(rpc.BlockNumberOrHash)
	// Try to unmarshal bNrOrHash
	if err := bnOrHash.UnmarshalJSON(data); err == nil {
		sb.BlockNumberOrHash = bnOrHash
		return nil
	}
	// Try to unmarshal RLP
	var input string
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	sb.RLP = hexutil.MustDecode(input)
	return nil
}

// ListValidators retrieves the list of authorized signers at the specified block.
func (api *API) ListValidators(blockNumberHex string) ([]*ValidatorDetails, error) {
	var blockNumber uint64
	var err error
	if blockNumberHex == "" || len(blockNumberHex) == 0 {
		blockNumber = api.chain.CurrentHeader().Number.Uint64()
	} else {
		blockNumber, err = hexutil.DecodeUint64(blockNumberHex)
		if err != nil {
			return nil, err
		}
	}

	var header = api.chain.GetHeaderByNumber(blockNumber)
	if header == nil {
		return nil, errUnknownBlock
	}
	validators, err := api.proofofstake.ListValidators(header.Hash())
	if err != nil {
		return nil, err
	}
	return validators, nil
}

type StakingData struct {
	TotalDepositedBalance string              `json:"totalDepositedBalance"     gencodec:"required"`
	Validators            []*ValidatorDetails `json:"validators"     gencodec:"required"`
}

func (api *API) GetStakingDetailsByValidatorAddress(validator common.Address, blockNumberHex string) (*ValidatorDetails, error) {
	var blockNumber uint64
	var err error
	if blockNumberHex == "" || len(blockNumberHex) == 0 {
		blockNumber = api.chain.CurrentHeader().Number.Uint64()
	} else {
		blockNumber, err = hexutil.DecodeUint64(blockNumberHex)
		if err != nil {
			return nil, err
		}
	}

	// Retrieve the requested block number (or current if none requested)
	var header = api.chain.GetHeaderByNumber(blockNumber)
	if header == nil {
		return nil, errUnknownBlock
	}

	return api.proofofstake.GetStakingDetailsByValidatorAddress(validator, header.Hash())
}

func (api *API) GetStakingDetailsByDepositorAddress(depositor common.Address, blockNumberHex string) (*ValidatorDetails, error) {
	var blockNumber uint64
	var err error
	if blockNumberHex == "" || len(blockNumberHex) == 0 {
		blockNumber = api.chain.CurrentHeader().Number.Uint64()
	} else {
		blockNumber, err = hexutil.DecodeUint64(blockNumberHex)
		if err != nil {
			return nil, err
		}
	}

	// Retrieve the requested block number (or current if none requested)
	var header = api.chain.GetHeaderByNumber(blockNumber)
	if header == nil {
		return nil, errUnknownBlock
	}

	validator, err := api.proofofstake.GetValidatorOfDepositor(depositor, header.Hash())
	if err != nil {
		return nil, err
	}

	return api.proofofstake.GetStakingDetailsByValidatorAddress(validator, header.Hash())
}

// GetStakingDetails retrieves the total deposited quantity.
func (api *API) GetStakingDetails(blockNumberHex string) (*StakingData, error) {
	var blockNumber uint64
	var err error
	if blockNumberHex == "" || len(blockNumberHex) == 0 {
		blockNumber = api.chain.CurrentHeader().Number.Uint64()
	} else {
		blockNumber, err = hexutil.DecodeUint64(blockNumberHex)
		if err != nil {
			return nil, err
		}
	}

	// Retrieve the requested block number (or current if none requested)
	var header = api.chain.GetHeaderByNumber(blockNumber)
	if header == nil {
		return nil, errUnknownBlock
	}
	balance, err := api.proofofstake.GetTotalDepositedBalance(header.Hash())
	if err != nil {
		return nil, err
	}

	validators, err := api.proofofstake.ListValidators(header.Hash())
	if err != nil {
		return nil, err
	}

	return &StakingData{
		TotalDepositedBalance: hexutil.EncodeBig(balance),
		Validators:            validators,
	}, nil
}

type ExtendedConsensusPacket struct {
	Signer     common.Address `json:"signer"     gencodec:"required"`
	PacketType byte           `json:"packetType" gencodec:"required"`
	Round      byte           `json:"round"      gencodec:"required"`
}

type ConsensusData struct {
	Data                     *BlockConsensusData           `json:"data"     gencodec:"required"`
	AdditionalData           *BlockAdditionalConsensusData `json:"additionalData"     gencodec:"required"`
	ExtendedConsensusPackets []*ExtendedConsensusPacket    `json:"extendedConsensusPackets"     gencodec:"required"`
	BlockProposerRewards     string                        `json:"blockProposerRewards"     gencodec:"required"`
}

// GetBlockConsensusData retrieves proofofstake consensus data of the block.
func (api *API) GetBlockConsensusData(blockNumberHex string) (*ConsensusData, error) {
	var blockNumber uint64
	var err error
	if blockNumberHex == "" || len(blockNumberHex) == 0 {
		blockNumber = api.chain.CurrentHeader().Number.Uint64()
	} else {
		blockNumber, err = hexutil.DecodeUint64(blockNumberHex)
		if err != nil {
			return nil, err
		}
	}

	blockConsensusData := &BlockConsensusData{}
	header := api.chain.GetHeaderByNumber(blockNumber)
	if header == nil {
		return nil, errUnknownBlock
	}

	err = rlp.DecodeBytes(header.ConsensusData, &blockConsensusData)
	if err != nil {
		return nil, err
	}

	blockAdditionalConsensusData := &BlockAdditionalConsensusData{}
	err = rlp.DecodeBytes(header.UnhashedConsensusData, blockAdditionalConsensusData)
	if err != nil {
		return nil, err
	}

	consensusData := &ConsensusData{
		Data:           blockConsensusData,
		AdditionalData: blockAdditionalConsensusData,
	}

	consensusData.ExtendedConsensusPackets = make([]*ExtendedConsensusPacket, 0)
	for i := 0; i < len(blockAdditionalConsensusData.ConsensusPackets); i++ {
		packet := blockAdditionalConsensusData.ConsensusPackets[i]
		round, signer, err := parsePacket(&packet)
		if err != nil {
			consensusData.ExtendedConsensusPackets = append(consensusData.ExtendedConsensusPackets, &ExtendedConsensusPacket{})
			continue
		}
		ePacket := ExtendedConsensusPacket{
			PacketType: packet.ConsensusData[0],
			Round:      round,
		}
		ePacket.Signer.CopyFrom(signer)

		consensusData.ExtendedConsensusPackets = append(consensusData.ExtendedConsensusPackets, &ePacket)
	}

	if blockConsensusData.VoteType == VOTE_TYPE_OK {
		blockRewards := GetReward(header.Number)
		consensusData.BlockProposerRewards = hexutil.EncodeBig(blockRewards)
	} else {
		consensusData.BlockProposerRewards = hexutil.EncodeUint64(0)
	}

	return consensusData, err
}

type ConversionDetails struct {
	EthAddress     common.Address `json:"ethAddress"     gencodec:"required"`
	QuantumAddress common.Address `json:"quantumAddress"     gencodec:"required"`
	IsConverted    bool           `json:"isConverted"     gencodec:"required"`
	Coins          *big.Int       `json:"coins"     gencodec:"required"`
}

// GetConversionDetails returns whether the ethereum address is converted or not and details on the conversion
func (api *API) GetConversionDetails(ethAddressHex string) (*ConversionDetails, error) {
	var header = api.chain.CurrentHeader()
	if header == nil {
		return nil, errUnknownBlock
	}

	ethAddress := common.HexToAddress(ethAddressHex)
	isConverted, err := api.getConversionStatus(ethAddress, header.Hash())
	if err != nil {
		return nil, err
	}

	coins, err := api.GetCoinsForEthereumAddress(ethAddress, header.Hash())
	if err != nil {
		return nil, err
	}

	var quantumAddress common.Address
	if isConverted {
		quantumAddress, err = api.getConversionQuantumAddress(ethAddress, header.Hash())
		if err != nil {
			return nil, err
		}
	} else {
		quantumAddress = ZERO_ADDRESS
	}

	conversionDetails := &ConversionDetails{
		EthAddress:     ethAddress,
		IsConverted:    isConverted,
		QuantumAddress: quantumAddress,
		Coins:          coins,
	}

	return conversionDetails, nil
}

func (api *API) getConversionStatus(ethAddress common.Address, blockHash common.Hash) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := conversion.GetContract_Method_getConversionStatus()

	abiData, err := conversion.GetConversionContract_ABI()
	if err != nil {
		log.Error("getConversionStatus abi error", "err", err)
		return false, err
	}
	contractAddress := common.HexToAddress(conversion.CONVERSION_CONTRACT)

	// call
	data, err := abiData.Pack(method, ethAddress)
	if err != nil {
		log.Error("Unable to pack tx for getConversionStatus", "error", err)
		return false, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	result, err := api.proofofstake.ethAPI.Call(ctx, ethapi.TransactionArgs{
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		log.Error("Call", "err", err)
		return false, err
	}
	if len(result) == 0 {
		return false, errors.New("getConversionStatus result is 0")
	}

	var out bool

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Debug("UnpackIntoInterface", "err", err, "ethAddress", ethAddress)
		return false, err
	}

	return out, nil
}

func (api *API) getConversionQuantumAddress(ethAddress common.Address, blockHash common.Hash) (common.Address, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := conversion.GetContract_Method_getQuantumAddress()

	abiData, err := conversion.GetConversionContract_ABI()
	if err != nil {
		log.Error("getConversionQuantumAddress abi error", "err", err)
		return ZERO_ADDRESS, err
	}
	contractAddress := common.HexToAddress(conversion.CONVERSION_CONTRACT)

	// call
	data, err := abiData.Pack(method, ethAddress)
	if err != nil {
		log.Error("Unable to pack tx for getConversionQuantumAddress", "error", err)
		return ZERO_ADDRESS, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	result, err := api.proofofstake.ethAPI.Call(ctx, ethapi.TransactionArgs{
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		log.Error("Call", "err", err)
		return ZERO_ADDRESS, err
	}
	if len(result) == 0 {
		return ZERO_ADDRESS, errors.New("getConversionQuantumAddress result is 0")
	}

	var (
		ret0 = new(common.Address)
	)
	out := ret0

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Debug("UnpackIntoInterface", "err", err, "ethAddress", ethAddress)
		return ZERO_ADDRESS, err
	}

	return *out, nil
}

func (api *API) GetCoinsForEthereumAddress(ethAddress common.Address, blockHash common.Hash) (*big.Int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := conversion.GetContract_Method_getAmount()

	abiData, err := conversion.GetConversionContract_ABI()
	if err != nil {
		log.Error("GetCoinsForEthereumAddress abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(conversion.CONVERSION_CONTRACT)

	// call
	data, err := abiData.Pack(method, ethAddress)
	if err != nil {
		log.Error("Unable to pack tx for GetCoinsForEthereumAddress", "error", err)
		return nil, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	result, err := api.proofofstake.ethAPI.Call(ctx, ethapi.TransactionArgs{
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		log.Error("Call", "err", err)
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("GetCoinsForEthereumAddress result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Debug("UnpackIntoInterface", "err", err, "ethAddress", ethAddress)
		return nil, err
	}

	return out, nil
}
