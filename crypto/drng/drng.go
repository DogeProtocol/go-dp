package drng

type DRNG interface {
	NextByte() byte
}

type DRNGInitializer interface {
	InitializeWithSeed(seed [32]byte) (*DRNG, error)
}
