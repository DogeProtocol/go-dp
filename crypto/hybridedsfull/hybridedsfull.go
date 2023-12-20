package hybridedsfull

/*
#cgo pkg-config: libhybridpqc
#include <dilithium/hybrid.h>
*/
import "C"

import (
	"bytes"
	"errors"
	"unsafe"
)

const (
	OK                     = 0
	CRYPTO_SECRETKEY_BYTES = 64 + 2560 + 1312 + 128
	CRYPTO_PUBLICKEY_BYTES = 32 + 1312 + 64
	CRYPTO_MESSAGE_LEN     = 32
	CRYPTO_SIGNATURE_BYTES = 2 + 64 + 2420 + 49856 + CRYPTO_MESSAGE_LEN //2558
	SIG_NAME               = "dilithium-ed25519-sphincs-full"
)

var (
	ErrSignatureInitial       = errors.New("signature mechanism is not supported by OQS")
	ErrInvalidMsgLen          = errors.New("invalid message length, need 32 bytes")
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
	ErrVerifyFailed           = errors.New("verify failed")
	ErrRecoverPublicKeyFailed = errors.New("recover public key length")
)

func GenerateKey() (publicKey []byte, secretKey []byte, err error) {
	publicKey = make([]byte, CRYPTO_PUBLICKEY_BYTES)
	secretKey = make([]byte, CRYPTO_SECRETKEY_BYTES)

	rv := C.crypto_sign_dilithium_ed25519_sphincs_keypair(
		(*C.uchar)(unsafe.Pointer(&publicKey[0])),
		(*C.uchar)(unsafe.Pointer(&secretKey[0])))

	if rv != OK {
		return nil, nil, errors.New("GenerateKey failed")
	}

	if bytes.Compare(publicKey[:32], secretKey[32:64]) != 0 {
		return nil, nil, ErrKeypairFailed
	}

	if bytes.Compare(publicKey[32:32+1312], secretKey[64+2560:64+2560+1312]) != 0 {
		return nil, nil, ErrKeypairFailed
	}

	if bytes.Compare(publicKey[32+1312:], secretKey[64+2560+1312+64:]) != 0 {
		return nil, nil, ErrKeypairFailed
	}

	return publicKey[:], secretKey[:], nil
}

func Sign(secretKey []byte, message []byte) ([]byte, error) {
	if len(secretKey) != CRYPTO_SECRETKEY_BYTES {
		return nil, ErrInvalidPrivateKeyLen
	}

	if len(message) != CRYPTO_MESSAGE_LEN {
		return nil, ErrInvalidMsgLen
	}

	signature := make([]byte, CRYPTO_SIGNATURE_BYTES)

	var lenSig uint64

	rv := C.crypto_sign_dilithium_ed25519_sphincs((*C.uchar)(unsafe.Pointer(&signature[0])),
		(*C.ulonglong)(unsafe.Pointer(&lenSig)),
		(*C.uchar)(unsafe.Pointer(&message[0])),
		(C.ulonglong)(uint64(len(message))),
		(*C.uchar)(unsafe.Pointer(&secretKey[0])))

	if rv != OK {
		return nil, ErrSignFailed
	}

	if lenSig != CRYPTO_SIGNATURE_BYTES {
		return nil, ErrInvalidSignatureLen
	}

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

	rv := C.crypto_verify_dilithium_ed25519_sphincs((*C.uchar)(unsafe.Pointer(&message[0])),
		(C.ulonglong)(uint64(len(message))),
		(*C.uchar)(unsafe.Pointer(&signature[0])),
		(C.ulonglong)(uint64(len(signature))),
		(*C.uchar)(unsafe.Pointer(&publicKey[0])))

	if rv != OK {
		return ErrVerifyFailed
	}

	return nil
}

func PrivateAndPublicFromPrivateKey(compositePrivateKey []byte) (privateBytes []byte, publicBytes []byte, err error) {

	if len(compositePrivateKey) != CRYPTO_SECRETKEY_BYTES {
		return nil, nil, ErrInvalidPrivateKeyLen
	}

	pub1 := make([]byte, len(compositePrivateKey[32:64]))
	copy(pub1, compositePrivateKey[32:64])

	pub2 := make([]byte, len(compositePrivateKey[64+2560:64+2560+1312]))
	copy(pub2, compositePrivateKey[64+2560:64+2560+1312])

	pub3 := make([]byte, len(compositePrivateKey[64+2560+1312+64:]))
	copy(pub3, compositePrivateKey[64+2560+1312+64:])

	pubKeyBytes := make([]byte, CRYPTO_PUBLICKEY_BYTES)
	pubKeyBytes = append(pub1, pub2...)
	pubKeyBytes = append(pubKeyBytes, pub3...)

	return compositePrivateKey, pubKeyBytes, nil
}
