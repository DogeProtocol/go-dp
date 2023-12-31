package conversion

import (
	"github.com/DogeProtocol/dp/accounts/abi"
	"github.com/DogeProtocol/dp/common"
	"strings"
)

// Steps after Contract is modified
// 1) solc --bin --bin-runtime --abi c:\github\go-dp\systemcontracts\conversion\ConversionContract.sol  -o c:\github\go-dp\systemcontracts\conversion
// 2) abigen --bin=c:\github\go-dp\systemcontracts\conversion\ConversionContract.bin --abi=c:\github\go-dp\systemcontracts\conversion\ConversionContract.abi --pkg=conversion --out=c:\github\go-dp\systemcontracts\conversion\conversion.go
// 3) copy ConversionContract-runtime.bin into conversionbin.go CONVERSION_RUNTIME_BIN field
const CONVERSION_CONTRACT = "0x0000000000000000000000000000000000000000000000000000000000002000"

var CONVERSION_CONTRACT_ADDRESS = common.HexToAddress(CONVERSION_CONTRACT)

func GetContract_Method_getAmount() string {
	return "getAmount"
}

func GetContract_Method_getConversionStatus() string {
	return "getConversionStatus"
}

func GetContract_Method_getQuantumAddress() string {
	return "getQuantumAddress"
}

func GetContract_Method_setConverted() string {
	return "setConverted"
}

func GetConversionContract_ABI() (abi.ABI, error) {
	a, err := abi.JSON(strings.NewReader(ConversionMetaData.ABI))
	return a, err
}
