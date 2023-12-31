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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"netBalance\",\"type\":\"uint256\"}],\"name\":\"OnCompleteWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"OnInitiateWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnNewDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"OnPauseValidation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"OnResumeValidation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewardAmount\",\"type\":\"uint256\"}],\"name\":\"OnReward\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"slashedAmount\",\"type\":\"uint256\"}],\"name\":\"OnSlashing\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"rewardAmount\",\"type\":\"uint256\"}],\"name\":\"addDepositorReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"slashAmount\",\"type\":\"uint256\"}],\"name\":\"addDepositorSlashing\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeWithdrawal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"didDepositorEverExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"didValidatorEverExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"doesDepositorExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"doesValidatorExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getBalanceOfDepositor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDepositorCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"getDepositorOfValidator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getDepositorRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getDepositorSlashings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getNetBalanceOfDepositor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalDepositedBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getValidatorOfDepositor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositorAddress\",\"type\":\"address\"}],\"name\":\"getWithdrawalBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initiateWithdrawal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"isValidationPaused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"listValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"newDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseValidation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resumeValidation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526000600255600060035534801561001a57600080fd5b5061231f8061002a6000396000f3fe6080604052600436106101355760003560e01c806377c06fdc116100ab578063c200baf91161006f578063c200baf91461049f578063c97ab777146104dc578063dd77e5cc146104f3578063e03ff7cb14610530578063f17bb4621461055b578063ff9205ab1461058657610135565b806377c06fdc146103805780637942317c146103bd578063a7113fee146103fa578063b112861214610437578063b51d1d4f1461047457610135565b806351cb11ab116100fd57806351cb11ab1461026b57806368d4e544146102a85780636d1e33cd146102d35780636d727bd0146103105780636e7f5bd31461034d578063731f750d1461036457610135565b80632ca3c0411461013a57806338c70a60146101775780633e3bc1a7146101b45780634f4af09e146101f157806351ca53171461022e575b600080fd5b34801561014657600080fd5b50610161600480360381019061015c919061176f565b6105b1565b60405161016e91906121c6565b60405180910390f35b34801561018357600080fd5b5061019e6004803603810190610199919061176f565b6105ce565b6040516101ab9190611f2b565b60405180910390f35b3480156101c057600080fd5b506101db60048036038101906101d6919061176f565b6105f8565b6040516101e89190611f2b565b60405180910390f35b3480156101fd57600080fd5b506102186004803603810190610213919061176f565b610622565b60405161022591906121c6565b60405180910390f35b34801561023a57600080fd5b5061025560048036038101906102509190611798565b61070e565b60405161026291906121c6565b60405180910390f35b34801561027757600080fd5b50610292600480360381019061028d919061176f565b61087c565b60405161029f91906121c6565b60405180910390f35b3480156102b457600080fd5b506102bd610899565b6040516102ca9190611f09565b60405180910390f35b3480156102df57600080fd5b506102fa60048036038101906102f5919061176f565b6108f1565b6040516103079190611f2b565b60405180910390f35b34801561031c57600080fd5b506103376004803603810190610332919061176f565b61091b565b6040516103449190611e9c565b60405180910390f35b34801561035957600080fd5b5061036261093d565b005b61037e6004803603810190610379919061176f565b610a91565b005b34801561038c57600080fd5b506103a760048036038101906103a2919061176f565b610ef0565b6040516103b491906121c6565b60405180910390f35b3480156103c957600080fd5b506103e460048036038101906103df9190611798565b610f0d565b6040516103f191906121c6565b60405180910390f35b34801561040657600080fd5b50610421600480360381019061041c919061176f565b610fe2565b60405161042e9190611e9c565b60405180910390f35b34801561044357600080fd5b5061045e6004803603810190610459919061176f565b611004565b60405161046b9190611f2b565b60405180910390f35b34801561048057600080fd5b5061048961102e565b60405161049691906121c6565b60405180910390f35b3480156104ab57600080fd5b506104c660048036038101906104c1919061176f565b611289565b6040516104d391906121c6565b60405180910390f35b3480156104e857600080fd5b506104f16112a6565b005b3480156104ff57600080fd5b5061051a6004803603810190610515919061176f565b6113fa565b6040516105279190611f2b565b60405180910390f35b34801561053c57600080fd5b50610545611424565b60405161055291906121c6565b60405180910390f35b34801561056757600080fd5b506105706116e9565b60405161057d91906121c6565b60405180910390f35b34801561059257600080fd5b5061059b6116f3565b6040516105a891906121c6565b60405180910390f35b6000600b6000838152602001908152602001600020549050919050565b60006007600083815260200190815260200160002060009054906101000a900460ff169050919050565b6000600d600083815260200190815260200160002060009054906101000a900460ff169050919050565b60008015156005600084815260200190815260200160002060009054906101000a900460ff16151514156106595760009050610709565b6000600c600084815260200190815260200160002054111561067e5760009050610709565b60006106b9600b60008581526020019081526020016000205460016000868152602001908152602001600020546116fd90919063ffffffff16565b9050600a60008481526020019081526020016000205481116106df576000915050610709565b610705600a6000858152602001908152602001600020548261171990919063ffffffff16565b9150505b919050565b6000803314610752576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161074990611fa6565b60405180910390fd5b61077882600a6000868152602001908152602001600020546116fd90919063ffffffff16565b600a60008581526020019081526020016000208190555060008081846040516107a090611e87565b60006040518083038185875af1925050503d80600081146107dd576040519150601f19603f3d011682016040523d82523d6000602084013e6107e2565b606091505b5050905080610826576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161081d90612186565b60405180910390fd5b847fcadc6c149d7c30ba433e0a526c9f018a1c4dc5b32099790e4dd9d78a930218108560405161085691906121c6565b60405180910390a2600a6000868152602001908152602001600020549250505092915050565b6000600c6000838152602001908152602001600020549050919050565b606060008054806020026020016040519081016040528092919081815260200182805480156108e757602002820191906000526020600020905b8154815260200190600101908083116108d3575b5050505050905090565b60006006600083815260200190815260200160002060009054906101000a900460ff169050919050565b6000806008600084815260200190815260200160002054905080915050919050565b6000339050600115156005600083815260200190815260200160002060009054906101000a900460ff161515146109a9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109a090611f86565b60405180910390fd5b60006008600083815260200190815260200160002054905060011515600d600083815260200190815260200160002060009054906101000a900460ff16151514610a28576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a1f90612126565b60405180910390fd5b6000600d600083815260200190815260200160002060006101000a81548160ff0219169083151502179055507f94b4f7e95a6cbc76e8ee615e2921a9f77aba94b68b52301f22a5844b433daadb8282604051610a85929190611eb7565b60405180910390a15050565b600033905060003490506a0422ca8b0a00a425000000811015610ae9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ae0906120c6565b60405180910390fd5b82821415610b2c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b2390612046565b60405180910390fd5b6000831415610b70576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b6790612066565b60405180910390fd5b600015156004600085815260200190815260200160002060009054906101000a900460ff16151514610bd7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bce906121a6565b60405180910390fd5b600015156006600085815260200190815260200160002060009054906101000a900460ff16151514610c3e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c3590612026565b60405180910390fd5b60008331905060008114610c87576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c7e90611f46565b60405180910390fd5b600015156005600085815260200190815260200160002060009054906101000a900460ff16151514610cee576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ce5906120e6565b60405180910390fd5b600015156007600085815260200190815260200160002060009054906101000a900460ff16151514610d55576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d4c906120a6565b60405180910390fd5b6000849080600181540180825580915050600190039060005260206000200160009091909190915055610d93826002546116fd90919063ffffffff16565b600281905550610daf60016003546116fd90919063ffffffff16565b60038190555081600160008581526020019081526020016000208190555060016004600086815260200190815260200160002060006101000a81548160ff02191690831515021790555060016005600085815260200190815260200160002060006101000a81548160ff02191690831515021790555060016006600086815260200190815260200160002060006101000a81548160ff02191690831515021790555060016007600085815260200190815260200160002060006101000a81548160ff02191690831515021790555082600860008681526020019081526020016000208190555083600960008581526020019081526020016000208190555083837fbe02029a5af0c964ebee7370f030cf18a026aae3a5d66f8107aee23f226d9ada844342604051610ee2939291906121e1565b60405180910390a350505050565b600060016000838152602001908152602001600020549050919050565b6000803314610f51576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f4890611fa6565b60405180910390fd5b610f7782600b6000868152602001908152602001600020546116fd90919063ffffffff16565b600b600085815260200190815260200160002081905550827fd1072bb52c3131d0c96197b73fb8a45637e30f8b6664fc142310cc9b242859b483604051610fbe91906121c6565b60405180910390a2600b600084815260200190815260200160002054905092915050565b6000806009600084815260200190815260200160002054905080915050919050565b60006004600083815260200190815260200160002060009054906101000a900460ff169050919050565b600080339050600115156005600083815260200190815260200160002060009054906101000a900460ff1615151461109b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161109290611f86565b60405180910390fd5b6000600c600083815260200190815260200160002054146110f1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110e890612146565b60405180910390fd5b6000600160008381526020019081526020016000205411611147576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161113e90612166565b60405180910390fd5b600030634f4af09e836040518263ffffffff1660e01b815260040161116c9190611e9c565b60206040518083038186803b15801561118457600080fd5b505afa158015611198573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111bc91906117d4565b905060008111611201576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111f890611fe6565b60405180910390fd5b6203e8004301600c60008481526020019081526020016000208190555060006005600084815260200190815260200160002060006101000a81548160ff0219169083151502179055507f47e2f9085249de3b62accafda3451074e283e2c6f30a39ae0b9952f3a0f8ecf7826040516112799190611e9c565b60405180910390a1809250505090565b6000600a6000838152602001908152602001600020549050919050565b6000339050600115156005600083815260200190815260200160002060009054906101000a900460ff16151514611312576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161130990611f86565b60405180910390fd5b60006008600083815260200190815260200160002054905060001515600d600083815260200190815260200160002060009054906101000a900460ff16151514611391576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161138890612086565b60405180910390fd5b6001600d600083815260200190815260200160002060006101000a81548160ff0219169083151502179055507f805619aaa3e6ed885faf94910ba4a75ca1ca3b3e0c539f6dd93b4d1e48b6782682826040516113ee929190611eb7565b60405180910390a15050565b60006005600083815260200190815260200160002060009054906101000a900460ff169050919050565b6000803390506000600c60008381526020019081526020016000205411611480576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161147790612006565b60405180910390fd5b600c60008281526020019081526020016000205443116114d5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016114cc90612106565b60405180910390fd5b6000611510600b60008481526020019081526020016000205460016000858152602001908152602001600020546116fd90919063ffffffff16565b9050600a6000838152602001908152602001600020548111611567576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161155e90611f66565b60405180910390fd5b600061158f600a6000858152602001908152602001600020548361171990919063ffffffff16565b90506001600084815260200190815260200160002060009055600b600084815260200190815260200160002060009055600a6000848152602001908152602001600020600090556005600084815260200190815260200160002060006101000a81549060ff0219169055600c6000848152602001908152602001600020600090556000838260405161162090611e87565b60006040518083038185875af1925050503d806000811461165d576040519150601f19603f3d011682016040523d82523d6000602084013e611662565b606091505b50509050806116a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161169d90611fc6565b60405180910390fd5b7f3b7625df2fe98121022a93d8d1c02ff11f383651160f046ffc4c1e379d4242ef84836040516116d7929190611ee0565b60405180910390a18194505050505090565b6000600354905090565b6000600254905090565b60008082840190508381101561170f57fe5b8091505092915050565b60008282111561172557fe5b818303905092915050565b60008135905061173f81612295565b92915050565b600081359050611754816122ac565b92915050565b600081519050611769816122ac565b92915050565b60006020828403121561178157600080fd5b600061178f84828501611730565b91505092915050565b600080604083850312156117ab57600080fd5b60006117b985828601611730565b92505060206117ca85828601611745565b9150509250929050565b6000602082840312156117e657600080fd5b60006117f48482850161175a565b91505092915050565b60006118098383611815565b60208301905092915050565b61181e8161226d565b82525050565b61182d8161226d565b82525050565b600061183e82612228565b6118488185612240565b935061185383612218565b8060005b8381101561188457815161186b88826117fd565b975061187683612233565b925050600181019050611857565b5085935050505092915050565b61189a8161227f565b82525050565b60006118ad60208361225c565b91507f76616c696461746f722062616c616e63652073686f756c64206265207a65726f6000830152602082019050919050565b60006118ed60138361225c565b91507f62616c616e6365206973206e65676174697665000000000000000000000000006000830152602082019050919050565b600061192d60188361225c565b91507f4465706f7369746f7220646f6573206e6f7420657869737400000000000000006000830152602082019050919050565b600061196d60198361225c565b91507f4f6e6c7920564d2063616c6c732061726520616c6c6f776564000000000000006000830152602082019050919050565b60006119ad600f8361225c565b91507f5769746864726177206661696c656400000000000000000000000000000000006000830152602082019050919050565b60006119ed601d8361225c565b91507f4465706f7369746f72206e65742062616c616e6365206973207a65726f0000006000830152602082019050919050565b6000611a2d602b8361225c565b91507f4465706f7369746f72207769746864726177616c207265717565737420646f6560008301527f73206e6f742065786973740000000000000000000000000000000000000000006020830152604082019050919050565b6000611a9360168361225c565b91507f56616c696461746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b6000611ad360358361225c565b91507f4465706f7369746f7220616464726573732063616e6e6f742062652073616d6560008301527f2061732056616c696461746f72206164647265737300000000000000000000006020830152604082019050919050565b6000611b3960118361225c565b91507f496e76616c69642076616c696461746f720000000000000000000000000000006000830152602082019050919050565b6000611b79601c8361225c565b91507f56616c69646174696f6e20697320616c726561647920706175736564000000006000830152602082019050919050565b6000611bb960168361225c565b91507f4465706f7369746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b6000611bf9600083612251565b9150600082019050919050565b6000611c13602b8361225c565b91507f4465706f73697420616d6f756e742062656c6f77206d696e696d756d2064657060008301527f6f73697420616d6f756e740000000000000000000000000000000000000000006020830152604082019050919050565b6000611c7960188361225c565b91507f4465706f7369746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b6000611cb960248361225c565b91507f4465706f7369746f72207769746864726177616c20726571756573742070656e60008301527f64696e67000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000611d1f60188361225c565b91507f56616c69646174696f6e206973206e6f742070617573656400000000000000006000830152602082019050919050565b6000611d5f60238361225c565b91507f4465706f7369746f72207769746864726177616c20726571756573742065786960008301527f73747300000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000611dc560198361225c565b91507f4465706f7369746f722062616c616e6365206973207a65726f000000000000006000830152602082019050919050565b6000611e05601e8361225c565b91507f7472616e7366657220746f207a65726f41646472657373206661696c656400006000830152602082019050919050565b6000611e4560188361225c565b91507f56616c696461746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b611e818161228b565b82525050565b6000611e9282611bec565b9150819050919050565b6000602082019050611eb16000830184611824565b92915050565b6000604082019050611ecc6000830185611824565b611ed96020830184611824565b9392505050565b6000604082019050611ef56000830185611824565b611f026020830184611e78565b9392505050565b60006020820190508181036000830152611f238184611833565b905092915050565b6000602082019050611f406000830184611891565b92915050565b60006020820190508181036000830152611f5f816118a0565b9050919050565b60006020820190508181036000830152611f7f816118e0565b9050919050565b60006020820190508181036000830152611f9f81611920565b9050919050565b60006020820190508181036000830152611fbf81611960565b9050919050565b60006020820190508181036000830152611fdf816119a0565b9050919050565b60006020820190508181036000830152611fff816119e0565b9050919050565b6000602082019050818103600083015261201f81611a20565b9050919050565b6000602082019050818103600083015261203f81611a86565b9050919050565b6000602082019050818103600083015261205f81611ac6565b9050919050565b6000602082019050818103600083015261207f81611b2c565b9050919050565b6000602082019050818103600083015261209f81611b6c565b9050919050565b600060208201905081810360008301526120bf81611bac565b9050919050565b600060208201905081810360008301526120df81611c06565b9050919050565b600060208201905081810360008301526120ff81611c6c565b9050919050565b6000602082019050818103600083015261211f81611cac565b9050919050565b6000602082019050818103600083015261213f81611d12565b9050919050565b6000602082019050818103600083015261215f81611d52565b9050919050565b6000602082019050818103600083015261217f81611db8565b9050919050565b6000602082019050818103600083015261219f81611df8565b9050919050565b600060208201905081810360008301526121bf81611e38565b9050919050565b60006020820190506121db6000830184611e78565b92915050565b60006060820190506121f66000830186611e78565b6122036020830185611e78565b6122106040830184611e78565b949350505050565b6000819050602082019050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b60006122788261228b565b9050919050565b60008115159050919050565b6000819050919050565b61229e8161226d565b81146122a957600080fd5b50565b6122b58161228b565b81146122c057600080fd5b5056fea2646970667358221220d4c990bd487e82780b32d6f5884acf01c9c4374ead29a1ab12cdae32acf0f7ae64736f6c637828302e372e362d646576656c6f702e323032332e31322e33302b636f6d6d69742e37326538396665320059",
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

// DidDepositorEverExist is a free data retrieval call binding the contract method 0x38c70a60.
//
// Solidity: function didDepositorEverExist(address depositorAddress) view returns(bool)
func (_Staking *StakingCaller) DidDepositorEverExist(opts *bind.CallOpts, depositorAddress common.Address) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "didDepositorEverExist", depositorAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DidDepositorEverExist is a free data retrieval call binding the contract method 0x38c70a60.
//
// Solidity: function didDepositorEverExist(address depositorAddress) view returns(bool)
func (_Staking *StakingSession) DidDepositorEverExist(depositorAddress common.Address) (bool, error) {
	return _Staking.Contract.DidDepositorEverExist(&_Staking.CallOpts, depositorAddress)
}

// DidDepositorEverExist is a free data retrieval call binding the contract method 0x38c70a60.
//
// Solidity: function didDepositorEverExist(address depositorAddress) view returns(bool)
func (_Staking *StakingCallerSession) DidDepositorEverExist(depositorAddress common.Address) (bool, error) {
	return _Staking.Contract.DidDepositorEverExist(&_Staking.CallOpts, depositorAddress)
}

// DidValidatorEverExist is a free data retrieval call binding the contract method 0x6d1e33cd.
//
// Solidity: function didValidatorEverExist(address validatorAddress) view returns(bool)
func (_Staking *StakingCaller) DidValidatorEverExist(opts *bind.CallOpts, validatorAddress common.Address) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "didValidatorEverExist", validatorAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DidValidatorEverExist is a free data retrieval call binding the contract method 0x6d1e33cd.
//
// Solidity: function didValidatorEverExist(address validatorAddress) view returns(bool)
func (_Staking *StakingSession) DidValidatorEverExist(validatorAddress common.Address) (bool, error) {
	return _Staking.Contract.DidValidatorEverExist(&_Staking.CallOpts, validatorAddress)
}

// DidValidatorEverExist is a free data retrieval call binding the contract method 0x6d1e33cd.
//
// Solidity: function didValidatorEverExist(address validatorAddress) view returns(bool)
func (_Staking *StakingCallerSession) DidValidatorEverExist(validatorAddress common.Address) (bool, error) {
	return _Staking.Contract.DidValidatorEverExist(&_Staking.CallOpts, validatorAddress)
}

// DoesDepositorExist is a free data retrieval call binding the contract method 0xdd77e5cc.
//
// Solidity: function doesDepositorExist(address depositorAddress) view returns(bool)
func (_Staking *StakingCaller) DoesDepositorExist(opts *bind.CallOpts, depositorAddress common.Address) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "doesDepositorExist", depositorAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DoesDepositorExist is a free data retrieval call binding the contract method 0xdd77e5cc.
//
// Solidity: function doesDepositorExist(address depositorAddress) view returns(bool)
func (_Staking *StakingSession) DoesDepositorExist(depositorAddress common.Address) (bool, error) {
	return _Staking.Contract.DoesDepositorExist(&_Staking.CallOpts, depositorAddress)
}

// DoesDepositorExist is a free data retrieval call binding the contract method 0xdd77e5cc.
//
// Solidity: function doesDepositorExist(address depositorAddress) view returns(bool)
func (_Staking *StakingCallerSession) DoesDepositorExist(depositorAddress common.Address) (bool, error) {
	return _Staking.Contract.DoesDepositorExist(&_Staking.CallOpts, depositorAddress)
}

// DoesValidatorExist is a free data retrieval call binding the contract method 0xb1128612.
//
// Solidity: function doesValidatorExist(address validatorAddress) view returns(bool)
func (_Staking *StakingCaller) DoesValidatorExist(opts *bind.CallOpts, validatorAddress common.Address) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "doesValidatorExist", validatorAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DoesValidatorExist is a free data retrieval call binding the contract method 0xb1128612.
//
// Solidity: function doesValidatorExist(address validatorAddress) view returns(bool)
func (_Staking *StakingSession) DoesValidatorExist(validatorAddress common.Address) (bool, error) {
	return _Staking.Contract.DoesValidatorExist(&_Staking.CallOpts, validatorAddress)
}

// DoesValidatorExist is a free data retrieval call binding the contract method 0xb1128612.
//
// Solidity: function doesValidatorExist(address validatorAddress) view returns(bool)
func (_Staking *StakingCallerSession) DoesValidatorExist(validatorAddress common.Address) (bool, error) {
	return _Staking.Contract.DoesValidatorExist(&_Staking.CallOpts, validatorAddress)
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

// GetWithdrawalBlock is a free data retrieval call binding the contract method 0x51cb11ab.
//
// Solidity: function getWithdrawalBlock(address depositorAddress) view returns(uint256)
func (_Staking *StakingCaller) GetWithdrawalBlock(opts *bind.CallOpts, depositorAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getWithdrawalBlock", depositorAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetWithdrawalBlock is a free data retrieval call binding the contract method 0x51cb11ab.
//
// Solidity: function getWithdrawalBlock(address depositorAddress) view returns(uint256)
func (_Staking *StakingSession) GetWithdrawalBlock(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetWithdrawalBlock(&_Staking.CallOpts, depositorAddress)
}

// GetWithdrawalBlock is a free data retrieval call binding the contract method 0x51cb11ab.
//
// Solidity: function getWithdrawalBlock(address depositorAddress) view returns(uint256)
func (_Staking *StakingCallerSession) GetWithdrawalBlock(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetWithdrawalBlock(&_Staking.CallOpts, depositorAddress)
}

// IsValidationPaused is a free data retrieval call binding the contract method 0x3e3bc1a7.
//
// Solidity: function isValidationPaused(address validatorAddress) view returns(bool)
func (_Staking *StakingCaller) IsValidationPaused(opts *bind.CallOpts, validatorAddress common.Address) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "isValidationPaused", validatorAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidationPaused is a free data retrieval call binding the contract method 0x3e3bc1a7.
//
// Solidity: function isValidationPaused(address validatorAddress) view returns(bool)
func (_Staking *StakingSession) IsValidationPaused(validatorAddress common.Address) (bool, error) {
	return _Staking.Contract.IsValidationPaused(&_Staking.CallOpts, validatorAddress)
}

// IsValidationPaused is a free data retrieval call binding the contract method 0x3e3bc1a7.
//
// Solidity: function isValidationPaused(address validatorAddress) view returns(bool)
func (_Staking *StakingCallerSession) IsValidationPaused(validatorAddress common.Address) (bool, error) {
	return _Staking.Contract.IsValidationPaused(&_Staking.CallOpts, validatorAddress)
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

// CompleteWithdrawal is a paid mutator transaction binding the contract method 0xe03ff7cb.
//
// Solidity: function completeWithdrawal() returns(uint256)
func (_Staking *StakingTransactor) CompleteWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "completeWithdrawal")
}

// CompleteWithdrawal is a paid mutator transaction binding the contract method 0xe03ff7cb.
//
// Solidity: function completeWithdrawal() returns(uint256)
func (_Staking *StakingSession) CompleteWithdrawal() (*types.Transaction, error) {
	return _Staking.Contract.CompleteWithdrawal(&_Staking.TransactOpts)
}

// CompleteWithdrawal is a paid mutator transaction binding the contract method 0xe03ff7cb.
//
// Solidity: function completeWithdrawal() returns(uint256)
func (_Staking *StakingTransactorSession) CompleteWithdrawal() (*types.Transaction, error) {
	return _Staking.Contract.CompleteWithdrawal(&_Staking.TransactOpts)
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0xb51d1d4f.
//
// Solidity: function initiateWithdrawal() returns(uint256)
func (_Staking *StakingTransactor) InitiateWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "initiateWithdrawal")
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0xb51d1d4f.
//
// Solidity: function initiateWithdrawal() returns(uint256)
func (_Staking *StakingSession) InitiateWithdrawal() (*types.Transaction, error) {
	return _Staking.Contract.InitiateWithdrawal(&_Staking.TransactOpts)
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0xb51d1d4f.
//
// Solidity: function initiateWithdrawal() returns(uint256)
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

// PauseValidation is a paid mutator transaction binding the contract method 0xc97ab777.
//
// Solidity: function pauseValidation() returns()
func (_Staking *StakingTransactor) PauseValidation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "pauseValidation")
}

// PauseValidation is a paid mutator transaction binding the contract method 0xc97ab777.
//
// Solidity: function pauseValidation() returns()
func (_Staking *StakingSession) PauseValidation() (*types.Transaction, error) {
	return _Staking.Contract.PauseValidation(&_Staking.TransactOpts)
}

// PauseValidation is a paid mutator transaction binding the contract method 0xc97ab777.
//
// Solidity: function pauseValidation() returns()
func (_Staking *StakingTransactorSession) PauseValidation() (*types.Transaction, error) {
	return _Staking.Contract.PauseValidation(&_Staking.TransactOpts)
}

// ResumeValidation is a paid mutator transaction binding the contract method 0x6e7f5bd3.
//
// Solidity: function resumeValidation() returns()
func (_Staking *StakingTransactor) ResumeValidation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "resumeValidation")
}

// ResumeValidation is a paid mutator transaction binding the contract method 0x6e7f5bd3.
//
// Solidity: function resumeValidation() returns()
func (_Staking *StakingSession) ResumeValidation() (*types.Transaction, error) {
	return _Staking.Contract.ResumeValidation(&_Staking.TransactOpts)
}

// ResumeValidation is a paid mutator transaction binding the contract method 0x6e7f5bd3.
//
// Solidity: function resumeValidation() returns()
func (_Staking *StakingTransactorSession) ResumeValidation() (*types.Transaction, error) {
	return _Staking.Contract.ResumeValidation(&_Staking.TransactOpts)
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
	NetBalance       *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnCompleteWithdrawal is a free log retrieval operation binding the contract event 0x3b7625df2fe98121022a93d8d1c02ff11f383651160f046ffc4c1e379d4242ef.
//
// Solidity: event OnCompleteWithdrawal(address depositorAddress, uint256 netBalance)
func (_Staking *StakingFilterer) FilterOnCompleteWithdrawal(opts *bind.FilterOpts) (*StakingOnCompleteWithdrawalIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnCompleteWithdrawal")
	if err != nil {
		return nil, err
	}
	return &StakingOnCompleteWithdrawalIterator{contract: _Staking.contract, event: "OnCompleteWithdrawal", logs: logs, sub: sub}, nil
}

// WatchOnCompleteWithdrawal is a free log subscription operation binding the contract event 0x3b7625df2fe98121022a93d8d1c02ff11f383651160f046ffc4c1e379d4242ef.
//
// Solidity: event OnCompleteWithdrawal(address depositorAddress, uint256 netBalance)
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

// ParseOnCompleteWithdrawal is a log parse operation binding the contract event 0x3b7625df2fe98121022a93d8d1c02ff11f383651160f046ffc4c1e379d4242ef.
//
// Solidity: event OnCompleteWithdrawal(address depositorAddress, uint256 netBalance)
func (_Staking *StakingFilterer) ParseOnCompleteWithdrawal(log types.Log) (*StakingOnCompleteWithdrawal, error) {
	event := new(StakingOnCompleteWithdrawal)
	if err := _Staking.contract.UnpackLog(event, "OnCompleteWithdrawal", log); err != nil {
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

// StakingOnPauseValidationIterator is returned from FilterOnPauseValidation and is used to iterate over the raw logs and unpacked data for OnPauseValidation events raised by the Staking contract.
type StakingOnPauseValidationIterator struct {
	Event *StakingOnPauseValidation // Event containing the contract specifics and raw log

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
func (it *StakingOnPauseValidationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOnPauseValidation)
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
		it.Event = new(StakingOnPauseValidation)
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
func (it *StakingOnPauseValidationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOnPauseValidationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOnPauseValidation represents a OnPauseValidation event raised by the Staking contract.
type StakingOnPauseValidation struct {
	DepositorAddress common.Address
	ValidatorAddress common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnPauseValidation is a free log retrieval operation binding the contract event 0x805619aaa3e6ed885faf94910ba4a75ca1ca3b3e0c539f6dd93b4d1e48b67826.
//
// Solidity: event OnPauseValidation(address depositorAddress, address validatorAddress)
func (_Staking *StakingFilterer) FilterOnPauseValidation(opts *bind.FilterOpts) (*StakingOnPauseValidationIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnPauseValidation")
	if err != nil {
		return nil, err
	}
	return &StakingOnPauseValidationIterator{contract: _Staking.contract, event: "OnPauseValidation", logs: logs, sub: sub}, nil
}

// WatchOnPauseValidation is a free log subscription operation binding the contract event 0x805619aaa3e6ed885faf94910ba4a75ca1ca3b3e0c539f6dd93b4d1e48b67826.
//
// Solidity: event OnPauseValidation(address depositorAddress, address validatorAddress)
func (_Staking *StakingFilterer) WatchOnPauseValidation(opts *bind.WatchOpts, sink chan<- *StakingOnPauseValidation) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OnPauseValidation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOnPauseValidation)
				if err := _Staking.contract.UnpackLog(event, "OnPauseValidation", log); err != nil {
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

// ParseOnPauseValidation is a log parse operation binding the contract event 0x805619aaa3e6ed885faf94910ba4a75ca1ca3b3e0c539f6dd93b4d1e48b67826.
//
// Solidity: event OnPauseValidation(address depositorAddress, address validatorAddress)
func (_Staking *StakingFilterer) ParseOnPauseValidation(log types.Log) (*StakingOnPauseValidation, error) {
	event := new(StakingOnPauseValidation)
	if err := _Staking.contract.UnpackLog(event, "OnPauseValidation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingOnResumeValidationIterator is returned from FilterOnResumeValidation and is used to iterate over the raw logs and unpacked data for OnResumeValidation events raised by the Staking contract.
type StakingOnResumeValidationIterator struct {
	Event *StakingOnResumeValidation // Event containing the contract specifics and raw log

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
func (it *StakingOnResumeValidationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOnResumeValidation)
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
		it.Event = new(StakingOnResumeValidation)
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
func (it *StakingOnResumeValidationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOnResumeValidationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOnResumeValidation represents a OnResumeValidation event raised by the Staking contract.
type StakingOnResumeValidation struct {
	DepositorAddress common.Address
	ValidatorAddress common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnResumeValidation is a free log retrieval operation binding the contract event 0x94b4f7e95a6cbc76e8ee615e2921a9f77aba94b68b52301f22a5844b433daadb.
//
// Solidity: event OnResumeValidation(address depositorAddress, address validatorAddress)
func (_Staking *StakingFilterer) FilterOnResumeValidation(opts *bind.FilterOpts) (*StakingOnResumeValidationIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnResumeValidation")
	if err != nil {
		return nil, err
	}
	return &StakingOnResumeValidationIterator{contract: _Staking.contract, event: "OnResumeValidation", logs: logs, sub: sub}, nil
}

// WatchOnResumeValidation is a free log subscription operation binding the contract event 0x94b4f7e95a6cbc76e8ee615e2921a9f77aba94b68b52301f22a5844b433daadb.
//
// Solidity: event OnResumeValidation(address depositorAddress, address validatorAddress)
func (_Staking *StakingFilterer) WatchOnResumeValidation(opts *bind.WatchOpts, sink chan<- *StakingOnResumeValidation) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OnResumeValidation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOnResumeValidation)
				if err := _Staking.contract.UnpackLog(event, "OnResumeValidation", log); err != nil {
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

// ParseOnResumeValidation is a log parse operation binding the contract event 0x94b4f7e95a6cbc76e8ee615e2921a9f77aba94b68b52301f22a5844b433daadb.
//
// Solidity: event OnResumeValidation(address depositorAddress, address validatorAddress)
func (_Staking *StakingFilterer) ParseOnResumeValidation(log types.Log) (*StakingOnResumeValidation, error) {
	event := new(StakingOnResumeValidation)
	if err := _Staking.contract.UnpackLog(event, "OnResumeValidation", log); err != nil {
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
