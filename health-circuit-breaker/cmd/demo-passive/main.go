package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/janiussyafiq/health-circuit-breaker/pkg/healthcheck"
)

func main() {
	log.Println("Starting passive health check demo...")

	// Create backends
	backends := []*healthcheck.Backend{
		healthcheck.NewBackend("http://backend-1:8081"),
		healthcheck.NewBackend("http://backend-2:8082"),
	}

	// Create passive health checker
	config := healthcheck.PassiveHealthCheckConfig{
		FailureThreshold: 3,
		SuccessThreshold: 2,
	}
	checker := healthcheck.NewPassiveHealthChecker(backends, config)

	log.Println("Simulating request patterns...")
	log.Println("- Backend 1: Will have increasing failures")
	log.Println("- Backend 2: Will stay mostly healthy\n")

	// Simulate requests for 20 seconds
	for i := 0; i < 40; i++ {
		time.Sleep(500 * time.Millisecond)

		// Backend 1: Gradually becomes unhealthy
		if i < 10 {
			// First 10 requests: all succeed
			simulateRequest(checker, backends[0], true)
		} else if i < 15 {
			// Next 5 requests: all fail (should mark unhealthy)
			simulateRequest(checker, backends[0], false)
		} else if i < 18 {
			// Next 3 requests: success (should mark healthy again)
			simulateRequest(checker, backends[0], true)
		} else {
			// Rest: random (to show continuos monitoring)
			simulateRequest(checker, backends[0], rand.Float32() > 0.5)
		}

		// Backend 2: Mostly healthy with occasional failures
		simulateRequest(checker, backends[1], rand.Float32() > 0.2) // 80% success

		// Print status every 5 iterations
		if (i+5)%5 == 0 {
			printStatus(backends, checker)
		}
	}

	log.Println("\nDemo completed!")
}

func simulateRequest(checker *healthcheck.PassiveHealthChecker, backend *healthcheck.Backend, success bool) {
	if success {
		checker.RecordSuccess(backend.URL)
		fmt.Printf("✓ %s - Success\n", backend.URL)
	} else {
		checker.RecordFailure(backend.URL)
		fmt.Printf("✗ %s - Failure\n", backend.URL)
	}
}

func printStatus(backends []*healthcheck.Backend, checker *healthcheck.PassiveHealthChecker) {
	fmt.Println("\n--- Backend Status ---")
	for _, backend := range backends {
		failures := checker.GetConsecutiveFailures(backend.URL)
		successes := checker.GetConsecutiveSuccesses(backend.URL)
		status := "healthy"
		if !backend.IsHealthy() {
			status = "UNHEALTHY"
		}
		fmt.Printf("%s: %s (failures: %d, successes: %d)\n",
			backend.URL, status, failures, successes)
	}
	fmt.Println("----------------------\n")
}
