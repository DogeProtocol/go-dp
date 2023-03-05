// Package oqs provides a GO wrapper for the C liboqs quantum-resistant library.
//This file was added for go-dogep project (Doge Protocol Platform)

package oqs

/*
#cgo pkg-config: liboqs
#include <oqs/oqs.h>
*/
import "C"

import (
	"errors"
	"github.com/DogeProtocol/dp/crypto/keyestablishmentalgorithm"
	"math/big"
	"unsafe"
)

const KemName = "Kyber512" //sntrup761

var (
	ErrKemInitial              = errors.New("kem is not supported by OQS")
	ErrInvalidKemCiphertextLen = errors.New("invalid ciphertext length")
	ErrKemKeypairFailed        = errors.New("can not generate keypair")
	ErrEncapsulate             = errors.New("can not encapsulate secret")
	ErrDecapsulate             = errors.New("can not decapsulate secret")
	ErrInvalidKemPrivateKeyLen = errors.New("incorrect secret key length, make sure you " +
		"specify one in Init() or run GenerateKemKeyPair()")
	ErrInvalidKemPublicKeyLen = errors.New("invalid public key length")
)

// IsKEMEnabled returns true if a KEM algorithm is enabled, and false otherwise.
func IsKEMEnabled(algName string) bool {
	result := C.OQS_KEM_alg_is_enabled(C.CString(algName))
	return result != 0
}

// KeyEncapsulation defines the KEM main data structure.
type KeyEncapsulation struct {
	kem        *C.OQS_KEM
	secretKey  []byte
	AlgDetails KeyEncapsulationDetails
}

// KeyEncapsulationDetails defines the KEM algorithm details.
type KeyEncapsulationDetails struct {
	ClaimedNISTLevel   int
	IsINDCCA           bool
	LengthCiphertext   int
	LengthPublicKey    int
	LengthSecretKey    int
	LengthSharedSecret int
	Name               string
	Version            string
}

func (kem *KeyEncapsulation) Init(algName string, secretKey []byte) error {
	if !IsKEMEnabled(algName) {
		return ErrKemInitial
	}
	kem.kem = C.OQS_KEM_new(C.CString(algName))
	kem.secretKey = secretKey
	kem.AlgDetails.Name = C.GoString(kem.kem.method_name)
	kem.AlgDetails.Version = C.GoString(kem.kem.alg_version)
	kem.AlgDetails.ClaimedNISTLevel = int(kem.kem.claimed_nist_level)
	kem.AlgDetails.IsINDCCA = bool(kem.kem.ind_cca)
	kem.AlgDetails.LengthPublicKey = int(kem.kem.length_public_key)
	kem.AlgDetails.LengthSecretKey = int(kem.kem.length_secret_key)
	kem.AlgDetails.LengthCiphertext = int(kem.kem.length_ciphertext)
	kem.AlgDetails.LengthSharedSecret = int(kem.kem.length_shared_secret)
	return nil
}

// Details returns the KEM algorithm details.
func (kem *KeyEncapsulation) Details() KeyEncapsulationDetails {
	return kem.AlgDetails
}

func GenerateKemKeyPair() (*keyestablishmentalgorithm.PrivateKey, error) {
	kem := KeyEncapsulation{}
	defer kem.Clean() // clean up even in case of panic
	err := kem.Init(KemName, nil)
	if err != nil {
		return nil, err
	}
	privKey, err := kem.GenerateKemKeyPair()
	return privKey, err
}

func EncapSecret(publicKey []byte) (ciphertext, sharedSecret []byte, err error) {
	kem := KeyEncapsulation{}
	defer kem.Clean() // clean up even in case of panic
	err = kem.Init(KemName, nil)
	if err != nil {
		return nil, nil, err
	}
	ciphertext, sharedSecret, err = kem.EncapsulateSecret(publicKey)
	return ciphertext, sharedSecret, err
}

func DecapSecret(seckey, ciphertext []byte) ([]byte, error) {
	kem := KeyEncapsulation{}
	defer kem.Clean() // clean up even in case of panic
	err := kem.Init(KemName, seckey)
	if err != nil {
		return nil, err
	}
	sharedSecret, err := kem.DecapsulateSecret(ciphertext)
	return sharedSecret, err
}

func (kem *KeyEncapsulation) GenerateKemKeyPair() (*keyestablishmentalgorithm.PrivateKey, error) {
	publicKey := make([]byte, kem.AlgDetails.LengthPublicKey)
	kem.secretKey = make([]byte, kem.AlgDetails.LengthSecretKey)

	rv := C.OQS_KEM_keypair(kem.kem,
		(*C.uint8_t)(unsafe.Pointer(&publicKey[0])),
		(*C.uint8_t)(unsafe.Pointer(&kem.secretKey[0])))

	if rv != C.OQS_SUCCESS {
		return nil, ErrKemKeypairFailed
	}

	privy := new(keyestablishmentalgorithm.PrivateKey)
	privy.D = new(big.Int).SetBytes(kem.secretKey)
	privy.PublicKey.N = new(big.Int).SetBytes(publicKey)

	return privy, nil
}

// encapSecret encapsulates a secret using a public key and returns the
// corresponding ciphertext and shared secret.
func (kem *KeyEncapsulation) EncapsulateSecret(publicKey []byte) (ciphertext,
	sharedSecret []byte, err error) {
	if len(publicKey) != kem.AlgDetails.LengthPublicKey {
		return nil, nil, ErrInvalidKemPublicKeyLen
	}

	ciphertext = make([]byte, kem.AlgDetails.LengthCiphertext)
	sharedSecret = make([]byte, kem.AlgDetails.LengthSharedSecret)

	rv := C.OQS_KEM_encaps(kem.kem,
		(*C.uint8_t)(unsafe.Pointer(&ciphertext[0])),
		(*C.uint8_t)(unsafe.Pointer(&sharedSecret[0])),
		(*C.uint8_t)(unsafe.Pointer(&publicKey[0])))

	if rv != C.OQS_SUCCESS {
		return nil, nil, ErrEncapsulate
	}
	return ciphertext, sharedSecret, nil
}

// decapSecret decapsulates a ciphertexts and returns the corresponding shared
// secret.
func (kem *KeyEncapsulation) DecapsulateSecret(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) != kem.AlgDetails.LengthCiphertext {
		return nil, ErrInvalidKemCiphertextLen
	}
	if len(kem.secretKey) != kem.AlgDetails.LengthSecretKey {
		return nil, ErrInvalidKemPrivateKeyLen
	}

	sharedSecret := make([]byte, kem.AlgDetails.LengthSharedSecret)
	rv := C.OQS_KEM_decaps(kem.kem,
		(*C.uint8_t)(unsafe.Pointer(&sharedSecret[0])),
		(*C.uchar)(unsafe.Pointer(&ciphertext[0])),
		(*C.uint8_t)(unsafe.Pointer(&kem.secretKey[0])))

	if rv != C.OQS_SUCCESS {
		return nil, ErrDecapsulate
	}

	return sharedSecret, nil
}

func (kem *KeyEncapsulation) Clean() {
	if len(kem.secretKey) > 0 {
		MemCleanse(kem.secretKey)
	}
	C.OQS_KEM_free(kem.kem)
	*kem = KeyEncapsulation{}
}
