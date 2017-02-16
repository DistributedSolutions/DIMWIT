package random

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

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

func RandomHexStringOfSize(n int) (string, error) {
	if n%2 != 0 {
		return "", fmt.Errorf("Must be even length")
	}

	data := RandByteSliceOfSize(n / 2)
	return hex.EncodeToString(data), nil
}

func RandomUInt32Between(min uint32, max uint32) uint32 {
	return uint32(RandomIntBetween(int(min), int(max)))
}

func RandomUInt32() uint32 {
	return uint32(RandomInt())
}

func RandomInt() int {
	return rand.Int()
}

func RandomIntBetween(min int, max int) int {
	if min == max {
		max = min + 1
	}
	return rand.Intn(max-min) + min
}

func RandomInt64() int64 {
	return rand.Int63()
}

func RandomInt64Between(min int64, max int64) int64 {
	if min == max {
		max = min + 1
	}
	return rand.Int63n(max-min) + min
}

func RandByteSlice() []byte {
	l := RandomInt() % 64
	answer := make([]byte, l)
	_, err := rand.Read(answer)
	if err != nil {
		return nil
	}
	return answer
}

func RandByteSliceOfSize(l int) []byte {
	if l <= 0 {
		return nil
	}
	answer := make([]byte, l)
	_, err := rand.Read(answer)
	if err != nil {
		return nil
	}
	return answer
}

func RandomBool() bool {
	b := RandomIntBetween(0, 100)
	if b > 50 {
		return true
	} else {
		return false
	}
}
