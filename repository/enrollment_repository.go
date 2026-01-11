package repository

import (
	"fmt"
	"sync"
	"time"

	"grademanagement-demo/models"
)

// EnrollmentRepository manages enrollment data with thread safety
type EnrollmentRepository struct {
	enrollments map[int]*models.Enrollment
	nextID      int
	mu          sync.RWMutex
}

// NewEnrollmentRepository creates a new repository instance
func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{
		enrollments: make(map[int]*models.Enrollment),
		nextID:      1,
	}
}

// Create adds a new enrollment
func (r *EnrollmentRepository) Create(input models.EnrollmentInput) (*models.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate status
	if !models.ValidateStatus(input.Status) {
		return nil, fmt.Errorf("invalid status: must be one of pending, active, or completed")
	}

	// Validate required fields
	if input.StudentID <= 0 {
		return nil, fmt.Errorf("student_id must be a positive integer")
	}
	if input.CourseID <= 0 {
		return nil, fmt.Errorf("course_id must be a positive integer")
	}

	// Parse enrollment date or use current time
	var enrollmentDate time.Time
	var err error
	if input.EnrollmentDate != "" {
		enrollmentDate, err = time.Parse("2006-01-02", input.EnrollmentDate)
		if err != nil {
			return nil, fmt.Errorf("invalid enrollment_date format: use YYYY-MM-DD")
		}
	} else {
		enrollmentDate = time.Now()
	}

	now := time.Now()
	enrollment := &models.Enrollment{
		ID:             r.nextID,
		StudentID:      input.StudentID,
		CourseID:       input.CourseID,
		EnrollmentDate: enrollmentDate,
		Status:         input.Status,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	r.enrollments[r.nextID] = enrollment
	r.nextID++

	return enrollment, nil
}

// GetByID retrieves an enrollment by ID
func (r *EnrollmentRepository) GetByID(id int) (*models.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollment, exists := r.enrollments[id]
	if !exists {
		return nil, fmt.Errorf("enrollment not found")
	}

	return enrollment, nil
}

// GetAll retrieves all enrollments
func (r *EnrollmentRepository) GetAll() []*models.Enrollment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollments := make([]*models.Enrollment, 0, len(r.enrollments))
	for _, enrollment := range r.enrollments {
		enrollments = append(enrollments, enrollment)
	}

	return enrollments
}

// Update modifies an existing enrollment
func (r *EnrollmentRepository) Update(id int, input models.EnrollmentInput) (*models.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	enrollment, exists := r.enrollments[id]
	if !exists {
		return nil, fmt.Errorf("enrollment not found")
	}

	// Validate status if provided
	if input.Status != "" && !models.ValidateStatus(input.Status) {
		return nil, fmt.Errorf("invalid status: must be one of pending, active, or completed")
	}

	// Update fields if provided
	if input.StudentID != 0 {
		if input.StudentID <= 0 {
			return nil, fmt.Errorf("student_id must be a positive integer")
		}
		enrollment.StudentID = input.StudentID
	}
	if input.CourseID != 0 {
		if input.CourseID <= 0 {
			return nil, fmt.Errorf("course_id must be a positive integer")
		}
		enrollment.CourseID = input.CourseID
	}
	if input.EnrollmentDate != "" {
		enrollmentDate, err := time.Parse("2006-01-02", input.EnrollmentDate)
		if err != nil {
			return nil, fmt.Errorf("invalid enrollment_date format: use YYYY-MM-DD")
		}
		enrollment.EnrollmentDate = enrollmentDate
	}
	if input.Status != "" {
		enrollment.Status = input.Status
	}

	enrollment.UpdatedAt = time.Now()

	return enrollment, nil
}

// Delete removes an enrollment
func (r *EnrollmentRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[id]; !exists {
		return fmt.Errorf("enrollment not found")
	}

	delete(r.enrollments, id)
	return nil
}
