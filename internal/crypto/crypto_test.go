package crypto

import (
	"strings"
	"testing"
)

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name      string
		secretKey string
		wantErr   bool
	}{
		{
			name:      "valid 32-byte key",
			secretKey: "12345678901234567890123456789012",
			wantErr:   false,
		},
		{
			name:      "invalid short key",
			secretKey: "short",
			wantErr:   true,
		},
		{
			name:      "invalid long key",
			secretKey: "123456789012345678901234567890123",
			wantErr:   true,
		},
		{
			name:      "empty key",
			secretKey: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := NewEncryptor(tt.secretKey)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewEncryptor() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("NewEncryptor() unexpected error: %v", err)
				}
				if enc == nil {
					t.Errorf("NewEncryptor() returned nil encryptor")
				}
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	secretKey := "12345678901234567890123456789012"
	enc, err := NewEncryptor(secretKey)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	tests := []struct {
		name        string
		originalURL string
	}{
		{
			name:        "simple URL",
			originalURL: "https://example.com",
		},
		{
			name:        "URL with query params",
			originalURL: "https://example.com/page?param1=value1&param2=value2",
		},
		{
			name:        "long URL",
			originalURL: "https://example.com/very/long/path/with/many/segments/and/parameters?key1=value1&key2=value2&key3=value3",
		},
		{
			name:        "URL with special characters",
			originalURL: "https://example.com/path?query=hello%20world&special=!@#$%",
		},
		{
			name:        "empty string",
			originalURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := enc.Encrypt(tt.originalURL)
			if err != nil {
				t.Fatalf("Encrypt() error: %v", err)
			}

			if encrypted == "" {
				t.Errorf("Encrypt() returned empty string")
			}

			decrypted, err := enc.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt() error: %v", err)
			}

			if decrypted != tt.originalURL {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.originalURL)
			}
		})
	}
}

func TestEncryptDifferentOutputs(t *testing.T) {
	secretKey := "12345678901234567890123456789012"
	enc, err := NewEncryptor(secretKey)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	url := "https://example.com"
	encrypted1, err := enc.Encrypt(url)
	if err != nil {
		t.Fatalf("First Encrypt() error: %v", err)
	}

	encrypted2, err := enc.Encrypt(url)
	if err != nil {
		t.Fatalf("Second Encrypt() error: %v", err)
	}

	// Encrypted outputs should be different due to random IV
	if encrypted1 == encrypted2 {
		t.Errorf("Encrypt() produced same output twice, expected different due to random IV")
	}

	// But both should decrypt to the same original URL
	decrypted1, _ := enc.Decrypt(encrypted1)
	decrypted2, _ := enc.Decrypt(encrypted2)

	if decrypted1 != url || decrypted2 != url {
		t.Errorf("Decrypted URLs don't match original")
	}
}

func TestDecryptInvalidInput(t *testing.T) {
	secretKey := "12345678901234567890123456789012"
	enc, err := NewEncryptor(secretKey)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	tests := []struct {
		name      string
		encrypted string
		wantErr   bool
	}{
		{
			name:      "invalid hex string",
			encrypted: "not-a-hex-string",
			wantErr:   true,
		},
		{
			name:      "too short ciphertext",
			encrypted: "abcd",
			wantErr:   true,
		},
		{
			name:      "empty string",
			encrypted: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := enc.Decrypt(tt.encrypted)
			if tt.wantErr && err == nil {
				t.Errorf("Decrypt() expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Decrypt() unexpected error: %v", err)
			}
		})
	}
}

func TestEncryptorWithDifferentKeys(t *testing.T) {
	key1 := "12345678901234567890123456789012"
	key2 := "abcdefghijklmnopqrstuvwxyz123456"

	enc1, _ := NewEncryptor(key1)
	enc2, _ := NewEncryptor(key2)

	url := "https://example.com"

	encrypted, err := enc1.Encrypt(url)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	// Trying to decrypt with wrong key should fail or produce garbage
	decrypted, err := enc2.Decrypt(encrypted)
	if err == nil && decrypted == url {
		t.Errorf("Decrypt() with wrong key should not produce original URL")
	}
}

func TestEncryptOutputFormat(t *testing.T) {
	secretKey := "12345678901234567890123456789012"
	enc, err := NewEncryptor(secretKey)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	url := "https://example.com"
	encrypted, err := enc.Encrypt(url)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	// Check that output is hex encoded (only contains 0-9, a-f)
	validHex := true
	for _, c := range strings.ToLower(encrypted) {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			validHex = false
			break
		}
	}

	if !validHex {
		t.Errorf("Encrypt() output is not valid hex string")
	}

	// Length should be even (hex encoding)
	if len(encrypted)%2 != 0 {
		t.Errorf("Encrypt() output length should be even, got %d", len(encrypted))
	}
}
