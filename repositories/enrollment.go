package repositories

import (
    "context"
    "errors"

    "grademanagement-demo/models"
    "github.com/google/uuid"
)

// ErrNotFound indicates the requested enrollment does not exist
var ErrNotFound = errors.New("enrollment not found")

// EnrollmentRepository defines CRUD operations for Enrollment aggregate
type EnrollmentRepository interface {
    Create(ctx context.Context, e *models.Enrollment) (*models.Enrollment, error)
    GetByID(ctx context.Context, id uuid.UUID) (*models.Enrollment, error)
    Update(ctx context.Context, e *models.Enrollment) (*models.Enrollment, error)
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context) ([]*models.Enrollment, error)
}
