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
	Bin: "0x60806040526000600255600060035534801561001a57600080fd5b50611f528061002a6000396000f3fe6080604052600436106100e85760003560e01c80637942317c1161008a578063e03ff7cb11610059578063e03ff7cb14610333578063f17bb4621461034a578063f6abfc7614610375578063ff9205ab14610391576100e8565b80637942317c14610265578063a7113fee146102a2578063b51d1d4f146102df578063c200baf9146102f6576100e8565b806368d4e544116100c657806368d4e544146101a45780636d727bd0146101cf578063731f750d1461020c57806377c06fdc14610228576100e8565b80632ca3c041146100ed5780634f4af09e1461012a57806351ca531714610167575b600080fd5b3480156100f957600080fd5b50610114600480360381019061010f9190611487565b6103bc565b6040516101219190611e02565b60405180910390f35b34801561013657600080fd5b50610151600480360381019061014c9190611487565b6103d9565b60405161015e9190611e02565b60405180910390f35b34801561017357600080fd5b5061018e600480360381019061018991906114b0565b6104c5565b60405161019b9190611e02565b60405180910390f35b3480156101b057600080fd5b506101b961059a565b6040516101c69190611b80565b60405180910390f35b3480156101db57600080fd5b506101f660048036038101906101f19190611487565b6105f2565b6040516102039190611b65565b60405180910390f35b61022660048036038101906102219190611487565b610614565b005b34801561023457600080fd5b5061024f600480360381019061024a9190611487565b610a73565b60405161025c9190611e02565b60405180910390f35b34801561027157600080fd5b5061028c600480360381019061028791906114b0565b610a90565b6040516102999190611e02565b60405180910390f35b3480156102ae57600080fd5b506102c960048036038101906102c49190611487565b610b65565b6040516102d69190611b65565b60405180910390f35b3480156102eb57600080fd5b506102f4610b87565b005b34801561030257600080fd5b5061031d60048036038101906103189190611487565b610d1a565b60405161032a9190611e02565b60405180910390f35b34801561033f57600080fd5b50610348610d37565b005b34801561035657600080fd5b5061035f610f98565b60405161036c9190611e02565b60405180910390f35b61038f600480360381019061038a9190611487565b610fa2565b005b34801561039d57600080fd5b506103a661140b565b6040516103b39190611e02565b60405180910390f35b6000600b6000838152602001908152602001600020549050919050565b60008015156005600084815260200190815260200160002060009054906101000a900460ff161515141561041057600090506104c0565b6000600c600084815260200190815260200160002054111561043557600090506104c0565b6000610470600b600085815260200190815260200160002054600160008681526020019081526020016000205461141590919063ffffffff16565b9050600a60008481526020019081526020016000205481116104965760009150506104c0565b6104bc600a6000858152602001908152602001600020548261143190919063ffffffff16565b9150505b919050565b6000803314610509576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161050090611be2565b60405180910390fd5b61052f82600a60008681526020019081526020016000205461141590919063ffffffff16565b600a600085815260200190815260200160002081905550827fcadc6c149d7c30ba433e0a526c9f018a1c4dc5b32099790e4dd9d78a93021810836040516105769190611e02565b60405180910390a2600a600084815260200190815260200160002054905092915050565b606060008054806020026020016040519081016040528092919081815260200182805480156105e857602002820191906000526020600020905b8154815260200190600101908083116105d4575b5050505050905090565b6000806008600084815260200190815260200160002054905080915050919050565b600033905060003490506acecb8f27f4200f3a00000081101561066c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161066390611d02565b60405180910390fd5b828214156106af576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106a690611ca2565b60405180910390fd5b60008314156106f3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106ea90611cc2565b60405180910390fd5b600015156004600085815260200190815260200160002060009054906101000a900460ff1615151461075a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161075190611de2565b60405180910390fd5b600015156006600085815260200190815260200160002060009054906101000a900460ff161515146107c1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107b890611c82565b60405180910390fd5b6000833190506000811461080a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161080190611ba2565b60405180910390fd5b600015156005600085815260200190815260200160002060009054906101000a900460ff16151514610871576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161086890611d62565b60405180910390fd5b600015156007600085815260200190815260200160002060009054906101000a900460ff161515146108d8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108cf90611ce2565b60405180910390fd5b60008490806001815401808255809150506001900390600052602060002001600090919091909150556109168260025461141590919063ffffffff16565b600281905550610932600160035461141590919063ffffffff16565b60038190555081600160008581526020019081526020016000208190555060016004600086815260200190815260200160002060006101000a81548160ff02191690831515021790555060016005600085815260200190815260200160002060006101000a81548160ff02191690831515021790555060016006600086815260200190815260200160002060006101000a81548160ff02191690831515021790555060016007600085815260200190815260200160002060006101000a81548160ff02191690831515021790555082600860008681526020019081526020016000208190555083600960008581526020019081526020016000208190555083837fbe02029a5af0c964ebee7370f030cf18a026aae3a5d66f8107aee23f226d9ada844342604051610a6593929190611e1d565b60405180910390a350505050565b600060016000838152602001908152602001600020549050919050565b6000803314610ad4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610acb90611be2565b60405180910390fd5b610afa82600b60008681526020019081526020016000205461141590919063ffffffff16565b600b600085815260200190815260200160002081905550827fd1072bb52c3131d0c96197b73fb8a45637e30f8b6664fc142310cc9b242859b483604051610b419190611e02565b60405180910390a2600b600084815260200190815260200160002054905092915050565b6000806008600084815260200190815260200160002054905080915050919050565b6000339050600115156005600083815260200190815260200160002060009054906101000a900460ff16151514610bf3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bea90611bc2565b60405180910390fd5b6000600c60008381526020019081526020016000205414610c49576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c4090611da2565b60405180910390fd5b6000600160008381526020019081526020016000205411610c9f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c9690611dc2565b60405180910390fd5b6203d0904301600c600083815260200190815260200160002081905550600015156005600083815260200190815260200160002060009054906101000a905050507f47e2f9085249de3b62accafda3451074e283e2c6f30a39ae0b9952f3a0f8ecf781604051610d0f9190611b65565b60405180910390a150565b6000600a6000838152602001908152602001600020549050919050565b60003390506000600c60008381526020019081526020016000205411610d92576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d8990611c42565b60405180910390fd5b600c6000828152602001908152602001600020544311610de7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610dde90611d82565b60405180910390fd5b600030634f4af09e836040518263ffffffff1660e01b8152600401610e0c9190611b65565b60206040518083038186803b158015610e2457600080fd5b505afa158015610e38573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e5c91906114ec565b90506001600083815260200190815260200160002060009055600b600083815260200190815260200160002060009055600a6000838152602001908152602001600020600090556005600083815260200190815260200160002060006101000a81549060ff021916905560008282604051610ed690611b50565b60006040518083038185875af1925050503d8060008114610f13576040519150601f19603f3d011682016040523d82523d6000602084013e610f18565b606091505b5050905080610f5c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f5390611c22565b60405180910390fd5b7fda6373ecbed97803ca40cc1b7ed282476253b5aa7cd093dbb61d6990d5efcde483604051610f8b9190611b65565b60405180910390a1505050565b6000600354905090565b600015156004600083815260200190815260200160002060009054906101000a900460ff16151514611009576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161100090611de2565b60405180910390fd5b600015156005600083815260200190815260200160002060009054906101000a900460ff16151514611070576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161106790611c02565b60405180910390fd5b600015156006600083815260200190815260200160002060009054906101000a900460ff161515146110d7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110ce90611d42565b60405180910390fd5b600015156007600083815260200190815260200160002060009054906101000a900460ff1615151461113e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161113590611d22565b60405180910390fd5b6000813114611182576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161117990611ba2565b60405180910390fd5b60008114156111c6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111bd90611cc2565b60405180910390fd5b60003390508181141561120e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161120590611ca2565b60405180910390fd5b600115156005600083815260200190815260200160002060009054906101000a900460ff16151514611275576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161126c90611bc2565b60405180910390fd5b6000600c600083815260200190815260200160002054146112cb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112c290611c62565b60405180910390fd5b60016004600084815260200190815260200160002060006101000a81548160ff02191690831515021790555060016006600084815260200190815260200160002060006101000a81548160ff021916908315150217905550806008600084815260200190815260200160002081905550816009600083815260200190815260200160002081905550600082908060018154018082558091505060019003906000526020600020016000909190919091505560006009600083815260200190815260200160002054905060006004600083815260200190815260200160002060006101000a81548160ff02191690831515021790555060086000828152602001908152602001600020600090558281837f66d7a4dea74851a2dcc039f4c17dd2862081083c29daceb1a9346783de9185ce60405160405180910390a4505050565b6000600254905090565b60008082840190508381101561142757fe5b8091505092915050565b60008282111561143d57fe5b818303905092915050565b60008135905061145781611ec5565b92915050565b60008135905061146c81611edc565b92915050565b60008151905061148181611edc565b92915050565b60006020828403121561149957600080fd5b60006114a784828501611448565b91505092915050565b600080604083850312156114c357600080fd5b60006114d185828601611448565b92505060206114e28582860161145d565b9150509250929050565b6000602082840312156114fe57600080fd5b600061150c84828501611472565b91505092915050565b6000611521838361152d565b60208301905092915050565b61153681611ea9565b82525050565b61154581611ea9565b82525050565b600061155682611e64565b6115608185611e7c565b935061156b83611e54565b8060005b8381101561159c5781516115838882611515565b975061158e83611e6f565b92505060018101905061156f565b5085935050505092915050565b60006115b6602083611e98565b91507f76616c696461746f722062616c616e63652073686f756c64206265207a65726f6000830152602082019050919050565b60006115f6601883611e98565b91507f4465706f7369746f7220646f6573206e6f7420657869737400000000000000006000830152602082019050919050565b6000611636601983611e98565b91507f4f6e6c7920564d2063616c6c732061726520616c6c6f776564000000000000006000830152602082019050919050565b6000611676601883611e98565b91507f56616c696461746f722069732061206465706f7369746f7200000000000000006000830152602082019050919050565b60006116b6600f83611e98565b91507f5769746864726177206661696c656400000000000000000000000000000000006000830152602082019050919050565b60006116f6602b83611e98565b91507f4465706f7369746f72207769746864726177616c207265717565737420646f6560008301527f73206e6f742065786973740000000000000000000000000000000000000000006020830152604082019050919050565b600061175c601583611e98565b91507f5769746864726177616c2069732070656e64696e6700000000000000000000006000830152602082019050919050565b600061179c601683611e98565b91507f56616c696461746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b60006117dc603583611e98565b91507f4465706f7369746f7220616464726573732063616e6e6f742062652073616d6560008301527f2061732056616c696461746f72206164647265737300000000000000000000006020830152604082019050919050565b6000611842601183611e98565b91507f496e76616c69642076616c696461746f720000000000000000000000000000006000830152602082019050919050565b6000611882601683611e98565b91507f4465706f7369746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b60006118c2600083611e8d565b9150600082019050919050565b60006118dc602b83611e98565b91507f4465706f73697420616d6f756e742062656c6f77206d696e696d756d2064657060008301527f6f73697420616d6f756e740000000000000000000000000000000000000000006020830152604082019050919050565b6000611942601983611e98565b91507f4465706f7369746f7220616c72656164792065786973746564000000000000006000830152602082019050919050565b6000611982601983611e98565b91507f56616c696461746f7220616c72656164792065786973746564000000000000006000830152602082019050919050565b60006119c2601883611e98565b91507f4465706f7369746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b6000611a02602483611e98565b91507f4465706f7369746f72207769746864726177616c20726571756573742070656e60008301527f64696e67000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000611a68602383611e98565b91507f4465706f7369746f72207769746864726177616c20726571756573742065786960008301527f73747300000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000611ace601983611e98565b91507f4465706f7369746f722062616c616e6365206973207a65726f000000000000006000830152602082019050919050565b6000611b0e601883611e98565b91507f56616c696461746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b611b4a81611ebb565b82525050565b6000611b5b826118b5565b9150819050919050565b6000602082019050611b7a600083018461153c565b92915050565b60006020820190508181036000830152611b9a818461154b565b905092915050565b60006020820190508181036000830152611bbb816115a9565b9050919050565b60006020820190508181036000830152611bdb816115e9565b9050919050565b60006020820190508181036000830152611bfb81611629565b9050919050565b60006020820190508181036000830152611c1b81611669565b9050919050565b60006020820190508181036000830152611c3b816116a9565b9050919050565b60006020820190508181036000830152611c5b816116e9565b9050919050565b60006020820190508181036000830152611c7b8161174f565b9050919050565b60006020820190508181036000830152611c9b8161178f565b9050919050565b60006020820190508181036000830152611cbb816117cf565b9050919050565b60006020820190508181036000830152611cdb81611835565b9050919050565b60006020820190508181036000830152611cfb81611875565b9050919050565b60006020820190508181036000830152611d1b816118cf565b9050919050565b60006020820190508181036000830152611d3b81611935565b9050919050565b60006020820190508181036000830152611d5b81611975565b9050919050565b60006020820190508181036000830152611d7b816119b5565b9050919050565b60006020820190508181036000830152611d9b816119f5565b9050919050565b60006020820190508181036000830152611dbb81611a5b565b9050919050565b60006020820190508181036000830152611ddb81611ac1565b9050919050565b60006020820190508181036000830152611dfb81611b01565b9050919050565b6000602082019050611e176000830184611b41565b92915050565b6000606082019050611e326000830186611b41565b611e3f6020830185611b41565b611e4c6040830184611b41565b949350505050565b6000819050602082019050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b6000611eb482611ebb565b9050919050565b6000819050919050565b611ece81611ea9565b8114611ed957600080fd5b50565b611ee581611ebb565b8114611ef057600080fd5b5056fea2646970667358221220dc2fb22fa3c137fd3e3656221c4fb242ec02b1065185161a20cb611741e9f3ef64736f6c63782b302e372e362d646576656c6f702e323032332e31312e362b636f6d6d69742e33313838663336632e6d6f64005c",
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

// GetBalanceOfDepositor is a free data retrieval call binding the contract method 0x83c65dc1.
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

// GetBalanceOfDepositor is a free data retrieval call binding the contract method 0x83c65dc1.
//
// Solidity: function getBalanceOfDepositor(address depositorAddress) view returns(uint256)
func (_Staking *StakingSession) GetBalanceOfDepositor(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetBalanceOfDepositor(&_Staking.CallOpts, depositorAddress)
}

// GetBalanceOfDepositor is a free data retrieval call binding the contract method 0x83c65dc1.
//
// Solidity: function getBalanceOfDepositor(address depositorAddress) view returns(uint256)
func (_Staking *StakingCallerSession) GetBalanceOfDepositor(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetBalanceOfDepositor(&_Staking.CallOpts, depositorAddress)
}

// GetDepositorCount is a free data retrieval call binding the contract method 0x337ce811.
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

// GetDepositorCount is a free data retrieval call binding the contract method 0x337ce811.
//
// Solidity: function getDepositorCount() view returns(uint256)
func (_Staking *StakingSession) GetDepositorCount() (*big.Int, error) {
	return _Staking.Contract.GetDepositorCount(&_Staking.CallOpts)
}

// GetDepositorCount is a free data retrieval call binding the contract method 0x337ce811.
//
// Solidity: function getDepositorCount() view returns(uint256)
func (_Staking *StakingCallerSession) GetDepositorCount() (*big.Int, error) {
	return _Staking.Contract.GetDepositorCount(&_Staking.CallOpts)
}

// GetDepositorOfValidator is a free data retrieval call binding the contract method 0x8d0d793d.
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

// GetDepositorOfValidator is a free data retrieval call binding the contract method 0x8d0d793d.
//
// Solidity: function getDepositorOfValidator(address validatorAddress) view returns(address)
func (_Staking *StakingSession) GetDepositorOfValidator(validatorAddress common.Address) (common.Address, error) {
	return _Staking.Contract.GetDepositorOfValidator(&_Staking.CallOpts, validatorAddress)
}

// GetDepositorOfValidator is a free data retrieval call binding the contract method 0x8d0d793d.
//
// Solidity: function getDepositorOfValidator(address validatorAddress) view returns(address)
func (_Staking *StakingCallerSession) GetDepositorOfValidator(validatorAddress common.Address) (common.Address, error) {
	return _Staking.Contract.GetDepositorOfValidator(&_Staking.CallOpts, validatorAddress)
}

// GetDepositorRewards is a free data retrieval call binding the contract method 0x29180463.
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

// GetDepositorRewards is a free data retrieval call binding the contract method 0x29180463.
//
// Solidity: function getDepositorRewards(address depositorAddress) view returns(uint256)
func (_Staking *StakingSession) GetDepositorRewards(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDepositorRewards(&_Staking.CallOpts, depositorAddress)
}

// GetDepositorRewards is a free data retrieval call binding the contract method 0x29180463.
//
// Solidity: function getDepositorRewards(address depositorAddress) view returns(uint256)
func (_Staking *StakingCallerSession) GetDepositorRewards(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDepositorRewards(&_Staking.CallOpts, depositorAddress)
}

// GetDepositorSlashings is a free data retrieval call binding the contract method 0x633fb13a.
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

// GetDepositorSlashings is a free data retrieval call binding the contract method 0x633fb13a.
//
// Solidity: function getDepositorSlashings(address depositorAddress) view returns(uint256)
func (_Staking *StakingSession) GetDepositorSlashings(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDepositorSlashings(&_Staking.CallOpts, depositorAddress)
}

// GetDepositorSlashings is a free data retrieval call binding the contract method 0x633fb13a.
//
// Solidity: function getDepositorSlashings(address depositorAddress) view returns(uint256)
func (_Staking *StakingCallerSession) GetDepositorSlashings(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDepositorSlashings(&_Staking.CallOpts, depositorAddress)
}

// GetNetBalanceOfDepositor is a free data retrieval call binding the contract method 0xc50873bd.
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

// GetNetBalanceOfDepositor is a free data retrieval call binding the contract method 0xc50873bd.
//
// Solidity: function getNetBalanceOfDepositor(address depositorAddress) view returns(uint256)
func (_Staking *StakingSession) GetNetBalanceOfDepositor(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetNetBalanceOfDepositor(&_Staking.CallOpts, depositorAddress)
}

// GetNetBalanceOfDepositor is a free data retrieval call binding the contract method 0xc50873bd.
//
// Solidity: function getNetBalanceOfDepositor(address depositorAddress) view returns(uint256)
func (_Staking *StakingCallerSession) GetNetBalanceOfDepositor(depositorAddress common.Address) (*big.Int, error) {
	return _Staking.Contract.GetNetBalanceOfDepositor(&_Staking.CallOpts, depositorAddress)
}

// GetTotalDepositedBalance is a free data retrieval call binding the contract method 0x3dff8b2e.
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

// GetTotalDepositedBalance is a free data retrieval call binding the contract method 0x3dff8b2e.
//
// Solidity: function getTotalDepositedBalance() view returns(uint256)
func (_Staking *StakingSession) GetTotalDepositedBalance() (*big.Int, error) {
	return _Staking.Contract.GetTotalDepositedBalance(&_Staking.CallOpts)
}

// GetTotalDepositedBalance is a free data retrieval call binding the contract method 0x3dff8b2e.
//
// Solidity: function getTotalDepositedBalance() view returns(uint256)
func (_Staking *StakingCallerSession) GetTotalDepositedBalance() (*big.Int, error) {
	return _Staking.Contract.GetTotalDepositedBalance(&_Staking.CallOpts)
}

// GetValidatorOfDepositor is a free data retrieval call binding the contract method 0xcbe25864.
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

// GetValidatorOfDepositor is a free data retrieval call binding the contract method 0xcbe25864.
//
// Solidity: function getValidatorOfDepositor(address depositorAddress) view returns(address)
func (_Staking *StakingSession) GetValidatorOfDepositor(depositorAddress common.Address) (common.Address, error) {
	return _Staking.Contract.GetValidatorOfDepositor(&_Staking.CallOpts, depositorAddress)
}

// GetValidatorOfDepositor is a free data retrieval call binding the contract method 0xcbe25864.
//
// Solidity: function getValidatorOfDepositor(address depositorAddress) view returns(address)
func (_Staking *StakingCallerSession) GetValidatorOfDepositor(depositorAddress common.Address) (common.Address, error) {
	return _Staking.Contract.GetValidatorOfDepositor(&_Staking.CallOpts, depositorAddress)
}

// ListValidators is a free data retrieval call binding the contract method 0x98d5635f.
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

// ListValidators is a free data retrieval call binding the contract method 0x98d5635f.
//
// Solidity: function listValidators() view returns(address[])
func (_Staking *StakingSession) ListValidators() ([]common.Address, error) {
	return _Staking.Contract.ListValidators(&_Staking.CallOpts)
}

// ListValidators is a free data retrieval call binding the contract method 0x98d5635f.
//
// Solidity: function listValidators() view returns(address[])
func (_Staking *StakingCallerSession) ListValidators() ([]common.Address, error) {
	return _Staking.Contract.ListValidators(&_Staking.CallOpts)
}

// AddDepositorReward is a paid mutator transaction binding the contract method 0x6d6b83d1.
//
// Solidity: function addDepositorReward(address depositorAddress, uint256 rewardAmount) returns(uint256)
func (_Staking *StakingTransactor) AddDepositorReward(opts *bind.TransactOpts, depositorAddress common.Address, rewardAmount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "addDepositorReward", depositorAddress, rewardAmount)
}

// AddDepositorReward is a paid mutator transaction binding the contract method 0x6d6b83d1.
//
// Solidity: function addDepositorReward(address depositorAddress, uint256 rewardAmount) returns(uint256)
func (_Staking *StakingSession) AddDepositorReward(depositorAddress common.Address, rewardAmount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.AddDepositorReward(&_Staking.TransactOpts, depositorAddress, rewardAmount)
}

// AddDepositorReward is a paid mutator transaction binding the contract method 0x6d6b83d1.
//
// Solidity: function addDepositorReward(address depositorAddress, uint256 rewardAmount) returns(uint256)
func (_Staking *StakingTransactorSession) AddDepositorReward(depositorAddress common.Address, rewardAmount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.AddDepositorReward(&_Staking.TransactOpts, depositorAddress, rewardAmount)
}

// AddDepositorSlashing is a paid mutator transaction binding the contract method 0xc8efe694.
//
// Solidity: function addDepositorSlashing(address depositorAddress, uint256 slashAmount) returns(uint256)
func (_Staking *StakingTransactor) AddDepositorSlashing(opts *bind.TransactOpts, depositorAddress common.Address, slashAmount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "addDepositorSlashing", depositorAddress, slashAmount)
}

// AddDepositorSlashing is a paid mutator transaction binding the contract method 0xc8efe694.
//
// Solidity: function addDepositorSlashing(address depositorAddress, uint256 slashAmount) returns(uint256)
func (_Staking *StakingSession) AddDepositorSlashing(depositorAddress common.Address, slashAmount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.AddDepositorSlashing(&_Staking.TransactOpts, depositorAddress, slashAmount)
}

// AddDepositorSlashing is a paid mutator transaction binding the contract method 0xc8efe694.
//
// Solidity: function addDepositorSlashing(address depositorAddress, uint256 slashAmount) returns(uint256)
func (_Staking *StakingTransactorSession) AddDepositorSlashing(depositorAddress common.Address, slashAmount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.AddDepositorSlashing(&_Staking.TransactOpts, depositorAddress, slashAmount)
}

// ChangeValidator is a paid mutator transaction binding the contract method 0x2b6c79ac.
//
// Solidity: function changeValidator(address newValidatorAddress) payable returns()
func (_Staking *StakingTransactor) ChangeValidator(opts *bind.TransactOpts, newValidatorAddress common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "changeValidator", newValidatorAddress)
}

// ChangeValidator is a paid mutator transaction binding the contract method 0x2b6c79ac.
//
// Solidity: function changeValidator(address newValidatorAddress) payable returns()
func (_Staking *StakingSession) ChangeValidator(newValidatorAddress common.Address) (*types.Transaction, error) {
	return _Staking.Contract.ChangeValidator(&_Staking.TransactOpts, newValidatorAddress)
}

// ChangeValidator is a paid mutator transaction binding the contract method 0x2b6c79ac.
//
// Solidity: function changeValidator(address newValidatorAddress) payable returns()
func (_Staking *StakingTransactorSession) ChangeValidator(newValidatorAddress common.Address) (*types.Transaction, error) {
	return _Staking.Contract.ChangeValidator(&_Staking.TransactOpts, newValidatorAddress)
}

// CompleteWithdrawal is a paid mutator transaction binding the contract method 0xd35319a7.
//
// Solidity: function completeWithdrawal() returns()
func (_Staking *StakingTransactor) CompleteWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "completeWithdrawal")
}

// CompleteWithdrawal is a paid mutator transaction binding the contract method 0xd35319a7.
//
// Solidity: function completeWithdrawal() returns()
func (_Staking *StakingSession) CompleteWithdrawal() (*types.Transaction, error) {
	return _Staking.Contract.CompleteWithdrawal(&_Staking.TransactOpts)
}

// CompleteWithdrawal is a paid mutator transaction binding the contract method 0xd35319a7.
//
// Solidity: function completeWithdrawal() returns()
func (_Staking *StakingTransactorSession) CompleteWithdrawal() (*types.Transaction, error) {
	return _Staking.Contract.CompleteWithdrawal(&_Staking.TransactOpts)
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0x75c8a157.
//
// Solidity: function initiateWithdrawal() returns()
func (_Staking *StakingTransactor) InitiateWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "initiateWithdrawal")
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0x75c8a157.
//
// Solidity: function initiateWithdrawal() returns()
func (_Staking *StakingSession) InitiateWithdrawal() (*types.Transaction, error) {
	return _Staking.Contract.InitiateWithdrawal(&_Staking.TransactOpts)
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0x75c8a157.
//
// Solidity: function initiateWithdrawal() returns()
func (_Staking *StakingTransactorSession) InitiateWithdrawal() (*types.Transaction, error) {
	return _Staking.Contract.InitiateWithdrawal(&_Staking.TransactOpts)
}

// NewDeposit is a paid mutator transaction binding the contract method 0x53447914.
//
// Solidity: function newDeposit(address validatorAddress) payable returns()
func (_Staking *StakingTransactor) NewDeposit(opts *bind.TransactOpts, validatorAddress common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "newDeposit", validatorAddress)
}

// NewDeposit is a paid mutator transaction binding the contract method 0x53447914.
//
// Solidity: function newDeposit(address validatorAddress) payable returns()
func (_Staking *StakingSession) NewDeposit(validatorAddress common.Address) (*types.Transaction, error) {
	return _Staking.Contract.NewDeposit(&_Staking.TransactOpts, validatorAddress)
}

// NewDeposit is a paid mutator transaction binding the contract method 0x53447914.
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

// FilterOnChangeValidator is a free log retrieval operation binding the contract event 0x955d76f15f2c7186fa66d87610def07550c910c010df50653927e4a708c7ba4d.
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

// WatchOnChangeValidator is a free log subscription operation binding the contract event 0x955d76f15f2c7186fa66d87610def07550c910c010df50653927e4a708c7ba4d.
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

// ParseOnChangeValidator is a log parse operation binding the contract event 0x955d76f15f2c7186fa66d87610def07550c910c010df50653927e4a708c7ba4d.
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

// FilterOnCompleteWithdrawal is a free log retrieval operation binding the contract event 0x6bbffc7f86996c75e5e2f4d2f237156a1db1eae0be2c9dfe495ca3709fd00d31.
//
// Solidity: event OnCompleteWithdrawal(address depositorAddress)
func (_Staking *StakingFilterer) FilterOnCompleteWithdrawal(opts *bind.FilterOpts) (*StakingOnCompleteWithdrawalIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnCompleteWithdrawal")
	if err != nil {
		return nil, err
	}
	return &StakingOnCompleteWithdrawalIterator{contract: _Staking.contract, event: "OnCompleteWithdrawal", logs: logs, sub: sub}, nil
}

// WatchOnCompleteWithdrawal is a free log subscription operation binding the contract event 0x6bbffc7f86996c75e5e2f4d2f237156a1db1eae0be2c9dfe495ca3709fd00d31.
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

// ParseOnCompleteWithdrawal is a log parse operation binding the contract event 0x6bbffc7f86996c75e5e2f4d2f237156a1db1eae0be2c9dfe495ca3709fd00d31.
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

// FilterOnIncreaseDeposit is a free log retrieval operation binding the contract event 0x028315a731be42717de3e7cf182e1b0d8a3f6112892c6304cddbeb647dcecfff.
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

// WatchOnIncreaseDeposit is a free log subscription operation binding the contract event 0x028315a731be42717de3e7cf182e1b0d8a3f6112892c6304cddbeb647dcecfff.
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

// ParseOnIncreaseDeposit is a log parse operation binding the contract event 0x028315a731be42717de3e7cf182e1b0d8a3f6112892c6304cddbeb647dcecfff.
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

// FilterOnInitiateWithdrawal is a free log retrieval operation binding the contract event 0xece30a6dff490a8d8300956508d30a7ccff98d80fea0d31985511ad8d5cd8d45.
//
// Solidity: event OnInitiateWithdrawal(address depositorAddress)
func (_Staking *StakingFilterer) FilterOnInitiateWithdrawal(opts *bind.FilterOpts) (*StakingOnInitiateWithdrawalIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OnInitiateWithdrawal")
	if err != nil {
		return nil, err
	}
	return &StakingOnInitiateWithdrawalIterator{contract: _Staking.contract, event: "OnInitiateWithdrawal", logs: logs, sub: sub}, nil
}

// WatchOnInitiateWithdrawal is a free log subscription operation binding the contract event 0xece30a6dff490a8d8300956508d30a7ccff98d80fea0d31985511ad8d5cd8d45.
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

// ParseOnInitiateWithdrawal is a log parse operation binding the contract event 0xece30a6dff490a8d8300956508d30a7ccff98d80fea0d31985511ad8d5cd8d45.
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

// FilterOnNewDeposit is a free log retrieval operation binding the contract event 0x3841518e49e6f4d4a82b0be1c23ded807f824d9c83ed7ec8d8c4530fe7e89cd3.
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

// WatchOnNewDeposit is a free log subscription operation binding the contract event 0x3841518e49e6f4d4a82b0be1c23ded807f824d9c83ed7ec8d8c4530fe7e89cd3.
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

// ParseOnNewDeposit is a log parse operation binding the contract event 0x3841518e49e6f4d4a82b0be1c23ded807f824d9c83ed7ec8d8c4530fe7e89cd3.
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

// FilterOnReward is a free log retrieval operation binding the contract event 0x05a6922ae8ed4fd1cdcf21d1f416eab35ef3bd50db7519644637e013f3bae74a.
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

// WatchOnReward is a free log subscription operation binding the contract event 0x05a6922ae8ed4fd1cdcf21d1f416eab35ef3bd50db7519644637e013f3bae74a.
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

// ParseOnReward is a log parse operation binding the contract event 0x05a6922ae8ed4fd1cdcf21d1f416eab35ef3bd50db7519644637e013f3bae74a.
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

// FilterOnSlashing is a free log retrieval operation binding the contract event 0x76a8947ad65877232453fd0ef3761a23a41155c575a7c12834bdda2711756f99.
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

// WatchOnSlashing is a free log subscription operation binding the contract event 0x76a8947ad65877232453fd0ef3761a23a41155c575a7c12834bdda2711756f99.
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

// ParseOnSlashing is a log parse operation binding the contract event 0x76a8947ad65877232453fd0ef3761a23a41155c575a7c12834bdda2711756f99.
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
