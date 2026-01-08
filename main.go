package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Enrollment represents a student enrollment record
type Enrollment struct {
	ID             string    `json:"id"`
	StudentID      string    `json:"student_id"`
	CourseID       string    `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// EnrollmentRepository handles data persistence for enrollments
type EnrollmentRepository interface {
	Create(e Enrollment) (Enrollment, error)
	GetAll() []Enrollment
	GetByID(id string) (Enrollment, bool)
	Update(id string, e Enrollment) (Enrollment, error)
	Delete(id string) error
}

// InMemoryEnrollmentRepository implements EnrollmentRepository with in-memory storage
type InMemoryEnrollmentRepository struct {
	mu          sync.RWMutex
	enrollments map[string]Enrollment
	nextID      int
}

// NewInMemoryEnrollmentRepository creates a new in-memory enrollment repository
func NewInMemoryEnrollmentRepository() *InMemoryEnrollmentRepository {
	return &InMemoryEnrollmentRepository{
		enrollments: make(map[string]Enrollment),
		nextID:      1,
	}
}

// Create adds a new enrollment
func (r *InMemoryEnrollmentRepository) Create(e Enrollment) (Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if e.StudentID == "" {
		return Enrollment{}, fmt.Errorf("student_id is required")
	}
	if e.CourseID == "" {
		return Enrollment{}, fmt.Errorf("course_id is required")
	}
	if e.Status == "" {
		return Enrollment{}, fmt.Errorf("status is required")
	}
	if !isValidStatus(e.Status) {
		return Enrollment{}, fmt.Errorf("status must be one of: pending, active, completed")
	}

	now := time.Now()
	e.ID = fmt.Sprintf("ENR%03d", r.nextID)
	r.nextID++
	e.EnrollmentDate = now
	e.CreatedAt = now
	e.UpdatedAt = now

	r.enrollments[e.ID] = e
	return e, nil
}

// GetAll returns all enrollments
func (r *InMemoryEnrollmentRepository) GetAll() []Enrollment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	enrollments := make([]Enrollment, 0, len(r.enrollments))
	for _, e := range r.enrollments {
		enrollments = append(enrollments, e)
	}
	return enrollments
}

// GetByID returns a specific enrollment
func (r *InMemoryEnrollmentRepository) GetByID(id string) (Enrollment, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	e, ok := r.enrollments[id]
	return e, ok
}

// Update updates an existing enrollment
func (r *InMemoryEnrollmentRepository) Update(id string, e Enrollment) (Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.enrollments[id]
	if !ok {
		return Enrollment{}, fmt.Errorf("enrollment not found")
	}

	if e.StudentID != "" {
		existing.StudentID = e.StudentID
	}
	if e.CourseID != "" {
		existing.CourseID = e.CourseID
	}
	if e.Status != "" {
		if !isValidStatus(e.Status) {
			return Enrollment{}, fmt.Errorf("status must be one of: pending, active, completed")
		}
		existing.Status = e.Status
	}

	existing.UpdatedAt = time.Now()
	r.enrollments[id] = existing
	return existing, nil
}

// Delete removes an enrollment
func (r *InMemoryEnrollmentRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.enrollments[id]; !ok {
		return fmt.Errorf("enrollment not found")
	}

	delete(r.enrollments, id)
	return nil
}

// isValidStatus validates enrollment status
func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":   true,
		"active":    true,
		"completed": true,
	}
	return validStatuses[status]
}

// EnrollmentHandler handles enrollment API requests
type EnrollmentHandler struct {
	repo EnrollmentRepository
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(repo EnrollmentRepository) *EnrollmentHandler {
	return &EnrollmentHandler{repo: repo}
}

// CreateEnrollment handles POST /api/enrollments
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var enrollment Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, `{"error": "Invalid request format"}`, http.StatusBadRequest)
		return
	}

	created, err := h.repo.Create(enrollment)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetAllEnrollments handles GET /api/enrollments
func (h *EnrollmentHandler) GetAllEnrollments(w http.ResponseWriter, r *http.Request) {
	enrollments := h.repo.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollments)
}

// GetEnrollment handles GET /api/enrollments/{id}
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	enrollment, ok := h.repo.GetByID(id)
	if !ok {
		http.Error(w, `{"error": "Enrollment not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollment)
}

// UpdateEnrollment handles PUT /api/enrollments/{id}
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var enrollment Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, `{"error": "Invalid request format"}`, http.StatusBadRequest)
		return
	}

	updated, err := h.repo.Update(id, enrollment)
	if err != nil {
		// Check if it's a validation error or not found error
		if err.Error() == "enrollment not found" {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteEnrollment handles DELETE /api/enrollments/{id}
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()

	// Initialize enrollment repository and handler
	repo := NewInMemoryEnrollmentRepository()
	handler := NewEnrollmentHandler(repo)

	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API - Ready for AI delegation!", "status": "healthy"}`)
	}).Methods("GET")

	// API subrouter with /api prefix
	api := r.PathPrefix("/api").Subrouter()

	// Enrollment API endpoints
	api.HandleFunc("/enrollments", handler.CreateEnrollment).Methods("POST")
	api.HandleFunc("/enrollments", handler.GetAllEnrollments).Methods("GET")
	api.HandleFunc("/enrollments/{id}", handler.GetEnrollment).Methods("GET")
	api.HandleFunc("/enrollments/{id}", handler.UpdateEnrollment).Methods("PUT")
	api.HandleFunc("/enrollments/{id}", handler.DeleteEnrollment).Methods("DELETE")

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Ready for Copilot Agent delegation!")
	fmt.Println("📚 Enrollment API endpoints:")
	fmt.Println("   POST   /api/enrollments       - Create new enrollment")
	fmt.Println("   GET    /api/enrollments       - Get all enrollments")
	fmt.Println("   GET    /api/enrollments/{id}  - Get enrollment by ID")
	fmt.Println("   PUT    /api/enrollments/{id}  - Update enrollment")
	fmt.Println("   DELETE /api/enrollments/{id}  - Delete enrollment")

	log.Fatal(http.ListenAndServe(port, r))
}