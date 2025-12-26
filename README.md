# RevProx - A Local Caching Reverse Proxy in Go

RevProxy is a high-performance local caching reverse proxy server written in Go. It forwards HTTP requests to an upstream origin server, caches the responses, and serves cached content for subsequent requests to improve performance and reduce latency.

## Core Functionalities

- **Transparent Proxying**: Forwards HTTP requests to a configurable upstream origin server.
- **Response Caching**:
  - Caches responses based on the request method and URL.
  - Respects configurable TTL defaults.
  - Adds `X-Cache: HIT|MISS` header to responses.
  - Obeys `Cache-Control: no-store` and `no-cache` directives from the origin.
- **CLI Interface**:
  - Launch the proxy with `caching-proxy --port <number> --origin <url>`.
  - Clear the entire cache with `caching-proxy --clear-cache`.
- **Metrics & Logging**: Provides structured logs and tracks key metrics like cache hits, misses, and upstream latency for observability.

## Technology Stack

- **Language & Framework**: Go (>=1.22) with Gin for HTTP routing.
- **CLI Toolkit**: Cobra for command parsing and Viper for configuration management.
- **Caching Layer**: Ristretto for a high-performance, in-memory cache.
- **Logging**: Zap for structured, high-performance logging.
- **Testing**: Go’s native testing package with `httptest` and Testify.

## High-Level Architecture

This project follows the principles of Clean Architecture to separate concerns and ensure testability.

![Architecture Diagram](https://raw.githubusercontent.com/mikiasgoitom/Caching-proxy/main/doc/high-arch.png)

## Local Development

### Prerequisites

- Go >=1.22

### Running the Proxy

You can run the proxy using the following command:

```bash
go run ./cmd/caching-proxy --port 3000 --origin http://dummyjson.com
```


### Clearing the Cache

To clear the cache (if a persistent cache is configured):

```bash
task clear-cache
```

### Running Tests

```bash
task test
```

