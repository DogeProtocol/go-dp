package rlpx

import (
	"bytes"
	crypto2 "crypto"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/DogeProtocol/dp/common"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"
)

const (
	derivedLabelName                = "derived"
	clientHandshakeTrafficLabelName = "c hs traffic"
	serverHandshakeTrafficLabelName = "s hs traffic"
	secretKeyLabelName              = "key"
	secretIvLabelName               = "iv"

	clientApplicationTrafficLabelName = "c ap traffic"
	serverApplicationTrafficLabelName = "s ap traffic"
)

type SessionSecret struct {
	handshakeSecret []byte

	clientHandshakeTrafficSecret []byte
	serverHandshakeTrafficSecret []byte
	ClientHandshakeKey           []byte
	ServerHandshakeKey           []byte
	ClientHandshakeIv            []byte
	ServerHandshakeIv            []byte

	clientApplicationTrafficSecret []byte
	serverApplicationTrafficSecret []byte
	ClientApplicationKey           []byte
	ServerApplicationKey           []byte
	ClientApplicationIv            []byte
	ServerApplicationIv            []byte

	ClientHandshakeCipher cipher.AEAD
	ServerHandshakeCipher cipher.AEAD

	ClientApplicationCipher cipher.AEAD
	ServerApplicationCipher cipher.AEAD

	masterSecret   []byte
	TranscriptHash []byte
}

func NewSessionSecret(transcriptHash []byte, sharedSecret []byte) (*SessionSecret, error) {
	//Create early secrets
	zeroKey := bytes.Repeat([]byte{0}, common.HashLength)
	earlySecret := hkdf.Extract(sha3.New256, zeroKey, transcriptHash)

	var hash crypto2.Hash
	hash = crypto2.SHA3_256
	emptyHash := hash.New().Sum(nil)

	derivedSecret, err := HkdfExpandLabel(
		earlySecret,
		derivedLabelName,
		emptyHash,
		shaLength)
	if err != nil {
		return nil, err
	}

	handshakeSecret := hkdf.Extract(sha3.New256, sharedSecret, derivedSecret)

	clientHandshakeTrafficSecret, err := HkdfExpandLabel(
		handshakeSecret,
		clientHandshakeTrafficLabelName,
		transcriptHash,
		shaLength)
	if err != nil {
		return nil, err
	}

	serverHandshakeTrafficSecret, err := HkdfExpandLabel(
		handshakeSecret,
		serverHandshakeTrafficLabelName,
		transcriptHash,
		shaLength)
	if err != nil {
		return nil, err
	}

	clientHandshakeKey, err := HkdfExpandLabel(
		clientHandshakeTrafficSecret,
		secretKeyLabelName,
		nil,
		symmetricKeySize)
	if err != nil {
		return nil, err
	}

	serverHandshakeKey, err := HkdfExpandLabel(
		serverHandshakeTrafficSecret,
		secretKeyLabelName,
		nil,
		symmetricKeySize)
	if err != nil {
		return nil, err
	}

	clientHandshakeIv, err := HkdfExpandLabel(
		clientHandshakeTrafficSecret,
		secretIvLabelName,
		nil,
		ivSize)
	if err != nil {
		return nil, err
	}

	serverHandshakeIv, err := HkdfExpandLabel(
		serverHandshakeTrafficSecret,
		secretIvLabelName,
		nil,
		ivSize)
	if err != nil {
		return nil, err
	}

	secret := &SessionSecret{
		handshakeSecret:              handshakeSecret,
		clientHandshakeTrafficSecret: clientHandshakeTrafficSecret,
		serverHandshakeTrafficSecret: serverHandshakeTrafficSecret,
		ClientHandshakeKey:           clientHandshakeKey,
		ServerHandshakeKey:           serverHandshakeKey,
		ClientHandshakeIv:            clientHandshakeIv,
		ServerHandshakeIv:            serverHandshakeIv,
	}

	//Create the Client Handshake Cipher
	blockHandshakeClient, err := aes.NewCipher(clientHandshakeKey)
	if err != nil {
		return nil, err
	}

	secret.ClientHandshakeCipher, err = cipher.NewGCM(blockHandshakeClient)
	if err != nil {
		return nil, err
	}

	//Create the Server Handshake Cipher
	blockHandshakeServer, err := aes.NewCipher(serverHandshakeKey)
	if err != nil {
		return nil, err
	}

	secret.ServerHandshakeCipher, err = cipher.NewGCM(blockHandshakeServer)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (ss *SessionSecret) CreateApplicationSecrets(transcriptHash []byte) error {
	var hash crypto2.Hash
	hash = crypto2.SHA3_256
	emptyHash := hash.New().Sum(nil)

	derivedSecret, err := HkdfExpandLabel(
		ss.handshakeSecret,
		derivedLabelName,
		emptyHash,
		shaLength)
	if err != nil {
		return err
	}

	zeroKey := bytes.Repeat([]byte{0}, common.HashLength)
	masterSecret := hkdf.Extract(sha3.New256, zeroKey, derivedSecret)
	ss.masterSecret = masterSecret
	ss.TranscriptHash = transcriptHash

	clientApplicationTrafficSecret, err := HkdfExpandLabel(
		masterSecret,
		clientApplicationTrafficLabelName,
		transcriptHash,
		shaLength)
	if err != nil {
		return err
	}
	ss.clientApplicationTrafficSecret = clientApplicationTrafficSecret

	serverApplicationTrafficSecret, err := HkdfExpandLabel(
		masterSecret,
		serverApplicationTrafficLabelName,
		transcriptHash,
		shaLength)
	if err != nil {
		return err
	}
	ss.serverApplicationTrafficSecret = serverApplicationTrafficSecret

	clientApplicationKey, err := HkdfExpandLabel(
		clientApplicationTrafficSecret,
		secretKeyLabelName,
		nil,
		symmetricKeySize)
	if err != nil {
		return err
	}
	ss.ClientApplicationKey = clientApplicationKey

	serverApplicationKey, err := HkdfExpandLabel(
		serverApplicationTrafficSecret,
		secretKeyLabelName,
		nil,
		symmetricKeySize)
	if err != nil {
		return err
	}
	ss.ServerApplicationKey = serverApplicationKey

	clientApplicationIv, err := HkdfExpandLabel(
		clientApplicationTrafficSecret,
		secretIvLabelName,
		nil,
		ivSize)
	if err != nil {
		return err
	}
	ss.ClientApplicationIv = clientApplicationIv

	serverApplicationIv, err := HkdfExpandLabel(
		serverApplicationTrafficSecret,
		secretIvLabelName,
		nil,
		ivSize)
	if err != nil {
		return err
	}
	ss.ServerApplicationIv = serverApplicationIv

	//Create the Client Application Cipher
	blockApplicationClient, err := aes.NewCipher(clientApplicationKey)
	if err != nil {
		return err
	}

	ss.ClientApplicationCipher, err = cipher.NewGCM(blockApplicationClient)
	if err != nil {
		return err
	}

	//Create the Server Application Cipher
	blockApplicationServer, err := aes.NewCipher(serverApplicationKey)
	if err != nil {
		return nil
	}

	ss.ServerApplicationCipher, err = cipher.NewGCM(blockApplicationServer)
	if err != nil {
		return nil
	}

	return nil
}

func HkdfExpandLabel(secret []byte, label string, hashVal []byte, outputLength int) ([]byte, error) {
	hkdfLabel := hkdfEncodeLabel(label, hashVal, outputLength)

	reader := hkdf.Expand(sha3.New256, secret, hkdfLabel)
	output := make([]byte, outputLength)
	n, err := reader.Read(output)
	if err != nil {
		return nil, err
	}
	if n != outputLength {
		return nil, errors.New("invalid output length")
	}

	return output, err
}

func hkdfEncodeLabel(label string, hashVal []byte, outputLength int) []byte {
	fullLabel := "pqkem " + label

	fullLabelLen := len(fullLabel)
	hashLen := len(hashVal)
	hkdfLabel := make([]byte, 2+1+fullLabelLen+1+hashLen)
	hkdfLabel[0] = byte(outputLength >> 8)
	hkdfLabel[1] = byte(outputLength)
	hkdfLabel[2] = byte(fullLabelLen)
	copy(hkdfLabel[3:3+fullLabelLen], []byte(label))
	hkdfLabel[3+fullLabelLen] = byte(hashLen)
	copy(hkdfLabel[3+fullLabelLen+1:], hashVal)

	return hkdfLabel
}
