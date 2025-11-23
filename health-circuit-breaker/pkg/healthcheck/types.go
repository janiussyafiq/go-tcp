package healthcheck

import (
	"sync"
	"time"
)

// HealthStatus represents the health state of a backend
type HealthStatus int

const (
	Healthy HealthStatus = iota
	Unhealthy
)

func (h HealthStatus) String() string {
	switch h {
	case Healthy:
		return "healthy"
	case Unhealthy:
		return "unhealthy"
	default:
		return "unknown"
	}
}

// Backend represents a backend server
type Backend struct {
	URL          string
	Status       HealthStatus
	LastCheck    time.Time
	FailureCount int
	mu           sync.RWMutex
}

// NewBackend creates a new backend
func NewBackend(url string) *Backend {
	return &Backend{
		URL:       url,
		Status:    Healthy,
		LastCheck: time.Now(),
	}
}

// IsHealthy returns true if backend is healthy (thread-safe)
func (b *Backend) IsHealthy() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Status == Healthy
}

// MarkUnhealthy marks backend as unhealthy (thread-safe)
func (b *Backend) MarkHealthy() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Status = Healthy
	b.FailureCount = 0
	b.LastCheck = time.Now()
}

// MarkUnhealthy marks backend as unhealthy (thread-safe)
func (b *Backend) MarkUnhealthy() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Status = Unhealthy
	b.FailureCount++
	b.LastCheck = time.Now()
}

// GetStatus returns current status (thread-safe)
func (b *Backend) GetStatus() HealthStatus {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Status
}
