package models

import "time"

// Grade represents a student's grade for a course
type Grade struct {
	ID           string    `json:"id"`
	StudentID    string    `json:"student_id"`
	CourseID     string    `json:"course_id"`
	Grade        string    `json:"grade"` // e.g., "A", "B", "C", etc.
	Score        float64   `json:"score"` // Numerical score
	Semester     string    `json:"semester"`
	AcademicYear string    `json:"academic_year"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateGradeRequest represents the request body for creating a grade
type CreateGradeRequest struct {
	StudentID    string  `json:"student_id"`
	CourseID     string  `json:"course_id"`
	Grade        string  `json:"grade"`
	Score        float64 `json:"score"`
	Semester     string  `json:"semester"`
	AcademicYear string  `json:"academic_year"`
}

// UpdateGradeRequest represents the request body for updating a grade
type UpdateGradeRequest struct {
	Grade        *string  `json:"grade,omitempty"`
	Score        *float64 `json:"score,omitempty"`
	Semester     *string  `json:"semester,omitempty"`
	AcademicYear *string  `json:"academic_year,omitempty"`
}
