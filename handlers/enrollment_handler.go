package handlers

import (
	"encoding/json"
	"grademanagement-demo/models"
	"grademanagement-demo/repository"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// EnrollmentHandler handles HTTP requests for enrollment operations
type EnrollmentHandler struct {
	repo *repository.EnrollmentRepository
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(repo *repository.EnrollmentRepository) *EnrollmentHandler {
	return &EnrollmentHandler{
		repo: repo,
	}
}

// CreateEnrollment handles POST requests to create a new enrollment
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.StudentID == "" || req.CourseID == "" {
		http.Error(w, "Missing required fields: student_id or course_id", http.StatusBadRequest)
		return
	}

	// Validate status if provided
	if req.Status != "" && !models.ValidateStatus(req.Status) {
		http.Error(w, "Invalid status. Must be 'pending', 'active', or 'completed'", http.StatusBadRequest)
		return
	}

	// Default status to "pending" if not provided
	status := req.Status
	if status == "" {
		status = "pending"
	}

	// Parse enrollment date or use current time
	var enrollmentDate time.Time
	if req.EnrollmentDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EnrollmentDate)
		if err != nil {
			http.Error(w, "Invalid enrollment_date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		enrollmentDate = parsed
	} else {
		enrollmentDate = time.Now()
	}

	enrollment := &models.Enrollment{
		StudentID:      req.StudentID,
		CourseID:       req.CourseID,
		EnrollmentDate: enrollmentDate,
		Status:         status,
	}

	if err := h.repo.Create(enrollment); err != nil {
		http.Error(w, "Failed to create enrollment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(enrollment)
}

// GetEnrollment handles GET requests to retrieve a specific enrollment
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	enrollment, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Enrollment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollment)
}

// GetAllEnrollments handles GET requests to retrieve all enrollments
func (h *EnrollmentHandler) GetAllEnrollments(w http.ResponseWriter, r *http.Request) {
	enrollments := h.repo.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollments)
}

// UpdateEnrollment handles PUT requests to update an enrollment
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.UpdateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate status if provided
	if req.Status != "" && !models.ValidateStatus(req.Status) {
		http.Error(w, "Invalid status. Must be 'pending', 'active', or 'completed'", http.StatusBadRequest)
		return
	}

	enrollment := &models.Enrollment{
		Status: req.Status,
	}

	if err := h.repo.Update(id, enrollment); err != nil {
		http.Error(w, "Enrollment not found", http.StatusNotFound)
		return
	}

	// Retrieve updated enrollment
	updatedEnrollment, _ := h.repo.GetByID(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedEnrollment)
}

// DeleteEnrollment handles DELETE requests to remove an enrollment
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "Enrollment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Enrollment deleted successfully",
	})
}

// GetEnrollmentsByStudent handles GET requests to retrieve enrollments by student ID
func (h *EnrollmentHandler) GetEnrollmentsByStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]

	enrollments := h.repo.GetByStudentID(studentID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollments)
}

// GetEnrollmentsByCourse handles GET requests to retrieve enrollments by course ID
func (h *EnrollmentHandler) GetEnrollmentsByCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["course_id"]

	enrollments := h.repo.GetByCourseID(courseID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollments)
}
