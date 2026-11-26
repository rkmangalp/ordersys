# Interview Notes

## Session 1: HTTP Server Fundamentals

### Questions to Be Able to Answer

#### net/http Internals
1. **How does `http.Server` handle incoming connections?**
   - Listens on TCP socket, accepts connections, spawns goroutine per connection, parses HTTP, calls Handler

2. **What's the difference between `ReadTimeout`, `WriteTimeout`, and `IdleTimeout`?**
   - Read: max time to read request (headers + body)
   - Write: max time to write response
   - Idle: max time for keep-alive connection to wait for next request

3. **How does `http.ServeMux` match routes?**
   - Longest pattern match wins; Go 1.22+ supports method + path patterns

4. **Why wrap `http.ResponseWriter`?**
   - To capture status code, response size, or modify headers after handler runs

#### Concurrency & Context
5. **What happens if a handler panics?**
   - Server recovers, logs, returns 500 — unless custom `Handler` with `recover()`

6. **How does graceful shutdown work?**
   - Stop accepting new connections → wait for in-flight (with timeout) → force close

7. **Why use `signal.NotifyContext` instead of a channel?**
   - Integrates with context cancellation, composable, cleaner API

#### JSON & Encoding
8. **How does `json.NewEncoder(w).Encode(v)` differ from `json.Marshal(v)`?**
   - Encoder streams directly to Writer (less allocation), Marshal returns []byte

9. **What happens if struct fields are unexported?**
   - Ignored by encoder (unless custom MarshalJSON)

#### Testing & Debugging
10. **How would you test the logging middleware?**
    - Create `httptest.ResponseRecorder`, call middleware with test handler, assert log output

---

## Template for Future Sessions

### Session N: Topic

#### Questions to Be Able to Answer

**Category 1: Concept**
1. Question
2. Question

**Category 2: Implementation**
1. Question
2. Question

**Category 3: Trade-offs**
1. Question
2. Question

#### "Explain Like I'm 5" Versions
- Complex concept → simple analogy

#### Code Walkthrough Prep
- Be able to trace a request through your code line by line