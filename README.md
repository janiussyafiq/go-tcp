# Rate Limiter in Go

A production-ready rate limiting system implementing three industry-standard algorithms with HTTP middleware support.

## Overview

This project implements concurrent, goroutine-safe rate limiting algorithms commonly used in API gateways like AWS API Gateway, Google Cloud, APISIX, and Kong. Built with Go, it demonstrates core distributed systems concepts including token bucket algorithms, sliding window techniques, and HTTP middleware patterns.

## Features

- **3 Rate Limiting Algorithms**
  - Fixed Window - Simple counter-based limiting
  - Sliding Window - Accurate timestamp-based limiting
  - Token Bucket - Production-standard with burst support
  
- **HTTP Middleware** - Drop-in rate limiting for any HTTP endpoint
- **Goroutine-Safe** - Concurrent request handling with mutex locks
- **Memory Efficient** - O(1) space complexity for Token Bucket
- **IP Detection** - Supports proxied requests (X-Forwarded-For, X-Real-IP)

## 🏗️ Architecture

```
HTTP Request → Middleware (Extract IP) → Rate Limiter → API Handler
                                              ↓
                                         Allow/Block
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher

### Installation

```bash
git clone https://github.com/yourusername/go-tcp
cd rate-limiter
go mod tidy
```

### Run Tests

```bash
# Test individual algorithms
# (uncomment any test you want to try)
go run *.go

# Start HTTP server
# (uncomment testHTTPServer() in main.go)
go run *.go
```

### Basic Usage

```go
package main

import (
    "time"
    "rate-limiter/ratelimiter"
)

func main() {
    // Create a token bucket: 5 requests per minute
    limiter := ratelimiter.NewTokenBucket(5, 5, time.Minute)
    
    // Check if request is allowed
    if limiter.Allow("192.168.1.1") {
        // Process request
    } else {
        // Return 429 Too Many Requests
    }
}
```

### HTTP Server Example

```go
package main

import (
    "net/http"
    "time"
    "rate-limiter/middleware"
    "rate-limiter/ratelimiter"
)

func main() {
    limiter := ratelimiter.NewTokenBucket(100, 100, time.Hour)
    mux := http.NewServeMux()
    
    // Add rate limiting middleware
    mux.Handle("/api/data", 
        middleware.RateLimitMiddleware(limiter)(
            http.HandlerFunc(dataHandler),
        ),
    )
    
    http.ListenAndServe(":8080", mux)
}
```

## 📊 Algorithm Comparison

| Algorithm      | Time Complexity | Space per User | Burst Support | Edge Cases |
|----------------|----------------|----------------|---------------|------------|
| Fixed Window   | O(1)           | O(1)           | ❌            | Window boundary exploit |
| Sliding Window | O(n)           | O(n)           | ❌            | None |
| Token Bucket   | O(1)           | O(1)           | ✅            | None |

**Recommendation:** Use Token Bucket for production (industry standard).

## 🧪 Testing

The project includes comprehensive tests demonstrating:
- Normal rate limiting behavior
- Window boundary edge cases
- Burst traffic handling
- Token refill mechanisms
- Concurrent request safety

## 🔑 Key Concepts

### Token Bucket Algorithm
```
Bucket Capacity: 5 tokens
Refill Rate: 1 token/second

Start:        ○○○○○ (5 tokens)
5 requests:   _____ (0 tokens) → All allowed (burst!)
Wait 1s:      ○     (1 token refilled)
1 request:    _____ (0 tokens) → Allowed
```

### Fixed Window Edge Case
```
Window 1 [0-60s]:  5 requests at 59s ✓
Window 2 [60-120s]: 5 requests at 61s ✓
Result: 10 requests in 2 seconds! (violates 5/min limit)
```

### Sliding Window Fix
```
At 61s, look back 60 seconds → window [1s-61s]
Sees all 10 requests → Blocks correctly ✓
```

## 🛠️ Technical Details

### Concurrency Safety
- Uses `sync.Mutex` for goroutine-safe operations
- Lock/defer unlock pattern prevents race conditions
- Supports concurrent HTTP requests

### IP Detection
- Checks `X-Forwarded-For` header (proxy/load balancer)
- Falls back to `X-Real-IP` header
- Extracts IP from `RemoteAddr` as last resort

### Memory Efficiency
- **Token Bucket**: 16 bytes per user (float64 + time.Time)
- **Fixed Window**: 24 bytes per user (int + time.Time)
- **Sliding Window**: 24 * n bytes per user (slice of timestamps)

## 🎯 Use Cases

- **API Rate Limiting** - Protect backend services from abuse
- **DDoS Prevention** - Block excessive requests from single IPs
- **Cost Control** - Limit expensive operations (e.g., AI API calls)
- **Fair Usage** - Ensure equal access for all users

## 📚 Learning Outcomes

This project demonstrates:
- Rate limiting algorithms used in production systems
- HTTP middleware patterns
- Concurrent programming in Go
- Interface-based design for swappable implementations
- Time-based algorithm implementations

## 🔗 References

- [AWS API Gateway Rate Limiting](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-request-throttling.html)
- [Token Bucket Algorithm](https://en.wikipedia.org/wiki/Token_bucket)
- [APISIX Rate Limiting Plugin](https://apisix.apache.org/docs/apisix/plugins/limit-count/)

## 📝 License

MIT License - Feel free to use this for learning and production!

## 👤 Author

janiussyafiq. Built as a learning project to understand production-grade rate limiting systems.

---

⭐ If you found this helpful, please star the repository!