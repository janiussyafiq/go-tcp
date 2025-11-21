package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"rate-limiter/middleware"
	"rate-limiter/ratelimiter"
)

func testHTTPServer() {
	fmt.Println("🌐 HTTP Server with Rate Limiting")
	fmt.Println("==================================\n")

	fixedWindow := ratelimiter.NewFixedWindow(5, time.Minute)
	slidingWindow := ratelimiter.NewSlidingWindow(5, time.Minute)
	tokenBucket := ratelimiter.NewTokenBucket(5, 5, time.Minute)

	// Create HTTP multiplexer (router)
	mux := http.NewServeMux()

	// Health check endpoint (no rate limiting)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server is healthy!\n")
	})

	// API endpoints with DIFFERENT rate limiting strategies

	// Fixed Window endpoint
	mux.Handle("/api/fixed",
		middleware.RateLimitMiddleware(fixedWindow)(
			http.HandlerFunc(apiHandler("Fixed Window")),
		),
	)

	// Sliding Window endpoint
	mux.Handle("/api/sliding",
		middleware.RateLimitMiddleware(slidingWindow)(
			http.HandlerFunc(apiHandler("Sliding Window")),
		),
	)

	// Token Bucket endpoint
	mux.Handle("/api/token",
		middleware.RateLimitMiddleware(tokenBucket)(
			http.HandlerFunc(apiHandler("Token Bucket")),
		),
	)

	// Print server info
	port := ":8080"
	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("\nEndpoints:")
	fmt.Println("  • GET /health           - Health check (no rate limit)")
	fmt.Println("  • GET /api/fixed        - Fixed Window (5 req/min)")
	fmt.Println("  • GET /api/sliding      - Sliding Window (5 req/min)")
	fmt.Println("  • GET /api/token        - Token Bucket (5 req/min)")
	fmt.Println("\nTest with:")
	fmt.Println("  curl http://localhost:8080/api/fixed")
	fmt.Println("\nOr run multiple requests:")
	fmt.Println("  for i in {1..7}; do curl http://localhost:8080/api/fixed; echo \"\"; done")
	fmt.Println("\nPress Ctrl+C to stop the server")
	fmt.Println(repeatString("=", 60))

	// Start the server
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

// apiHandler creates a simple API handler that returns success
func apiHandler(strategy string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "✅ Success!\n")
		fmt.Fprintf(w, "Strategy: %s\n", strategy)
		fmt.Fprintf(w, "Time: %s\n", time.Now().Format(time.RFC3339))
		fmt.Fprintf(w, "Your IP: %s\n", r.RemoteAddr)
	}
}
