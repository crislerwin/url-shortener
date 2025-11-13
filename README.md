# URL Shortener

A simple URL shortening service built with Go that encrypts stored URLs using AES-256 encryption.

## Features

- Shorten long URLs to 6-character alphanumeric short codes
- AES-256 CTR encryption for stored URLs
- In-memory storage with thread-safe operations
- Simple HTTP API
- Automatic redirection from short URLs to original URLs

## Prerequisites

- Go 1.25.4 or higher

## Installation

Clone the repository:

```bash
git clone https://github.com/crislerwin/url-shortener.git
cd url-shortener
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
The shorted URL of this original URL is: http://localhost:8080/abc123
```

#### Redirect to Original URL

**Endpoint:** `GET /<short-id>`

**Example:**
```bash
curl "http://localhost:8080/abc123"
```

This will redirect you to the original URL.

## How It Works

1. When you submit a URL to shorten, the service:
   - Encrypts the original URL using AES-256 CTR mode
   - Generates a random 6-character alphanumeric short ID
   - Stores the mapping in memory
   - Returns the shortened URL

2. When accessing a short URL:
   - The service looks up the encrypted URL by short ID
   - Decrypts the URL
   - Redirects to the original URL

## Security Note

The current implementation uses a hardcoded encryption key for demonstration purposes. For production use, you should:
- Use environment variables for the encryption key
- Implement persistent storage (database)
- Add proper authentication and rate limiting
- Use HTTPS

## License

MIT
