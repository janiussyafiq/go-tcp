package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindow struct {
	limit     int
	window    time.Duration
	counters  map[string]int
	windowEnd map[string]time.Time
	mu        sync.Mutex
}

func NewFixedWindow(limit int, window time.Duration) *FixedWindow {
	return &FixedWindow{
		limit:     limit,
		window:    window,
		counters:  make(map[string]int),
		windowEnd: make(map[string]time.Time),
	}
}

func (fw *FixedWindow) Allow(identifier string) bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()

	windowEnd, exists := fw.windowEnd[identifier]

	if !exists || now.After(windowEnd) {
		fw.counters[identifier] = 0
		fw.windowEnd[identifier] = now.Add(fw.window)
	}

	if fw.counters[identifier] >= fw.limit {
		return false
	}

	fw.counters[identifier]++
	return true
}
