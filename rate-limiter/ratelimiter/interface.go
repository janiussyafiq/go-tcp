package ratelimiter

type RateLimiter interface {
	Allow(identifier string) bool
}
