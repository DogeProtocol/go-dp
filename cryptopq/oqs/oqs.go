// Package oqs provides a GO wrapper for the C liboqs quantum-resistant library.
//This file was added for go-dogep project (Doge Protocol Platform)

package oqs

/*
#cgo pkg-config: liboqs
#include <oqs/oqs.h>
*/
import "C"

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"unsafe"
)

const sigName = "Falcon-512"

const (
	SignerLength          = 8 // sign length store(dynamic length)
	PublicKeyLen          = 897
	PrivateKeyLen         = 1281
	SignatureLen          = 690
	SignPublicKeyLen      = 897 + 690 + SignerLength
	SignPublicKeyLenAbove = 1400
	PublicKeyStartVal     = 0x00 + 9
	SignKeyStartVal       = 0x30 + 9
)

var (
	ErrSignatureInitial     = errors.New("signature mechanism is not supported by OQS")
	ErrInvalidMsgLen        = errors.New("invalid message length, need 32 bytes")
	ErrInvalidSignatureLen  = errors.New("invalid signature length")
	ErrInvalidPublicKeyLen  = errors.New("invalid public key length")
	ErrInvalidPrivateKeyLen = errors.New("invalid private key length")
	ErrInvalidRecoveryID    = errors.New("invalid signature recovery id")
	ErrInvalidKey           = errors.New("invalid private key")
	ErrInvalidPubkey        = errors.New("invalid public key")
	ErrMismatchPublicKey    = errors.New("mismatch public key")
	ErrSignFailed           = errors.New("signing failed")
	ErrRecoverFailed        = errors.New("recovery failed")
	ErrKeypairFailed        = errors.New("can not generate keypair")
	ErrInvalidLen           = errors.New("invalid length")
)

// IsSigEnabled returns true if a signature algorithm is enabled, and false
// otherwise.
func IsSigEnabled(algName string) bool {
	result := C.OQS_SIG_alg_is_enabled(C.CString(algName))
	return result != 0
}

// Initializes the lists enabledSigs and supportedSigs.
func init() {

}

func (sig *Signature) Init(algName string, secretKey []byte) error {
	if !IsSigEnabled(algName) {
		return ErrSignatureInitial
	}

	sig.sig = C.OQS_SIG_new(C.CString(algName))
	sig.secretKey = secretKey
	sig.algDetails.Name = C.GoString(sig.sig.method_name)
	sig.algDetails.Version = C.GoString(sig.sig.alg_version)
	sig.algDetails.ClaimedNISTLevel = int(sig.sig.claimed_nist_level)
	sig.algDetails.IsEUFCMA = bool(sig.sig.euf_cma)
	sig.algDetails.LengthPublicKey = int(sig.sig.length_public_key)
	sig.algDetails.LengthSecretKey = int(sig.sig.length_secret_key)
	sig.algDetails.MaxLengthSignature = int(sig.sig.length_signature)

	return nil
}

// Details returns the signature algorithm details.
func (sig *Signature) Details() SignatureDetails {
	return sig.algDetails
}

// MemCleanse sets to zero the content of a byte slice by invoking the liboqs
// OQS_MEM_cleanse() function. Use it to clean "hot" memory areas, such as
// secret keys etc.
func MemCleanse(v []byte) {
	C.OQS_MEM_cleanse(unsafe.Pointer(&v[0]), C.size_t(len(v)))
}

// Clean zeroes-in the stored secret key and resets the sig receiver. One can
// reuse the signature by re-initializing it with the Signature.Init method.
func (sig *Signature) Clean() {
	if len(sig.secretKey) > 0 {
		MemCleanse(sig.secretKey)
	}
	C.OQS_SIG_free(sig.sig)
	*sig = Signature{}
}

// SignatureDetails defines the signature algorithm details.
type SignatureDetails struct {
	ClaimedNISTLevel   int
	IsEUFCMA           bool
	LengthPublicKey    int
	LengthSecretKey    int
	MaxLengthSignature int
	Name               string
	Version            string
}

// Signature defines the signature main data structure.
type Signature struct {
	sig        *C.OQS_SIG
	secretKey  []byte
	algDetails SignatureDetails
}

func SignWithKey(digestHash []byte, prv *PrivateKey) (sig []byte, err error) {
	seckey, err := ExportPrivateKey(prv)
	if err != nil {
		return nil, err
	}
	return Sign(digestHash, seckey)
}

func Sign(msg []byte, seckey []byte) ([]byte, error) {
	signer := Signature{}
	defer signer.Clean() // clean up even in case of panic
	err := signer.Init(sigName, seckey)
	if err != nil {
		return nil, err
	}

	sig, err := signer.sign(msg)
	return sig, err
}

func RecoverPubkey(msg []byte, sig []byte) ([]byte, error) {
	signer := Signature{}
	defer signer.Clean() // clean up even in case of panic
	err := signer.Init(sigName, nil)
	if err != nil {
		return nil, err
	}
	pubkey, err := signer.recoverPubkey(msg, sig)
	if err != nil {
		return nil, err
	}

	if VerifySignature(pubkey, msg, sig) == false {
		return nil, ErrInvalidPubkey
	}

	return pubkey, err
}

func VerifySignature(pubkey, msg, signature []byte) bool {
	signer := Signature{}
	defer signer.Clean() // clean up even in case of panic
	err := signer.Init(sigName, nil)
	if err != nil {
		return false
	}

	ok, err := signer.verify(msg, signature, pubkey)
	if err != nil {
		return false
	}

	return ok
}

func RecoverPubkeyByPrivate(seckey []byte) ([]byte, error) {
	signer := Signature{}
	defer signer.Clean() // clean up even in case of panic
	err := signer.Init(sigName, seckey)
	if err != nil {
		return nil, err
	}
	msg := []byte("dogeprotocoldogeprotocoldogeprotocol")
	sig, err := signer.sign(msg)
	pubkey, err := signer.recoverPubkey(msg, sig)
	return pubkey, err
}

// DecompressPubkey is NO-OP for now
func DecompressPubkey(pubkey []byte) (N *big.Int, err error) {
	signer := Signature{}
	defer signer.Clean() // clean up even in case of panic
	err = signer.Init(sigName, nil)
	if err != nil {
		return nil, err
	}
	pub, err := signer.convertBytesToPublic(pubkey)
	if err != nil {
		return nil, err
	}
	if pub == nil {
		return nil, nil
	}
	return pub.N, nil
}

// CompressPubkey is NO-OP for now
func CompressPubkey(N *big.Int) []byte {
	var (
		pubkey = N.Bytes()
	)
	return pubkey
}

func checkSignature(sig []byte) error {
	if len(sig) != SignatureLen {
		return ErrInvalidSignatureLen
	}
	if sig[SignatureLen-1] >= 4 {
		return ErrInvalidRecoveryID
	}
	return nil
}

func (sig *Signature) recoverPubkey(message []byte, signature []byte) ([]byte, error) {
	if len(message) == 0 || len(signature) == 0 || len(signature) < SignatureLen {
		return nil, ErrInvalidLen
	}

	var pubKey []byte
	startPos := len(signature) - sig.algDetails.LengthPublicKey
	if startPos < 0 {
		return nil, errors.New("recoverpubkey failed")
	}
	pubKey = signature[len(signature)-sig.algDetails.LengthPublicKey:]

	if len(pubKey) != sig.algDetails.LengthPublicKey {
		return nil, ErrInvalidPublicKeyLen
	}
	if len(signature) < sig.algDetails.MaxLengthSignature {
		return nil, ErrInvalidSignatureLen
	}

	signature = signature[SignerLength:]

	rv := C.OQS_SIG_verify_with_key(sig.sig, (*C.uint8_t)(unsafe.Pointer(&message[0])),
		C.size_t(len(message)), (*C.uint8_t)(unsafe.Pointer(&signature[0])),
		C.size_t(len(signature)), (*C.uint8_t)(unsafe.Pointer(&pubKey[0])))

	if rv != C.OQS_SUCCESS {
		return nil, ErrSignFailed
	}
	return pubKey, nil
}

// Sign signs a message and returns the corresponding signature.
func (sig *Signature) sign(message []byte) ([]byte, error) {
	if len(message) == 0 {
		return nil, ErrInvalidMsgLen
	}
	if len(sig.secretKey) != sig.algDetails.LengthSecretKey {
		return nil, ErrInvalidPrivateKeyLen
	}
	signature := make([]byte, sig.algDetails.MaxLengthSignature+
		sig.algDetails.LengthPublicKey)

	var lenSig uint64

	rv := C.OQS_SIG_sign_with_key(sig.sig, (*C.uint8_t)(unsafe.Pointer(&signature[0])),
		(*C.size_t)(unsafe.Pointer(&lenSig)),
		(*C.uint8_t)(unsafe.Pointer(&message[0])),
		C.size_t(len(message)), (*C.uint8_t)(unsafe.Pointer(&sig.secretKey[0])))

	if rv != C.OQS_SUCCESS {
		fmt.Print("rv", rv)
		return nil, ErrSignFailed
	}

	b := make([]byte, SignerLength)
	binary.LittleEndian.PutUint64(b, lenSig)
	signature = append(b, signature[:lenSig]...)

	return signature, nil
}

// Verify verifies the validity of a signed message, returning true if the
// signature is valid, and false otherwise.
func (sig *Signature) verify(message []byte, signature []byte,
	publicKey []byte) (bool, error) {

	if len(message) == 0 || len(signature) == 0 || len(publicKey) == 0 {
		return false, ErrInvalidLen
	}
	if len(publicKey) != sig.algDetails.LengthPublicKey {
		return false, ErrInvalidPublicKeyLen
	}
	if len(signature) < sig.algDetails.MaxLengthSignature {
		return false, ErrInvalidSignatureLen
	}

	signature = signature[SignerLength:]

	rv := C.OQS_SIG_verify_with_key(sig.sig, (*C.uint8_t)(unsafe.Pointer(&message[0])),
		C.size_t(len(message)), (*C.uint8_t)(unsafe.Pointer(&signature[0])),
		C.size_t(len(signature)), (*C.uint8_t)(unsafe.Pointer(&publicKey[0])))

	if rv != C.OQS_SUCCESS {
		return false, ErrSignFailed
	}

	return true, nil
}

// KeyDetails defines the signature main data structure.

// A PublicKey represents the public part
type PublicKey struct {
	N *big.Int // public key bytes
}

// A PrivateKey represents an assymetric key
type PrivateKey struct {
	PublicKey          // public part.
	D         *big.Int // private key bytes
}

// GenerateKey exports the corresponding secret key from the sig receiver.
func GenerateKey() (*PrivateKey, error) {
	signer := Signature{}
	defer signer.Clean() // clean up even in case of panic
	err := signer.Init(sigName, nil)
	if err != nil {
		return nil, err
	}
	privKey, err := signer.generateKey()
	return privKey, err
}

// ConvertBytesToPrivate exports the corresponding secret key from the sig receiver.
func ConvertBytesToPrivate(privy []byte) (*PrivateKey, error) {
	signer := Signature{}
	defer signer.Clean() // clean up even in case of panic
	err := signer.Init(sigName, privy)
	if err != nil {
		return nil, err
	}
	privKey, err := signer.convertBytesToPrivate(privy)
	return privKey, err
}

// ConvertBytesToPublic exports the corresponding secret key from the sig receiver.
func ConvertBytesToPublic(pub []byte) (*PublicKey, error) {
	signer := Signature{}
	defer signer.Clean() // clean up even in case of panic
	err := signer.Init(sigName, nil)
	if err != nil {
		return nil, err
	}
	pubKey, err := signer.convertBytesToPublic(pub)
	return pubKey, err
}

// ExportPrivateKey exports the corresponding secret key from the sig receiver.
func ExportPrivateKey(privy *PrivateKey) ([]byte, error) {
	signer := Signature{}
	defer signer.Clean() // clean up even in case of panic
	err := signer.Init(sigName, nil)
	if err != nil {
		return nil, err
	}
	privKey, err := signer.exportPrivateKey(privy)
	if err != nil {
		return nil, err
	}
	if privKey == nil {
		return nil, nil
	}
	return privKey, nil
}

// ExportPublicKey exports the corresponding secret key from the sig receiver.
func ExportPublicKey(pub *PublicKey) ([]byte, error) {
	signer := Signature{}
	defer signer.Clean() // clean up even in case of panic
	err := signer.Init(sigName, nil)
	if err != nil {
		return nil, err
	}
	pubKey, err := signer.exportPublicKey(pub)
	if err != nil {
		return nil, err
	}
	if pubKey == nil {
		return nil, nil
	}
	return pubKey, nil
}

// generateKey exports the corresponding secret key from the sig receiver.
func (sig *Signature) generateKey() (*PrivateKey, error) {
	publicKey := make([]byte, sig.algDetails.LengthPublicKey)
	sig.secretKey = make([]byte, sig.algDetails.LengthSecretKey)

	rv := C.OQS_SIG_keypair(sig.sig,
		(*C.uint8_t)(unsafe.Pointer(&publicKey[0])),
		(*C.uint8_t)(unsafe.Pointer(&sig.secretKey[0])))

	if rv != C.OQS_SUCCESS {
		return nil, ErrKeypairFailed
	}
	privy := new(PrivateKey)
	privy.D = new(big.Int).SetBytes(sig.secretKey)
	privy.PublicKey.N = new(big.Int).SetBytes(publicKey)
	return privy, nil
}

// convertBytesToPrivate exports the corresponding secret key from the sig receiver.
func (sig *Signature) convertBytesToPrivate(privy []byte) (*PrivateKey, error) {
	if len(privy) != int(sig.sig.length_secret_key) {
		return nil, ErrInvalidPrivateKeyLen
	}
	privKey := new(PrivateKey)
	privKey.D = new(big.Int).SetBytes(privy)
	return privKey, nil
}

// convertBytesToPublic exports the corresponding secret key from the sig receiver.
func (sig *Signature) convertBytesToPublic(pub []byte) (*PublicKey, error) {
	if len(pub) != int(sig.sig.length_public_key) {
		return nil, ErrInvalidPublicKeyLen
	}
	pubKey := new(PublicKey)
	pubKey.N = new(big.Int).SetBytes(pub)
	return pubKey, nil
}

// exportPrivateKey exports a private key into a binary dump.
func (sig *Signature) exportPrivateKey(privy *PrivateKey) ([]byte, error) {
	if len(privy.D.Bytes()) != int(sig.sig.length_secret_key) {
		return nil, ErrInvalidPrivateKeyLen
	}
	return privy.D.Bytes(), nil
}

// exportPublicKey exports a public key into a binary dump.
func (sig *Signature) exportPublicKey(pub *PublicKey) ([]byte, error) {
	if len(pub.N.Bytes()) != int(sig.sig.length_public_key) {
		return nil, ErrInvalidPublicKeyLen
	}
	return pub.N.Bytes(), nil
}
