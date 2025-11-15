// Package storage provides PostgreSQL-backed storage for URL mappings.
package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresURLStore manages URL storage with PostgreSQL
type PostgresURLStore struct {
	pool *pgxpool.Pool
}

// NewPostgresURLStore creates a new PostgresURLStore with connection pooling
func NewPostgresURLStore(ctx context.Context, databaseURL string) (*PostgresURLStore, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresURLStore{pool: pool}, nil
}

// Set stores a URL with the given short ID
func (s *PostgresURLStore) Set(shortID, url string) error {
	if shortID == "" {
		return fmt.Errorf("short ID cannot be empty")
	}
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO urls (short_id, encrypted_url, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (short_id) DO UPDATE
		SET encrypted_url = EXCLUDED.encrypted_url
	`

	_, err := s.pool.Exec(ctx, query, shortID, url, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert URL: %w", err)
	}

	return nil
}

// Get retrieves a URL by its short ID and updates access statistics
func (s *PostgresURLStore) Get(shortID string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use a transaction to get URL and update statistics atomically
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", false
	}
	defer tx.Rollback(ctx)

	// Get the URL
	var encryptedURL string
	query := `SELECT encrypted_url FROM urls WHERE short_id = $1`
	err = tx.QueryRow(ctx, query, shortID).Scan(&encryptedURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false
		}
		return "", false
	}

	// Update access statistics
	updateQuery := `
		UPDATE urls
		SET click_count = click_count + 1,
		    last_accessed = $1
		WHERE short_id = $2
	`
	_, err = tx.Exec(ctx, updateQuery, time.Now(), shortID)
	if err != nil {
		// If update fails, still return the URL but don't commit stats
		return encryptedURL, true
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return encryptedURL, true
	}

	return encryptedURL, true
}

// Delete removes a URL by its short ID
func (s *PostgresURLStore) Delete(shortID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM urls WHERE short_id = $1`
	s.pool.Exec(ctx, query, shortID)
}

// Exists checks if a short ID exists
func (s *PostgresURLStore) Exists(shortID string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_id = $1)`
	err := s.pool.QueryRow(ctx, query, shortID).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

// Count returns the number of stored URLs
func (s *PostgresURLStore) Count() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	query := `SELECT COUNT(*) FROM urls`
	err := s.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0
	}

	return count
}

// Close closes the database connection pool
func (s *PostgresURLStore) Close() {
	s.pool.Close()
}

// GetStats returns statistics for a specific short URL
func (s *PostgresURLStore) GetStats(shortID string) (clickCount int64, lastAccessed *time.Time, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT click_count, last_accessed FROM urls WHERE short_id = $1`
	err = s.pool.QueryRow(ctx, query, shortID).Scan(&clickCount, &lastAccessed)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil, fmt.Errorf("URL not found")
		}
		return 0, nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return clickCount, lastAccessed, nil
}

// CleanExpiredURLs removes URLs that have passed their expiration time
func (s *PostgresURLStore) CleanExpiredURLs() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `DELETE FROM urls WHERE expires_at IS NOT NULL AND expires_at < $1`
	result, err := s.pool.Exec(ctx, query, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to clean expired URLs: %w", err)
	}

	return result.RowsAffected(), nil
}

// Ping checks if the database connection is alive
func (s *PostgresURLStore) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.pool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
