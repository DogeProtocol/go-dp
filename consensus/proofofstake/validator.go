package proofofstake

import (
	"context"
	"errors"
	"github.com/DogeProtocol/dp/accounts/abi"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/core/state"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/internal/ethapi"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/rpc"
	"github.com/DogeProtocol/dp/systemcontracts/staking"
	"math"
	"math/big"
)

type ValidatorDetails struct {
	Depositor          common.Address `json:"depositor"     gencodec:"required"`
	Validator          common.Address `json:"validator"     gencodec:"required"`
	Balance            *big.Int       `json:"balance"       gencodec:"required"`
	NetBalance         *big.Int       `json:"netBalance"    gencodec:"required"`
	BlockRewards       *big.Int       `json:"blockRewards"  gencodec:"required"`
	Slashings          *big.Int       `json:"slashings"  	gencodec:"required"`
	IsValidationPaused bool           `json:"isValidationPaused"  gencodec:"required"`
	WithdrawalBlock    *big.Int       `json:"withdrawalBlock"  gencodec:"required"`
}

func (p *ProofOfStake) GetValidators(blockHash common.Hash) (map[common.Address]*big.Int, error) {
	depositorCount, err := p.GetDepositorCount(blockHash)
	if err != nil {
		return nil, err
	} else {
		log.Debug("depositorCount", "depositorCount", depositorCount)
	}
	totalDepositedBalance, err := p.GetTotalDepositedBalance(blockHash)
	if err != nil {
		log.Debug("totalDepositedBalance error", "err", err)
		return nil, err
	} else {
		log.Debug("totalDepositedBalance", "totalDepositedBalance", totalDepositedBalance)
	}

	err = staking.IsStakingContract()
	if err != nil {
		log.Warn("GETH_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_ListValidators()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("GetValidators error getting abidata", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())
	// call
	data, err := abiData.Pack(method)
	if err != nil {
		log.Error("Unable to pack tx for get filteredValidatorsDepositMap", "error", err)
		return nil, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)

	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		log.Debug("result 0 length")
		return nil, nil
	}
	_, err = abiData.Unpack(method, result)
	if err != nil {
		log.Error("Unpack", "err", err)
		return nil, err
	}

	var (
		ret0 = new([]common.Address)
	)
	out := ret0

	if err := abiData.UnpackIntoInterface(out, method, result); err != nil {
		log.Info("UnpackIntoInterface error")
		return nil, err
	}

	proposalsTxnsMap := make(map[common.Address]*big.Int)
	for _, val := range *out {
		if val.IsEqualTo(ZERO_ADDRESS) {
			return nil, errors.New("invalid validator")
		}
		log.Debug("GetValidators Validator", "val", val)
	}

	for _, val := range *out {
		isPaused, err := p.IsValidatorPaused(val, blockHash)
		if err != nil {
			log.Debug("IsValidatorPaused failed", "err", err)
			return nil, err
		}

		if isPaused {
			log.Debug("Validator is paused, skipping", "val", isPaused)
			continue
		}

		depositor, err := p.GetDepositorOfValidator(val, blockHash)
		if err != nil {
			log.Debug("GetDepositorOfValidator failed", "err", err)
			continue
		}

		if depositor.IsEqualTo(ZERO_ADDRESS) {
			return nil, errors.New("invalid depositor")
		}

		balance, err := p.GetNetBalanceOfDepositor(depositor, blockHash)
		if err != nil {
			log.Debug("GetBalanceOfDepositor failed", "err", err)
			continue
		}

		proposalsTxnsMap[val] = balance
		log.Debug("GetValidators", "validator", val, "depositor", depositor, "depositAmount", balance)
	}

	return proposalsTxnsMap, nil
}

func (p *ProofOfStake) GetValidatorOfDepositor(depositor common.Address, blockHash common.Hash) (common.Address, error) {
	log.Trace("GetValidatorOfDepositor depositor", "depositor", depositor)
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("GETH_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return common.Address{}, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_GetValidatorOfDepositor()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("GetValidatorOfDepositor abi error", "err", err)
		return common.Address{}, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, depositor)
	if err != nil {
		log.Error("Unable to pack tx for get validator", "error", err)
		return common.Address{}, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		//Gas:  &gas,
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		return common.Address{}, err
	}
	if len(result) == 0 {
		return common.Address{}, errors.New("no depositor found")
	}

	var (
		ret0 = new(common.Address)
	)
	out := ret0

	if err := abiData.UnpackIntoInterface(out, method, result); err != nil {
		return common.Address{}, err
	}

	return *out, nil
}

func (p *ProofOfStake) GetDepositorOfValidator(validator common.Address, blockHash common.Hash) (common.Address, error) {
	log.Trace("GetDepositorOfValidator validator", "validator", validator)
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("GETH_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return common.Address{}, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_GetDepositorOfValidator()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("GetDepositorOfValidator abi error", "err", err)
		return common.Address{}, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, validator)
	if err != nil {
		log.Error("Unable to pack tx for get depositor", "error", err)
		return common.Address{}, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		//Gas:  &gas,
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		return common.Address{}, err
	}
	if len(result) == 0 {
		return common.Address{}, errors.New("no depositor found")
	}

	var (
		ret0 = new(common.Address)
	)
	out := ret0

	if err := abiData.UnpackIntoInterface(out, method, result); err != nil {
		return common.Address{}, err
	}

	return *out, nil
}

func (p *ProofOfStake) GetNetBalanceOfDepositor(depositor common.Address, blockHash common.Hash) (*big.Int, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("DP_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_GetNetBalanceOfDepositor() //todo: change once initial storage is set
	//method := staking.GetContract_Method_GetBalanceOfDepositor()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("GetNetBalanceOfDepositor abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, depositor)
	if err != nil {
		log.Error("Unable to pack tx for GetNetBalanceOfDepositor", "error", err)
		return nil, err
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
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("GetNetBalanceOfDepositor result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Debug("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func (p *ProofOfStake) GetDepositorCount(blockHash common.Hash) (*big.Int, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("DP_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_GetDepositorCount()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Trace("GetDepositorCount abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method)
	if err != nil {
		log.Error("Unable to pack tx for get filteredValidatorsDepositMap", "error", err)
		return nil, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)

	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		log.Trace("Call", "err", err)
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("GetDepositorCount result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err)
		return nil, err
	}

	return out, nil
}

func (p *ProofOfStake) GetTotalDepositedBalance(blockHash common.Hash) (*big.Int, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("DP_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_GetTotalDepositedBalance()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Trace("GetTotalDepositedBalance abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method)
	if err != nil {
		log.Error("Unable to pack tx for get filteredValidatorsDepositMap", "error", err)
		return nil, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)

	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		log.Trace("Call", "err", err)
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("GetTotalDepositedBalance result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err)
		return nil, err
	}

	return out, nil
}

func (p *ProofOfStake) DoesDepositorExist(address common.Address, blockHash common.Hash) (bool, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("GETH_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return false, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_DoesDepositorExist()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("DoesDepositorExist abi error", "err", err)
		return false, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, address)
	if err != nil {
		log.Error("Unable to pack tx for get depositor exist", "error", err)
		return false, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		//Gas:  &gas,
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		return false, err
	}
	if len(result) == 0 {
		return false, errors.New("no depositor exist found")
	}

	var ret0 bool
	out := ret0
	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		return false, err
	}

	return out, nil
}

func (p *ProofOfStake) DidDepositorEverExists(address common.Address, blockHash common.Hash) (bool, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("GETH_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return false, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_DidDepositorEverExist()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("DidDepositorEverExists abi error", "err", err)
		return false, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, address)
	if err != nil {
		log.Error("Unable to pack tx for get depositor ever exists", "error", err)
		return false, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		//Gas:  &gas,
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		return false, err
	}
	if len(result) == 0 {
		return false, errors.New("no depositor ever exists found")
	}

	var ret0 bool
	out := ret0
	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		return false, err
	}

	return out, nil
}

func (p *ProofOfStake) DoesValidatorExist(address common.Address, blockHash common.Hash) (bool, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("GETH_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return false, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_DoesValidatorExist()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("DoesValidatorExist abi error", "err", err)
		return false, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, address)
	if err != nil {
		log.Error("Unable to pack tx for get validator exist", "error", err)
		return false, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		//Gas:  &gas,
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		return false, err
	}
	if len(result) == 0 {
		return false, errors.New("no validator exist found")
	}

	var ret0 bool
	out := ret0

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		return false, err
	}

	return out, nil
}

func (p *ProofOfStake) DidValidatorEverExists(address common.Address, blockHash common.Hash) (bool, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("GETH_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return false, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_DidValidatorEverExist()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("DidValidatorEverExists abi error", "err", err)
		return false, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, address)
	if err != nil {
		log.Error("Unable to pack tx for get validator ever exists", "error", err)
		return false, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)
	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		//Gas:  &gas,
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		return false, err
	}
	if len(result) == 0 {
		return false, errors.New("no validator ever exists found")
	}

	var ret0 bool
	out := ret0
	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		return false, err
	}

	return out, nil
}

func encodeCall(abi *abi.ABI, method string, args ...interface{}) ([]byte, error) {
	return abi.Pack(method, args...)
}

func (p *ProofOfStake) AddDepositorSlashing(blockHash common.Hash,
	depositor common.Address, slashedAmount *big.Int,
	state *state.StateDB, header *types.Header) (*big.Int, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("DP_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}

	method := staking.GetContract_Method_AddDepositorSlashing()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("AddDepositorSlashing abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := encodeCall(&abiData, method, depositor, slashedAmount)
	if err != nil {
		log.Error("Unable to pack AddDepositorSlashing", "error", err)
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
		return nil, errors.New("AddDepositorSlashing result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func (p *ProofOfStake) AddDepositorReward(blockHash common.Hash,
	depositor common.Address, rewardAmount *big.Int,
	state *state.StateDB, header *types.Header) (*big.Int, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("DP_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}

	method := staking.GetContract_Method_AddDepositorReward()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("AddDepositorReward abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := encodeCall(&abiData, method, depositor, rewardAmount)
	if err != nil {
		log.Error("Unable to pack AddDepositorReward", "error", err)
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
		return nil, errors.New("AddDepositorReward result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func (p *ProofOfStake) IsValidatorPaused(validator common.Address, blockHash common.Hash) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_IsValidationPaused()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("IsValidatorPaused abi error", "err", err)
		return false, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, validator)
	if err != nil {
		log.Error("IsValidatorPaused Unable to pack", "error", err)
		return false, err
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
		return false, err
	}
	if len(result) == 0 {
		return false, errors.New("IsValidatorPaused result is 0")
	}

	var out bool

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Debug("IsValidatorPaused UnpackIntoInterface", "err", err, "validator", validator)
		return false, err
	}

	return out, nil
}

func (p *ProofOfStake) GetBalanceOfDepositor(depositor common.Address, blockHash common.Hash) (*big.Int, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("DP_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_GetBalanceOfDepositor()

	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("GetBalanceOfDepositor abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, depositor)
	if err != nil {
		log.Error("Unable to pack tx for GetBalanceOfDepositor", "error", err)
		return nil, err
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
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("GetBalanceOfDepositor result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Debug("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func (p *ProofOfStake) GetDepositorRewards(depositor common.Address, blockHash common.Hash) (*big.Int, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("DP_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_GetDepositorRewards()

	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("GetDepositorRewards abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, depositor)
	if err != nil {
		log.Error("Unable to pack tx for GetDepositorRewards", "error", err)
		return nil, err
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
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("GetDepositorRewards result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Debug("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func (p *ProofOfStake) GetDepositorSlashings(depositor common.Address, blockHash common.Hash) (*big.Int, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("DP_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_GetDepositorSlashings()

	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("GetDepositorSlashings abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, depositor)
	if err != nil {
		log.Error("Unable to pack tx for GetDepositorSlashings", "error", err)
		return nil, err
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
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("GetDepositorSlashings result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Debug("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func (p *ProofOfStake) GetWithdrawalBlock(depositor common.Address, blockHash common.Hash) (*big.Int, error) {
	err := staking.IsStakingContract()
	if err != nil {
		log.Warn("DP_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_GetWithdrawalBlock()

	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("GetWithdrawalBlock abi error", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())

	// call
	data, err := abiData.Pack(method, depositor)
	if err != nil {
		log.Error("Unable to pack tx for GetWithdrawalBlock", "error", err)
		return nil, err
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
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("GetWithdrawalBlock result is 0")
	}

	var out *big.Int

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Debug("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func (p *ProofOfStake) ListValidators(blockHash common.Hash) ([]*ValidatorDetails, error) {
	depositorCount, err := p.GetDepositorCount(blockHash)
	if err != nil {
		return nil, err
	} else {
		log.Debug("depositorCount", "depositorCount", depositorCount)
	}

	err = staking.IsStakingContract()
	if err != nil {
		log.Warn("GETH_STAKING_CONTRACT_ADDRESS: Contract1 address is empty")
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	method := staking.GetContract_Method_ListValidators()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		log.Error("GetValidators error getting abidata", "err", err)
		return nil, err
	}
	contractAddress := common.HexToAddress(staking.GetStakingContract_Address_String())
	// call
	data, err := abiData.Pack(method)
	if err != nil {
		log.Error("Unable to pack tx for get filteredValidatorsDepositMap", "error", err)
		return nil, err
	}
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(blockHash, false)

	msgData := (hexutil.Bytes)(data)

	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		To:   &contractAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		log.Debug("result 0 length")
		return nil, nil
	}
	_, err = abiData.Unpack(method, result)
	if err != nil {
		log.Error("Unpack", "err", err)
		return nil, err
	}

	var (
		ret0 = new([]common.Address)
	)
	out := ret0

	if err := abiData.UnpackIntoInterface(out, method, result); err != nil {
		log.Info("UnpackIntoInterface error")
		return nil, err
	}

	var validatorList []*ValidatorDetails
	for _, val := range *out {
		if val.IsEqualTo(ZERO_ADDRESS) {
			return nil, errors.New("invalid validator")
		}
		log.Debug("GetValidators Validator", "val", val)
	}

	for _, val := range *out {
		isPaused, err := p.IsValidatorPaused(val, blockHash)
		if err != nil {
			log.Debug("IsValidatorPaused failed", "err", err)
			return nil, err
		}

		depositor, err := p.GetDepositorOfValidator(val, blockHash)
		if err != nil {
			log.Debug("GetDepositorOfValidator failed", "err", err)
			continue
		}

		if depositor.IsEqualTo(ZERO_ADDRESS) {
			return nil, errors.New("invalid depositor")
		}

		balance, err := p.GetBalanceOfDepositor(depositor, blockHash)
		if err != nil {
			log.Debug("GetBalanceOfDepositor failed", "err", err)
			continue
		}

		netBalance, err := p.GetNetBalanceOfDepositor(depositor, blockHash)
		if err != nil {
			log.Debug("GetNetBalanceOfDepositor failed", "err", err)
			continue
		}

		depositorRewards, err := p.GetDepositorRewards(depositor, blockHash)
		if err != nil {
			log.Debug("GetDepositorRewards failed", "err", err)
			continue
		}

		depositorSlashings, err := p.GetDepositorSlashings(depositor, blockHash)
		if err != nil {
			log.Debug("GetDepositorSlashings failed", "err", err)
			continue
		}

		withdrawalBlock, err := p.GetWithdrawalBlock(depositor, blockHash)
		if err != nil {
			log.Debug("GetWithdrawalBlock failed", "err", err)
			continue
		}

		validatorDetails := &ValidatorDetails{
			Validator:          val,
			Depositor:          depositor,
			Balance:            balance,
			NetBalance:         netBalance,
			BlockRewards:       depositorRewards,
			Slashings:          depositorSlashings,
			IsValidationPaused: isPaused,
			WithdrawalBlock:    withdrawalBlock,
		}

		validatorList = append(validatorList, validatorDetails)
	}

	return validatorList, nil
}
