package models

import (
	"time"

	"github.com/google/uuid"
)

// Enrollment represents a student's enrollment in a course
// JSON tags use snake_case per guidelines; db tags reserved for future persistence layer
type Enrollment struct {
	ID             uuid.UUID `json:"id" db:"id"`
	StudentID      uuid.UUID `json:"student_id" validate:"required" db:"student_id"`
	CourseID       uuid.UUID `json:"course_id" validate:"required" db:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date" db:"enrollment_date"`
	Status         string    `json:"status" validate:"required,oneof=pending active completed" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
