package staking

import (
	"encoding/hex"
	"testing"
)

func TestStakingContract(t *testing.T) {
	newContractCode, err := hex.DecodeString(STAKING_BIN)
	if err != nil {
		t.Fatal(err)
	}
	if len(newContractCode)%32 != 0 {
		t.Errorf("len(data) is %d, want a multiple of 32", len(newContractCode))
	}
}
