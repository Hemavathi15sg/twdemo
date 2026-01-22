package repository

import (
	"fmt"
	"sync"
	"time"

	"grademanagement-demo/models"
)

//go:generate mockgen -destination=../mocks/mock_grade_repository.go -package=mocks grademanagement-demo/repository GradeRepositoryInterface

// GradeRepositoryInterface defines the contract for grade repository
type GradeRepositoryInterface interface {
	Create(grade *models.Grade) (*models.Grade, error)
	GetByID(id int) (*models.Grade, error)
	GetByStudentAndCourse(studentID, courseID int) (*models.Grade, error)
	GetAll() []*models.Grade
	Update(id int, grade *models.Grade) (*models.Grade, error)
	Delete(id int) error
}

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
func (r *GradeRepository) Create(grade *models.Grade) (*models.Grade, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if grade.StudentID <= 0 {
		return nil, fmt.Errorf("student_id must be a positive integer")
	}
	if grade.CourseID <= 0 {
		return nil, fmt.Errorf("course_id must be a positive integer")
	}

	now := time.Now()
	grade.ID = r.nextID
	grade.CreatedAt = now
	grade.UpdatedAt = now

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

// GetByStudentAndCourse retrieves a grade by student and course
func (r *GradeRepository) GetByStudentAndCourse(studentID, courseID int) (*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, grade := range r.grades {
		if grade.StudentID == studentID && grade.CourseID == courseID {
			return grade, nil
		}
	}

	return nil, fmt.Errorf("grade not found for student %d in course %d", studentID, courseID)
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
func (r *GradeRepository) Update(id int, updatedGrade *models.Grade) (*models.Grade, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	grade, exists := r.grades[id]
	if !exists {
		return nil, fmt.Errorf("grade not found")
	}

	// Update fields
	grade.NumericGrade = updatedGrade.NumericGrade
	grade.LetterGrade = updatedGrade.LetterGrade
	grade.GradeColor = updatedGrade.GradeColor
	grade.GradeStatus = updatedGrade.GradeStatus
	grade.WeightedAverage = updatedGrade.WeightedAverage
	grade.CurveApplied = updatedGrade.CurveApplied
	grade.CurveAmount = updatedGrade.CurveAmount
	grade.UpdatedAt = time.Now()

	return grade, nil
}

// Delete removes a grade by ID
func (r *GradeRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.grades[id]; !exists {
		return fmt.Errorf("grade not found")
	}

	delete(r.grades, id)
	return nil
}
