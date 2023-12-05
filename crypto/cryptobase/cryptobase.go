package cryptobase

import (
	"github.com/DogeProtocol/dp/crypto/drng/ChaCha20"
	"github.com/DogeProtocol/dp/crypto/hybrideds"
)

var SigAlg = hybrideds.CreateHybridedsSig(true)

var DRNG = &ChaCha20.ChaCha20DRNGInitializer{}

//var SigAlg = mocksignaturealgorithm.CreateMockSig()
