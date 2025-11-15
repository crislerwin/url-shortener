# URL Shortener

A simple URL shortening service built with Go that encrypts stored URLs using AES-256 encryption.

## Features

- Shorten long URLs to 6-character alphanumeric short codes
- AES-256 CTR encryption for stored URLs
- In-memory storage with thread-safe operations
- Simple HTTP API
- Automatic redirection from short URLs to original URLs
- Modular architecture with clean separation of concerns
- Comprehensive test coverage (34+ tests)

## Prerequisites

- Go 1.25.4 or higher

## Installation

Clone the repository:

```bash
git clone https://github.com/crislerwin/url-shortener.git
cd url-shortener
```

## Project Structure

```
url-shortener/
├── main.go                    # Application entry point
├── internal/
│   ├── crypto/               # Encryption/decryption logic
│   │   ├── crypto.go
│   │   └── crypto_test.go
│   ├── handlers/             # HTTP request handlers
│   │   ├── handlers.go
│   │   └── handlers_test.go
│   ├── storage/              # Thread-safe URL storage
│   │   ├── storage.go
│   │   └── storage_test.go
│   └── utils/                # Utility functions (ID generation)
│       ├── generator.go
│       └── generator_test.go
├── go.mod
└── README.md
```

## Usage

### Starting the Server

Run the application:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

### API Endpoints

#### Shorten a URL

**Endpoint:** `GET /shorten?url=<your-url>`

**Example:**
```bash
curl "http://localhost:8080/shorten?url=https://www.example.com/very/long/url"
```

**Response:**
```
The shortened URL of this original URL is: http://localhost:8080/abc123
```

#### Redirect to Original URL

**Endpoint:** `GET /<short-id>`

**Example:**
```bash
curl -L "http://localhost:8080/abc123"
```

This will redirect you to the original URL.

## How It Works

1. **URL Shortening Process:**
   - Original URL is encrypted using AES-256 CTR mode
   - A cryptographically secure random 6-character alphanumeric short ID is generated
   - The encrypted URL is stored in-memory mapped to the short ID
   - A shortened URL is returned to the user

2. **URL Redirection Process:**
   - The service looks up the encrypted URL by short ID
   - Decrypts the URL using the secret key
   - Issues an HTTP 302 redirect to the original URL

3. **Security Features:**
   - All URLs are encrypted at rest using AES-256
   - Thread-safe concurrent access to URL store
   - Cryptographically secure random ID generation
   - Proper error handling throughout the application

## Development

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test ./... -v
```

Run tests with coverage:
```bash
go test ./... -cover
```

Run tests with coverage report:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Run benchmarks:
```bash
go test ./internal/utils -bench=. -benchmem
```

### Package Overview

#### `internal/crypto`
Handles AES-256 encryption and decryption of URLs.

**Key Features:**
- Encryptor struct with dependency injection
- Validates secret key length (must be 32 bytes for AES-256)
- Returns errors instead of panicking
- Unique IV (Initialization Vector) for each encryption

#### `internal/storage`
Thread-safe in-memory storage for URL mappings.

**Key Features:**
- RWMutex for concurrent read/write access
- CRUD operations (Set, Get, Delete, Exists, Count)
- Input validation
- Safe for concurrent use

#### `internal/handlers`
HTTP request handlers for the API.

**Key Features:**
- Dependency injection pattern
- Configurable base URL
- Proper HTTP status codes
- Error handling and validation

#### `internal/utils`
Utility functions for the application.

**Key Features:**
- Cryptographically secure random ID generation
- Configurable ID length
- Alphanumeric character set (a-z, A-Z, 0-9)

## Architecture Decisions

### Why Internal Package?
The `internal/` directory prevents external packages from importing these modules, enforcing encapsulation.

### Why Dependency Injection?
All packages use dependency injection for better testability and flexibility. This allows easy mocking in tests and configuration changes without modifying code.

### Why In-Memory Storage?
For simplicity and performance. In production, you would replace this with a database (Redis, PostgreSQL, etc.) by implementing the same interface.

### Why Return Errors Instead of Panicking?
Following Go best practices, all functions return errors instead of using `log.Fatal` or `panic`, making the code more robust and testable.

## Testing

The project includes comprehensive test coverage:

- **Crypto Package**: 6 test functions covering encryption, decryption, error cases
- **Storage Package**: 8 test functions including concurrent access tests
- **Handlers Package**: 8 test functions including integration and concurrent request tests
- **Utils Package**: 8 test functions + 3 benchmarks covering randomness and uniqueness

All tests use table-driven tests for better coverage and maintainability.

## Security Notes

The current implementation uses a hardcoded encryption key for demonstration purposes. 

### For Production Use:

1. **Environment Variables:**
   ```go
   secretKey := os.Getenv("URL_SHORTENER_SECRET_KEY")
   if len(secretKey) != 32 {
       log.Fatal("SECRET_KEY must be 32 bytes")
   }
   ```

2. **Persistent Storage:**
   - Replace in-memory storage with Redis or PostgreSQL
   - Implement proper data persistence
   - Add database migration support

3. **Additional Security:**
   - Implement rate limiting
   - Add authentication/authorization
   - Use HTTPS in production
   - Add request validation and sanitization
   - Implement logging and monitoring
   - Add expiration for shortened URLs

4. **Infrastructure:**
   - Use reverse proxy (nginx)
   - Add health check endpoints
   - Implement graceful shutdown
   - Add metrics and observability

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

MIT
