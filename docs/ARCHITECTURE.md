# Architecture Document

## Current Architecture (Phase 1)

```
┌─────────────────────────────────────┐
│           main.go                   │
│  ┌───────────────────────────────┐  │
│  │      http.Server              │  │
│  │  ┌─────────────────────────┐  │  │
│  │  │    http.ServeMux        │  │  │
│  │  │  ┌───────────────────┐  │  │  │
│  │  │  │ loggingMiddleware │  │  │  │
│  │  │  └─────────┬─────────┘  │  │  │
│  │  │            ▼            │  │  │
│  │  │  ┌───────────────────┐  │  │  │
│  │  │  │ healthHandler     │  │  │  │
│  │  │  │ productsHandler   │  │  │  │
│  │  │  └───────────────────┘  │  │  │
│  │  └─────────────────────────┘  │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility |
|-----------|----------------|
| `main()` | Server lifecycle, signal handling, graceful shutdown |
| `http.Server` | TCP listener, connection management, timeouts |
| `http.ServeMux` | Request routing (method + path → handler) |
| `loggingMiddleware` | Cross-cutting concern: request/response logging |
| `healthHandler` | Liveness/readiness endpoint |
| `productsHandler` | Placeholder for product listing |

### Data Flow

```
HTTP Request
    │
    ▼
http.Server (read timeout)
    │
    ▼
http.ServeMux (route match)
    │
    ▼
loggingMiddleware (wrap ResponseWriter)
    │
    ▼
Handler (business logic)
    │
    ▼
JSON Response
    │
    ▼
loggingMiddleware (log: method, path, status, duration)
    │
    ▼
HTTP Response
```

## Planned Evolution

### Phase 2: Layered Architecture
```
cmd/server/main.go
internal/
├── handler/       # HTTP layer (request/response)
├── service/       # Business logic
├── repository/    # Data access
└── model/         # Domain entities
```

### Phase 3+: Modular Monolith
```
internal/
├── user/
├── product/
├── warehouse/
├── inventory/
├── order/
├── shipment/
└── shared/        # Common: config, logging, middleware
```

## Design Principles

1. **Stdlib first** — Understand before abstracting
2. **Explicit over implicit** — No magic, clear dependencies
3. **Single responsibility** — Each package has one reason to change
4. **Dependency direction** — Handler → Service → Repository (never reverse)
5. **Interfaces at boundaries** — Define contracts where packages meet