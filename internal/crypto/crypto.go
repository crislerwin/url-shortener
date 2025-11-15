// Package crypto provides AES-256 encryption and decryption functionality for URL shortener.
//
// This package implements secure URL encryption using AES-256 in CTR (Counter) mode.
// Each encryption operation generates a unique initialization vector (IV) to ensure
// that encrypting the same URL multiple times produces different ciphertext outputs.
//
// Example usage:
//
//	encryptor, err := crypto.NewEncryptor("12345678901234567890123456789012") // 32-byte key
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	encrypted, err := encryptor.Encrypt("https://example.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	decrypted, err := encryptor.Decrypt(encrypted)
//	if err != nil {
//	    log.Fatal(err)
//	}
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Encryptor handles URL encryption and decryption
type Encryptor struct {
	secretKey []byte
}

// NewEncryptor creates a new Encryptor with the provided secret key
func NewEncryptor(secretKey string) (*Encryptor, error) {
	if len(secretKey) != 32 {
		return nil, fmt.Errorf("secret key must be 32 bytes for AES-256, got %d bytes", len(secretKey))
	}
	return &Encryptor{
		secretKey: []byte(secretKey),
	}, nil
}

// Decrypt decrypts an encrypted URL string
func (e *Encryptor) Decrypt(encryptedUrl string) (string, error) {
	block, err := aes.NewCipher(e.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	cipherText, err := hex.DecodeString(encryptedUrl)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %w", err)
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("cipher text too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// Encrypt encrypts a URL string
func (e *Encryptor) Encrypt(originalUrl string) (string, error) {
	block, err := aes.NewCipher(e.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	plainText := []byte(originalUrl)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %w", err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return hex.EncodeToString(cipherText), nil
}
