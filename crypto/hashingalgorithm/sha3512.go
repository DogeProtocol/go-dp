package hashingalgorithm

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/sha3"
)

type Sha3512HashState struct {
	sha3 HashState
	buff *bytes.Buffer
}

func NewSha3512HashState() Sha3512HashState {
	return Sha3512HashState{
		sha3: sha3.New512().(HashState),
		buff: new(bytes.Buffer),
	}
}

func (s Sha3512HashState) Write(p []byte) (n int, err error) {
	return s.buff.Write(p)
}

func (s Sha3512HashState) Sum(b []byte) []byte {
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
	fmt.Println("sizeinner", s.Size())
	hashBytes := make([]byte, s.Size())
	_, err = s.sha3.Read(hashBytes)
	if err != nil {
		return nil
	}

	return hashBytes
}

func (s Sha3512HashState) Reset() {
	s.buff.Reset()
}

func (s Sha3512HashState) Size() int {
	return s.sha3.Size()
}

func (s Sha3512HashState) BlockSize() int {
	return s.sha3.BlockSize()
}

func (s Sha3512HashState) Read(b []byte) (int, error) {
	s.sha3.Reset()

	_, err := s.sha3.Write(s.buff.Bytes())
	if err != nil {
		return 0, err
	}

	return s.sha3.Read(b)
}
