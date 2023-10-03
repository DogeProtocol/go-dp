// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package staking

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/DogeProtocol/dp"
	"github.com/DogeProtocol/dp/accounts/abi"
	"github.com/DogeProtocol/dp/accounts/abi/bind"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// StakingMetaData contains all meta data concerning the Staking contract.
var StakingMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldValidatorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newValidatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnChangeValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnCompleteWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnIncreaseDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnInitiateWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnNewDeposit\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newValidatorAddress\",\"type\":\"address\"}],\"name\":\"changeValidator\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeWithdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getBalanceOfDepositor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDepositorCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"getDepositorOfValidator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalDepositedBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getValidatorOfDepositor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"increaseDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initiateWithdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"listValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"newDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60806040526000600255600060035534801561001a57600080fd5b506110a78061002a6000396000f3fe60806040526004361061009c5760003560e01c8063a7113fee11610064578063a7113fee1461016c578063b51d1d4f146101a9578063e03ff7cb146101c0578063f17bb462146101d7578063f6abfc7614610202578063ff9205ab1461021e5761009c565b806305b050de146100a157806368d4e544146100ab5780636d727bd0146100d6578063731f750d1461011357806377c06fdc1461012f575b600080fd5b6100a9610249565b005b3480156100b757600080fd5b506100c061024b565b6040516100cd9190610e58565b60405180910390f35b3480156100e257600080fd5b506100fd60048036038101906100f89190610b65565b6102d9565b60405161010a9190610e3d565b60405180910390f35b61012d60048036038101906101289190610b65565b610347565b005b34801561013b57600080fd5b5061015660048036038101906101519190610b65565b610a62565b6040516101639190610f5a565b60405180910390f35b34801561017857600080fd5b50610193600480360381019061018e9190610b65565b610aab565b6040516101a09190610e3d565b60405180910390f35b3480156101b557600080fd5b506101be610b19565b005b3480156101cc57600080fd5b506101d5610b1b565b005b3480156101e357600080fd5b506101ec610b1d565b6040516101f99190610f5a565b60405180910390f35b61021c60048036038101906102179190610b65565b610b27565b005b34801561022a57600080fd5b50610233610b2a565b6040516102409190610f5a565b60405180910390f35b565b606060008054806020026020016040519081016040528092919081815260200182805480156102cf57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610285575b5050505050905090565b600080600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905080915050919050565b600033905060003490506acecb8f27f4200f3a00000081101561039f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161039690610efa565b60405180910390fd5b8273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561040e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161040590610eba565b60405180910390fd5b60001515600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515146104a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161049890610f3a565b60405180910390fd5b60001515600660008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514610534576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161052b90610e9a565b60405180910390fd5b60008373ffffffffffffffffffffffffffffffffffffffff1631905060008114610593576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161058a90610e7a565b60405180910390fd5b60001515600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514610626576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161061d90610f1a565b60405180910390fd5b60001515600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515146106b9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106b090610eda565b60405180910390fd5b6000849080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555061073182600254610b3490919063ffffffff16565b60028190555061074d6001600354610b3490919063ffffffff16565b60038190555081600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506001600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506001600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506001600660008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506001600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555082600860008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555083600960008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508373ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fbe02029a5af0c964ebee7370f030cf18a026aae3a5d66f8107aee23f226d9ada844342604051610a5493929190610f75565b60405180910390a350505050565b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600080600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905080915050919050565b565b565b6000600354905090565b50565b6000600254905090565b600080828401905083811015610b4657fe5b8091505092915050565b600081359050610b5f81611032565b92915050565b600060208284031215610b7757600080fd5b6000610b8584828501610b50565b91505092915050565b6000610b9a8383610ba6565b60208301905092915050565b610baf81610ff6565b82525050565b610bbe81610ff6565b82525050565b6000610bcf82610fbc565b610bd98185610fd4565b9350610be483610fac565b8060005b83811015610c15578151610bfc8882610b8e565b9750610c0783610fc7565b925050600181019050610be8565b5085935050505092915050565b6000610c2f602083610fe5565b91507f76616c696461746f722062616c616e63652073686f756c64206265207a65726f6000830152602082019050919050565b6000610c6f601683610fe5565b91507f56616c696461746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b6000610caf603583610fe5565b91507f4465706f7369746f7220616464726573732063616e6e6f742062652073616d6560008301527f2061732056616c696461746f72206164647265737300000000000000000000006020830152604082019050919050565b6000610d15601683610fe5565b91507f4465706f7369746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b6000610d55602b83610fe5565b91507f4465706f73697420616d6f756e742062656c6f77206d696e696d756d2064657060008301527f6f73697420616d6f756e740000000000000000000000000000000000000000006020830152604082019050919050565b6000610dbb601883610fe5565b91507f4465706f7369746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b6000610dfb601883610fe5565b91507f56616c696461746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b610e3781611028565b82525050565b6000602082019050610e526000830184610bb5565b92915050565b60006020820190508181036000830152610e728184610bc4565b905092915050565b60006020820190508181036000830152610e9381610c22565b9050919050565b60006020820190508181036000830152610eb381610c62565b9050919050565b60006020820190508181036000830152610ed381610ca2565b9050919050565b60006020820190508181036000830152610ef381610d08565b9050919050565b60006020820190508181036000830152610f1381610d48565b9050919050565b60006020820190508181036000830152610f3381610dae565b9050919050565b60006020820190508181036000830152610f5381610dee565b9050919050565b6000602082019050610f6f6000830184610e2e565b92915050565b6000606082019050610f8a6000830186610e2e565b610f976020830185610e2e565b610fa46040830184610e2e565b949350505050565b6000819050602082019050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b600061100182611008565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b61103b81610ff6565b811461104657600080fd5b5056fea2646970667358221220121b3e8dfac96e3ea722773d264f9db17c677be6f59a50506f5c598fb4dc3b3d64736f6c63782a302e372e362d646576656c6f702e323032332e352e312b636f6d6d69742e37333338323935662e6d6f64005b",
}

// StakingABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingMetaData.ABI instead.
var StakingABI = StakingMetaData.ABI

// StakingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingMetaData.Bin instead.
var StakingBin = StakingMetaData.Bin

// DeployStaking deploys a new Ethereum contract, binding an instance of Staking to it.
func DeployStaking(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Staking, error) {
	parsed, err := StakingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// Staking is an auto generated Go binding around an Ethereum contract.
type Staking struct {
	StakingCaller     // Read-only binding to the contract
	StakingTransactor // Write-only binding to the contract
	StakingFilterer   // Log filterer for contract events
}

// StakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingSession struct {
	Contract     *Staking          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingCallerSession struct {
	Contract *StakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingTransactorSession struct {
	Contract     *StakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingRaw struct {
	Contract *Staking // Generic contract binding to access the raw methods on
}

// StakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingCallerRaw struct {
	Contract *StakingCaller // Generic read-only contract binding to access the raw methods on
}

// StakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingTransactorRaw struct {
	Contract *StakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStaking creates a new instance of Staking, bound to a specific deployed contract.
func NewStaking(address common.Address, backend bind.ContractBackend) (*Staking, error) {
	contract, err := bindStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// NewStakingCaller creates a new read-only instance of Staking, bound to a specific deployed contract.
func NewStakingCaller(address common.Address, caller bind.ContractCaller) (*StakingCaller, error) {
	contract, err := bindStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingCaller{contract: contract}, nil
}

// NewStakingTransactor creates a new write-only instance of Staking, bound to a specific deployed contract.
func NewStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingTransactor, error) {
	contract, err := bindStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTransactor{contract: contract}, nil
}

// NewStakingFilterer creates a new log filterer instance of Staking, bound to a specific deployed contract.
func NewStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingFilterer, error) {
	contract, err := bindStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingFilterer{contract: contract}, nil
}

// bindStaking binds a generic wrapper to an already deployed contract.
func bindStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.StakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transact(opts, method, params...)
}

// GetBalanceOfDepositor is a free data retrieval call binding the contract method 0x77c06fdc.
//
// Solidity: function getBalanceOfDepositor(address depositorAddress) view returns(uint256)
func (_Staking *StakingCaller) GetBalanceOfDepositor(opts *bind.CallOpts, depositorAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getBalanceOfDepositor", depositorAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalanceOfDepositor is a free data retrieval call binding the contract method 0x77c06fdc.
//
// Solidity: function getBalanceOfDepositor(address depositorAddress) view returns(uint256)
func (_Staking *StakingSession) GetBalanceOfDepositor(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetBalanceOfDepositor(&_Staking.CallOpts, depositorAddress)
}

// GetBalanceOfDepositor is a free data retrieval call binding the contract method 0x77c06fdc.
//
// Solidity: function getBalanceOfDepositor(address depositorAddress) view returns(uint256)
func (_Staking *StakingCallerSession) GetBalanceOfDepositor(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetBalanceOfDepositor(&_Staking.CallOpts, depositorAddress)
}

// GetDepositorCount is a free data retrieval call binding the contract method 0xf17bb462.
//
// Solidity: function getDepositorCount() view returns(uint256)
func (_Staking *StakingCaller) GetDepositorCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getDepositorCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDepositorCount is a free data retrieval call binding the contract method 0xf17bb462.
//
// Solidity: function getDepositorCount() view returns(uint256)
func (_Staking *StakingSession) GetDepositorCount() (*big.Int, error) {
	return _Staking.Contract.GetDepositorCount(&_Staking.CallOpts)
}

// GetDepositorCount is a free data retrieval call binding the contract method 0xf17bb462.
//
// Solidity: function getDepositorCount() view returns(uint256)
func (_Staking *StakingCallerSession) GetDepositorCount() (*big.Int, error) {
	return _Staking.Contract.GetDepositorCount(&_Staking.CallOpts)
}

// GetDepositorOfValidator is a free data retrieval call binding the contract method 0x6d727bd0.
//
// Solidity: function getDepositorOfValidator(address validatorAddress) view returns(address)
func (_Staking *StakingCaller) GetDepositorOfValidator(opts *bind.CallOpts, validatorAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getDepositorOfValidator", validatorAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDepositorOfValidator is a free data retrieval call binding the contract method 0x6d727bd0.
//
// Solidity: function getDepositorOfValidator(address validatorAddress) view returns(address)
func (_Staking *StakingSession) GetDepositorOfValidator(validatorAddress common.Address) (common.Address, error) {
	return _Staking.Contract.GetDepositorOfValidator(&_Staking.CallOpts, validatorAddress)
}

// GetDepositorOfValidator is a free data retrieval call binding the contract method 0x6d727bd0.
//
// Solidity: function getDepositorOfValidator(address validatorAddress) view returns(address)
func (_Staking *StakingCallerSession) GetDepositorOfValidator(validatorAddress common.Address) (common.Address, error) {
	return _Staking.Contract.GetDepositorOfValidator(&_Staking.CallOpts, validatorAddress)
}

// GetTotalDepositedBalance is a free data retrieval call binding the contract method 0xff9205ab.
//
// Solidity: function getTotalDepositedBalance() view returns(uint256)
func (_Staking *StakingCaller) GetTotalDepositedBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getTotalDepositedBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalDepositedBalance is a free data retrieval call binding the contract method 0xff9205ab.
//
// Solidity: function getTotalDepositedBalance() view returns(uint256)
func (_Staking *StakingSession) GetTotalDepositedBalance() (*big.Int, error) {
	return _Staking.Contract.GetTotalDepositedBalance(&_Staking.CallOpts)
}

// GetTotalDepositedBalance is a free data retrieval call binding the contract method 0xff9205ab.
//
// Solidity: function getTotalDepositedBalance() view returns(uint256)
func (_Staking *StakingCallerSession) GetTotalDepositedBalance() (*big.Int, error) {
	return _Staking.Contract.GetTotalDepositedBalance(&_Staking.CallOpts)
}

// GetValidatorOfDepositor is a free data retrieval call binding the contract method 0xa7113fee.
//
// Solidity: function getValidatorOfDepositor(address depositorAddress) view returns(address)
func (_Staking *StakingCaller) GetValidatorOfDepositor(opts *bind.CallOpts, depositorAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getValidatorOfDepositor", depositorAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetValidatorOfDepositor is a free data retrieval call binding the contract method 0xa7113fee.
//
// Solidity: function getValidatorOfDepositor(address depositorAddress) view returns(address)
func (_Staking *StakingSession) GetValidatorOfDepositor(depositorAddress common.Address) (common.Address, error) {
	return _Staking.Contract.GetValidatorOfDepositor(&_Staking.CallOpts, depositorAddress)
}

// GetValidatorOfDepositor is a free data retrieval call binding the contract method 0xa7113fee.
//
// Solidity: function getValidatorOfDepositor(address depositorAddress) view returns(address)
func (_Staking *StakingCallerSession) GetValidatorOfDepositor(depositorAddress common.Address) (common.Address, error) {
	return _Staking.Contract.GetValidatorOfDepositor(&_Staking.CallOpts, depositorAddress)
}

// ListValidators is a free data retrieval call binding the contract method 0x68d4e544.
//
// Solidity: function listValidators() view returns(address[])
func (_Staking *StakingCaller) ListValidators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "listValidators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ListValidators is a free data retrieval call binding the contract method 0x68d4e544.
//
// Solidity: function listValidators() view returns(address[])
func (_Staking *StakingSession) ListValidators() ([]common.Address, error) {
	return _Staking.Contract.ListValidators(&_Staking.CallOpts)
}

// ListValidators is a free data retrieval call binding the contract method 0x68d4e544.
//
// Solidity: function listValidators() view returns(address[])
func (_Staking *StakingCallerSession) ListValidators() ([]common.Address, error) {
	return _Staking.Contract.ListValidators(&_Staking.CallOpts)
}

// ChangeValidator is a paid mutator transaction binding the contract method 0xf6abfc76.
//
// Solidity: function changeValidator(address newValidatorAddress) payable returns()
func (_Staking *StakingTransactor) ChangeValidator(opts *bind.TransactOpts, newValidatorAddress common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "changeValidator", newValidatorAddress)
}

// ChangeValidator is a paid mutator transaction binding the contract method 0xf6abfc76.
//
// Solidity: function changeValidator(address newValidatorAddress) payable returns()
func (_Staking *StakingSession) ChangeValidator(newValidatorAddress common.Address) (*types.Transaction, error) {
	return _Staking.Contract.ChangeValidator(&_Staking.TransactOpts, newValidatorAddress)
}

// ChangeValidator is a paid mutator transaction binding the contract method 0xf6abfc76.
//
// Solidity: function changeValidator(address newValidatorAddress) payable returns()
func (_Staking *StakingTransactorSession) ChangeValidator(newValidatorAddress common.Address) (*types.Transaction, error) {
	return _Staking.Contract.ChangeValidator(&_Staking.TransactOpts, newValidatorAddress)
}

// CompleteWithdrawal is a paid mutator transaction binding the contract method 0xe03ff7cb.
//
// Solidity: function completeWithdrawal() returns()
func (_Staking *StakingTransactor) CompleteWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "completeWithdrawal")
}

// CompleteWithdrawal is a paid mutator transaction binding the contract method 0xe03ff7cb.
//
// Solidity: function completeWithdrawal() returns()
func (_Staking *StakingSession) CompleteWithdrawal() (*types.Transaction, error) {
	return _Staking.Contract.CompleteWithdrawal(&_Staking.TransactOpts)
}

// CompleteWithdrawal is a paid mutator transaction binding the contract method 0xe03ff7cb.
//
// Solidity: function completeWithdrawal() returns()
func (_Staking *StakingTransactorSession) CompleteWithdrawal() (*types.Transaction, error) {
	return _Staking.Contract.CompleteWithdrawal(&_Staking.TransactOpts)
}

// IncreaseDeposit is a paid mutator transaction binding the contract method 0x05b050de.
//
// Solidity: function increaseDeposit() payable returns()
func (_Staking *StakingTransactor) IncreaseDeposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "increaseDeposit")
}

// IncreaseDeposit is a paid mutator transaction binding the contract method 0x05b050de.
//
// Solidity: function increaseDeposit() payable returns()
func (_Staking *StakingSession) IncreaseDeposit() (*types.Transaction, error) {
	return _Staking.Contract.IncreaseDeposit(&_Staking.TransactOpts)
}

// IncreaseDeposit is a paid mutator transaction binding the contract method 0x05b050de.
//
// Solidity: function increaseDeposit() payable returns()
func (_Staking *StakingTransactorSession) IncreaseDeposit() (*types.Transaction, error) {
	return _Staking.Contract.IncreaseDeposit(&_Staking.TransactOpts)
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0xb51d1d4f.
//
// Solidity: function initiateWithdrawal() returns()
func (_Staking *StakingTransactor) InitiateWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "initiateWithdrawal")
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0xb51d1d4f.
//
// Solidity: function initiateWithdrawal() returns()
func (_Staking *StakingSession) InitiateWithdrawal() (*types.Transaction, error) {
	return _Staking.Contract.InitiateWithdrawal(&_Staking.TransactOpts)
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0xb51d1d4f.
//
// Solidity: function initiateWithdrawal() returns()
func (_Staking *StakingTransactorSession) InitiateWithdrawal() (*types.Transaction, error) {
	return _Staking.Contract.InitiateWithdrawal(&_Staking.TransactOpts)
}

// NewDeposit is a paid mutator transaction binding the contract method 0x731f750d.
//
// Solidity: function newDeposit(address validatorAddress) payable returns()
func (_Staking *StakingTransactor) NewDeposit(opts *bind.TransactOpts, validatorAddress common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "newDeposit", validatorAddress)
}

// NewDeposit is a paid mutator transaction binding the contract method 0x731f750d.
//
// Solidity: function newDeposit(address validatorAddress) payable returns()
func (_Staking *StakingSession) NewDeposit(validatorAddress common.Address) (*types.Transaction, error) {
	return _Staking.Contract.NewDeposit(&_Staking.TransactOpts, validatorAddress)
}

// NewDeposit is a paid mutator transaction binding the contract method 0x731f750d.
//
// Solidity: function newDeposit(address validatorAddress) payable returns()
func (_Staking *StakingTransactorSession) NewDeposit(validatorAddress common.Address) (*types.Transaction, error) {
	return _Staking.Contract.NewDeposit(&_Staking.TransactOpts, validatorAddress)
}

// StakingOnChangeValidatorIterator is returned from FilterOnChangeValidator and is used to iterate over the raw logs and unpacked data for OnChangeValidator events raised by the Staking contract.
type StakingOnChangeValidatorIterator struct {
	Event *StakingOnChangeValidator // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingOnChangeValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOnChangeValidator)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingOnChangeValidator)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingOnChangeValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOnChangeValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOnChangeValidator represents a OnChangeValidator event raised by the Staking contract.
type StakingOnChangeValidator struct {
	DepositorAddress    common.Address
	OldValidatorAddress common.Address
	NewValidatorAddress common.Address
	BlockNumber         *big.Int
	BlockTime           *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterOnChangeValidator is a free log retrieval operation binding the contract event 0x18050f198ff9118bad8fd08ab78228d2d5787a07299c7a267dc5c3dedb0388f8.
//
// Solidity: event OnChangeValidator(address indexed depositorAddress, address indexed oldValidatorAddress, address indexed newValidatorAddress, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) FilterOnChangeValidator(opts *bind.FilterOpts, depositorAddress []common.Address, oldValidatorAddress []common.Address, newValidatorAddress []common.Address) (*StakingOnChangeValidatorIterator, error) {

	var depositorAddressRule []interface{}
	for _, depositorAddressItem := range depositorAddress {
		depositorAddressRule = append(depositorAddressRule, depositorAddressItem)
	}
	var oldValidatorAddressRule []interface{}
	for _, oldValidatorAddressItem := range oldValidatorAddress {
		oldValidatorAddressRule = append(oldValidatorAddressRule, oldValidatorAddressItem)
	}
	var newValidatorAddressRule []interface{}
	for _, newValidatorAddressItem := range newValidatorAddress {
		newValidatorAddressRule = append(newValidatorAddressRule, newValidatorAddressItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnChangeValidator", depositorAddressRule, oldValidatorAddressRule, newValidatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &StakingOnChangeValidatorIterator{contract: _Staking.contract, event: "OnChangeValidator", logs: logs, sub: sub}, nil
}

// WatchOnChangeValidator is a free log subscription operation binding the contract event 0x18050f198ff9118bad8fd08ab78228d2d5787a07299c7a267dc5c3dedb0388f8.
//
// Solidity: event OnChangeValidator(address indexed depositorAddress, address indexed oldValidatorAddress, address indexed newValidatorAddress, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) WatchOnChangeValidator(opts *bind.WatchOpts, sink chan<- *StakingOnChangeValidator, depositorAddress []common.Address, oldValidatorAddress []common.Address, newValidatorAddress []common.Address) (event.Subscription, error) {

	var depositorAddressRule []interface{}
	for _, depositorAddressItem := range depositorAddress {
		depositorAddressRule = append(depositorAddressRule, depositorAddressItem)
	}
	var oldValidatorAddressRule []interface{}
	for _, oldValidatorAddressItem := range oldValidatorAddress {
		oldValidatorAddressRule = append(oldValidatorAddressRule, oldValidatorAddressItem)
	}
	var newValidatorAddressRule []interface{}
	for _, newValidatorAddressItem := range newValidatorAddress {
		newValidatorAddressRule = append(newValidatorAddressRule, newValidatorAddressItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OnChangeValidator", depositorAddressRule, oldValidatorAddressRule, newValidatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOnChangeValidator)
				if err := _Staking.contract.UnpackLog(event, "OnChangeValidator", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOnChangeValidator is a log parse operation binding the contract event 0x18050f198ff9118bad8fd08ab78228d2d5787a07299c7a267dc5c3dedb0388f8.
//
// Solidity: event OnChangeValidator(address indexed depositorAddress, address indexed oldValidatorAddress, address indexed newValidatorAddress, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) ParseOnChangeValidator(log types.Log) (*StakingOnChangeValidator, error) {
	event := new(StakingOnChangeValidator)
	if err := _Staking.contract.UnpackLog(event, "OnChangeValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingOnCompleteWithdrawalIterator is returned from FilterOnCompleteWithdrawal and is used to iterate over the raw logs and unpacked data for OnCompleteWithdrawal events raised by the Staking contract.
type StakingOnCompleteWithdrawalIterator struct {
	Event *StakingOnCompleteWithdrawal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingOnCompleteWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOnCompleteWithdrawal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingOnCompleteWithdrawal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingOnCompleteWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOnCompleteWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOnCompleteWithdrawal represents a OnCompleteWithdrawal event raised by the Staking contract.
type StakingOnCompleteWithdrawal struct {
	DepositorAddress common.Address
	BlockNumber      *big.Int
	BlockTime        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnCompleteWithdrawal is a free log retrieval operation binding the contract event 0xfe14ddf0dbfd6c9bb2f9e52e3325fb42ee3789d7c599b0c846631c53a4600cc5.
//
// Solidity: event OnCompleteWithdrawal(address depositorAddress, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) FilterOnCompleteWithdrawal(opts *bind.FilterOpts) (*StakingOnCompleteWithdrawalIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnCompleteWithdrawal")
	if err != nil {
		return nil, err
	}
	return &StakingOnCompleteWithdrawalIterator{contract: _Staking.contract, event: "OnCompleteWithdrawal", logs: logs, sub: sub}, nil
}

// WatchOnCompleteWithdrawal is a free log subscription operation binding the contract event 0xfe14ddf0dbfd6c9bb2f9e52e3325fb42ee3789d7c599b0c846631c53a4600cc5.
//
// Solidity: event OnCompleteWithdrawal(address depositorAddress, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) WatchOnCompleteWithdrawal(opts *bind.WatchOpts, sink chan<- *StakingOnCompleteWithdrawal) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OnCompleteWithdrawal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOnCompleteWithdrawal)
				if err := _Staking.contract.UnpackLog(event, "OnCompleteWithdrawal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOnCompleteWithdrawal is a log parse operation binding the contract event 0xfe14ddf0dbfd6c9bb2f9e52e3325fb42ee3789d7c599b0c846631c53a4600cc5.
//
// Solidity: event OnCompleteWithdrawal(address depositorAddress, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) ParseOnCompleteWithdrawal(log types.Log) (*StakingOnCompleteWithdrawal, error) {
	event := new(StakingOnCompleteWithdrawal)
	if err := _Staking.contract.UnpackLog(event, "OnCompleteWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingOnIncreaseDepositIterator is returned from FilterOnIncreaseDeposit and is used to iterate over the raw logs and unpacked data for OnIncreaseDeposit events raised by the Staking contract.
type StakingOnIncreaseDepositIterator struct {
	Event *StakingOnIncreaseDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingOnIncreaseDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOnIncreaseDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingOnIncreaseDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingOnIncreaseDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOnIncreaseDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOnIncreaseDeposit represents a OnIncreaseDeposit event raised by the Staking contract.
type StakingOnIncreaseDeposit struct {
	DepositorAddress common.Address
	Amount           *big.Int
	BlockNumber      *big.Int
	BlockTime        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnIncreaseDeposit is a free log retrieval operation binding the contract event 0xa3800194a102007a5793b1a612edd6205cd2cbc182f2bd6af553519c20ab06fc.
//
// Solidity: event OnIncreaseDeposit(address indexed depositorAddress, uint256 amount, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) FilterOnIncreaseDeposit(opts *bind.FilterOpts, depositorAddress []common.Address) (*StakingOnIncreaseDepositIterator, error) {

	var depositorAddressRule []interface{}
	for _, depositorAddressItem := range depositorAddress {
		depositorAddressRule = append(depositorAddressRule, depositorAddressItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnIncreaseDeposit", depositorAddressRule)
	if err != nil {
		return nil, err
	}
	return &StakingOnIncreaseDepositIterator{contract: _Staking.contract, event: "OnIncreaseDeposit", logs: logs, sub: sub}, nil
}

// WatchOnIncreaseDeposit is a free log subscription operation binding the contract event 0xa3800194a102007a5793b1a612edd6205cd2cbc182f2bd6af553519c20ab06fc.
//
// Solidity: event OnIncreaseDeposit(address indexed depositorAddress, uint256 amount, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) WatchOnIncreaseDeposit(opts *bind.WatchOpts, sink chan<- *StakingOnIncreaseDeposit, depositorAddress []common.Address) (event.Subscription, error) {

	var depositorAddressRule []interface{}
	for _, depositorAddressItem := range depositorAddress {
		depositorAddressRule = append(depositorAddressRule, depositorAddressItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OnIncreaseDeposit", depositorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOnIncreaseDeposit)
				if err := _Staking.contract.UnpackLog(event, "OnIncreaseDeposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOnIncreaseDeposit is a log parse operation binding the contract event 0xa3800194a102007a5793b1a612edd6205cd2cbc182f2bd6af553519c20ab06fc.
//
// Solidity: event OnIncreaseDeposit(address indexed depositorAddress, uint256 amount, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) ParseOnIncreaseDeposit(log types.Log) (*StakingOnIncreaseDeposit, error) {
	event := new(StakingOnIncreaseDeposit)
	if err := _Staking.contract.UnpackLog(event, "OnIncreaseDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingOnInitiateWithdrawalIterator is returned from FilterOnInitiateWithdrawal and is used to iterate over the raw logs and unpacked data for OnInitiateWithdrawal events raised by the Staking contract.
type StakingOnInitiateWithdrawalIterator struct {
	Event *StakingOnInitiateWithdrawal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingOnInitiateWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOnInitiateWithdrawal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingOnInitiateWithdrawal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingOnInitiateWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOnInitiateWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOnInitiateWithdrawal represents a OnInitiateWithdrawal event raised by the Staking contract.
type StakingOnInitiateWithdrawal struct {
	DepositorAddress common.Address
	BlockNumber      *big.Int
	BlockTime        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnInitiateWithdrawal is a free log retrieval operation binding the contract event 0xa3360e02f2be7df4763d047bb98447e7d835ec965e7b88f697fba878f56447d1.
//
// Solidity: event OnInitiateWithdrawal(address depositorAddress, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) FilterOnInitiateWithdrawal(opts *bind.FilterOpts) (*StakingOnInitiateWithdrawalIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnInitiateWithdrawal")
	if err != nil {
		return nil, err
	}
	return &StakingOnInitiateWithdrawalIterator{contract: _Staking.contract, event: "OnInitiateWithdrawal", logs: logs, sub: sub}, nil
}

// WatchOnInitiateWithdrawal is a free log subscription operation binding the contract event 0xa3360e02f2be7df4763d047bb98447e7d835ec965e7b88f697fba878f56447d1.
//
// Solidity: event OnInitiateWithdrawal(address depositorAddress, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) WatchOnInitiateWithdrawal(opts *bind.WatchOpts, sink chan<- *StakingOnInitiateWithdrawal) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OnInitiateWithdrawal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOnInitiateWithdrawal)
				if err := _Staking.contract.UnpackLog(event, "OnInitiateWithdrawal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOnInitiateWithdrawal is a log parse operation binding the contract event 0xa3360e02f2be7df4763d047bb98447e7d835ec965e7b88f697fba878f56447d1.
//
// Solidity: event OnInitiateWithdrawal(address depositorAddress, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) ParseOnInitiateWithdrawal(log types.Log) (*StakingOnInitiateWithdrawal, error) {
	event := new(StakingOnInitiateWithdrawal)
	if err := _Staking.contract.UnpackLog(event, "OnInitiateWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingOnNewDepositIterator is returned from FilterOnNewDeposit and is used to iterate over the raw logs and unpacked data for OnNewDeposit events raised by the Staking contract.
type StakingOnNewDepositIterator struct {
	Event *StakingOnNewDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingOnNewDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOnNewDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingOnNewDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingOnNewDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOnNewDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOnNewDeposit represents a OnNewDeposit event raised by the Staking contract.
type StakingOnNewDeposit struct {
	DepositorAddress common.Address
	ValidatorAddress common.Address
	Amount           *big.Int
	BlockNumber      *big.Int
	BlockTime        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnNewDeposit is a free log retrieval operation binding the contract event 0xbe02029a5af0c964ebee7370f030cf18a026aae3a5d66f8107aee23f226d9ada.
//
// Solidity: event OnNewDeposit(address indexed depositorAddress, address indexed validatorAddress, uint256 amount, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) FilterOnNewDeposit(opts *bind.FilterOpts, depositorAddress []common.Address, validatorAddress []common.Address) (*StakingOnNewDepositIterator, error) {

	var depositorAddressRule []interface{}
	for _, depositorAddressItem := range depositorAddress {
		depositorAddressRule = append(depositorAddressRule, depositorAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnNewDeposit", depositorAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &StakingOnNewDepositIterator{contract: _Staking.contract, event: "OnNewDeposit", logs: logs, sub: sub}, nil
}

// WatchOnNewDeposit is a free log subscription operation binding the contract event 0xbe02029a5af0c964ebee7370f030cf18a026aae3a5d66f8107aee23f226d9ada.
//
// Solidity: event OnNewDeposit(address indexed depositorAddress, address indexed validatorAddress, uint256 amount, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) WatchOnNewDeposit(opts *bind.WatchOpts, sink chan<- *StakingOnNewDeposit, depositorAddress []common.Address, validatorAddress []common.Address) (event.Subscription, error) {

	var depositorAddressRule []interface{}
	for _, depositorAddressItem := range depositorAddress {
		depositorAddressRule = append(depositorAddressRule, depositorAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OnNewDeposit", depositorAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOnNewDeposit)
				if err := _Staking.contract.UnpackLog(event, "OnNewDeposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOnNewDeposit is a log parse operation binding the contract event 0xbe02029a5af0c964ebee7370f030cf18a026aae3a5d66f8107aee23f226d9ada.
//
// Solidity: event OnNewDeposit(address indexed depositorAddress, address indexed validatorAddress, uint256 amount, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) ParseOnNewDeposit(log types.Log) (*StakingOnNewDeposit, error) {
	event := new(StakingOnNewDeposit)
	if err := _Staking.contract.UnpackLog(event, "OnNewDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
