package healthcheck

import (
	"testing"
)

func TestPassiveHealthChecker_RecordFailure(t *testing.T) {
	backends := []*Backend{
		NewBackend("http://localhost:8081"),
	}

	config := PassiveHealthCheckConfig{
		FailureThreshold: 3,
		SuccessThreshold: 2,
	}

	checker := NewPassiveHealthChecker(backends, config)

	// Backend should start healthy
	if !backends[0].IsHealthy() {
		t.Fatal("Backend should start healthy")
	}

	// Record 2 failures (below threshold)
	checker.RecordFailure(backends[0].URL)
	checker.RecordFailure(backends[0].URL)

	if !backends[0].IsHealthy() {
		t.Error("Backend should still be healthy after 2 failures")
	}
}

func TestPassiveHealthChecker_RecordSuccess(t *testing.T) {
	backends := []*Backend{
		NewBackend("http://localhost:8081"),
	}

	config := PassiveHealthCheckConfig{
		FailureThreshold: 3,
		SuccessThreshold: 2,
	}

	checker := NewPassiveHealthChecker(backends, config)

	// Mark backend unhealthy first
	backends[0].MarkUnhealthy()

	// Record 1 success (below threshold)
	checker.RecordSuccess(backends[0].URL)

	if backends[0].IsHealthy() {
		t.Error("Backend should still be unhealthy after 1 success")
	}

	// Record 2nd success (meets threshold)
	checker.RecordSuccess(backends[0].URL)

	if !backends[0].IsHealthy() {
		t.Error("Backend should be healthy after 2 successes")
	}
}

func TestPassiveHealthChecker_FailureResetBySuccess(t *testing.T) {
	backends := []*Backend{
		NewBackend("http://localhost:8081"),
	}

	config := PassiveHealthCheckConfig{
		FailureThreshold: 3,
		SuccessThreshold: 2,
	}

	checker := NewPassiveHealthChecker(backends, config)

	// Record 2 failures
	checker.RecordFailure(backends[0].URL)
	checker.RecordFailure(backends[0].URL)

	// Record 1 success (should reset failure count)
	checker.RecordSuccess(backends[0].URL)

	// Verify failure count was reset
	if checker.GetConsecutiveFailures(backends[0].URL) != 0 {
		t.Error("Consecutive failures should be reset after success")
	}

	// Backend should still be healthy
	if !backends[0].IsHealthy() {
		t.Error("Backend should still be healthy")
	}
}
