# Concurrent Load-Balancing Reverse Proxy (Go)

A **concurrent load-balancing reverse proxy** built in Go that distributes incoming HTTP requests across multiple backend servers while continuously monitoring their health.

## Features

- **Reverse Proxy** - `httputil.ReverseProxy` for request forwarding
- **Round-Robin Load Balancing** - Atomic counter for thread-safe distribution
- **Health Checks** - Background goroutine with configurable interval
- **Thread-Safe** - `sync.RWMutex` and `sync/atomic` for concurrent access
- **Admin API** - Dynamic backend management at runtime
- **Context Propagation** - Request timeout and client cancellation support
- **Graceful Shutdown** - Clean termination on SIGINT/SIGTERM

---

## Project Structure

```
ReverseProxy/
├── reverse-proxy/          # Main application
│   ├── main.go             # Entry point, server setup
│   ├── handler.go          # Proxy HTTP handler
│   ├── config.json         # Backend servers list
│   └── proxyConfig.json    # Proxy settings
├── models/                 # Data structures
│   ├── backend.go          # Backend struct
│   ├── serverpool.go       # ServerPool struct
│   ├── loadBalancer.go     # LoadBalancer interface
│   └── proxyConfig.go      # ProxyConfig struct
├── config/                 # Configuration loading
├── admin/                  # Admin API handlers
├── healthCheker/           # Health check service
└── slow_backend/           # Test server for timeouts
```

---

## Quick Start

### 1. Start Backend Servers (for testing)

Start multiple backend servers on different ports using Python's built-in HTTP server:

```bash
# Terminal 1
python3 -m http.server 8082

# Terminal 2
python3 -m http.server 8083

# Terminal 3
python3 -m http.server 8084
```

### 2. Start Proxy

From the project root directory:

```bash
cd ReverseProxy
go run reverse-proxy/*.go
```

**Expected Output:**
```
Loaded 4 backends:
starting proxy server on :8080
starting admin server on :8081
```

The proxy will start on port `:8080` and the admin API on `:8081`.

---

## Configuration

**config.json** - Backend servers:
```json
{
  "backends": ["http://localhost:8082", "http://localhost:8083", "http://localhost:8084"]
}
```

**proxyConfig.json** - Proxy settings:
```json
{
  "port": 8080,
  "strategy": "round-robin",
  "health_check_frequency": "5s"
}
```

---

## Architecture

```
                    ┌─────────────────┐
                    │     Client      │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │  Reverse Proxy  │  :8080
                    │  (Load Balancer)│
                    └────────┬────────┘
                             │
           ┌─────────────────┼─────────────────┐
           │                 │                 │
           ▼                 ▼                 ▼
    ┌──────────┐      ┌──────────┐      ┌──────────┐
    │ Backend1 │      │ Backend2 │      │ Backend3 │
    │  :8082   │      │  :8083   │      │  :8084   │
    └──────────┘      └──────────┘      └──────────┘

                    ┌─────────────────┐
                    │   Admin API     │  :8081
                    │ (Status/Manage) │
                    └─────────────────┘
```

---

## API Reference

### Proxy (Port 8080)
```bash
curl http://localhost:8080/
```

### Admin API (Port 8081)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/status` | GET | View all backends status |
| `/backends` | POST | Add new backend |
| `/backends` | DELETE | Remove backend |

---

## Testing & Expected Outputs

### GET /status

Command: `curl -s http://localhost:8081/status`

**Response:**
```json
{
    "total_backends": 4,
    "active_backends": 3,
    "backends": [
        {"url": {"Host": "localhost:8082"}, "alive": true, "current_connections": 0},
        {"url": {"Host": "localhost:8083"}, "alive": true, "current_connections": 0},
        {"url": {"Host": "localhost:8084"}, "alive": true, "current_connections": 0},
        {"url": {"Host": "localhost:8085"}, "alive": false, "current_connections": 0}
    ]
}
```

### POST /backends - Add Backend

Command: `curl -X POST http://localhost:8081/backends -H "Content-Type: application/json" -d '{"url": "http://localhost:8086"}'`

| Scenario | Response |
|----------|----------|
| Success | `{"message":"Backend added successfully","url":"http://localhost:8086"}` |
| Duplicate | `{"message":"Backend already exists","url":"http://localhost:8086"}` |

### DELETE /backends - Remove Backend

Command: `curl -X DELETE http://localhost:8081/backends -H "Content-Type: application/json" -d '{"url": "http://localhost:8086"}'`

| Scenario | Response |
|----------|----------|
| Success | `{"message":"Backend deleted successfully","url":"http://localhost:8086"}` |
| Not found | `{"message":"Backend does not exists","url":"http://localhost:9999"}` |

### Error Handling

| Test | Response |
|------|----------|
| Wrong HTTP method (`PUT /status`) | 405 Method Not Allowed |

---

## Health Checker

Runs every 5 seconds (configurable in `proxyConfig.json`). Logs status changes:

```
Backend http://localhost:8085 is DOWN
Backend http://localhost:8085 is UP
```

**Test:** Stop a backend server → wait 5s → check proxy logs for status change.

---

## Load Balancing (Round-Robin)

Send multiple requests and observe proxy logs rotating between backends:

```
Forwarding request to http://localhost:8082
Forwarding request to http://localhost:8083
Forwarding request to http://localhost:8084
Forwarding request to http://localhost:8082
```

### Connection Tracking

Each backend tracks active connections using `sync/atomic`:

```bash
# Send concurrent requests
for i in {1..10}; do curl -s http://localhost:8080/ & done

# Check connections during load
curl -s http://localhost:8081/status | grep current_connections
```

The `CurrentConnections` counter increments before forwarding and decrements after completion.

---

## Error Responses

| Scenario | HTTP Code | Message |
|----------|-----------|---------|
| No backends available | 503 | Service Unavailable No backends available |
| Backend failure | 502 | Bad Gateway |
| Wrong HTTP method | 405 | Method not allowed |

### Test 503 - No Backends Available

1. Stop ALL backend servers (Ctrl+C on each terminal)
2. Wait for health check to mark them as DOWN (5 seconds)
3. Send a request: `curl -v http://localhost:8080/`

**Expected:** HTTP 503 Service Unavailable

### Test 502 - Backend Failure

When a backend fails during a request, the ErrorHandler returns 502 and marks the backend as DOWN.

**Proxy logs:** `Backend http://localhost:8085 error: dial tcp [::1]:8085: connect: connection refused`

---

## Graceful Shutdown

Press `Ctrl+C` to trigger graceful shutdown:

```
Shutting down servers...
proxy server stopped gracefully
admin server stopped gracefully
All servers stopped!
```

---

## Context Timeout & Client Cancellation

The proxy uses `context.WithTimeout` to set a **30-second timeout** for backend requests. This ensures slow backends don't block indefinitely.

### Test Request Timeout

1. Start the slow backend: `go run slow_backend/slow_backend.go`
2. Send request through proxy: `curl http://localhost:8080/`
3. After 30 seconds, the request will timeout with: `context deadline exceeded`

### Test Client Cancellation

When a client cancels their request (Ctrl+C), the context is cancelled and the backend request is aborted.


---

## Conclusion

We have learned a lot throughout this project and I really enjoyed working on it, especially in Go.

### Key Concepts Learned

- Concurrency in Go (goroutines, channels)
- Synchronization (`sync.RWMutex`, `sync/atomic`)
- HTTP Networking (`net/http`, `httputil.ReverseProxy`)
- Load Balancing (Round-Robin)
- Health Monitoring (`time.Ticker`)
- Context Propagation (`context.Context`)
- Graceful Shutdown (SIGINT/SIGTERM)
- RESTful API Design

---

### Acknowledgments

I want to thank **Mr. Abdelghafour Mourchid** for his guidance and support throughout all the bootcamp. His teaching made learning Go an enjoyable and rewarding experience.

---

Walid MAMZE - 2026 :) 
