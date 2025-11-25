# Go TCP & Distributed Systems

A collection of Go implementations demonstrating production-grade distributed systems patterns and API gateway middleware.

## Projects

### `tcp-echo/`
Basic TCP echo server for learning TCP networking fundamentals.

### `reverse-proxy/`
Lightweight reverse proxy with round-robin load balancing, request forwarding, and header manipulation.

### `rate-limiter/`
Production-ready rate limiting implementing Fixed Window, Sliding Window, and Token Bucket algorithms. Includes HTTP middleware and supports burst traffic handling.

### `health-circuit-breaker/`
Resilient reverse proxy with active/passive health checking and circuit breaker pattern for preventing cascading failures in distributed systems.

### `apisix-plugins/`
Custom Lua plugins for Apache APISIX demonstrating API gateway middleware patterns (request logging, IP blocking, header injection).

## Learning Focus

These projects cover essential backend engineering concepts:
- Concurrent programming with goroutines and mutexes
- HTTP middleware patterns
- Load balancing and failover strategies
- Rate limiting algorithms
- Circuit breaker and health check patterns
- API gateway plugin development

Each folder contains its own README with detailed documentation, quick start guides, and implementation details.
