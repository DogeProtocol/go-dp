package systemcontracts

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSystemContracts(t *testing.T) {
	for _, c := range systemContracts {
		fmt.Println(c)
		assert.Equal(t, c, "0x0000000000000000000000000000000000001000")
	}
}

func TestSystemContractsFail(t *testing.T) {
	for _, c := range systemContracts {
		fmt.Println(c)
		assert.NotEqual(t, c, "0x0000000000000000000000000000000000000000")
	}
}

func TestSystemContractsData(t *testing.T) {
	c := systemContractsData["0x0000000000000000000000000000000000001000"]
	fmt.Println(c)
	assert.Equal(t, c, systemContractsData["0x0000000000000000000000000000000000001000"])
}

func TestSystemContractsDataFail(t *testing.T) {
	c := systemContractsData["0x0000000000000000000000000000000000001000"]
	fmt.Println(c)
	assert.NotEqual(t, c, systemContractsData["0x0000000000000000000000000000000000000000"])
}

func TestSystemContractVerify(t *testing.T) {
	s := systemContractVerify[common.HexToAddress(systemContractsData["0x0000000000000000000000000000000000001000"].ContractAddressString)]
	fmt.Println(s)
	assert.Equal(t, s, true)
}

func TestSystemContractVerifyFail(t *testing.T) {
	s := systemContractVerify[common.HexToAddress("0x0000000000000000000000000000000000000000")]
	fmt.Println(s)
	assert.Equal(t, s, false)
}
