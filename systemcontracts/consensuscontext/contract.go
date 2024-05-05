package consensuscontext

import (
	"github.com/DogeProtocol/dp/common"
	"strings"
)
import "github.com/DogeProtocol/dp/accounts/abi"

// Steps after Contract is modified
// 1) solc --bin --bin-runtime --abi c:\github\go-dp\systemcontracts\consensuscontext\ConsensusContextContract.sol  -o c:\github\go-dp\systemcontracts\consensuscontext
// 2) abigen --bin=c:\github\go-dp\systemcontracts\consensuscontext\ConsensusContextContract.bin --abi=c:\github\go-dp\systemcontracts\consensuscontext\ConsensusContextContract.abi --pkg=consensuscontext --out=c:\github\go-dp\systemcontracts\consensuscontext\ConsensusContext.go
// 3) copy ConsensusContextContract-runtime.bin into consensuscontextbin.go CONTEXT_RUNTIME_BIN field

const CONSENSUS_CONTEXT_CONTRACT = "0x0000000000000000000000000000000000000000000000000000000000003000"

var CONSENSUS_CONTEXT_CONTRACT_ADDRESS = common.HexToAddress(CONSENSUS_CONTEXT_CONTRACT)

const SET_CONTEXT_METHOD = "setContext"
const GET_CONTEXT_METHOD = "getContext"
const DELETE_CONTEXT_METHOD = "deleteContext"

func GetConsensusContract_ABI() (abi.ABI, error) {
	s := ConsensuscontextMetaData.ABI
	abi, err := abi.JSON(strings.NewReader(s))
	return abi, err
}
