package main

import (
	"fmt"
	"net/url"
	"sync"
)

// LoadBalancer handles distributing requests across backends
type LoadBalancer struct {
	backends []*url.URL
	current  int
	mu       sync.Mutex
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(backendURLs []string) (*LoadBalancer, error) {
	lb := &LoadBalancer{
		backends: make([]*url.URL, 0, len(backendURLs)),
		current:  0,
	}

	// Parse all backend URLs
	for _, backendURL := range backendURLs {
		parsedURL, err := url.Parse(backendURL)
		if err != nil {
			return nil, fmt.Errorf("invalid backend URL %s: %v", backendURL, err)
		}
		lb.backends = append(lb.backends, parsedURL)
	}

	return lb, nil
}

// NextBackend returns the next backend in round-robin fashion
func (lb *LoadBalancer) NextBackend() *url.URL {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	backend := lb.backends[lb.current]
	lb.current = (lb.current + 1) % len(lb.backends)

	return backend
}
