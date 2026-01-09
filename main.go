package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Enrollment model with required fields
type Enrollment struct {
	ID             int       `json:"id"`
	StudentID      int       `json:"student_id"`
	CourseID       int       `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// EnrollmentInput for creating/updating enrollments
type EnrollmentInput struct {
	StudentID      int    `json:"student_id"`
	CourseID       int    `json:"course_id"`
	EnrollmentDate string `json:"enrollment_date,omitempty"`
	Status         string `json:"status"`
}

// ErrorResponse for consistent error handling
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse for consistent success responses
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// EnrollmentRepository manages enrollment data with thread safety
type EnrollmentRepository struct {
	enrollments map[int]*Enrollment
	nextID      int
	mu          sync.RWMutex
}

// NewEnrollmentRepository creates a new repository instance
func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{
		enrollments: make(map[int]*Enrollment),
		nextID:      1,
	}
}

// ValidateStatus checks if status is one of the allowed values
func ValidateStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":   true,
		"active":    true,
		"completed": true,
	}
	return validStatuses[status]
}

// Create adds a new enrollment
func (r *EnrollmentRepository) Create(input EnrollmentInput) (*Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate status
	if !ValidateStatus(input.Status) {
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
	enrollment := &Enrollment{
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
func (r *EnrollmentRepository) GetByID(id int) (*Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollment, exists := r.enrollments[id]
	if !exists {
		return nil, fmt.Errorf("enrollment not found")
	}

	return enrollment, nil
}

// GetAll retrieves all enrollments
func (r *EnrollmentRepository) GetAll() []*Enrollment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollments := make([]*Enrollment, 0, len(r.enrollments))
	for _, enrollment := range r.enrollments {
		enrollments = append(enrollments, enrollment)
	}

	return enrollments
}

// Update modifies an existing enrollment
func (r *EnrollmentRepository) Update(id int, input EnrollmentInput) (*Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	enrollment, exists := r.enrollments[id]
	if !exists {
		return nil, fmt.Errorf("enrollment not found")
	}

	// Validate status if provided
	if input.Status != "" && !ValidateStatus(input.Status) {
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

// EnrollmentHandler handles HTTP requests for enrollments
type EnrollmentHandler struct {
	repo *EnrollmentRepository
}

// NewEnrollmentHandler creates a new handler instance
func NewEnrollmentHandler(repo *EnrollmentRepository) *EnrollmentHandler {
	return &EnrollmentHandler{repo: repo}
}

// CreateEnrollment handles POST /api/enrollments
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var input EnrollmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	enrollment, err := h.repo.Create(input)
	if err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, enrollment, http.StatusCreated)
}

// GetEnrollment handles GET /api/enrollments/{id}
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "invalid enrollment ID", http.StatusBadRequest)
		return
	}

	enrollment, err := h.repo.GetByID(id)
	if err != nil {
		respondError(w, err.Error(), http.StatusNotFound)
		return
	}

	respondJSON(w, enrollment, http.StatusOK)
}

// ListEnrollments handles GET /api/enrollments
func (h *EnrollmentHandler) ListEnrollments(w http.ResponseWriter, r *http.Request) {
	enrollments := h.repo.GetAll()
	respondJSON(w, enrollments, http.StatusOK)
}

// UpdateEnrollment handles PUT /api/enrollments/{id}
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "invalid enrollment ID", http.StatusBadRequest)
		return
	}

	var input EnrollmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	enrollment, err := h.repo.Update(id, input)
	if err != nil {
		if err.Error() == "enrollment not found" {
			respondError(w, err.Error(), http.StatusNotFound)
		} else {
			respondError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	respondJSON(w, enrollment, http.StatusOK)
}

// DeleteEnrollment handles DELETE /api/enrollments/{id}
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "invalid enrollment ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		respondError(w, err.Error(), http.StatusNotFound)
		return
	}

	respondJSON(w, SuccessResponse{Message: "enrollment deleted successfully"}, http.StatusOK)
}

// Helper functions for consistent JSON responses
func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	respondJSON(w, ErrorResponse{Error: message}, statusCode)
}

func main() {
	r := mux.NewRouter()

	// Initialize repository and handler
	repo := NewEnrollmentRepository()
	handler := NewEnrollmentHandler(repo)

	// API routes with /api prefix
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/enrollments", handler.CreateEnrollment).Methods("POST")
	api.HandleFunc("/enrollments", handler.ListEnrollments).Methods("GET")
	api.HandleFunc("/enrollments/{id}", handler.GetEnrollment).Methods("GET")
	api.HandleFunc("/enrollments/{id}", handler.UpdateEnrollment).Methods("PUT")
	api.HandleFunc("/enrollments/{id}", handler.DeleteEnrollment).Methods("DELETE")

	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API - Ready for AI delegation!", "status": "healthy"}`)
	}).Methods("GET")

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Ready for Copilot Agent delegation!")
	fmt.Println("📚 Enrollment API available at /api/enrollments")

	log.Fatal(http.ListenAndServe(port, r))
}