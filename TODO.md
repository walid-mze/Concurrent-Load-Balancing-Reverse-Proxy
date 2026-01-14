## Phase 1: Project Setup & Structure

 Initialize Go module 

 Create project structure

 Load configuration from config.json

 Parse backend URLs from config

## Phase 2: Data Models

 Implement Backend struct

 Implement ServerPool struct

 Implement ProxyConfig struct

## Phase 3: Load Balancer Logic

 Define LoadBalancer interface

 Implement Round-Robin strategy

 Handle case: no alive backends : return 503

## Phase 4: Concurrency & Thread Safety

 Protect server pool with sync.RWMutex

 Use sync/atomic


 Ensure no race conditions (run go test -race)

## Phase 5: Reverse Proxy Handler

 Create HTTP handler using net/http

 Select backend using LoadBalancer

 Use httputil.ReverseProxy

 Forward request with original context.Context

 Increment CurrentConns before forwarding

 Decrement CurrentConns after request finishes

 Implement custom ErrorHandler

## Phase 6: Health Checker (Background Goroutine)

 Create health checker service

 Use time.Ticker with configurable interval

 Periodically ping each backend

 Update Alive status

 Log status changes (UP / DOWN)

 Ensure thread-safe updates

## Phase 7: Admin API

(run on separate port, e.g. :8081)

 GET /status

 POST /backends

 DELETE /backends

 Validate input & handle errors properly

## Phase 8: Context, Timeouts & Shutdown

 Add timeouts for backend requests

 Cancel backend request if client disconnects

 Implement graceful shutdown


## Phase 9: Testing & Validation

 Test with multiple dummy backend servers

 Simulate backend failure

 Verify load balancing behavior

 Test concurrent requests

 Test Admin API endpoints

## Phase 10: Documentation & Cleanup

 Write README.md

 Ensure clean package separation

## Optional Enhancements