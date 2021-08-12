package utils

import (
	"crypto/rand"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return nil if error occurs
func GenerateRandomBytes(n int) []byte {

	// Init byte slice
	b := make([]byte, n)

	// Generate random bytes and write to b
	_, err := rand.Read(b)

	if err != nil {
		return nil
	}

	return b
}
