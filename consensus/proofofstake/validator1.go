package proofofstake

import (
	"context"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/internal/ethapi"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/rpc"
	"github.com/DogeProtocol/dp/systemcontracts1"
)

func (p *ProofOfStake) GetValidatorsAddress1(number uint64, blockHash common.Hash) ([]common.Address, error) {
	err := systemcontracts1.IsStakingContract()
	if err != nil {
		log.Warn("GETH_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}
	//blockNumber = new(big.Int).SetUint64(172)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := systemcontracts1.GetContract_Method_ListValidator()
	abiData := systemcontracts1.GetStakingContract_ABI()
	contractAddress := common.HexToAddress(systemcontracts1.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method)
	if err != nil {
		log.Error("Unable to pack tx for get validators", "error", err)
		return nil, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	//gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		//Gas:  &gas,
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, err
	}
	///fmt.Println("result...", result)

	var (
		ret0 = new([]common.Address)
	)
	out := ret0

	if err := abiData.UnpackIntoInterface(out, method, result); err != nil {
		return nil, err
	}

	valz := make([]common.Address, len(*ret0))
	for i, a := range *out {
		///fmt.Println("Get validator1 ID len(a), len(*ret0), a", len(a), len(*ret0), a)
		valz[i] = a
	}
	return valz, nil
}

func (p *ProofOfStake) GetDepositor(validator common.Address, blockHash common.Hash) (common.Address, error) {
	err := systemcontracts1.IsStakingContract()
	if err != nil {
		log.Warn("GETH_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return common.Address{}, err
	}
	//blockNumber = new(big.Int).SetUint64(172)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := systemcontracts1.GetContract_Method_GetDepositor()
	abiData := systemcontracts1.GetStakingContract_ABI()
	contractAddress := common.HexToAddress(systemcontracts1.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, validator)
	if err != nil {
		log.Error("Unable to pack tx for get depositor", "error", err)
		return common.Address{}, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	//gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		//Gas:  &gas,
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		return common.Address{}, err
	}
	if len(result) == 0 {
		return common.Address{}, err
	}
	///fmt.Println("result...", result)

	var (
		ret0 = new(common.Address)
	)
	out := ret0

	if err := abiData.UnpackIntoInterface(out, method, result); err != nil {
		return common.Address{}, err
	}

	return *out, nil
}
