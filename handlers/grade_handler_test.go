package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"grademanagement-demo/models"
	"grademanagement-demo/tests"

	"github.com/gorilla/mux"
)

// MockGradeRepository is a mock implementation for testing
type MockGradeRepository struct {
	grades map[int]*models.Grade
	nextID int
}

func NewMockGradeRepository() *MockGradeRepository {
	return &MockGradeRepository{
		grades: make(map[int]*models.Grade),
		nextID: 1,
	}
}

func (m *MockGradeRepository) Create(grade *models.Grade) (*models.Grade, error) {
	grade.ID = m.nextID
	grade.CreatedAt = time.Now()
	grade.UpdatedAt = time.Now()
	m.grades[m.nextID] = grade
	m.nextID++
	return grade, nil
}

func (m *MockGradeRepository) GetByID(id int) (*models.Grade, error) {
	if grade, exists := m.grades[id]; exists {
		return grade, nil
	}
	return nil, nil
}

func (m *MockGradeRepository) GetByStudentAndCourse(studentID, courseID int) (*models.Grade, error) {
	return nil, nil
}

func (m *MockGradeRepository) GetAll() []*models.Grade {
	grades := make([]*models.Grade, 0, len(m.grades))
	for _, grade := range m.grades {
		grades = append(grades, grade)
	}
	return grades
}

func (m *MockGradeRepository) Update(id int, grade *models.Grade) (*models.Grade, error) {
	return grade, nil
}

func (m *MockGradeRepository) Delete(id int) error {
	delete(m.grades, id)
	return nil
}

// MockGradeCache is a mock cache for testing
type MockGradeCache struct {
	cache map[string]*models.Grade
}

func NewMockGradeCache() *MockGradeCache {
	return &MockGradeCache{
		cache: make(map[string]*models.Grade),
	}
}

func (m *MockGradeCache) Set(grade *models.Grade) error {
	return nil
}

func (m *MockGradeCache) GetByID(id int) (*models.Grade, error) {
	return nil, nil
}

func (m *MockGradeCache) Delete(id int) error {
	return nil
}

func (m *MockGradeCache) GetByStudentAndCourse(studentID, courseID int) (*models.Grade, error) {
	return nil, nil
}

func (m *MockGradeCache) SetByStudentAndCourse(grade *models.Grade) error {
	return nil
}

func TestCalculateGrade_Success(t *testing.T) {
	repo := NewMockGradeRepository()
	cache := NewMockGradeCache()
	handler := NewGradeHandler(repo, cache)

	input := tests.NewGradeCalculationInputFactory().Build()
	body, _ := json.Marshal(input)

	req := httptest.NewRequest("POST", "/api/grades/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CalculateGrade(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var result models.Grade
	json.NewDecoder(w.Body).Decode(&result)

	if result.StudentID != input.StudentID {
		t.Errorf("Expected student ID %d, got %d", input.StudentID, result.StudentID)
	}

	if result.LetterGrade == "" {
		t.Error("Expected letter grade to be set")
	}

	if result.GradeColor == "" {
		t.Error("Expected grade color (Figma token) to be set")
	}
}

func TestCalculateGrade_WithCurve(t *testing.T) {
	repo := NewMockGradeRepository()
	cache := NewMockGradeCache()
	handler := NewGradeHandler(repo, cache)

	input := tests.NewGradeCalculationInputFactory().
		WithCurve(5).
		Build()

	body, _ := json.Marshal(input)

	req := httptest.NewRequest("POST", "/api/grades/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CalculateGrade(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var result models.Grade
	json.NewDecoder(w.Body).Decode(&result)

	if !result.CurveApplied {
		t.Error("Expected curve to be applied")
	}

	if result.CurveAmount != 5 {
		t.Errorf("Expected curve amount 5, got %f", result.CurveAmount)
	}

	// Verify curve increased the grade
	if result.NumericGrade <= result.WeightedAverage {
		t.Error("Expected numeric grade to be higher than weighted average after curve")
	}
}

func TestCalculateGrade_FigmaDesignCompliance(t *testing.T) {
	repo := NewMockGradeRepository()
	cache := NewMockGradeCache()
	handler := NewGradeHandler(repo, cache)

	tests := []struct {
		name           string
		assignments    []models.Assignment
		expectedLetter string
		expectedColor  string
		expectedStatus string
	}{
		{
			name: "A+ grade - success token",
			assignments: []models.Assignment{
				{Name: "Test", Score: 98, Weight: 1.0},
			},
			expectedLetter: "A+",
			expectedColor:  "#28A745",
			expectedStatus: "success",
		},
		{
			name: "B grade - info token",
			assignments: []models.Assignment{
				{Name: "Test", Score: 85, Weight: 1.0},
			},
			expectedLetter: "B",
			expectedColor:  "#007BFF",
			expectedStatus: "info",
		},
		{
			name: "C grade - warning token",
			assignments: []models.Assignment{
				{Name: "Test", Score: 75, Weight: 1.0},
			},
			expectedLetter: "C",
			expectedColor:  "#FFC107",
			expectedStatus: "warning",
		},
		{
			name: "F grade - critical token",
			assignments: []models.Assignment{
				{Name: "Test", Score: 55, Weight: 1.0},
			},
			expectedLetter: "F",
			expectedColor:  "#DC3545",
			expectedStatus: "critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tests.NewGradeCalculationInputFactory().
				WithAssignments(tt.assignments).
				Build()

			body, _ := json.Marshal(input)
			req := httptest.NewRequest("POST", "/api/grades/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CalculateGrade(w, req)

			var result models.Grade
			json.NewDecoder(w.Body).Decode(&result)

			if result.LetterGrade != tt.expectedLetter {
				t.Errorf("Expected letter grade %s, got %s", tt.expectedLetter, result.LetterGrade)
			}
			if result.GradeColor != tt.expectedColor {
				t.Errorf("Expected Figma color %s, got %s", tt.expectedColor, result.GradeColor)
			}
			if result.GradeStatus != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.GradeStatus)
			}
		})
	}
}

func TestCalculateGrade_PerformanceRequirement(t *testing.T) {
	repo := NewMockGradeRepository()
	cache := NewMockGradeCache()
	handler := NewGradeHandler(repo, cache)

	// Simulate 100 students
	for i := 0; i < 100; i++ {
		input := tests.NewGradeCalculationInputFactory().
			WithStudentID(1000 + i).
			Build()

		body, _ := json.Marshal(input)

		start := time.Now()
		req := httptest.NewRequest("POST", "/api/grades/calculate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CalculateGrade(w, req)
		elapsed := time.Since(start)

		// Performance requirement: <200ms per request
		if elapsed > 200*time.Millisecond {
			t.Errorf("Performance requirement violated: request %d took %v (expected <200ms)", i, elapsed)
		}

		if w.Code != http.StatusCreated {
			t.Errorf("Request %d failed with status %d", i, w.Code)
		}
	}
}

func TestGetGrade(t *testing.T) {
	repo := NewMockGradeRepository()
	cache := NewMockGradeCache()
	handler := NewGradeHandler(repo, cache)

	// Create a grade first
	grade := tests.NewGradeFactory().Build()
	repo.Create(grade)

	req := httptest.NewRequest("GET", "/api/grades/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetGrade(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCalculateGrade_InvalidInput(t *testing.T) {
	repo := NewMockGradeRepository()
	cache := NewMockGradeCache()
	handler := NewGradeHandler(repo, cache)

	tests := []struct {
		name       string
		input      models.GradeCalculationInput
		wantStatus int
	}{
		{
			name: "Invalid student ID",
			input: models.GradeCalculationInput{
				StudentID: 0,
				CourseID:  2001,
				Assignments: []models.Assignment{
					{Name: "Test", Score: 85, Weight: 1.0},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "No assignments",
			input: models.GradeCalculationInput{
				StudentID:   1001,
				CourseID:    2001,
				Assignments: []models.Assignment{},
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/api/grades/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CalculateGrade(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
