// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
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

// StakingContractAddress1MetaData contains all meta data concerning the StakingContractAddress1 contract.
var StakingContractAddress1MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"validatorId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnNewDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnWithdrawKey\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"depositBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"getDepositor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"listValidator\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"}],\"name\":\"newDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalDepositBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600080819055506000600181905550610f478061002f6000396000f3fe6080604052600436106100705760003560e01c806375697e661161004e57806375697e6614610106578063dfcd068f14610131578063e8c0a0df1461014d578063fba13bd01461017857610070565b80632dfdf0b5146100755780632e1a7d4d146100a05780636e2baf48146100c9575b600080fd5b34801561008157600080fd5b5061008a6101b5565b6040516100979190610d9d565b60405180910390f35b3480156100ac57600080fd5b506100c760048036038101906100c29190610a67565b6101be565b005b3480156100d557600080fd5b506100f060048036038101906100eb91906109f9565b610377565b6040516100fd9190610c6d565b60405180910390f35b34801561011257600080fd5b5061011b6103c7565b6040516101289190610ccd565b60405180910390f35b61014b60048036038101906101469190610a22565b610455565b005b34801561015957600080fd5b50610162610820565b60405161016f9190610d9d565b60405180910390f35b34801561018457600080fd5b5061019f600480360381019061019a91906109f9565b61082a565b6040516101ac9190610d9d565b60405180910390f35b60008054905090565b80600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015610240576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161023790610d5d565b60405180910390fd5b6102558160015461087390919063ffffffff16565b6001819055506102ad81600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461087390919063ffffffff16565b600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503373ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610336573d6000803e3d6000fd5b507f4d4666331ec61727075c5624fde25f5510c566e528d0565f2a2263a23b70d81a3382434260405161036c9493929190610c88565b60405180910390a150565b6000806103838361088a565b905060006004600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508092505050919050565b6060600680548060200260200160405190810160405280929190818152602001828054801561044b57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610401575b5050505050905090565b6000828290501161049b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161049290610d3d565b60405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff1660046000600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141561057c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161057390610d7d565b60405180910390fd5b61059260016000546108b190919063ffffffff16565b6000819055506105ad346001546108b190919063ffffffff16565b60018190555061060534600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546108b190919063ffffffff16565b600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555060008282600190809261065d93929190610e1e565b60405161066b929190610c54565b604051809103902090506000610680826108cd565b9050600061068d8261088a565b905084846003600084815260200190815260200160002091906106b19291906108da565b50336004600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506006829080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff16813373ffffffffffffffffffffffffffffffffffffffff167f9a1f4f083763f8508b19d4301c0110d2b47d99a8c5cf52c825c9e8cfea17f89c8888344342604051610811959493929190610cef565b60405180910390a45050505050565b6000600154905090565b6000600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b60008282111561087f57fe5b818303905092915050565b600060608273ffffffffffffffffffffffffffffffffffffffff16901b60001b9050919050565b6000808284019050838110156108c357fe5b8091505092915050565b60008160001c9050919050565b828054600181600116156101000203166002900490600052602060002090601f0160209004810192826109105760008555610957565b82601f1061092957803560ff1916838001178555610957565b82800160010185558215610957579182015b8281111561095657823582559160200191906001019061093b565b5b5090506109649190610968565b5090565b5b80821115610981576000816000905550600101610969565b5090565b60008135905061099481610ee3565b92915050565b60008083601f8401126109ac57600080fd5b8235905067ffffffffffffffff8111156109c557600080fd5b6020830191508360018202830111156109dd57600080fd5b9250929050565b6000813590506109f381610efa565b92915050565b600060208284031215610a0b57600080fd5b6000610a1984828501610985565b91505092915050565b60008060208385031215610a3557600080fd5b600083013567ffffffffffffffff811115610a4f57600080fd5b610a5b8582860161099a565b92509250509250929050565b600060208284031215610a7957600080fd5b6000610a87848285016109e4565b91505092915050565b6000610a9c8383610ab7565b60208301905092915050565b610ab181610e8d565b82525050565b610ac081610e51565b82525050565b610acf81610e51565b82525050565b6000610ae082610dc8565b610aea8185610de0565b9350610af583610db8565b8060005b83811015610b26578151610b0d8882610a90565b9750610b1883610dd3565b925050600181019050610af9565b5085935050505092915050565b6000610b3f8385610df1565b9350610b4c838584610ec3565b610b5583610ed2565b840190509392505050565b6000610b6c8385610e02565b9350610b79838584610ec3565b82840190509392505050565b6000610b92601583610e0d565b91507f5075626c6963206b657920697320696e76616c696400000000000000000000006000830152602082019050919050565b6000610bd2601283610e0d565b91507f496e73756666696369656e742066756e647300000000000000000000000000006000830152602082019050919050565b6000610c12601583610e0d565b91507f53656e64657220616c72656164792065786973747300000000000000000000006000830152602082019050919050565b610c4e81610e83565b82525050565b6000610c61828486610b60565b91508190509392505050565b6000602082019050610c826000830184610ac6565b92915050565b6000608082019050610c9d6000830187610aa8565b610caa6020830186610c45565b610cb76040830185610c45565b610cc46060830184610c45565b95945050505050565b60006020820190508181036000830152610ce78184610ad5565b905092915050565b60006080820190508181036000830152610d0a818789610b33565b9050610d196020830186610c45565b610d266040830185610c45565b610d336060830184610c45565b9695505050505050565b60006020820190508181036000830152610d5681610b85565b9050919050565b60006020820190508181036000830152610d7681610bc5565b9050919050565b60006020820190508181036000830152610d9681610c05565b9050919050565b6000602082019050610db26000830184610c45565b92915050565b6000819050602082019050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b60008085851115610e2e57600080fd5b83861115610e3b57600080fd5b6001850283019150848603905094509492505050565b6000610e5c82610e63565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b6000610e9882610e9f565b9050919050565b6000610eaa82610eb1565b9050919050565b6000610ebc82610e63565b9050919050565b82818337600083830152505050565b6000601f19601f8301169050919050565b610eec81610e51565b8114610ef757600080fd5b50565b610f0381610e83565b8114610f0e57600080fd5b5056fea2646970667358221220181b27743bf08caf1acd6da3c4bdbd6904a61bcc46be748c38e4a0ca4f2b4e5964736f6c63430007060033",
}

// StakingContractAddress1ABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingContractAddress1MetaData.ABI instead.
var StakingContractAddress1ABI = StakingContractAddress1MetaData.ABI

// StakingContractAddress1Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingContractAddress1MetaData.Bin instead.
var StakingContractAddress1Bin = StakingContractAddress1MetaData.Bin

// DeployStakingContractAddress1 deploys a new Ethereum contract, binding an instance of StakingContractAddress1 to it.
func DeployStakingContractAddress1(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StakingContractAddress1, error) {
	parsed, err := StakingContractAddress1MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingContractAddress1Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StakingContractAddress1{StakingContractAddress1Caller: StakingContractAddress1Caller{contract: contract}, StakingContractAddress1Transactor: StakingContractAddress1Transactor{contract: contract}, StakingContractAddress1Filterer: StakingContractAddress1Filterer{contract: contract}}, nil
}

// StakingContractAddress1 is an auto generated Go binding around an Ethereum contract.
type StakingContractAddress1 struct {
	StakingContractAddress1Caller     // Read-only binding to the contract
	StakingContractAddress1Transactor // Write-only binding to the contract
	StakingContractAddress1Filterer   // Log filterer for contract events
}

// StakingContractAddress1Caller is an auto generated read-only Go binding around an Ethereum contract.
type StakingContractAddress1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingContractAddress1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingContractAddress1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingContractAddress1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingContractAddress1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingContractAddress1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingContractAddress1Session struct {
	Contract     *StakingContractAddress1 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StakingContractAddress1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingContractAddress1CallerSession struct {
	Contract *StakingContractAddress1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// StakingContractAddress1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingContractAddress1TransactorSession struct {
	Contract     *StakingContractAddress1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// StakingContractAddress1Raw is an auto generated low-level Go binding around an Ethereum contract.
type StakingContractAddress1Raw struct {
	Contract *StakingContractAddress1 // Generic contract binding to access the raw methods on
}

// StakingContractAddress1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingContractAddress1CallerRaw struct {
	Contract *StakingContractAddress1Caller // Generic read-only contract binding to access the raw methods on
}

// StakingContractAddress1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingContractAddress1TransactorRaw struct {
	Contract *StakingContractAddress1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingContractAddress1 creates a new instance of StakingContractAddress1, bound to a specific deployed contract.
func NewStakingContractAddress1(address common.Address, backend bind.ContractBackend) (*StakingContractAddress1, error) {
	contract, err := bindStakingContractAddress1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingContractAddress1{StakingContractAddress1Caller: StakingContractAddress1Caller{contract: contract}, StakingContractAddress1Transactor: StakingContractAddress1Transactor{contract: contract}, StakingContractAddress1Filterer: StakingContractAddress1Filterer{contract: contract}}, nil
}

// NewStakingContractAddress1Caller creates a new read-only instance of StakingContractAddress1, bound to a specific deployed contract.
func NewStakingContractAddress1Caller(address common.Address, caller bind.ContractCaller) (*StakingContractAddress1Caller, error) {
	contract, err := bindStakingContractAddress1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingContractAddress1Caller{contract: contract}, nil
}

// NewStakingContractAddress1Transactor creates a new write-only instance of StakingContractAddress1, bound to a specific deployed contract.
func NewStakingContractAddress1Transactor(address common.Address, transactor bind.ContractTransactor) (*StakingContractAddress1Transactor, error) {
	contract, err := bindStakingContractAddress1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingContractAddress1Transactor{contract: contract}, nil
}

// NewStakingContractAddress1Filterer creates a new log filterer instance of StakingContractAddress1, bound to a specific deployed contract.
func NewStakingContractAddress1Filterer(address common.Address, filterer bind.ContractFilterer) (*StakingContractAddress1Filterer, error) {
	contract, err := bindStakingContractAddress1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingContractAddress1Filterer{contract: contract}, nil
}

// bindStakingContractAddress1 binds a generic wrapper to an already deployed contract.
func bindStakingContractAddress1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingContractAddress1ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingContractAddress1 *StakingContractAddress1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingContractAddress1.Contract.StakingContractAddress1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingContractAddress1 *StakingContractAddress1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingContractAddress1.Contract.StakingContractAddress1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingContractAddress1 *StakingContractAddress1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingContractAddress1.Contract.StakingContractAddress1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingContractAddress1 *StakingContractAddress1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingContractAddress1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingContractAddress1 *StakingContractAddress1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingContractAddress1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingContractAddress1 *StakingContractAddress1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingContractAddress1.Contract.contract.Transact(opts, method, params...)
}

// DepositBalanceOf is a free data retrieval call binding the contract method 0xfba13bd0.
//
// Solidity: function depositBalanceOf(address depositor) view returns(uint256)
func (_StakingContractAddress1 *StakingContractAddress1Caller) DepositBalanceOf(opts *bind.CallOpts, depositor common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingContractAddress1.contract.Call(opts, &out, "depositBalanceOf", depositor)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositBalanceOf is a free data retrieval call binding the contract method 0xfba13bd0.
//
// Solidity: function depositBalanceOf(address depositor) view returns(uint256)
func (_StakingContractAddress1 *StakingContractAddress1Session) DepositBalanceOf(depositor common.Address) (*big.Int, error) {
	return _StakingContractAddress1.Contract.DepositBalanceOf(&_StakingContractAddress1.CallOpts, depositor)
}

// DepositBalanceOf is a free data retrieval call binding the contract method 0xfba13bd0.
//
// Solidity: function depositBalanceOf(address depositor) view returns(uint256)
func (_StakingContractAddress1 *StakingContractAddress1CallerSession) DepositBalanceOf(depositor common.Address) (*big.Int, error) {
	return _StakingContractAddress1.Contract.DepositBalanceOf(&_StakingContractAddress1.CallOpts, depositor)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_StakingContractAddress1 *StakingContractAddress1Caller) DepositCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingContractAddress1.contract.Call(opts, &out, "depositCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_StakingContractAddress1 *StakingContractAddress1Session) DepositCount() (*big.Int, error) {
	return _StakingContractAddress1.Contract.DepositCount(&_StakingContractAddress1.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_StakingContractAddress1 *StakingContractAddress1CallerSession) DepositCount() (*big.Int, error) {
	return _StakingContractAddress1.Contract.DepositCount(&_StakingContractAddress1.CallOpts)
}

// GetDepositor is a free data retrieval call binding the contract method 0x6e2baf48.
//
// Solidity: function getDepositor(address validator) view returns(address)
func (_StakingContractAddress1 *StakingContractAddress1Caller) GetDepositor(opts *bind.CallOpts, validator common.Address) (common.Address, error) {
	var out []interface{}
	err := _StakingContractAddress1.contract.Call(opts, &out, "getDepositor", validator)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDepositor is a free data retrieval call binding the contract method 0x6e2baf48.
//
// Solidity: function getDepositor(address validator) view returns(address)
func (_StakingContractAddress1 *StakingContractAddress1Session) GetDepositor(validator common.Address) (common.Address, error) {
	return _StakingContractAddress1.Contract.GetDepositor(&_StakingContractAddress1.CallOpts, validator)
}

// GetDepositor is a free data retrieval call binding the contract method 0x6e2baf48.
//
// Solidity: function getDepositor(address validator) view returns(address)
func (_StakingContractAddress1 *StakingContractAddress1CallerSession) GetDepositor(validator common.Address) (common.Address, error) {
	return _StakingContractAddress1.Contract.GetDepositor(&_StakingContractAddress1.CallOpts, validator)
}

// ListValidator is a free data retrieval call binding the contract method 0x75697e66.
//
// Solidity: function listValidator() view returns(address[])
func (_StakingContractAddress1 *StakingContractAddress1Caller) ListValidator(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _StakingContractAddress1.contract.Call(opts, &out, "listValidator")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ListValidator is a free data retrieval call binding the contract method 0x75697e66.
//
// Solidity: function listValidator() view returns(address[])
func (_StakingContractAddress1 *StakingContractAddress1Session) ListValidator() ([]common.Address, error) {
	return _StakingContractAddress1.Contract.ListValidator(&_StakingContractAddress1.CallOpts)
}

// ListValidator is a free data retrieval call binding the contract method 0x75697e66.
//
// Solidity: function listValidator() view returns(address[])
func (_StakingContractAddress1 *StakingContractAddress1CallerSession) ListValidator() ([]common.Address, error) {
	return _StakingContractAddress1.Contract.ListValidator(&_StakingContractAddress1.CallOpts)
}

// TotalDepositBalance is a free data retrieval call binding the contract method 0xe8c0a0df.
//
// Solidity: function totalDepositBalance() view returns(uint256)
func (_StakingContractAddress1 *StakingContractAddress1Caller) TotalDepositBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingContractAddress1.contract.Call(opts, &out, "totalDepositBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalDepositBalance is a free data retrieval call binding the contract method 0xe8c0a0df.
//
// Solidity: function totalDepositBalance() view returns(uint256)
func (_StakingContractAddress1 *StakingContractAddress1Session) TotalDepositBalance() (*big.Int, error) {
	return _StakingContractAddress1.Contract.TotalDepositBalance(&_StakingContractAddress1.CallOpts)
}

// TotalDepositBalance is a free data retrieval call binding the contract method 0xe8c0a0df.
//
// Solidity: function totalDepositBalance() view returns(uint256)
func (_StakingContractAddress1 *StakingContractAddress1CallerSession) TotalDepositBalance() (*big.Int, error) {
	return _StakingContractAddress1.Contract.TotalDepositBalance(&_StakingContractAddress1.CallOpts)
}

// NewDeposit is a paid mutator transaction binding the contract method 0xdfcd068f.
//
// Solidity: function newDeposit(bytes pubkey) payable returns()
func (_StakingContractAddress1 *StakingContractAddress1Transactor) NewDeposit(opts *bind.TransactOpts, pubkey []byte) (*types.Transaction, error) {
	return _StakingContractAddress1.contract.Transact(opts, "newDeposit", pubkey)
}

// NewDeposit is a paid mutator transaction binding the contract method 0xdfcd068f.
//
// Solidity: function newDeposit(bytes pubkey) payable returns()
func (_StakingContractAddress1 *StakingContractAddress1Session) NewDeposit(pubkey []byte) (*types.Transaction, error) {
	return _StakingContractAddress1.Contract.NewDeposit(&_StakingContractAddress1.TransactOpts, pubkey)
}

// NewDeposit is a paid mutator transaction binding the contract method 0xdfcd068f.
//
// Solidity: function newDeposit(bytes pubkey) payable returns()
func (_StakingContractAddress1 *StakingContractAddress1TransactorSession) NewDeposit(pubkey []byte) (*types.Transaction, error) {
	return _StakingContractAddress1.Contract.NewDeposit(&_StakingContractAddress1.TransactOpts, pubkey)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns()
func (_StakingContractAddress1 *StakingContractAddress1Transactor) Withdraw(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _StakingContractAddress1.contract.Transact(opts, "withdraw", value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns()
func (_StakingContractAddress1 *StakingContractAddress1Session) Withdraw(value *big.Int) (*types.Transaction, error) {
	return _StakingContractAddress1.Contract.Withdraw(&_StakingContractAddress1.TransactOpts, value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns()
func (_StakingContractAddress1 *StakingContractAddress1TransactorSession) Withdraw(value *big.Int) (*types.Transaction, error) {
	return _StakingContractAddress1.Contract.Withdraw(&_StakingContractAddress1.TransactOpts, value)
}

// StakingContractAddress1OnNewDepositIterator is returned from FilterOnNewDeposit and is used to iterate over the raw logs and unpacked data for OnNewDeposit events raised by the StakingContractAddress1 contract.
type StakingContractAddress1OnNewDepositIterator struct {
	Event *StakingContractAddress1OnNewDeposit // Event containing the contract specifics and raw log

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
func (it *StakingContractAddress1OnNewDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingContractAddress1OnNewDeposit)
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
		it.Event = new(StakingContractAddress1OnNewDeposit)
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
func (it *StakingContractAddress1OnNewDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingContractAddress1OnNewDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingContractAddress1OnNewDeposit represents a OnNewDeposit event raised by the StakingContractAddress1 contract.
type StakingContractAddress1OnNewDeposit struct {
	Sender           common.Address
	ValidatorId      [32]byte
	ValidatorAddress common.Address
	Pubkey           []byte
	Value            *big.Int
	BlockNumber      *big.Int
	BlockTime        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnNewDeposit is a free log retrieval operation binding the contract event 0x9a1f4f083763f8508b19d4301c0110d2b47d99a8c5cf52c825c9e8cfea17f89c.
//
// Solidity: event OnNewDeposit(address indexed sender, bytes32 indexed validatorId, address indexed validatorAddress, bytes pubkey, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContractAddress1 *StakingContractAddress1Filterer) FilterOnNewDeposit(opts *bind.FilterOpts, sender []common.Address, validatorId [][32]byte, validatorAddress []common.Address) (*StakingContractAddress1OnNewDepositIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _StakingContractAddress1.contract.FilterLogs(opts, "OnNewDeposit", senderRule, validatorIdRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &StakingContractAddress1OnNewDepositIterator{contract: _StakingContractAddress1.contract, event: "OnNewDeposit", logs: logs, sub: sub}, nil
}

// WatchOnNewDeposit is a free log subscription operation binding the contract event 0x9a1f4f083763f8508b19d4301c0110d2b47d99a8c5cf52c825c9e8cfea17f89c.
//
// Solidity: event OnNewDeposit(address indexed sender, bytes32 indexed validatorId, address indexed validatorAddress, bytes pubkey, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContractAddress1 *StakingContractAddress1Filterer) WatchOnNewDeposit(opts *bind.WatchOpts, sink chan<- *StakingContractAddress1OnNewDeposit, sender []common.Address, validatorId [][32]byte, validatorAddress []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _StakingContractAddress1.contract.WatchLogs(opts, "OnNewDeposit", senderRule, validatorIdRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingContractAddress1OnNewDeposit)
				if err := _StakingContractAddress1.contract.UnpackLog(event, "OnNewDeposit", log); err != nil {
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

// ParseOnNewDeposit is a log parse operation binding the contract event 0x9a1f4f083763f8508b19d4301c0110d2b47d99a8c5cf52c825c9e8cfea17f89c.
//
// Solidity: event OnNewDeposit(address indexed sender, bytes32 indexed validatorId, address indexed validatorAddress, bytes pubkey, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContractAddress1 *StakingContractAddress1Filterer) ParseOnNewDeposit(log types.Log) (*StakingContractAddress1OnNewDeposit, error) {
	event := new(StakingContractAddress1OnNewDeposit)
	if err := _StakingContractAddress1.contract.UnpackLog(event, "OnNewDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingContractAddress1OnWithdrawKeyIterator is returned from FilterOnWithdrawKey and is used to iterate over the raw logs and unpacked data for OnWithdrawKey events raised by the StakingContractAddress1 contract.
type StakingContractAddress1OnWithdrawKeyIterator struct {
	Event *StakingContractAddress1OnWithdrawKey // Event containing the contract specifics and raw log

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
func (it *StakingContractAddress1OnWithdrawKeyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingContractAddress1OnWithdrawKey)
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
		it.Event = new(StakingContractAddress1OnWithdrawKey)
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
func (it *StakingContractAddress1OnWithdrawKeyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingContractAddress1OnWithdrawKeyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingContractAddress1OnWithdrawKey represents a OnWithdrawKey event raised by the StakingContractAddress1 contract.
type StakingContractAddress1OnWithdrawKey struct {
	Sender      common.Address
	Value       *big.Int
	BlockNumber *big.Int
	BlockTime   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOnWithdrawKey is a free log retrieval operation binding the contract event 0x4d4666331ec61727075c5624fde25f5510c566e528d0565f2a2263a23b70d81a.
//
// Solidity: event OnWithdrawKey(address sender, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContractAddress1 *StakingContractAddress1Filterer) FilterOnWithdrawKey(opts *bind.FilterOpts) (*StakingContractAddress1OnWithdrawKeyIterator, error) {

	logs, sub, err := _StakingContractAddress1.contract.FilterLogs(opts, "OnWithdrawKey")
	if err != nil {
		return nil, err
	}
	return &StakingContractAddress1OnWithdrawKeyIterator{contract: _StakingContractAddress1.contract, event: "OnWithdrawKey", logs: logs, sub: sub}, nil
}

// WatchOnWithdrawKey is a free log subscription operation binding the contract event 0x4d4666331ec61727075c5624fde25f5510c566e528d0565f2a2263a23b70d81a.
//
// Solidity: event OnWithdrawKey(address sender, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContractAddress1 *StakingContractAddress1Filterer) WatchOnWithdrawKey(opts *bind.WatchOpts, sink chan<- *StakingContractAddress1OnWithdrawKey) (event.Subscription, error) {

	logs, sub, err := _StakingContractAddress1.contract.WatchLogs(opts, "OnWithdrawKey")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingContractAddress1OnWithdrawKey)
				if err := _StakingContractAddress1.contract.UnpackLog(event, "OnWithdrawKey", log); err != nil {
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

// ParseOnWithdrawKey is a log parse operation binding the contract event 0x4d4666331ec61727075c5624fde25f5510c566e528d0565f2a2263a23b70d81a.
//
// Solidity: event OnWithdrawKey(address sender, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContractAddress1 *StakingContractAddress1Filterer) ParseOnWithdrawKey(log types.Log) (*StakingContractAddress1OnWithdrawKey, error) {
	event := new(StakingContractAddress1OnWithdrawKey)
	if err := _StakingContractAddress1.contract.UnpackLog(event, "OnWithdrawKey", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
