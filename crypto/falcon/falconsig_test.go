package falcon

import (
	"github.com/ethereum/go-ethereum/crypto/signaturealgorithm"
	"testing"
)

func TestFalconSig_Basic(t *testing.T) {

	var sig signaturealgorithm.SignatureAlgorithm
	sig = CreateFalconSig()

	signaturealgorithm.SignatureAlgorithmTest(t, sig)
}
