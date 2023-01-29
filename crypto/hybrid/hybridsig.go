package hybrid

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
	"math/rand"
	"os"
)

type HybridSig struct {
	sigName                      string
	publicKeyBytesIndexStart     int
	publicKeyLength              int
	privateKeyLength             int
	signatureLength              int
	signatureWithPublicKeyLength int
}

func CreateHybridSig() HybridSig {
	return HybridSig{sigName: SIG_NAME,
		publicKeyBytesIndexStart:     12,
		publicKeyLength:              CRYPTO_PUBLICKEY_BYTES,
		privateKeyLength:             CRYPTO_SECRETKEY_BYTES,
		signatureLength:              CRYPTO_SIGNATURE_BYTES,
		signatureWithPublicKeyLength: CRYPTO_PUBLICKEY_BYTES + CRYPTO_SIGNATURE_BYTES + common.LengthByteSize + common.LengthByteSize,
	}
}

func (s HybridSig) SignatureName() string {
	return s.sigName
}

func (s HybridSig) PublicKeyLength() int {
	return s.publicKeyLength
}

func (s HybridSig) PrivateKeyLength() int {
	return s.privateKeyLength
}

func (s HybridSig) SignatureLength() int {
	return s.signatureLength
}

func (s HybridSig) SignatureWithPublicKeyLength() int {
	return s.signatureWithPublicKeyLength
}

func (s HybridSig) GenerateKey() (*signaturealgorithm.PrivateKey, error) {
	pubKey, priKey, err := GenerateKey()
	if err != nil {
		return nil, err
	}

	privy := new(signaturealgorithm.PrivateKey)
	privy.D = new(big.Int).SetBytes(priKey)
	privy.PublicKey.N = new(big.Int).SetBytes(pubKey)

	return privy, nil
}

func (s HybridSig) SerializePrivateKey(priv *signaturealgorithm.PrivateKey) ([]byte, error) {
	priBytes, err := s.exportPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	pubBytes, err := s.SerializePublicKey(&priv.PublicKey)
	if err != nil {
		return nil, err
	}

	return common.CombineTwoParts(priBytes, pubBytes), nil
}

func (s HybridSig) DeserializePrivateKey(priv []byte) (*signaturealgorithm.PrivateKey, error) {
	privKeyBytes, pubKeyBytes, err := common.ExtractTwoParts(priv)
	if err != nil {
		return nil, err
	}

	if s.doesPrivateMatchPublic(privKeyBytes, pubKeyBytes) == false {
		return nil, errors.New("publicKey does not match privateKey")
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

func (s HybridSig) doesPrivateMatchPublic(privKeyBytes []byte, pubKeyBytes []byte) bool {
	tempPrivBytes := make([]byte, len(privKeyBytes))
	copy(tempPrivBytes, privKeyBytes)

	digestHash := make([]byte, 32)
	rand.Read(digestHash)
	signature, err := Sign(tempPrivBytes, digestHash)
	if err != nil {
		return false
	}

	err = Verify(digestHash, signature, pubKeyBytes)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (s HybridSig) SerializePublicKey(pub *signaturealgorithm.PublicKey) ([]byte, error) {
	return s.exportPublicKey(pub)
}

func (s HybridSig) DeserializePublicKey(pub []byte) (*signaturealgorithm.PublicKey, error) {
	pubKey, error := s.convertBytesToPublic(pub)
	return pubKey, error
}

func (s HybridSig) HexToPrivateKey(hexkey string) (*signaturealgorithm.PrivateKey, error) {
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

func (s HybridSig) HexToPrivateKeyNoError(hexkey string) *signaturealgorithm.PrivateKey {
	p, err := s.HexToPrivateKey(hexkey)
	if err != nil {
		panic("HexToPrivateKey")
	}
	return p
}

func (s HybridSig) PrivateKeyToHex(priv *signaturealgorithm.PrivateKey) (string, error) {
	data, err := s.SerializePrivateKey(priv)
	if err != nil {
		return "", err
	}
	k := hex.EncodeToString(data)
	return k, nil
}

func (s HybridSig) PublicKeyToHex(pub *signaturealgorithm.PublicKey) (string, error) {
	data, err := s.SerializePublicKey(pub)
	if err != nil {
		return "", err
	}
	k := hex.EncodeToString(data)
	return k, nil
}

func (s HybridSig) HexToPublicKey(hexkey string) (*signaturealgorithm.PublicKey, error) {
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

func (s HybridSig) LoadPrivateKeyFromFile(file string) (*signaturealgorithm.PrivateKey, error) {
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

func (s HybridSig) SavePrivateKeyToFile(file string, key *signaturealgorithm.PrivateKey) error {
	k, err := s.PrivateKeyToHex(key)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(k), 0600)
}

func (s HybridSig) PublicKeyToAddress(p *signaturealgorithm.PublicKey) (common.Address, error) {
	pubBytes, err := s.SerializePublicKey(p)
	tempAddr := common.Address{}
	if err != nil {
		return tempAddr, err
	}
	return common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[s.publicKeyBytesIndexStart:]), nil
}

func (s HybridSig) PublicKeyToAddressNoError(p *signaturealgorithm.PublicKey) common.Address {
	addr, err := s.PublicKeyToAddress(p)
	if err != nil {
		panic("PublicKeyToAddress failed")
	}
	return addr
}

func (s HybridSig) Sign(digestHash []byte, prv *signaturealgorithm.PrivateKey) (sig []byte, err error) {
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

	return common.CombineTwoParts(sigBytes, pubBytes), nil
}

func (s HybridSig) Verify(pubKey []byte, digestHash []byte, signature []byte) bool {
	sigBytes, pubKeyBytes, err := common.ExtractTwoParts(signature)
	if err != nil {
		return false
	}

	if !bytes.Equal(pubKey, pubKeyBytes) {
		return false
	}

	err = Verify(digestHash, sigBytes, pubKey)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (s HybridSig) PublicKeyAndSignatureFromCombinedSignature(digestHash []byte, sig []byte) (signature []byte, pubKey []byte, err error) {
	signature, pubKey, err = common.ExtractTwoParts(sig)
	if err != nil {
		return nil, nil, err
	}

	if digestHash != nil {
		err = Verify(digestHash, signature, pubKey)

		if err != nil {
			return nil, nil, err
		}
	}

	return signature, pubKey, nil
}

func (s HybridSig) CombinePublicKeySignature(sigBytes []byte, pubKeyBytes []byte) (combinedSignature []byte, err error) {
	if len(sigBytes) < s.signatureLength {
		return nil, errors.New("invalid signature length")
	}

	if len(pubKeyBytes) != s.publicKeyLength {
		return nil, errors.New("invalid public key length")
	}

	return common.CombineTwoParts(sigBytes, pubKeyBytes), nil
}

func (s HybridSig) PublicKeyBytesFromSignature(digestHash []byte, sig []byte) ([]byte, error) {
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

func (s HybridSig) PublicKeyFromSignature(digestHash []byte, sig []byte) (*signaturealgorithm.PublicKey, error) {
	b, err := s.PublicKeyBytesFromSignature(digestHash, sig)
	if err != nil {
		return nil, err
	}
	return s.DeserializePublicKey(b)
}

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func (osig HybridSig) ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	if v == 0 || v == 1 {
		// encode the signature in uncompressed format
		R, S := r.Bytes(), s.Bytes()

		if len(R) != osig.PublicKeyLength() {
			return false
		}

		if len(S) < osig.SignatureLength() {
			return false
		}

		return true
	}
	return false
}

func (s HybridSig) PublicKeyStartValue() byte {
	return 0x00 + 9
}

func (s HybridSig) SignatureStartValue() byte {
	return 0x30 + 9
}

func (s HybridSig) Zeroize(prv *signaturealgorithm.PrivateKey) {
	b := prv.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

func (s HybridSig) PrivateKeyAsBigInt(prv *signaturealgorithm.PrivateKey) *big.Int {
	privKeyBytes, err := s.SerializePrivateKey(prv)
	if err != nil {
		panic(err) //todo: no panic
	}

	return new(big.Int).SetBytes(privKeyBytes)
}

func (s HybridSig) PublicKeyAsBigInt(pub *signaturealgorithm.PublicKey) *big.Int {
	return pub.N
}

func (s HybridSig) EncodePublicKey(pubKey *signaturealgorithm.PublicKey) []byte {
	encoded := make([]byte, s.publicKeyLength)
	math.ReadBits(s.PublicKeyAsBigInt(pubKey), encoded[:])
	return encoded
}

func (s HybridSig) DecodePublicKey(encoded []byte) (*signaturealgorithm.PublicKey, error) {
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

// convertBytesToPrivate exports the corresponding secret key from the sig receiver.
func (s HybridSig) convertBytesToPrivate(privy []byte) (*signaturealgorithm.PrivateKey, error) {
	if len(privy) != int(s.privateKeyLength) {
		return nil, ErrInvalidPrivateKeyLen
	}
	privKey := new(signaturealgorithm.PrivateKey)
	privKey.D = new(big.Int).SetBytes(privy)
	return privKey, nil
}

// convertBytesToPublic exports the corresponding secret key from the sig receiver.
func (s HybridSig) convertBytesToPublic(pub []byte) (*signaturealgorithm.PublicKey, error) {
	if len(pub) != int(s.publicKeyLength) {
		fmt.Println("convertBytesToPublic", len(pub), s.publicKeyLength)
		return nil, ErrInvalidPublicKeyLen
	}
	pubKey := new(signaturealgorithm.PublicKey)
	pubKey.N = new(big.Int).SetBytes(pub)
	return pubKey, nil
}

// exportPrivateKey exports a private key into a binary dump.
func (s HybridSig) exportPrivateKey(privy *signaturealgorithm.PrivateKey) ([]byte, error) {
	if len(privy.D.Bytes()) != int(s.privateKeyLength) {
		return nil, ErrInvalidPrivateKeyLen
	}
	return privy.D.Bytes(), nil
}

// exportPublicKey exports a public key into a binary dump.
func (s HybridSig) exportPublicKey(pub *signaturealgorithm.PublicKey) ([]byte, error) {
	if len(pub.N.Bytes()) != int(s.publicKeyLength) {
		fmt.Println("exportPublicKey", len(pub.N.Bytes()), s.publicKeyLength)
		return nil, ErrInvalidPublicKeyLen
	}
	return pub.N.Bytes(), nil
}
