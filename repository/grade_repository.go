package repository

import (
	"fmt"
	"sync"
	"time"

	"grademanagement-demo/models"
)

// GradeRepository manages grade data with thread safety
type GradeRepository struct {
	grades map[int]*models.Grade
	nextID int
	mu     sync.RWMutex
}

// NewGradeRepository creates a new repository instance
func NewGradeRepository() *GradeRepository {
	return &GradeRepository{
		grades: make(map[int]*models.Grade),
		nextID: 1,
	}
}

// Create adds a new grade
func (r *GradeRepository) Create(input models.GradeInput) (*models.Grade, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate status
	if !models.ValidateGradeStatus(input.Status) {
		return nil, fmt.Errorf("invalid status: must be one of draft, submitted, or final")
	}

	// Validate grade value
	if !models.ValidateGradeValue(input.Grade) {
		return nil, fmt.Errorf("invalid grade: must be one of A+, A, A-, B+, B, B-, C+, C, C-, D+, D, D-, F")
	}

	// Validate required fields
	if input.StudentID <= 0 {
		return nil, fmt.Errorf("student_id must be a positive integer")
	}
	if input.CourseID <= 0 {
		return nil, fmt.Errorf("course_id must be a positive integer")
	}

	// Parse grade date or use current time
	var gradeDate time.Time
	var err error
	if input.GradeDate != "" {
		gradeDate, err = time.Parse("2006-01-02", input.GradeDate)
		if err != nil {
			return nil, fmt.Errorf("invalid grade_date format: use YYYY-MM-DD")
		}
	} else {
		gradeDate = time.Now()
	}

	now := time.Now()
	grade := &models.Grade{
		ID:        r.nextID,
		StudentID: input.StudentID,
		CourseID:  input.CourseID,
		Grade:     input.Grade,
		GradeDate: gradeDate,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.grades[r.nextID] = grade
	r.nextID++

	return grade, nil
}

// GetByID retrieves a grade by ID
func (r *GradeRepository) GetByID(id int) (*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	grade, exists := r.grades[id]
	if !exists {
		return nil, fmt.Errorf("grade not found")
	}

	return grade, nil
}

// GetAll retrieves all grades
func (r *GradeRepository) GetAll() []*models.Grade {
	r.mu.RLock()
	defer r.mu.RUnlock()

	grades := make([]*models.Grade, 0, len(r.grades))
	for _, grade := range r.grades {
		grades = append(grades, grade)
	}

	return grades
}

// Update modifies an existing grade
func (r *GradeRepository) Update(id int, input models.GradeInput) (*models.Grade, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	grade, exists := r.grades[id]
	if !exists {
		return nil, fmt.Errorf("grade not found")
	}

	// Validate status if provided
	if input.Status != "" && !models.ValidateGradeStatus(input.Status) {
		return nil, fmt.Errorf("invalid status: must be one of draft, submitted, or final")
	}

	// Validate grade value if provided
	if input.Grade != "" && !models.ValidateGradeValue(input.Grade) {
		return nil, fmt.Errorf("invalid grade: must be one of A+, A, A-, B+, B, B-, C+, C, C-, D+, D, D-, F")
	}

	// Update fields if provided
	if input.StudentID != 0 {
		if input.StudentID <= 0 {
			return nil, fmt.Errorf("student_id must be a positive integer")
		}
		grade.StudentID = input.StudentID
	}
	if input.CourseID != 0 {
		if input.CourseID <= 0 {
			return nil, fmt.Errorf("course_id must be a positive integer")
		}
		grade.CourseID = input.CourseID
	}
	if input.Grade != "" {
		grade.Grade = input.Grade
	}
	if input.GradeDate != "" {
		gradeDate, err := time.Parse("2006-01-02", input.GradeDate)
		if err != nil {
			return nil, fmt.Errorf("invalid grade_date format: use YYYY-MM-DD")
		}
		grade.GradeDate = gradeDate
	}
	if input.Status != "" {
		grade.Status = input.Status
	}

	grade.UpdatedAt = time.Now()

	return grade, nil
}

// Delete removes a grade
func (r *GradeRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.grades[id]; !exists {
		return fmt.Errorf("grade not found")
	}

	delete(r.grades, id)
	return nil
}
