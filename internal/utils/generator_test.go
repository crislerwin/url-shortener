package utils

import (
	"testing"
	"unicode"
)

func TestGenerateShortID(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "length 6",
			length:  6,
			wantErr: false,
		},
		{
			name:    "length 1",
			length:  1,
			wantErr: false,
		},
		{
			name:    "length 10",
			length:  10,
			wantErr: false,
		},
		{
			name:    "length 100",
			length:  100,
			wantErr: false,
		},
		{
			name:    "zero length",
			length:  0,
			wantErr: true,
		},
		{
			name:    "negative length",
			length:  -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GenerateShortID(tt.length)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GenerateShortID() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateShortID() unexpected error: %v", err)
			}

			if len(id) != tt.length {
				t.Errorf("GenerateShortID() length = %d, want %d", len(id), tt.length)
			}
		})
	}
}

func TestGenerateShortID_ValidCharacters(t *testing.T) {
	length := 100
	id, err := GenerateShortID(length)
	if err != nil {
		t.Fatalf("GenerateShortID() error: %v", err)
	}

	validChars := "abcdefghijlmnopqrstuvwxyzABCDEFGHIJLMNOPQRSTUVWXYZ123456789"
	validCharMap := make(map[rune]bool)
	for _, c := range validChars {
		validCharMap[c] = true
	}

	for i, c := range id {
		if !validCharMap[c] {
			t.Errorf("GenerateShortID() contains invalid character '%c' at position %d", c, i)
		}
	}
}

func TestGenerateShortID_OnlyAlphanumeric(t *testing.T) {
	length := 50
	id, err := GenerateShortID(length)
	if err != nil {
		t.Fatalf("GenerateShortID() error: %v", err)
	}

	for i, c := range id {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			t.Errorf("GenerateShortID() contains non-alphanumeric character '%c' at position %d", c, i)
		}
	}
}

func TestGenerateShortID_Uniqueness(t *testing.T) {
	length := 6
	numIDs := 1000
	ids := make(map[string]bool)

	for i := 0; i < numIDs; i++ {
		id, err := GenerateShortID(length)
		if err != nil {
			t.Fatalf("GenerateShortID() error on iteration %d: %v", i, err)
		}
		ids[id] = true
	}

	// With 6 characters from 61 possible chars, collision is very unlikely
	// We expect most IDs to be unique
	uniqueCount := len(ids)
	expectedMinUnique := int(float64(numIDs) * 0.95) // At least 95% should be unique

	if uniqueCount < expectedMinUnique {
		t.Errorf("GenerateShortID() uniqueness too low: got %d unique IDs out of %d, want at least %d",
			uniqueCount, numIDs, expectedMinUnique)
	}
}

func TestGenerateShortID_Randomness(t *testing.T) {
	length := 6
	id1, err1 := GenerateShortID(length)
	id2, err2 := GenerateShortID(length)

	if err1 != nil || err2 != nil {
		t.Fatalf("GenerateShortID() errors: %v, %v", err1, err2)
	}

	// Two consecutive calls should produce different results (with very high probability)
	if id1 == id2 {
		// Try one more time to be sure
		id3, _ := GenerateShortID(length)
		if id1 == id3 {
			t.Error("GenerateShortID() produced same result three times in a row, likely not random")
		}
	}
}

func TestGenerateShortID_CharacterDistribution(t *testing.T) {
	length := 1
	numSamples := 1000
	charCount := make(map[rune]int)

	for i := 0; i < numSamples; i++ {
		id, err := GenerateShortID(length)
		if err != nil {
			t.Fatalf("GenerateShortID() error: %v", err)
		}
		if len(id) != 1 {
			t.Fatalf("Expected length 1, got %d", len(id))
		}
		charCount[rune(id[0])]++
	}

	// With 61 possible characters and 1000 samples, we expect some distribution
	// At minimum, we should see more than one unique character
	if len(charCount) < 2 {
		t.Errorf("GenerateShortID() poor character distribution: only %d unique characters in %d samples",
			len(charCount), numSamples)
	}
}

func TestGenerateShortID_EmptyString(t *testing.T) {
	id, err := GenerateShortID(0)
	if err == nil {
		t.Error("GenerateShortID(0) expected error but got none")
	}
	if id != "" {
		t.Errorf("GenerateShortID(0) returned non-empty string: %s", id)
	}
}

func TestGenerateShortID_LargeLength(t *testing.T) {
	length := 10000
	id, err := GenerateShortID(length)
	if err != nil {
		t.Fatalf("GenerateShortID() error: %v", err)
	}

	if len(id) != length {
		t.Errorf("GenerateShortID() length = %d, want %d", len(id), length)
	}

	// Verify all characters are valid
	validChars := "abcdefghijlmnopqrstuvwxyzABCDEFGHIJLMNOPQRSTUVWXYZ123456789"
	validCharMap := make(map[rune]bool)
	for _, c := range validChars {
		validCharMap[c] = true
	}

	for i, c := range id {
		if !validCharMap[c] {
			t.Errorf("GenerateShortID() contains invalid character '%c' at position %d", c, i)
			break // Only report first invalid char
		}
	}
}

func BenchmarkGenerateShortID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateShortID(6)
	}
}

func BenchmarkGenerateShortID_Length10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateShortID(10)
	}
}

func BenchmarkGenerateShortID_Length100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateShortID(100)
	}
}
