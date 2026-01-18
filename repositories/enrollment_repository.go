package repositories

import (
	"context"
	"errors"
	"sync"
	"time"

	"grademanagement-demo/models"
)

var (
	ErrEnrollmentNotFound = errors.New("enrollment not found")
	ErrInvalidStatus      = errors.New("invalid status: must be pending, active, or completed")
)

// EnrollmentRepository defines the interface for enrollment data operations
type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error)
	GetByID(ctx context.Context, id int64) (*models.Enrollment, error)
	List(ctx context.Context) ([]*models.Enrollment, error)
	Update(ctx context.Context, id int64, enrollment *models.Enrollment) (*models.Enrollment, error)
	Delete(ctx context.Context, id int64) error
}

// InMemoryEnrollmentRepository implements thread-safe in-memory storage
type InMemoryEnrollmentRepository struct {
	mu          sync.RWMutex
	enrollments map[int64]*models.Enrollment
	nextID      int64
}

// NewInMemoryEnrollmentRepository creates a new in-memory enrollment repository
func NewInMemoryEnrollmentRepository() *InMemoryEnrollmentRepository {
	return &InMemoryEnrollmentRepository{
		enrollments: make(map[int64]*models.Enrollment),
		nextID:      1,
	}
}

// Create adds a new enrollment to the repository
func (r *InMemoryEnrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error) {
	// Validate status
	if !models.IsValidStatus(string(enrollment.Status)) {
		return nil, ErrInvalidStatus
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	enrollment.ID = r.nextID
	r.nextID++
	enrollment.CreatedAt = time.Now()
	enrollment.UpdatedAt = time.Now()
	enrollment.EnrollmentDate = time.Now()

	r.enrollments[enrollment.ID] = enrollment
	return enrollment, nil
}

// GetByID retrieves an enrollment by its ID
func (r *InMemoryEnrollmentRepository) GetByID(ctx context.Context, id int64) (*models.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollment, exists := r.enrollments[id]
	if !exists {
		return nil, ErrEnrollmentNotFound
	}

	return enrollment, nil
}

// List retrieves all enrollments
func (r *InMemoryEnrollmentRepository) List(ctx context.Context) ([]*models.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollments := make([]*models.Enrollment, 0, len(r.enrollments))
	for _, enrollment := range r.enrollments {
		enrollments = append(enrollments, enrollment)
	}

	return enrollments, nil
}

// Update modifies an existing enrollment
func (r *InMemoryEnrollmentRepository) Update(ctx context.Context, id int64, enrollment *models.Enrollment) (*models.Enrollment, error) {
	// Validate status if provided
	if enrollment.Status != "" && !models.IsValidStatus(string(enrollment.Status)) {
		return nil, ErrInvalidStatus
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.enrollments[id]
	if !exists {
		return nil, ErrEnrollmentNotFound
	}

	// Update only provided fields
	if enrollment.Status != "" {
		existing.Status = enrollment.Status
	}
	existing.UpdatedAt = time.Now()

	r.enrollments[id] = existing
	return existing, nil
}

// Delete removes an enrollment from the repository
func (r *InMemoryEnrollmentRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[id]; !exists {
		return ErrEnrollmentNotFound
	}

	delete(r.enrollments, id)
	return nil
}
