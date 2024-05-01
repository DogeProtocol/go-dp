package proofofstake

import (
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
	"github.com/DogeProtocol/dp/params"
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

func getNoGasEVM(data []byte, from common.Address, state *state.StateDB, header *types.Header) (*vm.EVM, func() error, error) {
	msgData := (hexutil.Bytes)(data)

	args := ethapi.TransactionArgs{
		From: &from,
		To:   &ContractAddress,
		Data: &msgData,
	}

	msg, err := args.ToMessage(math.MaxUint64)
	if err != nil {
		return nil, nil, err
	}

	vmError := func() error { return nil }
	vmConfig := &vm.Config{OverrideGasFailure: true}
	chainConfig := &params.ChainConfig{
		ChainID:        big.NewInt(1),
		HomesteadBlock: new(big.Int),
		EIP155Block:    new(big.Int),
		EIP150Block:    new(big.Int),
		EIP158Block:    big.NewInt(2),
	}
	engine := mockconsensus.New(chainConfig, nil, common.HexToHash(GENESIS_BLOCK_HASH))

	tcc := &TestChainContext{
		Eng: engine,
	}

	txContext := core.NewEVMTxContext(msg)
	context := core.NewEVMBlockContext(header, tcc, nil)
	return vm.NewEVM(context, txContext, state, chainConfig, *vmConfig), vmError, nil
}

// revertEr

func setupTest() {
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	statedb.CreateAccount(ContractAddress)
	statedb.SetCode(ContractAddress, common.FromHex(stakingv2.STAKING_RUNTIME_BIN))

	statedb.Finalise(true) // Push the state into the "original" slot

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

func Test_GetBalance(t *testing.T) {
	setupTest()
}
