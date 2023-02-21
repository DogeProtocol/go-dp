package falcon

/*
#cgo pkg-config: libhybridpqc
#include <hybridpqc/api.h>
*/
import "C"
import (
	"bytes"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"unsafe"
)

const (
	OK                                 = 0
	CRYPTO_SECRETKEY_BYTES             = 1281
	CRYPTO_PUBLICKEY_BYTES             = 897
	CRYPTO_MESSAGE_LEN                 = 32                                //todo: validate this
	CRYPTO_SIGNATURE_BYTES             = 690 + 40 + 2 + CRYPTO_MESSAGE_LEN //Nonce + 2 for size
	CRYPTO_SIGNATURE_BYTES_WITH_LENGTH = CRYPTO_SIGNATURE_BYTES + common.LengthByteSize
	SIG_NAME                           = "Falcon-512"
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
	ErrVerifyFailed           = errors.New("verify length")
	ErrRecoverPublicKeyFailed = errors.New("recover public key length")
)

func GenerateKey() (publicKey []byte, secretKey []byte, err error) {
	publicKey = make([]byte, CRYPTO_PUBLICKEY_BYTES)
	secretKey = make([]byte, CRYPTO_SECRETKEY_BYTES)

	rv := C.crypto_sign_falcon_keypair(
		(*C.uchar)(unsafe.Pointer(&publicKey[0])),
		(*C.uchar)(unsafe.Pointer(&secretKey[0])))

	if rv != OK {
		return nil, nil, errors.New("GenerateKey failed")
	}
	return publicKey, secretKey, nil
}

func Sign(secretKey []byte, message []byte) ([]byte, error) {
	if len(secretKey) != CRYPTO_SECRETKEY_BYTES {
		return nil, ErrInvalidPrivateKeyLen
	}

	if len(message) == 0 || len(message) != CRYPTO_MESSAGE_LEN {
		return nil, ErrInvalidMsgLen
	}

	signature := make([]byte, CRYPTO_SIGNATURE_BYTES)

	var lenSig uint64

	rv := C.crypto_sign_falcon((*C.uchar)(unsafe.Pointer(&signature[0])),
		(*C.size_t)(unsafe.Pointer(&lenSig)),
		(*C.uchar)(unsafe.Pointer(&message[0])),
		(C.size_t)(uint64(len(message))),
		(*C.uchar)(unsafe.Pointer(&secretKey[0])))

	if rv != OK {
		return nil, ErrSignFailed
	}

	if lenSig > CRYPTO_SIGNATURE_BYTES {
		return nil, ErrInvalidSignatureLen
	}

	b := common.LenToBytes(int(lenSig))

	signature = append(b[:], signature...)

	return signature, nil
}

// Verify verifies the validity of a signed message, returning true if the
// signature is valid, and false otherwise.
func Verify(message []byte, signature []byte, publicKey []byte) error {

	if len(message) == 0 || len(signature) == 0 || len(publicKey) == 0 {
		return ErrInvalidLen
	}
	if len(publicKey) != CRYPTO_PUBLICKEY_BYTES {
		return ErrInvalidPublicKeyLen
	}
	if len(signature) != CRYPTO_SIGNATURE_BYTES_WITH_LENGTH {
		return ErrInvalidSignatureLen
	}

	lenSig := common.BytesToLen(signature[:common.LengthByteSize])
	if lenSig > CRYPTO_SIGNATURE_BYTES_WITH_LENGTH {
		return ErrInvalidSignatureLen
	}
	sigExtracted := signature[common.LengthByteSize : common.LengthByteSize+lenSig]
	msgLenCheck := 0

	messageCheck := make([]byte, CRYPTO_SIGNATURE_BYTES+len(message))

	rv := C.crypto_sign_falcon_open((*C.uchar)(unsafe.Pointer(&messageCheck[0])),
		(*C.size_t)(unsafe.Pointer(&msgLenCheck)),
		(*C.uchar)(unsafe.Pointer(&sigExtracted[0])),
		(C.size_t)(uint64(len(sigExtracted))),
		(*C.uchar)(unsafe.Pointer(&publicKey[0])))

	if rv != OK {
		return ErrVerifyFailed
	}

	if msgLenCheck != len(message) {
		return ErrVerifyFailed
	}
	if bytes.Compare(message, messageCheck[:msgLenCheck]) != 0 {
		return ErrVerifyFailed
	}

	return nil
}
