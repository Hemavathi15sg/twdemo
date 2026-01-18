package repositories

import (
	"context"
	"errors"
	"grademanagement-demo/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

// InMemoryEnrollmentRepository implements EnrollmentRepository using in-memory storage
type InMemoryEnrollmentRepository struct {
	mu          sync.RWMutex
	enrollments map[string]*models.Enrollment
}

// NewInMemoryEnrollmentRepository creates a new in-memory enrollment repository
func NewInMemoryEnrollmentRepository() *InMemoryEnrollmentRepository {
	return &InMemoryEnrollmentRepository{
		enrollments: make(map[string]*models.Enrollment),
	}
}

// Create creates a new enrollment
func (r *InMemoryEnrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if enrollment.ID == "" {
		enrollment.ID = uuid.New().String()
	}

	now := time.Now()
	enrollment.CreatedAt = now
	enrollment.UpdatedAt = now
	
	if enrollment.EnrollmentDate.IsZero() {
		enrollment.EnrollmentDate = now
	}

	r.enrollments[enrollment.ID] = enrollment
	return enrollment, nil
}

// GetByID retrieves an enrollment by its ID
func (r *InMemoryEnrollmentRepository) GetByID(ctx context.Context, id string) (*models.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollment, exists := r.enrollments[id]
	if !exists {
		return nil, errors.New("enrollment not found")
	}

	return enrollment, nil
}

// Update updates an existing enrollment
func (r *InMemoryEnrollmentRepository) Update(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[enrollment.ID]; !exists {
		return nil, errors.New("enrollment not found")
	}

	enrollment.UpdatedAt = time.Now()
	r.enrollments[enrollment.ID] = enrollment
	return enrollment, nil
}

// Delete deletes an enrollment by its ID
func (r *InMemoryEnrollmentRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[id]; !exists {
		return errors.New("enrollment not found")
	}

	delete(r.enrollments, id)
	return nil
}

// List retrieves all enrollments with optional filtering
func (r *InMemoryEnrollmentRepository) List(ctx context.Context, studentID, courseID string) ([]*models.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Enrollment
	for _, enrollment := range r.enrollments {
		if studentID != "" && enrollment.StudentID != studentID {
			continue
		}
		if courseID != "" && enrollment.CourseID != courseID {
			continue
		}
		result = append(result, enrollment)
	}

	return result, nil
}

// GetByStudentID retrieves all enrollments for a specific student
func (r *InMemoryEnrollmentRepository) GetByStudentID(ctx context.Context, studentID string) ([]*models.Enrollment, error) {
	return r.List(ctx, studentID, "")
}

// GetByCourseID retrieves all enrollments for a specific course
func (r *InMemoryEnrollmentRepository) GetByCourseID(ctx context.Context, courseID string) ([]*models.Enrollment, error) {
	return r.List(ctx, "", courseID)
}
