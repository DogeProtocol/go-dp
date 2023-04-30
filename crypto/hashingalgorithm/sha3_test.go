package hashingalgorithm

import "testing"

func TestHashStateSum_sha3(t *testing.T) {
	h := NewSha3HashState()
	HashStateSumTest(t, h)
}

func TestHashState_sha3(t *testing.T) {
	h1 := NewSha3HashState()
	h2 := NewSha3HashState()
	HashStateTest(t, h1, h2)
}
