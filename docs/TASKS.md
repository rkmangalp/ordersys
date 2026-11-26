# Task Tracker

## Legend
- ✅ Done (independently)
- 🟡 Done (with guidance)
- 🔄 In Progress
- ⏳ Pending
- 🚫 Blocked

---

## Phase 1 — Backend Fundamentals

| Task | Status | Notes |
|------|--------|-------|
| Initialize git repo | ✅ | |
| Create go.mod | ✅ | |
| Create main.go with net/http server | 🟡 | Guided implementation |
| Implement /health endpoint | 🟡 | |
| Implement /api/v1/products endpoint | 🟡 | |
| Add logging middleware | 🟡 | ResponseWriter wrapping |
| Add graceful shutdown | 🟡 | signal.NotifyContext + Shutdown |
| Add server timeouts | 🟡 | Read/Write/Idle |
| Create README.md | ✅ | |
| Create docs/ROADMAP.md | ✅ | |
| Create docs/LEARNING_LOG.md | ✅ | |
| Create docs/ARCHITECTURE.md | ✅ | |
| Create docs/DECISIONS.md | ✅ | |
| Create docs/TASKS.md | ✅ | This file |
| Create docs/INTERVIEW_NOTES.md | ⏳ | |
| **Test server manually with curl** | ✅ | Run after copying to go/src |
| **Add request validation (query params)** | ⏳ | Next session |
| **Add context values (trace ID)** | ⏳ | Next session |
| **Separate handler package** | ⏳ | Phase 2 |
| **Write table-driven tests** | ⏳ | Phase 1 completion |

---

## Phase 2 — REST APIs and Persistence

| Task | Status | Notes |
|------|--------|-------|
| Design REST API contracts | 🟡 | Done with guidance |
| Model database schema | 🟡 | MySQL schema in migrations |
| Write raw SQL queries | 🟡 | In repository layer |
| Implement database migrations | 🟡 | Migration tool + 001_products |
| Implement CRUD for products | 🟡 | Handler → Service → Repository |
| Connection pooling | 🟡 | Configured in database.go |
| Transactions | ⏳ | For multi-table operations |
| Pagination/filtering/sorting | 🟡 | Offset/limit implemented |
| Introduce GORM | ⏳ | After raw SQL |

---

## Recurring Tasks (Every Session)

| Task | Frequency |
|------|-----------|
| Update LEARNING_LOG.md | Each session |
| Update TASKS.md | Each session |
| Add DECISIONS.md entries | When architectural choices made |
| Add INTERVIEW_NOTES.md | Each session |
| Run `go test -race` | Before commit |
| Run `go vet ./...` | Before commit |