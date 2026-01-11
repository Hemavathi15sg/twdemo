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
	// Seed with current time to get different behavior on each run
	rand.Seed(time.Now().UnixNano())

	// Simulate concurrent enrollment operations
	results := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func(id int) {
			// FLAKY: Random timing to simulate network/database variability
			randomDelay := rand.Intn(15) // 0-15ms random delay
			time.Sleep(time.Duration(randomDelay) * time.Millisecond)

			// Simulate enrollment operation
			enrollment := createTestEnrollment(id)

			// FLAKY: Race condition on shared state
			if enrollment.ID > 0 {
				results <- true
			} else {
				results <- false
			}
		}(i)
	}

	// FLAKY: Not waiting for all goroutines
	successCount := 0
	// FLAKY: Timeout is borderline - will timeout about 40% of the time
	timeout := time.After(60 * time.Millisecond)

	for i := 0; i < 5; i++ {
		select {
		case success := <-results:
			if success {
				successCount++
			}
		case <-timeout:
			// FLAKY: Times out before all complete (sometimes)
			t.Logf("Timeout after %d results", successCount)
			goto done
		}
	}
done:

	// FLAKY: Assertion may fail if timeout occurs early
	assert.Equal(t, 5, successCount, "Expected all 5 enrollments to succeed")
}

func createTestEnrollment(id int) struct{ ID int } {
	// Simulate database operation with variable timing
	randomDelay := rand.Intn(25) // 0-25ms random delay
	time.Sleep(time.Duration(randomDelay) * time.Millisecond)
	return struct{ ID int }{ID: id + 1}
}
