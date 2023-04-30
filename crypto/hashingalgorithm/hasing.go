package hashingalgorithm

import "hash"

type HashState interface {
	hash.Hash
	Read([]byte) (int, error)
}

func NewHashState() HashState {
	return NewSha3HashState()
}
