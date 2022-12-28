package rlpx

import (
	"crypto/cipher"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"time"
)

// Constants for the handshake.
const (
	//pubLen          = oqs.PublicKeyLen
	shaLength           = 32 // hash length (for nonce etc)
	kemPublicKeyLen     = 1138
	symmetricKeySize    = 32
	ivSize              = 12
	KemName             = "NTRU-HRSS-701"
	kemCipherTextLength = 1138
	kemSecretLength     = 32

	majorVersion = 1
	minorVersion = 1

	padLen = 0
	shaLen = 32
)

type PacketType byte

const (
	PacketTypeHandshake       PacketType = 21
	PacketTypeApplicationData PacketType = 23
	ReadTimeout                          = time.Second * 10
	WriteTimeout                         = time.Second * 20
)

type DataPacket struct {
	packetType PacketType
	seqNum     uint
	fragment   []byte
	context    uint64
}

type Header struct {
	PacketType     uint
	MinorVersion   uint
	MajorVersion   uint
	RecordLength   uint
	Context        uint64
	AdditionalData [common.HashLength]byte
	Rest           []rlp.RawValue `rlp:"tail"`
}

type AdditionalData struct {
	PacketType   uint
	MinorVersion uint
	MajorVersion uint
	DataLength   uint
	Rest         []rlp.RawValue `rlp:"tail"`
}

func CalculateNonce(recordCount uint, input []byte) []byte {
	inputLen := len(input)
	output := make([]byte, inputLen)
	copy(output, input)

	rec := recordCount

	for i := 0; i < 8; i++ {
		output[(inputLen-i)-1] ^= byte(rec & 0xff)
		rec >>= 8
	}

	return output
}

func Encrypt(cipher1 cipher.AEAD, fragment []byte, additionalData []byte, packetType PacketType, handshakeIv []byte, seqNum uint) (encrypted []byte, err error) {
	dataLen := len(fragment)

	nonce := CalculateNonce(seqNum, handshakeIv)

	//Calculate packet overhead
	beforeEncryptLen := dataLen + 1 + padLen
	encryptedLen := beforeEncryptLen + cipher1.Overhead()

	//Create array to store encrypted data with overhead
	buffer := make([]byte, encryptedLen)
	copy(buffer, fragment)
	buffer[dataLen] = byte(packetType)
	for i := 1; i <= padLen; i++ {
		buffer[dataLen+i] = 0
	}

	//Encrypt the data
	payload := buffer[:beforeEncryptLen]
	encryptedData := cipher1.Seal(payload[:0], nonce, payload, additionalData)

	return encryptedData, nil
}

func Decrypt(cipher1 cipher.AEAD, encryptedData []byte, additionalData []byte, packetType PacketType, iv []byte, seqNum uint) (*DataPacket, error) {
	if len(encryptedData) < cipher1.Overhead() {
		return nil, errors.New("invalid data")
	}

	dataLen := len(encryptedData) - cipher1.Overhead()
	dataPacket := &DataPacket{
		packetType: packetType,
		fragment:   make([]byte, dataLen),
	}

	//Compute the nonce
	nonce := CalculateNonce(seqNum, iv)

	// Decrypt
	_, err := cipher1.Open(dataPacket.fragment[:0], nonce, encryptedData, additionalData)
	if err != nil {
		return nil, err
	}

	// Find the padding boundary
	padLen1 := padLen
	for ; padLen1 < dataLen+1 && dataPacket.fragment[dataLen-padLen1-1] == 0; padLen1++ {
	}

	// Transfer the content type
	newLen := dataLen - padLen1 - 1
	dataPacket.packetType = PacketType(dataPacket.fragment[newLen])

	dataPacket.fragment = dataPacket.fragment[:newLen]
	dataPacket.seqNum = seqNum

	return dataPacket, nil
}
