# Project Plan

## Month 1 — Foundations
**Theme: Get uncomfortable in Go. Build habits.**

### Go Project (Phase 1)
- Define the `Task` struct. Think hard about it before writing code.
- Implement the in-memory queue with a mutex. Read about `sync.Mutex` properly.
- Write your first goroutines and channels — producers and consumers.
- Run `go test -race` and fix every race condition you find.
- Add structured logging with `slog` (stdlib, Go 1.21+).

---

## Month 2 — Persistence & Reliability
**Theme: Learn how real systems handle failure.**

### Go Project (Phase 2)
- Integrate **BoltDB or Badger** for persistence. Read their docs, understand why an embedded DB makes sense here.
- Implement crash recovery — kill your process mid-run and verify tasks aren't lost.
- Build exponential backoff retry logic from scratch (don't import a library).
- Implement Dead Letter Queue.
- Add idempotency keys — understand *why* they exist, not just how.

---

## Month 3 — Networking & Distribution
**Theme: Your system talks to the outside world.**

### Go Project (Phase 3)
- Expose an **HTTP API** first (simpler), then consider gRPC.
- Implement worker registration with heartbeat — think about what happens when a worker dies silently.
- Read the **Raft paper** before implementing leader election. Understand it conceptually first.
- Add a `/metrics` endpoint. Manually, not with a library yet — so you understand what you're exporting.
- Write integration tests, not just unit tests.

---

## Month 4 — AI Layer & Job Preparation Begins
**Theme: Bridge into AI infrastructure.**

### Go Project (Phase 4)
- Define AI task types. Build the provider abstraction layer properly — this is an interface design exercise.
- Integrate one real LLM API (Anthropic or OpenAI). Keep it simple.
- Implement rate limiting — build a token bucket from scratch.
- Add cost tracking per task.

---

## Month 5 — Observability & Polish
**Theme: Make the project production-worthy.**

### Go Project (Phase 6 — skip Phase 5 for now)
- Add **Prometheus metrics export**.
- Implement task tracing and execution history.
- Write a proper **architecture document** — explain every design decision as if presenting to a senior engineer.
- Record a short **walkthrough video** of the system. This becomes part of your portfolio.

---

## Month 6 — Finish the Project

### Go Project (Phase 7 — your DevOps background)
- Containerize with Docker (multi-stage builds).
- Deploy to AKS — you already know this, it costs you almost no time.
- Set up GitHub Actions CI/CD.
- Add a basic Grafana dashboard.
- **Project is now complete and deployed.** Link in your CV.
