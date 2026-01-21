package repository

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"grademanagement-demo/models"
)

// EnrollmentRepository manages enrollment data with thread safety
type EnrollmentRepository struct {
	enrollments map[string]*models.Enrollment
	nextID      int
	mu          sync.RWMutex
}

// NewEnrollmentRepository creates a new repository instance
func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{
		enrollments: make(map[string]*models.Enrollment),
		nextID:      1,
	}
}

// Create adds a new enrollment
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate status
	if !models.ValidateStatus(enrollment.Status) {
		return fmt.Errorf("invalid status: must be one of pending, active, or completed")
	}

	// Validate required fields
	if enrollment.StudentID == "" {
		return fmt.Errorf("student_id is required")
	}
	if enrollment.CourseID == "" {
		return fmt.Errorf("course_id is required")
	}

	// Generate ID if not set
	if enrollment.ID == "" {
		enrollment.ID = strconv.Itoa(r.nextID)
		r.nextID++
	}

	now := time.Now()
	enrollment.CreatedAt = now
	enrollment.UpdatedAt = now

	r.enrollments[enrollment.ID] = enrollment

	return nil
}

// GetByID retrieves an enrollment by ID
func (r *EnrollmentRepository) GetByID(id string) (*models.Enrollment, error) {
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
func (r *EnrollmentRepository) Update(id string, enrollment *models.Enrollment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.enrollments[id]
	if !exists {
		return fmt.Errorf("enrollment not found")
	}

	// Validate status if provided
	if enrollment.Status != "" && !models.ValidateStatus(enrollment.Status) {
		return fmt.Errorf("invalid status: must be one of pending, active, or completed")
	}

	// Update fields if provided
	if enrollment.StudentID != "" {
		existing.StudentID = enrollment.StudentID
	}
	if enrollment.CourseID != "" {
		existing.CourseID = enrollment.CourseID
	}
	if enrollment.Status != "" {
		existing.Status = enrollment.Status
	}

	existing.UpdatedAt = time.Now()

	return nil
}

// Delete removes an enrollment
func (r *EnrollmentRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.enrollments[id]; !exists {
		return fmt.Errorf("enrollment not found")
	}

	delete(r.enrollments, id)
	return nil
}

// GetByStudentID retrieves all enrollments for a specific student
func (r *EnrollmentRepository) GetByStudentID(studentID string) []*models.Enrollment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var enrollments []*models.Enrollment
	for _, enrollment := range r.enrollments {
		if enrollment.StudentID == studentID {
			enrollments = append(enrollments, enrollment)
		}
	}

	return enrollments
}

// GetByCourseID retrieves all enrollments for a specific course
func (r *EnrollmentRepository) GetByCourseID(courseID string) []*models.Enrollment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var enrollments []*models.Enrollment
	for _, enrollment := range r.enrollments {
		if enrollment.CourseID == courseID {
			enrollments = append(enrollments, enrollment)
		}
	}

	return enrollments
}
