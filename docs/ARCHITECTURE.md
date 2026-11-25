# Architecture

Phase 1: single `main.go` process — mux, handlers, logging middleware, graceful shutdown.
No frameworks. Allowed: `net/http`, `encoding/json`, `log`, `time`, `context`, `os`, `os/signal`, `syscall`.
