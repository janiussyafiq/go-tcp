package resilient

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/janiussyafiq/health-circuit-breaker/pkg/circuitbreaker"
	"github.com/janiussyafiq/health-circuit-breaker/pkg/healthcheck"
)

// BackendServer represents a backend with health checking and circuit breaker
type BackendServer struct {
	Backend        *healthcheck.Backend
	CircuitBreaker *circuitbreaker.CircuitBreaker
	ReverseProxy   *httputil.ReverseProxy
}

// ResilientProxy is a reverse proxy with health checking and circuit breaking
type ResilientProxy struct {
	backends           []*BackendServer
	activeHealthCheck  *healthcheck.ActiveHealthChecker
	passiveHealthCheck *healthcheck.PassiveHealthChecker
	currentIndex       int
	mu                 sync.Mutex
}

// Config holds configuration for the resilient proxy
type Config struct {
	Backends                  []string
	ActiveHealthCheckInterval time.Duration
	ActiveHealthCheckTimeout  time.Duration
	PassiveFailureThreshold   int
	PassiveSuccessThreshold   int
	CircuitBreakerTimeout     time.Duration
	CircuitBreakerMaxRequests uint32
}

// New ResilientProxy creates a new resilient reverse proxy
func NewResilientProxy(config Config) (*ResilientProxy, error) {
	if len(config.Backends) == 0 {
		return nil, fmt.Errorf("at least one backend required")
	}

	// Create backend servers
	backends := make([]*BackendServer, 0, len(config.Backends))
	healthCheckBackends := make([]*healthcheck.Backend, 0, len(config.Backends))

	for _, backendURL := range config.Backends {
		// Parse URL
		target, err := url.Parse(backendURL)
		if err != nil {
			return nil, fmt.Errorf("invalid backend URL %s: %w", backendURL, err)
		}

		// Create health check backend
		backend := healthcheck.NewBackend(backendURL)
		healthCheckBackends = append(healthCheckBackends, backend)

		// Create circuit breaker
		cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.Config{
			Name:        backendURL,
			MaxRequests: config.CircuitBreakerMaxRequests,
			Timeout:     config.CircuitBreakerTimeout,
			ReadyToTrip: func(counts circuitbreaker.Counts) bool {
				// Open circuit after 5 consecutive failures
				return counts.ConsecutiveFailures >= 5
			},
			OnStateChange: func(name string, from circuitbreaker.State, to circuitbreaker.State) {
				log.Printf("Circuit breaker [%s]: %s -> %s", name, from, to)
			},
		})

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(target)

		backends = append(backends, &BackendServer{
			Backend:        backend,
			CircuitBreaker: cb,
			ReverseProxy:   proxy,
		})
	}

	// Create active health checker
	activeChecker := healthcheck.NewActiveHealthChecker(healthCheckBackends, healthcheck.ActiveHealthCheckConfig{
		CheckInterval: config.ActiveHealthCheckInterval,
		Timeout:       config.ActiveHealthCheckTimeout,
	})

	// Create passive health checker
	passiveChecker := healthcheck.NewPassiveHealthChecker(healthCheckBackends, healthcheck.PassiveHealthCheckConfig{
		FailureThreshold: config.PassiveFailureThreshold,
		SuccessThreshold: config.PassiveSuccessThreshold,
	})

	return &ResilientProxy{
		backends:           backends,
		activeHealthCheck:  activeChecker,
		passiveHealthCheck: passiveChecker,
	}, nil
}

// Start begins health checking
func (rp *ResilientProxy) Start() {
	rp.activeHealthCheck.Start()
	log.Println("Resilient proxy started with health checking")
}

// Stop gracefully stops the proxy
func (rp *ResilientProxy) Stop() {
	rp.activeHealthCheck.Stop()
	log.Println("Resilient proxy stopped")
}

// ServeHTTP handles incoming requests with load balancing, health checking, and circuit breaking
func (rp *ResilientProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := rp.selectHealthyBackend()
	if backend == nil {
		log.Println("No healthy backends available")
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	// Try to execute request through circuit breaker
	_, err := backend.CircuitBreaker.Execute(func() (interface{}, error) {
		// Create a custom response writer to capture status
		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		backend.ReverseProxy.ServeHTTP(recorder, r)

		// Check if request was successful
		if recorder.statusCode >= 500 {
			return nil, fmt.Errorf("backend returned %d", recorder.statusCode)
		}
		return nil, nil
	})

	// Record result in passive health checker
	if err != nil {
		if err == circuitbreaker.ErrCircuitOpen {
			log.Printf("Circuit breaker open for %s", backend.Backend.URL)
			http.Error(w, "Service Unavailable - Circuit Open", http.StatusServiceUnavailable)
		} else if err == circuitbreaker.ErrTooManyRequests {
			log.Printf("Too many requests for %s in half-open state", backend.Backend.URL)
			http.Error(w, "Service Unavailable - Too Many Requests", http.StatusServiceUnavailable)
		}
		rp.passiveHealthCheck.RecordFailure(backend.Backend.URL)
	} else {
		rp.passiveHealthCheck.RecordSuccess(backend.Backend.URL)
	}
}

// selectHealthyBackend uses round-robin to select a healthy backend
func (rp *ResilientProxy) selectHealthyBackend() *BackendServer {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	// Try all backends starting from current index
	for i := 0; i < len(rp.backends); i++ {
		idx := (rp.currentIndex + i) % len(rp.backends)
		backend := rp.backends[idx]

		// Check if backend is healthy and circuit is not open
		if backend.Backend.IsHealthy() && backend.CircuitBreaker.State() != circuitbreaker.StateOpen {
			rp.currentIndex = (idx + 1) % len(rp.backends)
			return backend
		}
	}

	return nil
}

// statusRecorder captures HTTP response status
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
