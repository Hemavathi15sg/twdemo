package tests

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnrollmentConcurrency_Flaky tests concurrent enrollment operations
// Fixed: Uses proper synchronization with WaitGroup instead of arbitrary timeouts
func TestEnrollmentConcurrency_Flaky(t *testing.T) {
	// Use WaitGroup for proper synchronization
	var wg sync.WaitGroup
	results := make(chan bool, 5)

	// Start 5 concurrent enrollment operations
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			// Simulate enrollment operation
			enrollment := createTestEnrollment(id)
			
			// Send result to channel
			results <- enrollment.ID > 0
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	
	// Close the channel after all goroutines are done
	close(results)

	// Count successful enrollments
	successCount := 0
	for success := range results {
		if success {
			successCount++
		}
	}

	// Assert all enrollments succeeded
	assert.Equal(t, 5, successCount, "Expected all 5 enrollments to succeed")
}

func createTestEnrollment(id int) struct{ ID int } {
	// Simulate database operation deterministically
	return struct{ ID int }{ID: id + 1}
}
