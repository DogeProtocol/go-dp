package proofofstake

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/systemcontracts"
)

func (p *ProofOfStake) GetValidators1(blockHash common.Hash) ([]common.Address, error) {
	valz := make([]common.Address, 3)
	valz[0] = common.HexToAddress("4643635a54Ca29C1E803B9c0Eca489426757c4C2")
	valz[1] = common.HexToAddress("46f8c16c50b122a568c96FB5e97E44ca9cD205CE")
	valz[2] = common.HexToAddress("5836eFD181459F6b3BF90816c13d1e27E4F346Ad")
	return valz, nil
}

func (p *ProofOfStake) GetValidators(number uint64, blockHash common.Hash) ([]common.Address, error) {

	err := systemcontracts.IsStakingContract()
	if err != nil {
		log.Warn("GETH_STAKING_CONTRACT : Contract address is empty")
		return nil, err
	}
	//blockNumber = new(big.Int).SetUint64(172)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := systemcontracts.GetContract_Method_ListValidator()
	abiData := systemcontracts.GetStakingContract_ABI()
	contractAddress := common.HexToAddress(systemcontracts.GetStakingContract_Address_String())

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
		///fmt.Println("Get validator ID len(a), len(*ret0), a", len(a), len(*ret0), a)
		valz[i] = a
	}
	return valz, nil
}
