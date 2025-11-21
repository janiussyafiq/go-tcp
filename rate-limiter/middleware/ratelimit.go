package middleware

import (
	"fmt"
	"net"
	"net/http"

	"rate-limiter/ratelimiter"
)

func RateLimitMiddleware(limiter ratelimiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)

			if !limiter.Allow(clientIP) {
				w.Header().Set("X-RateLimit-Exceeded", "true")
				w.WriteHeader(http.StatusTooManyRequests) // 429
				fmt.Fprintf(w, "Rate limit exceeded for IP: %s\n", clientIP)
				fmt.Fprintf(w, "Please try again later.\n")
				return
			}

			w.Header().Set("X-RateLimit-Allowed", "true")
			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		for i, c := range forwarded {
			if c == ',' {
				forwarded = forwarded[:i]
				break
			}
		}

		forwarded = trimSpace(forwarded)

		if ip := net.ParseIP(forwarded); ip != nil {
			return forwarded
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && s[start] == ' ' {
		start++
	}

	end := len(s)
	for end > start && s[end-1] == ' ' {
		end--
	}

	return s[start:end]
}
