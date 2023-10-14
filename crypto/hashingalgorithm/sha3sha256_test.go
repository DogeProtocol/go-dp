package hashingalgorithm

import (
	"testing"
)

func TestHashStateSum_sha3sha256(t *testing.T) {
	h := NewSha3Sha256HashState()
	HashStateSumTest(t, h)
}

func TestHashState_sha3sha256(t *testing.T) {
	h1 := NewSha3Sha256HashState()
	h2 := NewSha3Sha256HashState()
	HashStateTest(t, h1, h2)
}
