package commontest

import (
	"bytes"
	"fmt"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/hybrideds"
	"github.com/DogeProtocol/dp/crypto/hybridedsfull"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"testing"
)

var testmsg1 = "HELLO WORLD"

func TestCompactAndFullInterop(t *testing.T) {
	var sigFull signaturealgorithm.SignatureAlgorithm
	sigFull = hybridedsfull.CreateHybridedsfullSig()

	var sigCompact signaturealgorithm.SignatureAlgorithm
	sigCompact = hybrideds.CreateHybridedsSig(true)

	keyCompact, err := sigCompact.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}

	serializedCompact, err := sigCompact.SerializePrivateKey(keyCompact)
	if err != nil {
		t.Fatalf(err.Error())
	}

	keyFull, err := sigFull.DeserializePrivateKey(serializedCompact)
	if err != nil {
		t.Fatalf(err.Error())
	}

	addrCompact, err := sigCompact.PublicKeyToAddress(&keyCompact.PublicKey)
	if err != nil {
		t.Fatalf(err.Error())
	}

	addrFull, err := sigFull.PublicKeyToAddress(&keyFull.PublicKey)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if addrFull.IsEqualTo(addrCompact) == false {
		t.Fatalf("failed")
	}

	if bytes.Compare(keyCompact.PubData, keyFull.PubData) != 0 {
		t.Fatalf("failed")
	}

	digestHash1 := crypto.Keccak256([]byte(testmsg1))
	fmt.Println("digestHash1", len(digestHash1))
	signature1, err := sigFull.Sign(digestHash1, keyFull)
	if err != nil {
		fmt.Println(err)
		t.Fatal("Sign failed")
	}

	if sigFull.Verify(keyCompact.PubData, digestHash1, signature1) != true { //compact pub
		t.Fatal("Verify failed")
	}

	signature2, err := sigCompact.Sign(digestHash1, keyFull)
	if err != nil {
		fmt.Println(err)
		t.Fatal("Sign failed")
	}

	if sigCompact.Verify(keyFull.PubData, digestHash1, signature2) != true { //full pub
		t.Fatal("Verify failed")
	}
}
