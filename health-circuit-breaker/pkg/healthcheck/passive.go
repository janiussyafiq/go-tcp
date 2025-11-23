package healthcheck

import (
	"log"
	"sync"
)

// PassiveHealthChecker monitors request failures
type PassiveHealthChecker struct {
	backends             []*Backend
	failureThreshold     int
	successThreshold     int
	mu                   sync.RWMutex
	consecutiveFailures  map[string]int
	consecutiveSuccesses map[string]int
}

// PassiveHealthCheckConfig holds configuration
type PassiveHealthCheckConfig struct {
	FailureThreshold int // Mark unhealthy after N consecutive failures
	SuccessThreshold int // Mark healthy after N consecutive successes
}

// NewPassiveHealthChecker creates a new passive health checker
func NewPassiveHealthChecker(backends []*Backend, config PassiveHealthCheckConfig) *PassiveHealthChecker {
	return &PassiveHealthChecker{
		backends:             backends,
		failureThreshold:     config.FailureThreshold,
		successThreshold:     config.SuccessThreshold,
		consecutiveFailures:  make(map[string]int),
		consecutiveSuccesses: make(map[string]int),
	}
}

// RecordSuccess records a successful request to a backend
func (phc *PassiveHealthChecker) RecordSuccess(backendURL string) {
	phc.mu.Lock()
	defer phc.mu.Unlock()

	// Reset failure count
	phc.consecutiveFailures[backendURL] = 0

	// Increment success count
	phc.consecutiveSuccesses[backendURL]++

	// Check if we should mark backend as healthy
	if phc.consecutiveSuccesses[backendURL] >= phc.successThreshold {
		backend := phc.findBackend(backendURL)
		if backend != nil && backend.GetStatus() == Unhealthy {
			backend.MarkHealthy()
			log.Printf("Passive: Backend %s marked healthy after %d successes",
				backendURL, phc.successThreshold)
		}
		// Reset success count after marking healthy
		phc.consecutiveSuccesses[backendURL] = 0
	}
}

// RecordFailure records a failed request to a backend
func (phc *PassiveHealthChecker) RecordFailure(backendURL string) {
	phc.mu.Lock()
	defer phc.mu.Unlock()

	// Reset success count
	phc.consecutiveSuccesses[backendURL] = 0

	// Increment failure count
	phc.consecutiveFailures[backendURL]++

	// Check if we should mark unhealthy
	if phc.consecutiveFailures[backendURL] >= phc.failureThreshold {
		backend := phc.findBackend(backendURL)
		if backend != nil && backend.GetStatus() == Healthy {
			backend.MarkUnhealthy()
			log.Printf("Passive: Backend %s marked unhealthy after %d failures",
				backendURL, phc.failureThreshold)
		}
		// Reset failure count after marking unhealthy
		phc.consecutiveFailures[backendURL] = 0
	}
}

// findBackend finds a backend by URL (caller must hold lock)
func (phc *PassiveHealthChecker) findBackend(url string) *Backend {
	for _, backend := range phc.backends {
		if backend.URL == url {
			return backend
		}
	}
	return nil
}

// GetConsecutiveFailures returns failure count for a backend (for testing)
func (phc *PassiveHealthChecker) GetConsecutiveFailures(backendURL string) int {
	phc.mu.RLock()
	defer phc.mu.RUnlock()
	return phc.consecutiveFailures[backendURL]
}

// GetConsecutiveSuccesses success failure count for a backend (for testing)
func (phc *PassiveHealthChecker) GetConsecutiveSuccesses(backendURL string) int {
	phc.mu.RLock()
	defer phc.mu.RUnlock()
	return phc.consecutiveSuccesses[backendURL]
}
