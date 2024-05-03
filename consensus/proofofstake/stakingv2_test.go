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
		Number:      largeNumber(2),
		Difficulty:  largeNumber(2),
	}

	return header
}

var chainConfig = &params.ChainConfig{
	ChainID:        big.NewInt(1),
	HomesteadBlock: new(big.Int),
	EIP155Block:    new(big.Int),
	EIP150Block:    new(big.Int),
	EIP158Block:    big.NewInt(2),
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

	fmt.Println("value", value)

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
	/*vmctx := vm.BlockContext{
		CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
		Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
	}
	_ = vm.NewEVM(vmctx, vm.TxContext{}, statedb, params.AllProofOfStakeProtocolChanges, vm.Config{OverrideGasFailure: true})*/

	/*_, gas, err := vmenv.Call(vm.AccountRef(common.Address{}), address, nil, tt.gaspool, new(big.Int))
	if err != tt.failure {
		t.Errorf("test %d: failure mismatch: have %v, want %v", i, err, tt.failure)
	}
	if used := tt.gaspool - gas; used != tt.used {
		t.Errorf("test %d: gas used mismatch: have %v, want %v", i, used, tt.used)
	}
	if refund := vmenv.StateDB.GetRefund(); refund != tt.refund {
		t.Errorf("test %d: gas refund mismatch: have %v, want %v", i, refund, tt.refund)
	}*/
}

func AddDeposit(state *state.StateDB, depositor common.Address, validator common.Address, amount *big.Int) error {
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

	result, err := execute(tcc, data, depositor, state, header, amount)
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return errors.New("AddDeposit result is 0")
	}

	var out *big.Int

	if err = abiData.UnpackIntoInterface(&out, method, result); err != nil {
		log.Trace("UnpackIntoInterface", "err", err, "depositor", depositor)
		return err
	}

	if out.Cmp(amount) != 0 {
		return errors.New("adddeposit return value failed")
	}

	return nil
}

func Test_AddDeposit_Basic(t *testing.T) {
	depositor := common.RandomAddress()
	validator := common.RandomAddress()
	state := newStakingStateDb()

	balance := params.EtherToWei(big.NewInt(10000000))
	state.SetBalance(depositor, balance)
	//state.Finalise(true)

	err := AddDeposit(state, depositor, validator, MIN_VALIDATOR_DEPOSIT)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetBalance(t *testing.T) {

}
