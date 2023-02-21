package oqs

type Falcon struct {
	OqsSig
}

func InitFalcon() Falcon {
	return Falcon{
		CreateOqs("Falcon-512"),
	}
}
