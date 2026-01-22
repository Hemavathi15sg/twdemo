package models

import (
	"testing"
)

func TestCalculateWeightedAverage(t *testing.T) {
	tests := []struct {
		name        string
		assignments []Assignment
		expected    float64
		wantErr     bool
	}{
		{
			name: "Standard weighted average",
			assignments: []Assignment{
				{Name: "Midterm", Score: 85, Weight: 0.4},
				{Name: "Final", Score: 90, Weight: 0.4},
				{Name: "Homework", Score: 95, Weight: 0.2},
			},
			expected: 88.5,
			wantErr:  false,
		},
		{
			name: "All same weight",
			assignments: []Assignment{
				{Name: "Test1", Score: 80, Weight: 0.33},
				{Name: "Test2", Score: 90, Weight: 0.33},
				{Name: "Test3", Score: 100, Weight: 0.34},
			},
			expected: 89.8,
			wantErr:  false,
		},
		{
			name:        "Empty assignments",
			assignments: []Assignment{},
			expected:    0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateWeightedAverage(tt.assignments)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateWeightedAverage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !almostEqual(got, tt.expected, 0.5) {
				t.Errorf("CalculateWeightedAverage() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestApplyGradeCurve(t *testing.T) {
	tests := []struct {
		name        string
		grade       float64
		curveAmount float64
		expected    float64
	}{
		{
			name:        "Apply 5% curve",
			grade:       85,
			curveAmount: 5,
			expected:    90,
		},
		{
			name:        "Curve capped at 100",
			grade:       98,
			curveAmount: 5,
			expected:    100,
		},
		{
			name:        "No curve",
			grade:       85,
			curveAmount: 0,
			expected:    85,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyGradeCurve(tt.grade, tt.curveAmount)
			if got != tt.expected {
				t.Errorf("ApplyGradeCurve() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConvertToLetterGrade_FigmaCompliance(t *testing.T) {
	tests := []struct {
		name           string
		numericGrade   float64
		expectedLetter string
		expectedColor  string
		expectedStatus string
	}{
		{
			name:           "A+ grade - Figma success token",
			numericGrade:   98,
			expectedLetter: "A+",
			expectedColor:  "#28A745",
			expectedStatus: "success",
		},
		{
			name:           "A grade - Figma success token",
			numericGrade:   95,
			expectedLetter: "A",
			expectedColor:  "#28A745",
			expectedStatus: "success",
		},
		{
			name:           "B grade - Figma info token",
			numericGrade:   85,
			expectedLetter: "B",
			expectedColor:  "#007BFF",
			expectedStatus: "info",
		},
		{
			name:           "C grade - Figma warning token",
			numericGrade:   75,
			expectedLetter: "C",
			expectedColor:  "#FFC107",
			expectedStatus: "warning",
		},
		{
			name:           "F grade - Figma critical token",
			numericGrade:   55,
			expectedLetter: "F",
			expectedColor:  "#DC3545",
			expectedStatus: "critical",
		},
		{
			name:           "Boundary A+ min",
			numericGrade:   97,
			expectedLetter: "A+",
			expectedColor:  "#28A745",
			expectedStatus: "success",
		},
		{
			name:           "Boundary A max",
			numericGrade:   96.99,
			expectedLetter: "A",
			expectedColor:  "#28A745",
			expectedStatus: "success",
		},
		{
			name:           "Boundary F max",
			numericGrade:   59.99,
			expectedLetter: "F",
			expectedColor:  "#DC3545",
			expectedStatus: "critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			letter, color, status := ConvertToLetterGrade(tt.numericGrade)
			if letter != tt.expectedLetter {
				t.Errorf("ConvertToLetterGrade() letter = %v, want %v", letter, tt.expectedLetter)
			}
			if color != tt.expectedColor {
				t.Errorf("ConvertToLetterGrade() color = %v, want %v (Figma token)", color, tt.expectedColor)
			}
			if status != tt.expectedStatus {
				t.Errorf("ConvertToLetterGrade() status = %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestValidateGradeInput(t *testing.T) {
	tests := []struct {
		name    string
		input   GradeCalculationInput
		wantErr bool
	}{
		{
			name: "Valid input",
			input: GradeCalculationInput{
				StudentID: 1001,
				CourseID:  2001,
				Assignments: []Assignment{
					{Name: "Test", Score: 85, Weight: 1.0},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid student ID",
			input: GradeCalculationInput{
				StudentID: 0,
				CourseID:  2001,
				Assignments: []Assignment{
					{Name: "Test", Score: 85, Weight: 1.0},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid score",
			input: GradeCalculationInput{
				StudentID: 1001,
				CourseID:  2001,
				Assignments: []Assignment{
					{Name: "Test", Score: 105, Weight: 1.0},
				},
			},
			wantErr: true,
		},
		{
			name: "Weights don't sum to 1",
			input: GradeCalculationInput{
				StudentID: 1001,
				CourseID:  2001,
				Assignments: []Assignment{
					{Name: "Test", Score: 85, Weight: 0.5},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGradeInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGradeInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func almostEqual(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}
