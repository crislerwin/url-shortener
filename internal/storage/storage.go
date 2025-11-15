// Package storage provides storage implementations for URL mappings.
//
// This package defines the Storage interface and provides both in-memory
// and PostgreSQL implementations for storing URL mappings.
//
// Example usage:
//
//	// In-memory storage
//	store := storage.NewURLStore()
//
//	// PostgreSQL storage
//	store, err := storage.NewPostgresURLStore(ctx, databaseURL)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Both implement the same interface
//	err = store.Set("abc123", "encrypted-url-data")
//	url, ok := store.Get("abc123")
package storage

import (
	"fmt"
	"sync"
)

// Storage defines the interface for URL storage implementations
type Storage interface {
	Set(shortID, url string) error
	Get(shortID string) (string, bool)
	Delete(shortID string)
	Exists(shortID string) bool
	Count() int
}

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
