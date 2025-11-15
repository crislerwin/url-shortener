// Package handlers provides HTTP request handlers for the URL shortener service.
//
// This package implements the HTTP handlers for URL shortening and redirection.
// It follows the dependency injection pattern, allowing for easy testing and
// configuration of storage, encryption, and base URL settings.
//
// Example usage:
//
//	// Initialize dependencies
//	store := storage.NewURLStore()
//	encryptor, _ := crypto.NewEncryptor("12345678901234567890123456789012")
//	handler := handlers.NewHandler(store, encryptor, "http://localhost:8080")
//
//	// Register routes
//	http.HandleFunc("/shorten", handler.ShortenURL)
//	http.HandleFunc("/", handler.RedirectHandler)
//
//	// Start server
//	http.ListenAndServe(":8080", nil)
//
// The package provides two main endpoints:
//   - /shorten?url=<original-url> - Creates a shortened URL
//   - /<short-id> - Redirects to the original URL
package handlers

import (
	"fmt"
	"net/http"

	"github.com/crislerwin/url-shortener/internal/crypto"
	"github.com/crislerwin/url-shortener/internal/metrics"
	"github.com/crislerwin/url-shortener/internal/storage"
	"github.com/crislerwin/url-shortener/internal/utils"
)

// Handler manages HTTP handlers with dependencies
type Handler struct {
	store     storage.Storage
	encryptor *crypto.Encryptor
	baseURL   string
}

// NewHandler creates a new Handler
func NewHandler(store storage.Storage, encryptor *crypto.Encryptor, baseURL string) *Handler {
	return &Handler{
		store:     store,
		encryptor: encryptor,
		baseURL:   baseURL,
	}
}

// RedirectHandler handles redirects from short URLs to original URLs
func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortID := r.URL.Path[1:]

	// Don't handle empty paths
	if shortID == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	encryptedURL, ok := h.store.Get(shortID)
	if !ok {
		metrics.URLRedirectsTotal.WithLabelValues("not_found").Inc()
		http.Error(w, "This URL doesn't exist", http.StatusNotFound)
		return
	}

	decryptedURL, err := h.encryptor.Decrypt(encryptedURL)
	if err != nil {
		metrics.URLRedirectsTotal.WithLabelValues("error").Inc()
		metrics.EncryptionOperationsTotal.WithLabelValues("decrypt", "error").Inc()
		http.Error(w, "Failed to decrypt URL", http.StatusInternalServerError)
		return
	}

	metrics.EncryptionOperationsTotal.WithLabelValues("decrypt", "success").Inc()
	metrics.URLRedirectsTotal.WithLabelValues("success").Inc()
	http.Redirect(w, r, decryptedURL, http.StatusFound)
}

// ShortenURL handles URL shortening requests
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	encryptedURL, err := h.encryptor.Encrypt(originalURL)
	if err != nil {
		metrics.EncryptionOperationsTotal.WithLabelValues("encrypt", "error").Inc()
		http.Error(w, "Failed to encrypt URL", http.StatusInternalServerError)
		return
	}
	metrics.EncryptionOperationsTotal.WithLabelValues("encrypt", "success").Inc()

	shortID, err := utils.GenerateShortID(6)
	if err != nil {
		http.Error(w, "Failed to generate short ID", http.StatusInternalServerError)
		return
	}

	if err := h.store.Set(shortID, encryptedURL); err != nil {
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}

	// Increment metrics
	metrics.URLsShortenedTotal.Inc()

	shortURL := fmt.Sprintf("%s/%s", h.baseURL, shortID)
	fmt.Fprintf(w, "The shortened URL of this original URL is: %s", shortURL)
}

// Health handles health check requests
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	// Check if storage is accessible
	if err := h.store.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Database health check failed: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
