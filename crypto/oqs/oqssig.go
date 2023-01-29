package oqs

import "C"
import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/signaturealgorithm"
	"io"
	"io/ioutil"
	"math/big"
	"os"
)

type OqsSig struct {
	sigName                      string
	publicKeyBytesIndexStart     int
	publicKeyLength              int
	privateKeyLength             int
	signatureLength              int
	signatureWithPublicKeyLength int
}

func CreateOqs(sigName string) OqsSig {
	sigDetails, err := GetSignatureDetails(sigName)
	if err != nil {
		panic("unable to initialize " + sigName)
	}

	return OqsSig{sigName: sigName,
		publicKeyBytesIndexStart:     12,
		publicKeyLength:              sigDetails.LengthPublicKey,
		privateKeyLength:             sigDetails.LengthSecretKey,
		signatureLength:              sigDetails.MaxLengthSignature,
		signatureWithPublicKeyLength: sigDetails.LengthPublicKey + sigDetails.MaxLengthSignature + common.LengthByteSize + common.LengthByteSize,
	}
}

func (s OqsSig) SignatureName() string {
	return s.sigName
}

func (s OqsSig) PublicKeyLength() int {
	return s.publicKeyLength
}

func (s OqsSig) PrivateKeyLength() int {
	return s.privateKeyLength
}

func (s OqsSig) SignatureLength() int {
	return s.signatureLength
}

func (s OqsSig) SignatureWithPublicKeyLength() int {
	return s.signatureWithPublicKeyLength
}

func (s OqsSig) GenerateKey() (*signaturealgorithm.PrivateKey, error) {
	return GenerateKey(s.sigName)
}

func (s OqsSig) SerializePrivateKey(priv *signaturealgorithm.PrivateKey) ([]byte, error) {
	priBytes, err := ExportPrivateKey(s.sigName, priv)
	if err != nil {
		return nil, err
	}

	pubBytes, err := s.SerializePublicKey(&priv.PublicKey)
	if err != nil {
		return nil, err
	}

	return common.CombineTwoParts(priBytes, pubBytes), nil
}

func (s OqsSig) DeserializePrivateKey(priv []byte) (*signaturealgorithm.PrivateKey, error) {
	privKeyBytes, pubKeyBytes, err := common.ExtractTwoParts(priv)
	if err != nil {
		return nil, err
	}

	if s.doesPrivateMatchPublic(privKeyBytes, pubKeyBytes) == false {
		return nil, errors.New("publicKey does not match privateKey")
	}

	privKey, err := ConvertBytesToPrivate(s.sigName, privKeyBytes)
	if err != nil {
		return nil, err
	}

	pubkey, err := ConvertBytesToPublic(s.sigName, pubKeyBytes)
	if err != nil {
		return nil, err
	}

	privKey.PublicKey = *pubkey
	//get private key
	return privKey, err
}

func (s OqsSig) doesPrivateMatchPublic(privKeyBytes []byte, pubKeyBytes []byte) bool {
	tempPrivBytes := make([]byte, len(privKeyBytes))
	copy(tempPrivBytes, privKeyBytes)

	digestHash := []byte("verify digest hash")
	signature, err := Sign(s.sigName, digestHash, tempPrivBytes)
	if err != nil {
		return false
	}

	return VerifySignature(s.sigName, pubKeyBytes, digestHash, signature)
}

func (s OqsSig) SerializePublicKey(pub *signaturealgorithm.PublicKey) ([]byte, error) {
	return ExportPublicKey(s.sigName, pub)
}

func (s OqsSig) DeserializePublicKey(pub []byte) (*signaturealgorithm.PublicKey, error) {
	pubKey, error := ConvertBytesToPublic(s.sigName, pub)
	return pubKey, error
}

func (s OqsSig) HexToPrivateKey(hexkey string) (*signaturealgorithm.PrivateKey, error) {
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

func (s OqsSig) HexToPrivateKeyNoError(hexkey string) *signaturealgorithm.PrivateKey {
	p, err := s.HexToPrivateKey(hexkey)
	if err != nil {
		panic("HexToPrivateKey")
	}
	return p
}

func (s OqsSig) PrivateKeyToHex(priv *signaturealgorithm.PrivateKey) (string, error) {
	data, err := s.SerializePrivateKey(priv)
	if err != nil {
		return "", err
	}
	k := hex.EncodeToString(data)
	return k, nil
}

func (s OqsSig) PublicKeyToHex(pub *signaturealgorithm.PublicKey) (string, error) {
	data, err := s.SerializePublicKey(pub)
	if err != nil {
		return "", err
	}
	k := hex.EncodeToString(data)
	return k, nil
}

func (s OqsSig) HexToPublicKey(hexkey string) (*signaturealgorithm.PublicKey, error) {
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

func (s OqsSig) LoadPrivateKeyFromFile(file string) (*signaturealgorithm.PrivateKey, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	r := bufio.NewReader(fd)
	buf := make([]byte, (s.privateKeyLength+s.publicKeyLength+common.LengthByteSize+common.LengthByteSize)*2)
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

func (s OqsSig) SavePrivateKeyToFile(file string, key *signaturealgorithm.PrivateKey) error {
	k, err := s.PrivateKeyToHex(key)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(k), 0600)
}

func (s OqsSig) PublicKeyToAddress(p *signaturealgorithm.PublicKey) (common.Address, error) {
	pubBytes, err := s.SerializePublicKey(p)
	tempAddr := common.Address{}
	if err != nil {
		return tempAddr, err
	}
	return common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[s.publicKeyBytesIndexStart:]), nil
}

func (s OqsSig) PublicKeyToAddressNoError(p *signaturealgorithm.PublicKey) common.Address {
	addr, err := s.PublicKeyToAddress(p)
	if err != nil {
		panic("PublicKeyToAddress failed")
	}
	return addr
}

func (s OqsSig) Sign(digestHash []byte, prv *signaturealgorithm.PrivateKey) (sig []byte, err error) {
	seckey, err := ExportPrivateKey(s.sigName, prv)
	if err != nil {
		return nil, err
	}

	sigBytes, err := Sign(s.sigName, digestHash, seckey)
	if err != nil {
		return nil, err
	}

	pubBytes, err := s.SerializePublicKey(&prv.PublicKey)
	if err != nil {
		return nil, err
	}

	return common.CombineTwoParts(sigBytes, pubBytes), nil
}

func (s OqsSig) Verify(pubKey []byte, digestHash []byte, signature []byte) bool {
	sigBytes, pubKeyBytes, err := common.ExtractTwoParts(signature)
	if err != nil {
		return false
	}

	if !bytes.Equal(pubKey, pubKeyBytes) {
		return false
	}

	return VerifySignature(s.sigName, pubKey, digestHash, sigBytes)
}

func (s OqsSig) PublicKeyAndSignatureFromCombinedSignature(digestHash []byte, sig []byte) (signature []byte, pubKey []byte, err error) {
	signature, pubKey, err = common.ExtractTwoParts(sig)
	if err != nil {
		return nil, nil, err
	}

	if digestHash != nil {
		if VerifySignature(s.sigName, pubKey, digestHash, signature) == false {
			return nil, nil, errors.New("verify failed")
		}
	}

	return signature, pubKey, nil
}

func (s OqsSig) CombinePublicKeySignature(sigBytes []byte, pubKeyBytes []byte) (combinedSignature []byte, err error) {
	if len(sigBytes) != s.signatureLength {
		return nil, errors.New("invalid signature length")
	}

	if len(pubKeyBytes) != s.publicKeyLength {
		return nil, errors.New("invalid public key length")
	}

	return common.CombineTwoParts(sigBytes, pubKeyBytes), nil
}

func (s OqsSig) PublicKeyBytesFromSignature(digestHash []byte, sig []byte) ([]byte, error) {
	sigBytes, pubKeyBytes, err := common.ExtractTwoParts(sig)
	if err != nil {
		return nil, err
	}

	if VerifySignature(s.sigName, pubKeyBytes, digestHash, sigBytes) == false {
		return nil, errors.New("verify failed")
	}

	return pubKeyBytes, nil
}

func (s OqsSig) PublicKeyFromSignature(digestHash []byte, sig []byte) (*signaturealgorithm.PublicKey, error) {
	b, err := s.PublicKeyBytesFromSignature(digestHash, sig)
	if err != nil {
		return nil, err
	}
	return s.DeserializePublicKey(b)
}

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func (osig OqsSig) ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	if v == 0 || v == 1 {
		// encode the signature in uncompressed format
		R, S := r.Bytes(), s.Bytes()

		if len(R) != osig.PublicKeyLength() {
			return false
		}

		if len(S) != osig.SignatureLength() {
			return false
		}

		return true
	}
	return false
}

func (s OqsSig) PublicKeyStartValue() byte {
	return 0x00 + 9
}

func (s OqsSig) SignatureStartValue() byte {
	return 0x30 + 9
}

func (s OqsSig) Zeroize(prv *signaturealgorithm.PrivateKey) {
	b := prv.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

func (s OqsSig) PrivateKeyAsBigInt(prv *signaturealgorithm.PrivateKey) *big.Int {
	privKeyBytes, err := s.SerializePrivateKey(prv)
	if err != nil {
		panic(err) //todo: no panic
	}

	return new(big.Int).SetBytes(privKeyBytes)
}

func (s OqsSig) PublicKeyAsBigInt(pub *signaturealgorithm.PublicKey) *big.Int {
	return pub.N
}

func (s OqsSig) EncodePublicKey(pubKey *signaturealgorithm.PublicKey) []byte {
	encoded := make([]byte, s.publicKeyLength)
	math.ReadBits(s.PublicKeyAsBigInt(pubKey), encoded[:])
	return encoded
}

func (s OqsSig) DecodePublicKey(encoded []byte) (*signaturealgorithm.PublicKey, error) {
	if len(encoded) != s.publicKeyLength {
		return nil, errors.New("wrong size public key data")
	}
	p := &signaturealgorithm.PublicKey{N: new(big.Int)}
	p.N.SetBytes(encoded[:])

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
