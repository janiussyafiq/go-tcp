package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/janiussyafiq/health-circuit-breaker/pkg/healthcheck"
)

func main() {
	log.Println("Starting health check demo...")

	// Start test backend servers
	go startTestServer(":8081", true)  // Healthy server
	go startTestServer(":8082", true)  // Healthy server
	go startTestServer(":8083", false) // Unhealthy server

	// Give servers time to start
	time.Sleep(100 * time.Millisecond)

	// Create backends
	backends := []*healthcheck.Backend{
		healthcheck.NewBackend("http://localhost:8081/health"),
		healthcheck.NewBackend("http://localhost:8082/health"),
		healthcheck.NewBackend("http://localhost:8083/health"),
		healthcheck.NewBackend("http://localhost:9999/health"), // Non-existent server
	}

	// Create health checker
	config := healthcheck.ActiveHealthCheckConfig{
		CheckInterval: 3 * time.Second,
		Timeout:       1 * time.Second,
	}
	checker := healthcheck.NewActiveHealthChecker(backends, config)

	// Start health checking
	checker.Start()

	// Monitor backend status
	go monitorBackends(backends)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\nShutting down...")
	checker.Stop()
	log.Println("Demo finished.")
}

func startTestServer(port string, healthy bool) {
	mux := http.NewServeMux() // Create a separate handler for each server

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if healthy {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ERROR"))
		}
	})

	log.Printf("Test server started on %s (healthy=%v)", port, healthy)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Printf("Server %s error: %v", port, err)
	}
}

func monitorBackends(backends []*healthcheck.Backend) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("\n--- Backend Status ---")
		for i, backend := range backends {
			status := backend.GetStatus()
			fmt.Printf("Backend %d: %s - %s\n", i+1, backend.URL, status)
		}
		fmt.Println("----------------------")
	}
}
