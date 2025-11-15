// Package storage provides thread-safe in-memory storage for URL mappings.
//
// This package implements a concurrent-safe key-value store that maps short IDs
// to encrypted URLs. It uses sync.RWMutex to allow multiple concurrent readers
// while ensuring exclusive access for write operations.
//
// Example usage:
//
//	store := storage.NewURLStore()
//
//	// Store a URL
//	err := store.Set("abc123", "encrypted-url-data")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Retrieve a URL
//	url, ok := store.Get("abc123")
//	if !ok {
//	    log.Println("URL not found")
//	}
//
//	// Check if URL exists
//	if store.Exists("abc123") {
//	    fmt.Println("URL exists")
//	}
//
//	// Delete a URL
//	store.Delete("abc123")
//
//	// Get total count
//	count := store.Count()
package storage

import (
	"fmt"
	"sync"
)

// URLStore manages URL storage with thread-safe operations
type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex
}

// NewURLStore creates a new URLStore
func NewURLStore() *URLStore {
	return &URLStore{
		urls: make(map[string]string),
	}
}

// Set stores a URL with the given short ID
func (s *URLStore) Set(shortID, url string) error {
	if shortID == "" {
		return fmt.Errorf("short ID cannot be empty")
	}
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.urls[shortID] = url
	return nil
}

// Get retrieves a URL by its short ID
func (s *URLStore) Get(shortID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.urls[shortID]
	return url, ok
}

// Delete removes a URL by its short ID
func (s *URLStore) Delete(shortID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.urls, shortID)
}

// Exists checks if a short ID exists
func (s *URLStore) Exists(shortID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.urls[shortID]
	return ok
}

// Count returns the number of stored URLs
func (s *URLStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.urls)
}
