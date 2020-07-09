package artisan

import (
	"crypto/rand"
)

type UUID = []byte

func init() {

}

func NewUUID() (UUID, error) {
	var dest = make([]byte, 16)
	if _, err := rand.Read(dest); err != nil {
		return nil, err
	}
	return dest, nil
}

func MustUUID() UUID {
	x, err := NewUUID()
	if err != nil {
		panic(err)
	}
	return x
}
