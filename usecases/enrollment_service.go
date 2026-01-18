package usecases

import (
    "context"
    "errors"
    "time"

    "grademanagement-demo/models"
    "grademanagement-demo/repositories"
    "github.com/google/uuid"
)

// ErrValidation indicates the input failed business rules
var ErrValidation = errors.New("validation error")

var allowedStatuses = map[string]struct{}{
    "pending":  {},
    "active":   {},
    "completed":{},
}

type EnrollmentService struct {
    repo repositories.EnrollmentRepository
}

func NewEnrollmentService(repo repositories.EnrollmentRepository) *EnrollmentService {
    return &EnrollmentService{repo: repo}
}

func (s *EnrollmentService) Create(ctx context.Context, e *models.Enrollment) (*models.Enrollment, error) {
    if err := validateEnrollment(e); err != nil {
        return nil, err
    }
    now := time.Now().UTC()
    e.ID = uuid.New()
    if e.EnrollmentDate.IsZero() {
        e.EnrollmentDate = now
    }
    e.CreatedAt = now
    e.UpdatedAt = now
    return s.repo.Create(ctx, e)
}

func (s *EnrollmentService) GetByID(ctx context.Context, id uuid.UUID) (*models.Enrollment, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *EnrollmentService) Update(ctx context.Context, e *models.Enrollment) (*models.Enrollment, error) {
    if e.ID == uuid.Nil {
        return nil, ErrValidation
    }
    if err := validateEnrollment(e); err != nil {
        return nil, err
    }
    // Preserve CreatedAt if present, ensure UpdatedAt
    e.UpdatedAt = time.Now().UTC()
    return s.repo.Update(ctx, e)
}

func (s *EnrollmentService) Delete(ctx context.Context, id uuid.UUID) error {
    return s.repo.Delete(ctx, id)
}

func (s *EnrollmentService) List(ctx context.Context) ([]*models.Enrollment, error) {
    return s.repo.List(ctx)
}

func validateEnrollment(e *models.Enrollment) error {
    if e.StudentID == uuid.Nil || e.CourseID == uuid.Nil || e.Status == "" {
        return ErrValidation
    }
    if _, ok := allowedStatuses[e.Status]; !ok {
        return ErrValidation
    }
    return nil
}
