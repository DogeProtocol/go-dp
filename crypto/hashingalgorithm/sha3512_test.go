package hashingalgorithm

import (
	"testing"
)

func TestHashStateSum_sha3512(t *testing.T) {
	h := NewSha3512HashState()
	HashStateSumTest(t, h)
}

func TestHashState_sha3512(t *testing.T) {
	h1 := NewSha3512HashState()
	h2 := NewSha3512HashState()
	HashStateTest(t, h1, h2)
}
