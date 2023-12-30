// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package conversion

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

// ConversionMetaData contains all meta data concerning the Conversion contract.
var ConversionMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"quantumAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"ethAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OnConversion\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"quantumAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ethereumSignature\",\"type\":\"string\"}],\"name\":\"OnRequestConversion\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ethAddress\",\"type\":\"address\"}],\"name\":\"getAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ethAddress\",\"type\":\"address\"}],\"name\":\"getConversionStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ethAddress\",\"type\":\"address\"}],\"name\":\"getQuantumAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"ethSignature\",\"type\":\"string\"}],\"name\":\"requestConversion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ethAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"quantumAddress\",\"type\":\"address\"}],\"name\":\"setConverted\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061097a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80631947f47d1461005c57806331f8a7151461008c5780634ce89f58146100bc578063f5a79767146100ec578063fa2ff79c1461011c575b600080fd5b61007660048036038101906100719190610513565b61014c565b604051610083919061087b565b60405180910390f35b6100a660048036038101906100a191906104ae565b610198565b6040516100b39190610746565b60405180910390f35b6100d660048036038101906100d191906104ae565b6101b5565b6040516100e3919061078a565b60405180910390f35b610106600480360381019061010191906104ae565b6101df565b6040516101139190610860565b60405180910390f35b610136600480360381019061013191906104d7565b6101fb565b6040516101439190610860565b60405180910390f35b6000337fe3c11e89286c3875683a51a5862ab8536bc29ca10fc579954cf06bfd1a50635d8686868660405161018494939291906107a5565b60405180910390a260009050949350505050565b600060016000838152602001908152602001600020549050919050565b60006002600083815260200190815260200160002060009054906101000a900460ff169050919050565b6000806000838152602001908152602001600020549050919050565b600080331461023f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161023690610800565b60405180910390fd5b60008060008581526020019081526020016000205411610294576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161028b90610840565b60405180910390fd5b600015156002600085815260200190815260200160002060009054906101000a900460ff161515146102fb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102f290610820565b60405180910390fd5b60016002600085815260200190815260200160002060006101000a81548160ff0219169083151502179055508160016000858152602001908152602001600020819055506000826000808681526020019081526020016000205460405161036190610731565b60006040518083038185875af1925050503d806000811461039e576040519150601f19603f3d011682016040523d82523d6000602084013e6103a3565b606091505b50509050806103e7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103de906107e0565b60405180910390fd5b827f07ba754d52e1422be3ad73643dfaf20685c9fd3f05a45ddb82b563b3da095c6e856000808881526020019081526020016000205460405161042b929190610761565b60405180910390a26000808581526020019081526020016000205491505092915050565b60008135905061045e81610907565b92915050565b60008083601f84011261047657600080fd5b8235905067ffffffffffffffff81111561048f57600080fd5b6020830191508360018202830111156104a757600080fd5b9250929050565b6000602082840312156104c057600080fd5b60006104ce8482850161044f565b91505092915050565b600080604083850312156104ea57600080fd5b60006104f88582860161044f565b92505060206105098582860161044f565b9150509250929050565b6000806000806040858703121561052957600080fd5b600085013567ffffffffffffffff81111561054357600080fd5b61054f87828801610464565b9450945050602085013567ffffffffffffffff81111561056e57600080fd5b61057a87828801610464565b925092505092959194509250565b610591816108b2565b82525050565b6105a0816108c4565b82525050565b60006105b283856108a1565b93506105bf8385846108e7565b6105c8836108f6565b840190509392505050565b60006105e06017836108a1565b91507f5472616e736665722042616c616e6365206661696c65640000000000000000006000830152602082019050919050565b60006106206019836108a1565b91507f4f6e6c7920564d2063616c6c732061726520616c6c6f776564000000000000006000830152602082019050919050565b60006106606011836108a1565b91507f416c726561647920636f6e7665727465640000000000000000000000000000006000830152602082019050919050565b60006106a06024836108a1565b91507f6574684164647265737320446f65736e277420657869737420696e20736e617060008301527f73686f74000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000610706600083610896565b9150600082019050919050565b61071c816108d0565b82525050565b61072b816108da565b82525050565b600061073c826106f9565b9150819050919050565b600060208201905061075b6000830184610588565b92915050565b60006040820190506107766000830185610588565b6107836020830184610713565b9392505050565b600060208201905061079f6000830184610597565b92915050565b600060408201905081810360008301526107c08186886105a6565b905081810360208301526107d58184866105a6565b905095945050505050565b600060208201905081810360008301526107f9816105d3565b9050919050565b6000602082019050818103600083015261081981610613565b9050919050565b6000602082019050818103600083015261083981610653565b9050919050565b6000602082019050818103600083015261085981610693565b9050919050565b60006020820190506108756000830184610713565b92915050565b60006020820190506108906000830184610722565b92915050565b600081905092915050565b600082825260208201905092915050565b60006108bd826108d0565b9050919050565b60008115159050919050565b6000819050919050565b600060ff82169050919050565b82818337600083830152505050565b6000601f19601f8301169050919050565b610910816108b2565b811461091b57600080fd5b5056fea26469706673582212208fed4f0439b626b16caa29697df391b3bac3c9c5997e2de07aab7c6004b9bb5d64736f6c637828302e372e362d646576656c6f702e323032332e31322e33302b636f6d6d69742e37326538396665320059",
}

// ConversionABI is the input ABI used to generate the binding from.
// Deprecated: Use ConversionMetaData.ABI instead.
var ConversionABI = ConversionMetaData.ABI

// ConversionBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConversionMetaData.Bin instead.
var ConversionBin = ConversionMetaData.Bin

// DeployConversion deploys a new Ethereum contract, binding an instance of Conversion to it.
func DeployConversion(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Conversion, error) {
	parsed, err := ConversionMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConversionBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Conversion{ConversionCaller: ConversionCaller{contract: contract}, ConversionTransactor: ConversionTransactor{contract: contract}, ConversionFilterer: ConversionFilterer{contract: contract}}, nil
}

// Conversion is an auto generated Go binding around an Ethereum contract.
type Conversion struct {
	ConversionCaller     // Read-only binding to the contract
	ConversionTransactor // Write-only binding to the contract
	ConversionFilterer   // Log filterer for contract events
}

// ConversionCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConversionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConversionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConversionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConversionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConversionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConversionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConversionSession struct {
	Contract     *Conversion       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConversionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConversionCallerSession struct {
	Contract *ConversionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ConversionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConversionTransactorSession struct {
	Contract     *ConversionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ConversionRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConversionRaw struct {
	Contract *Conversion // Generic contract binding to access the raw methods on
}

// ConversionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConversionCallerRaw struct {
	Contract *ConversionCaller // Generic read-only contract binding to access the raw methods on
}

// ConversionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConversionTransactorRaw struct {
	Contract *ConversionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConversion creates a new instance of Conversion, bound to a specific deployed contract.
func NewConversion(address common.Address, backend bind.ContractBackend) (*Conversion, error) {
	contract, err := bindConversion(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Conversion{ConversionCaller: ConversionCaller{contract: contract}, ConversionTransactor: ConversionTransactor{contract: contract}, ConversionFilterer: ConversionFilterer{contract: contract}}, nil
}

// NewConversionCaller creates a new read-only instance of Conversion, bound to a specific deployed contract.
func NewConversionCaller(address common.Address, caller bind.ContractCaller) (*ConversionCaller, error) {
	contract, err := bindConversion(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConversionCaller{contract: contract}, nil
}

// NewConversionTransactor creates a new write-only instance of Conversion, bound to a specific deployed contract.
func NewConversionTransactor(address common.Address, transactor bind.ContractTransactor) (*ConversionTransactor, error) {
	contract, err := bindConversion(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConversionTransactor{contract: contract}, nil
}

// NewConversionFilterer creates a new log filterer instance of Conversion, bound to a specific deployed contract.
func NewConversionFilterer(address common.Address, filterer bind.ContractFilterer) (*ConversionFilterer, error) {
	contract, err := bindConversion(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConversionFilterer{contract: contract}, nil
}

// bindConversion binds a generic wrapper to an already deployed contract.
func bindConversion(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConversionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Conversion *ConversionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Conversion.Contract.ConversionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Conversion *ConversionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Conversion.Contract.ConversionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Conversion *ConversionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Conversion.Contract.ConversionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Conversion *ConversionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Conversion.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Conversion *ConversionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Conversion.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Conversion *ConversionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Conversion.Contract.contract.Transact(opts, method, params...)
}

// GetAmount is a free data retrieval call binding the contract method 0xf5a79767.
//
// Solidity: function getAmount(address ethAddress) view returns(uint256)
func (_Conversion *ConversionCaller) GetAmount(opts *bind.CallOpts, ethAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Conversion.contract.Call(opts, &out, "getAmount", ethAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmount is a free data retrieval call binding the contract method 0xf5a79767.
//
// Solidity: function getAmount(address ethAddress) view returns(uint256)
func (_Conversion *ConversionSession) GetAmount(ethAddress common.Address) (*big.Int, error) {
	return _Conversion.Contract.GetAmount(&_Conversion.CallOpts, ethAddress)
}

// GetAmount is a free data retrieval call binding the contract method 0xf5a79767.
//
// Solidity: function getAmount(address ethAddress) view returns(uint256)
func (_Conversion *ConversionCallerSession) GetAmount(ethAddress common.Address) (*big.Int, error) {
	return _Conversion.Contract.GetAmount(&_Conversion.CallOpts, ethAddress)
}

// GetConversionStatus is a free data retrieval call binding the contract method 0x4ce89f58.
//
// Solidity: function getConversionStatus(address ethAddress) view returns(bool)
func (_Conversion *ConversionCaller) GetConversionStatus(opts *bind.CallOpts, ethAddress common.Address) (bool, error) {
	var out []interface{}
	err := _Conversion.contract.Call(opts, &out, "getConversionStatus", ethAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetConversionStatus is a free data retrieval call binding the contract method 0x4ce89f58.
//
// Solidity: function getConversionStatus(address ethAddress) view returns(bool)
func (_Conversion *ConversionSession) GetConversionStatus(ethAddress common.Address) (bool, error) {
	return _Conversion.Contract.GetConversionStatus(&_Conversion.CallOpts, ethAddress)
}

// GetConversionStatus is a free data retrieval call binding the contract method 0x4ce89f58.
//
// Solidity: function getConversionStatus(address ethAddress) view returns(bool)
func (_Conversion *ConversionCallerSession) GetConversionStatus(ethAddress common.Address) (bool, error) {
	return _Conversion.Contract.GetConversionStatus(&_Conversion.CallOpts, ethAddress)
}

// GetQuantumAddress is a free data retrieval call binding the contract method 0x31f8a715.
//
// Solidity: function getQuantumAddress(address ethAddress) view returns(address)
func (_Conversion *ConversionCaller) GetQuantumAddress(opts *bind.CallOpts, ethAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _Conversion.contract.Call(opts, &out, "getQuantumAddress", ethAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetQuantumAddress is a free data retrieval call binding the contract method 0x31f8a715.
//
// Solidity: function getQuantumAddress(address ethAddress) view returns(address)
func (_Conversion *ConversionSession) GetQuantumAddress(ethAddress common.Address) (common.Address, error) {
	return _Conversion.Contract.GetQuantumAddress(&_Conversion.CallOpts, ethAddress)
}

// GetQuantumAddress is a free data retrieval call binding the contract method 0x31f8a715.
//
// Solidity: function getQuantumAddress(address ethAddress) view returns(address)
func (_Conversion *ConversionCallerSession) GetQuantumAddress(ethAddress common.Address) (common.Address, error) {
	return _Conversion.Contract.GetQuantumAddress(&_Conversion.CallOpts, ethAddress)
}

// RequestConversion is a paid mutator transaction binding the contract method 0x1947f47d.
//
// Solidity: function requestConversion(string ethAddress, string ethSignature) returns(uint8)
func (_Conversion *ConversionTransactor) RequestConversion(opts *bind.TransactOpts, ethAddress string, ethSignature string) (*types.Transaction, error) {
	return _Conversion.contract.Transact(opts, "requestConversion", ethAddress, ethSignature)
}

// RequestConversion is a paid mutator transaction binding the contract method 0x1947f47d.
//
// Solidity: function requestConversion(string ethAddress, string ethSignature) returns(uint8)
func (_Conversion *ConversionSession) RequestConversion(ethAddress string, ethSignature string) (*types.Transaction, error) {
	return _Conversion.Contract.RequestConversion(&_Conversion.TransactOpts, ethAddress, ethSignature)
}

// RequestConversion is a paid mutator transaction binding the contract method 0x1947f47d.
//
// Solidity: function requestConversion(string ethAddress, string ethSignature) returns(uint8)
func (_Conversion *ConversionTransactorSession) RequestConversion(ethAddress string, ethSignature string) (*types.Transaction, error) {
	return _Conversion.Contract.RequestConversion(&_Conversion.TransactOpts, ethAddress, ethSignature)
}

// SetConverted is a paid mutator transaction binding the contract method 0xfa2ff79c.
//
// Solidity: function setConverted(address ethAddress, address quantumAddress) returns(uint256)
func (_Conversion *ConversionTransactor) SetConverted(opts *bind.TransactOpts, ethAddress common.Address, quantumAddress common.Address) (*types.Transaction, error) {
	return _Conversion.contract.Transact(opts, "setConverted", ethAddress, quantumAddress)
}

// SetConverted is a paid mutator transaction binding the contract method 0xfa2ff79c.
//
// Solidity: function setConverted(address ethAddress, address quantumAddress) returns(uint256)
func (_Conversion *ConversionSession) SetConverted(ethAddress common.Address, quantumAddress common.Address) (*types.Transaction, error) {
	return _Conversion.Contract.SetConverted(&_Conversion.TransactOpts, ethAddress, quantumAddress)
}

// SetConverted is a paid mutator transaction binding the contract method 0xfa2ff79c.
//
// Solidity: function setConverted(address ethAddress, address quantumAddress) returns(uint256)
func (_Conversion *ConversionTransactorSession) SetConverted(ethAddress common.Address, quantumAddress common.Address) (*types.Transaction, error) {
	return _Conversion.Contract.SetConverted(&_Conversion.TransactOpts, ethAddress, quantumAddress)
}

// ConversionOnConversionIterator is returned from FilterOnConversion and is used to iterate over the raw logs and unpacked data for OnConversion events raised by the Conversion contract.
type ConversionOnConversionIterator struct {
	Event *ConversionOnConversion // Event containing the contract specifics and raw log

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
func (it *ConversionOnConversionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConversionOnConversion)
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
		it.Event = new(ConversionOnConversion)
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
func (it *ConversionOnConversionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConversionOnConversionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConversionOnConversion represents a OnConversion event raised by the Conversion contract.
type ConversionOnConversion struct {
	QuantumAddress common.Address
	EthAddress     common.Address
	Amount         *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterOnConversion is a free log retrieval operation binding the contract event 0x07ba754d52e1422be3ad73643dfaf20685c9fd3f05a45ddb82b563b3da095c6e.
//
// Solidity: event OnConversion(address indexed quantumAddress, address ethAddress, uint256 amount)
func (_Conversion *ConversionFilterer) FilterOnConversion(opts *bind.FilterOpts, quantumAddress []common.Address) (*ConversionOnConversionIterator, error) {

	var quantumAddressRule []interface{}
	for _, quantumAddressItem := range quantumAddress {
		quantumAddressRule = append(quantumAddressRule, quantumAddressItem)
	}

	logs, sub, err := _Conversion.contract.FilterLogs(opts, "OnConversion", quantumAddressRule)
	if err != nil {
		return nil, err
	}
	return &ConversionOnConversionIterator{contract: _Conversion.contract, event: "OnConversion", logs: logs, sub: sub}, nil
}

// WatchOnConversion is a free log subscription operation binding the contract event 0x07ba754d52e1422be3ad73643dfaf20685c9fd3f05a45ddb82b563b3da095c6e.
//
// Solidity: event OnConversion(address indexed quantumAddress, address ethAddress, uint256 amount)
func (_Conversion *ConversionFilterer) WatchOnConversion(opts *bind.WatchOpts, sink chan<- *ConversionOnConversion, quantumAddress []common.Address) (event.Subscription, error) {

	var quantumAddressRule []interface{}
	for _, quantumAddressItem := range quantumAddress {
		quantumAddressRule = append(quantumAddressRule, quantumAddressItem)
	}

	logs, sub, err := _Conversion.contract.WatchLogs(opts, "OnConversion", quantumAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConversionOnConversion)
				if err := _Conversion.contract.UnpackLog(event, "OnConversion", log); err != nil {
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

// ParseOnConversion is a log parse operation binding the contract event 0x07ba754d52e1422be3ad73643dfaf20685c9fd3f05a45ddb82b563b3da095c6e.
//
// Solidity: event OnConversion(address indexed quantumAddress, address ethAddress, uint256 amount)
func (_Conversion *ConversionFilterer) ParseOnConversion(log types.Log) (*ConversionOnConversion, error) {
	event := new(ConversionOnConversion)
	if err := _Conversion.contract.UnpackLog(event, "OnConversion", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConversionOnRequestConversionIterator is returned from FilterOnRequestConversion and is used to iterate over the raw logs and unpacked data for OnRequestConversion events raised by the Conversion contract.
type ConversionOnRequestConversionIterator struct {
	Event *ConversionOnRequestConversion // Event containing the contract specifics and raw log

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
func (it *ConversionOnRequestConversionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConversionOnRequestConversion)
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
		it.Event = new(ConversionOnRequestConversion)
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
func (it *ConversionOnRequestConversionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConversionOnRequestConversionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConversionOnRequestConversion represents a OnRequestConversion event raised by the Conversion contract.
type ConversionOnRequestConversion struct {
	QuantumAddress    common.Address
	EthAddress        string
	EthereumSignature string
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterOnRequestConversion is a free log retrieval operation binding the contract event 0xe3c11e89286c3875683a51a5862ab8536bc29ca10fc579954cf06bfd1a50635d.
//
// Solidity: event OnRequestConversion(address indexed quantumAddress, string ethAddress, string ethereumSignature)
func (_Conversion *ConversionFilterer) FilterOnRequestConversion(opts *bind.FilterOpts, quantumAddress []common.Address) (*ConversionOnRequestConversionIterator, error) {

	var quantumAddressRule []interface{}
	for _, quantumAddressItem := range quantumAddress {
		quantumAddressRule = append(quantumAddressRule, quantumAddressItem)
	}

	logs, sub, err := _Conversion.contract.FilterLogs(opts, "OnRequestConversion", quantumAddressRule)
	if err != nil {
		return nil, err
	}
	return &ConversionOnRequestConversionIterator{contract: _Conversion.contract, event: "OnRequestConversion", logs: logs, sub: sub}, nil
}

// WatchOnRequestConversion is a free log subscription operation binding the contract event 0xe3c11e89286c3875683a51a5862ab8536bc29ca10fc579954cf06bfd1a50635d.
//
// Solidity: event OnRequestConversion(address indexed quantumAddress, string ethAddress, string ethereumSignature)
func (_Conversion *ConversionFilterer) WatchOnRequestConversion(opts *bind.WatchOpts, sink chan<- *ConversionOnRequestConversion, quantumAddress []common.Address) (event.Subscription, error) {

	var quantumAddressRule []interface{}
	for _, quantumAddressItem := range quantumAddress {
		quantumAddressRule = append(quantumAddressRule, quantumAddressItem)
	}

	logs, sub, err := _Conversion.contract.WatchLogs(opts, "OnRequestConversion", quantumAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConversionOnRequestConversion)
				if err := _Conversion.contract.UnpackLog(event, "OnRequestConversion", log); err != nil {
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

// ParseOnRequestConversion is a log parse operation binding the contract event 0xe3c11e89286c3875683a51a5862ab8536bc29ca10fc579954cf06bfd1a50635d.
//
// Solidity: event OnRequestConversion(address indexed quantumAddress, string ethAddress, string ethereumSignature)
func (_Conversion *ConversionFilterer) ParseOnRequestConversion(log types.Log) (*ConversionOnRequestConversion, error) {
	event := new(ConversionOnRequestConversion)
	if err := _Conversion.contract.UnpackLog(event, "OnRequestConversion", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
