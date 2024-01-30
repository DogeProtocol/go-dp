package keystore

import (
	"errors"
)

var (
	ErrLocked               = errors.New("password or unlock")
	ErrNoMatch              = errors.New("no key for given address or file")
	ErrDecrypt              = errors.New("could not decrypt key with given password")
	ErrAccountAlreadyExists = errors.New("account already exists")
)
