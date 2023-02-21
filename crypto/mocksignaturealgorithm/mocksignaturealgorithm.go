package mocksignaturealgorithm

import "C"
import (
	"bytes"
	"errors"
	"math/rand"
)

const (
	OK                     = 0
	CRYPTO_SECRETKEY_BYTES = 64 + 1281 + 897
	CRYPTO_PUBLICKEY_BYTES = 32 + 897
	CRYPTO_SIGNATURE_BYTES = 64 + 1281 + 897 //Nonce + 2 for size
	CRYPTO_MESSAGE_LEN     = 32              //todo: validate this
	SIG_NAME               = "Falcon-512-ed25519"
)

var (
	ErrSignatureInitial       = errors.New("signature mechanism is not supported by OQS")
	ErrInvalidMsgLen          = errors.New("invalid message length, need max 64 bytes")
	ErrInvalidSignatureLen    = errors.New("invalid signature length")
	ErrInvalidPublicKeyLen    = errors.New("invalid public key length")
	ErrInvalidPrivateKeyLen   = errors.New("invalid private key length")
	ErrInvalidRecoveryID      = errors.New("invalid signature recovery id")
	ErrInvalidKey             = errors.New("invalid private key")
	ErrInvalidPubkey          = errors.New("invalid public key")
	ErrMismatchPublicKey      = errors.New("mismatch public key")
	ErrSignFailed             = errors.New("signing failed")
	ErrRecoverFailed          = errors.New("recovery failed")
	ErrKeypairFailed          = errors.New("can not generate keypair")
	ErrInvalidLen             = errors.New("invalid length")
	ErrVerifyFailed           = errors.New("verify length")
	ErrRecoverPublicKeyFailed = errors.New("recover public key length")
)

func GenerateKey() (publicKey []byte, secretKey []byte, err error) {
	publicKey = make([]byte, CRYPTO_PUBLICKEY_BYTES)
	secretKey = make([]byte, CRYPTO_SECRETKEY_BYTES)

	rand.Read(publicKey)
	copy(secretKey, publicKey)

	return publicKey[:CRYPTO_PUBLICKEY_BYTES], secretKey[:CRYPTO_SECRETKEY_BYTES], nil
}

func Sign(secretKey []byte, message []byte) ([]byte, error) {
	if len(secretKey) != CRYPTO_SECRETKEY_BYTES {
		return nil, ErrInvalidPrivateKeyLen
	}

	if len(message) != CRYPTO_MESSAGE_LEN {
		return nil, ErrInvalidMsgLen
	}

	signature := make([]byte, CRYPTO_SIGNATURE_BYTES)
	copy(signature, append(message, secretKey...))

	return signature, nil
}

// Verify verifies the validity of a signed message, returning true if the
// signature is valid, and false otherwise.
func Verify(message []byte, signature []byte, publicKey []byte) error {
	if len(message) != CRYPTO_MESSAGE_LEN || len(signature) == 0 || len(publicKey) == 0 {
		return ErrInvalidLen
	}
	if len(publicKey) != CRYPTO_PUBLICKEY_BYTES {
		return ErrInvalidPublicKeyLen
	}
	if len(signature) != CRYPTO_SIGNATURE_BYTES {
		return ErrInvalidSignatureLen
	}

	if bytes.Compare(signature[:len(message)], message) != 0 {
		return ErrVerifyFailed
	}

	if bytes.Compare(signature[len(message):len(message)+CRYPTO_PUBLICKEY_BYTES], publicKey) != 0 {
		return ErrVerifyFailed
	}

	return nil
}
