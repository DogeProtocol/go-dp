package oqs

type Dilithium struct {
	OqsSig
}

func InitDilithium() Dilithium {
	return Dilithium{
		CreateOqs("Dilithium2"),
	}
}
