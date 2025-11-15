package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crislerwin/url-shortener/internal/crypto"
	"github.com/crislerwin/url-shortener/internal/handlers"
	"github.com/crislerwin/url-shortener/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get configuration from environment
	secretKey := getEnv("SECRET_KEY", "shhhhh_this_is_an_dumb_key123456") // 32 bytes for AES-256
	databaseURL := getEnv("DATABASE_URL", "postgres://urlshortener:urlshortener_dev_password@localhost:5432/urlshortener?sslmode=disable")
	serverPort := getEnv("SERVER_PORT", "8080")
	baseURL := getEnv("BASE_URL", "http://localhost:8080")

	// Validate secret key
	if len(secretKey) != 32 {
		log.Fatalf("SECRET_KEY must be exactly 32 bytes, got %d bytes", len(secretKey))
	}

	// Initialize encryptor
	encryptor, err := crypto.NewEncryptor(secretKey)
	if err != nil {
		log.Fatalf("Failed to create encryptor: %v", err)
	}

	// Initialize PostgreSQL storage
	ctx := context.Background()
	store, err := storage.NewPostgresURLStore(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer store.Close()

	log.Println("Successfully connected to PostgreSQL database")

	// Initialize handler
	handler := handlers.NewHandler(store, encryptor, baseURL)

	// Register routes
	http.HandleFunc("/shorten", handler.ShortenURL)
	http.HandleFunc("/", handler.RedirectHandler)

	// Setup HTTP server
	server := &http.Server{
		Addr:         ":" + serverPort,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("Server running on %s\n", baseURL)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
