package systemcontracts1

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"strings"
)

var (
	testContract    = "0x0000000000000000000000000000000000000000"
	stakingContract = os.Getenv("GETH_STAKING_CONTRACT_ADDRESS")

	stakingContractABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"validatorId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnNewDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"OnWithdrawKey\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"depositBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"getDepositor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"listValidator\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"}],\"name\":\"newDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalDepositBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	stakingContractBIN = "0x608060405234801561001057600080fd5b50600080819055506000600181905550610f478061002f6000396000f3fe6080604052600436106100705760003560e01c806375697e661161004e57806375697e6614610106578063dfcd068f14610131578063e8c0a0df1461014d578063fba13bd01461017857610070565b80632dfdf0b5146100755780632e1a7d4d146100a05780636e2baf48146100c9575b600080fd5b34801561008157600080fd5b5061008a6101b5565b6040516100979190610d9d565b60405180910390f35b3480156100ac57600080fd5b506100c760048036038101906100c29190610a67565b6101be565b005b3480156100d557600080fd5b506100f060048036038101906100eb91906109f9565b610377565b6040516100fd9190610c6d565b60405180910390f35b34801561011257600080fd5b5061011b6103c7565b6040516101289190610ccd565b60405180910390f35b61014b60048036038101906101469190610a22565b610455565b005b34801561015957600080fd5b50610162610820565b60405161016f9190610d9d565b60405180910390f35b34801561018457600080fd5b5061019f600480360381019061019a91906109f9565b61082a565b6040516101ac9190610d9d565b60405180910390f35b60008054905090565b80600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015610240576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161023790610d5d565b60405180910390fd5b6102558160015461087390919063ffffffff16565b6001819055506102ad81600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461087390919063ffffffff16565b600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503373ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610336573d6000803e3d6000fd5b507f4d4666331ec61727075c5624fde25f5510c566e528d0565f2a2263a23b70d81a3382434260405161036c9493929190610c88565b60405180910390a150565b6000806103838361088a565b905060006004600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508092505050919050565b6060600680548060200260200160405190810160405280929190818152602001828054801561044b57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610401575b5050505050905090565b6000828290501161049b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161049290610d3d565b60405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff1660046000600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141561057c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161057390610d7d565b60405180910390fd5b61059260016000546108b190919063ffffffff16565b6000819055506105ad346001546108b190919063ffffffff16565b60018190555061060534600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546108b190919063ffffffff16565b600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555060008282600190809261065d93929190610e1e565b60405161066b929190610c54565b604051809103902090506000610680826108cd565b9050600061068d8261088a565b905084846003600084815260200190815260200160002091906106b19291906108da565b50336004600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506006829080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff16813373ffffffffffffffffffffffffffffffffffffffff167f9a1f4f083763f8508b19d4301c0110d2b47d99a8c5cf52c825c9e8cfea17f89c8888344342604051610811959493929190610cef565b60405180910390a45050505050565b6000600154905090565b6000600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b60008282111561087f57fe5b818303905092915050565b600060608273ffffffffffffffffffffffffffffffffffffffff16901b60001b9050919050565b6000808284019050838110156108c357fe5b8091505092915050565b60008160001c9050919050565b828054600181600116156101000203166002900490600052602060002090601f0160209004810192826109105760008555610957565b82601f1061092957803560ff1916838001178555610957565b82800160010185558215610957579182015b8281111561095657823582559160200191906001019061093b565b5b5090506109649190610968565b5090565b5b80821115610981576000816000905550600101610969565b5090565b60008135905061099481610ee3565b92915050565b60008083601f8401126109ac57600080fd5b8235905067ffffffffffffffff8111156109c557600080fd5b6020830191508360018202830111156109dd57600080fd5b9250929050565b6000813590506109f381610efa565b92915050565b600060208284031215610a0b57600080fd5b6000610a1984828501610985565b91505092915050565b60008060208385031215610a3557600080fd5b600083013567ffffffffffffffff811115610a4f57600080fd5b610a5b8582860161099a565b92509250509250929050565b600060208284031215610a7957600080fd5b6000610a87848285016109e4565b91505092915050565b6000610a9c8383610ab7565b60208301905092915050565b610ab181610e8d565b82525050565b610ac081610e51565b82525050565b610acf81610e51565b82525050565b6000610ae082610dc8565b610aea8185610de0565b9350610af583610db8565b8060005b83811015610b26578151610b0d8882610a90565b9750610b1883610dd3565b925050600181019050610af9565b5085935050505092915050565b6000610b3f8385610df1565b9350610b4c838584610ec3565b610b5583610ed2565b840190509392505050565b6000610b6c8385610e02565b9350610b79838584610ec3565b82840190509392505050565b6000610b92601583610e0d565b91507f5075626c6963206b657920697320696e76616c696400000000000000000000006000830152602082019050919050565b6000610bd2601283610e0d565b91507f496e73756666696369656e742066756e647300000000000000000000000000006000830152602082019050919050565b6000610c12601583610e0d565b91507f53656e64657220616c72656164792065786973747300000000000000000000006000830152602082019050919050565b610c4e81610e83565b82525050565b6000610c61828486610b60565b91508190509392505050565b6000602082019050610c826000830184610ac6565b92915050565b6000608082019050610c9d6000830187610aa8565b610caa6020830186610c45565b610cb76040830185610c45565b610cc46060830184610c45565b95945050505050565b60006020820190508181036000830152610ce78184610ad5565b905092915050565b60006080820190508181036000830152610d0a818789610b33565b9050610d196020830186610c45565b610d266040830185610c45565b610d336060830184610c45565b9695505050505050565b60006020820190508181036000830152610d5681610b85565b9050919050565b60006020820190508181036000830152610d7681610bc5565b9050919050565b60006020820190508181036000830152610d9681610c05565b9050919050565b6000602082019050610db26000830184610c45565b92915050565b6000819050602082019050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b60008085851115610e2e57600080fd5b83861115610e3b57600080fd5b6001850283019150848603905094509492505050565b6000610e5c82610e63565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b6000610e9882610e9f565b9050919050565b6000610eaa82610eb1565b9050919050565b6000610ebc82610e63565b9050919050565b82818337600083830152505050565b6000601f19601f8301169050919050565b610eec81610e51565b8114610ef757600080fd5b50565b610f0381610e83565b8114610f0e57600080fd5b5056fea2646970667358221220181b27743bf08caf1acd6da3c4bdbd6904a61bcc46be748c38e4a0ca4f2b4e5964736f6c63430007060033"

	systemContracts      []string
	systemContractsData  = make(map[string]*Contracts)
	systemContractVerify map[common.Address]bool
)

type Contracts struct {
	ContractAddressString string    `json:"ContractAddressString"`
	Contracts             *Contract `json:"Contracts"`
}

type Contract struct {
	ContractAddress common.Address `json:"ContractAddress"`
	ABI             string         `json:"ABI"`
	BIN             string         `json:"BIN"`
	Methods         *Method        `json:"Methods"`
}

type Method struct {
	Deposits   *Deposit   `json:"Deposits"`
	Validators *Validator `json:"Validators"`
}

type Deposit struct {
	GetDepositCount        string `json:"GetDepositCount"`
	GetTotalDepositBalance string `json:"GetTotalDepositBalance"`
}

type Validator struct {
	GetDepositBalanceOf string `json:"GetDepositBalanceOf"`
	ListValidator       string `json:"ListValidator"`
	GetDepositor        string `json:"GetDepositor"`
}

var (
	methods_collection = &Method{
		Deposits: &Deposit{
			GetDepositCount:        "depositCount",
			GetTotalDepositBalance: "totalDepositBalance",
		},
		Validators: &Validator{
			GetDepositBalanceOf: "depositBalanceOf",
			ListValidator:       "listValidator",
			GetDepositor:        "getDepositor",
		},
	}
)

func init() {
	systemContracts = []string{
		stakingContract,
	}

	systemContractsData[stakingContract] = &Contracts{
		ContractAddressString: stakingContract,
		Contracts: &Contract{
			ContractAddress: common.HexToAddress(stakingContract),
			ABI:             stakingContractABI,
			BIN:             stakingContractBIN,
			Methods:         methods_collection,
		},
	}

	systemContractVerify = map[common.Address]bool{
		common.HexToAddress(stakingContract): true,
	}
}

func GetContracts() []string {
	return systemContracts
}

func GetContract_Data(contract string) *Contract {
	return systemContractsData[contract].Contracts
}

func GetContractVerify(address common.Address) bool {
	return systemContractVerify[address]
}

func IsStakingContract() error {
	if len(stakingContract) < 40 {
		return fmt.Errorf("Staking contractor is not found")
	}
	return nil
}

func GetStakingContract_Address_String() string {
	return systemContractsData[stakingContract].ContractAddressString
}

func GetStakingContract_Address() common.Address {
	return systemContractsData[stakingContract].Contracts.ContractAddress
}

func GetStakingContract_ABI() abi.ABI {
	s := systemContractsData[stakingContract].Contracts.ABI
	abi, _ := abi.JSON(strings.NewReader(s))
	return abi
}

// Validators method

func GetContract_Method_ListValidator() string {
	return systemContractsData[stakingContract].Contracts.Methods.Validators.ListValidator
}

func GetContract_Method_GetDepositor() string {
	return systemContractsData[stakingContract].Contracts.Methods.Validators.GetDepositor
}
