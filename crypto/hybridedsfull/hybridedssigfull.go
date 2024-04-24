package hybridedsfull

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"io"
	"io/ioutil"
	"math/big"
	"os"
)

const CRYPTO_ED25519_PUBLICKEY_BYTES = 32
const CRYPTO_ED25519_SIGNATURE_BYTES = 64

const CRYPTO_DILITHIUM_PUBLICKEY_BYTES = 1312
const CRYPTO_DILITHIUM_SIGNATURE_BYTES = 2420

const CRYPTO_SPHINCS_PUBLICKEY_BYTES = 64
const CRYPTO_SPHINCS_SIGNATURE_BYTES = 49856

const CRYPTO_HYBRID_SIGNATURE_BYTES = 2 + 64 + 2420 + 49856 //+ MESSAGE_LEN
const SIGNATURE_ID = 2

type HybridedsfullSig struct {
	sigName                      string
	publicKeyBytesIndexStart     int
	publicKeyLength              int
	privateKeyLength             int
	signatureLength              int
	signatureWithPublicKeyLength int
}

func CreateHybridedsfullSig() HybridedsfullSig {
	return HybridedsfullSig{sigName: SIG_NAME,
		publicKeyBytesIndexStart:     12,
		publicKeyLength:              CRYPTO_PUBLICKEY_BYTES,
		privateKeyLength:             CRYPTO_SECRETKEY_BYTES,
		signatureLength:              CRYPTO_SIGNATURE_BYTES,
		signatureWithPublicKeyLength: CRYPTO_PUBLICKEY_BYTES + CRYPTO_SIGNATURE_BYTES + common.LengthByteSize + common.LengthByteSize,
	}
}

func (s HybridedsfullSig) SignatureName() string {
	return s.sigName
}

func (s HybridedsfullSig) PublicKeyLength() int {
	return s.publicKeyLength
}

func (s HybridedsfullSig) PrivateKeyLength() int {
	return s.privateKeyLength
}

func (s HybridedsfullSig) SignatureLength() int {
	return s.signatureLength
}

func (s HybridedsfullSig) SignatureWithPublicKeyLength() int {
	return s.signatureWithPublicKeyLength
}

func (s HybridedsfullSig) GenerateKey() (*signaturealgorithm.PrivateKey, error) {
	pubKey, priKey, err := GenerateKey()
	if err != nil {
		return nil, err
	}

	if len(pubKey) != s.publicKeyLength || len(priKey) != s.privateKeyLength {
		panic("keygen basic check failed")
	}

	privy := new(signaturealgorithm.PrivateKey)
	privy.PriData = make([]byte, len(priKey))
	copy(privy.PriData, priKey)

	privy.PublicKey.PubData = make([]byte, len(pubKey))
	copy(privy.PublicKey.PubData, pubKey)

	return privy, nil
}

func (s HybridedsfullSig) SerializePrivateKey(priv *signaturealgorithm.PrivateKey) ([]byte, error) {
	priBytes, err := s.exportPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	return priBytes, err
}

func (s HybridedsfullSig) DeserializePrivateKey(priv []byte) (*signaturealgorithm.PrivateKey, error) {

	privKeyBytes, pubKeyBytes, err := PrivateAndPublicFromPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	privKey, err := s.convertBytesToPrivate(privKeyBytes)
	if err != nil {
		return nil, err
	}

	pubkey, err := s.convertBytesToPublic(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	privKey.PublicKey = *pubkey

	return privKey, err
}

func (s HybridedsfullSig) SerializePublicKey(pub *signaturealgorithm.PublicKey) ([]byte, error) {
	return s.exportPublicKey(pub)
}

func (s HybridedsfullSig) DeserializePublicKey(pub []byte) (*signaturealgorithm.PublicKey, error) {
	pubKey, error := s.convertBytesToPublic(pub)
	return pubKey, error
}

func (s HybridedsfullSig) HexToPrivateKey(hexkey string) (*signaturealgorithm.PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, err
	}

	if byteErr, ok := err.(hex.InvalidByteError); ok {
		return nil, fmt.Errorf("invalid hex character %q in private key", byte(byteErr))
	} else if err != nil {
		return nil, errors.New("invalid hex data for private key")
	}
	return s.DeserializePrivateKey(b)
}

func (s HybridedsfullSig) HexToPrivateKeyNoError(hexkey string) *signaturealgorithm.PrivateKey {
	p, err := s.HexToPrivateKey(hexkey)
	if err != nil {
		panic("HexToPrivateKey")
	}
	return p
}

func (s HybridedsfullSig) PrivateKeyToHex(priv *signaturealgorithm.PrivateKey) (string, error) {
	data, err := s.SerializePrivateKey(priv)
	if err != nil {
		return "", err
	}
	k := hex.EncodeToString(data)
	return k, nil
}

func (s HybridedsfullSig) PublicKeyToHex(pub *signaturealgorithm.PublicKey) (string, error) {
	data, err := s.SerializePublicKey(pub)
	if err != nil {
		return "", err
	}
	k := hex.EncodeToString(data)
	return k, nil
}

func (s HybridedsfullSig) HexToPublicKey(hexkey string) (*signaturealgorithm.PublicKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, err
	}

	if byteErr, ok := err.(hex.InvalidByteError); ok {
		return nil, fmt.Errorf("invalid hex character %q in private key", byte(byteErr))
	} else if err != nil {
		return nil, errors.New("invalid hex data for private key")
	}
	return s.DeserializePublicKey(b)
}

func (s HybridedsfullSig) LoadPrivateKeyFromFile(file string) (*signaturealgorithm.PrivateKey, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	r := bufio.NewReader(fd)
	buf := make([]byte, (s.privateKeyLength)*2)
	n, err := readASCII(buf, r)
	if err != nil {
		return nil, err
	} else if n != len(buf) {
		return nil, fmt.Errorf("key file too short, want oqs hex character")
	}
	if err := checkKeyFileEnd(r); err != nil {
		return nil, err
	}
	return s.HexToPrivateKey(string(buf))
}

func (s HybridedsfullSig) SavePrivateKeyToFile(file string, key *signaturealgorithm.PrivateKey) error {
	k, err := s.PrivateKeyToHex(key)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(k), 0600)
}

func (s HybridedsfullSig) PublicKeyToAddress(p *signaturealgorithm.PublicKey) (common.Address, error) {
	pubBytes, err := s.SerializePublicKey(p)
	tempAddr := common.Address{}
	if err != nil {
		return tempAddr, err
	}
	return crypto.PublicKeyBytesToAddress(pubBytes), nil
}

func (s HybridedsfullSig) PublicKeyToAddressNoError(p *signaturealgorithm.PublicKey) common.Address {
	addr, err := s.PublicKeyToAddress(p)
	if err != nil {
		panic("PublicKeyToAddress failed")
	}
	return addr
}

func (s HybridedsfullSig) Sign(digestHash []byte, prv *signaturealgorithm.PrivateKey) (sig []byte, err error) {
	seckey, err := s.exportPrivateKey(prv)
	if err != nil {
		return nil, err
	}

	sigBytes, err := Sign(seckey, digestHash)
	if err != nil {
		return nil, err
	}

	pubBytes, err := s.SerializePublicKey(&prv.PublicKey)
	if err != nil {
		return nil, err
	}

	combinedSignature := common.CombineTwoParts(sigBytes, pubBytes)

	if !s.Verify(pubBytes, digestHash, combinedSignature) {
		return nil, errors.New("Verify failed after signing")
	}

	return combinedSignature, nil
}

func (s HybridedsfullSig) SignWithContext(digestHash []byte, prv *signaturealgorithm.PrivateKey, context []byte) (sig []byte, err error) {
	if context == nil || len(context) < 1 {
		return nil, errors.New("SignWithContext failed context")
	}

	if context[0] == crypto.DILITHIUM_ED25519_SPHINCS_FULL_ID {
		newDigestHash := crypto.Keccak256(digestHash, context)
		return s.Sign(newDigestHash, prv)
	}

	return nil, errors.New("SignWithContext failed invalid context")
}

func (s HybridedsfullSig) Verify(pubKey []byte, digestHash []byte, signature []byte) bool {
	sigBytes, pubKeyBytes, err := common.ExtractTwoParts(signature)
	if err != nil {
		return false
	}

	if !bytes.Equal(pubKey, pubKeyBytes) {
		return false
	}

	err = Verify(digestHash, sigBytes, pubKey)
	if err != nil {
		return false
	}

	//Important! Verify the original message
	for i := 0; i < len(digestHash); i++ {
		if sigBytes[2+CRYPTO_ED25519_SIGNATURE_BYTES+i] != digestHash[i] {
			return false
		}
	}

	return true
}

func (s HybridedsfullSig) VerifyWithContext(pubKey []byte, digestHash []byte, signature []byte, context []byte) bool {
	if context == nil || len(context) < 1 {
		return false
	}

	if context[0] == crypto.DILITHIUM_ED25519_SPHINCS_FULL_ID {
		newDigestHash := crypto.Keccak256(digestHash, context)
		return s.Verify(pubKey, newDigestHash, signature)
	}

	return false
}

func (s HybridedsfullSig) PublicKeyAndSignatureFromCombinedSignature(digestHash []byte, sig []byte) (signature []byte, pubKey []byte, err error) {
	signature, pubKey, err = common.ExtractTwoParts(sig)
	if err != nil {
		return nil, nil, err
	}

	err = Verify(digestHash, signature, pubKey)

	if err != nil {
		return nil, nil, err
	}

	return signature, pubKey, nil
}

func (s HybridedsfullSig) CombinePublicKeySignature(sigBytes []byte, pubKeyBytes []byte) (combinedSignature []byte, err error) {
	if len(sigBytes) < s.signatureLength {
		return nil, errors.New("invalid signature length")
	}

	if len(pubKeyBytes) != s.publicKeyLength {
		return nil, errors.New("invalid public key length")
	}

	return common.CombineTwoParts(sigBytes, pubKeyBytes), nil
}

func (s HybridedsfullSig) PublicKeyBytesFromSignature(digestHash []byte, sig []byte) ([]byte, error) {
	sigBytes, pubKeyBytes, err := common.ExtractTwoParts(sig)
	if err != nil {
		return nil, err
	}

	err = Verify(digestHash, sigBytes, pubKeyBytes)
	if err != nil {
		return nil, err
	}

	return pubKeyBytes, nil
}

func (s HybridedsfullSig) PublicKeyFromSignature(digestHash []byte, sig []byte) (*signaturealgorithm.PublicKey, error) {
	b, err := s.PublicKeyBytesFromSignature(digestHash, sig)
	if err != nil {
		return nil, err
	}
	return s.DeserializePublicKey(b)
}

func (s HybridedsfullSig) PublicKeyFromSignatureWithContext(digestHash []byte, sig []byte, context []byte) (*signaturealgorithm.PublicKey, error) {
	if context[0] != crypto.DILITHIUM_ED25519_SPHINCS_FULL_ID {
		return nil, errors.New("invalid context")
	}

	sigBytes, pubKeyBytes, err := common.ExtractTwoParts(sig)
	if err != nil {
		return nil, err
	}

	newDigestHash := crypto.Keccak256(digestHash, context)
	err = Verify(newDigestHash, sigBytes, pubKeyBytes)
	if err != nil {
		return nil, err
	}

	return s.DeserializePublicKey(pubKeyBytes)
}

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func (osig HybridedsfullSig) ValidateSignatureValues(digestHash []byte, v byte, r, s *big.Int) bool {
	if v == 0 || v == 1 {
		pubKey, signature := r.Bytes(), s.Bytes()

		if len(pubKey) != osig.PublicKeyLength() {
			return false
		}

		if len(signature) < osig.SignatureLength() {
			return false
		}

		combinedSignature := common.CombineTwoParts(signature, pubKey)
		if !osig.Verify(pubKey, digestHash, combinedSignature) {
			return false
		}

		return true
	}
	return false
}

func (s HybridedsfullSig) PublicKeyStartValue() byte {
	return 0x00 + 9
}

func (s HybridedsfullSig) SignatureStartValue() byte {
	return 0x30 + 9
}

func (s HybridedsfullSig) Zeroize(prv *signaturealgorithm.PrivateKey) {
	b := prv.PriData
	for i := range b {
		b[i] = 0
	}
}

func (s HybridedsfullSig) EncodePublicKey(pubKey *signaturealgorithm.PublicKey) []byte {
	encoded := make([]byte, s.publicKeyLength)
	copy(encoded, pubKey.PubData)
	return encoded
}

func (s HybridedsfullSig) DecodePublicKey(encoded []byte) (*signaturealgorithm.PublicKey, error) {
	if len(encoded) != s.publicKeyLength {
		return nil, errors.New("wrong size public key data")
	}
	p := &signaturealgorithm.PublicKey{}
	p.PubData = make([]byte, s.publicKeyLength)
	copy(p.PubData, encoded)
	return p, nil
}

// readASCII reads into 'buf', stopping when the buffer is full or
// when a non-printable control character is encountered.
func readASCII(buf []byte, r *bufio.Reader) (n int, err error) {
	for ; n < len(buf); n++ {
		buf[n], err = r.ReadByte()
		switch {
		case err == io.EOF || buf[n] < '!':
			return n, nil
		case err != nil:
			return n, err
		}
	}
	return n, nil
}

// checkKeyFileEnd skips over additional newlines at the end of a key file.
func checkKeyFileEnd(r *bufio.Reader) error {
	for i := 0; ; i++ {
		b, err := r.ReadByte()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case b != '\n' && b != '\r':
			return fmt.Errorf("invalid character %q at end of key file", b)
		case i >= 2:
			return errors.New("key file too long, want 64 hex characters")
		}
	}
}

// convertBytesToPrivate exports the corresponding secret key from the sig receiver.
func (s HybridedsfullSig) convertBytesToPrivate(privy []byte) (*signaturealgorithm.PrivateKey, error) {
	if len(privy) != s.privateKeyLength {
		return nil, ErrInvalidPrivateKeyLen
	}
	privKey := new(signaturealgorithm.PrivateKey)
	privKey.PriData = make([]byte, s.privateKeyLength)
	copy(privKey.PriData, privy)

	return privKey, nil
}

// convertBytesToPublic exports the corresponding secret key from the sig receiver.
func (s HybridedsfullSig) convertBytesToPublic(pub []byte) (*signaturealgorithm.PublicKey, error) {
	if len(pub) != s.publicKeyLength {
		return nil, ErrInvalidPublicKeyLen
	}
	pubKey := new(signaturealgorithm.PublicKey)
	pubKey.PubData = make([]byte, s.publicKeyLength)
	copy(pubKey.PubData, pub)
	return pubKey, nil
}

// exportPrivateKey exports a private key into a binary dump.
func (s HybridedsfullSig) exportPrivateKey(privy *signaturealgorithm.PrivateKey) ([]byte, error) {
	if len(privy.PriData) != s.privateKeyLength {
		return nil, ErrInvalidPrivateKeyLen
	}

	buf := make([]byte, s.privateKeyLength)
	copy(buf, privy.PriData)
	return buf, nil
}

// exportPublicKey exports a public key into a binary dump.
func (s HybridedsfullSig) exportPublicKey(pub *signaturealgorithm.PublicKey) ([]byte, error) {
	if len(pub.PubData) != s.publicKeyLength {
		return nil, ErrInvalidPublicKeyLen
	}
	buf := make([]byte, s.publicKeyLength)
	copy(buf, pub.PubData)
	return buf, nil
}
