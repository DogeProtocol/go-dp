// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package consensuscontext

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

// ConsensuscontextMetaData contains all meta data concerning the Consensuscontext contract.
var ConsensuscontextMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contextId\",\"type\":\"string\"}],\"name\":\"deleteContext\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contextId\",\"type\":\"string\"}],\"name\":\"getContext\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contextId\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"context\",\"type\":\"bytes32\"}],\"name\":\"setContext\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061041a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806353949d4f1461004657806360eccdd414610076578063c1d971e414610092575b600080fd5b610060600480360381019061005b919061020f565b6100ae565b60405161006d9190610339565b60405180910390f35b610090600480360381019061008b919061020f565b6100d8565b005b6100ac60048036038101906100a79190610254565b610143565b005b60008083836040516100c1929190610320565b908152602001604051809103902054905092915050565b6000331461011b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161011290610354565b60405180910390fd5b6000828260405161012d929190610320565b9081526020016040518091039020600090555050565b60003314610186576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161017d90610354565b60405180910390fd5b8060008484604051610199929190610320565b908152602001604051809103902081905550505050565b6000813590506101bf816103a9565b92915050565b60008083601f8401126101d757600080fd5b8235905067ffffffffffffffff8111156101f057600080fd5b60208301915083600182028301111561020857600080fd5b9250929050565b6000806020838503121561022257600080fd5b600083013567ffffffffffffffff81111561023c57600080fd5b610248858286016101c5565b92509250509250929050565b60008060006040848603121561026957600080fd5b600084013567ffffffffffffffff81111561028357600080fd5b61028f868287016101c5565b935093505060206102a2868287016101b0565b9150509250925092565b6102b581610390565b82525050565b60006102c78385610385565b93506102d483858461039a565b82840190509392505050565b60006102ed601983610374565b91507f4f6e6c7920564d2063616c6c732061726520616c6c6f776564000000000000006000830152602082019050919050565b600061032d8284866102bb565b91508190509392505050565b600060208201905061034e60008301846102ac565b92915050565b6000602082019050818103600083015261036d816102e0565b9050919050565b600082825260208201905092915050565b600081905092915050565b6000819050919050565b82818337600083830152505050565b6103b281610390565b81146103bd57600080fd5b5056fea2646970667358221220064cef03d06c5c3a3ddcec1fffbc922205861979339d36970c508e7e29ef135664736f6c637826302e372e362d646576656c6f702e323032342e322e332b636f6d6d69742e33343239663030640057",
}

// ConsensuscontextABI is the input ABI used to generate the binding from.
// Deprecated: Use ConsensuscontextMetaData.ABI instead.
var ConsensuscontextABI = ConsensuscontextMetaData.ABI

// ConsensuscontextBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConsensuscontextMetaData.Bin instead.
var ConsensuscontextBin = ConsensuscontextMetaData.Bin

// DeployConsensuscontext deploys a new Ethereum contract, binding an instance of Consensuscontext to it.
func DeployConsensuscontext(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Consensuscontext, error) {
	parsed, err := ConsensuscontextMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConsensuscontextBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Consensuscontext{ConsensuscontextCaller: ConsensuscontextCaller{contract: contract}, ConsensuscontextTransactor: ConsensuscontextTransactor{contract: contract}, ConsensuscontextFilterer: ConsensuscontextFilterer{contract: contract}}, nil
}

// Consensuscontext is an auto generated Go binding around an Ethereum contract.
type Consensuscontext struct {
	ConsensuscontextCaller     // Read-only binding to the contract
	ConsensuscontextTransactor // Write-only binding to the contract
	ConsensuscontextFilterer   // Log filterer for contract events
}

// ConsensuscontextCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConsensuscontextCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensuscontextTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConsensuscontextTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensuscontextFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConsensuscontextFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensuscontextSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConsensuscontextSession struct {
	Contract     *Consensuscontext // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConsensuscontextCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConsensuscontextCallerSession struct {
	Contract *ConsensuscontextCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ConsensuscontextTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConsensuscontextTransactorSession struct {
	Contract     *ConsensuscontextTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ConsensuscontextRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConsensuscontextRaw struct {
	Contract *Consensuscontext // Generic contract binding to access the raw methods on
}

// ConsensuscontextCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConsensuscontextCallerRaw struct {
	Contract *ConsensuscontextCaller // Generic read-only contract binding to access the raw methods on
}

// ConsensuscontextTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConsensuscontextTransactorRaw struct {
	Contract *ConsensuscontextTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConsensuscontext creates a new instance of Consensuscontext, bound to a specific deployed contract.
func NewConsensuscontext(address common.Address, backend bind.ContractBackend) (*Consensuscontext, error) {
	contract, err := bindConsensuscontext(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Consensuscontext{ConsensuscontextCaller: ConsensuscontextCaller{contract: contract}, ConsensuscontextTransactor: ConsensuscontextTransactor{contract: contract}, ConsensuscontextFilterer: ConsensuscontextFilterer{contract: contract}}, nil
}

// NewConsensuscontextCaller creates a new read-only instance of Consensuscontext, bound to a specific deployed contract.
func NewConsensuscontextCaller(address common.Address, caller bind.ContractCaller) (*ConsensuscontextCaller, error) {
	contract, err := bindConsensuscontext(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensuscontextCaller{contract: contract}, nil
}

// NewConsensuscontextTransactor creates a new write-only instance of Consensuscontext, bound to a specific deployed contract.
func NewConsensuscontextTransactor(address common.Address, transactor bind.ContractTransactor) (*ConsensuscontextTransactor, error) {
	contract, err := bindConsensuscontext(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensuscontextTransactor{contract: contract}, nil
}

// NewConsensuscontextFilterer creates a new log filterer instance of Consensuscontext, bound to a specific deployed contract.
func NewConsensuscontextFilterer(address common.Address, filterer bind.ContractFilterer) (*ConsensuscontextFilterer, error) {
	contract, err := bindConsensuscontext(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConsensuscontextFilterer{contract: contract}, nil
}

// bindConsensuscontext binds a generic wrapper to an already deployed contract.
func bindConsensuscontext(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConsensuscontextABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Consensuscontext *ConsensuscontextRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Consensuscontext.Contract.ConsensuscontextCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Consensuscontext *ConsensuscontextRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consensuscontext.Contract.ConsensuscontextTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Consensuscontext *ConsensuscontextRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Consensuscontext.Contract.ConsensuscontextTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Consensuscontext *ConsensuscontextCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Consensuscontext.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Consensuscontext *ConsensuscontextTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consensuscontext.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Consensuscontext *ConsensuscontextTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Consensuscontext.Contract.contract.Transact(opts, method, params...)
}

// GetContext is a free data retrieval call binding the contract method 0x53949d4f.
//
// Solidity: function getContext(string contextId) view returns(bytes32)
func (_Consensuscontext *ConsensuscontextCaller) GetContext(opts *bind.CallOpts, contextId string) ([32]byte, error) {
	var out []interface{}
	err := _Consensuscontext.contract.Call(opts, &out, "getContext", contextId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetContext is a free data retrieval call binding the contract method 0x53949d4f.
//
// Solidity: function getContext(string contextId) view returns(bytes32)
func (_Consensuscontext *ConsensuscontextSession) GetContext(contextId string) ([32]byte, error) {
	return _Consensuscontext.Contract.GetContext(&_Consensuscontext.CallOpts, contextId)
}

// GetContext is a free data retrieval call binding the contract method 0x53949d4f.
//
// Solidity: function getContext(string contextId) view returns(bytes32)
func (_Consensuscontext *ConsensuscontextCallerSession) GetContext(contextId string) ([32]byte, error) {
	return _Consensuscontext.Contract.GetContext(&_Consensuscontext.CallOpts, contextId)
}

// DeleteContext is a paid mutator transaction binding the contract method 0x60eccdd4.
//
// Solidity: function deleteContext(string contextId) returns()
func (_Consensuscontext *ConsensuscontextTransactor) DeleteContext(opts *bind.TransactOpts, contextId string) (*types.Transaction, error) {
	return _Consensuscontext.contract.Transact(opts, "deleteContext", contextId)
}

// DeleteContext is a paid mutator transaction binding the contract method 0x60eccdd4.
//
// Solidity: function deleteContext(string contextId) returns()
func (_Consensuscontext *ConsensuscontextSession) DeleteContext(contextId string) (*types.Transaction, error) {
	return _Consensuscontext.Contract.DeleteContext(&_Consensuscontext.TransactOpts, contextId)
}

// DeleteContext is a paid mutator transaction binding the contract method 0x60eccdd4.
//
// Solidity: function deleteContext(string contextId) returns()
func (_Consensuscontext *ConsensuscontextTransactorSession) DeleteContext(contextId string) (*types.Transaction, error) {
	return _Consensuscontext.Contract.DeleteContext(&_Consensuscontext.TransactOpts, contextId)
}

// SetContext is a paid mutator transaction binding the contract method 0xc1d971e4.
//
// Solidity: function setContext(string contextId, bytes32 context) returns()
func (_Consensuscontext *ConsensuscontextTransactor) SetContext(opts *bind.TransactOpts, contextId string, context [32]byte) (*types.Transaction, error) {
	return _Consensuscontext.contract.Transact(opts, "setContext", contextId, context)
}

// SetContext is a paid mutator transaction binding the contract method 0xc1d971e4.
//
// Solidity: function setContext(string contextId, bytes32 context) returns()
func (_Consensuscontext *ConsensuscontextSession) SetContext(contextId string, context [32]byte) (*types.Transaction, error) {
	return _Consensuscontext.Contract.SetContext(&_Consensuscontext.TransactOpts, contextId, context)
}

// SetContext is a paid mutator transaction binding the contract method 0xc1d971e4.
//
// Solidity: function setContext(string contextId, bytes32 context) returns()
func (_Consensuscontext *ConsensuscontextTransactorSession) SetContext(contextId string, context [32]byte) (*types.Transaction, error) {
	return _Consensuscontext.Contract.SetContext(&_Consensuscontext.TransactOpts, contextId, context)
}
