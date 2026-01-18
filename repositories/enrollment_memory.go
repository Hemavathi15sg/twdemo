package repositories

import (
	"context"
	"sync"

	"grademanagement-demo/models"
)

// InMemoryEnrollmentRepository provides thread-safe in-memory storage for enrollments
type InMemoryEnrollmentRepository struct {
	mu     sync.RWMutex
	store  map[int64]*models.Enrollment
	nextID int64
}

// NewInMemoryEnrollmentRepository initializes a new in-memory repository with auto-incrementing IDs
func NewInMemoryEnrollmentRepository() *InMemoryEnrollmentRepository {
	return &InMemoryEnrollmentRepository{
		store:  make(map[int64]*models.Enrollment),
		nextID: 1,
	}
}

// Create adds a new enrollment and assigns an auto-incremented ID
func (r *InMemoryEnrollmentRepository) Create(ctx context.Context, e *models.Enrollment) (*models.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e.ID = r.nextID
	r.nextID++
	r.store[e.ID] = e
	// Return a copy to prevent external mutation
	copy := *e
	return &copy, nil
}

// GetByID retrieves an enrollment by ID
func (r *InMemoryEnrollmentRepository) GetByID(ctx context.Context, id int64) (*models.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.store[id]; ok {
		copy := *v
		return &copy, nil
	}
	return nil, ErrNotFound
}

// Update modifies an existing enrollment
func (r *InMemoryEnrollmentRepository) Update(ctx context.Context, e *models.Enrollment) (*models.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[e.ID]; !ok {
		return nil, ErrNotFound
	}
	r.store[e.ID] = e
	copy := *e
	return &copy, nil
}

// Delete removes an enrollment by ID
func (r *InMemoryEnrollmentRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}

// List returns all enrollments
func (r *InMemoryEnrollmentRepository) List(ctx context.Context) ([]*models.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*models.Enrollment, 0, len(r.store))
	for _, v := range r.store {
		copy := *v
		result = append(result, &copy)
	}
	return result, nil
}
