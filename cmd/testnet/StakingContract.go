// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

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

// StakingContractMetaData contains all meta data concerning the StakingContract contract.
var StakingContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"validatorId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnNewDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"validatorId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"reward\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnRewardDepositKey\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnWithdrawKey\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"depositBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"listValidator\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"}],\"name\":\"newDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalDepositBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506000808190555060006001819055506111ba8061002f6000396000f3fe6080604052600436106100705760003560e01c806375697e661161004e57806375697e66146100b4578063dfcd068f146100df578063e8c0a0df146100fb578063fba13bd01461012657610070565b8063116b5e47146100755780632dfdf0b51461007f5780633ccfd60b146100aa575b600080fd5b61007d610163565b005b34801561008b57600080fd5b506100946102da565b6040516100a19190610fec565b60405180910390f35b6100b26102e3565b005b3480156100c057600080fd5b506100c961049b565b6040516100d69190610edc565b60405180910390f35b6100f960048036038101906100f49190610ae1565b610529565b005b34801561010757600080fd5b506101106108f4565b60405161011d9190610fec565b60405180910390f35b34801561013257600080fd5b5061014d60048036038101906101489190610ab8565b6108fe565b60405161015a9190610fec565b60405180910390f35b600034116101a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161019d90610f4c565b60405180910390fd5b60006101b133610947565b905060006004600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600a816040516020016101fe9190610e02565b604051602081830303815290604052511161024e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161024590610fac565b60405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f19350505050158015610294573d6000803e3d6000fd5b507fe0b518260035297556cfeb160ef4b66aed5ba1606403b996e4102fdd87e133663383833443426040516102ce96959493929190610e36565b60405180910390a15050565b60008054905090565b34600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015610365576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161035c90610fcc565b60405180910390fd5b61037a3460015461096e90919063ffffffff16565b6001819055506103d234600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461096e90919063ffffffff16565b600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503373ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f1935050505015801561045b573d6000803e3d6000fd5b507f4d4666331ec61727075c5624fde25f5510c566e528d0565f2a2263a23b70d81a333443426040516104919493929190610e97565b60405180910390a1565b6060600680548060200260200160405190810160405280929190818152602001828054801561051f57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190600101908083116104d5575b5050505050905090565b6000828290501161056f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161056690610f6c565b60405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff1660046000600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415610650576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161064790610f8c565b60405180910390fd5b610666600160005461098590919063ffffffff16565b6000819055506106813460015461098590919063ffffffff16565b6001819055506106d934600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461098590919063ffffffff16565b600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506000828260019080926107319392919061106d565b60405161073f929190610e1d565b604051809103902090506000610754826109a1565b9050600061076182610947565b905084846003600084815260200190815260200160002091906107859291906109ae565b50336004600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506006829080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff16813373ffffffffffffffffffffffffffffffffffffffff167f9a1f4f083763f8508b19d4301c0110d2b47d99a8c5cf52c825c9e8cfea17f89c88883443426040516108e5959493929190610efe565b60405180910390a45050505050565b6000600154905090565b6000600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600060608273ffffffffffffffffffffffffffffffffffffffff16901b60001b9050919050565b60008282111561097a57fe5b818303905092915050565b60008082840190508381101561099757fe5b8091505092915050565b60008160001c9050919050565b828054600181600116156101000203166002900490600052602060002090601f0160209004810192826109e45760008555610a2b565b82601f106109fd57803560ff1916838001178555610a2b565b82800160010185558215610a2b579182015b82811115610a2a578235825591602001919060010190610a0f565b5b509050610a389190610a3c565b5090565b5b80821115610a55576000816000905550600101610a3d565b5090565b600081359050610a688161116d565b92915050565b60008083601f840112610a8057600080fd5b8235905067ffffffffffffffff811115610a9957600080fd5b602083019150836001820283011115610ab157600080fd5b9250929050565b600060208284031215610aca57600080fd5b6000610ad884828501610a59565b91505092915050565b60008060208385031215610af457600080fd5b600083013567ffffffffffffffff811115610b0e57600080fd5b610b1a85828601610a6e565b92509250509250929050565b6000610b328383610b4d565b60208301905092915050565b610b47816110e6565b82525050565b610b56816110a0565b82525050565b610b65816110a0565b82525050565b610b7c610b77826110a0565b61112b565b82525050565b6000610b8d82611017565b610b97818561102f565b9350610ba283611007565b8060005b83811015610bd3578151610bba8882610b26565b9750610bc583611022565b925050600181019050610ba6565b5085935050505092915050565b610be9816110b2565b82525050565b6000610bfb8385611040565b9350610c0883858461111c565b610c118361114f565b840190509392505050565b6000610c288385611051565b9350610c3583858461111c565b82840190509392505050565b6000610c4e60258361105c565b91507f5374616b696e67436f6e74726163743a207265776172642076616c756520746f60008301527f6f206c6f770000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000610cb460168361105c565b91507f5075626c69636b6579206973206e6f742076616c6964000000000000000000006000830152602082019050919050565b6000610cf460128361105c565b91507f53656e64657220686176652065786973747300000000000000000000000000006000830152602082019050919050565b6000610d3460238361105c565b91507f5374616b696e67436f6e74726163743a2076616c696461746f7220697320656d60008301527f70747900000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000610d9a60218361105c565b91507f5374616b696e67436f6e74726163743a20696e737566666963656e742066756e60008301527f64000000000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b610dfc816110dc565b82525050565b6000610e0e8284610b6b565b60148201915081905092915050565b6000610e2a828486610c1c565b91508190509392505050565b600060c082019050610e4b6000830189610b3e565b610e586020830188610be0565b610e656040830187610b5c565b610e726060830186610df3565b610e7f6080830185610df3565b610e8c60a0830184610df3565b979650505050505050565b6000608082019050610eac6000830187610b3e565b610eb96020830186610df3565b610ec66040830185610df3565b610ed36060830184610df3565b95945050505050565b60006020820190508181036000830152610ef68184610b82565b905092915050565b60006080820190508181036000830152610f19818789610bef565b9050610f286020830186610df3565b610f356040830185610df3565b610f426060830184610df3565b9695505050505050565b60006020820190508181036000830152610f6581610c41565b9050919050565b60006020820190508181036000830152610f8581610ca7565b9050919050565b60006020820190508181036000830152610fa581610ce7565b9050919050565b60006020820190508181036000830152610fc581610d27565b9050919050565b60006020820190508181036000830152610fe581610d8d565b9050919050565b60006020820190506110016000830184610df3565b92915050565b6000819050602082019050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b6000808585111561107d57600080fd5b8386111561108a57600080fd5b6001850283019150848603905094509492505050565b60006110ab826110bc565b9050919050565b6000819050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b60006110f1826110f8565b9050919050565b60006111038261110a565b9050919050565b6000611115826110bc565b9050919050565b82818337600083830152505050565b60006111368261113d565b9050919050565b600061114882611160565b9050919050565b6000601f19601f8301169050919050565b60008160601b9050919050565b611176816110a0565b811461118157600080fd5b5056fea26469706673582212209ee10c0da0938c488ad04180d56697c45357d5ac690edbd051ae5fce27c3767664736f6c63430007060033",
}

// StakingContractABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingContractMetaData.ABI instead.
var StakingContractABI = StakingContractMetaData.ABI

// StakingContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingContractMetaData.Bin instead.
var StakingContractBin = StakingContractMetaData.Bin

// DeployStakingContract deploys a new Ethereum contract, binding an instance of StakingContract to it.
func DeployStakingContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StakingContract, error) {
	parsed, err := StakingContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StakingContract{StakingContractCaller: StakingContractCaller{contract: contract}, StakingContractTransactor: StakingContractTransactor{contract: contract}, StakingContractFilterer: StakingContractFilterer{contract: contract}}, nil
}

// StakingContract is an auto generated Go binding around an Ethereum contract.
type StakingContract struct {
	StakingContractCaller     // Read-only binding to the contract
	StakingContractTransactor // Write-only binding to the contract
	StakingContractFilterer   // Log filterer for contract events
}

// StakingContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingContractSession struct {
	Contract     *StakingContract  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingContractCallerSession struct {
	Contract *StakingContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// StakingContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingContractTransactorSession struct {
	Contract     *StakingContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// StakingContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingContractRaw struct {
	Contract *StakingContract // Generic contract binding to access the raw methods on
}

// StakingContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingContractCallerRaw struct {
	Contract *StakingContractCaller // Generic read-only contract binding to access the raw methods on
}

// StakingContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingContractTransactorRaw struct {
	Contract *StakingContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingContract creates a new instance of StakingContract, bound to a specific deployed contract.
func NewStakingContract(address common.Address, backend bind.ContractBackend) (*StakingContract, error) {
	contract, err := bindStakingContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingContract{StakingContractCaller: StakingContractCaller{contract: contract}, StakingContractTransactor: StakingContractTransactor{contract: contract}, StakingContractFilterer: StakingContractFilterer{contract: contract}}, nil
}

// NewStakingContractCaller creates a new read-only instance of StakingContract, bound to a specific deployed contract.
func NewStakingContractCaller(address common.Address, caller bind.ContractCaller) (*StakingContractCaller, error) {
	contract, err := bindStakingContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingContractCaller{contract: contract}, nil
}

// NewStakingContractTransactor creates a new write-only instance of StakingContract, bound to a specific deployed contract.
func NewStakingContractTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingContractTransactor, error) {
	contract, err := bindStakingContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingContractTransactor{contract: contract}, nil
}

// NewStakingContractFilterer creates a new log filterer instance of StakingContract, bound to a specific deployed contract.
func NewStakingContractFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingContractFilterer, error) {
	contract, err := bindStakingContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingContractFilterer{contract: contract}, nil
}

// bindStakingContract binds a generic wrapper to an already deployed contract.
func bindStakingContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingContract *StakingContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingContract.Contract.StakingContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingContract *StakingContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingContract.Contract.StakingContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingContract *StakingContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingContract.Contract.StakingContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingContract *StakingContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingContract *StakingContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingContract *StakingContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingContract.Contract.contract.Transact(opts, method, params...)
}

// DepositBalanceOf is a free data retrieval call binding the contract method 0xfba13bd0.
//
// Solidity: function depositBalanceOf(address owner) view returns(uint256)
func (_StakingContract *StakingContractCaller) DepositBalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingContract.contract.Call(opts, &out, "depositBalanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositBalanceOf is a free data retrieval call binding the contract method 0xfba13bd0.
//
// Solidity: function depositBalanceOf(address owner) view returns(uint256)
func (_StakingContract *StakingContractSession) DepositBalanceOf(owner common.Address) (*big.Int, error) {
	return _StakingContract.Contract.DepositBalanceOf(&_StakingContract.CallOpts, owner)
}

// DepositBalanceOf is a free data retrieval call binding the contract method 0xfba13bd0.
//
// Solidity: function depositBalanceOf(address owner) view returns(uint256)
func (_StakingContract *StakingContractCallerSession) DepositBalanceOf(owner common.Address) (*big.Int, error) {
	return _StakingContract.Contract.DepositBalanceOf(&_StakingContract.CallOpts, owner)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_StakingContract *StakingContractCaller) DepositCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingContract.contract.Call(opts, &out, "depositCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_StakingContract *StakingContractSession) DepositCount() (*big.Int, error) {
	return _StakingContract.Contract.DepositCount(&_StakingContract.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_StakingContract *StakingContractCallerSession) DepositCount() (*big.Int, error) {
	return _StakingContract.Contract.DepositCount(&_StakingContract.CallOpts)
}

// ListValidator is a free data retrieval call binding the contract method 0x75697e66.
//
// Solidity: function listValidator() view returns(address[])
func (_StakingContract *StakingContractCaller) ListValidator(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _StakingContract.contract.Call(opts, &out, "listValidator")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ListValidator is a free data retrieval call binding the contract method 0x75697e66.
//
// Solidity: function listValidator() view returns(address[])
func (_StakingContract *StakingContractSession) ListValidator() ([]common.Address, error) {
	return _StakingContract.Contract.ListValidator(&_StakingContract.CallOpts)
}

// ListValidator is a free data retrieval call binding the contract method 0x75697e66.
//
// Solidity: function listValidator() view returns(address[])
func (_StakingContract *StakingContractCallerSession) ListValidator() ([]common.Address, error) {
	return _StakingContract.Contract.ListValidator(&_StakingContract.CallOpts)
}

// TotalDepositBalance is a free data retrieval call binding the contract method 0xe8c0a0df.
//
// Solidity: function totalDepositBalance() view returns(uint256)
func (_StakingContract *StakingContractCaller) TotalDepositBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingContract.contract.Call(opts, &out, "totalDepositBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalDepositBalance is a free data retrieval call binding the contract method 0xe8c0a0df.
//
// Solidity: function totalDepositBalance() view returns(uint256)
func (_StakingContract *StakingContractSession) TotalDepositBalance() (*big.Int, error) {
	return _StakingContract.Contract.TotalDepositBalance(&_StakingContract.CallOpts)
}

// TotalDepositBalance is a free data retrieval call binding the contract method 0xe8c0a0df.
//
// Solidity: function totalDepositBalance() view returns(uint256)
func (_StakingContract *StakingContractCallerSession) TotalDepositBalance() (*big.Int, error) {
	return _StakingContract.Contract.TotalDepositBalance(&_StakingContract.CallOpts)
}

// NewDeposit is a paid mutator transaction binding the contract method 0xdfcd068f.
//
// Solidity: function newDeposit(bytes pubkey) payable returns()
func (_StakingContract *StakingContractTransactor) NewDeposit(opts *bind.TransactOpts, pubkey []byte) (*types.Transaction, error) {
	return _StakingContract.contract.Transact(opts, "newDeposit", pubkey)
}

// NewDeposit is a paid mutator transaction binding the contract method 0xdfcd068f.
//
// Solidity: function newDeposit(bytes pubkey) payable returns()
func (_StakingContract *StakingContractSession) NewDeposit(pubkey []byte) (*types.Transaction, error) {
	return _StakingContract.Contract.NewDeposit(&_StakingContract.TransactOpts, pubkey)
}

// NewDeposit is a paid mutator transaction binding the contract method 0xdfcd068f.
//
// Solidity: function newDeposit(bytes pubkey) payable returns()
func (_StakingContract *StakingContractTransactorSession) NewDeposit(pubkey []byte) (*types.Transaction, error) {
	return _StakingContract.Contract.NewDeposit(&_StakingContract.TransactOpts, pubkey)
}

// RewardDeposit is a paid mutator transaction binding the contract method 0x116b5e47.
//
// Solidity: function rewardDeposit() payable returns()
func (_StakingContract *StakingContractTransactor) RewardDeposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingContract.contract.Transact(opts, "rewardDeposit")
}

// RewardDeposit is a paid mutator transaction binding the contract method 0x116b5e47.
//
// Solidity: function rewardDeposit() payable returns()
func (_StakingContract *StakingContractSession) RewardDeposit() (*types.Transaction, error) {
	return _StakingContract.Contract.RewardDeposit(&_StakingContract.TransactOpts)
}

// RewardDeposit is a paid mutator transaction binding the contract method 0x116b5e47.
//
// Solidity: function rewardDeposit() payable returns()
func (_StakingContract *StakingContractTransactorSession) RewardDeposit() (*types.Transaction, error) {
	return _StakingContract.Contract.RewardDeposit(&_StakingContract.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() payable returns()
func (_StakingContract *StakingContractTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingContract.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() payable returns()
func (_StakingContract *StakingContractSession) Withdraw() (*types.Transaction, error) {
	return _StakingContract.Contract.Withdraw(&_StakingContract.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() payable returns()
func (_StakingContract *StakingContractTransactorSession) Withdraw() (*types.Transaction, error) {
	return _StakingContract.Contract.Withdraw(&_StakingContract.TransactOpts)
}

// StakingContractOnNewDepositIterator is returned from FilterOnNewDeposit and is used to iterate over the raw logs and unpacked data for OnNewDeposit events raised by the StakingContract contract.
type StakingContractOnNewDepositIterator struct {
	Event *StakingContractOnNewDeposit // Event containing the contract specifics and raw log

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
func (it *StakingContractOnNewDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingContractOnNewDeposit)
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
		it.Event = new(StakingContractOnNewDeposit)
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
func (it *StakingContractOnNewDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingContractOnNewDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingContractOnNewDeposit represents a OnNewDeposit event raised by the StakingContract contract.
type StakingContractOnNewDeposit struct {
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
func (_StakingContract *StakingContractFilterer) FilterOnNewDeposit(opts *bind.FilterOpts, sender []common.Address, validatorId [][32]byte, validatorAddress []common.Address) (*StakingContractOnNewDepositIterator, error) {

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

	logs, sub, err := _StakingContract.contract.FilterLogs(opts, "OnNewDeposit", senderRule, validatorIdRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &StakingContractOnNewDepositIterator{contract: _StakingContract.contract, event: "OnNewDeposit", logs: logs, sub: sub}, nil
}

// WatchOnNewDeposit is a free log subscription operation binding the contract event 0x9a1f4f083763f8508b19d4301c0110d2b47d99a8c5cf52c825c9e8cfea17f89c.
//
// Solidity: event OnNewDeposit(address indexed sender, bytes32 indexed validatorId, address indexed validatorAddress, bytes pubkey, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContract *StakingContractFilterer) WatchOnNewDeposit(opts *bind.WatchOpts, sink chan<- *StakingContractOnNewDeposit, sender []common.Address, validatorId [][32]byte, validatorAddress []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _StakingContract.contract.WatchLogs(opts, "OnNewDeposit", senderRule, validatorIdRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingContractOnNewDeposit)
				if err := _StakingContract.contract.UnpackLog(event, "OnNewDeposit", log); err != nil {
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
func (_StakingContract *StakingContractFilterer) ParseOnNewDeposit(log types.Log) (*StakingContractOnNewDeposit, error) {
	event := new(StakingContractOnNewDeposit)
	if err := _StakingContract.contract.UnpackLog(event, "OnNewDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingContractOnRewardDepositKeyIterator is returned from FilterOnRewardDepositKey and is used to iterate over the raw logs and unpacked data for OnRewardDepositKey events raised by the StakingContract contract.
type StakingContractOnRewardDepositKeyIterator struct {
	Event *StakingContractOnRewardDepositKey // Event containing the contract specifics and raw log

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
func (it *StakingContractOnRewardDepositKeyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingContractOnRewardDepositKey)
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
		it.Event = new(StakingContractOnRewardDepositKey)
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
func (it *StakingContractOnRewardDepositKeyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingContractOnRewardDepositKeyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingContractOnRewardDepositKey represents a OnRewardDepositKey event raised by the StakingContract contract.
type StakingContractOnRewardDepositKey struct {
	Sender      common.Address
	ValidatorId [32]byte
	Reward      common.Address
	Value       *big.Int
	BlockNumber *big.Int
	BlockTime   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOnRewardDepositKey is a free log retrieval operation binding the contract event 0xe0b518260035297556cfeb160ef4b66aed5ba1606403b996e4102fdd87e13366.
//
// Solidity: event OnRewardDepositKey(address sender, bytes32 validatorId, address reward, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContract *StakingContractFilterer) FilterOnRewardDepositKey(opts *bind.FilterOpts) (*StakingContractOnRewardDepositKeyIterator, error) {

	logs, sub, err := _StakingContract.contract.FilterLogs(opts, "OnRewardDepositKey")
	if err != nil {
		return nil, err
	}
	return &StakingContractOnRewardDepositKeyIterator{contract: _StakingContract.contract, event: "OnRewardDepositKey", logs: logs, sub: sub}, nil
}

// WatchOnRewardDepositKey is a free log subscription operation binding the contract event 0xe0b518260035297556cfeb160ef4b66aed5ba1606403b996e4102fdd87e13366.
//
// Solidity: event OnRewardDepositKey(address sender, bytes32 validatorId, address reward, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContract *StakingContractFilterer) WatchOnRewardDepositKey(opts *bind.WatchOpts, sink chan<- *StakingContractOnRewardDepositKey) (event.Subscription, error) {

	logs, sub, err := _StakingContract.contract.WatchLogs(opts, "OnRewardDepositKey")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingContractOnRewardDepositKey)
				if err := _StakingContract.contract.UnpackLog(event, "OnRewardDepositKey", log); err != nil {
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

// ParseOnRewardDepositKey is a log parse operation binding the contract event 0xe0b518260035297556cfeb160ef4b66aed5ba1606403b996e4102fdd87e13366.
//
// Solidity: event OnRewardDepositKey(address sender, bytes32 validatorId, address reward, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContract *StakingContractFilterer) ParseOnRewardDepositKey(log types.Log) (*StakingContractOnRewardDepositKey, error) {
	event := new(StakingContractOnRewardDepositKey)
	if err := _StakingContract.contract.UnpackLog(event, "OnRewardDepositKey", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingContractOnWithdrawKeyIterator is returned from FilterOnWithdrawKey and is used to iterate over the raw logs and unpacked data for OnWithdrawKey events raised by the StakingContract contract.
type StakingContractOnWithdrawKeyIterator struct {
	Event *StakingContractOnWithdrawKey // Event containing the contract specifics and raw log

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
func (it *StakingContractOnWithdrawKeyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingContractOnWithdrawKey)
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
		it.Event = new(StakingContractOnWithdrawKey)
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
func (it *StakingContractOnWithdrawKeyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingContractOnWithdrawKeyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingContractOnWithdrawKey represents a OnWithdrawKey event raised by the StakingContract contract.
type StakingContractOnWithdrawKey struct {
	Sender      common.Address
	Value       *big.Int
	BlockNumber *big.Int
	BlockTime   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOnWithdrawKey is a free log retrieval operation binding the contract event 0x4d4666331ec61727075c5624fde25f5510c566e528d0565f2a2263a23b70d81a.
//
// Solidity: event OnWithdrawKey(address sender, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContract *StakingContractFilterer) FilterOnWithdrawKey(opts *bind.FilterOpts) (*StakingContractOnWithdrawKeyIterator, error) {

	logs, sub, err := _StakingContract.contract.FilterLogs(opts, "OnWithdrawKey")
	if err != nil {
		return nil, err
	}
	return &StakingContractOnWithdrawKeyIterator{contract: _StakingContract.contract, event: "OnWithdrawKey", logs: logs, sub: sub}, nil
}

// WatchOnWithdrawKey is a free log subscription operation binding the contract event 0x4d4666331ec61727075c5624fde25f5510c566e528d0565f2a2263a23b70d81a.
//
// Solidity: event OnWithdrawKey(address sender, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_StakingContract *StakingContractFilterer) WatchOnWithdrawKey(opts *bind.WatchOpts, sink chan<- *StakingContractOnWithdrawKey) (event.Subscription, error) {

	logs, sub, err := _StakingContract.contract.WatchLogs(opts, "OnWithdrawKey")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingContractOnWithdrawKey)
				if err := _StakingContract.contract.UnpackLog(event, "OnWithdrawKey", log); err != nil {
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
func (_StakingContract *StakingContractFilterer) ParseOnWithdrawKey(log types.Log) (*StakingContractOnWithdrawKey, error) {
	event := new(StakingContractOnWithdrawKey)
	if err := _StakingContract.contract.UnpackLog(event, "OnWithdrawKey", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
