package main

import (
	"fmt"
	"time"

	"rate-limiter/ratelimiter"
)

func test_fixed_window() {
	fmt.Println("Rate Limiter Test Program")
	fmt.Println("-------------------------")

	limiter := ratelimiter.NewFixedWindow(5, 10*time.Second)

	userIP := "192.168.1.1"

	fmt.Printf("Configuration: %d requests per 10 seconds\n", 5)
	fmt.Printf("Testing with IP: %s\n\n", userIP)

	fmt.Println("Making 8 rapid requests:")
	fmt.Println("------------------------")

	for i := 1; i <= 8; i++ {
		allowed := limiter.Allow(userIP)

		if allowed {
			fmt.Printf("Request %d: Allowed\n", i)
		} else {
			fmt.Printf("Request %d: Denied (Rate limit exceeded)\n", i)
		}

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nWaiting for 11 seconds to reset the window...")
	time.Sleep(11 * time.Second)

	fmt.Println("\nAfter window reset:")
	fmt.Println("---------------------")

	for i := 1; i <= 3; i++ {
		allowed := limiter.Allow(userIP)

		if allowed {
			fmt.Printf("Request %d: Allowed (fresh window!)\n", i)
		} else {
			fmt.Printf("Request %d: Denied\n", i)
		}
	}

	fmt.Println("\nTest completed!")
}
