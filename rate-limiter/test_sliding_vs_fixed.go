package main

import (
	"fmt"
	"strings"
	"time"

	"rate-limiter/ratelimiter"
)

func testSlidingVsFixed() {
	fmt.Println("Rate Limiter Comparison Test")
	fmt.Println("----------------------------")

	userIP := "192.168.1.1"

	// Test 1: Normal rate limiting
	fmt.Println("TEST 1: Normal Rate Limiting")
	fmt.Println("Configuration: 5 requests per 10 seconds")

	testNormalRateLimiting(userIP)

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Test 2: The edge case
	fmt.Println("TEST 2: Edge Case at Window Boudary")
	fmt.Println("(This is where Fixed Window has a problem)")

	testEdgeCase(userIP)

	fmt.Println("\nAll tests completed!")
}

// Test normal rate limiting behavior
func testNormalRateLimiting(userIP string) {
	fixedWindow := ratelimiter.NewFixedWindow(5, 10*time.Second)
	slidingWindow := ratelimiter.NewSlidingWindow(5, 10*time.Second)

	fmt.Println("Making 8 rapid requests:")
	fmt.Println("------------------------")

	for i := 1; i <= 8; i++ {
		fixedAllowed := fixedWindow.Allow(userIP)
		slidingAllowed := slidingWindow.Allow(userIP)

		fmt.Printf("Request %d:\n", i)
		fmt.Printf("  Fixed Window:   %s\n", formatResult(fixedAllowed))
		fmt.Printf("  Sliding Window: %s\n", formatResult(slidingAllowed))

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nResult: Both algorithms block after 5 requests (as expected)")
}

// Test the edge case where Fixed Window fails
func testEdgeCase(userIP string) {
	fixedWindow := ratelimiter.NewFixedWindow(5, 5*time.Second)
	slidingWindow := ratelimiter.NewSlidingWindow(5, 5*time.Second)

	fmt.Println("Allow once to start the window...")
	fmt.Println("----------------------------------------------------------")

	fixAllowed := fixedWindow.Allow(userIP)
	slidAllowed := slidingWindow.Allow(userIP)

	fmt.Printf("  Fixed Window:   %s\n", formatResult(fixAllowed))
	fmt.Printf("  Sliding Window: %s\n", formatResult(slidAllowed))

	fmt.Println("\nWaiting 4 seconds first...")
	time.Sleep(4 * time.Second)

	fmt.Println("Scenario: Make 4 other requests near end of window, then 5 more after reset")
	fmt.Println("----------------------------------------------------------")

	fmt.Println("\nTime ~4s: Making 4 requests...")
	for i := 1; i <= 4; i++ {
		fixedWindow.Allow(userIP)
		slidingWindow.Allow(userIP)
	}
	fmt.Println("	Both: 5/5 requests used")

	fmt.Println("\nWaiting 1.5 seconds (window almost over)...")
	time.Sleep(1500 * time.Millisecond)

	fmt.Println("\nTime ~5.5s: Making 5 more requests...")
	fmt.Println("(Fixed Window: Boundary was at 5s - NEW window!)")
	fmt.Println("(Sliding Window: Looks back 5s from 5.5s = [0.5s to 5.5s])")
	fmt.Println("                 Requests from 4s are still visible!)")

	fixedCount := 0
	slidingCount := 0

	for i := 1; i <= 5; i++ {
		fixedAllowed := fixedWindow.Allow(userIP)
		slidingAllowed := slidingWindow.Allow(userIP)

		if fixedAllowed {
			fixedCount++
		}

		if slidingAllowed {
			slidingCount++
		}

		fmt.Printf("Request %d:\n", i)
		fmt.Printf("  Fixed Window:   %s\n", formatResult(fixedAllowed))
		fmt.Printf("  Sliding Window: %s\n", formatResult(slidingAllowed))
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("-", 50))
	fmt.Println("📊 SUMMARY:")
	fmt.Printf("  Fixed Window:   Allowed %d/5 requests in this batch\n", fixedCount)
	fmt.Printf("  Sliding Window: Allowed %d/5 requests in this batch\n", slidingCount)
	fmt.Println()

	if fixedCount == 5 {
		fmt.Println("  ⚠️  Fixed Window: Allowed 10 total requests in ~6 seconds!")
		fmt.Println("      (5 at start + 5 after boundary = 10 in 6 seconds)")
		fmt.Println("      This violates the '5 per 5 seconds' limit!")
	}

	if slidingCount == 1 {
		fmt.Println("  ✅ Sliding Window: Correctly blocked 4/5 requests")
		fmt.Println("      (Looks back 5 seconds and sees the first 5 requests)")
		fmt.Println("      This properly enforces the '5 per 5 seconds' limit!")
	}
}

func formatResult(allowed bool) string {
	if allowed {
		return "✅ ALLOWED"
	}
	return "❌ BLOCKED"
}
