package falcon

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"testing"
)

var (
	testmsg1 = hexutil.MustDecode("0x68692074686572656f636b636861696e62626262626262626262626262626262")
	testmsg2 = hexutil.MustDecode("0x68692074686572656f636b636861696e62626262626262626262626262626261")
)

func TestFalcon_Basic(t *testing.T) {
	pubKey, priKey, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	digestHash1 := []byte(testmsg1)
	signature, err := Sign(priKey, digestHash1)
	if err != nil {
		t.Fatal(err)
	}

	err = Verify(digestHash1, signature, pubKey)
	if err != nil {
		t.Fatal(err)
	}

}
