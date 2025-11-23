package main

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// )

// // Test that LoadBalancer distributes requests evenly
// func TestRoundRobin(t *testing.T) {
// 	backendURLs := []string{
// 		"http://backend1.com",
// 		"http://backend2.com",
// 		"http://backend3.com",
// 	}

// 	lb, err := NewLoadBalancer(backendURLs)
// 	if err != nil {
// 		t.Fatalf("Failed to create load balancer: %v", err)
// 	}

// 	// Make 9 requests and track which backend is selected
// 	results := make([]string, 9)
// 	for i := 0; i < 9; i++ {
// 		backend := lb.NextBackend()
// 		results[i] = backend.Host
// 	}

// 	// Verify round-robin pattern: 1, 2, 3, 1, 2, 3, 1, 2, 3
// 	expected := []string{
// 		"backend1.com", "backend2.com", "backend3.com",
// 		"backend1.com", "backend2.com", "backend3.com",
// 		"backend1.com", "backend2.com", "backend3.com",
// 	}

// 	for i, host := range results {
// 		if host != expected[i] {
// 			t.Errorf("Request %d: expected %s, got %s", i, expected[i], host)
// 		}
// 	}
// }

// // Test that headers are added correctly
// func TestHeaderModification(t *testing.T) {
// 	// Create a mock backend server
// 	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Check if our custom headear was added
// 		proxyHeader := r.Header.Get("X-Proxy-By")
// 		if proxyHeader != "GoReverseProxy" {
// 			t.Errorf("Expected X-Proxy-By header, got: %s", proxyHeader)
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("OK"))
// 	}))
// 	defer backend.Close()

// 	// Create a test request
// 	req := httptest.NewRequest("GET", "http://example.com/test", nil)
// 	modifyRequest(req)

// 	// Check if the header was added
// 	proxyHeader := req.Header.Get("X-Proxy-By")
// 	if proxyHeader != "GoReverseProxy" {
// 		t.Errorf("Expected X-Proxy-By header 'GoReverseProxy', got: '%s'", proxyHeader)
// 	}
// }

// // Test load balancer with single backend
// func TestSingleBackend(t *testing.T) {
// 	lb, err := NewLoadBalancer([]string{"http://localhost:8081"})
// 	if err != nil {
// 		t.Fatalf("Failed to create load balancer: %v", err)
// 	}

// 	// All requests should go to the same backend
// 	for i := 0; i < 5; i++ {
// 		backend := lb.NextBackend()
// 		if backend.Host != "localhost:8081" {
// 			t.Errorf("Expected localhost:8081, got %s", backend.Host)
// 		}
// 	}
// }
