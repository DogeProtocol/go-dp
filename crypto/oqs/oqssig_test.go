package oqs

import (
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"testing"
)

func TestOqsSig_Basic(t *testing.T) {
	InitOqs()

	var sig signaturealgorithm.SignatureAlgorithm
	sig = InitFalcon()

	signaturealgorithm.SignatureAlgorithmTest(t, sig)
}
