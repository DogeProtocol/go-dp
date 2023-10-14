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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldValidatorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newValidatorAddress\",\"type\":\"address\"}],\"name\":\"OnChangeValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"OnCompleteWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnIncreaseDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"OnInitiateWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnNewDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewardAmount\",\"type\":\"uint256\"}],\"name\":\"OnReward\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"slashedAmount\",\"type\":\"uint256\"}],\"name\":\"OnSlashing\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"rewardAmount\",\"type\":\"uint256\"}],\"name\":\"addDepositorReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"slashAmount\",\"type\":\"uint256\"}],\"name\":\"addDepositorSlashing\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newValidatorAddress\",\"type\":\"address\"}],\"name\":\"changeValidator\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeWithdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getBalanceOfDepositor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDepositorCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"getDepositorOfValidator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getDepositorRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getDepositorSlashings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getNetBalanceOfDepositor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalDepositedBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getValidatorOfDepositor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initiateWithdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"listValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"newDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60806040526000600255600060035534801561001a57600080fd5b50612c6e8061002a6000396000f3fe6080604052600436106100e85760003560e01c80637942317c1161008a578063e03ff7cb11610059578063e03ff7cb14610333578063f17bb4621461034a578063f6abfc7614610375578063ff9205ab14610391576100e8565b80637942317c14610265578063a7113fee146102a2578063b51d1d4f146102df578063c200baf9146102f6576100e8565b806368d4e544116100c657806368d4e544146101a45780636d727bd0146101cf578063731f750d1461020c57806377c06fdc14610228576100e8565b80632ca3c041146100ed5780634f4af09e1461012a57806351ca531714610167575b600080fd5b3480156100f957600080fd5b50610114600480360381019061010f9190612182565b6103bc565b6040516101219190612afd565b60405180910390f35b34801561013657600080fd5b50610151600480360381019061014c9190612182565b610405565b60405161015e9190612afd565b60405180910390f35b34801561017357600080fd5b5061018e600480360381019061018991906121ab565b6105f9565b60405161019b9190612afd565b60405180910390f35b3480156101b057600080fd5b506101b9610794565b6040516101c6919061287b565b60405180910390f35b3480156101db57600080fd5b506101f660048036038101906101f19190612182565b610822565b6040516102039190612860565b60405180910390f35b61022660048036038101906102219190612182565b610890565b005b34801561023457600080fd5b5061024f600480360381019061024a9190612182565b61101b565b60405161025c9190612afd565b60405180910390f35b34801561027157600080fd5b5061028c600480360381019061028791906121ab565b611064565b6040516102999190612afd565b60405180910390f35b3480156102ae57600080fd5b506102c960048036038101906102c49190612182565b6111ff565b6040516102d69190612860565b60405180910390f35b3480156102eb57600080fd5b506102f461126d565b005b34801561030257600080fd5b5061031d60048036038101906103189190612182565b6114dc565b60405161032a9190612afd565b60405180910390f35b34801561033f57600080fd5b50610348611525565b005b34801561035657600080fd5b5061035f6118ba565b60405161036c9190612afd565b60405180910390f35b61038f600480360381019061038a9190612182565b6118c4565b005b34801561039d57600080fd5b506103a6612106565b6040516103b39190612afd565b60405180910390f35b6000600b60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6000801515600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515141561046857600090506105f4565b6000600c60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205411156104b957600090506105f4565b600061054c600b60008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461211090919063ffffffff16565b9050600a60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054811161059e5760009150506105f4565b6105f0600a60008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020548261212c90919063ffffffff16565b9150505b919050565b60008073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610669576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610660906128dd565b60405180910390fd5b6106bb82600a60008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461211090919063ffffffff16565b600a60008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff167fcadc6c149d7c30ba433e0a526c9f018a1c4dc5b32099790e4dd9d78a93021810836040516107449190612afd565b60405180910390a2600a60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6060600080548060200260200160405190810160405280929190818152602001828054801561081857602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190600101908083116107ce575b5050505050905090565b600080600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905080915050919050565b600033905060003490506acecb8f27f4200f3a0000008110156108e8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108df906129fd565b60405180910390fd5b8273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610957576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161094e9061299d565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614156109c7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109be906129bd565b60405180910390fd5b60001515600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514610a5a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a5190612add565b60405180910390fd5b60001515600660008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514610aed576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ae49061297d565b60405180910390fd5b60008373ffffffffffffffffffffffffffffffffffffffff1631905060008114610b4c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b439061289d565b60405180910390fd5b60001515600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514610bdf576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bd690612a5d565b60405180910390fd5b60001515600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514610c72576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c69906129dd565b60405180910390fd5b6000849080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550610cea8260025461211090919063ffffffff16565b600281905550610d06600160035461211090919063ffffffff16565b60038190555081600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506001600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506001600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506001600660008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506001600760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555082600860008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555083600960008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508373ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fbe02029a5af0c964ebee7370f030cf18a026aae3a5d66f8107aee23f226d9ada84434260405161100d93929190612b18565b60405180910390a350505050565b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b60008073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146110d4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110cb906128dd565b60405180910390fd5b61112682600b60008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461211090919063ffffffff16565b600b60008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff167fd1072bb52c3131d0c96197b73fb8a45637e30f8b6664fc142310cc9b242859b4836040516111af9190612afd565b60405180910390a2600b60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600080600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905080915050919050565b600033905060011515600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514611305576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112fc906128bd565b60405180910390fd5b6000600c60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414611387576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161137e90612a9d565b60405180910390fd5b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205411611409576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161140090612abd565b60405180910390fd5b6203d0904301600c60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555060001515600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a905050507f47e2f9085249de3b62accafda3451074e283e2c6f30a39ae0b9952f3a0f8ecf7816040516114d19190612860565b60405180910390a150565b6000600a60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b60003390506000600c60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054116115ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016115a39061293d565b60405180910390fd5b600c60008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054431161162d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161162490612a7d565b60405180910390fd5b60003073ffffffffffffffffffffffffffffffffffffffff16634f4af09e836040518263ffffffff1660e01b81526004016116689190612860565b60206040518083038186803b15801561168057600080fd5b505afa158015611694573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116b891906121e7565b9050600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009055600b60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009055600a60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009055600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81549060ff021916905560008273ffffffffffffffffffffffffffffffffffffffff16826040516117f89061284b565b60006040518083038185875af1925050503d8060008114611835576040519150601f19603f3d011682016040523d82523d6000602084013e61183a565b606091505b505090508061187e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016118759061291d565b60405180910390fd5b7fda6373ecbed97803ca40cc1b7ed282476253b5aa7cd093dbb61d6990d5efcde4836040516118ad9190612860565b60405180910390a1505050565b6000600354905090565b60001515600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514611957576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161194e90612add565b60405180910390fd5b60001515600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515146119ea576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016119e1906128fd565b60405180910390fd5b60001515600660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514611a7d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611a7490612a3d565b60405180910390fd5b60001515600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514611b10576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611b0790612a1d565b60405180910390fd5b60008173ffffffffffffffffffffffffffffffffffffffff163114611b6a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611b619061289d565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415611bda576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611bd1906129bd565b60405180910390fd5b60003390508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415611c4e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611c459061299d565b60405180910390fd5b60011515600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514611ce1576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611cd8906128bd565b60405180910390fd5b6000600c60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414611d63576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611d5a9061295d565b60405180910390fd5b6001600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506001600660008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555080600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000829080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506000600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550600860008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690558273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f66d7a4dea74851a2dcc039f4c17dd2862081083c29daceb1a9346783de9185ce60405160405180910390a4505050565b6000600254905090565b60008082840190508381101561212257fe5b8091505092915050565b60008282111561213857fe5b818303905092915050565b60008135905061215281612be0565b92915050565b60008135905061216781612bf7565b92915050565b60008151905061217c81612bf7565b92915050565b60006020828403121561219457600080fd5b60006121a284828501612143565b91505092915050565b600080604083850312156121be57600080fd5b60006121cc85828601612143565b92505060206121dd85828601612158565b9150509250929050565b6000602082840312156121f957600080fd5b60006122078482850161216d565b91505092915050565b600061221c8383612228565b60208301905092915050565b61223181612ba4565b82525050565b61224081612ba4565b82525050565b600061225182612b5f565b61225b8185612b77565b935061226683612b4f565b8060005b8381101561229757815161227e8882612210565b975061228983612b6a565b92505060018101905061226a565b5085935050505092915050565b60006122b1602083612b93565b91507f76616c696461746f722062616c616e63652073686f756c64206265207a65726f6000830152602082019050919050565b60006122f1601883612b93565b91507f4465706f7369746f7220646f6573206e6f7420657869737400000000000000006000830152602082019050919050565b6000612331601983612b93565b91507f4f6e6c7920564d2063616c6c732061726520616c6c6f776564000000000000006000830152602082019050919050565b6000612371601883612b93565b91507f56616c696461746f722069732061206465706f7369746f7200000000000000006000830152602082019050919050565b60006123b1600f83612b93565b91507f5769746864726177206661696c656400000000000000000000000000000000006000830152602082019050919050565b60006123f1602b83612b93565b91507f4465706f7369746f72207769746864726177616c207265717565737420646f6560008301527f73206e6f742065786973740000000000000000000000000000000000000000006020830152604082019050919050565b6000612457601583612b93565b91507f5769746864726177616c2069732070656e64696e6700000000000000000000006000830152602082019050919050565b6000612497601683612b93565b91507f56616c696461746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b60006124d7603583612b93565b91507f4465706f7369746f7220616464726573732063616e6e6f742062652073616d6560008301527f2061732056616c696461746f72206164647265737300000000000000000000006020830152604082019050919050565b600061253d601183612b93565b91507f496e76616c69642076616c696461746f720000000000000000000000000000006000830152602082019050919050565b600061257d601683612b93565b91507f4465706f7369746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b60006125bd600083612b88565b9150600082019050919050565b60006125d7602b83612b93565b91507f4465706f73697420616d6f756e742062656c6f77206d696e696d756d2064657060008301527f6f73697420616d6f756e740000000000000000000000000000000000000000006020830152604082019050919050565b600061263d601983612b93565b91507f4465706f7369746f7220616c72656164792065786973746564000000000000006000830152602082019050919050565b600061267d601983612b93565b91507f56616c696461746f7220616c72656164792065786973746564000000000000006000830152602082019050919050565b60006126bd601883612b93565b91507f4465706f7369746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b60006126fd602483612b93565b91507f4465706f7369746f72207769746864726177616c20726571756573742070656e60008301527f64696e67000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000612763602383612b93565b91507f4465706f7369746f72207769746864726177616c20726571756573742065786960008301527f73747300000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b60006127c9601983612b93565b91507f4465706f7369746f722062616c616e6365206973207a65726f000000000000006000830152602082019050919050565b6000612809601883612b93565b91507f56616c696461746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b61284581612bd6565b82525050565b6000612856826125b0565b9150819050919050565b60006020820190506128756000830184612237565b92915050565b600060208201905081810360008301526128958184612246565b905092915050565b600060208201905081810360008301526128b6816122a4565b9050919050565b600060208201905081810360008301526128d6816122e4565b9050919050565b600060208201905081810360008301526128f681612324565b9050919050565b6000602082019050818103600083015261291681612364565b9050919050565b60006020820190508181036000830152612936816123a4565b9050919050565b60006020820190508181036000830152612956816123e4565b9050919050565b600060208201905081810360008301526129768161244a565b9050919050565b600060208201905081810360008301526129968161248a565b9050919050565b600060208201905081810360008301526129b6816124ca565b9050919050565b600060208201905081810360008301526129d681612530565b9050919050565b600060208201905081810360008301526129f681612570565b9050919050565b60006020820190508181036000830152612a16816125ca565b9050919050565b60006020820190508181036000830152612a3681612630565b9050919050565b60006020820190508181036000830152612a5681612670565b9050919050565b60006020820190508181036000830152612a76816126b0565b9050919050565b60006020820190508181036000830152612a96816126f0565b9050919050565b60006020820190508181036000830152612ab681612756565b9050919050565b60006020820190508181036000830152612ad6816127bc565b9050919050565b60006020820190508181036000830152612af6816127fc565b9050919050565b6000602082019050612b12600083018461283c565b92915050565b6000606082019050612b2d600083018661283c565b612b3a602083018561283c565b612b47604083018461283c565b949350505050565b6000819050602082019050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b6000612baf82612bb6565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b612be981612ba4565b8114612bf457600080fd5b50565b612c0081612bd6565b8114612c0b57600080fd5b5056fea2646970667358221220f333f40fa8db191d3820d25deb8869da7247ef7ecb11da1166e8b45b323811e264736f6c63782c302e372e362d646576656c6f702e323032332e31302e31322b636f6d6d69742e37333338323935662e6d6f64005d",
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

// GetDepositorRewards is a free data retrieval call binding the contract method 0x2ca3c041.
//
// Solidity: function getDepositorRewards(address depositorAddress) view returns(uint256)
func (_Staking *StakingCaller) GetDepositorRewards(opts *bind.CallOpts, depositorAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getDepositorRewards", depositorAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDepositorRewards is a free data retrieval call binding the contract method 0x2ca3c041.
//
// Solidity: function getDepositorRewards(address depositorAddress) view returns(uint256)
func (_Staking *StakingSession) GetDepositorRewards(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDepositorRewards(&_Staking.CallOpts, depositorAddress)
}

// GetDepositorRewards is a free data retrieval call binding the contract method 0x2ca3c041.
//
// Solidity: function getDepositorRewards(address depositorAddress) view returns(uint256)
func (_Staking *StakingCallerSession) GetDepositorRewards(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDepositorRewards(&_Staking.CallOpts, depositorAddress)
}

// GetDepositorSlashings is a free data retrieval call binding the contract method 0xc200baf9.
//
// Solidity: function getDepositorSlashings(address depositorAddress) view returns(uint256)
func (_Staking *StakingCaller) GetDepositorSlashings(opts *bind.CallOpts, depositorAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getDepositorSlashings", depositorAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDepositorSlashings is a free data retrieval call binding the contract method 0xc200baf9.
//
// Solidity: function getDepositorSlashings(address depositorAddress) view returns(uint256)
func (_Staking *StakingSession) GetDepositorSlashings(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDepositorSlashings(&_Staking.CallOpts, depositorAddress)
}

// GetDepositorSlashings is a free data retrieval call binding the contract method 0xc200baf9.
//
// Solidity: function getDepositorSlashings(address depositorAddress) view returns(uint256)
func (_Staking *StakingCallerSession) GetDepositorSlashings(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDepositorSlashings(&_Staking.CallOpts, depositorAddress)
}

// GetNetBalanceOfDepositor is a free data retrieval call binding the contract method 0x4f4af09e.
//
// Solidity: function getNetBalanceOfDepositor(address depositorAddress) view returns(uint256)
func (_Staking *StakingCaller) GetNetBalanceOfDepositor(opts *bind.CallOpts, depositorAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getNetBalanceOfDepositor", depositorAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNetBalanceOfDepositor is a free data retrieval call binding the contract method 0x4f4af09e.
//
// Solidity: function getNetBalanceOfDepositor(address depositorAddress) view returns(uint256)
func (_Staking *StakingSession) GetNetBalanceOfDepositor(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetNetBalanceOfDepositor(&_Staking.CallOpts, depositorAddress)
}

// GetNetBalanceOfDepositor is a free data retrieval call binding the contract method 0x4f4af09e.
//
// Solidity: function getNetBalanceOfDepositor(address depositorAddress) view returns(uint256)
func (_Staking *StakingCallerSession) GetNetBalanceOfDepositor(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetNetBalanceOfDepositor(&_Staking.CallOpts, depositorAddress)
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

// AddDepositorReward is a paid mutator transaction binding the contract method 0x7942317c.
//
// Solidity: function addDepositorReward(address depositorAddress, uint256 rewardAmount) returns(uint256)
func (_Staking *StakingTransactor) AddDepositorReward(opts *bind.TransactOpts, depositorAddress common.Address, rewardAmount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "addDepositorReward", depositorAddress, rewardAmount)
}

// AddDepositorReward is a paid mutator transaction binding the contract method 0x7942317c.
//
// Solidity: function addDepositorReward(address depositorAddress, uint256 rewardAmount) returns(uint256)
func (_Staking *StakingSession) AddDepositorReward(depositorAddress common.Address, rewardAmount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.AddDepositorReward(&_Staking.TransactOpts, depositorAddress, rewardAmount)
}

// AddDepositorReward is a paid mutator transaction binding the contract method 0x7942317c.
//
// Solidity: function addDepositorReward(address depositorAddress, uint256 rewardAmount) returns(uint256)
func (_Staking *StakingTransactorSession) AddDepositorReward(depositorAddress common.Address, rewardAmount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.AddDepositorReward(&_Staking.TransactOpts, depositorAddress, rewardAmount)
}

// AddDepositorSlashing is a paid mutator transaction binding the contract method 0x51ca5317.
//
// Solidity: function addDepositorSlashing(address depositorAddress, uint256 slashAmount) returns(uint256)
func (_Staking *StakingTransactor) AddDepositorSlashing(opts *bind.TransactOpts, depositorAddress common.Address, slashAmount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "addDepositorSlashing", depositorAddress, slashAmount)
}

// AddDepositorSlashing is a paid mutator transaction binding the contract method 0x51ca5317.
//
// Solidity: function addDepositorSlashing(address depositorAddress, uint256 slashAmount) returns(uint256)
func (_Staking *StakingSession) AddDepositorSlashing(depositorAddress common.Address, slashAmount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.AddDepositorSlashing(&_Staking.TransactOpts, depositorAddress, slashAmount)
}

// AddDepositorSlashing is a paid mutator transaction binding the contract method 0x51ca5317.
//
// Solidity: function addDepositorSlashing(address depositorAddress, uint256 slashAmount) returns(uint256)
func (_Staking *StakingTransactorSession) AddDepositorSlashing(depositorAddress common.Address, slashAmount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.AddDepositorSlashing(&_Staking.TransactOpts, depositorAddress, slashAmount)
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
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterOnChangeValidator is a free log retrieval operation binding the contract event 0x66d7a4dea74851a2dcc039f4c17dd2862081083c29daceb1a9346783de9185ce.
//
// Solidity: event OnChangeValidator(address indexed depositorAddress, address indexed oldValidatorAddress, address indexed newValidatorAddress)
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

// WatchOnChangeValidator is a free log subscription operation binding the contract event 0x66d7a4dea74851a2dcc039f4c17dd2862081083c29daceb1a9346783de9185ce.
//
// Solidity: event OnChangeValidator(address indexed depositorAddress, address indexed oldValidatorAddress, address indexed newValidatorAddress)
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

// ParseOnChangeValidator is a log parse operation binding the contract event 0x66d7a4dea74851a2dcc039f4c17dd2862081083c29daceb1a9346783de9185ce.
//
// Solidity: event OnChangeValidator(address indexed depositorAddress, address indexed oldValidatorAddress, address indexed newValidatorAddress)
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
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnCompleteWithdrawal is a free log retrieval operation binding the contract event 0xda6373ecbed97803ca40cc1b7ed282476253b5aa7cd093dbb61d6990d5efcde4.
//
// Solidity: event OnCompleteWithdrawal(address depositorAddress)
func (_Staking *StakingFilterer) FilterOnCompleteWithdrawal(opts *bind.FilterOpts) (*StakingOnCompleteWithdrawalIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnCompleteWithdrawal")
	if err != nil {
		return nil, err
	}
	return &StakingOnCompleteWithdrawalIterator{contract: _Staking.contract, event: "OnCompleteWithdrawal", logs: logs, sub: sub}, nil
}

// WatchOnCompleteWithdrawal is a free log subscription operation binding the contract event 0xda6373ecbed97803ca40cc1b7ed282476253b5aa7cd093dbb61d6990d5efcde4.
//
// Solidity: event OnCompleteWithdrawal(address depositorAddress)
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

// ParseOnCompleteWithdrawal is a log parse operation binding the contract event 0xda6373ecbed97803ca40cc1b7ed282476253b5aa7cd093dbb61d6990d5efcde4.
//
// Solidity: event OnCompleteWithdrawal(address depositorAddress)
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
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnInitiateWithdrawal is a free log retrieval operation binding the contract event 0x47e2f9085249de3b62accafda3451074e283e2c6f30a39ae0b9952f3a0f8ecf7.
//
// Solidity: event OnInitiateWithdrawal(address depositorAddress)
func (_Staking *StakingFilterer) FilterOnInitiateWithdrawal(opts *bind.FilterOpts) (*StakingOnInitiateWithdrawalIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnInitiateWithdrawal")
	if err != nil {
		return nil, err
	}
	return &StakingOnInitiateWithdrawalIterator{contract: _Staking.contract, event: "OnInitiateWithdrawal", logs: logs, sub: sub}, nil
}

// WatchOnInitiateWithdrawal is a free log subscription operation binding the contract event 0x47e2f9085249de3b62accafda3451074e283e2c6f30a39ae0b9952f3a0f8ecf7.
//
// Solidity: event OnInitiateWithdrawal(address depositorAddress)
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

// ParseOnInitiateWithdrawal is a log parse operation binding the contract event 0x47e2f9085249de3b62accafda3451074e283e2c6f30a39ae0b9952f3a0f8ecf7.
//
// Solidity: event OnInitiateWithdrawal(address depositorAddress)
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

// StakingOnRewardIterator is returned from FilterOnReward and is used to iterate over the raw logs and unpacked data for OnReward events raised by the Staking contract.
type StakingOnRewardIterator struct {
	Event *StakingOnReward // Event containing the contract specifics and raw log

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
func (it *StakingOnRewardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOnReward)
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
		it.Event = new(StakingOnReward)
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
func (it *StakingOnRewardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOnRewardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOnReward represents a OnReward event raised by the Staking contract.
type StakingOnReward struct {
	DepositorAddress common.Address
	RewardAmount     *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnReward is a free log retrieval operation binding the contract event 0xd1072bb52c3131d0c96197b73fb8a45637e30f8b6664fc142310cc9b242859b4.
//
// Solidity: event OnReward(address indexed depositorAddress, uint256 rewardAmount)
func (_Staking *StakingFilterer) FilterOnReward(opts *bind.FilterOpts, depositorAddress []common.Address) (*StakingOnRewardIterator, error) {

	var depositorAddressRule []interface{}
	for _, depositorAddressItem := range depositorAddress {
		depositorAddressRule = append(depositorAddressRule, depositorAddressItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnReward", depositorAddressRule)
	if err != nil {
		return nil, err
	}
	return &StakingOnRewardIterator{contract: _Staking.contract, event: "OnReward", logs: logs, sub: sub}, nil
}

// WatchOnReward is a free log subscription operation binding the contract event 0xd1072bb52c3131d0c96197b73fb8a45637e30f8b6664fc142310cc9b242859b4.
//
// Solidity: event OnReward(address indexed depositorAddress, uint256 rewardAmount)
func (_Staking *StakingFilterer) WatchOnReward(opts *bind.WatchOpts, sink chan<- *StakingOnReward, depositorAddress []common.Address) (event.Subscription, error) {

	var depositorAddressRule []interface{}
	for _, depositorAddressItem := range depositorAddress {
		depositorAddressRule = append(depositorAddressRule, depositorAddressItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OnReward", depositorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOnReward)
				if err := _Staking.contract.UnpackLog(event, "OnReward", log); err != nil {
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

// ParseOnReward is a log parse operation binding the contract event 0xd1072bb52c3131d0c96197b73fb8a45637e30f8b6664fc142310cc9b242859b4.
//
// Solidity: event OnReward(address indexed depositorAddress, uint256 rewardAmount)
func (_Staking *StakingFilterer) ParseOnReward(log types.Log) (*StakingOnReward, error) {
	event := new(StakingOnReward)
	if err := _Staking.contract.UnpackLog(event, "OnReward", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingOnSlashingIterator is returned from FilterOnSlashing and is used to iterate over the raw logs and unpacked data for OnSlashing events raised by the Staking contract.
type StakingOnSlashingIterator struct {
	Event *StakingOnSlashing // Event containing the contract specifics and raw log

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
func (it *StakingOnSlashingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOnSlashing)
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
		it.Event = new(StakingOnSlashing)
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
func (it *StakingOnSlashingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOnSlashingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOnSlashing represents a OnSlashing event raised by the Staking contract.
type StakingOnSlashing struct {
	DepositorAddress common.Address
	SlashedAmount    *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnSlashing is a free log retrieval operation binding the contract event 0xcadc6c149d7c30ba433e0a526c9f018a1c4dc5b32099790e4dd9d78a93021810.
//
// Solidity: event OnSlashing(address indexed depositorAddress, uint256 slashedAmount)
func (_Staking *StakingFilterer) FilterOnSlashing(opts *bind.FilterOpts, depositorAddress []common.Address) (*StakingOnSlashingIterator, error) {

	var depositorAddressRule []interface{}
	for _, depositorAddressItem := range depositorAddress {
		depositorAddressRule = append(depositorAddressRule, depositorAddressItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnSlashing", depositorAddressRule)
	if err != nil {
		return nil, err
	}
	return &StakingOnSlashingIterator{contract: _Staking.contract, event: "OnSlashing", logs: logs, sub: sub}, nil
}

// WatchOnSlashing is a free log subscription operation binding the contract event 0xcadc6c149d7c30ba433e0a526c9f018a1c4dc5b32099790e4dd9d78a93021810.
//
// Solidity: event OnSlashing(address indexed depositorAddress, uint256 slashedAmount)
func (_Staking *StakingFilterer) WatchOnSlashing(opts *bind.WatchOpts, sink chan<- *StakingOnSlashing, depositorAddress []common.Address) (event.Subscription, error) {

	var depositorAddressRule []interface{}
	for _, depositorAddressItem := range depositorAddress {
		depositorAddressRule = append(depositorAddressRule, depositorAddressItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OnSlashing", depositorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOnSlashing)
				if err := _Staking.contract.UnpackLog(event, "OnSlashing", log); err != nil {
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

// ParseOnSlashing is a log parse operation binding the contract event 0xcadc6c149d7c30ba433e0a526c9f018a1c4dc5b32099790e4dd9d78a93021810.
//
// Solidity: event OnSlashing(address indexed depositorAddress, uint256 slashedAmount)
func (_Staking *StakingFilterer) ParseOnSlashing(log types.Log) (*StakingOnSlashing, error) {
	event := new(StakingOnSlashing)
	if err := _Staking.contract.UnpackLog(event, "OnSlashing", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
