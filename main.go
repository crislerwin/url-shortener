package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/crislerwin/url-shortener/internal/crypto"
	"github.com/crislerwin/url-shortener/internal/handlers"
	"github.com/crislerwin/url-shortener/internal/storage"
)

func main() {
	// Initialize dependencies
	secretKey := "shhhhh_this_is_an_dumb_key123456" // 32 bytes for AES-256

	encryptor, err := crypto.NewEncryptor(secretKey)
	if err != nil {
		log.Fatalf("Failed to create encryptor: %v", err)
	}

	store := storage.NewURLStore()
	handler := handlers.NewHandler(store, encryptor, "http://localhost:8080")

	// Register routes
	http.HandleFunc("/shorten", handler.ShortenURL)
	http.HandleFunc("/", handler.RedirectHandler)

	// Start server
	fmt.Println("Running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
