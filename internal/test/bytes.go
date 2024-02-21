package test

import (
	"bytes"
	"math/rand"
	"testing"
)

func RandomInRange(start int32, end int32) int32 {
	return start + rand.Int31n(end-start)
}

func RandomTestBytes() []byte {
	result := make([]byte, RandomInRange(1024, 4096))
	rand.Read(result)
	return result
}

func BytesCompare(t *testing.T, expected []byte, actual []byte) {
	if !bytes.Equal(expected, actual) {
		t.Fatalf(
			"data not equal, expected %v...%v(%d), got %v...%v(%d)",
			expected[:16], expected[len(expected)-16:], len(expected),
			actual[:16], actual[len(actual)-16:], len(actual),
		)
	}
}
