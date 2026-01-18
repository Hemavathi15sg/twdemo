package usecases

import (
	"context"
	"errors"
	"grademanagement-demo/models"
	"grademanagement-demo/repositories"
)

// EnrollmentUseCase handles business logic for enrollment operations
type EnrollmentUseCase struct {
	repo repositories.EnrollmentRepository
}

// NewEnrollmentUseCase creates a new enrollment use case
func NewEnrollmentUseCase(repo repositories.EnrollmentRepository) *EnrollmentUseCase {
	return &EnrollmentUseCase{
		repo: repo,
	}
}

// CreateEnrollment creates a new enrollment with validation
func (uc *EnrollmentUseCase) CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error) {
	// Validate required fields
	if enrollment.StudentID == "" {
		return nil, errors.New("student_id is required")
	}
	if enrollment.StudentName == "" {
		return nil, errors.New("student_name is required")
	}
	if enrollment.CourseID == "" {
		return nil, errors.New("course_id is required")
	}
	if enrollment.CourseName == "" {
		return nil, errors.New("course_name is required")
	}
	if enrollment.Credits < 1 || enrollment.Credits > 10 {
		return nil, errors.New("credits must be between 1 and 10")
	}

	// Set default status if not provided
	if enrollment.Status == "" {
		enrollment.Status = models.StatusActive
	}

	// Validate status
	if !isValidStatus(enrollment.Status) {
		return nil, errors.New("invalid status: must be active, inactive, completed, or withdrawn")
	}

	return uc.repo.Create(ctx, enrollment)
}

// GetEnrollment retrieves an enrollment by ID
func (uc *EnrollmentUseCase) GetEnrollment(ctx context.Context, id string) (*models.Enrollment, error) {
	if id == "" {
		return nil, errors.New("enrollment id is required")
	}

	return uc.repo.GetByID(ctx, id)
}

// UpdateEnrollment updates an existing enrollment
func (uc *EnrollmentUseCase) UpdateEnrollment(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error) {
	if enrollment.ID == "" {
		return nil, errors.New("enrollment id is required")
	}

	// Validate status if provided
	if enrollment.Status != "" && !isValidStatus(enrollment.Status) {
		return nil, errors.New("invalid status: must be active, inactive, completed, or withdrawn")
	}

	// Validate credits if provided and not zero
	if enrollment.Credits != 0 && (enrollment.Credits < 1 || enrollment.Credits > 10) {
		return nil, errors.New("credits must be between 1 and 10")
	}

	// Verify enrollment exists
	existing, err := uc.repo.GetByID(ctx, enrollment.ID)
	if err != nil {
		return nil, err
	}

	// Merge updates with existing data
	if enrollment.StudentID != "" {
		existing.StudentID = enrollment.StudentID
	}
	if enrollment.StudentName != "" {
		existing.StudentName = enrollment.StudentName
	}
	if enrollment.CourseID != "" {
		existing.CourseID = enrollment.CourseID
	}
	if enrollment.CourseName != "" {
		existing.CourseName = enrollment.CourseName
	}
	if enrollment.Status != "" {
		existing.Status = enrollment.Status
	}
	if enrollment.Grade != nil {
		existing.Grade = enrollment.Grade
	}
	if enrollment.Credits > 0 {
		existing.Credits = enrollment.Credits
	}
	if !enrollment.EnrollmentDate.IsZero() {
		existing.EnrollmentDate = enrollment.EnrollmentDate
	}

	return uc.repo.Update(ctx, existing)
}

// DeleteEnrollment deletes an enrollment by ID
func (uc *EnrollmentUseCase) DeleteEnrollment(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("enrollment id is required")
	}

	return uc.repo.Delete(ctx, id)
}

// ListEnrollments retrieves all enrollments with optional filtering
func (uc *EnrollmentUseCase) ListEnrollments(ctx context.Context, studentID, courseID string) ([]*models.Enrollment, error) {
	return uc.repo.List(ctx, studentID, courseID)
}

// GetEnrollmentsByStudent retrieves all enrollments for a student
func (uc *EnrollmentUseCase) GetEnrollmentsByStudent(ctx context.Context, studentID string) ([]*models.Enrollment, error) {
	if studentID == "" {
		return nil, errors.New("student_id is required")
	}

	return uc.repo.GetByStudentID(ctx, studentID)
}

// GetEnrollmentsByCourse retrieves all enrollments for a course
func (uc *EnrollmentUseCase) GetEnrollmentsByCourse(ctx context.Context, courseID string) ([]*models.Enrollment, error) {
	if courseID == "" {
		return nil, errors.New("course_id is required")
	}

	return uc.repo.GetByCourseID(ctx, courseID)
}

// isValidStatus checks if the status is valid
func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		models.StatusActive:    true,
		models.StatusInactive:  true,
		models.StatusCompleted: true,
		models.StatusWithdrawn: true,
	}
	return validStatuses[status]
}
