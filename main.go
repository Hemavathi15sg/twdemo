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
	StudentName    string    `json:"student_name"`
	Course         string    `json:"course"`
	Grade          string    `json:"grade,omitempty"`
	EnrollmentDate time.Time `json:"enrollment_date"`
}

// EnrollmentStore handles in-memory storage for enrollments
type EnrollmentStore struct {
	mu          sync.RWMutex
	enrollments map[string]Enrollment
	nextID      int
}

// NewEnrollmentStore creates a new enrollment store
func NewEnrollmentStore() *EnrollmentStore {
	return &EnrollmentStore{
		enrollments: make(map[string]Enrollment),
		nextID:      1,
	}
}

// Create adds a new enrollment
func (s *EnrollmentStore) Create(e Enrollment) (Enrollment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if e.StudentName == "" {
		return Enrollment{}, fmt.Errorf("student_name is required")
	}
	if e.Course == "" {
		return Enrollment{}, fmt.Errorf("course is required")
	}
	
	e.ID = fmt.Sprintf("ENR%03d", s.nextID)
	s.nextID++
	e.EnrollmentDate = time.Now()
	
	s.enrollments[e.ID] = e
	return e, nil
}

// GetAll returns all enrollments
func (s *EnrollmentStore) GetAll() []Enrollment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	enrollments := make([]Enrollment, 0, len(s.enrollments))
	for _, e := range s.enrollments {
		enrollments = append(enrollments, e)
	}
	return enrollments
}

// GetByID returns a specific enrollment
func (s *EnrollmentStore) GetByID(id string) (Enrollment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	e, ok := s.enrollments[id]
	return e, ok
}

// Update updates an existing enrollment
func (s *EnrollmentStore) Update(id string, e Enrollment) (Enrollment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	existing, ok := s.enrollments[id]
	if !ok {
		return Enrollment{}, fmt.Errorf("enrollment not found")
	}
	
	if e.StudentName != "" {
		existing.StudentName = e.StudentName
	}
	if e.Course != "" {
		existing.Course = e.Course
	}
	if e.Grade != "" {
		existing.Grade = e.Grade
	}
	
	s.enrollments[id] = existing
	return existing, nil
}

// Delete removes an enrollment
func (s *EnrollmentStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, ok := s.enrollments[id]; !ok {
		return fmt.Errorf("enrollment not found")
	}
	
	delete(s.enrollments, id)
	return nil
}

// EnrollmentHandler handles enrollment API requests
type EnrollmentHandler struct {
	store *EnrollmentStore
}

// CreateEnrollment handles POST /enrollments
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var enrollment Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request body: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	
	created, err := h.store.Create(enrollment)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetAllEnrollments handles GET /enrollments
func (h *EnrollmentHandler) GetAllEnrollments(w http.ResponseWriter, r *http.Request) {
	enrollments := h.store.GetAll()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollments)
}

// GetEnrollment handles GET /enrollments/{id}
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	enrollment, ok := h.store.GetByID(id)
	if !ok {
		http.Error(w, `{"error": "Enrollment not found"}`, http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollment)
}

// UpdateEnrollment handles PUT /enrollments/{id}
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	var enrollment Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request body: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	
	updated, err := h.store.Update(id, enrollment)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteEnrollment handles DELETE /enrollments/{id}
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	if err := h.store.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()
	
	// Initialize enrollment store and handler
	store := NewEnrollmentStore()
	handler := &EnrollmentHandler{store: store}
	
	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API - Ready for AI delegation!", "status": "healthy"}`)
	}).Methods("GET")
	
	// Enrollment API endpoints
	r.HandleFunc("/enrollments", handler.CreateEnrollment).Methods("POST")
	r.HandleFunc("/enrollments", handler.GetAllEnrollments).Methods("GET")
	r.HandleFunc("/enrollments/{id}", handler.GetEnrollment).Methods("GET")
	r.HandleFunc("/enrollments/{id}", handler.UpdateEnrollment).Methods("PUT")
	r.HandleFunc("/enrollments/{id}", handler.DeleteEnrollment).Methods("DELETE")

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Ready for Copilot Agent delegation!")
	fmt.Println("📚 Enrollment API endpoints:")
	fmt.Println("   POST   /enrollments       - Create new enrollment")
	fmt.Println("   GET    /enrollments       - Get all enrollments")
	fmt.Println("   GET    /enrollments/{id}  - Get enrollment by ID")
	fmt.Println("   PUT    /enrollments/{id}  - Update enrollment")
	fmt.Println("   DELETE /enrollments/{id}  - Delete enrollment")
	
	log.Fatal(http.ListenAndServe(port, r))
}