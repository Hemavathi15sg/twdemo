package models

import "time"

// Enrollment model with required fields
type Enrollment struct {
	ID             int       `json:"id"`
	StudentID      int       `json:"student_id"`
	CourseID       int       `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// EnrollmentInput for creating/updating enrollments
type EnrollmentInput struct {
	StudentID      int    `json:"student_id"`
	CourseID       int    `json:"course_id"`
	EnrollmentDate string `json:"enrollment_date,omitempty"`
	Status         string `json:"status"`
}

// ValidateStatus checks if status is one of the allowed values
func ValidateStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":   true,
		"active":    true,
		"completed": true,
	}
	return validStatuses[status]
}
