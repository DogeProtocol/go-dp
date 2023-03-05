package mocksignaturealgorithm

import (
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"testing"
)

func TestMockSig_Basic(t *testing.T) {

	var sig signaturealgorithm.SignatureAlgorithm
	sig = CreateMockSig()

	signaturealgorithm.SignatureAlgorithmTest(t, sig)
}
