package repositories

import (
	"context"
	"grademanagement-demo/models"
)

// EnrollmentRepository defines the interface for enrollment data access operations
type EnrollmentRepository interface {
	// Create creates a new enrollment
	Create(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error)
	
	// GetByID retrieves an enrollment by its ID
	GetByID(ctx context.Context, id string) (*models.Enrollment, error)
	
	// Update updates an existing enrollment
	Update(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error)
	
	// Delete deletes an enrollment by its ID
	Delete(ctx context.Context, id string) error
	
	// List retrieves all enrollments with optional filtering
	List(ctx context.Context, studentID, courseID string) ([]*models.Enrollment, error)
	
	// GetByStudentID retrieves all enrollments for a specific student
	GetByStudentID(ctx context.Context, studentID string) ([]*models.Enrollment, error)
	
	// GetByCourseID retrieves all enrollments for a specific course
	GetByCourseID(ctx context.Context, courseID string) ([]*models.Enrollment, error)
}
