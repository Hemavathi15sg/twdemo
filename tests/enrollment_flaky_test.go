package tests

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

// TestEnrollmentConcurrency_Flaky demonstrates a flaky test with race condition
// This test intentionally has timing issues to demonstrate auto-fix capability
func TestEnrollmentConcurrency_Flaky(t *testing.T) {
	// Simulate concurrent enrollment operations
	results := make(chan bool, 5)
	
	for i := 0; i < 5; i++ {
		go func(id int) {
			// FLAKY: No synchronization, random timing
			time.Sleep(time.Duration(id*10) * time.Millisecond)
			
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
	timeout := time.After(50 * time.Millisecond) // Too short!
	
	for i := 0; i < 5; i++ {
		select {
		case success := <-results:
			if success {
				successCount++
			}
		case <-timeout:
			// FLAKY: Times out before all complete
			t.Logf("Timeout after %d results", successCount)
			break
		}
	}
	
	// FLAKY: Assertion may fail if timeout occurs early
	assert.Equal(t, 5, successCount, "Expected all 5 enrollments to succeed")
}

func createTestEnrollment(id int) struct{ ID int } {
	// Simulate database operation with variable timing
	time.Sleep(time.Duration(id*15) * time.Millisecond)
	return struct{ ID int }{ID: id + 1}
}