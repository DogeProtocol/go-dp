package hashingalgorithm

import (
	"golang.org/x/crypto/sha3"
)

type Sha3HashState struct {
	sha3 HashState
}

func NewSha3HashState() Sha3HashState {
	return Sha3HashState{
		sha3: sha3.New256().(HashState),
	}
}

func (s Sha3HashState) Write(p []byte) (n int, err error) {
	return s.sha3.Write(p)
}

func (s Sha3HashState) Sum(b []byte) []byte {
	return s.sha3.Sum(b)
}

func (s Sha3HashState) Reset() {
	s.sha3.Reset()
}

func (s Sha3HashState) Size() int {
	return s.sha3.Size()
}

func (s Sha3HashState) BlockSize() int {
	return s.sha3.BlockSize()
}

func (s Sha3HashState) Read(b []byte) (int, error) {
	return s.sha3.Read(b)
}
