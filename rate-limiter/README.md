# Rate Limiter in Go

Production-ready rate limiting implementing three algorithms: Fixed Window, Sliding Window, and Token Bucket.

## Overview

Concurrent, goroutine-safe rate limiting for APIs. Implements algorithms used in AWS API Gateway, APISIX, and Kong.

**Features:**
- 3 rate limiting algorithms (Fixed Window, Sliding Window, Token Bucket)
- HTTP middleware for easy integration
- Thread-safe with mutex locks
- Memory efficient O(1) space for Token Bucket
- IP detection with proxy support

## Quick Start

### Installation
```bash
git clone https://github.com/yourusername/go-tcp
cd rate-limiter
go mod tidy
```

### Basic Usage
```go
limiter := ratelimiter.NewTokenBucket(5, 5, time.Minute)

if limiter.Allow("192.168.1.1") {
    // Process request
} else {
    // Return 429 Too Many Requests
}
```

### HTTP Middleware
```go
limiter := ratelimiter.NewTokenBucket(100, 100, time.Hour)
mux.Handle("/api", middleware.RateLimitMiddleware(limiter)(handler))
http.ListenAndServe(":8080", mux)
```

## Algorithm Comparison

| Algorithm      | Time | Space/User | Burst | Edge Cases |
|----------------|------|------------|-------|------------|
| Fixed Window   | O(1) | O(1)       | ❌    | Window boundary exploit |
| Sliding Window | O(n) | O(n)       | ❌    | None |
| Token Bucket   | O(1) | O(1)       | ✅    | None |

**Recommendation:** Token Bucket (industry standard, supports burst traffic)

## Key Concepts

### Fixed Window Edge Case
```
Window 1 [0-60s]:  5 requests at t=59s ✓
Window 2 [60-120s]: 5 requests at t=61s ✓
Result: 10 requests in 2 seconds (violates 5/min limit!)
```

### Token Bucket (Production Standard)
- Capacity: Maximum burst size
- Refill rate: Sustained rate limit
- Allows bursts while enforcing average rate

## Technical Details

- **Concurrency:** `sync.Mutex` for goroutine safety
- **IP Detection:** X-Forwarded-For → X-Real-IP → RemoteAddr
- **Memory:** Token Bucket uses 16 bytes per user

## Learning Outcomes

✅ Production rate limiting algorithms  
✅ HTTP middleware patterns  
✅ Concurrent programming in Go  
✅ Interface-based design  
✅ Time-based algorithms  

## Related Projects

- **Reverse Proxy** - Load balancing with health checks
- **Circuit Breaker** - Failure handling and resilience
- **APISIX Plugins** - API gateway middleware

---

Built to understand production-grade distributed systems patterns.