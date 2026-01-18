package models

import (
	"time"
)

// Enrollment represents a student enrollment in a course
type Enrollment struct {
	ID           string    `json:"id" db:"id"`
	StudentID    string    `json:"student_id" db:"student_id" validate:"required"`
	StudentName  string    `json:"student_name" db:"student_name" validate:"required"`
	CourseID     string    `json:"course_id" db:"course_id" validate:"required"`
	CourseName   string    `json:"course_name" db:"course_name" validate:"required"`
	EnrollmentDate time.Time `json:"enrollment_date" db:"enrollment_date"`
	Status       string    `json:"status" db:"status" validate:"required,oneof=active inactive completed withdrawn"`
	Grade        *string   `json:"grade,omitempty" db:"grade"`
	Credits      int       `json:"credits" db:"credits" validate:"required,min=1,max=10"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// EnrollmentStatus constants
const (
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusCompleted = "completed"
	StatusWithdrawn = "withdrawn"
)
