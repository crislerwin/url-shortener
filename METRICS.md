# Prometheus Metrics

This document describes the Prometheus metrics exposed by the URL Shortener service.

## Endpoint

The metrics are exposed at:
```
GET /metrics
```

## Available Metrics

### URL Operations

#### `url_shortener_urls_shortened_total`
- **Type**: Counter
- **Description**: The total number of URLs that have been shortened
- **Labels**: None

#### `url_shortener_redirects_total`
- **Type**: Counter
- **Description**: The total number of URL redirects
- **Labels**:
  - `status`: The result status (`success`, `not_found`, `error`)

#### `url_shortener_urls_stored`
- **Type**: Gauge
- **Description**: The current number of URLs stored in the database
- **Labels**: None

### HTTP Metrics

#### `url_shortener_http_request_duration_seconds`
- **Type**: Histogram
- **Description**: Duration of HTTP requests in seconds
- **Labels**:
  - `path`: The request path
  - `method`: The HTTP method (GET, POST, etc.)
  - `status`: The HTTP status code

#### `url_shortener_http_requests_total`
- **Type**: Counter
- **Description**: The total number of HTTP requests
- **Labels**:
  - `path`: The request path
  - `method`: The HTTP method
  - `status`: The HTTP status code

### Database Metrics

#### `url_shortener_db_operation_duration_seconds`
- **Type**: Histogram
- **Description**: Duration of database operations in seconds
- **Labels**:
  - `operation`: The database operation type (`get`, `set`, `delete`, `exists`, `count`)

#### `url_shortener_db_errors_total`
- **Type**: Counter
- **Description**: The total number of database errors
- **Labels**:
  - `operation`: The database operation type

### Encryption Metrics

#### `url_shortener_encryption_operations_total`
- **Type**: Counter
- **Description**: The total number of encryption/decryption operations
- **Labels**:
  - `operation`: The operation type (`encrypt`, `decrypt`)
  - `status`: The result status (`success`, `error`)

## Prometheus Configuration

To scrape these metrics, add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'url-shortener'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

## Example Queries

### Request Rate
```promql
rate(url_shortener_http_requests_total[5m])
```

### Error Rate
```promql
rate(url_shortener_redirects_total{status="error"}[5m])
```

### 95th Percentile Request Duration
```promql
histogram_quantile(0.95, rate(url_shortener_http_request_duration_seconds_bucket[5m]))
```

### Success Rate
```promql
sum(rate(url_shortener_redirects_total{status="success"}[5m])) / sum(rate(url_shortener_redirects_total[5m]))
```

### Database Error Rate
```promql
rate(url_shortener_db_errors_total[5m])
```

## Grafana Dashboard

The metrics can be visualized using Grafana. Import the following panels:

1. **Request Rate**: Total requests per second
2. **Success vs Error Rate**: Pie chart of successful vs failed redirects
3. **Response Time**: P50, P95, P99 latencies
4. **URLs Shortened**: Total URLs shortened over time
5. **Database Performance**: Database operation duration

## Docker Compose Integration

If you're running with Docker Compose, you can add Prometheus:

```yaml
version: '3.8'
services:
  url-shortener:
    # ... your existing config
    
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
```
