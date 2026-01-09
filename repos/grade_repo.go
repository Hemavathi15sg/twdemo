package repos

import (
	"errors"
	"grademanagement-demo/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

// GradeRepository defines the interface for grade data operations
type GradeRepository interface {
	Create(grade *models.Grade) error
	GetByID(id string) (*models.Grade, error)
	GetAll() ([]*models.Grade, error)
	GetByStudentID(studentID string) ([]*models.Grade, error)
	GetByCourseID(courseID string) ([]*models.Grade, error)
	Update(id string, grade *models.Grade) error
	Delete(id string) error
}

// InMemoryGradeRepository implements GradeRepository using in-memory storage
type InMemoryGradeRepository struct {
	grades map[string]*models.Grade
	mu     sync.RWMutex
}

// NewInMemoryGradeRepository creates a new in-memory grade repository
func NewInMemoryGradeRepository() *InMemoryGradeRepository {
	return &InMemoryGradeRepository{
		grades: make(map[string]*models.Grade),
	}
}

// Create adds a new grade to the repository
func (r *InMemoryGradeRepository) Create(grade *models.Grade) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if grade.ID == "" {
		grade.ID = uuid.New().String()
	}
	grade.CreatedAt = time.Now()
	grade.UpdatedAt = time.Now()

	r.grades[grade.ID] = grade
	return nil
}

// GetByID retrieves a grade by its ID
func (r *InMemoryGradeRepository) GetByID(id string) (*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	grade, exists := r.grades[id]
	if !exists {
		return nil, errors.New("grade not found")
	}
	return grade, nil
}

// GetAll retrieves all grades
func (r *InMemoryGradeRepository) GetAll() ([]*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	grades := make([]*models.Grade, 0, len(r.grades))
	for _, grade := range r.grades {
		grades = append(grades, grade)
	}
	return grades, nil
}

// GetByStudentID retrieves all grades for a specific student
func (r *InMemoryGradeRepository) GetByStudentID(studentID string) ([]*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	grades := make([]*models.Grade, 0)
	for _, grade := range r.grades {
		if grade.StudentID == studentID {
			grades = append(grades, grade)
		}
	}
	return grades, nil
}

// GetByCourseID retrieves all grades for a specific course
func (r *InMemoryGradeRepository) GetByCourseID(courseID string) ([]*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	grades := make([]*models.Grade, 0)
	for _, grade := range r.grades {
		if grade.CourseID == courseID {
			grades = append(grades, grade)
		}
	}
	return grades, nil
}

// Update modifies an existing grade
func (r *InMemoryGradeRepository) Update(id string, updatedGrade *models.Grade) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	grade, exists := r.grades[id]
	if !exists {
		return errors.New("grade not found")
	}

	// Update fields
	if updatedGrade.Grade != "" {
		grade.Grade = updatedGrade.Grade
	}
	if updatedGrade.Score != 0 {
		grade.Score = updatedGrade.Score
	}
	if updatedGrade.Semester != "" {
		grade.Semester = updatedGrade.Semester
	}
	if updatedGrade.AcademicYear != "" {
		grade.AcademicYear = updatedGrade.AcademicYear
	}
	grade.UpdatedAt = time.Now()

	return nil
}

// Delete removes a grade from the repository
func (r *InMemoryGradeRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.grades[id]; !exists {
		return errors.New("grade not found")
	}

	delete(r.grades, id)
	return nil
}
