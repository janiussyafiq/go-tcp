package healthcheck

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

// ActiveHealthChecker periodically pings backends
type ActiveHealthChecker struct {
	backends      []*Backend
	checkInterval time.Duration
	timeout       time.Duration
	httpClient    *http.Client
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// ActivateHealthCheckConfig holds configuration for active health checks
type ActiveHealthCheckConfig struct {
	CheckInterval time.Duration // How often to check backends
	Timeout       time.Duration // HTTP request timeout
}

// NewActiveHealthChecker creates a new active health checker
func NewActiveHealthChecker(backends []*Backend, config ActiveHealthCheckConfig) *ActiveHealthChecker {
	ctx, cancel := context.WithCancel(context.Background())

	return &ActiveHealthChecker{
		backends:      backends,
		checkInterval: config.CheckInterval,
		timeout:       config.Timeout,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins health checking in the background
func (hc *ActiveHealthChecker) Start() {
	log.Println("Active health checker started")

	hc.wg.Add(1)
	go hc.healthCheckLoop()
}

// Stop gracefully stops the health checker
func (hc *ActiveHealthChecker) Stop() {
	log.Println("Stopping active health checker...")
	hc.cancel()
	hc.wg.Wait()
	log.Println("Active health checker stopped")
}

// healtCheckLoop runs the periodic health check
func (hc *ActiveHealthChecker) healthCheckLoop() {
	defer hc.wg.Done()

	ticker := time.NewTicker(hc.checkInterval)
	defer ticker.Stop()

	// Do initial check immediately
	hc.checkAllBackends()

	for {
		select {
		case <-ticker.C:
			hc.checkAllBackends()
		case <-hc.ctx.Done():
			return
		}
	}
}

// checkAllBackends checks all backends concurrently
func (hc *ActiveHealthChecker) checkAllBackends() {
	var wg sync.WaitGroup

	for _, backend := range hc.backends {
		wg.Add(1)
		go func(b *Backend) {
			defer wg.Done()
			hc.checkBackend(b)
		}(backend)
	}

	wg.Wait()
}

// checkBackend performs a health check on a single backend
func (hc *ActiveHealthChecker) checkBackend(backend *Backend) {
	req, err := http.NewRequestWithContext(hc.ctx, http.MethodGet, backend.URL, nil)
	if err != nil {
		log.Printf("Failed to create request for %s: %v", backend.URL, err)
		backend.MarkUnhealthy()
		return
	}

	resp, err := hc.httpClient.Do(req)
	if err != nil {
		log.Printf("Health check failed for %s: %v", backend.URL, err)
		backend.MarkUnhealthy()
		return
	}
	defer resp.Body.Close()

	// Consider 2xx and 3xx as healthy
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		wasUnhealthy := backend.GetStatus() == Unhealthy
		backend.MarkHealthy()
		if wasUnhealthy {
			log.Printf("Backend %s recovered", backend.URL)
		}
	} else {
		log.Printf("Backend %s returned status %d", backend.URL, resp.StatusCode)
		backend.MarkUnhealthy()
	}
}
