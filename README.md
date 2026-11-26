# Order & Inventory Management Platform

A production-style backend system built from first principles — learning in public.

## Project Overview

This is a **modular monolith** implementing an Order and Inventory Management Platform, built progressively following the [Backend from First Principles](https://youtube.com/playlist?list=PLui3EUkuMTPgZcV0QhQrOcwMPcBCcd_Q1) course.

### Current Phase: Phase 1 — Backend Fundamentals

- [x] Basic HTTP server with `net/http`
- [x] Request routing with `http.ServeMux`
- [x] JSON serialization
- [x] Request logging middleware
- [x] Graceful shutdown with signal handling
- [x] Server timeouts (read, write, idle)

### Planned Phases

| Phase | Focus |
|-------|-------|
| 1 | HTTP server, routing, middleware, JSON |
| 2 | REST APIs, SQL, database, migrations, GORM |
| 3 | Auth (sessions, JWT, RBAC), security |
| 4 | Caching (Redis), async (Kafka, workers), Temporal |
| 5 | Observability, resilience, config |
| 6 | Testing, profiling, load testing |
| 7 | Distributed systems, Docker, Kubernetes |

## Quick Start

```bash
# Clone
git clone https://github.com/rkmangalp/ordersys.git
cd ordersys

# Initialize module (first time only)
go mod tidy

# Run server
go run main.go

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/products
```

## Project Structure (Evolving)

```
ordersys/
├── main.go           # Entry point, server setup
├── go.mod            # Module definition
├── .gitignore
├── README.md
└── docs/             # Documentation (grows with project)
    ├── ROADMAP.md
    ├── LEARNING_LOG.md
    ├── ARCHITECTURE.md
    ├── DECISIONS.md
    ├── INTERVIEW_NOTES.md
    └── TASKS.md
```

## Learning Log

| Date | Concept | Exercise |
|------|---------|----------|
| 2026-11-25 | HTTP server, net/http, graceful shutdown | Minimal server with health + products endpoints |

## Engineering Principles

1. **Learn by building** — No copy-paste without understanding
2. **Start simple** — Stdlib first, frameworks later
3. **Document decisions** — Every architectural choice in `docs/DECISIONS.md`
4. **Test as you go** — Table-driven tests, integration tests
5. **Observe everything** — Structured logging, metrics, tracing

## License

MIT — Learning project, use freely.