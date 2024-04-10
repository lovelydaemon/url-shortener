package random

import (
	"math/rand"
	"time"
)

var _defaultStringLength = 9

// NewRandomString generates a random string
func NewRandomString() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune(
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"abcdefghijklmnopqrstuvwxyz" +
			"0123456789",
	)

	b := make([]rune, _defaultStringLength)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
