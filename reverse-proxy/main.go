package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

// Logging middleware that wraps our proxy
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the incoming request
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call the actual proxy handler
		next.ServeHTTP(w, r)

		// Log how long it took
		log.Printf("Request completed in %v", time.Since(start))
	})
}

// Modify headers going to the backend
func modifyRequest(r *http.Request) {
	// Add a custome header
	r.Header.Add("X-Proxy-By", "GoReverseProxy")

	// Remove a header (example: remove any existing Authorization)
	r.Header.Del("X-Remove-This")

	log.Printf("Added headers to request")
}

func main() {
	// Parse configuration
	config, err := ParseConfig()
	if err != nil {
		log.Fatal("Configuration error:", err)
	}

	lb, err := NewLoadBalancer(config.Backends)
	if err != nil {
		log.Fatal("Failed to create load balancer:", err)
	}

	// Create a handler that picks a backend for each request
	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backend := lb.NextBackend()
		log.Printf("Forwarding to backend: %s", backend.Host)

		// Create a proxy for this specific backend
		proxy := httputil.NewSingleHostReverseProxy(backend)

		// Add header modification
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			modifyRequest(req)
		}

		// Forward the request
		proxy.ServeHTTP(w, r)
	})

	// Wrap the proxy with logging middleware
	handler := loggingMiddleware(proxyHandler)

	fmt.Printf("Reverse proxy starting on port %s", config.ProxyPort)
	fmt.Printf("Load balancing across %d backends:\n", len(config.Backends))
	for i, backend := range config.Backends {
		fmt.Printf(" Backend %d: %s\n", i+1, backend)
	}
	fmt.Println()

	log.Fatal(http.ListenAndServe(":"+config.ProxyPort, handler))
}
