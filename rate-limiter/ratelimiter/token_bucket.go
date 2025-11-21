package ratelimiter

import (
	"sync"
	"time"
)

type bucket struct {
	tokens     float64
	lastRefill time.Time
}

type TokenBucket struct {
	capacity   float64
	refillRate float64
	buckets    map[string]*bucket
	mu         sync.Mutex
}

func NewTokenBucket(capacity int, tokensPerWindow int, refillInterval time.Duration) *TokenBucket {
	refillRate := float64(tokensPerWindow) / refillInterval.Seconds()

	return &TokenBucket{
		capacity:   float64(capacity),
		refillRate: refillRate,
		buckets:    make(map[string]*bucket),
	}
}

func (tb *TokenBucket) Allow(identifier string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()

	b, exists := tb.buckets[identifier]
	if !exists {
		b = &bucket{
			tokens:     tb.capacity,
			lastRefill: now,
		}
		tb.buckets[identifier] = b
	}

	tb.refill(b, now)

	if b.tokens < 1 {
		return false
	}

	b.tokens--
	return true
}

func (tb *TokenBucket) refill(b *bucket, now time.Time) {
	elapsed := now.Sub(b.lastRefill).Seconds()

	tokensToAdd := elapsed * tb.refillRate

	b.tokens = min(b.tokens+tokensToAdd, tb.capacity)

	b.lastRefill = now
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
