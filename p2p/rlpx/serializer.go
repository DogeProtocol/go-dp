package rlpx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/DogeProtocol/dp/rlp"
	"io"
	mrand "math/rand"
	"sync"
)

type Serializer interface {
	Serialize(msg interface{}) ([]byte, error)
	SerializeDeterministic(msg interface{}, padLen int) ([]byte, error)
	Deserialize(msg interface{}, reader io.Reader) ([]byte, error)
	SetContext(context string)
}

type RlpxSerializer struct {
	rbuf    ReadBuffer
	wbuf    WriteBuffer
	context string
	mutex   sync.Mutex
}

func NewRlpxSerializer() Serializer {
	return &RlpxSerializer{
		mutex: sync.Mutex{},
	}
}

func (rs *RlpxSerializer) SetContext(context string) {
	rs.context = context
}

func (rs *RlpxSerializer) SerializeDeterministic(msg interface{}, padLen int) ([]byte, error) {
	rs.wbuf.Reset()

	// Write the message plaintext.
	if err := rlp.Encode(&rs.wbuf, msg); err != nil {
		return nil, err
	}

	// Pad with random amount of data. the amount needs to be at least 100 bytes to make
	// the message distinguishable from pre-EIP-8 handshakes.
	rs.wbuf.AppendZero(padLen)

	prefix := make([]byte, 2)

	binary.BigEndian.PutUint16(prefix, uint16(len(rs.wbuf.Data)))

	return append(prefix, rs.wbuf.Data...), nil
}

func (rs *RlpxSerializer) Serialize(msg interface{}) ([]byte, error) {
	padLen := mrand.Intn(100) + 100
	return rs.SerializeDeterministic(msg, padLen)
}

func (rs *RlpxSerializer) Deserialize(msg interface{}, reader io.Reader) ([]byte, error) {
	//rs.rbuf.Reset()

	// Read the size prefix.

	prefixSize := 2
	prefix := make([]byte, prefixSize)
	bytesRead, err := io.ReadAtLeast(reader, prefix, prefixSize)
	if err != nil {

		return nil, err
	}
	if bytesRead != prefixSize {

		return nil, errors.New("prefix size less")
	}

	size := binary.BigEndian.Uint16(prefix)

	packet := make([]byte, int(size))
	bytesRead, err = io.ReadAtLeast(reader, packet, int(size))
	if err != nil {

		return nil, err
	}

	if bytesRead != int(size) {

		return nil, errors.New("prefix size less")
	}

	if len(packet) != int(size) {

	}

	s := rlp.NewStream(bytes.NewReader(packet), 0)
	err = s.Decode(msg)

	return packet, err
}
