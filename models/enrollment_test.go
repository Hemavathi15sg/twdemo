package models

import (
	"testing"
	"time"
)

func TestEnrollmentModel(t *testing.T) {
	t.Run("Create enrollment with all fields", func(t *testing.T) {
		grade := "A"
		enrollment := Enrollment{
			ID:             "test-id-1",
			StudentID:      "student-123",
			StudentName:    "John Doe",
			CourseID:       "course-456",
			CourseName:     "Introduction to Go",
			EnrollmentDate: time.Now(),
			Status:         StatusActive,
			Grade:          &grade,
			Credits:        3,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if enrollment.StudentID != "student-123" {
			t.Errorf("Expected StudentID to be 'student-123', got '%s'", enrollment.StudentID)
		}
		if enrollment.Status != StatusActive {
			t.Errorf("Expected Status to be '%s', got '%s'", StatusActive, enrollment.Status)
		}
		if enrollment.Credits != 3 {
			t.Errorf("Expected Credits to be 3, got %d", enrollment.Credits)
		}
		if *enrollment.Grade != "A" {
			t.Errorf("Expected Grade to be 'A', got '%s'", *enrollment.Grade)
		}
	})

	t.Run("Status constants are correct", func(t *testing.T) {
		if StatusActive != "active" {
			t.Errorf("Expected StatusActive to be 'active', got '%s'", StatusActive)
		}
		if StatusInactive != "inactive" {
			t.Errorf("Expected StatusInactive to be 'inactive', got '%s'", StatusInactive)
		}
		if StatusCompleted != "completed" {
			t.Errorf("Expected StatusCompleted to be 'completed', got '%s'", StatusCompleted)
		}
		if StatusWithdrawn != "withdrawn" {
			t.Errorf("Expected StatusWithdrawn to be 'withdrawn', got '%s'", StatusWithdrawn)
		}
	})
}
