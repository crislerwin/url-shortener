# URL Shortener

A simple URL shortening service built with Go that encrypts stored URLs using AES-256 encryption.

## Features

- Shorten long URLs to 6-character alphanumeric short codes
- AES-256 CTR encryption for stored URLs
- PostgreSQL database for durable, long-term storage
- Click tracking and analytics
- Simple HTTP API
- Automatic redirection from short URLs to original URLs
- Modular architecture with clean separation of concerns
- Comprehensive test coverage (34+ tests)
- Graceful shutdown with connection cleanup
- Environment-based configuration

## Prerequisites

- Go 1.25.4 or higher
- Docker and Docker Compose
- PostgreSQL 16 (via Docker)

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
├── docker-compose.yml         # Docker configuration for PostgreSQL
├── .env.example               # Environment variables template
├── .env                       # Environment variables (git-ignored)
├── migrations/                # Database migrations
│   └── 001_init.sql          # Initial schema
├── internal/
│   ├── crypto/               # Encryption/decryption logic
│   │   ├── crypto.go
│   │   └── crypto_test.go
│   ├── handlers/             # HTTP request handlers
│   │   ├── handlers.go
│   │   └── handlers_test.go
│   ├── storage/              # URL storage implementations
│   │   ├── storage.go        # In-memory storage (legacy)
│   │   ├── postgres.go       # PostgreSQL storage
│   │   └── storage_test.go
│   └── utils/                # Utility functions (ID generation)
│       ├── generator.go
│       └── generator_test.go
├── go.mod
└── README.md
```

## Quick Start

### 1. Setup Environment Variables

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` if you need to change any configuration (optional for development).

### 2. Start PostgreSQL with Docker

Start the PostgreSQL database:

```bash
docker-compose up -d
```

Verify the database is running:

```bash
docker-compose ps
```

You should see the `url-shortener-db` container running.

### 3. Run the Application

Run the application:

```bash
go run main.go
```

The server will start on `http://localhost:8080` and automatically connect to PostgreSQL.

### 4. Stop Services

To stop the application, press `Ctrl+C`.

To stop the database:

```bash
docker-compose down
```

To stop and remove all data (including stored URLs):

```bash
docker-compose down -v
```

## Usage

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
   - The encrypted URL is stored in PostgreSQL mapped to the short ID
   - A shortened URL is returned to the user

2. **URL Redirection Process:**
   - The service looks up the encrypted URL by short ID from PostgreSQL
   - Decrypts the URL using the secret key
   - Updates click count and last accessed timestamp
   - Issues an HTTP 302 redirect to the original URL

3. **Security & Persistence Features:**
   - All URLs are encrypted at rest using AES-256
   - PostgreSQL ensures ACID compliance and data durability
   - Connection pooling for efficient database access
   - Cryptographically secure random ID generation
   - Proper error handling throughout the application
   - Graceful shutdown with proper connection cleanup

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
PostgreSQL-backed storage for URL mappings with durability guarantees.

**Key Features:**
- PostgreSQL connection pooling for high performance
- CRUD operations (Set, Get, Delete, Exists, Count)
- Click tracking and analytics (GetStats)
- Automatic expiration cleanup (CleanExpiredURLs)
- Input validation
- Context-based timeouts for all operations
- ACID compliance for data durability

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

### Why PostgreSQL?
For data durability and long-term persistence. PostgreSQL provides:
- **ACID guarantees**: No data loss on crashes/restarts
- **Durability**: URLs are permanent assets that users expect to work forever
- **Performance**: 10K-50K reads/second with proper indexing
- **Analytics**: Built-in support for click tracking and reporting
- **Production-ready**: Battle-tested database used by major URL shorteners

### Why Connection Pooling?
Using `pgxpool` for efficient connection management:
- Reuses database connections instead of creating new ones
- Configurable pool size (5-25 connections)
- Health checks and automatic connection recovery
- Significant performance improvement under load

### Why Return Errors Instead of Panicking?
Following Go best practices, all functions return errors instead of using `log.Fatal` or `panic`, making the code more robust and testable.

## Testing

The project includes comprehensive test coverage:

- **Crypto Package**: 6 test functions covering encryption, decryption, error cases
- **Storage Package**: 8 test functions including concurrent access tests
- **Handlers Package**: 8 test functions including integration and concurrent request tests
- **Utils Package**: 8 test functions + 3 benchmarks covering randomness and uniqueness

All tests use table-driven tests for better coverage and maintainability.

## Environment Variables

The application uses the following environment variables (defined in `.env`):

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://urlshortener:urlshortener_dev_password@localhost:5432/urlshortener?sslmode=disable` |
| `SECRET_KEY` | AES-256 encryption key (must be 32 bytes) | `shhhhh_this_is_an_dumb_key123456` |
| `SERVER_PORT` | HTTP server port | `8080` |
| `BASE_URL` | Base URL for shortened links | `http://localhost:8080` |
| `POSTGRES_USER` | PostgreSQL username | `urlshortener` |
| `POSTGRES_PASSWORD` | PostgreSQL password | `urlshortener_dev_password` |
| `POSTGRES_DB` | PostgreSQL database name | `urlshortener` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |

## Database Management

### View Database Logs

```bash
docker-compose logs -f postgres
```

### Connect to PostgreSQL

Using `psql`:
```bash
docker exec -it url-shortener-db psql -U urlshortener -d urlshortener
```

Useful queries:
```sql
-- View all shortened URLs
SELECT short_id, created_at, click_count, last_accessed FROM urls;

-- Get total count
SELECT COUNT(*) FROM urls;

-- Find most clicked URLs
SELECT short_id, click_count FROM urls ORDER BY click_count DESC LIMIT 10;

-- Clean up old URLs (example: older than 1 year)
DELETE FROM urls WHERE created_at < NOW() - INTERVAL '1 year';
```

### Backup Database

```bash
docker exec url-shortener-db pg_dump -U urlshortener urlshortener > backup.sql
```

### Restore Database

```bash
cat backup.sql | docker exec -i url-shortener-db psql -U urlshortener urlshortener
```

## Security Notes

### Current Implementation
- Uses environment-based configuration
- AES-256 encryption for all stored URLs
- PostgreSQL with ACID guarantees
- Graceful shutdown with connection cleanup

### For Production Use:

1. **Encryption Key:**
   - Generate a secure 32-byte random key
   - Store in secure secret management (AWS Secrets Manager, HashiCorp Vault)
   - Never commit to version control
   ```bash
   # Generate a secure key
   openssl rand -base64 32
   ```

2. **Database Security:**
   - Use strong passwords (not the default dev password)
   - Enable SSL/TLS for database connections (`sslmode=require`)
   - Use managed PostgreSQL (AWS RDS, Google Cloud SQL)
   - Configure firewall rules to restrict database access
   - Regular automated backups

3. **Additional Security:**
   - Implement rate limiting (prevent abuse)
   - Add authentication/authorization for API endpoints
   - Use HTTPS in production (Let's Encrypt)
   - Add request validation and sanitization
   - Implement structured logging and monitoring
   - Add URL expiration support
   - Validate and sanitize input URLs
   - Implement CORS policies

4. **Infrastructure:**
   - Use reverse proxy (nginx, Caddy)
   - Add health check endpoints (`/health`, `/ready`)
   - Implement metrics and observability (Prometheus, Grafana)
   - Use containerization (Docker) for deployment
   - Set up CI/CD pipeline
   - Configure log aggregation
   - Add distributed tracing

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
