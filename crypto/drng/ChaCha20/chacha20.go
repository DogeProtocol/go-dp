package ChaCha20

import (
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto"
	"golang.org/x/crypto/chacha20"
)

type ChaCha20DRNG struct {
	cipher *chacha20.Cipher // underlying ChaCha cipher
}

type ChaCha20DRNGInitializer struct {
}

func (g *ChaCha20DRNGInitializer) InitializeWithSeed(seed [common.HashLength]byte) (*ChaCha20DRNG, error) {

	// bounds restriction
	nonce := []byte("DogeProtocol")
	nonce = nonce[:chacha20.NonceSize]

	// create the underlying chacha20 cipher instance
	var key []byte
	if common.HashLength == chacha20.KeySize {
		key = seed[:]
	} else {
		key = crypto.Sha256(seed[:])
	}
	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}

	c := &ChaCha20DRNG{
		cipher: cipher,
	}

	return c, nil
}

func (g *ChaCha20DRNG) NextByte() byte {
	result := make([]byte, 1)
	g.cipher.XORKeyStream(result, result)
	return result[0]
}
