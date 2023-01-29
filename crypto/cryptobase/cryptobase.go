package cryptobase

import "github.com/ethereum/go-ethereum/crypto/falcon"

// var SigAlg = hybrid.CreateHybridSig()
var SigAlg = falcon.CreateFalconSig()
