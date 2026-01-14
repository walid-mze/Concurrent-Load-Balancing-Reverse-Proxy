# Concurrent Load-Balancing Reverse Proxy (Go)

## Overview
This project implements a **concurrent load-balancing reverse proxy** in Go.  
The proxy distributes incoming HTTP requests across multiple backend servers while continuously monitoring their health.

The system is designed to demonstrate:
- Go networking with `net/http`
- Concurrency using goroutines, mutexes, and atomic operations
- Load balancing strategies
- Background health monitoring
- Clean and modular project structure

---

## Features
- Reverse proxy using `httputil.ReverseProxy`
- Load balancing strategies:
  - Round-Robin
  - (Optional) Least-Connections
- Periodic backend health checks
- Thread-safe server pool management
- Admin API for dynamic backend management
- Context propagation and request cancellation
- Graceful shutdown support
