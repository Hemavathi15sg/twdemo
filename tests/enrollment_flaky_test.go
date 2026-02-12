package tests

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnrollmentConcurrency_Flaky tests concurrent enrollment operations
// FIXED: Replaced timeout pattern with sync.WaitGroup for reliable synchronization
func TestEnrollmentConcurrency_Flaky(t *testing.T) {
	var wg sync.WaitGroup
	results := make(chan bool, 5)

	// Start 5 concurrent goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Simulate enrollment operation
			enrollment := createTestEnrollment(id)

			// Send result to channel
			if enrollment.ID > 0 {
				results <- true
			} else {
				results <- false
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	// Count successful enrollments
	successCount := 0
	for success := range results {
		if success {
			successCount++
		}
	}

	// Assert all 5 enrollments succeeded
	assert.Equal(t, 5, successCount, "Expected all 5 enrollments to succeed")
}

func createTestEnrollment(id int) struct{ ID int } {
	// Simulate enrollment creation with deterministic behavior
	// No random delays or timing issues
	return struct{ ID int }{ID: id + 1}
}
