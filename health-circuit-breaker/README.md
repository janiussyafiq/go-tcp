# Health Check + Circuit Breaker Middleware

A production-ready Go implementation of health checking and circuit breaker patterns for building resilient distributed systems.

## Overview

This project demonstrates critical reliability patterns used in production systems to prevent cascading failures and ensure high availability. It combines active health checking, passive failure detection, and circuit breakers to create a fault-tolerant reverse proxy.

## Components

### 1. Active Health Checker
Periodically pings backend servers to detect failures proactively.
- Configurable check interval and timeout
- Thread-safe status management with `sync.RWMutex`
- Graceful shutdown with context cancellation
- **Location:** `pkg/healthcheck/active.go`

### 2. Passive Health Checker
Monitors real request failures for faster detection.
- Tracks consecutive failures/successes
- Configurable thresholds for marking unhealthy/healthy
- No polling delay - reacts to actual traffic patterns
- **Location:** `pkg/healthcheck/passive.go`

### 3. Circuit Breaker
Implements the circuit breaker pattern to prevent cascading failures.
- **States:** Closed (normal) → Open (failing fast) → Half-Open (testing recovery)
- Configurable failure thresholds and recovery timeouts
- Generation counter to ignore stale results
- **Location:** `pkg/circuitbreaker/circuitbreaker.go`

### 4. Resilient Proxy
Integrates all components into a production-ready reverse proxy.
- Round-robin load balancing across healthy backends
- Automatic failover when backends fail
- Zero user impact during failures
- Auto-recovery when backends heal
- **Location:** `pkg/resilient/proxy.go`

## Architecture
```
Request → Resilient Proxy
            ↓
        Select Healthy Backend
        (Skip if health check failed OR circuit open)
            ↓
        Circuit Breaker Execute
        (Closed: allow | Open: fail-fast | Half-Open: test)
            ↓
        Reverse Proxy → Backend
            ↓
        Record Success/Failure
        (Update passive health checker)
            ↓
    Background: Active Health Checker
    (Periodic health checks every N seconds)
```

## Quick Start

### Run Tests
```bash
# Test all components
go test ./...

# Test specific component
go test ./pkg/circuitbreaker -v
go test ./pkg/healthcheck -v
```

### Run Demos

**Active Health Checker:**
```bash
go run cmd/demo/main.go
```

**Passive Health Checker:**
```bash
go run cmd/demo-passive/main.go
```

**Full Integration:**
```bash
go run cmd/demo-integrated/main.go
```

The integrated demo simulates:
1. Normal operation with 3 backends
2. Backend failures and automatic failover
3. Recovery detection and traffic resumption
4. Zero user impact during failures

## Key Patterns

### Thread Safety
- `sync.RWMutex` for read-heavy operations (health status checks)
- `sync.Mutex` for write-heavy operations (failure counters)

### Concurrency
- Goroutines for background health checking
- `context.Context` for graceful cancellation
- `sync.WaitGroup` for coordinated shutdown

### State Management
- Circuit breaker state machine (Closed/Open/Half-Open)
- Generation counter prevents stale operation results
- Atomic operations where appropriate

### Failure Handling
- **Active checks:** Detect persistent backend failures
- **Passive checks:** Detect transient request failures
- **Circuit breakers:** Prevent cascading failures with fail-fast

## Production Concepts

**Why Both Health Checks AND Circuit Breakers?**
- Health checks detect long-term failures (server down, network partition)
- Circuit breakers detect short-term bursts (overload, transient errors)
- Together they provide comprehensive protection

**Fail-Fast Pattern:**
When circuit opens, requests fail immediately without trying the backend. This prevents:
- Thread pool exhaustion
- Cascading timeouts
- Resource contention
- User experience degradation

**Gradual Recovery:**
Half-open state allows limited test requests before fully reopening circuit, preventing thundering herd problems.

## Project Structure
```
health-circuit-breaker/
├── pkg/
│   ├── healthcheck/
│   │   ├── active.go          # Active health checker
│   │   ├── passive.go         # Passive health checker
│   │   ├── passive_test.go    # Passive tests
│   │   └── types.go           # Shared types
│   ├── circuitbreaker/
│   │   ├── circuitbreaker.go  # Circuit breaker implementation
│   │   ├── circuitbreaker_test.go
│   │   └── types.go
│   └── resilient/
│       └── proxy.go           # Integrated resilient proxy
├── cmd/
│   ├── demo/                  # Active health checker demo
│   ├── demo-passive/          # Passive health checker demo
│   └── demo-integrated/       # Full integration demo
└── README.md
```

## Learning Outcomes

✅ Implemented production reliability patterns  
✅ Mastered Go concurrency primitives (goroutines, channels, mutexes)  
✅ Built thread-safe concurrent systems  
✅ Designed state machines (circuit breaker states)  
✅ Created comprehensive unit tests  
✅ Understood distributed systems failure modes

## Related Projects

This project complements:
- **Rate Limiter** - Protects backends from traffic overload
- **Reverse Proxy** - Forwards requests with load balancing
- **API Gateway Plugins** - Middleware for request processing

Together, these projects demonstrate complete backend systems knowledge from basic patterns to production-ready implementations.