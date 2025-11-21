package ratelimiter

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	limit    int
	window   time.Duration
	requests map[string][]time.Time
	mu       sync.Mutex
}

func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		limit:    limit,
		window:   window,
		requests: make(map[string][]time.Time),
	}
}

func (sw *SlidingWindow) Allow(identifier string) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()

	windowStart := now.Add(-sw.window)

	timestamps := sw.requests[identifier]

	validTimestamps := make([]time.Time, 0, len(timestamps))

	for _, ts := range timestamps {
		if !ts.Before(windowStart) {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	if len(validTimestamps) >= sw.limit {
		sw.requests[identifier] = validTimestamps
		return false
	}

	validTimestamps = append(validTimestamps, now)
	sw.requests[identifier] = validTimestamps
	return true
}
