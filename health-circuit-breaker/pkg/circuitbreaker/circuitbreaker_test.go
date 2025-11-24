package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_SuccessfulRequests(t *testing.T) {
	cb := NewCircuitBreaker(Config{
		Name:        "test",
		MaxRequests: 1,
		Timeout:     10 * time.Millisecond,
	})

	// All successful requests should pass
	for i := 0; i < 10; i++ {
		_, err := cb.Execute(func() (interface{}, error) {
			return "success", nil
		})
		if err != nil {
			t.Errorf("Request %d failed: %v", i, err)
		}
	}

	if cb.State() != StateClosed {
		t.Errorf("Expected stateClosed, got %s", cb.State())
	}
}

func TestCircuitBreaker_OpenAfterFailures(t *testing.T) {
	cb := NewCircuitBreaker(Config{
		Name:        "test",
		MaxRequests: 1,
		Timeout:     100 * time.Millisecond,
		ReadyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	})

	// Circuit should be closed initially
	if cb.State() != StateClosed {
		t.Errorf("Expected initial state Closed, got %s", cb.State())
	}

	// Trigger 3 failures
	for i := 0; i < 3; i++ {
		cb.Execute(func() (interface{}, error) {
			return nil, errors.New("failed")
		})
	}

	// Circuit should now be open
	if cb.State() != StateOpen {
		t.Errorf("Expected state Open after failures, got %s", cb.State())
	}

	// Requests should fail immediately
	_, err := cb.Execute(func() (interface{}, error) {
		return "should not execute", nil
	})

	if err != ErrCircuitOpen {
		t.Errorf("Expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_HalfOpenTransition(t *testing.T) {
	cb := NewCircuitBreaker(Config{
		Name:        "test",
		MaxRequests: 1,
		Timeout:     50 * time.Millisecond,
		ReadyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures >= 2
		},
	})

	// Trigger failures to open circuit
	for i := 0; i < 2; i++ {
		cb.Execute(func() (interface{}, error) {
			return nil, errors.New("failed")
		})
	}

	if cb.State() != StateOpen {
		t.Errorf("Expected state Open, got %s", cb.State())
	}

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Next request should trigger half-open state
	_, err := cb.Execute(func() (interface{}, error) {
		return "success", nil
	})

	if err != nil {
		t.Errorf("Half-open request should succeed, got error: %v", err)
	}

	// After successful reqeust in half-open, should close
	if cb.State() != StateClosed {
		t.Errorf("Expected state Closed after successful half-open, got %s", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenBackToOpen(t *testing.T) {
	cb := NewCircuitBreaker(Config{
		Name:        "test",
		MaxRequests: 1,
		Timeout:     50 * time.Millisecond,
		ReadyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures >= 2
		},
	})

	// Open the circuit
	for i := 0; i < 2; i++ {
		cb.Execute(func() (interface{}, error) {
			return nil, errors.New("failed")
		})
	}

	// Wait for timeout to enter half-open
	time.Sleep(60 * time.Millisecond)

	// Fail the test request in half-open
	cb.Execute(func() (interface{}, error) {
		return nil, errors.New("still failing")
	})

	// Should go back to open
	if cb.State() != StateOpen {
		t.Errorf("Expected state Open after failed half-open request, got %s", cb.State())
	}
}

func TestCircuitBreaker_TooManyRequestsInHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(Config{
		Name:        "test",
		MaxRequests: 2,
		Timeout:     50 * time.Millisecond,
		ReadyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures >= 2
		},
	})

	// Open the circuit
	for i := 0; i < 2; i++ {
		cb.Execute(func() (interface{}, error) {
			return nil, errors.New("failed")
		})
	}

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Start 2 slow requests in half-open (don't complete yet)
	done := make(chan bool, 2)

	go func() {
		cb.Execute(func() (interface{}, error) {
			time.Sleep(100 * time.Millisecond) // Slow request
			return "ok", nil
		})
		done <- true
	}()

	go func() {
		cb.Execute(func() (interface{}, error) {
			time.Sleep(100 * time.Millisecond) // Slow request
			return "ok", nil
		})
		done <- true
	}()

	// Give goroutines time to start
	time.Sleep(10 * time.Millisecond)

	// Third request should be rejected9while first 2 are still running
	_, err := cb.Execute(func() (interface{}, error) {
		return "should not execute", nil
	})

	if err != ErrTooManyRequests {
		t.Errorf("Expected ErrTooManyRequests, got %v", err)
	}

	// Wait for slow requests to complete
	<-done
	<-done
}
