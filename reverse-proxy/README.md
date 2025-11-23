# Go Reverse Proxy

A lightweight reverse proxy with round-robin load balancing built in Go.

## Features

- Request forwarding to backend services
- Header manipulation (add/remove)
- Request logging with timing metrics
- Round-robin load balancing across multiple backends
- Thread-safe concurrent request handling

## Quick Start

**Start test backends:**
```bash
cd backend
go run backend.go 8081
go run backend.go 8082
go run backend.go 8083
```

**Start proxy with configuration:**
```bash
# Single backend
go run . -port=8080 -backends="http://localhost:8081"

# Multiple backends (load balancing)
go run . -port=8080 -backends="http://localhost:8081,http://localhost:8082,http://localhost:8083"
```

**Test it:**
```bash
curl http://localhost:8080/test
curl http://localhost:8080/test
curl http://localhost:8080/test
curl http://localhost:8080/test
curl http://localhost:8080/test
```

Watch the logs to see requests distributed across backends!

## Run Tests
Uncomment the content of the file first, and comment any main block that is declared
```bash
go test -v
```

## How It Works

The proxy uses round-robin load balancing to distribute requests evenly across backends. A mutex ensures thread-safety when multiple requests arrive concurrently.

Each request gets logged with method, path, client IP, and completion time. Custom headers (`X-Proxy-By`) are automatically added.

## What I Learned

- Go's `httputil.ReverseProxy` and HTTP handling
- Thread-safe concurrent programming with `sync.Mutex`
- Middleware patterns for logging
- Round-robin algorithm implementation
- Writing unit tests in Go