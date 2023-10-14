package hashingalgorithm

import (
	"bytes"
	"golang.org/x/crypto/sha3"
)

type Sha3HashState struct {
	sha3 HashState
	buff *bytes.Buffer
}

func NewSha3HashState() Sha3HashState {
	return Sha3HashState{
		sha3: sha3.NewLegacyKeccak256().(HashState),
		buff: new(bytes.Buffer),
	}
}

func (s Sha3HashState) Write(p []byte) (n int, err error) {
	return s.buff.Write(p)
}

func (s Sha3HashState) Sum(b []byte) []byte {
	s.sha3.Reset()

	var totalBuffer []byte

	if b == nil {
		totalBuffer = s.buff.Bytes()
	} else {
		totalBuffer = CopyArrays(s.buff.Bytes(), b)
	}

	_, err := s.sha3.Write(totalBuffer)
	if err != nil {
		return nil
	}
	hashBytes := make([]byte, 32)
	_, err = s.sha3.Read(hashBytes)
	if err != nil {
		return nil
	}

	return hashBytes
}

func (s Sha3HashState) Reset() {
	s.buff.Reset()
}

func (s Sha3HashState) Size() int {
	return s.sha3.Size()
}

func (s Sha3HashState) BlockSize() int {
	return s.sha3.BlockSize()
}

func (s Sha3HashState) Read(b []byte) (int, error) {
	s.sha3.Reset()

	_, err := s.sha3.Write(s.buff.Bytes())
	if err != nil {
		return 0, err
	}

	return s.sha3.Read(b)
}
