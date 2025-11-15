// Package utils provides utility functions for the URL shortener service.
//
// This package contains helper functions for generating cryptographically secure
// random short IDs used in URL shortening. The IDs are composed of alphanumeric
// characters (a-z, A-Z, 0-9) and use crypto/rand for secure random generation.
//
// Example usage:
//
//	// Generate a 6-character short ID
//	shortID, err := utils.GenerateShortID(6)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(shortID) // Output: "aBc123" (example, will be random)
//
//	// Generate a 10-character short ID
//	longID, err := utils.GenerateShortID(10)
//	if err != nil {
//	    log.Fatal(err)
//	}
package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var lettersRune = []rune("abcdefghijlmnopqrstuvwxyzABCDEFGHIJLMNOPQRSTUVWXYZ123456789")

// GenerateShortID generates a random short ID of the specified length
func GenerateShortID(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive, got %d", length)
	}

	b := make([]rune, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(lettersRune))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		b[i] = lettersRune[num.Int64()]
	}
	return string(b), nil
}
