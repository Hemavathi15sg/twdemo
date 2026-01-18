package usecases

import (
	"context"
	"fmt"
	"time"

	"grademanagement-demo/models"
	"grademanagement-demo/repositories"
)

// EnrollmentService handles business logic for enrollments
type EnrollmentService struct {
	repo repositories.EnrollmentRepository
}

// NewEnrollmentService creates a new enrollment service
func NewEnrollmentService(repo repositories.EnrollmentRepository) *EnrollmentService {
	return &EnrollmentService{repo: repo}
}

// CreateEnrollment creates a new enrollment with validation
func (s *EnrollmentService) CreateEnrollment(ctx context.Context, req *CreateEnrollmentRequest) (*models.Enrollment, error) {
	// Validate input
	if err := s.validateEnrollmentRequest(req); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	enrollment := &models.Enrollment{
		StudentID:      req.StudentID,
		CourseID:       req.CourseID,
		EnrollmentDate: req.EnrollmentDate,
		Status:         req.Status,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return s.repo.Create(ctx, enrollment)
}

// GetEnrollmentByID retrieves an enrollment by ID
func (s *EnrollmentService) GetEnrollmentByID(ctx context.Context, id int64) (*models.Enrollment, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid enrollment ID: %d", id)
	}
	return s.repo.GetByID(ctx, id)
}

// UpdateEnrollment updates an existing enrollment
func (s *EnrollmentService) UpdateEnrollment(ctx context.Context, id int64, req *UpdateEnrollmentRequest) (*models.Enrollment, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid enrollment ID: %d", id)
	}

	// Validate input
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get existing enrollment
	enrollment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Status != "" {
		enrollment.Status = req.Status
	}
	if req.EnrollmentDate != nil {
		enrollment.EnrollmentDate = *req.EnrollmentDate
	}
	enrollment.UpdatedAt = time.Now().UTC()

	return s.repo.Update(ctx, enrollment)
}

// DeleteEnrollment deletes an enrollment
func (s *EnrollmentService) DeleteEnrollment(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid enrollment ID: %d", id)
	}
	return s.repo.Delete(ctx, id)
}

// ListEnrollments lists all enrollments
func (s *EnrollmentService) ListEnrollments(ctx context.Context) ([]*models.Enrollment, error) {
	return s.repo.List(ctx)
}

// validateEnrollmentRequest validates create request
func (s *EnrollmentService) validateEnrollmentRequest(req *CreateEnrollmentRequest) error {
	if req.StudentID <= 0 {
		return fmt.Errorf("student_id must be positive")
	}
	if req.CourseID <= 0 {
		return fmt.Errorf("course_id must be positive")
	}
	if req.Status == "" {
		return fmt.Errorf("status is required")
	}
	if !isValidStatus(req.Status) {
		return fmt.Errorf("status must be one of: pending, active, completed")
	}
	if req.EnrollmentDate.IsZero() {
		return fmt.Errorf("enrollment_date is required")
	}
	return nil
}

// validateUpdateRequest validates update request
func (s *EnrollmentService) validateUpdateRequest(req *UpdateEnrollmentRequest) error {
	if req.Status != "" && !isValidStatus(req.Status) {
		return fmt.Errorf("status must be one of: pending, active, completed")
	}
	return nil
}

// isValidStatus checks if status is valid
func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":   true,
		"active":    true,
		"completed": true,
	}
	return validStatuses[status]
}

// CreateEnrollmentRequest is the DTO for creating enrollments
type CreateEnrollmentRequest struct {
	StudentID      int64     `json:"student_id"`
	CourseID       int64     `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
}

// UpdateEnrollmentRequest is the DTO for updating enrollments
type UpdateEnrollmentRequest struct {
	Status         string     `json:"status,omitempty"`
	EnrollmentDate *time.Time `json:"enrollment_date,omitempty"`
}
