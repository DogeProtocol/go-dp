package params

import "math/big"
import "testing"

func Test_EtherToWei(t *testing.T) {
	expected := big.NewInt(1000000000000000000)
	if EtherToWei(big.NewInt(1)).Cmp(expected) != 0 {
		t.Fatalf("failed")
	}
}

func Test_WeiToEther(t *testing.T) {
	expected := big.NewInt(1)
	if WeiToEther(big.NewInt(1000000000000000000)).Cmp(expected) != 0 {
		t.Fatalf("failed")
	}
}
