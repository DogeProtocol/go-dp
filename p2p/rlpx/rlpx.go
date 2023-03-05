// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package rlpx implements the RLPx transport protocol.
package rlpx

import (
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"net"
	"time"
)

// Conn is an RLPx network connection. It wraps a low-level network connection. The
// underlying connection should not be used for other activity when it is wrapped by Conn.
//
// Before sending messages, a handshake must be performed by calling the Handshake method.
// This type is not generally safe for concurrent use, but reading and writing of messages
// may happen concurrently after the handshake.
type Conn struct {
	dialDest *signaturealgorithm.PublicKey
	conn     net.Conn

	// These are the buffers for snappy compression.
	// Compression is enabled if they are non-nil.
	snappyReadBuffer  []byte
	snappyWriteBuffer []byte

	client *Client

	server *Server

	context string
}

// NewConn wraps the given network connection. If dialDest is non-nil, the connection
// behaves as the initiator during the handshake.
func NewConn(conn net.Conn, dialDest *signaturealgorithm.PublicKey, context string) *Conn {
	connection := &Conn{
		dialDest: dialDest,
		conn:     conn,
		context:  context,
	}

	if dialDest == nil {
		connection.server = NewServer(conn, nil, context)
	} else {
		connection.client = NewClient(conn, nil, dialDest, context)
	}

	return connection
}

// SetSnappy enables or disables snappy compression of messages. This is usually called
// after the devp2p Hello message exchange when the negotiated version indicates that
// compression is available on both ends of the connection.
func (c *Conn) SetSnappy(snappy bool) {
	if snappy {

		c.snappyReadBuffer = []byte{}
		c.snappyWriteBuffer = []byte{}
	} else {

		c.snappyReadBuffer = nil
		c.snappyWriteBuffer = nil
	}
}

// SetReadDeadline sets the deadline for all future read operations.
func (c *Conn) SetReadDeadline(deadlineTime time.Time) error {

	return c.conn.SetReadDeadline(deadlineTime)
}

// SetWriteDeadline sets the deadline for all future write operations.
func (c *Conn) SetWriteDeadline(deadlineTime time.Time) error {
	return c.conn.SetWriteDeadline(deadlineTime)
}

// SetDeadline sets the deadline for all future read and write operations.
func (c *Conn) SetDeadline(deadlineTime time.Time) error {
	return c.conn.SetDeadline(deadlineTime)
}

// Read reads a message from the connection.
// The returned data buffer is valid until the next call to Read.
func (c *Conn) Read() (code uint64, data []byte, wireSize int, err error) {

	if c.client != nil {
		dataPacket, err := c.client.ReadAndDecrypt(PacketTypeApplicationData)
		if err != nil {

			return 0, nil, 0, err
		}

		return dataPacket.context, dataPacket.fragment, len(dataPacket.fragment), nil
	} else {
		dataPacket, err := c.server.ReadAndDecrypt(PacketTypeApplicationData)
		if err != nil {

			return 0, nil, 0, err
		}

		return dataPacket.context, dataPacket.fragment, len(dataPacket.fragment), nil
	}
}

// Write writes a message to the connection.
//
// Write returns the written size of the message data. This may be less than or equal to
// len(data) depending on whether snappy compression is enabled.
func (c *Conn) Write(code uint64, data []byte) (uint32, error) {

	size := uint32(len(data))

	if c.client != nil {
		err := c.client.WriteEncrypted(data, code, PacketTypeApplicationData)
		if err != nil {

			return size, err
		}
	} else {
		err := c.server.WriteEncrypted(data, code, PacketTypeApplicationData)
		if err != nil {

			return size, err
		}
	}

	return size, nil
}

// Handshake performs the handshake. This must be called before any data is written
// or read from the connection.
func (c *Conn) Handshake(prv *signaturealgorithm.PrivateKey) (*signaturealgorithm.PublicKey, error) {
	if c.client != nil {
		c.client.SetClientSigningPrivateKey(prv)
		err := c.client.PerformHandshake()
		if err != nil {

			return nil, err
		}
		return c.client.serverSigningPublicKey, nil
	} else {
		c.server.SetServerSigningPrivateKey(prv)
		err := c.server.PerformHandshake()
		if err != nil {

			return nil, err
		}

		return c.server.clientSigningPublicKey, nil
	}
	return nil, nil
}

// Close closes the underlying network connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) InitWithSecrets(secret SessionSecret) {
	if c.client != nil {
		c.client.InitWithSecrets(secret)
	} else {
		c.server.InitWithSecrets(secret)
	}
}
