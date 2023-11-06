package hashingalgorithm

import (
	"bytes"
)
import "testing"

func HashStateSumTest(t *testing.T, h HashState) {
	msg := []byte("abc")
	d1 := make([]byte, h.Size())

	_, err := h.Write(msg)
	if err != nil {
		t.Fatal(err)
	}

	h.Read(d1)
	d2 := h.Sum(nil)

	if bytes.Compare(d1, d2) != 0 {
		t.Fatal(err)
	}

	h.Reset()
	_, err = h.Write(msg)
	if err != nil {
		t.Fatal(err)
	}

	h.Read(d1)
	d2 = h.Sum(nil)

	if bytes.Compare(d1, d2) != 0 {
		t.Fatal(err)
	}

	h.Reset()

	_, err = h.Write(msg)
	if err != nil {
		t.Fatal(err)
	}

	h.Read(d1)
	d2 = h.Sum(msg)

	if bytes.Compare(d1, d2) == 0 {
		t.Fatal(err)
	}

	h.Read(d1)
	d2 = h.Sum(msg)
	d3 := h.Sum(msg)

	if bytes.Compare(d1, d2) == 0 {
		t.Fatal(err)
	}

	if bytes.Compare(d2, d3) != 0 {
		t.Fatal(err)
	}
}

func HashStateTest(t *testing.T, h1 HashState, h2 HashState) {
	msg := []byte("abc")

	d1 := make([]byte, h1.Size())
	d2 := make([]byte, h1.Size())
	d3 := make([]byte, h1.Size())
	d4 := make([]byte, h1.Size())

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

	h1.Read(d3)
	h2.Read(d4)

	if bytes.Compare(d3, d4) != 0 {
		t.Fatal(err)
	}

	d5 := h1.Sum(msg)
	d6 := h1.Sum(msg)

	if bytes.Compare(d5, d6) != 0 {
		t.Fatal(err)
	}
}
