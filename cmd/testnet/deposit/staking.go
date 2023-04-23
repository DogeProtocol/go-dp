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

// StakingMetaData contains all meta data concerning the Staking contract.
var StakingMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"validatorId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnNewDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnWithdrawKey\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"depositBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"getDepositor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"listValidator\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyhash\",\"type\":\"bytes32\"}],\"name\":\"newDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalDepositBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405260006000600050909055600060016000509090553480156100255760006000fd5b5061002b565b610d958061003a6000396000f3fe6080604052600436106100745760003560e01c806375697e661161004e57806375697e661461010e578063993aeb851461013a578063e8c0a0df14610156578063fba13bd01461018257610074565b80632dfdf0b51461007a5780632e1a7d4d146100a65780636e2baf48146100d057610074565b60006000fd5b3480156100875760006000fd5b506100906101c0565b60405161009d9190610be7565b60405180910390f35b3480156100b35760006000fd5b506100ce60048036038101906100c991906109ba565b6101d2565b005b3480156100dd5760006000fd5b506100f860048036038101906100f39190610964565b6103ab565b6040516101059190610b20565b60405180910390f35b34801561011b5760006000fd5b50610124610414565b6040516101319190610b82565b60405180910390f35b610154600480360381019061014f919061098f565b6104aa565b005b3480156101635760006000fd5b5061016c610837565b6040516101799190610be7565b60405180910390f35b34801561018f5760006000fd5b506101aa60048036038101906101a59190610964565b610849565b6040516101b79190610be7565b60405180910390f35b600060006000505490506101cf565b90565b80600260005060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050541015151561025c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161025390610ba5565b60405180910390fd5b6102748160016000505461089d90919063ffffffff16565b60016000508190909055506102d781600260005060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000505461089d90919063ffffffff16565b600260005060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000508190909055503373ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610369573d600060003e3d6000fd5b507f4d4666331ec61727075c5624fde25f5510c566e528d0565f2a2263a23b70d81a3382434260405161039f9493929190610b3c565b60405180910390a15b50565b600060006103be836108bb63ffffffff16565b9050600060036000506000836000191660001916815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050809250505061040f5650505b919050565b6060600560005080548060200260200160405190810160405280929190818152602001828054801561049b57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610451575b505050505090506104a7565b90565b3373ffffffffffffffffffffffffffffffffffffffff1660036000506000600460005060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050546000191660001916815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415151561059e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161059590610bc6565b60405180910390fd5b6105b760016000600050546108e790919063ffffffff16565b60006000508190909055506105da346001600050546108e790919063ffffffff16565b600160005081909090555061063d34600260005060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050546108e790919063ffffffff16565b600260005060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005081909090555060006106998261090c63ffffffff16565b905060006106ac826108bb63ffffffff16565b90503360036000506000836000191660001916815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600460005060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005081909060001916905550600560005082908060018154018082558091505060019003906000526020600020900160005b9091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff1681600019163373ffffffffffffffffffffffffffffffffffffffff167fef144c8f165094acc8d6e390b954b1644ba9de277d7be885d4511c39c5124d0b34434260405161082993929190610c03565b60405180910390a450505b50565b60006001600050549050610846565b90565b6000600260005060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050549050610898565b919050565b60008282111515156108ab57fe5b81830390506108b5565b92915050565b600060608273ffffffffffffffffffffffffffffffffffffffff16901b60001b90506108e2565b919050565b6000600082840190508381101515156108fc57fe5b8091505061090656505b92915050565b60008160001c9050610919565b91905056610d5e565b60008135905061093181610d0d565b5b92915050565b60008135905061094781610d28565b5b92915050565b60008135905061095d81610d43565b5b92915050565b6000602082840312156109775760006000fd5b600061098584828501610922565b9150505b92915050565b6000602082840312156109a25760006000fd5b60006109b084828501610938565b9150505b92915050565b6000602082840312156109cd5760006000fd5b60006109db8482850161094e565b9150505b92915050565b60006109f18383610a0e565b6020830190505b92915050565b610a0781610cd4565b82525b5050565b610a1781610c8a565b82525b5050565b610a2781610c8a565b82525b5050565b6000610a3982610c4c565b610a438185610c66565b9350610a4e83610c3b565b8060005b83811015610a80578151610a6688826109e5565b9750610a7183610c58565b9250505b600181019050610a52565b508593505050505b92915050565b6000610a9b601283610c78565b91507f496e73756666696369656e742066756e6473000000000000000000000000000060008301526020820190505b919050565b6000610adc601583610c78565b91507f53656e64657220616c726561647920657869737473000000000000000000000060008301526020820190505b919050565b610b1981610cc9565b82525b5050565b6000602082019050610b356000830184610a1e565b5b92915050565b6000608082019050610b5160008301876109fe565b610b5e6020830186610b10565b610b6b6040830185610b10565b610b786060830184610b10565b5b95945050505050565b60006020820190508181036000830152610b9c8184610a2e565b90505b92915050565b60006020820190508181036000830152610bbe81610a8e565b90505b919050565b60006020820190508181036000830152610bdf81610acf565b90505b919050565b6000602082019050610bfc6000830184610b10565b5b92915050565b6000606082019050610c186000830186610b10565b610c256020830185610b10565b610c326040830184610b10565b5b949350505050565b60008190506020820190505b919050565b6000815190505b919050565b60006020820190505b919050565b60008282526020820190505b92915050565b60008282526020820190505b92915050565b6000610c9582610ca8565b90505b919050565b60008190505b919050565b600073ffffffffffffffffffffffffffffffffffffffff821690505b919050565b60008190505b919050565b6000610cdf82610ce7565b90505b919050565b6000610cf282610cfa565b90505b919050565b6000610d0582610ca8565b90505b919050565b610d1681610c8a565b81141515610d245760006000fd5b5b50565b610d3181610c9d565b81141515610d3f5760006000fd5b5b50565b610d4c81610cc9565b81141515610d5a5760006000fd5b5b50565bfea26469706673582212205826fa46de57a867761296f414a31bc9be71d3c82e965d592ac0e34b5dd5b78864736f6c63430007060033",
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

// DepositBalanceOf is a free data retrieval call binding the contract method 0xfba13bd0.
//
// Solidity: function depositBalanceOf(address depositor) view returns(uint256)
func (_Staking *StakingCaller) DepositBalanceOf(opts *bind.CallOpts, depositor common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "depositBalanceOf", depositor)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositBalanceOf is a free data retrieval call binding the contract method 0xfba13bd0.
//
// Solidity: function depositBalanceOf(address depositor) view returns(uint256)
func (_Staking *StakingSession) DepositBalanceOf(depositor common.Address) (*big.Int, error) {
	return _Staking.Contract.DepositBalanceOf(&_Staking.CallOpts, depositor)
}

// DepositBalanceOf is a free data retrieval call binding the contract method 0xfba13bd0.
//
// Solidity: function depositBalanceOf(address depositor) view returns(uint256)
func (_Staking *StakingCallerSession) DepositBalanceOf(depositor common.Address) (*big.Int, error) {
	return _Staking.Contract.DepositBalanceOf(&_Staking.CallOpts, depositor)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Staking *StakingCaller) DepositCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "depositCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Staking *StakingSession) DepositCount() (*big.Int, error) {
	return _Staking.Contract.DepositCount(&_Staking.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Staking *StakingCallerSession) DepositCount() (*big.Int, error) {
	return _Staking.Contract.DepositCount(&_Staking.CallOpts)
}

// GetDepositor is a free data retrieval call binding the contract method 0x6e2baf48.
//
// Solidity: function getDepositor(address validator) view returns(address)
func (_Staking *StakingCaller) GetDepositor(opts *bind.CallOpts, validator common.Address) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getDepositor", validator)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDepositor is a free data retrieval call binding the contract method 0x6e2baf48.
//
// Solidity: function getDepositor(address validator) view returns(address)
func (_Staking *StakingSession) GetDepositor(validator common.Address) (common.Address, error) {
	return _Staking.Contract.GetDepositor(&_Staking.CallOpts, validator)
}

// GetDepositor is a free data retrieval call binding the contract method 0x6e2baf48.
//
// Solidity: function getDepositor(address validator) view returns(address)
func (_Staking *StakingCallerSession) GetDepositor(validator common.Address) (common.Address, error) {
	return _Staking.Contract.GetDepositor(&_Staking.CallOpts, validator)
}

// ListValidator is a free data retrieval call binding the contract method 0x75697e66.
//
// Solidity: function listValidator() view returns(address[])
func (_Staking *StakingCaller) ListValidator(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "listValidator")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ListValidator is a free data retrieval call binding the contract method 0x75697e66.
//
// Solidity: function listValidator() view returns(address[])
func (_Staking *StakingSession) ListValidator() ([]common.Address, error) {
	return _Staking.Contract.ListValidator(&_Staking.CallOpts)
}

// ListValidator is a free data retrieval call binding the contract method 0x75697e66.
//
// Solidity: function listValidator() view returns(address[])
func (_Staking *StakingCallerSession) ListValidator() ([]common.Address, error) {
	return _Staking.Contract.ListValidator(&_Staking.CallOpts)
}

// TotalDepositBalance is a free data retrieval call binding the contract method 0xe8c0a0df.
//
// Solidity: function totalDepositBalance() view returns(uint256)
func (_Staking *StakingCaller) TotalDepositBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "totalDepositBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalDepositBalance is a free data retrieval call binding the contract method 0xe8c0a0df.
//
// Solidity: function totalDepositBalance() view returns(uint256)
func (_Staking *StakingSession) TotalDepositBalance() (*big.Int, error) {
	return _Staking.Contract.TotalDepositBalance(&_Staking.CallOpts)
}

// TotalDepositBalance is a free data retrieval call binding the contract method 0xe8c0a0df.
//
// Solidity: function totalDepositBalance() view returns(uint256)
func (_Staking *StakingCallerSession) TotalDepositBalance() (*big.Int, error) {
	return _Staking.Contract.TotalDepositBalance(&_Staking.CallOpts)
}

// NewDeposit is a paid mutator transaction binding the contract method 0x993aeb85.
//
// Solidity: function newDeposit(bytes32 keyhash) payable returns()
func (_Staking *StakingTransactor) NewDeposit(opts *bind.TransactOpts, keyhash [32]byte) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "newDeposit", keyhash)
}

// NewDeposit is a paid mutator transaction binding the contract method 0x993aeb85.
//
// Solidity: function newDeposit(bytes32 keyhash) payable returns()
func (_Staking *StakingSession) NewDeposit(keyhash [32]byte) (*types.Transaction, error) {
	return _Staking.Contract.NewDeposit(&_Staking.TransactOpts, keyhash)
}

// NewDeposit is a paid mutator transaction binding the contract method 0x993aeb85.
//
// Solidity: function newDeposit(bytes32 keyhash) payable returns()
func (_Staking *StakingTransactorSession) NewDeposit(keyhash [32]byte) (*types.Transaction, error) {
	return _Staking.Contract.NewDeposit(&_Staking.TransactOpts, keyhash)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns()
func (_Staking *StakingTransactor) Withdraw(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdraw", value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns()
func (_Staking *StakingSession) Withdraw(value *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Withdraw(&_Staking.TransactOpts, value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns()
func (_Staking *StakingTransactorSession) Withdraw(value *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Withdraw(&_Staking.TransactOpts, value)
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
	Sender           common.Address
	ValidatorId      [32]byte
	ValidatorAddress common.Address
	Value            *big.Int
	BlockNumber      *big.Int
	BlockTime        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnNewDeposit is a free log retrieval operation binding the contract event 0xef144c8f165094acc8d6e390b954b1644ba9de277d7be885d4511c39c5124d0b.
//
// Solidity: event OnNewDeposit(address indexed sender, bytes32 indexed validatorId, address indexed validatorAddress, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) FilterOnNewDeposit(opts *bind.FilterOpts, sender []common.Address, validatorId [][32]byte, validatorAddress []common.Address) (*StakingOnNewDepositIterator, error) {

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

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnNewDeposit", senderRule, validatorIdRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &StakingOnNewDepositIterator{contract: _Staking.contract, event: "OnNewDeposit", logs: logs, sub: sub}, nil
}

// WatchOnNewDeposit is a free log subscription operation binding the contract event 0xef144c8f165094acc8d6e390b954b1644ba9de277d7be885d4511c39c5124d0b.
//
// Solidity: event OnNewDeposit(address indexed sender, bytes32 indexed validatorId, address indexed validatorAddress, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) WatchOnNewDeposit(opts *bind.WatchOpts, sink chan<- *StakingOnNewDeposit, sender []common.Address, validatorId [][32]byte, validatorAddress []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OnNewDeposit", senderRule, validatorIdRule, validatorAddressRule)
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

// ParseOnNewDeposit is a log parse operation binding the contract event 0xef144c8f165094acc8d6e390b954b1644ba9de277d7be885d4511c39c5124d0b.
//
// Solidity: event OnNewDeposit(address indexed sender, bytes32 indexed validatorId, address indexed validatorAddress, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) ParseOnNewDeposit(log types.Log) (*StakingOnNewDeposit, error) {
	event := new(StakingOnNewDeposit)
	if err := _Staking.contract.UnpackLog(event, "OnNewDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingOnWithdrawKeyIterator is returned from FilterOnWithdrawKey and is used to iterate over the raw logs and unpacked data for OnWithdrawKey events raised by the Staking contract.
type StakingOnWithdrawKeyIterator struct {
	Event *StakingOnWithdrawKey // Event containing the contract specifics and raw log

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
func (it *StakingOnWithdrawKeyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOnWithdrawKey)
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
		it.Event = new(StakingOnWithdrawKey)
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
func (it *StakingOnWithdrawKeyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOnWithdrawKeyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOnWithdrawKey represents a OnWithdrawKey event raised by the Staking contract.
type StakingOnWithdrawKey struct {
	Sender      common.Address
	Value       *big.Int
	BlockNumber *big.Int
	BlockTime   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOnWithdrawKey is a free log retrieval operation binding the contract event 0x4d4666331ec61727075c5624fde25f5510c566e528d0565f2a2263a23b70d81a.
//
// Solidity: event OnWithdrawKey(address sender, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) FilterOnWithdrawKey(opts *bind.FilterOpts) (*StakingOnWithdrawKeyIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnWithdrawKey")
	if err != nil {
		return nil, err
	}
	return &StakingOnWithdrawKeyIterator{contract: _Staking.contract, event: "OnWithdrawKey", logs: logs, sub: sub}, nil
}

// WatchOnWithdrawKey is a free log subscription operation binding the contract event 0x4d4666331ec61727075c5624fde25f5510c566e528d0565f2a2263a23b70d81a.
//
// Solidity: event OnWithdrawKey(address sender, uint256 value, uint256 blockNumber, uint256 blockTime)
func (_Staking *StakingFilterer) WatchOnWithdrawKey(opts *bind.WatchOpts, sink chan<- *StakingOnWithdrawKey) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OnWithdrawKey")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOnWithdrawKey)
				if err := _Staking.contract.UnpackLog(event, "OnWithdrawKey", log); err != nil {
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
func (_Staking *StakingFilterer) ParseOnWithdrawKey(log types.Log) (*StakingOnWithdrawKey, error) {
	event := new(StakingOnWithdrawKey)
	if err := _Staking.contract.UnpackLog(event, "OnWithdrawKey", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
