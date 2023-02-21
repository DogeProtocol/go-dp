package mocksignaturealgorithm

import (
	"github.com/ethereum/go-ethereum/crypto/signaturealgorithm"
	"testing"
)

func TestMockSig_Basic(t *testing.T) {

	var sig signaturealgorithm.SignatureAlgorithm
	sig = CreateMockSig()

	signaturealgorithm.SignatureAlgorithmTest(t, sig)
}
