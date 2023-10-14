package ChaCha20

import (
	"fmt"
	"testing"
)

func TestChaCha20DRNG(t *testing.T) {
	for ctr := 0; ctr <= 10; ctr++ {
		var seed1 [32]byte
		var i byte
		for i = 0; i < 32; i++ {
			seed1[i] = i
		}

		g1 := ChaCha20DRNGInitializer{}
		rng1, err := g1.InitializeWithSeed(seed1)
		if err != nil {
			t.Fail()
		}

		var seed2 [32]byte
		for i = 0; i < 32; i++ {
			seed1[i] = i + 32
		}

		g2 := ChaCha20DRNGInitializer{}
		rng2, err := g2.InitializeWithSeed(seed2)
		if err != nil {
			t.Fail()
		}

		expected1 := [8]byte{
			73, 25, 207, 89, 253, 231, 30, 214,
		}

		expected2 := [8]byte{
			240, 54, 186, 159, 50, 248, 142, 193,
		}

		for i = 0; i < 8; i++ {
			r1 := rng1.NextByte()
			r2 := rng2.NextByte()
			fmt.Println("rng", i, r1, r2)
			if r1 != expected1[i] {
				t.Fail()
			}
			if r2 != expected2[i] {
				t.Fail()
			}
			if r1 == r2 {
				t.Fail()
			}
		}
	}
}
