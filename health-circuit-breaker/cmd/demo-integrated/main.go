package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/janiussyafiq/health-circuit-breaker/pkg/resilient"
)

func main() {
	log.Println("=== Resilient Proxy Demo ===")
	log.Println("Starting test backedn servers...")

	// Start test backends
	backend1 := startBackend(":8081", "Backend-1")
	backend2 := startBackend(":8082", "Backend-2")
	backend3 := startBackend(":8083", "Backend-3")

	// Give servers time to start
	time.Sleep(100 * time.Millisecond)

	// Create resilient proxy
	proxy, err := resilient.NewResilientProxy(resilient.Config{
		Backends: []string{
			"http://localhost:8081",
			"http://localhost:8082",
			"http://localhost:8083",
		},
		ActiveHealthCheckInterval: 3 * time.Second,
		ActiveHealthCheckTimeout:  1 * time.Second,
		PassiveFailureThreshold:   3,
		PassiveSuccessThreshold:   2,
		CircuitBreakerTimeout:     10 * time.Second,
		CircuitBreakerMaxRequests: 1,
	})
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
	}

	proxy.Start()
	defer proxy.Stop()

	// Start proxy server
	proxyServer := &http.Server{
		Addr:    ":9000",
		Handler: proxy,
	}

	go func() {
		log.Println("Proxy listening on :9000")
		if err := proxyServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Proxy server error: %v", err)
		}
	}()

	// Simulate different failure scenarios
	go simulateFailures(backend1, backend2, backend3)

	// Send test traffic
	go sendTraffic()

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\nShutting down...")
	proxyServer.Close()
}

type testBackend struct {
	name    string
	healthy atomic.Bool
	port    string
}

func startBackend(port, name string) *testBackend {
	backend := &testBackend{
		name: name,
		port: port,
	}
	backend.healthy.Store(true)

	mux := http.NewServeMux()

	// Healthy check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if backend.healthy.Load() {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s is healthy", backend.name)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s is unhealthy", backend.name)
		}
	})

	// API endpoint
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		if !backend.healthy.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s unavailable", backend.name)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Response from %s", backend.name)
	})

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	go func() {
		log.Printf("%s started on %s", backend.name, port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("%s error: %v", backend.name, err)
		}
	}()

	return backend
}

func simulateFailures(b1, b2, b3 *testBackend) {
	time.Sleep(5 * time.Second)
	log.Println("\n>>> Scenario 1: Backend-2 starts failing")
	b2.healthy.Store(false)

	time.Sleep(15 * time.Second)
	log.Println("\n>>> Scenario 2: Backend-2 recovers")
	b2.healthy.Store(true)

	time.Sleep(10 * time.Second)
	log.Println("\n>>> Scenario 3: Backend-1 and Backend-3 fail")
	b1.healthy.Store(false)
	b3.healthy.Store(false)

	time.Sleep(15 * time.Second)
	log.Println("\n>>> Scenario 3.5: All backends fail")
	b1.healthy.Store(false)
	b2.healthy.Store(false)
	b3.healthy.Store(false)

	time.Sleep(15 * time.Second)
	log.Println("\n>>> Scenario 4: All backends recover")
	b1.healthy.Store(true)
	b2.healthy.Store(true)
	b3.healthy.Store(true)
}

func sendTraffic() {
	time.Sleep(2 * time.Second)
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	requestNum := 1
	for range ticker.C {
		resp, err := client.Get("http://localhost:9000/api")
		if err != nil {
			log.Printf("Request #%d: ERROR - %v", requestNum, err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				log.Printf("Request #%d: ✓ SUCCESS (%d)", requestNum, resp.StatusCode)
			} else {
				log.Printf("Request #%d: ✗ FAILED (%d)", requestNum, resp.StatusCode)
			}
		}
		requestNum++
	}
}
