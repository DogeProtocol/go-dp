package keyestablishmentalgorithm

import "math/big"

type PublicKey struct {
	N *big.Int // public key bytes
}

type PrivateKey struct {
	PublicKey          // public part.
	D         *big.Int // private key bytes
}
