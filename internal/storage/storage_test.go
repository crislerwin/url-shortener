package storage

import (
	"sync"
	"testing"
)

func TestNewURLStore(t *testing.T) {
	store := NewURLStore()
	if store == nil {
		t.Fatal("NewURLStore() returned nil")
	}
	if store.urls == nil {
		t.Error("NewURLStore() did not initialize urls map")
	}
	if store.Count() != 0 {
		t.Errorf("NewURLStore() Count() = %d, want 0", store.Count())
	}
}

func TestURLStore_Set(t *testing.T) {
	tests := []struct {
		name    string
		shortID string
		url     string
		wantErr bool
	}{
		{
			name:    "valid set",
			shortID: "abc123",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "empty shortID",
			shortID: "",
			url:     "https://example.com",
			wantErr: true,
		},
		{
			name:    "empty URL",
			shortID: "abc123",
			url:     "",
			wantErr: true,
		},
		{
			name:    "both empty",
			shortID: "",
			url:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewURLStore()
			err := store.Set(tt.shortID, tt.url)
			if tt.wantErr && err == nil {
				t.Errorf("Set() expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Set() unexpected error: %v", err)
			}
		})
	}
}

func TestURLStore_Get(t *testing.T) {
	store := NewURLStore()
	shortID := "abc123"
	url := "https://example.com"

	// Test getting non-existent key
	_, ok := store.Get(shortID)
	if ok {
		t.Error("Get() returned ok=true for non-existent key")
	}

	// Set and then get
	err := store.Set(shortID, url)
	if err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	gotURL, ok := store.Get(shortID)
	if !ok {
		t.Error("Get() returned ok=false for existing key")
	}
	if gotURL != url {
		t.Errorf("Get() = %v, want %v", gotURL, url)
	}
}

func TestURLStore_Delete(t *testing.T) {
	store := NewURLStore()
	shortID := "abc123"
	url := "https://example.com"

	// Set a value
	err := store.Set(shortID, url)
	if err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	// Verify it exists
	if !store.Exists(shortID) {
		t.Error("Exists() returned false for existing key")
	}

	// Delete it
	store.Delete(shortID)

	// Verify it's gone
	if store.Exists(shortID) {
		t.Error("Exists() returned true for deleted key")
	}

	_, ok := store.Get(shortID)
	if ok {
		t.Error("Get() returned ok=true for deleted key")
	}
}

func TestURLStore_Exists(t *testing.T) {
	store := NewURLStore()
	shortID := "abc123"
	url := "https://example.com"

	if store.Exists(shortID) {
		t.Error("Exists() returned true for non-existent key")
	}

	err := store.Set(shortID, url)
	if err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	if !store.Exists(shortID) {
		t.Error("Exists() returned false for existing key")
	}
}

func TestURLStore_Count(t *testing.T) {
	store := NewURLStore()

	if store.Count() != 0 {
		t.Errorf("Count() = %d, want 0", store.Count())
	}

	// Add some URLs
	urls := map[string]string{
		"abc123": "https://example.com",
		"def456": "https://google.com",
		"ghi789": "https://github.com",
	}

	for id, url := range urls {
		err := store.Set(id, url)
		if err != nil {
			t.Fatalf("Set() error: %v", err)
		}
	}

	if store.Count() != len(urls) {
		t.Errorf("Count() = %d, want %d", store.Count(), len(urls))
	}

	// Delete one
	store.Delete("abc123")

	if store.Count() != len(urls)-1 {
		t.Errorf("Count() after delete = %d, want %d", store.Count(), len(urls)-1)
	}
}

func TestURLStore_Update(t *testing.T) {
	store := NewURLStore()
	shortID := "abc123"
	url1 := "https://example.com"
	url2 := "https://newexample.com"

	// Set initial value
	err := store.Set(shortID, url1)
	if err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	// Update with new value
	err = store.Set(shortID, url2)
	if err != nil {
		t.Fatalf("Set() update error: %v", err)
	}

	// Verify new value
	gotURL, ok := store.Get(shortID)
	if !ok {
		t.Error("Get() returned ok=false")
	}
	if gotURL != url2 {
		t.Errorf("Get() = %v, want %v", gotURL, url2)
	}

	// Count should still be 1
	if store.Count() != 1 {
		t.Errorf("Count() = %d, want 1", store.Count())
	}
}

func TestURLStore_ConcurrentAccess(t *testing.T) {
	store := NewURLStore()
	numGoroutines := 100
	numOperations := 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // For Set, Get, and Delete operations

	// Concurrent Set operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				shortID := string(rune('a' + (id % 26)))
				url := "https://example.com"
				_ = store.Set(shortID, url)
			}
		}(i)
	}

	// Concurrent Get operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				shortID := string(rune('a' + (id % 26)))
				_, _ = store.Get(shortID)
			}
		}(i)
	}

	// Concurrent Delete operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				shortID := string(rune('a' + (id % 26)))
				store.Delete(shortID)
			}
		}(i)
	}

	wg.Wait()

	// Test completed without race conditions
	// The actual count doesn't matter, just that we didn't panic
	_ = store.Count()
}

func TestURLStore_MultipleEntries(t *testing.T) {
	store := NewURLStore()

	urls := map[string]string{
		"short1": "https://example1.com",
		"short2": "https://example2.com",
		"short3": "https://example3.com",
		"short4": "https://example4.com",
		"short5": "https://example5.com",
	}

	// Set all URLs
	for id, url := range urls {
		err := store.Set(id, url)
		if err != nil {
			t.Fatalf("Set() error for %s: %v", id, err)
		}
	}

	// Verify all URLs
	for id, expectedURL := range urls {
		gotURL, ok := store.Get(id)
		if !ok {
			t.Errorf("Get(%s) returned ok=false", id)
		}
		if gotURL != expectedURL {
			t.Errorf("Get(%s) = %v, want %v", id, gotURL, expectedURL)
		}
	}

	// Verify count
	if store.Count() != len(urls) {
		t.Errorf("Count() = %d, want %d", store.Count(), len(urls))
	}
}
