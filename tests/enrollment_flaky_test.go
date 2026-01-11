package tests

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEnrollmentConcurrency_Flaky demonstrates a flaky test with race condition
// This test intentionally has timing issues to demonstrate auto-fix capability
func TestEnrollmentConcurrency_Flaky(t *testing.T) {
	// FLAKY: Use nanoseconds to create true randomness between runs
	rand.Seed(time.Now().UnixNano())
	
	// FLAKY: Randomly decide if this run will be slow (50% chance)
	isSlowRun := rand.Intn(2) == 0
	
	// Simulate concurrent enrollment operations
	results := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func(id int) {
			// FLAKY: Simulate network/database variability
			var delay int
			if isSlowRun {
				// Slow run: 30-70ms (will cause timeout)
				delay = 30 + rand.Intn(40)
			} else {
				// Fast run: 5-25ms (should complete in time)
				delay = 5 + rand.Intn(20)
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)

			// Simulate enrollment operation
			enrollment := createTestEnrollment(id, isSlowRun)

			// FLAKY: Race condition on shared state
			if enrollment.ID > 0 {
				results <- true
			} else {
				results <- false
			}
		}(i)
	}

	// FLAKY: Not waiting for all goroutines properly
	successCount := 0
	// FLAKY: Timeout is set for fast operations only
	timeout := time.After(80 * time.Millisecond)

	for i := 0; i < 5; i++ {
		select {
		case success := <-results:
			if success {
				successCount++
			}
		case <-timeout:
			// FLAKY: Times out when run is slow
			t.Logf("Timeout after %d results (isSlowRun=%v)", successCount, isSlowRun)
			goto done
		}
	}
done:

	// FLAKY: Assertion will fail on slow runs
	assert.Equal(t, 5, successCount, "Expected all 5 enrollments to succeed")
}

func createTestEnrollment(id int, isSlowRun bool) struct{ ID int } {
	// Simulate database operation with variable timing
	var delay int
	if isSlowRun {
		delay = rand.Intn(30) // 0-30ms additional delay
	} else {
		delay = rand.Intn(10) // 0-10ms delay
	}
	time.Sleep(time.Duration(delay) * time.Millisecond)
	return struct{ ID int }{ID: id + 1}
}
