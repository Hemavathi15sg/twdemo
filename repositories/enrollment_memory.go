package repositories

import (
    "context"
    "sync"

    "grademanagement-demo/models"
    "github.com/google/uuid"
)

// InMemoryEnrollmentRepository provides a thread-safe in-memory storage for enrollments
type InMemoryEnrollmentRepository struct {
    mu         sync.RWMutex
    store      map[uuid.UUID]*models.Enrollment
}

func NewInMemoryEnrollmentRepository() *InMemoryEnrollmentRepository {
    return &InMemoryEnrollmentRepository{
        store: make(map[uuid.UUID]*models.Enrollment),
    }
}

func (r *InMemoryEnrollmentRepository) Create(ctx context.Context, e *models.Enrollment) (*models.Enrollment, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.store[e.ID] = e
    // return a copy to avoid external mutation
    copy := *e
    return &copy, nil
}

func (r *InMemoryEnrollmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Enrollment, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    if v, ok := r.store[id]; ok {
        copy := *v
        return &copy, nil
    }
    return nil, ErrNotFound
}

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

func (r *InMemoryEnrollmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.store[id]; !ok {
        return ErrNotFound
    }
    delete(r.store, id)
    return nil
}

func (r *InMemoryEnrollmentRepository) List(ctx context.Context) ([]*models.Enrollment, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    res := make([]*models.Enrollment, 0, len(r.store))
    for _, v := range r.store {
        copy := *v
        res = append(res, &copy)
    }
    return res, nil
}
