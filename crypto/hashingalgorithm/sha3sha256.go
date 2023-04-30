package hashingalgorithm

import (
	"bytes"
	"crypto/sha256"
	"golang.org/x/crypto/sha3"
	"hash"
)

type Sha3Sha256HashState struct {
	sha3   HashState
	sha256 hash.Hash
	buff   *bytes.Buffer
}

func NewSha3Sha256HashState() Sha3Sha256HashState {
	return Sha3Sha256HashState{
		sha3:   sha3.NewLegacyKeccak256().(HashState),
		sha256: sha256.New(),
		buff:   new(bytes.Buffer),
	}
}

func (s Sha3Sha256HashState) Write(p []byte) (n int, err error) {
	return s.buff.Write(p)
}

func (s Sha3Sha256HashState) Sum(b []byte) []byte {
	s.sha3.Reset()

	if b == nil {
		_, err := s.sha3.Write(s.buff.Bytes())
		if err != nil {
			return nil
		}
	} else {
		totalBuffer := CopyArrays(s.buff.Bytes(), b)
		_, err := s.sha3.Write(totalBuffer)
		if err != nil {
			return nil
		}
	}

	hashBytes := make([]byte, 32)
	_, err := s.sha3.Read(hashBytes)
	if err != nil {
		return nil
	}

	return hashBytes

	/*s.sha3.Reset()
	s.sha256.Reset()

	totalBuffer := CopyArrays(s.buff.Bytes(), b)
	sha256Bytes := s.sha256.Sum(totalBuffer)[12:]
	tempBuffer := CopyArrays(totalBuffer, sha256Bytes)
	_, err := s.sha3.Write(tempBuffer)
	if err != nil {
		return nil
	}

	hashBytes := make([]byte, 32)
	_, err = s.sha3.Read(hashBytes)
	if err != nil {
		return nil
	}

	return hashBytes*/
}

func (s Sha3Sha256HashState) Reset() {
	s.buff.Reset()
}

func (s Sha3Sha256HashState) Size() int {
	return s.sha256.Size()
}

func (s Sha3Sha256HashState) BlockSize() int {
	return s.sha256.BlockSize()
}

func (s Sha3Sha256HashState) Read(b []byte) (int, error) {
	s.sha3.Reset()
	_, err := s.sha3.Write(s.buff.Bytes())
	if err != nil {
		return 0, err
	}

	return s.sha3.Read(b)

	/*s.sha3.Reset()
	s.sha256.Reset()
	sha256Bytes := s.sha256.Sum(s.buff.Bytes())[12:]
	tempBuffer := CopyArrays(s.buff.Bytes(), sha256Bytes)
	_, err := s.sha3.Write(tempBuffer)
	if err != nil {
		return 0, err
	}

	return s.sha3.Read(b)*/
}

func CopyArrays(array1 []byte, array2 []byte) []byte {
	la := len(array1)
	c := make([]byte, la, la+len(array2))
	_ = copy(c, array1)
	c = append(c, array2...)
	return c
}
