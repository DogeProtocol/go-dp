// Copyright 2021 The go-ethereum Authors
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

package rlpx

import (
	"io"
)

// ReadBuffer implements buffering for network reads. This type is similar to bufio.Reader,
// with two crucial differences: the buffer slice is exposed, and the buffer keeps all
// read data available until reset.
//
// How to use this type:
//
// Keep a ReadBuffer b alongside the underlying network connection. When reading a packet
// from the connection, first call b.reset(). This empties b.data. Now perform reads
// through b.read() until the end of the packet is reached. The complete packet data is
// now available in b.data.
type ReadBuffer struct {
	Data []byte
	end  int
}

// Reset removes all processed data which was read since the last call to Reset.
// After Reset, len(b.data) is zero.
func (b *ReadBuffer) Reset() {
	unprocessed := b.end - len(b.Data)
	copy(b.Data[:unprocessed], b.Data[len(b.Data):b.end])
	b.end = unprocessed
	b.Data = b.Data[:0]
}

// Read reads at least n bytes from r, returning the bytes.
// The returned slice is valid until the next call to Reset.
func (b *ReadBuffer) Read(r io.Reader, n int) ([]byte, error) {


	offset := len(b.Data)
	have := b.end - len(b.Data)

	// If n bytes are available in the buffer, there is no need to read from r at all.
	if have >= n {
		b.Data = b.Data[:offset+n]
		return b.Data[offset : offset+n], nil
	}

	// Make buffer space available.
	need := n - have
	b.Grow(need)


	// Read.
	rn, err := io.ReadAtLeast(r, b.Data[b.end:cap(b.Data)], need)
	if err != nil {
		return nil, err
	}
	b.end += rn
	b.Data = b.Data[:offset+n]
	return b.Data[offset : offset+n], nil
}

// Grow ensures the buffer has at least n bytes of unused space.
func (b *ReadBuffer) Grow(n int) {
	if cap(b.Data)-b.end >= n {
		return
	}
	need := n - (cap(b.Data) - b.end)
	offset := len(b.Data)
	b.Data = append(b.Data[:cap(b.Data)], make([]byte, need)...)
	b.Data = b.Data[:offset]
}

// WriteBuffer implements buffering for network writes. This is essentially
// a convenience wrapper around a byte slice.
type WriteBuffer struct {
	Data []byte
}

func (b *WriteBuffer) Reset() {
	b.Data = b.Data[:0]
}

func (b *WriteBuffer) AppendZero(n int) []byte {
	offset := len(b.Data)
	b.Data = append(b.Data, make([]byte, n)...)
	return b.Data[offset : offset+n]
}

func (b *WriteBuffer) Write(data []byte) (int, error) {
	b.Data = append(b.Data, data...)
	return len(data), nil
}

const maxUint24 = int(^uint32(0) >> 8)

func readUint24(b []byte) uint32 {
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}

func putUint24(v uint32, b []byte) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

// growslice ensures b has the wanted length by either expanding it to its capacity
// or allocating a new slice if b has insufficient capacity.
func growslice(b []byte, wantLength int) []byte {
	if len(b) >= wantLength {
		return b
	}
	if cap(b) >= wantLength {
		return b[:cap(b)]
	}
	return make([]byte, wantLength)
}
