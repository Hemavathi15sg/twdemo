package repository

import (
	"errors"
	"grademanagement-demo/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

// EnrollmentRepository defines the interface for enrollment operations
type EnrollmentRepository interface {
	Create(enrollment *models.Enrollment) (*models.Enrollment, error)
	GetByID(id string) (*models.Enrollment, error)
	GetAll() ([]*models.Enrollment, error)
	Update(id string, enrollment *models.Enrollment) (*models.Enrollment, error)
	Delete(id string) error
}

// InMemoryEnrollmentRepository implements EnrollmentRepository with in-memory storage
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

// Create adds a new enrollment
func (r *InMemoryEnrollmentRepository) Create(enrollment *models.Enrollment) (*models.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate ID if not provided
	if enrollment.ID == "" {
		enrollment.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	enrollment.CreatedAt = now
	enrollment.UpdatedAt = now

	// Set enrollment date if not provided
	if enrollment.EnrollmentDate.IsZero() {
		enrollment.EnrollmentDate = now
	}

	// Validate enrollment
	if err := enrollment.Validate(); err != nil {
		return nil, err
	}

	r.enrollments[enrollment.ID] = enrollment
	return enrollment, nil
}

// GetByID retrieves an enrollment by ID
func (r *InMemoryEnrollmentRepository) GetByID(id string) (*models.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollment, exists := r.enrollments[id]
	if !exists {
		return nil, errors.New("enrollment not found")
	}
	return enrollment, nil
}

// GetAll retrieves all enrollments
func (r *InMemoryEnrollmentRepository) GetAll() ([]*models.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollments := make([]*models.Enrollment, 0, len(r.enrollments))
	for _, enrollment := range r.enrollments {
		enrollments = append(enrollments, enrollment)
	}
	return enrollments, nil
}

// Update updates an existing enrollment
func (r *InMemoryEnrollmentRepository) Update(id string, enrollment *models.Enrollment) (*models.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.enrollments[id]
	if !exists {
		return nil, errors.New("enrollment not found")
	}

	// Validate the updated enrollment
	if err := enrollment.Validate(); err != nil {
		return nil, err
	}

	// Preserve ID, CreatedAt, and EnrollmentDate if not provided
	enrollment.ID = existing.ID
	enrollment.CreatedAt = existing.CreatedAt
	if enrollment.EnrollmentDate.IsZero() {
		enrollment.EnrollmentDate = existing.EnrollmentDate
	}
	enrollment.UpdatedAt = time.Now()

	r.enrollments[id] = enrollment
	return enrollment, nil
}

// Delete removes an enrollment
func (r *InMemoryEnrollmentRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[id]; !exists {
		return errors.New("enrollment not found")
	}

	delete(r.enrollments, id)
	return nil
}
