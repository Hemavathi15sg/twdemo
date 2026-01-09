package models

import (
	"errors"
	"time"
)

// Enrollment represents a student enrollment in a course
type Enrollment struct {
	ID             string    `json:"id"`
	StudentID      string    `json:"student_id"`
	CourseID       string    `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Valid status values
const (
	StatusPending   = "pending"
	StatusActive    = "active"
	StatusCompleted = "completed"
)

// ValidateStatus checks if the status is valid
func ValidateStatus(status string) error {
	switch status {
	case StatusPending, StatusActive, StatusCompleted:
		return nil
	default:
		return errors.New("status must be one of: pending, active, completed")
	}
}

// Validate validates the enrollment model
func (e *Enrollment) Validate() error {
	if e.StudentID == "" {
		return errors.New("student_id is required")
	}
	if e.CourseID == "" {
		return errors.New("course_id is required")
	}
	if e.Status == "" {
		return errors.New("status is required")
	}
	if err := ValidateStatus(e.Status); err != nil {
		return err
	}
	return nil
}
