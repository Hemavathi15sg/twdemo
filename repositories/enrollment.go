package repositories

import (
	"context"
	"errors"

	"grademanagement-demo/models"
)

// ErrNotFound indicates the requested enrollment does not exist
var ErrNotFound = errors.New("enrollment not found")

// EnrollmentRepository defines CRUD operations for Enrollment aggregate
type EnrollmentRepository interface {
	Create(ctx context.Context, e *models.Enrollment) (*models.Enrollment, error)
	GetByID(ctx context.Context, id int64) (*models.Enrollment, error)
	Update(ctx context.Context, e *models.Enrollment) (*models.Enrollment, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*models.Enrollment, error)
}
