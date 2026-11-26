# Architectural Decisions

## ADR Template
```
## [ADR-NNN] Title
**Date:** YYYY-MM-DD
**Status:** Proposed | Accepted | Superseded
**Context:** What problem are we solving?
**Decision:** What did we choose?
**Consequences:** Trade-offs, follow-up work
```

---

## ADR-001: Use Go Stdlib for HTTP Server (Phase 1)
**Date:** 2026-11-25
**Status:** Accepted
**Context:** Need a minimal HTTP server for learning fundamentals.
**Decision:** Use `net/http` only — no Gin, Chi, or other frameworks.
**Consequences:**
- ✅ Learn how HTTP actually works (routing, middleware, timeouts)
- ✅ Zero dependencies, fast builds
- ❌ More boilerplate for path params, validation, binding
- 🔄 Re-evaluate in Phase 2 when adding REST API complexity

---

## ADR-002: Go 1.22+ ServeMux Pattern Syntax
**Date:** 2026-11-25
**Status:** Accepted
**Context:** Go 1.22 introduced method + pattern syntax (`"GET /health"`).
**Decision:** Use new syntax instead of manual method checks in handlers.
**Consequences:**
- ✅ Cleaner, declarative routing
- ✅ Built-in method validation (405 Method Not Allowed)
- ❌ Requires Go 1.22+ (acceptable — current is 1.23)

---

## ADR-003: Middleware as Function Composition
**Date:** 2026-11-25
**Status:** Accepted
**Context:** Need cross-cutting concerns (logging, auth, tracing).
**Decision:** Use `func(http.HandlerFunc) http.HandlerFunc` pattern.
**Consequences:**
- ✅ Simple, composable, stdlib-compatible
- ✅ Easy to test (wrap a test handler)
- ❌ No built-in middleware chain (unlike Gin/Chi)
- 🔄 May introduce a lightweight chain helper later

---

## ADR-004: Graceful Shutdown with signal.NotifyContext
**Date:** 2026-11-25
**Status:** Accepted
**Context:** Server must handle SIGINT/SIGTERM cleanly.
**Decision:** Use `signal.NotifyContext` (Go 1.16+) + `srv.Shutdown(ctx)`.
**Consequences:**
- ✅ Clean, no manual signal channels
- ✅ Integrates with context cancellation
- ✅ Configurable drain timeout
- ❌ None significant

---

## ADR-005: Server Timeouts (Read/Write/Idle)
**Date:** 2026-11-25
**Status:** Accepted
**Context:** Prevent resource exhaustion from slow/malicious clients.
**Decision:** Set conservative defaults: Read=10s, Write=15s, Idle=60s.
**Consequences:**
- ✅ Bounds goroutine lifetime
- ✅ Protects against Slowloris-style attacks
- ❌ May need tuning for large payloads/streaming
- 🔄 Make configurable via env in Phase 5

---

## ADR-006: Project Structure — Flat First, Layered Later
**Date:** 2026-11-25
**Status:** Accepted
**Context:** Where to put code as it grows.
**Decision:** Single `main.go` for Phase 1; introduce `internal/` layers in Phase 2.
**Consequences:**
- ✅ No premature abstraction
- ✅ Structure reflects actual complexity
- 🔄 Refactor when handlers > 3 or logic mixes concerns