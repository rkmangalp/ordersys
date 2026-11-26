# Learning Log

## 2026-11-25 — Session 1: HTTP Server Foundations

### Concepts Covered
- [x] `net/http` package basics
- [x] `http.Server` configuration (timeouts, handlers)
- [x] `http.ServeMux` routing (Go 1.22+ pattern syntax)
- [x] Middleware pattern (wrapping `http.HandlerFunc`)
- [x] ResponseWriter wrapping for status capture
- [x] Graceful shutdown with `signal.NotifyContext`
- [x] Context propagation for cancellation
- [x] JSON encoding with `encoding/json`

### Code Implemented
- `main.go` — Minimal server with `/health` and `/api/v1/products` endpoints
- Structured logging middleware
- 10s graceful shutdown timeout

### Key Insights
1. **Timeouts prevent resource leaks** — Without `ReadTimeout`/`WriteTimeout`, slow clients hold goroutines indefinitely
2. **Middleware = function composition** — `func(next http.HandlerFunc) http.HandlerFunc` is the standard pattern
3. **ResponseWriter wrapping** — Capturing status code requires implementing `http.ResponseWriter` interface
4. **Signal handling** — `signal.NotifyContext` (Go 1.16+) is cleaner than manual channel handling

### Questions for Review
1. Why does `http.ServeMux` in Go 1.22+ support method-pattern syntax (`"GET /health"`)?
2. What happens if `srv.Shutdown()` is called without a timeout context?
3. How does `loggingResponseWriter` capture the status code without breaking the response?

### Next Session Goals
- Request validation (query params, JSON body)
- Context values for request-scoped data (trace IDs, user info)
- Separate handler/service/repository packages
- Table-driven tests for handlers

---

## 2026-11-25 — Session 1 Complete: Phase 1 Done ✅

**Status:** Phase 1 complete. Server runs, handles requests, shuts down gracefully.

**Run it:**
```bash
cd C:\Users\Rk\go\src\ordersys
go run main.go
# In another terminal:
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/products
```

---

## 2026-11-25 — Session 2: Phase 2 — REST APIs & MySQL Persistence

### Concepts Covered
- [x] REST API contract design (Products CRUD)
- [x] MySQL schema design with migrations
- [x] Database connection pooling (`database/sql`)
- [x] Repository pattern with raw SQL
- [x] Service layer (business logic, validation)
- [x] Handler layer (HTTP → Service mapping)
- [x] Chi router with middleware
- [x] Error mapping (domain errors → HTTP status codes)
- [x] Pagination (offset/limit)

### Code Implemented
- `migrations/001_create_products_table.up.sql` — Products table with indexes
- `cmd/migrate/main.go` — Embedded migration runner (up/down)
- `internal/database/database.go` — Connection pool config
- `internal/product/model.go` — Domain types, requests, validation, errors
- `internal/product/repository.go` — MySQLRepository with raw SQL
- `internal/product/service.go` — Business logic, SKU uniqueness checks
- `internal/product/handler.go` — Chi handlers, JSON encoding, error mapping
- `main.go` — Wired together with chi router

### Key Insights
1. **Layered architecture** — Handler → Service → Repository keeps concerns separated
2. **Repository interface** — Enables testing with mocks, swapping implementations
3. **Domain errors** — Typed errors (ValidationError, NotFoundError, ConflictError) allow precise HTTP mapping
4. **Chi router** — Lightweight, stdlib-compatible, supports path params (`/{id}`)
5. **Migration tooling** — Embedded SQL files + version tracking = reproducible schema changes
6. **Connection pooling** — `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime` prevent exhaustion

### Questions for Review
1. Why use `sql.NullString` for nullable TEXT columns?
2. What happens if `LIMIT` is used without `ORDER BY`?
3. How does the migration tool track which migrations are applied?
4. When would you use a transaction vs. individual queries?

### Next Session Goals
- Add filtering (category, price range) and sorting to ListProducts
- Write table-driven tests for service and repository
- Add transactions for multi-table operations (later: orders + inventory)
- Introduce context values (trace ID, user ID)

---

## Template for Future Entries

### YYYY-MM-DD — Session N: Topic

#### Concepts Covered
- [ ] Concept 1
- [ ] Concept 2

#### Code Implemented
- File/function — Purpose

#### Key Insights
1. Insight 1
2. Insight 2

#### Questions for Review
1. Question 1
2. Question 2

#### Next Session Goals
- Goal 1
- Goal 2