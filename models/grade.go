package models

import "time"

// Grade model with required fields
type Grade struct {
	ID         int       `json:"id"`
	StudentID  int       `json:"student_id"`
	CourseID   int       `json:"course_id"`
	Grade      string    `json:"grade"`
	GradeDate  time.Time `json:"grade_date"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GradeInput for creating/updating grades
type GradeInput struct {
	StudentID int    `json:"student_id"`
	CourseID  int    `json:"course_id"`
	Grade     string `json:"grade"`
	GradeDate string `json:"grade_date,omitempty"`
	Status    string `json:"status"`
}

// ValidateGradeStatus checks if status is one of the allowed values
func ValidateGradeStatus(status string) bool {
	validStatuses := map[string]bool{
		"draft":     true,
		"submitted": true,
		"final":     true,
	}
	return validStatuses[status]
}

// ValidateGradeValue checks if grade value is valid
func ValidateGradeValue(grade string) bool {
	validGrades := map[string]bool{
		"A+": true, "A": true, "A-": true,
		"B+": true, "B": true, "B-": true,
		"C+": true, "C": true, "C-": true,
		"D+": true, "D": true, "D-": true,
		"F": true,
	}
	return validGrades[grade]
}
