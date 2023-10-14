package hashingalgorithm

func CopyArrays(array1 []byte, array2 []byte) []byte {
	la := len(array1)
	c := make([]byte, la, la+len(array2))
	_ = copy(c, array1)
	c = append(c, array2...)
	return c
}
