package tests

import (
	"grademanagement-demo/models"
)

// GradeFactory creates test grades using the builder pattern
type GradeFactory struct {
	grade *models.Grade
}

// NewGradeFactory creates a new grade factory with defaults
func NewGradeFactory() *GradeFactory {
	return &GradeFactory{
		grade: &models.Grade{
			ID:              1,
			StudentID:       1001,
			CourseID:        2001,
			NumericGrade:    85.5,
			LetterGrade:     "B",
			GradeColor:      "#007BFF",
			GradeStatus:     "info",
			WeightedAverage: 85.5,
			CurveApplied:    false,
			CurveAmount:     0,
		},
	}
}

// WithID sets the grade ID
func (f *GradeFactory) WithID(id int) *GradeFactory {
	f.grade.ID = id
	return f
}

// WithStudentID sets the student ID
func (f *GradeFactory) WithStudentID(studentID int) *GradeFactory {
	f.grade.StudentID = studentID
	return f
}

// WithCourseID sets the course ID
func (f *GradeFactory) WithCourseID(courseID int) *GradeFactory {
	f.grade.CourseID = courseID
	return f
}

// WithNumericGrade sets the numeric grade
func (f *GradeFactory) WithNumericGrade(grade float64) *GradeFactory {
	f.grade.NumericGrade = grade
	// Auto-convert to letter grade
	letter, color, status := models.ConvertToLetterGrade(grade)
	f.grade.LetterGrade = letter
	f.grade.GradeColor = color
	f.grade.GradeStatus = status
	return f
}

// WithCurve applies a curve to the grade
func (f *GradeFactory) WithCurve(amount float64) *GradeFactory {
	f.grade.CurveApplied = true
	f.grade.CurveAmount = amount
	f.grade.NumericGrade = models.ApplyGradeCurve(f.grade.WeightedAverage, amount)
	// Update letter grade after curve
	letter, color, status := models.ConvertToLetterGrade(f.grade.NumericGrade)
	f.grade.LetterGrade = letter
	f.grade.GradeColor = color
	f.grade.GradeStatus = status
	return f
}

// WithAPlus creates an A+ grade (Figma design token)
func (f *GradeFactory) WithAPlus() *GradeFactory {
	f.grade.NumericGrade = 98.0
	f.grade.LetterGrade = "A+"
	f.grade.GradeColor = "#28A745"
	f.grade.GradeStatus = "success"
	return f
}

// WithF creates an F grade (Figma design token)
func (f *GradeFactory) WithF() *GradeFactory {
	f.grade.NumericGrade = 55.0
	f.grade.LetterGrade = "F"
	f.grade.GradeColor = "#DC3545"
	f.grade.GradeStatus = "critical"
	return f
}

// Build returns the constructed grade
func (f *GradeFactory) Build() *models.Grade {
	return f.grade
}

// GradeCalculationInputFactory creates test grade calculation inputs
type GradeCalculationInputFactory struct {
	input *models.GradeCalculationInput
}

// NewGradeCalculationInputFactory creates a new factory with defaults
func NewGradeCalculationInputFactory() *GradeCalculationInputFactory {
	return &GradeCalculationInputFactory{
		input: &models.GradeCalculationInput{
			StudentID: 1001,
			CourseID:  2001,
			Assignments: []models.Assignment{
				{Name: "Midterm", Score: 85, Weight: 0.4},
				{Name: "Final", Score: 90, Weight: 0.4},
				{Name: "Homework", Score: 95, Weight: 0.2},
			},
			ApplyCurve:  false,
			CurveAmount: 0,
		},
	}
}

// WithStudentID sets the student ID
func (f *GradeCalculationInputFactory) WithStudentID(id int) *GradeCalculationInputFactory {
	f.input.StudentID = id
	return f
}

// WithCourseID sets the course ID
func (f *GradeCalculationInputFactory) WithCourseID(id int) *GradeCalculationInputFactory {
	f.input.CourseID = id
	return f
}

// WithAssignments sets custom assignments
func (f *GradeCalculationInputFactory) WithAssignments(assignments []models.Assignment) *GradeCalculationInputFactory {
	f.input.Assignments = assignments
	return f
}

// WithCurve applies a grade curve
func (f *GradeCalculationInputFactory) WithCurve(amount float64) *GradeCalculationInputFactory {
	f.input.ApplyCurve = true
	f.input.CurveAmount = amount
	return f
}

// Build returns the constructed input
func (f *GradeCalculationInputFactory) Build() *models.GradeCalculationInput {
	return f.input
}
