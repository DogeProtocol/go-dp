package signaturealgorithm

import (
	"bytes"
	"fmt"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/crypto"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var (
	testmsg1 = hexutil.MustDecode("0x68692074686572656f636b636861696e62626262626262626262626262626262")
	testmsg2 = hexutil.MustDecode("0x68692074686572656f636b636861696e62626262626262626262626262626263")
)

func SignatureAlgorithmTest(t *testing.T, sig SignatureAlgorithm) {
	sig.SignatureName()
	sig.PublicKeyLength()
	sig.PrivateKeyLength()
	sig.SignatureLength()
	sig.SignatureWithPublicKeyLength()
	sig.PublicKeyStartValue()
	sig.SignatureStartValue()

	key1, err := sig.GenerateKey()
	if err != nil {
		t.Fatal("GenerateKey failed")
	}

	priBytes1, err := sig.SerializePrivateKey(key1)
	if err != nil {
		t.Fatal("SerializePrivateKey failed")
	}

	//Temp copy array, since DeserializePrivateKey will clean private-key
	priBytes1Temp := make([]byte, len(priBytes1))
	copy(priBytes1Temp, priBytes1)

	key2, err := sig.DeserializePrivateKey(priBytes1)
	if err != nil {
		t.Fatal("DeserializePrivateKey failed")
	}

	priBytes2, err := sig.SerializePrivateKey(key2)
	if err != nil {
		t.Fatal("SerializePrivateKey failed")
	}

	if bytes.Compare(priBytes1Temp, priBytes2) != 0 {
		t.Fatal("Issue in serialize / deserialize privateKey")
	}

	addr1, err := sig.PublicKeyToAddress(&key1.PublicKey)
	if err != nil {
		t.Fatal("SerializePublicKey failed")
	}

	pubBytes1, err := sig.SerializePublicKey(&key1.PublicKey)
	if err != nil {
		t.Fatal("SerializePublicKey failed")
	}

	pubKey1, err := sig.DeserializePublicKey(pubBytes1)
	if err != nil {
		t.Fatal("DeserializePublicKey failed")
	}

	addr2, err := sig.PublicKeyToAddress(pubKey1)
	if err != nil {
		t.Fatal("PublicKeyBytesToAddress failed")
	}

	if addr1 != addr2 {
		t.Fatal("address mismatch")
	}

	pubKeyDirect1 := PublicKey{PubData: pubBytes1}
	addr3, err := sig.PublicKeyToAddress(&pubKeyDirect1)
	if err != nil {
		t.Fatal("PublicKeyBytesToAddress failed")
	}

	if addr1 != addr3 {
		t.Fatal("address mismatch")
	}

	addr4, err := sig.PublicKeyToAddress(&key2.PublicKey)
	if err != nil {
		t.Fatal("PublicKeyBytesToAddress failed")
	}

	if addr1 != addr4 {
		t.Fatal("address mismatch")
	}

	pubKeyDirect2 := PublicKey{PubData: key2.PubData}
	addr5, err := sig.PublicKeyToAddress(&pubKeyDirect2)
	if err != nil {
		t.Fatal("PublicKeyBytesToAddress failed")
	}

	if addr1 != addr5 {
		t.Fatal("address mismatch")
	}

	pubBytes2, err := sig.SerializePublicKey(pubKey1)
	if err != nil {
		t.Fatal("SerializePublicKey failed")
	}

	if bytes.Compare(pubBytes1, pubBytes2) != 0 {
		t.Fatal("Issue in serialize / deserialize publicKey")
	}

	digestHash1 := []byte(testmsg1)
	signature1, err := sig.Sign(digestHash1, key1)
	if err != nil {
		t.Fatal("Sign failed")
	}

	if sig.Verify(pubBytes1, digestHash1, signature1) != true {
		t.Fatal("Verify failed")
	}

	signature1copy, err := sig.Sign(digestHash1, key1)
	if err != nil {
		t.Fatal("Sign failed")
	}
	if bytes.Compare(signature1, signature1copy) != 0 {
		fmt.Errorf("signature not deterministic")
	}

	digestHash2 := []byte(testmsg2)
	signature2, err := sig.Sign(digestHash2, key1)
	if err != nil {
		t.Fatal("Sign failed")
	}

	if sig.Verify(pubBytes1, digestHash2, signature2) != true {
		t.Fatal("Verify failed")
	}

	if sig.Verify(pubBytes1, digestHash1, signature2) != false {
		t.Fatal("Verify negative failed")
	}

	if sig.Verify(pubBytes1, digestHash2, signature1) != false {
		t.Fatal("Verify negative failed")
	}

	//Deep signature change test
	maxFalconSigSize := 600 //todo: in hybrid-pqc lib, add zero check for padded signatures
	for i := 0; i < maxFalconSigSize; i++ {
		sigTemp := make([]byte, len(signature2))
		copy(sigTemp, signature2)
		sigTemp[i] = sigTemp[i] + 1
		if sig.Verify(pubBytes1, digestHash2, sigTemp) != false {
			t.Fatal("Verify signature change negative failed", i)
		}
	}

	if sig.Verify(pubBytes1, digestHash2, signature2) != true {
		t.Fatal("Verify failed")
	}

	//Deep public key change test
	for i := 0; i < len(pubBytes1); i++ {
		pubTemp := make([]byte, len(pubBytes1))
		copy(pubTemp, pubBytes1)
		pubTemp[i] = pubTemp[i] + 1
		if sig.Verify(pubTemp, digestHash2, signature2) != false {
			t.Fatal("Verify signature change negative failed")
		}
	}

	if sig.Verify(pubBytes1, digestHash2, signature2) != true {
		t.Fatal("Verify failed")
	}

	sigExtracted, pubExtracted, err := sig.PublicKeyAndSignatureFromCombinedSignature(digestHash1, signature1)
	if err != nil {
		t.Fatal(err)
	}

	if len(sigExtracted) != sig.SignatureLength() {
		t.Fatal("invalid signature length")
	}

	if len(pubExtracted) != sig.PublicKeyLength() {
		t.Fatal("invalid public key length")
	}

	combinedSig, err := sig.CombinePublicKeySignature(sigExtracted, pubExtracted)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(signature1, combinedSig) != 0 {
		t.Fatal("invalid combined sig")
	}

	hex1, err := sig.PrivateKeyToHex(key1)
	if err != nil {
		t.Fatal("PrivateKeyToHex failed")
	}

	_, err = sig.PublicKeyToAddress(pubKey1)
	if err != nil {
		t.Fatal(err)
	}

	key3, err := sig.HexToPrivateKey(hex1)
	if err != nil {
		t.Fatal("HexToPrivateKey failed")
	}

	key31 := sig.HexToPrivateKeyNoError(hex1)
	if err != nil {
		t.Fatal("HexToPrivateKeyNoError failed")
	}

	if bytes.Compare(key3.PriData, key31.PriData) != 0 {
		t.Fatal("private key compare failed")
	}

	hex2, err := sig.PrivateKeyToHex(key3)
	if err != nil {
		t.Fatal("PrivateKeyToHex failed")
	}

	if strings.Compare(hex1, hex2) != 0 {
		t.Fatal("Hex compare failed")
	}

	f, err := ioutil.TempFile("", "saveOQS_test.*.txt")
	if err != nil {
		t.Fatal(err)
	}
	file := f.Name()
	f.Close()
	defer os.Remove(file)
	err = sig.SavePrivateKeyToFile(file, key1)
	if err != nil {
		t.Fatal(err)
	}

	key4, err := sig.LoadPrivateKeyFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	hex4, err := sig.PrivateKeyToHex(key4)
	if err != nil {
		t.Fatal("PrivateKeyToHex failed")
	}

	if strings.Compare(hex1, hex4) != 0 {
		t.Fatal("Hex compare failed")
	}

	pub1Key, err := sig.PublicKeyFromSignature(digestHash1, signature1)
	if err != nil {
		t.Fatal(err)
	}

	pub1BytesTemp1, err := sig.SerializePublicKey(pub1Key)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(pubBytes1, pub1BytesTemp1) != 0 {
		t.Fatal(err)
	}

	pub1BytesTemp2, err := sig.PublicKeyBytesFromSignature(digestHash1, signature1)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(pubBytes1, pub1BytesTemp2) != 0 {
		t.Fatal(err)
	}

	generatedAddress, err := sig.PublicKeyToAddress(pub1Key)
	if err != nil {
		t.Fatal(err)
	}

	generatedAddress2 := sig.PublicKeyToAddressNoError(pub1Key)
	if bytes.Compare(generatedAddress[:], generatedAddress2[:]) != 0 {
		t.Fatal("address compare failed")
	}

	addr := crypto.PublicKeyBytesToAddress(pubBytes1)

	if generatedAddress != addr {
		t.Fatal(err)
	}

	hexpub, err := sig.PublicKeyToHex(pub1Key)
	if err != nil {
		t.Fatal(err)
	}

	pubBackFromHex, err := sig.HexToPublicKey(hexpub)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(pub1Key.PubData, pubBackFromHex.PubData) != 0 {
		t.Fatal("public key compare failed")
	}

	encodedPubKey := sig.EncodePublicKey(pub1Key)
	pubkeyDecoded, err := sig.DecodePublicKey(encodedPubKey)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(pubkeyDecoded.PubData, pub1Key.PubData) != 0 {
		t.Fatal(err)
	}

	serializedKey1Bytes, err := sig.SerializePrivateKey(key1)
	if err != nil {
		t.Fatal(err)
	}
	priDeserialized, err := sig.DeserializePrivateKey(serializedKey1Bytes)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(key1.PriData, priDeserialized.PriData) != 0 {
		t.Fatal("pri data compare failed 2")
	}

	if bytes.Compare(key1.PublicKey.PubData, priDeserialized.PublicKey.PubData) != 0 {
		t.Fatal("pub data compare failed")
	}

	sig.Zeroize(key1)
}
