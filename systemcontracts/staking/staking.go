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
	Bin: "0x60806040526000600255600060035534801561001a57600080fd5b506130f98061002a6000396000f3fe6080604052600436106100e85760003560e01c80637942317c1161008a578063e03ff7cb11610059578063e03ff7cb14610333578063f17bb4621461034a578063f6abfc7614610375578063ff9205ab14610391576100e8565b80637942317c14610265578063a7113fee146102a2578063b51d1d4f146102df578063c200baf9146102f6576100e8565b806368d4e544116100c657806368d4e544146101a45780636d727bd0146101cf578063731f750d1461020c57806377c06fdc14610228576100e8565b80632ca3c041146100ed5780634f4af09e1461012a57806351ca531714610167575b600080fd5b3480156100f957600080fd5b50610114600480360381019061010f919061260a565b6103bc565b6040516101219190612f85565b60405180910390f35b34801561013657600080fd5b50610151600480360381019061014c919061260a565b610415565b60405161015e9190612f85565b60405180910390f35b34801561017357600080fd5b5061018e60048036038101906101899190612633565b610669565b60405161019b9190612f85565b60405180910390f35b3480156101b057600080fd5b506101b961084c565b6040516101c69190612d03565b60405180910390f35b3480156101db57600080fd5b506101f660048036038101906101f1919061260a565b6108ea565b6040516102039190612ce8565b60405180910390f35b6102266004803603810190610221919061260a565b610970565b005b34801561023457600080fd5b5061024f600480360381019061024a919061260a565b611213565b60405161025c9190612f85565b60405180910390f35b34801561027157600080fd5b5061028c60048036038101906102879190612633565b61126c565b6040516102999190612f85565b60405180910390f35b3480156102ae57600080fd5b506102c960048036038101906102c4919061260a565b61144f565b6040516102d69190612ce8565b60405180910390f35b3480156102eb57600080fd5b506102f46114d5565b005b34801561030257600080fd5b5061031d6004803603810190610318919061260a565b611794565b60405161032a9190612f85565b60405180910390f35b34801561033f57600080fd5b506103486117ed565b005b34801561035657600080fd5b5061035f611bf2565b60405161036c9190612f85565b60405180910390f35b61038f600480360381019061038a919061260a565b611bfc565b005b34801561039d57600080fd5b506103a661258e565b6040516103b39190612f85565b60405180910390f35b6000600b6000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600080151560056000847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514156104885760009050610664565b6000600c6000847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205411156104e95760009050610664565b600061059c600b6000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205460016000867bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461259890919063ffffffff16565b9050600a6000847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205481116105fe576000915050610664565b610660600a6000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054826125b490919063ffffffff16565b9150505b919050565b6000807bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16337bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16146106e9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106e090612d65565b60405180910390fd5b61074b82600a6000867bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461259890919063ffffffff16565b600a6000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fcadc6c149d7c30ba433e0a526c9f018a1c4dc5b32099790e4dd9d78a93021810836040516107ec9190612f85565b60405180910390a2600a6000847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b606060008054806020026020016040519081016040528092919081815260200182805480156108e057602002820191906000526020600020905b8160009054906101000a90047bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610886575b5050505050905090565b60008060086000847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a90047bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16905080915050919050565b600033905060003490506acecb8f27f4200f3a0000008110156109c8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109bf90612e85565b60405180910390fd5b827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff161415610a47576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a3e90612e25565b60405180910390fd5b60007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff161415610ac7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610abe90612e45565b60405180910390fd5b6000151560046000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514610b6a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b6190612f65565b60405180910390fd5b6000151560066000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514610c0d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c0490612e05565b60405180910390fd5b6000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1631905060008114610c74576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c6b90612d25565b60405180910390fd5b6000151560056000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514610d17576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d0e90612ee5565b60405180910390fd5b6000151560076000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514610dba576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610db190612e65565b60405180910390fd5b6000849080600181540180825580915050600190039060005260206000200160009091909190916101000a8154817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff02191690837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff160217905550610e428260025461259890919063ffffffff16565b600281905550610e5e600160035461259890919063ffffffff16565b6003819055508160016000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550600160046000867bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550600160056000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550600160066000867bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550600160076000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508260086000867bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a8154817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff02191690837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1602179055508360096000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a8154817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff02191690837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff160217905550837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fbe02029a5af0c964ebee7370f030cf18a026aae3a5d66f8107aee23f226d9ada84434260405161120593929190612fa0565b60405180910390a350505050565b600060016000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6000807bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16337bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16146112ec576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112e390612d65565b60405180910390fd5b61134e82600b6000867bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461259890919063ffffffff16565b600b6000857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fd1072bb52c3131d0c96197b73fb8a45637e30f8b6664fc142310cc9b242859b4836040516113ef9190612f85565b60405180910390a2600b6000847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b60008060086000847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a90047bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16905080915050919050565b60003390506001151560056000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615151461157d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161157490612d45565b60405180910390fd5b6000600c6000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541461160f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161160690612f25565b60405180910390fd5b600060016000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054116116a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161169890612f45565b60405180910390fd5b6203d0904301600c6000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506000151560056000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a905050507f47e2f9085249de3b62accafda3451074e283e2c6f30a39ae0b9952f3a0f8ecf7816040516117899190612ce8565b60405180910390a150565b6000600a6000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b60003390506000600c6000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205411611884576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161187b90612dc5565b60405180910390fd5b600c6000827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020544311611915576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161190c90612f05565b60405180910390fd5b6000307bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16634f4af09e836040518263ffffffff1660e01b81526004016119589190612ce8565b60206040518083038186803b15801561197057600080fd5b505afa158015611984573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119a8919061266f565b905060016000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009055600b6000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009055600a6000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000905560056000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81549060ff02191690556000827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1682604051611b3090612cd3565b60006040518083038185875af1925050503d8060008114611b6d576040519150601f19603f3d011682016040523d82523d6000602084013e611b72565b606091505b5050905080611bb6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611bad90612da5565b60405180910390fd5b7fda6373ecbed97803ca40cc1b7ed282476253b5aa7cd093dbb61d6990d5efcde483604051611be59190612ce8565b60405180910390a1505050565b6000600354905090565b6000151560046000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514611c9f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611c9690612f65565b60405180910390fd5b6000151560056000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514611d42576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611d3990612d85565b60405180910390fd5b6000151560066000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514611de5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611ddc90612ec5565b60405180910390fd5b6000151560076000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514611e88576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611e7f90612ea5565b60405180910390fd5b6000817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff163114611eea576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611ee190612d25565b60405180910390fd5b60007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff161415611f6a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611f6190612e45565b60405180910390fd5b6000339050817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff161415611fee576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611fe590612e25565b60405180910390fd5b6001151560056000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151514612091576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161208890612d45565b60405180910390fd5b6000600c6000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414612123576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161211a90612de5565b60405180910390fd5b600160046000847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550600160066000847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508060086000847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a8154817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff02191690837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1602179055508160096000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a8154817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff02191690837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1602179055506000829080600181540180825580915050600190039060005260206000200160009091909190916101000a8154817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff02191690837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff160217905550600060096000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a90047bffffffffffffffffffffffffffffffffffffffffffffffffffffffff169050600060046000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555060086000827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a8154907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff0219169055827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f66d7a4dea74851a2dcc039f4c17dd2862081083c29daceb1a9346783de9185ce60405160405180910390a4505050565b6000600254905090565b6000808284019050838110156125aa57fe5b8091505092915050565b6000828211156125c057fe5b818303905092915050565b6000813590506125da81613070565b92915050565b6000813590506125ef81613087565b92915050565b60008151905061260481613087565b92915050565b60006020828403121561261c57600080fd5b600061262a848285016125cb565b91505092915050565b6000806040838503121561264657600080fd5b6000612654858286016125cb565b9250506020612665858286016125e0565b9150509250929050565b60006020828403121561268157600080fd5b600061268f848285016125f5565b91505092915050565b60006126a483836126b0565b60208301905092915050565b6126b98161302c565b82525050565b6126c88161302c565b82525050565b60006126d982612fe7565b6126e38185612fff565b93506126ee83612fd7565b8060005b8381101561271f5781516127068882612698565b975061271183612ff2565b9250506001810190506126f2565b5085935050505092915050565b600061273960208361301b565b91507f76616c696461746f722062616c616e63652073686f756c64206265207a65726f6000830152602082019050919050565b600061277960188361301b565b91507f4465706f7369746f7220646f6573206e6f7420657869737400000000000000006000830152602082019050919050565b60006127b960198361301b565b91507f4f6e6c7920564d2063616c6c732061726520616c6c6f776564000000000000006000830152602082019050919050565b60006127f960188361301b565b91507f56616c696461746f722069732061206465706f7369746f7200000000000000006000830152602082019050919050565b6000612839600f8361301b565b91507f5769746864726177206661696c656400000000000000000000000000000000006000830152602082019050919050565b6000612879602b8361301b565b91507f4465706f7369746f72207769746864726177616c207265717565737420646f6560008301527f73206e6f742065786973740000000000000000000000000000000000000000006020830152604082019050919050565b60006128df60158361301b565b91507f5769746864726177616c2069732070656e64696e6700000000000000000000006000830152602082019050919050565b600061291f60168361301b565b91507f56616c696461746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b600061295f60358361301b565b91507f4465706f7369746f7220616464726573732063616e6e6f742062652073616d6560008301527f2061732056616c696461746f72206164647265737300000000000000000000006020830152604082019050919050565b60006129c560118361301b565b91507f496e76616c69642076616c696461746f720000000000000000000000000000006000830152602082019050919050565b6000612a0560168361301b565b91507f4465706f7369746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b6000612a45600083613010565b9150600082019050919050565b6000612a5f602b8361301b565b91507f4465706f73697420616d6f756e742062656c6f77206d696e696d756d2064657060008301527f6f73697420616d6f756e740000000000000000000000000000000000000000006020830152604082019050919050565b6000612ac560198361301b565b91507f4465706f7369746f7220616c72656164792065786973746564000000000000006000830152602082019050919050565b6000612b0560198361301b565b91507f56616c696461746f7220616c72656164792065786973746564000000000000006000830152602082019050919050565b6000612b4560188361301b565b91507f4465706f7369746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b6000612b8560248361301b565b91507f4465706f7369746f72207769746864726177616c20726571756573742070656e60008301527f64696e67000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000612beb60238361301b565b91507f4465706f7369746f72207769746864726177616c20726571756573742065786960008301527f73747300000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000612c5160198361301b565b91507f4465706f7369746f722062616c616e6365206973207a65726f000000000000006000830152602082019050919050565b6000612c9160188361301b565b91507f56616c696461746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b612ccd81613066565b82525050565b6000612cde82612a38565b9150819050919050565b6000602082019050612cfd60008301846126bf565b92915050565b60006020820190508181036000830152612d1d81846126ce565b905092915050565b60006020820190508181036000830152612d3e8161272c565b9050919050565b60006020820190508181036000830152612d5e8161276c565b9050919050565b60006020820190508181036000830152612d7e816127ac565b9050919050565b60006020820190508181036000830152612d9e816127ec565b9050919050565b60006020820190508181036000830152612dbe8161282c565b9050919050565b60006020820190508181036000830152612dde8161286c565b9050919050565b60006020820190508181036000830152612dfe816128d2565b9050919050565b60006020820190508181036000830152612e1e81612912565b9050919050565b60006020820190508181036000830152612e3e81612952565b9050919050565b60006020820190508181036000830152612e5e816129b8565b9050919050565b60006020820190508181036000830152612e7e816129f8565b9050919050565b60006020820190508181036000830152612e9e81612a52565b9050919050565b60006020820190508181036000830152612ebe81612ab8565b9050919050565b60006020820190508181036000830152612ede81612af8565b9050919050565b60006020820190508181036000830152612efe81612b38565b9050919050565b60006020820190508181036000830152612f1e81612b78565b9050919050565b60006020820190508181036000830152612f3e81612bde565b9050919050565b60006020820190508181036000830152612f5e81612c44565b9050919050565b60006020820190508181036000830152612f7e81612c84565b9050919050565b6000602082019050612f9a6000830184612cc4565b92915050565b6000606082019050612fb56000830186612cc4565b612fc26020830185612cc4565b612fcf6040830184612cc4565b949350505050565b6000819050602082019050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b60006130378261303e565b9050919050565b60007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b6130798161302c565b811461308457600080fd5b50565b61309081613066565b811461309b57600080fd5b5056fea26469706673582212202be7e9787cc35d13d13ef00139ede65d397360fed3e62a8414aa77df8ab85e6664736f6c637827302e372e362d646576656c6f702e323032332e31322e312b636f6d6d69742e30356365636362660058",
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
