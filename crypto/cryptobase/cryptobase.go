package cryptobase

import (
	"github.com/DogeProtocol/dp/crypto/drng/ChaCha20"
	"github.com/DogeProtocol/dp/crypto/hybrid"
)

var SigAlg = hybrid.CreateHybridSig(true)

var DRNG = &ChaCha20.ChaCha20DRNGInitializer{}

//var SigAlg = falcon.CreateFalconSig()

//var SigAlg = mocksignaturealgorithm.CreateMockSig()
