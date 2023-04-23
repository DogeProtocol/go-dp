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
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"unsafe"
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

// IsSigEnabled returns true if a signature algorithm is enabled, and false
// otherwise.
func IsSigEnabled(algName string) bool {
	result := C.OQS_SIG_alg_is_enabled(C.CString(algName))
	return result != 0
}

// List of enabled KEM algorithms, populated by init().
var enabledKEMs []string

// List of supported KEM algorithms, populated by init().
var supportedKEMs []string

// MaxNumberKEMs returns the maximum number of supported KEM algorithms.
func MaxNumberKEMs() int {
	return int(C.OQS_KEM_alg_count())
}

// IsKEMSupported returns true if a KEM algorithm is supported, and false
// otherwise.
func IsKEMSupported(algName string) bool {
	for i := range supportedKEMs {
		if supportedKEMs[i] == algName {
			return true
		}
	}
	return false
}

// KEMName returns the KEM algorithm name from its corresponding numerical ID.
func KEMName(algID int) (string, error) {
	if algID >= MaxNumberKEMs() {
		return "", errors.New("algorithm ID out of range")
	}
	return C.GoString(C.OQS_KEM_alg_identifier(C.size_t(algID))), nil
}

// SupportedKEMs returns the list of supported KEM algorithms.
func SupportedKEMs() []string {
	return supportedKEMs
}

// EnabledKEMs returns the list of enabled KEM algorithms.
func EnabledKEMs() []string {
	return enabledKEMs
}

// Initializes liboqs and the lists enabledKEMs and supportedKEMs.
func InitOqs() {
	C.OQS_init()
	for i := 0; i < MaxNumberKEMs(); i++ {
		KEMName, _ := KEMName(i)
		supportedKEMs = append(supportedKEMs, KEMName)
		if IsKEMEnabled(KEMName) {
			enabledKEMs = append(enabledKEMs, KEMName)
		}
	}
}

func (sig *Signature) Init(algName string, secretKey []byte) error {
	if !IsSigEnabled(algName) {
		return ErrSignatureInitial
	}

	sig.sig = C.OQS_SIG_new(C.CString(algName))
	sig.secretKey = secretKey
	sig.AlgDetails.Name = C.GoString(sig.sig.method_name)
	sig.AlgDetails.Version = C.GoString(sig.sig.alg_version)
	sig.AlgDetails.ClaimedNISTLevel = int(sig.sig.claimed_nist_level)
	sig.AlgDetails.IsEUFCMA = bool(sig.sig.euf_cma)
	sig.AlgDetails.LengthPublicKey = int(sig.sig.length_public_key)
	sig.AlgDetails.LengthSecretKey = int(sig.sig.length_secret_key)
	sig.AlgDetails.MaxLengthSignature = int(sig.sig.length_signature) + common.LengthByteSize
	sig.AlgDetails.maxLengthSignatureInternal = int(sig.sig.length_signature)

	return nil
}

// Details returns the signature algorithm details.
func (sig *Signature) Details() SignatureDetails {
	return sig.AlgDetails
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
	ClaimedNISTLevel           int
	IsEUFCMA                   bool
	LengthPublicKey            int
	LengthSecretKey            int
	MaxLengthSignature         int
	Name                       string
	Version                    string
	maxLengthSignatureInternal int
}

// Signature defines the signature main data structure.
type Signature struct {
	sig        *C.OQS_SIG
	secretKey  []byte
	AlgDetails SignatureDetails
}

func SignWithKey(sigName string, digestHash []byte, prv *signaturealgorithm.PrivateKey) (sig []byte, err error) {
	seckey, err := ExportPrivateKey(sigName, prv)
	if err != nil {
		return nil, err
	}
	return Sign(sigName, digestHash, seckey)
}

func Sign(sigName string, msg []byte, seckey []byte) ([]byte, error) {
	signer := Signature{}
	defer signer.Clean() // clean up even in case of panic
	err := signer.Init(sigName, seckey)
	if err != nil {
		return nil, err
	}

	sig, err := signer.sign(msg)
	return sig, err
}

func VerifySignature(sigName string, pubkey, msg, signature []byte) bool {
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

// Sign signs a message and returns the corresponding signature.
func (sig *Signature) sign(message []byte) ([]byte, error) {
	if len(message) == 0 {
		return nil, ErrInvalidMsgLen
	}
	if len(sig.secretKey) != sig.AlgDetails.LengthSecretKey {
		return nil, ErrInvalidPrivateKeyLen
	}
	signature := make([]byte, sig.AlgDetails.maxLengthSignatureInternal)

	var lenSig uint64

	rv := C.OQS_SIG_sign(sig.sig, (*C.uint8_t)(unsafe.Pointer(&signature[0])),
		(*C.size_t)(unsafe.Pointer(&lenSig)),
		(*C.uint8_t)(unsafe.Pointer(&message[0])),
		C.size_t(len(message)), (*C.uint8_t)(unsafe.Pointer(&sig.secretKey[0])))

	if rv != C.OQS_SUCCESS {
		return nil, ErrSignFailed
	}

	b := common.LenToBytes(int(lenSig))
	signature = append(b[:], signature...)

	return signature, nil
}

// Verify verifies the validity of a signed message, returning true if the
// signature is valid, and false otherwise.
func (sig *Signature) verify(message []byte, signature []byte,
	publicKey []byte) (bool, error) {

	if len(message) == 0 || len(signature) == 0 || len(publicKey) == 0 {
		return false, ErrInvalidLen
	}
	if len(publicKey) != sig.AlgDetails.LengthPublicKey {
		return false, ErrInvalidPublicKeyLen
	}
	if len(signature) != sig.AlgDetails.MaxLengthSignature {
		return false, ErrInvalidSignatureLen
	}

	lenSig := common.BytesToLen(signature[:common.LengthByteSize])

	if lenSig > sig.AlgDetails.maxLengthSignatureInternal {
		return false, errors.New("invalid length")
	}
	sigExtracted := signature[common.LengthByteSize : common.LengthByteSize+lenSig]

	rv := C.OQS_SIG_verify(sig.sig, (*C.uint8_t)(unsafe.Pointer(&message[0])),
		C.size_t(len(message)), (*C.uint8_t)(unsafe.Pointer(&sigExtracted[0])),
		C.size_t(len(sigExtracted)), (*C.uint8_t)(unsafe.Pointer(&publicKey[0])))

	if rv != C.OQS_SUCCESS {
		return false, ErrVerifyFailed
	}

	return true, nil
}

// KeyDetails defines the signature main data structure.

// GenerateKey exports the corresponding secret key from the sig receiver.
func GenerateKey(sigName string) (*signaturealgorithm.PrivateKey, error) {
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
func ConvertBytesToPrivate(sigName string, privy []byte) (*signaturealgorithm.PrivateKey, error) {
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
func ConvertBytesToPublic(sigName string, pub []byte) (*signaturealgorithm.PublicKey, error) {
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
func ExportPrivateKey(sigName string, privy *signaturealgorithm.PrivateKey) ([]byte, error) {
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

func GetSignatureDetails(sigName string) (SignatureDetails, error) {
	signer := Signature{}
	err := signer.Init(sigName, nil)
	return signer.AlgDetails, err
}

// ExportPublicKey exports the corresponding secret key from the sig receiver.
func ExportPublicKey(sigName string, pub *signaturealgorithm.PublicKey) ([]byte, error) {
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
func (sig *Signature) generateKey() (*signaturealgorithm.PrivateKey, error) {
	publicKey := make([]byte, sig.AlgDetails.LengthPublicKey)
	sig.secretKey = make([]byte, sig.AlgDetails.LengthSecretKey)

	rv := C.OQS_SIG_keypair(sig.sig,
		(*C.uint8_t)(unsafe.Pointer(&publicKey[0])),
		(*C.uint8_t)(unsafe.Pointer(&sig.secretKey[0])))

	if rv != C.OQS_SUCCESS {
		return nil, ErrKeypairFailed
	}

	privy := new(signaturealgorithm.PrivateKey)
	privy.PriData = make([]byte, len(sig.secretKey))
	copy(privy.PriData, sig.secretKey)

	privy.PublicKey.PubData = make([]byte, len(publicKey))
	copy(privy.PublicKey.PubData, publicKey)

	return privy, nil
}

// convertBytesToPrivate exports the corresponding secret key from the sig receiver.
func (sig *Signature) convertBytesToPrivate(privy []byte) (*signaturealgorithm.PrivateKey, error) {
	if len(privy) != sig.AlgDetails.LengthSecretKey {
		return nil, ErrInvalidPrivateKeyLen
	}
	privKey := new(signaturealgorithm.PrivateKey)
	privKey.PriData = make([]byte, sig.AlgDetails.LengthSecretKey)
	copy(privKey.PriData, privy)

	return privKey, nil
}

// convertBytesToPublic exports the corresponding secret key from the sig receiver.
func (sig *Signature) convertBytesToPublic(pub []byte) (*signaturealgorithm.PublicKey, error) {
	if len(pub) != sig.AlgDetails.LengthPublicKey {
		return nil, ErrInvalidPublicKeyLen
	}
	pubKey := new(signaturealgorithm.PublicKey)
	pubKey.PubData = make([]byte, sig.AlgDetails.LengthPublicKey)
	copy(pubKey.PubData, pub)
	return pubKey, nil
}

// exportPrivateKey exports a private key into a binary dump.
func (sig *Signature) exportPrivateKey(privy *signaturealgorithm.PrivateKey) ([]byte, error) {
	if len(privy.PriData) != sig.AlgDetails.LengthSecretKey {
		return nil, ErrInvalidPrivateKeyLen
	}

	buf := make([]byte, sig.AlgDetails.LengthSecretKey)
	copy(buf, privy.PriData)
	return buf, nil
}

// exportPublicKey exports a public key into a binary dump.
func (sig *Signature) exportPublicKey(pub *signaturealgorithm.PublicKey) ([]byte, error) {
	if len(pub.PubData) != sig.AlgDetails.LengthPublicKey {
		return nil, ErrInvalidPublicKeyLen
	}
	buf := make([]byte, sig.AlgDetails.LengthPublicKey)
	copy(buf, pub.PubData)
	return buf, nil
}
