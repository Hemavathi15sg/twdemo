package models

import "time"

// Grade model with required fields - Figma Design Token Compliance
type Grade struct {
	ID              int       `json:"id"`
	StudentID       int       `json:"student_id"`
	CourseID        int       `json:"course_id"`
	NumericGrade    float64   `json:"numeric_grade"`
	LetterGrade     string    `json:"letter_grade"`
	GradeColor      string    `json:"grade_color"`  // Figma color token
	GradeStatus     string    `json:"grade_status"` // success, warning, critical
	WeightedAverage float64   `json:"weighted_average"`
	CurveApplied    bool      `json:"curve_applied"`
	CurveAmount     float64   `json:"curve_amount,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Assignment represents a single assignment with grade and weight
type Assignment struct {
	Name   string  `json:"name"`
	Score  float64 `json:"score"`  // 0-100
	Weight float64 `json:"weight"` // 0-1 (e.g., 0.3 for 30%)
}

// GradeCalculationInput for calculating grades
type GradeCalculationInput struct {
	StudentID   int          `json:"student_id"`
	CourseID    int          `json:"course_id"`
	Assignments []Assignment `json:"assignments"`
	ApplyCurve  bool         `json:"apply_curve,omitempty"`
	CurveAmount float64      `json:"curve_amount,omitempty"` // Percentage to add (e.g., 5 for +5%)
}

// GradeScale represents grade thresholds from Figma design tokens
type GradeScale struct {
	Letter string
	Min    float64
	Max    float64
	Color  string // Hex color from Figma
	Status string // success, warning, critical, info
}

// Figma Design Tokens - Grade Scale Conventions
var GradeScales = []GradeScale{
	{Letter: "A+", Min: 97, Max: 100, Color: "#28A745", Status: "success"},
	{Letter: "A", Min: 93, Max: 96.99, Color: "#28A745", Status: "success"},
	{Letter: "A-", Min: 90, Max: 92.99, Color: "#28A745", Status: "success"},
	{Letter: "B+", Min: 87, Max: 89.99, Color: "#28A745", Status: "success"},
	{Letter: "B", Min: 83, Max: 86.99, Color: "#007BFF", Status: "info"},
	{Letter: "B-", Min: 80, Max: 82.99, Color: "#007BFF", Status: "info"},
	{Letter: "C+", Min: 77, Max: 79.99, Color: "#007BFF", Status: "info"},
	{Letter: "C", Min: 73, Max: 76.99, Color: "#FFC107", Status: "warning"},
	{Letter: "C-", Min: 70, Max: 72.99, Color: "#FFC107", Status: "warning"},
	{Letter: "D+", Min: 67, Max: 69.99, Color: "#FFC107", Status: "warning"},
	{Letter: "D", Min: 63, Max: 66.99, Color: "#FD7E14", Status: "warning"},
	{Letter: "D-", Min: 60, Max: 62.99, Color: "#FD7E14", Status: "warning"},
	{Letter: "F", Min: 0, Max: 59.99, Color: "#DC3545", Status: "critical"},
}

// CalculateWeightedAverage calculates weighted average from assignments
func CalculateWeightedAverage(assignments []Assignment) (float64, error) {
	if len(assignments) == 0 {
		return 0, nil
	}

	var totalWeighted float64
	var totalWeight float64

	for _, assignment := range assignments {
		totalWeighted += assignment.Score * assignment.Weight
		totalWeight += assignment.Weight
	}

	// Normalize if weights don't sum to 1.0
	if totalWeight > 0 {
		return totalWeighted / totalWeight, nil
	}

	return 0, nil
}

// ApplyGradeCurve adds a percentage curve to the grade
func ApplyGradeCurve(grade float64, curveAmount float64) float64 {
	curved := grade + curveAmount
	if curved > 100 {
		return 100
	}
	return curved
}

// ConvertToLetterGrade converts numeric grade to letter grade using Figma tokens
func ConvertToLetterGrade(numericGrade float64) (string, string, string) {
	for _, scale := range GradeScales {
		if numericGrade >= scale.Min && numericGrade <= scale.Max {
			return scale.Letter, scale.Color, scale.Status
		}
	}
	// Default to F if out of range
	return "F", "#DC3545", "critical"
}

// ValidateGradeInput validates the grade calculation input
func ValidateGradeInput(input GradeCalculationInput) error {
	if input.StudentID <= 0 {
		return &ValidationError{Field: "student_id", Message: "must be a positive integer"}
	}
	if input.CourseID <= 0 {
		return &ValidationError{Field: "course_id", Message: "must be a positive integer"}
	}
	if len(input.Assignments) == 0 {
		return &ValidationError{Field: "assignments", Message: "at least one assignment is required"}
	}

	var totalWeight float64
	for i, assignment := range input.Assignments {
		if assignment.Score < 0 || assignment.Score > 100 {
			return &ValidationError{
				Field:   "assignments[" + string(rune(i)) + "].score",
				Message: "score must be between 0 and 100",
			}
		}
		if assignment.Weight < 0 || assignment.Weight > 1 {
			return &ValidationError{
				Field:   "assignments[" + string(rune(i)) + "].weight",
				Message: "weight must be between 0 and 1",
			}
		}
		totalWeight += assignment.Weight
	}

	// Allow some tolerance for floating point arithmetic
	if totalWeight < 0.99 || totalWeight > 1.01 {
		return &ValidationError{
			Field:   "assignments",
			Message: "total weight must sum to approximately 1.0",
		}
	}

	if input.ApplyCurve && (input.CurveAmount < 0 || input.CurveAmount > 20) {
		return &ValidationError{
			Field:   "curve_amount",
			Message: "curve amount must be between 0 and 20",
		}
	}

	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
