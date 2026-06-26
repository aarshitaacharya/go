## Week 1 — TCP Foundations

- Set up a TCP server in Go listening on port 6379
- Accepted and read raw bytes from clients using `net.Listen` and `conn.Read`
- Handled client disconnects gracefully
- Handled multiple concurrent clients using goroutines (`go handleConnection(...)`)

**What `go` does here:** The `go` keyword spawns a new goroutine — a lightweight thread managed by the Go runtime. Each client connection runs in its own goroutine, so one slow or blocked client doesn't stall others. The Go scheduler multiplexes these goroutines across OS threads automatically.

---

## Week 2 — A Real Protocol

- Implemented a text-based command protocol over TCP
- Parsed raw input into command + arguments using `strings.Fields`
- Built a command dispatcher (`switch` on command name)
- Supported commands: `SET`, `GET`, `DEL`, `EXISTS`
- Handled malformed input — wrong argument counts return descriptive errors
- Storage backed by a simple `map[string]string`

---

## The Race Condition (intentional)

Two clients `SET` the same key simultaneously. Go's `map` is not safe for concurrent writes — the race detector (`go run -race`) flags this. Not fixed yet — that's Week 3's problem.

---

## What I Learned About TCP Servers

TCP gives you a stream of bytes, not messages. There's no built-in concept of where one command ends and the next begins — that's what a protocol is for. Building even a trivial text protocol forces you to think about framing, error cases, and what happens when a client sends garbage. Concurrency makes all of this harder: a map that works fine with one client silently corrupts data with five. The race detector is the tool that makes the invisible visible.