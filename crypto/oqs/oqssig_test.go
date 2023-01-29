package oqs

import (
	"github.com/ethereum/go-ethereum/crypto/signaturealgorithm"
	"testing"
)

func TestOqsSig_Basic(t *testing.T) {
	InitOqs()

	var sig signaturealgorithm.SignatureAlgorithm
	sig = InitFalcon()

	signaturealgorithm.SignatureAlgorithmTest(t, sig)
}
