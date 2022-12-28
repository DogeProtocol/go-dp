package discover

import (
	"bytes"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"net"
	"sync"
)

const (
	maxUdpPacketSize = 1000
	packetPrefix     = "ch1p" //Chunk version v1
	hashSize         = 32
	packetHeadSize   = hashSize + len(packetPrefix) + 1
	maxChunkDataSize = maxUdpPacketSize - packetHeadSize //- len(packetSuffix)
)

type DpUdpSession struct {
	BaseConn net.Conn
	addr     *net.UDPAddr
	mutex    sync.Mutex
	buff     []byte
	hash     []byte
}

func (c *DpUdpSession) Read(inBuff []byte) (int, error) {
	n, err := c.BaseConn.Read(inBuff)

	if err != nil {
		return n, err
	}

	inputSize := n
	if inputSize < packetHeadSize {
		return n, nil
	}

	incomingPacketPrefix := string(inBuff[:len(packetPrefix)])
	if incomingPacketPrefix != packetPrefix {
		return n, nil
	}

	hash := inBuff[len(packetPrefix) : len(packetPrefix)+hashSize]
	if len(hash) != hashSize {
		return n, errors.New("invalid hash size")
	}

	isLastPacket := inBuff[len(packetPrefix)+hashSize : len(packetPrefix)+hashSize+1][0]
	if isLastPacket != 0 && isLastPacket != 1 {
		return n, errors.New("invalid packet type read")
	}

	chunk := inBuff[packetHeadSize:n]

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.hash) == 0 {

		c.hash = make([]byte, len(hash))
		copy(c.hash, hash)

		c.buff = make([]byte, len(chunk))
		copy(c.buff, chunk)

	} else {
		if bytes.Compare(c.hash, hash) == 0 {
			c.buff = append(c.buff, chunk...)

		} else {
			c.hash = make([]byte, len(hash))
			copy(c.hash, hash)

			c.buff = make([]byte, len(chunk))
			copy(c.buff, chunk)
		}
	}

	if isLastPacket == 1 {
		if len(inBuff) < len(c.buff) {
			return 0, errors.New("buff size mismatch")
		}

		bufHash := crypto.Keccak256(c.buff)
		if bytes.Compare(bufHash, c.hash) != 0 {
			return 0, errors.New("hash mismatch")
		}

		copy(inBuff, c.buff)
		buffLen := len(c.buff)
		c.buff = c.buff[:0]
		c.hash = c.hash[:0]

		return buffLen, nil
	} else {
		return 0, nil
	}
}

func (c *DpUdpSession) Write(b []byte) (n int, err error) {
	return c.writeChunked(b)
}

func (c *DpUdpSession) Close() error {
	return c.BaseConn.Close()
}

func (c *DpUdpSession) writeChunked(b []byte) (n int, err error) {

	inputSize := len(b)

	hash := crypto.Keccak256(b)
	if len(hash) != hashSize {
		return 0, errors.New("invalid hash")
	}

	prefix := []byte(packetPrefix)
	prefixPacket := append(prefix, hash...)

	startPos := 0
	i := 0
	var writtenOutput []byte
	for {
		i = i + 1
		endPos := startPos + maxChunkDataSize
		if endPos > inputSize {
			endPos = inputSize
		}

		chunk := prefixPacket
		if endPos == inputSize {
			chunk = append(chunk, 1) //last packet
		} else {
			chunk = append(chunk, 0) //other packets
		}

		chunk = append(chunk, b[startPos:endPos]...)

		writtenOutput = append(writtenOutput, b[startPos:endPos]...)

		_, err := c.BaseConn.Write(chunk)
		if err != nil {
			return 0, err
		}

		if endPos == inputSize {
			break
		}

		startPos = endPos
	}

	if bytes.Compare(writtenOutput, b) != 0 {
		return 0, errors.New("write failed invalid output")
	}

	return inputSize, nil
}
