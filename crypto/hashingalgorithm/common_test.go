package hashingalgorithm

import "bytes"
import "testing"

func HashStateSumTest(t *testing.T, h HashState) {
	msg := []byte("abc")
	d1 := make([]byte, 32)

	_, err := h.Write(msg)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		h.Read(d1)
		d2 := h.Sum(nil)

		if bytes.Compare(d1, d2) != 0 {
			t.Fatal(err)
		}
	}

	for i := 0; i < 10; i++ {
		_, err := h.Write(msg)
		if err != nil {
			t.Fatal(err)
		}

		h.Read(d1)
		d2 := h.Sum(nil)

		if bytes.Compare(d1, d2) != 0 {
			t.Fatal(err)
		}
	}

	h.Reset()
	for i := 0; i < 10; i++ {
		_, err := h.Write(msg)
		if err != nil {
			t.Fatal(err)
		}

		h.Read(d1)
		d2 := h.Sum(msg)

		if bytes.Compare(d1, d2) == 0 {
			t.Fatal(err)
		}
	}

	for i := 0; i < 10; i++ {
		h.Read(d1)
		d2 := h.Sum(msg)
		d3 := h.Sum(msg)

		if bytes.Compare(d1, d2) == 0 {
			t.Fatal(err)
		}

		if bytes.Compare(d2, d3) != 0 {
			t.Fatal(err)
		}
	}
}

func HashStateTest(t *testing.T, h1 HashState, h2 HashState) {
	msg := []byte("abc")

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

			d5 := h1.Sum(msg)
			d6 := h1.Sum(msg)

			if bytes.Compare(d5, d6) != 0 {
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
