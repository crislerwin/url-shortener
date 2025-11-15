package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crislerwin/url-shortener/internal/crypto"
	"github.com/crislerwin/url-shortener/internal/storage"
)

func setupHandler(t *testing.T) *Handler {
	t.Helper()
	secretKey := "12345678901234567890123456789012"
	encryptor, err := crypto.NewEncryptor(secretKey)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	store := storage.NewURLStore()
	return NewHandler(store, encryptor, "http://localhost:8080")
}

func TestNewHandler(t *testing.T) {
	handler := setupHandler(t)
	if handler == nil {
		t.Fatal("NewHandler() returned nil")
	}
	if handler.store == nil {
		t.Error("Handler store is nil")
	}
	if handler.encryptor == nil {
		t.Error("Handler encryptor is nil")
	}
	if handler.baseURL == "" {
		t.Error("Handler baseURL is empty")
	}
}

func TestHandler_ShortenURL(t *testing.T) {
	handler := setupHandler(t)

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "valid URL",
			url:            "https://example.com",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "valid URL with path",
			url:            "https://example.com/path/to/page",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "valid URL with query params",
			url:            "https://example.com?param1=value1&param2=value2",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "missing URL parameter",
			url:            "",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlParam := ""
			if tt.url != "" {
				urlParam = "?url=" + tt.url
			}
			req := httptest.NewRequest(http.MethodGet, "/shorten"+urlParam, nil)
			w := httptest.NewRecorder()

			handler.ShortenURL(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("ShortenURL() status = %d, want %d", resp.StatusCode, tt.expectedStatus)
			}

			if !tt.expectError && resp.StatusCode == http.StatusOK {
				// Check that response contains shortened URL
				body := w.Body.String()
				if body == "" {
					t.Error("ShortenURL() returned empty body")
				}
				if len(body) < len("The shortened URL of this original URL is: http://localhost:8080/") {
					t.Errorf("ShortenURL() response too short: %s", body)
				}
			}
		})
	}
}

func TestHandler_RedirectHandler(t *testing.T) {
	handler := setupHandler(t)

	// First, create a shortened URL
	originalURL := "https://example.com/test"
	encrypted, _ := handler.encryptor.Encrypt(originalURL)
	shortID := "test123"
	_ = handler.store.Set(shortID, encrypted)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectRedirect bool
	}{
		{
			name:           "valid short ID",
			path:           "/" + shortID,
			expectedStatus: http.StatusFound,
			expectRedirect: true,
		},
		{
			name:           "non-existent short ID",
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
			expectRedirect: false,
		},
		{
			name:           "empty path",
			path:           "/",
			expectedStatus: http.StatusBadRequest,
			expectRedirect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.RedirectHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("RedirectHandler() status = %d, want %d", resp.StatusCode, tt.expectedStatus)
			}

			if tt.expectRedirect {
				location := resp.Header.Get("Location")
				if location != originalURL {
					t.Errorf("RedirectHandler() Location = %v, want %v", location, originalURL)
				}
			}
		})
	}
}

func TestHandler_ShortenURLAndRedirect(t *testing.T) {
	handler := setupHandler(t)

	originalURL := "https://example.com/integration-test"

	// Step 1: Shorten the URL
	req := httptest.NewRequest(http.MethodGet, "/shorten?url="+originalURL, nil)
	w := httptest.NewRecorder()
	handler.ShortenURL(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ShortenURL() failed with status %d", w.Code)
	}

	// Extract short ID from response
	body := w.Body.String()
	expectedPrefix := "The shortened URL of this original URL is: http://localhost:8080/"
	if len(body) <= len(expectedPrefix) {
		t.Fatalf("ShortenURL() response too short: %s", body)
	}
	shortID := body[len(expectedPrefix):]

	// Step 2: Use the short ID to redirect
	req2 := httptest.NewRequest(http.MethodGet, "/"+shortID, nil)
	w2 := httptest.NewRecorder()
	handler.RedirectHandler(w2, req2)

	resp := w2.Result()
	if resp.StatusCode != http.StatusFound {
		t.Errorf("RedirectHandler() status = %d, want %d", resp.StatusCode, http.StatusFound)
	}

	location := resp.Header.Get("Location")
	if location != originalURL {
		t.Errorf("RedirectHandler() Location = %v, want %v", location, originalURL)
	}
}

func TestHandler_RedirectHandlerWithInvalidEncryption(t *testing.T) {
	handler := setupHandler(t)

	// Store an invalid encrypted value
	shortID := "invalid123"
	_ = handler.store.Set(shortID, "invalid-encrypted-data")

	req := httptest.NewRequest(http.MethodGet, "/"+shortID, nil)
	w := httptest.NewRecorder()

	handler.RedirectHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("RedirectHandler() with invalid encryption status = %d, want %d",
			resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestHandler_MultipleShortenRequests(t *testing.T) {
	handler := setupHandler(t)

	urls := []string{
		"https://example1.com",
		"https://example2.com",
		"https://example3.com",
	}

	shortIDs := make([]string, 0, len(urls))

	// Shorten all URLs
	for _, url := range urls {
		req := httptest.NewRequest(http.MethodGet, "/shorten?url="+url, nil)
		w := httptest.NewRecorder()
		handler.ShortenURL(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("ShortenURL() failed for %s with status %d", url, w.Code)
		}

		body := w.Body.String()
		expectedPrefix := "The shortened URL of this original URL is: http://localhost:8080/"
		shortID := body[len(expectedPrefix):]
		shortIDs = append(shortIDs, shortID)
	}

	// Verify all shortened URLs redirect correctly
	for i, shortID := range shortIDs {
		req := httptest.NewRequest(http.MethodGet, "/"+shortID, nil)
		w := httptest.NewRecorder()
		handler.RedirectHandler(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusFound {
			t.Errorf("RedirectHandler() status = %d, want %d", resp.StatusCode, http.StatusFound)
		}

		location := resp.Header.Get("Location")
		if location != urls[i] {
			t.Errorf("RedirectHandler() Location = %v, want %v", location, urls[i])
		}
	}
}

func TestHandler_ShortenURLWithDifferentBaseURL(t *testing.T) {
	secretKey := "12345678901234567890123456789012"
	encryptor, _ := crypto.NewEncryptor(secretKey)
	store := storage.NewURLStore()
	baseURL := "https://short.link"
	handler := NewHandler(store, encryptor, baseURL)

	originalURL := "https://example.com"
	req := httptest.NewRequest(http.MethodGet, "/shorten?url="+originalURL, nil)
	w := httptest.NewRecorder()

	handler.ShortenURL(w, req)

	body := w.Body.String()
	expectedPrefix := fmt.Sprintf("The shortened URL of this original URL is: %s/", baseURL)

	if len(body) <= len(expectedPrefix) {
		t.Fatalf("ShortenURL() response doesn't contain baseURL: %s", body)
	}

	if body[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("ShortenURL() doesn't use correct baseURL, got: %s", body)
	}
}

func TestHandler_ConcurrentRequests(t *testing.T) {
	handler := setupHandler(t)

	numRequests := 50
	done := make(chan bool, numRequests)

	// Concurrent shorten requests
	for i := 0; i < numRequests; i++ {
		go func(id int) {
			url := fmt.Sprintf("https://example.com/page%d", id)
			req := httptest.NewRequest(http.MethodGet, "/shorten?url="+url, nil)
			w := httptest.NewRecorder()
			handler.ShortenURL(w, req)
			done <- w.Code == http.StatusOK
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		success := <-done
		if !success {
			t.Error("Concurrent request failed")
		}
	}
}
