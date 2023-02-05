package cryptobase

import "github.com/ethereum/go-ethereum/crypto/hybrid"

var SigAlg = hybrid.CreateHybridSig()

//var SigAlg = falcon.CreateFalconSig()

//var SigAlg = mocksignaturealgorithm.CreateMockSig()
