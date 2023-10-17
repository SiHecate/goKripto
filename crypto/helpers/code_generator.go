package helpers

import (
	"math/rand"
	"time"
)

func GenerateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var codes [6]byte
	for i := 0; i < 6; i++ {
		codes[i] = uint8(48 + r.Intn(10))
	}

	return string(codes[:])
}
