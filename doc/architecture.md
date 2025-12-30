# Caching Proxy System Design

## 1. Project Intent

- Build a production-like yet fully local caching proxy to deepen understanding of reverse proxies, HTTP caching semantics, and CLI tooling in Go.
- Focus on transparent request forwarding, response caching, cache observability, and developer ergonomics.

## 2. Core Functionalities

- **CLI entrypoint**: `caching-proxy --port <number> --origin <url>` to launch the proxy on the desired port and bind it to a single upstream origin.
- **Transparent proxying**: Accept any HTTP method, forward to `<origin>/<path>`, stream headers/body back to the client.
- **Response caching**:
  - Cache key = method + full request URL (path + sorted query string).
  - Respect configurable TTL defaults (e.g. 5 minutes) with optional overrides via CLI flags or config file.
  - Add `X-Cache: HIT|MISS` header on responses.
  - Support cache busting via `caching-proxy --clear-cache`.
- **Cache controls**: Honor `Cache-Control: no-store` or `no-cache` headers from origin to optionally bypass caching.
- **Metrics & logging**: Track hits, misses, upstream latency, and log structured events to help debugging.

## 3. Non-Functional Goals

- **Local-first**: Require no external managed services; rely on embedded or file-based storage options.
- **Performance**: Handle concurrent requests efficiently with minimal contention in the cache layer.
- **Resilience**: Fail fast on cache/storage errors while allowing proxy to continue serving live upstream responses.
- **Observability**: Provide verbose mode for tracing cache behavior during development.

## 4. Recommended Technology Stack

- **Language & Framework**: Go (>=1.22) with Gin for HTTP routing/middleware ergonomics.
- **CLI Toolkit**: Cobra (command parsing) + Viper (config/flags/env binding) to deliver a polished CLI experience.
- **Caching Layer**:
  - Primary: Ristretto (fast in-memory cache with TTL, metrics) or BigCache for pure in-memory needs.
  - Optional persistence: BadgerDB or BoltDB if experimentation with disk-backed caching is desired.
- **HTTP Client**: Native `http.Client` with custom transport for timeout/retry tuning.
- **Logging**: Zerolog or Zap for structured logs; fallback to Go log for simplicity.
- **Testing**: Go’s testing package + httptest servers for reproducible integration tests.

## 5. Technology Comparison

| Capability             | Primary Choice                         | Strengths                                                                                                        | Notable Alternatives                     | Trade-offs                                                                                                                                                                   |
| ---------------------- | -------------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ---------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| HTTP server            | Gin                                    | Mature ecosystem, middleware support, good ergonomics while staying performant                                   | net/http stdlib, Fiber, Echo             | Gin adds a minimal abstraction layer; stdlib offers lowest overhead but requires more boilerplate, Fiber/Echo faster but diverge from Go idioms                              |
| CLI tooling            | Cobra + Viper                          | Widely adopted, clean command/subcommand model, seamless flag/env/config merging                                 | urfave/cli, Kong                         | Cobra has steeper learning curve; urfave/cli simpler but weaker subcommand ergonomics, Kong offers declarative structs but smaller community                                 |
| Caching                | Ristretto                              | High hit ratio, TTL support, metrics hooks, concurrent friendly                                                  | BigCache, go-cache, Redis (local daemon) | Ristretto higher memory overhead than BigCache; go-cache simple but less performant; Redis adds network hop and requires separate process but offers persistence and tooling |
| Persistence (optional) | BadgerDB                               | Embedded, high write throughput, LSM-based suited for large values                                               | BoltDB, SQLite                           | Badger heavier operationally; Bolt simpler but slower for large datasets; SQLite requires schema management                                                                  |
| Logging                | Zerolog or Zap                         | Structured JSON, zero allocation (Zerolog) or fast sugared logging (Zap), rich ecosystem                         | Logrus, slog (stdlib exp)                | Zerolog more opinionated API; Logrus slower; slog evolving but still young                                                                                                   |
| Metrics                | Prometheus client                      | De facto standard, easy to expose `/metrics` endpoint                                                            | Expvar, StatsD                           | Prometheus needs pull-based setup; expvar limited metrics types                                                                                                              |
| Testing                | `testing` + httptest + Testify/Mockery | Standard toolkit with richer assertions (Testify) and interface mocks (Mockery) fitting clean architecture ports | Ginkgo/Gomega, GoMock                    | Ginkgo expressive BDD but heavier; GoMock integrates with go generate yet generates fragile expectations                                                                     |

> **Testing note**: Layered architecture works well with Testify for expressive assertions and suites. Mockery can autogenerate doubles for cache/origin ports, speeding up unit tests. Retain plain `testing` for table-driven tests and integration coverage.

## 6. High-Level Architecture

### Clean Architecture View

- **Domain (Entities + Value Objects)**: Request/response models, cache metadata, cache policy representations kept free of framework dependencies.
- **Use Cases (Application Services)**: Orchestrate proxy flow—cache lookup, origin fetch, cache write, cache clearing—exposed through interfaces.
- **Interface Adapters**: Gin handlers, CLI commands, DTO mappers translating external input to domain models and vice versa.
- **Infrastructure**: Concrete implementations (Gin server, Ristretto cache adapter, HTTP client, persistence, logging, metrics) wired via dependency injection in the composition root.

### Runtime Interaction

```mermaid
graph TD
  A[CLI Command] -->|flags/config| B[Config Builder]
  B --> C[Composition Root]
  C --> D[Proxy Server (Gin Adapter)]
  D --> E[Proxy Use Case]
  E -->|lookup| F[Cache Port]
  F -->|hit| G[HTTP Response DTO]
  F -->|miss| H[Origin Port]
  H -->|HTTP request| I[Origin Server]
  I --> H
  H -->|persist| F
  E --> J[Metrics & Logging Port]
  E --> G
```

### Component Responsibilities

- **CLI Command**: Parse args (`--port`, `--origin`, TTL overrides, cache size, verbose mode) and execute subcommands (`serve`, `clear-cache`).
- **Config Builder**: Merge defaults, config files (optional `config.yaml`), env vars, and CLI flags.
- **Composition Root**: Assemble dependencies, bind interfaces to concrete adapters per clean architecture boundaries.
- **Proxy Server**: Gin app exposing `/` wildcard route, layered with middlewares for logging, cache headers, error handling.
- **Proxy Use Case**: Pure application service coordinating cache lookup, origin fetch, policy checks, and response construction.
- **Cache Port + Adapter**: Interface defining cache operations with adapter implementations (Ristretto, BigCache, Badger).
- **Origin Port + Adapter**: Interface for outbound HTTP requests, allowing test doubles and alternate transports.
- **Metrics & Logging Port + Adapter**: Abstract telemetry contracts with pluggable structured logging/metrics backends.

## 7. Request Lifecycle

1. CLI launches server; configuration is validated (port availability, origin URL).
2. Incoming client request hits Gin router; middleware composes canonical cache key.
3. Cache Manager checks store:
   - On hit: return cached payload + headers, append `X-Cache: HIT`.
   - On miss: forward request to origin via Origin Client.
4. Origin response streamed back; Cache Manager persists body + headers subject to TTL and cacheability policy.
5. Response returned to client with `X-Cache: MISS` on first pass.
6. Metrics updated (hit/miss counters, upstream latency), logs emitted.

## 8. Cache Management Strategy

- **TTL policy**: Default TTL with optional header-based overrides (e.g. respect `Cache-Control: max-age`).
- **Memory bounds**: Configure cache size via CLI (e.g. `--cache-size=128MB`).
- **Eviction**: Ristretto/BigCache handle LRU-like eviction; expose stats via `caching-proxy stats` command.
- **Clearing cache**: `caching-proxy --clear-cache` locates cache directory (if disk-backed) or triggers in-memory purge through cache API.

## 9. Local Development Approach

- Provide `Makefile` or Taskfile with targets:
  - `task run -- --port 3000 --origin http://dummyjson.com`
  - `task clear-cache`
  - `task test`
- Include example `.env.example` covering optional config values.
- Ship Postman collection or `httpie` examples inside `docs/` for manual testing.

## 10. Extension Ideas

- Add offline mode serving only cached responses when origin is unavailable.
- Implement cache warming via preload list.
- Expose simple web UI at `/admin` for inspecting cache entries.
- Support request/response compression negotiation.
- Integrate automated benchmarks to compare cache on/off scenarios.
