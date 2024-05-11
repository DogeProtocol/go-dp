package proofofstake

import (
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/consensus"
	"github.com/DogeProtocol/dp/consensus/mockconsensus"
	"github.com/DogeProtocol/dp/core"
	"github.com/DogeProtocol/dp/core/rawdb"
	"github.com/DogeProtocol/dp/core/state"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/core/vm"
	"github.com/DogeProtocol/dp/internal/ethapi"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/params"
	"github.com/DogeProtocol/dp/systemcontracts/staking"
	"github.com/DogeProtocol/dp/systemcontracts/staking/stakingv2"
	"math"
	"math/big"
	"strconv"
	"testing"
)

const STAKING_CONTRACT_V2 = "0x0000000000000000000000000000000000000000000000000000000000001000"

var (
	ContractAddress = common.HexToAddress(STAKING_CONTRACT_V2)
)

type TestChainContext struct {
	Eng consensus.Engine
}

func (tcc *TestChainContext) Engine() consensus.Engine {
	return tcc.Eng
}

func (tcc *TestChainContext) GetHeader(lastKnownHash common.Hash, lastKnownNumber uint64) *types.Header {
	hash := common.BytesToHash([]byte(strconv.FormatUint(lastKnownNumber+1, 10)))
	blockNumber := big.NewInt(int64(lastKnownNumber + 1))

	header := &types.Header{
		MixDigest:   hash,
		ReceiptHash: hash,
		TxHash:      hash,
		Nonce:       types.BlockNonce{},
		Extra:       []byte{},
		Bloom:       types.Bloom{},
		GasUsed:     0,
		Coinbase:    common.Address{},
		GasLimit:    0,
		Time:        1337,
		ParentHash:  lastKnownHash,
		Root:        hash,
		Number:      blockNumber,
		Difficulty:  largeNumber(2),
	}

	return header
}

var chainConfig = &params.ChainConfig{
	ChainID:             big.NewInt(1),
	HomesteadBlock:      new(big.Int),
	EIP155Block:         new(big.Int),
	EIP150Block:         new(big.Int),
	EIP158Block:         new(big.Int),
	ByzantiumBlock:      new(big.Int),
	ConstantinopleBlock: new(big.Int),
	PetersburgBlock:     new(big.Int),
	IstanbulBlock:       new(big.Int),
	BerlinBlock:         new(big.Int),
	LondonBlock:         new(big.Int),
}

var engine = mockconsensus.New(chainConfig, nil, common.HexToHash(GENESIS_BLOCK_HASH))

var tcc = &TestChainContext{Eng: engine}

func execute(tcc *TestChainContext, data []byte, from common.Address, state *state.StateDB, header *types.Header, value *big.Int) (hexutil.Bytes, error) {
	msgData := (hexutil.Bytes)(data)

	args := ethapi.TransactionArgs{
		From:  &from,
		To:    &ContractAddress,
		Data:  &msgData,
		Value: (*hexutil.Big)(value),
	}

	msg, err := args.ToMessage(math.MaxUint64)
	if err != nil {
		return nil, err
	}

	vmError := func() error { return nil }
	vmConfig := &vm.Config{OverrideGasFailure: true}

	txContext := core.NewEVMTxContext(msg)
	context := core.NewEVMBlockContext(header, tcc, nil)
	evm := vm.NewEVM(context, txContext, state, chainConfig, *vmConfig)

	gp := new(core.GasPool).AddGas(math.MaxUint64)
	result, err := core.ApplyMessage(evm, msg, gp)
	if err != nil {
		return nil, err
	}

	if err = vmError(); err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("result is nil")
	}

	// If the result contains a revert reason, try to unpack and return it.
	if len(result.Revert()) > 0 {
		return nil, core.NewRevertError(result)
	}

	return result.Return(), result.Err
}

func newStakingStateDb() *state.StateDB {
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	statedb.CreateAccount(ContractAddress)
	statedb.SetCode(ContractAddress, common.FromHex(stakingv2.STAKING_RUNTIME_BIN))
	statedb.Finalise(true) // Push the state into the "original" slot

	return statedb
}

func NewDeposit(state *state.StateDB, depositor common.Address, validator common.Address, amount *big.Int) error {
	method := staking.GetContract_Method_NewDeposit()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("AddDeposit abi error", "err", err)
		return nil
	}
	// call
	data, err := encodeCall(&abiData, method, validator)
	if err != nil {
		log.Error("Unable to pack AddDeposit", "error", err)
		return nil
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	_, err = execute(tcc, data, depositor, state, header, amount)
	if err != nil {
		return err
	}

	return nil
}

func PauseValidation(state *state.StateDB, depositor common.Address) error {
	method := staking.GetContract_Method_PauseValidation()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("PauseValidation abi error", "err", err)
		return err
	}
	// call
	data, err := encodeCall(&abiData, method)
	if err != nil {
		log.Error("Unable to pack PauseValidation", "error", err)
		return err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return errors.New("PauseValidation result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err)
		return err
	}

	return nil
}

func ResumeValidation(state *state.StateDB, depositor common.Address) error {
	method := staking.GetContract_Method_ResumeValidation()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("ResumeValidation abi error", "err", err)
		return err
	}
	// call
	data, err := encodeCall(&abiData, method)
	if err != nil {
		log.Error("Unable to pack ResumeValidation", "error", err)
		return err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return errors.New("ResumeValidation result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err)
		return err
	}

	return nil
}

func CompleteWithdrawal(state *state.StateDB, depositor common.Address) error {
	method := staking.GetContract_Method_CompleteWithdrawal()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("CompleteWithdrawal abi error", "err", err)
		return err
	}
	// call
	data, err := encodeCall(&abiData, method)
	if err != nil {
		log.Error("Unable to pack CompleteWithdrawal", "error", err)
		return err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return errors.New("CompleteWithdrawal result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err)
		return err
	}

	return nil
}

func AddDepositorSlashing(state *state.StateDB, from common.Address, depositor common.Address, amount *big.Int) (*big.Int, error) {
	method := staking.GetContract_Method_AddDepositorSlashing()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("AddDepositorSlashing abi error", "err", err)
		return nil, err
	}
	// call
	data, err := encodeCall(&abiData, method, amount)
	if err != nil {
		log.Error("Unable to pack AddDepositorSlashing", "error", err)
		return nil, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, from, state, header, new(big.Int))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("AddDepositorSlashing result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func AddDepositorReward(state *state.StateDB, from common.Address, depositor common.Address, amount *big.Int) (*big.Int, error) {
	method := staking.GetContract_Method_AddDepositorReward()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("AddDepositorReward abi error", "err", err)
		return nil, err
	}
	// call
	data, err := encodeCall(&abiData, method, amount)
	if err != nil {
		log.Error("Unable to pack AddDepositorReward", "error", err)
		return nil, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, from, state, header, new(big.Int))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("AddDepositorReward result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func GetDepositorCount(state *state.StateDB) (*big.Int, error) {
	method := staking.GetContract_Method_GetDepositorCount()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("GetDepositorCount abi error", "err", err)
		return nil, err
	}
	// call
	data, err := encodeCall(&abiData, method)
	if err != nil {
		log.Error("Unable to pack GetDepositorCount", "error", err)
		return nil, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, ZERO_ADDRESS, state, header, new(big.Int))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("GetDepositorCount result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err)
		return nil, err
	}

	return out, nil
}

func GetTotalDepositedBalance(state *state.StateDB) (*big.Int, error) {
	method := staking.GetContract_Method_GetTotalDepositedBalance()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("GetTotalDepositedBalance abi error", "err", err)
		return nil, err
	}
	// call
	data, err := encodeCall(&abiData, method)
	if err != nil {
		log.Error("Unable to pack GetTotalDepositedBalance", "error", err)
		return nil, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, ZERO_ADDRESS, state, header, new(big.Int))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("GetTotalDepositedBalance result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err)
		return nil, err
	}

	return out, nil
}

func ListValidators(state *state.StateDB) ([]common.Address, error) {
	method := staking.GetContract_Method_ListValidators()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("ListValidators abi error", "err", err)
		return nil, err
	}
	// call
	data, err := encodeCall(&abiData, method)
	if err != nil {
		log.Error("Unable to pack ListValidators", "error", err)
		return nil, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, ZERO_ADDRESS, state, header, new(big.Int))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("ListValidators result is 0")
	}

	var (
		out = new([]common.Address)
	)

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err)
		return nil, err
	}

	return *out, nil
}

func GetDepositorOfValidator(state *state.StateDB, validator common.Address) (common.Address, error) {
	method := staking.GetContract_Method_GetDepositorOfValidator()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("GetDepositorOfValidator abi error", "err", err)
		return ZERO_ADDRESS, err
	}
	// call
	data, err := encodeCall(&abiData, method, validator)
	if err != nil {
		log.Error("Unable to pack GetDepositorOfValidator", "error", err)
		return ZERO_ADDRESS, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, validator, state, header, new(big.Int))
	if err != nil {
		return ZERO_ADDRESS, err
	}

	if len(result) == 0 {
		return ZERO_ADDRESS, errors.New("GetDepositorOfValidator result is 0")
	}

	var (
		out = new(common.Address)
	)

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "validator", validator)
		return ZERO_ADDRESS, err
	}

	return *out, nil
}

func GetValidatorOfDepositor(state *state.StateDB, depositor common.Address) (common.Address, error) {
	method := staking.GetContract_Method_GetValidatorOfDepositor()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("GetValidatorOfDepositor abi error", "err", err)
		return ZERO_ADDRESS, err
	}
	// call
	data, err := encodeCall(&abiData, method, depositor)
	if err != nil {
		log.Error("Unable to pack GetValidatorOfDepositor", "error", err)
		return ZERO_ADDRESS, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return ZERO_ADDRESS, err
	}

	if len(result) == 0 {
		return ZERO_ADDRESS, errors.New("GetValidatorOfDepositor result is 0")
	}

	var (
		out = new(common.Address)
	)

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return ZERO_ADDRESS, err
	}

	return *out, nil
}

func GetBalanceOfDepositor(state *state.StateDB, depositor common.Address) (*big.Int, error) {
	method := staking.GetContract_Method_GetBalanceOfDepositor()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("GetBalanceOfDepositor abi error", "err", err)
		return nil, err
	}
	// call
	data, err := encodeCall(&abiData, method, depositor)
	if err != nil {
		log.Error("Unable to pack GetBalanceOfDepositor", "error", err)
		return nil, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("GetBalanceOfDepositor result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func GetNetBalanceOfDepositor(state *state.StateDB, depositor common.Address) (*big.Int, error) {
	method := staking.GetContract_Method_GetNetBalanceOfDepositor()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("GetNetBalanceOfDepositor abi error", "err", err)
		return nil, err
	}
	// call
	data, err := encodeCall(&abiData, method, depositor)
	if err != nil {
		log.Error("Unable to pack GetNetBalanceOfDepositor", "error", err)
		return nil, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("GetNetBalanceOfDepositor result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func GetDepositorRewards(state *state.StateDB, depositor common.Address) (*big.Int, error) {
	method := staking.GetContract_Method_GetDepositorRewards()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("GetDepositorRewards abi error", "err", err)
		return nil, err
	}
	// call
	data, err := encodeCall(&abiData, method, depositor)
	if err != nil {
		log.Error("Unable to pack GetDepositorRewards", "error", err)
		return nil, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("GetDepositorRewards result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func GetDepositorSlashings(state *state.StateDB, depositor common.Address) (*big.Int, error) {
	method := staking.GetContract_Method_GetDepositorSlashings()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("GetDepositorSlashings abi error", "err", err)
		return nil, err
	}
	// call
	data, err := encodeCall(&abiData, method, depositor)
	if err != nil {
		log.Error("Unable to pack GetDepositorSlashings", "error", err)
		return nil, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("GetDepositorSlashings result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func GetWithdrawalBlock(state *state.StateDB, depositor common.Address) (*big.Int, error) {
	method := staking.GetContract_Method_GetWithdrawalBlock()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("GetWithdrawalBlock abi error", "err", err)
		return nil, err
	}
	// call
	data, err := encodeCall(&abiData, method, depositor)
	if err != nil {
		log.Error("Unable to pack GetWithdrawalBlock", "error", err)
		return nil, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("GetWithdrawalBlock result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return nil, err
	}

	return out, nil
}

func IsValidationPaused(state *state.StateDB, depositor common.Address) (bool, error) {
	var out bool

	method := staking.GetContract_Method_IsValidationPaused()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("IsValidationPaused abi error", "err", err)
		return out, err
	}
	// call
	data, err := encodeCall(&abiData, method, depositor)
	if err != nil {
		log.Error("Unable to pack IsValidationPaused", "error", err)
		return out, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return out, err
	}

	if len(result) == 0 {
		return out, errors.New("IsValidationPaused result is 0")
	}
	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return out, err
	}

	return out, nil
}

func DoesValidatorExist(state *state.StateDB, validator common.Address) (bool, error) {
	var out bool

	method := staking.GetContract_Method_DoesValidatorExist()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("DoesValidatorExist abi error", "err", err)
		return out, err
	}
	// call
	data, err := encodeCall(&abiData, method, validator)
	if err != nil {
		log.Error("Unable to pack DoesValidatorExist", "error", err)
		return out, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, validator, state, header, new(big.Int))
	if err != nil {
		return out, err
	}

	if len(result) == 0 {
		return out, errors.New("DoesValidatorExist result is 0")
	}
	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", validator)
		return out, err
	}

	return out, nil
}

func DidValidatorEverExist(state *state.StateDB, validator common.Address) (bool, error) {
	var out bool

	method := staking.GetContract_Method_DidValidatorEverExist()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("DidValidatorEverExist abi error", "err", err)
		return out, err
	}
	// call
	data, err := encodeCall(&abiData, method, validator)
	if err != nil {
		log.Error("Unable to pack DidValidatorEverExist", "error", err)
		return out, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, validator, state, header, new(big.Int))
	if err != nil {
		return out, err
	}

	if len(result) == 0 {
		return out, errors.New("DidValidatorEverExist result is 0")
	}
	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", validator)
		return out, err
	}

	return out, nil
}

func DoesDepositorExist(state *state.StateDB, depositor common.Address) (bool, error) {
	var out bool

	method := staking.GetContract_Method_DoesDepositorExist()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("DoesDepositorExist abi error", "err", err)
		return out, err
	}
	// call
	data, err := encodeCall(&abiData, method, depositor)
	if err != nil {
		log.Error("Unable to pack DoesDepositorExist", "error", err)
		return out, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return out, err
	}

	if len(result) == 0 {
		return out, errors.New("DoesDepositorExist result is 0")
	}
	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return out, err
	}

	return out, nil
}

func DidDepositorEverExist(state *state.StateDB, depositor common.Address) (bool, error) {
	var out bool

	method := staking.GetContract_Method_DidValidatorEverExist()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("DidDepositorEverExist abi error", "err", err)
		return out, err
	}
	// call
	data, err := encodeCall(&abiData, method, depositor)
	if err != nil {
		log.Error("Unable to pack DidDepositorEverExist", "error", err)
		return out, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return out, err
	}

	if len(result) == 0 {
		return out, errors.New("DidDepositorEverExist result is 0")
	}
	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return out, err
	}

	return out, nil
}

func ChangeValidator(state *state.StateDB, depositor common.Address, newValidatorAddress common.Address) error {
	method := staking.GetContract_Method_ChangeValidator()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("ChangeValidator abi error", "err", err)
		return err
	}
	// call
	data, err := encodeCall(&abiData, method, newValidatorAddress)
	if err != nil {
		log.Error("Unable to pack ChangeValidator", "error", err)
		return err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	_, err = execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return err
	}

	return nil
}

func InitiateChangeDepositor(state *state.StateDB, oldDepositorAddress common.Address, newDepositorAddress common.Address) error {
	method := staking.GetContract_Method_InitiateChangeDepositor()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("InitiateChangeDepositor abi error", "err", err)
		return err
	}
	// call
	data, err := encodeCall(&abiData, method, newDepositorAddress)
	if err != nil {
		log.Error("Unable to pack InitiateChangeDepositor", "error", err)
		return err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	_, err = execute(tcc, data, oldDepositorAddress, state, header, new(big.Int))
	if err != nil {
		return err
	}

	return nil
}

func CompleteChangeDepositor(state *state.StateDB, oldDepositorAddress common.Address, newDepositorAddress common.Address) error {
	method := staking.GetContract_Method_CompleteChangeDepositor()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("CompleteChangeDepositor abi error", "err", err)
		return err
	}
	// call
	data, err := encodeCall(&abiData, method, oldDepositorAddress)
	if err != nil {
		log.Error("Unable to pack CompleteChangeDepositor", "error", err)
		return err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	_, err = execute(tcc, data, newDepositorAddress, state, header, new(big.Int))
	if err != nil {
		return err
	}

	return nil
}

func IncreaseDeposit(state *state.StateDB, depositor common.Address, amount *big.Int) error {
	method := staking.GetContract_Method_IncreaseDeposit()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("IncreaseDeposit abi error", "err", err)
		return err
	}
	// call
	data, err := encodeCall(&abiData, method, amount)
	if err != nil {
		log.Error("Unable to pack IncreaseDeposit", "error", err)
		return err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	_, err = execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return err
	}

	return nil
}

func InitiatePartialWithdrawal(state *state.StateDB, depositor common.Address, amount *big.Int, currentBlockNumber uint64) error {
	method := staking.GetContract_Method_InitiatePartialWithdrawal()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("InitiatePartialWithdrawal abi error", "err", err)
		return err
	}
	// call
	data, err := encodeCall(&abiData, method, amount)
	if err != nil {
		log.Error("Unable to pack InitiatePartialWithdrawal", "error", err)
		return err
	}

	header := tcc.GetHeader(ZERO_HASH, currentBlockNumber)

	_, err = execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return err
	}

	return nil
}

func CompletePartialWithdrawal(state *state.StateDB, depositor common.Address) error {
	method := staking.GetContract_Method_CompletePartialWithdrawal()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("CompletePartialWithdrawal abi error", "err", err)
		return err
	}
	// call
	data, err := encodeCall(&abiData, method)
	if err != nil {
		log.Error("Unable to pack CompletePartialWithdrawal", "error", err)
		return err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	_, err = execute(tcc, data, depositor, state, header, new(big.Int))
	if err != nil {
		return err
	}

	return nil
}

func GetStakingDetails(state *state.StateDB, validator common.Address) (*ValidatorDetailsV2, error) {
	method := staking.GetContract_Method_GetStakingDetails()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("GetStakingDetails abi error", "err", err)
		return nil, err
	}
	// call
	data, err := encodeCall(&abiData, method, validator)
	if err != nil {
		log.Error("Unable to pack GetStakingDetails", "error", err)
		return nil, err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	result, err := execute(tcc, data, validator, state, header, new(big.Int))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("GetStakingDetails result is 0")
	}
	var out *ValidatorDetailsV2
	out = new(ValidatorDetailsV2)

	if err := abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Debug("UnpackIntoInterface", "err", err, "validator", validator)
		return nil, err
	}

	return out, nil
}

func SetNilBlock(state *state.StateDB, from common.Address, validator common.Address) error {
	method := staking.GetContract_Method_SetNilBlock()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("CompletePartialWithdrawal abi error", "err", err)
		return err
	}
	// call
	data, err := encodeCall(&abiData, method, validator)
	if err != nil {
		log.Error("Unable to pack CompletePartialWithdrawal", "error", err)
		return err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	_, err = execute(tcc, data, from, state, header, new(big.Int))
	if err != nil {
		return err
	}

	return nil
}

func ResetNilBlock(state *state.StateDB, from common.Address, validator common.Address) error {
	method := staking.GetContract_Method_ResetNilBlock()
	abiData, err := staking.GetStakingContractV2_ABI()
	if err != nil {
		log.Error("ResetNilBlock abi error", "err", err)
		return err
	}
	// call
	data, err := encodeCall(&abiData, method, validator)
	if err != nil {
		log.Error("Unable to pack ResetNilBlock", "error", err)
		return err
	}

	header := tcc.GetHeader(ZERO_HASH, uint64(1))

	_, err = execute(tcc, data, from, state, header, new(big.Int))
	if err != nil {
		return err
	}

	return nil
}

func TestStaking_Basic(t *testing.T) {
	depositor := common.RandomAddress()
	validator := common.RandomAddress()
	state := newStakingStateDb()

	balance := params.EtherToWei(big.NewInt(10000000))
	state.SetBalance(depositor, balance)
	//state.Finalise(true)

	depositAmount := MIN_VALIDATOR_DEPOSIT
	err := NewDeposit(state, depositor, validator, MIN_VALIDATOR_DEPOSIT)
	if err != nil {
		t.Fatal(err)
	}

	stakingBalance, err := GetBalanceOfDepositor(state, depositor)
	if err != nil {
		t.Fatal(err)
	}

	if stakingBalance.Cmp(depositAmount) != 0 {
		t.Fatalf("balanace compare failed")
	}

	stakingNetBalance, err := GetNetBalanceOfDepositor(state, depositor)
	if err != nil {
		t.Fatal(err)
	}

	if stakingNetBalance.Cmp(depositAmount) != 0 {
		t.Fatalf("net balanace compare failed")
	}

	totalDepositedBalance, err := GetTotalDepositedBalance(state)
	if err != nil {
		t.Fatal(err)
	}

	if totalDepositedBalance.Cmp(depositAmount) != 0 {
		t.Fatalf("totalDepositedBalance compare failed")
	}

	err = InitiatePartialWithdrawal(state, depositor, big.NewInt(1000), 10)
	if err != nil {
		t.Fatal(err)
	}

	stakingDetails, err := GetStakingDetails(state, validator)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("depositAmount", depositAmount, "stakingBalance", stakingBalance, "stakingNetBalance", stakingNetBalance, "totalDepositedBalance", totalDepositedBalance)

	fmt.Println("StakingDetails withdrawalblock", stakingDetails.WithdrawalBlock.Uint64())
	fmt.Println("StakingDetails withdrawalbamount", stakingDetails.WithdrawalAmount.Uint64())
}
