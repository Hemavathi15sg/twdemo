package models

import "time"

// Enrollment represents a student's enrollment in a course
type Enrollment struct {
	ID             string    `json:"id"`
	StudentID      string    `json:"student_id"`
	CourseID       string    `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"` // pending, active, completed
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateEnrollmentRequest represents the request body for creating an enrollment
type CreateEnrollmentRequest struct {
	StudentID      string `json:"student_id"`
	CourseID       string `json:"course_id"`
	EnrollmentDate string `json:"enrollment_date,omitempty"` // Optional, defaults to now
	Status         string `json:"status,omitempty"`          // Optional, defaults to "pending"
}

// UpdateEnrollmentRequest represents the request body for updating an enrollment
type UpdateEnrollmentRequest struct {
	Status string `json:"status,omitempty"`
}

// ValidateStatus checks if the status is valid (pending, active, or completed)
func ValidateStatus(status string) bool {
	return status == "pending" || status == "active" || status == "completed"
}
