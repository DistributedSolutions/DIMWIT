package random

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringOfSize(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Max size 1000
func RandString() string {
	return RandStringOfSize(rand.Intn(1000))
}

func RandomInt() int {
	return rand.Int()
}

func RandomIntBetween(min int, max int) int {
	return rand.Intn(max-min) + min
}
