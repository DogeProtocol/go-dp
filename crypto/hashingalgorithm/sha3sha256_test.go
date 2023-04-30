package hashingalgorithm

import (
	"bytes"
	"testing"
)

func TestHashState(t *testing.T) {
	msg := []byte("abc")
	h1 := NewHashState()
	h2 := NewHashState()

	d1 := make([]byte, 32)
	d2 := make([]byte, 32)
	d3 := make([]byte, 32)
	d4 := make([]byte, 32)

	for i := 0; i < 100; i++ {
		_, err := h1.Write(msg)
		if err != nil {
			t.Fatal(err)
		}
		_, err = h2.Write(msg)
		if err != nil {
			t.Fatal(err)
		}

		h1.Read(d1)
		h2.Read(d2)

		if bytes.Compare(d1, d2) != 0 {
			t.Fatal(err)
		}

		for j := 0; j < 100; j++ {
			h1.Read(d3)
			h2.Read(d4)

			if bytes.Compare(d3, d4) != 0 {
				t.Fatal(err)
			}

			if bytes.Compare(d1, d3) != 0 {
				t.Fatal(err)
			}
		}
	}

	for i := 0; i < 100; i++ {
		h1.Reset()
		h2.Reset()

		_, err := h1.Write(msg)
		if err != nil {
			t.Fatal(err)
		}
		_, err = h2.Write(msg)
		if err != nil {
			t.Fatal(err)
		}

		h1.Read(d1)
		h2.Read(d2)

		if bytes.Compare(d1, d2) != 0 {
			t.Fatal(err)
		}

		for j := 0; j < 100; j++ {
			h1.Read(d3)
			h2.Read(d4)

			if bytes.Compare(d3, d4) != 0 {
				t.Fatal(err)
			}

			if bytes.Compare(d1, d3) != 0 {
				t.Fatal(err)
			}
		}
	}

	_, err := h1.Write(msg)
	if err != nil {
		t.Fatal(err)
	}

	h1.Read(d1)
	h2.Read(d2)

	if bytes.Compare(d1, d2) == 0 {
		t.Fatal(err)
	}

}
