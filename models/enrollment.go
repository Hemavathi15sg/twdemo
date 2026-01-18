package models

import (
	"time"
)

// EnrollmentStatus represents the valid status values for enrollment
type EnrollmentStatus string

const (
	StatusPending   EnrollmentStatus = "pending"
	StatusActive    EnrollmentStatus = "active"
	StatusCompleted EnrollmentStatus = "completed"
)

// Enrollment represents a student enrollment in a course
type Enrollment struct {
	ID             int64            `json:"id" db:"id"`
	StudentID      int64            `json:"student_id" db:"student_id" validate:"required"`
	CourseID       int64            `json:"course_id" db:"course_id" validate:"required"`
	EnrollmentDate time.Time        `json:"enrollment_date" db:"enrollment_date"`
	Status         EnrollmentStatus `json:"status" db:"status" validate:"required,oneof=pending active completed"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
}

// IsValidStatus checks if the given status is valid
func IsValidStatus(status string) bool {
	switch EnrollmentStatus(status) {
	case StatusPending, StatusActive, StatusCompleted:
		return true
	default:
		return false
	}
}

// CreateEnrollmentRequest represents the request body for creating an enrollment
type CreateEnrollmentRequest struct {
	StudentID int64  `json:"student_id" validate:"required"`
	CourseID  int64  `json:"course_id" validate:"required"`
	Status    string `json:"status" validate:"required"`
}

// UpdateEnrollmentRequest represents the request body for updating an enrollment
type UpdateEnrollmentRequest struct {
	Status *string `json:"status,omitempty"`
}
