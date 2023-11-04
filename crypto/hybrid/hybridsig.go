package hybrid

import (
	"bufio"
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/falcon"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"io"
	"io/ioutil"
	"math/big"
	"os"
)

const CRYPTO_ED25519_PUBLICKEY_BYTES = 32
const CRYPTO_ED25519_SIGNATURE_BYTES = 64
const LEN_BYTES = 2
const CRYPTO_FALCON_PUBLICKEY_BYTES = 897
const CRYPTO_FALCON_SECRETKEY_BYTES = 1281
const CRYPTO_FALCON_SECRETKEY_WITH_PUBLIC_KEY_BYTES = 1281 + 897
const CRYPTO_FALCON_NONCE_LENGTH = 40
const CRYPTO_HYBRID_MIN_SIGNATURE_BYTES = 64 + 600 + 40 + 2
const CRYPTO_HYBRID_MAX_SIGNATURE_BYTES = 64 + 690 + 40 + 2
const CRYPTO_FALCON_MIN_SIGNATURE_BYTES = 600 + 40 + 2 //Signature + Nonce + 2 for size
const CRYPTO_FALCON_MAX_SIGNATURE_BYTES = 690 + 40 + 2 //Signature + Nonce + 2 for size

type HybridSig struct {
	sigName                      string
	publicKeyLength              int
	privateKeyLength             int
	signatureLength              int
	signatureWithPublicKeyLength int
	NativeGolangVerify           bool
}

func CreateHybridSig(mativeGolangVerify bool) HybridSig {
	return HybridSig{sigName: SIG_NAME,
		publicKeyLength:              CRYPTO_PUBLICKEY_BYTES,
		privateKeyLength:             CRYPTO_SECRETKEY_BYTES,
		signatureLength:              CRYPTO_SIGNATURE_BYTES,
		signatureWithPublicKeyLength: CRYPTO_PUBLICKEY_BYTES + CRYPTO_SIGNATURE_BYTES + common.LengthByteSize + common.LengthByteSize,
		NativeGolangVerify:           mativeGolangVerify,
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

func (s HybridSig) SerializePrivateKey(priv *signaturealgorithm.PrivateKey) ([]byte, error) {
	priBytes, err := s.exportPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	return priBytes, err
}

func (s HybridSig) DeserializePrivateKey(priv []byte) (*signaturealgorithm.PrivateKey, error) {

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
	addr := common.BytesToAddress(crypto.Keccak256(pubBytes[:])[:])
	return addr, nil
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

	combinedSignature := common.CombineTwoParts(sigBytes, pubBytes)

	if !s.Verify(pubBytes, digestHash, combinedSignature) {
		return nil, errors.New("Verify failed after signing")
	}

	return combinedSignature, nil
}

func (s HybridSig) Verify(pubKey []byte, digestHash []byte, signature []byte) bool {
	if s.NativeGolangVerify {
		return s.VerifyNative(pubKey, digestHash, signature)
	}

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

// Verify with GoLang's native ED25519 implementation, while using hybrid-pqc for Falcon verification
func (s HybridSig) VerifyNative(pubKey []byte, digestHash []byte, signature []byte) bool {
	sigBytes, pubKeyBytes, err := common.ExtractTwoParts(signature)
	if err != nil {
		return false
	}

	if !bytes.Equal(pubKey, pubKeyBytes) {
		return false
	}

	msgLen := binary.BigEndian.Uint16(sigBytes[LEN_BYTES : LEN_BYTES+LEN_BYTES])
	if msgLen != uint16(len(digestHash)) {
		return false
	}

	ed25519Signature := sigBytes[LEN_BYTES+LEN_BYTES : LEN_BYTES+LEN_BYTES+CRYPTO_ED25519_SIGNATURE_BYTES]
	ed25519PubKey := pubKey[:CRYPTO_ED25519_PUBLICKEY_BYTES]

	ok := ed25519.Verify(ed25519PubKey, digestHash, ed25519Signature)
	if ok == false {
		return false
	}

	totalLen := binary.BigEndian.Uint16(sigBytes[:])
	if totalLen < CRYPTO_HYBRID_MIN_SIGNATURE_BYTES+msgLen || totalLen > CRYPTO_HYBRID_MAX_SIGNATURE_BYTES+msgLen {
		return false
	}

	sig1Len := CRYPTO_ED25519_SIGNATURE_BYTES + msgLen
	sig2Len := totalLen - sig1Len
	if sig2Len < CRYPTO_FALCON_MIN_SIGNATURE_BYTES || sig2Len > CRYPTO_FALCON_MAX_SIGNATURE_BYTES {
		return false
	}

	actualSig2Len := sig2Len - LEN_BYTES - CRYPTO_FALCON_NONCE_LENGTH - msgLen
	var sig2 [2 + 40 + 64 + 690]byte //SIZE_LEN + CRYPTO_FALCON_NONCE_LENGTH + MAX_MSG_LEN + CRYPTO_FALCON_MAX_SIGNATURE_BYTES
	binary.BigEndian.PutUint16(sig2[:], actualSig2Len)

	//Copy Falcon nonce into falconSig
	for i := uint16(0); i < CRYPTO_FALCON_NONCE_LENGTH; i++ {
		sig2[LEN_BYTES+i] = sigBytes[LEN_BYTES+LEN_BYTES+sig1Len+i]
	}

	//Copy Message info falconSig
	for i := uint16(0); i < msgLen; i++ {
		sig2[LEN_BYTES+CRYPTO_FALCON_NONCE_LENGTH+i] = sigBytes[LEN_BYTES+LEN_BYTES+CRYPTO_ED25519_SIGNATURE_BYTES+i]
	}

	//Copy actual Sig2 from source
	for i := uint16(0); i < sig2Len-LEN_BYTES-CRYPTO_FALCON_NONCE_LENGTH-msgLen; i++ {
		sig2[LEN_BYTES+CRYPTO_FALCON_NONCE_LENGTH+msgLen+i] = sigBytes[LEN_BYTES+LEN_BYTES+sig1Len+CRYPTO_FALCON_NONCE_LENGTH+i]
	}

	falconPubKey := pubKeyBytes[CRYPTO_ED25519_PUBLICKEY_BYTES:]

	err = falcon.VerifyDirect(digestHash, sig2[:], falconPubKey, sig2Len)
	if err != nil {
		return false
	}

	return true
}

func (s HybridSig) PublicKeyAndSignatureFromCombinedSignature(digestHash []byte, sig []byte) (signature []byte, pubKey []byte, err error) {
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
func (osig HybridSig) ValidateSignatureValues(digestHash []byte, v byte, r, s *big.Int) bool {
	if v == 0 || v == 1 {
		pubKey, signature := r.Bytes(), s.Bytes()

		if len(pubKey) != osig.PublicKeyLength() {
			fmt.Println("ValidateSignatureValues 1", len(pubKey), osig.PublicKeyLength())
			return false
		}

		if len(signature) < osig.SignatureLength() {
			fmt.Println("ValidateSignatureValues 2")
			return false
		}

		combinedSignature := common.CombineTwoParts(signature, pubKey)
		if !osig.Verify(pubKey, digestHash, combinedSignature) {
			fmt.Println("ValidateSignatureValues 3")
			return false
		}

		return true
	}
	fmt.Println("ValidateSignatureValues 4")
	return false
}

func (s HybridSig) PublicKeyStartValue() byte {
	return 0x00 + 9
}

func (s HybridSig) SignatureStartValue() byte {
	return 0x30 + 9
}

func (s HybridSig) Zeroize(prv *signaturealgorithm.PrivateKey) {
	b := prv.PriData
	for i := range b {
		b[i] = 0
	}
}

func (s HybridSig) EncodePublicKey(pubKey *signaturealgorithm.PublicKey) []byte {
	encoded := make([]byte, s.publicKeyLength)
	copy(encoded, pubKey.PubData)
	return encoded
}

func (s HybridSig) DecodePublicKey(encoded []byte) (*signaturealgorithm.PublicKey, error) {
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
func (s HybridSig) convertBytesToPrivate(privy []byte) (*signaturealgorithm.PrivateKey, error) {
	if len(privy) != s.privateKeyLength {
		return nil, ErrInvalidPrivateKeyLen
	}
	privKey := new(signaturealgorithm.PrivateKey)
	privKey.PriData = make([]byte, s.privateKeyLength)
	copy(privKey.PriData, privy)

	return privKey, nil
}

// convertBytesToPublic exports the corresponding secret key from the sig receiver.
func (s HybridSig) convertBytesToPublic(pub []byte) (*signaturealgorithm.PublicKey, error) {
	if len(pub) != s.publicKeyLength {
		return nil, ErrInvalidPublicKeyLen
	}
	pubKey := new(signaturealgorithm.PublicKey)
	pubKey.PubData = make([]byte, s.publicKeyLength)
	copy(pubKey.PubData, pub)
	return pubKey, nil
}

// exportPrivateKey exports a private key into a binary dump.
func (s HybridSig) exportPrivateKey(privy *signaturealgorithm.PrivateKey) ([]byte, error) {
	if len(privy.PriData) != s.privateKeyLength {
		return nil, ErrInvalidPrivateKeyLen
	}

	buf := make([]byte, s.privateKeyLength)
	copy(buf, privy.PriData)
	return buf, nil
}

// exportPublicKey exports a public key into a binary dump.
func (s HybridSig) exportPublicKey(pub *signaturealgorithm.PublicKey) ([]byte, error) {
	if len(pub.PubData) != s.publicKeyLength {
		return nil, ErrInvalidPublicKeyLen
	}
	buf := make([]byte, s.publicKeyLength)
	copy(buf, pub.PubData)
	return buf, nil
}
