package systemcontracts_test

import (
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/systemcontracts"
	"testing"
)

func TestFunctions(t *testing.T) {
	//fmt.Println(common.HexToAddress("0xB5c2F2779716bBa6Ba9B4372501208110581EDec").Bytes())

	cont := systemcontracts.GetContracts()
	fmt.Println("Contracts : ", cont)
	verify := systemcontracts.GetContractVerify(common.HexToAddress("0x0000000000000000000000000000000000001000"))
	fmt.Println("Verify contract bool : ", verify)
	verify = systemcontracts.GetContractVerify(common.HexToAddress("0x0000000000000000000000000000000000000000"))
	fmt.Println("Verify contract bool : ", verify)
	c := systemcontracts.GetStakingContract_Address_String()
	fmt.Println("Contract address string : ", c)
	v := systemcontracts.GetStakingContract_Address()
	fmt.Println("Contract address ", v)
	s := systemcontracts.GetContract_Method_ListValidator()
	fmt.Println("Method", s)
	a := systemcontracts.GetStakingContract_ABI()
	fmt.Println("abi", a)

}
