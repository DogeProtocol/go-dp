package oqs

import (
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"testing"
)

func TestDilithiumSig_Basic(t *testing.T) {
	InitOqs()

	var sig signaturealgorithm.SignatureAlgorithm
	sig = InitDilithium()

	signaturealgorithm.SignatureAlgorithmTest(t, sig)
}
