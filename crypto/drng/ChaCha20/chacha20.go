package ChaCha20

import (
	"golang.org/x/crypto/chacha20"
)

type ChaCha20DRNG struct {
	cipher *chacha20.Cipher // underlying ChaCha cipher
}

type ChaCha20DRNGInitializer struct {
}

func (g *ChaCha20DRNGInitializer) InitializeWithSeed(seed [32]byte) (*ChaCha20DRNG, error) {

	// bounds restriction
	nonce := []byte("DogeProtocol")
	nonce = nonce[:chacha20.NonceSize]

	// create the underlying chacha20 cipher instance
	cipher, err := chacha20.NewUnauthenticatedCipher(seed[:], nonce)
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
