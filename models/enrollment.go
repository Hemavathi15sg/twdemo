package models

import "time"

// Enrollment represents a student's enrollment in a course
type Enrollment struct {
	ID             int64     `json:"id" db:"id"`
	StudentID      int64     `json:"student_id" validate:"required" db:"student_id"`
	CourseID       int64     `json:"course_id" validate:"required" db:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date" db:"enrollment_date"`
	Status         string    `json:"status" validate:"required,oneof=pending active completed" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
