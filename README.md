# Concurrent Load-Balancing Reverse Proxy (Go)

## Overview

A **concurrent load-balancing reverse proxy** built in Go that distributes incoming HTTP requests across multiple backend servers while continuously monitoring their health.

This project demonstrates:
- Go networking with `net/http` and `httputil.ReverseProxy`
- Concurrency using goroutines, mutexes, and atomic operations
- Load balancing with Round-Robin strategy
- Background health monitoring
- Graceful shutdown with signal handling
- Clean and modular project structure

---

## Features

- **Reverse Proxy** using `httputil.ReverseProxy`
- **Load Balancing** with Round-Robin strategy
- **Health Checks** - Periodic backend health monitoring
- **Thread-Safe** - Server pool management with `sync.RWMutex` and `sync/atomic`
- **Admin API** - Dynamic backend management (add/remove backends at runtime)
- **Context Propagation** - Request timeout and cancellation support
- **Graceful Shutdown** - Clean server termination on SIGINT/SIGTERM

---

## Project Structure

```
ReverseProxy/
├── reverse-proxy/          # Main application
│   ├── main.go             # Entry point, server setup
│   ├── handler.go          # Proxy HTTP handler
│   ├── config.json         # Backend servers configuration
│   └── proxyConfig.json    # Proxy settings (port, strategy, health check)
├── models/                 # Data structures
│   ├── backend.go          # Backend struct with connection tracking
│   ├── serverpool.go       # ServerPool struct
│   ├── loadBalancer.go     # LoadBalancer interface & Round-Robin implementation
│   └── proxyConfig.go      # ProxyConfig struct
├── config/                 # Configuration loading
│   └── config.go           # JSON config parsers
├── admin/                  # Admin API
│   └── handlers.go         # Status, Add/Delete backend handlers
├── healthCheker/           # Health check service
│   └── checker.go          # Background health monitoring
├── backend/                # Sample backend server
├── slow_backend/           # Slow backend for testing timeouts
├── go.mod
├── README.md
└── TODO.md
```

---

## Configuration

### `config.json` - Backend Servers
```json
{
  "backends": [
    "http://localhost:8082",
    "http://localhost:8083",
    "http://localhost:8084",
    "http://localhost:8085"
  ]
}
```

### `proxyConfig.json` - Proxy Settings
```json
{
  "port": 8080,
  "strategy": "round-robin",
  "health_check_frequency": "5s"
}
```

---

## Getting Started

### Prerequisites
- Go 1.21+ installed

### Installation

```bash
git clone <repository-url>
cd ReverseProxy
go mod tidy
```

### Running the Proxy

```bash
# From the project root (ReverseProxy/)
go run reverse-proxy/*.go
```

The proxy will start on the configured port (default: `:8080`) and the admin API on `:8081`.

### Running Backend Servers (for testing)

Start multiple backend servers on different ports using Python's built-in HTTP server:

```bash
# Terminal 1
python3 -m http.server 8082

# Terminal 2
python3 -m http.server 8083

# Terminal 3
python3 -m http.server 8084

# Terminal 4
python3 -m http.server 8085
```

---

## API Reference

### Proxy Server (Port 8080)

All incoming requests are load-balanced across healthy backends.

```bash
curl http://localhost:8080/
```

### Admin API (Port 8081)

#### Get Status
```bash
GET /status
```

Response:
```json
{
  "total_backends": 4,
  "active_backends": 3,
  "backends": [...]
}
```

#### Add Backend
```bash
POST /backends
Content-Type: application/json

{
  "url": "http://localhost:8086"
}
```

#### Delete Backend
```bash
DELETE /backends
Content-Type: application/json

{
  "url": "http://localhost:8086"
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

## Key Components

### LoadBalancer Interface
```go
type LoadBalancer interface {
    GetNextValidPeer() *Backend
    AddBackend(backend *Backend)
    SetBackendStatus(uri *url.URL, alive bool)
}
```

### Backend Struct
```go
type Backend struct {
    URL                *url.URL
    Alive              bool
    CurrentConnections int64
    Mux                sync.RWMutex
}
```

### Health Checker
- Runs in a background goroutine
- Configurable check interval
- Pings `/health` endpoint on each backend
- Thread-safe status updates

---

## Graceful Shutdown

The proxy handles `SIGINT` and `SIGTERM` signals:
1. Stops accepting new connections
2. Waits for in-flight requests (30s timeout)
3. Shuts down proxy and admin servers cleanly

```bash
# Press Ctrl+C to trigger graceful shutdown
^C
Shutting down servers...
proxy server stopped gracefully
admin server stopped gracefully
All servers stopped!
```

---

## Testing

### Setup Test Environment

Open multiple terminal windows and start the services:

```bash
# Terminal 1: Start the proxy server
cd ReverseProxy
go run reverse-proxy/*.go

# Terminal 2-4: Start backend servers
python3 -m http.server 8082
python3 -m http.server 8083  
python3 -m http.server 8084

# Terminal 5: Start slow backend (for timeout testing)
go run slow_backend/slow_backend.go
```

---

### 1. Admin API Testing

#### GET /status - View all backends status
```bash
curl http://localhost:8081/status | jq
```

Expected response:
```json
{
  "total_backends": 4,
  "active_backends": 4,
  "backends": [
    {"url": "http://localhost:8082", "alive": true, "current_connections": 0},
    {"url": "http://localhost:8083", "alive": true, "current_connections": 0},
    ...
  ]
}
```

#### POST /backends - Add a new backend
```bash
# Add a new backend
curl -X POST http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8086"}'

# Expected: {"message": "Backend added successfully", "url": "http://localhost:8086"}

# Try adding duplicate backend
curl -X POST http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8086"}'

# Expected: {"message": "Backend already exists", "url": "http://localhost:8086"}
```

#### DELETE /backends - Remove a backend
```bash
# Delete existing backend
curl -X DELETE http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8086"}'

# Expected: {"message": "Backend deleted successfully", "url": "http://localhost:8086"}

# Try deleting non-existent backend
curl -X DELETE http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:9999"}'

# Expected: {"message": "Backend does not exists", "url": "http://localhost:9999"}
```

#### Error Handling - Invalid requests
```bash
# Wrong HTTP method
curl -X PUT http://localhost:8081/status
# Expected: 405 Method Not Allowed

# Invalid JSON
curl -X POST http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d 'invalid json'
# Expected: 400 Bad Request - Invalid JSON format

# Invalid URL
curl -X POST http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "not a valid url :///"}'
# Expected: 400 Bad Request - Invalid URL
```

---

### 2. Health Checker Testing

The health checker runs in a background goroutine every 5 seconds (configurable in `proxyConfig.json`).

#### Test backend going DOWN
```bash
# 1. Check initial status - all backends should be UP
curl http://localhost:8081/status | jq '.active_backends'

# 2. Stop one of the backend servers (Ctrl+C on Terminal 2)

# 3. Wait 5-10 seconds for health check to run

# 4. Check proxy logs - you should see:
#    "Backend http://localhost:8082 is DOWN"

# 5. Verify status updated
curl http://localhost:8081/status | jq '.active_backends'
# Should show one less active backend
```

#### Test backend coming back UP
```bash
# 1. Restart the stopped backend
python3 -m http.server 8082

# 2. Wait for next health check cycle (5 seconds)

# 3. Check proxy logs - you should see:
#    "Backend http://localhost:8082 is UP"

# 4. Verify status
curl http://localhost:8081/status | jq '.active_backends'
```

---

### 3. Request Handling & Load Balancing Testing

#### Test Round-Robin distribution
```bash
# Send multiple requests and observe proxy logs
for i in {1..10}; do 
  curl -s http://localhost:8080/ > /dev/null
  echo "Request $i sent"
done

# Check proxy logs - requests should rotate between backends:
# Forwarding request to http://localhost:8082
# Forwarding request to http://localhost:8083
# Forwarding request to http://localhost:8084
# Forwarding request to http://localhost:8082
# ...
```

#### Test concurrent requests
```bash
# Send 20 concurrent requests
for i in {1..20}; do 
  curl -s http://localhost:8080/ &
done
wait

# Check /status to see CurrentConnections during load
curl http://localhost:8081/status | jq '.backends[].current_connections'
```

#### Test no backends available (503 Service Unavailable)
```bash
# 1. Stop ALL backend servers

# 2. Wait for health checks to mark them as DOWN

# 3. Send a request
curl -v http://localhost:8080/

# Expected: HTTP 503 Service Unavailable
# "Service Unavailable No backends available"
```

---

### 4. Error Handling Testing

#### Test backend failure during request (502 Bad Gateway)
```bash
# The ErrorHandler marks backend as DOWN and returns 502
# This happens when a backend fails mid-request

# 1. Start slow backend
go run slow_backend/slow_backend.go

# 2. Add it to pool
curl -X POST http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8085"}'

# 3. Send request to slow backend, then kill slow_backend.go mid-request
curl http://localhost:8080/

# Expected: 502 Bad Gateway
# Proxy logs: "Backend http://localhost:8085 error: ..."
```

#### Test request timeout (30 second timeout)
```bash
# The proxy has a 30s timeout for backend requests

# 1. Make sure slow_backend is running (sleeps 60s)
go run slow_backend/slow_backend.go

# 2. Send request through proxy
curl http://localhost:8080/

# After 30 seconds, the request will timeout
# Expected: context deadline exceeded error
```

#### Test client cancellation
```bash
# Start a request and cancel it with Ctrl+C
curl http://localhost:8080/
# Press Ctrl+C before response

# The request context is cancelled and backend request is aborted
```

---

### 5. Main Goroutine & Graceful Shutdown Testing

#### Test graceful shutdown
```bash
# 1. Start proxy server
go run reverse-proxy/*.go

# 2. Send a long-running request in background
curl http://localhost:8080/ &

# 3. Press Ctrl+C on proxy server

# Expected output:
# ^C
# Shutting down servers...
# proxy server stopped gracefully
# admin server stopped gracefully
# All servers stopped!
```

#### Test concurrent servers (proxy + admin)
```bash
# Both servers run in separate goroutines
# Verify both are accessible simultaneously:

# Terminal A: Keep sending requests to proxy
while true; do curl -s http://localhost:8080/; sleep 1; done

# Terminal B: Access admin API
curl http://localhost:8081/status
```

---

### Quick Test Checklist

| Test | Command | Expected |
|------|---------|----------|
| Proxy running | `curl http://localhost:8080/` | Response from backend |
| Admin status | `curl http://localhost:8081/status` | JSON with backends |
| Add backend | `POST /backends {"url":"..."}` | Success message |
| Delete backend | `DELETE /backends {"url":"..."}` | Success message |
| Health check DOWN | Stop a backend, wait 5s | Log: "Backend ... is DOWN" |
| Health check UP | Restart backend, wait 5s | Log: "Backend ... is UP" |
| No backends | Stop all backends | 503 Service Unavailable |
| Load balancing | Send 10 requests | Round-robin in logs |
| Graceful shutdown | Ctrl+C | Clean shutdown message |

---

## Expected Outputs for Grading

This section shows the **actual outputs** captured during testing for verification.

### Proxy Server Startup

When starting the proxy with `go run reverse-proxy/*.go`:

```
Loaded 4 backends:
starting proxy server on :8080
starting admin server on :8081
2026/02/23 10:12:35 Backend http://localhost:8085 is DOWN
```

The health checker immediately runs and reports any backends that are not reachable.

---

### Admin API - GET /status

**Command:**
```bash
curl -s http://localhost:8081/status | python3 -m json.tool
```

**Actual Output:**
```json
{
    "active_backends": 3,
    "backends": [
        {
            "url": {
                "Scheme": "http",
                "Host": "localhost:8082"
            },
            "alive": true,
            "current_connections": 0
        },
        {
            "url": {
                "Scheme": "http",
                "Host": "localhost:8083"
            },
            "alive": true,
            "current_connections": 0
        },
        {
            "url": {
                "Scheme": "http",
                "Host": "localhost:8084"
            },
            "alive": true,
            "current_connections": 0
        },
        {
            "url": {
                "Scheme": "http",
                "Host": "localhost:8085"
            },
            "alive": false,
            "current_connections": 0
        }
    ],
    "total_backends": 4
}
```

---

### Admin API - POST /backends (Add Backend)

**Command:**
```bash
curl -s -X POST http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8086"}'
```

**Actual Output:**
```json
{"message":"Backend added successfully","url":"http://localhost:8086"}
```

**Adding duplicate backend:**
```bash
curl -s -X POST http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8086"}'
```

**Actual Output:**
```json
{"message":"Backend already exists","url":"http://localhost:8086"}
```

---

### Admin API - DELETE /backends (Remove Backend)

**Command:**
```bash
curl -s -X DELETE http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8086"}'
```

**Actual Output:**
```json
{"message":"Backend deleted successfully","url":"http://localhost:8086"}
```

**Deleting non-existent backend:**
```bash
curl -s -X DELETE http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:9999"}'
```

**Actual Output:**
```json
{"message":"Backend does not exists","url":"http://localhost:9999"}
```

---

### Admin API - Error Handling

**Wrong HTTP Method (PUT on GET endpoint):**
```bash
curl -s -X PUT http://localhost:8081/status -w "\nHTTP Status: %{http_code}\n"
```

**Actual Output:**
```
Method not allowed

HTTP Status: 405
```

---

### Health Checker - Backend Status Changes

**Proxy logs when backend goes DOWN:**
```
2026/02/23 10:12:35 Backend http://localhost:8085 is DOWN
```

**Proxy logs when backend comes back UP:**
```
2026/02/23 10:15:40 Backend http://localhost:8085 is UP
```

---

### Load Balancing - Round-Robin Distribution

**Sending multiple requests:**
```bash
for i in {1..4}; do curl -s http://localhost:8080/ > /dev/null; done
```

**Proxy logs showing Round-Robin:**
```
2026/02/23 10:14:55 Forwarding request to http://localhost:8082
2026/02/23 10:15:22 Forwarding request to http://localhost:8083
2026/02/23 10:15:28 Forwarding request to http://localhost:8084
2026/02/23 10:15:35 Forwarding request to http://localhost:8082
```

Notice how requests rotate through backends 8082 → 8083 → 8084 → 8082 (skipping 8085 which is DOWN).

---

### No Backends Available - 503 Error

**When all backends are DOWN:**
```bash
curl -v http://localhost:8080/
```

**Actual Output:**
```
< HTTP/1.1 503 Service Unavailable
Service Unavailable No backends available
```

**Proxy logs:**
```
2026/02/23 10:20:15 No backends available
```

---

### Error Handler - 502 Bad Gateway

**When a backend fails mid-request:**

**Proxy logs:**
```
2026/02/23 10:25:30 Backend http://localhost:8085 error: dial tcp [::1]:8085: connect: connection refused
```

**Client receives:**
```
< HTTP/1.1 502 Bad Gateway
Bad Gateway
```

---

### Graceful Shutdown

**Pressing Ctrl+C on proxy server:**

**Actual Output:**
```
^C
Shutting down servers...
proxy server stopped gracefully
admin server stopped gracefully
All servers stopped!
```

