package handlers

import (
	"encoding/json"
	"grademanagement-demo/models"
	"grademanagement-demo/repository"
	"net/http"

	"github.com/gorilla/mux"
)

// EnrollmentHandler handles enrollment-related HTTP requests
type EnrollmentHandler struct {
	repo repository.EnrollmentRepository
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(repo repository.EnrollmentRepository) *EnrollmentHandler {
	return &EnrollmentHandler{repo: repo}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// CreateEnrollment handles POST /api/enrollments
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var enrollment models.Enrollment

	// Decode request body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&enrollment); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create enrollment
	created, err := h.repo.Create(&enrollment)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, created)
}

// GetEnrollment handles GET /api/enrollments/{id}
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Enrollment ID is required")
		return
	}

	enrollment, err := h.repo.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, enrollment)
}

// ListEnrollments handles GET /api/enrollments
func (h *EnrollmentHandler) ListEnrollments(w http.ResponseWriter, r *http.Request) {
	enrollments, err := h.repo.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, enrollments)
}

// UpdateEnrollment handles PUT /api/enrollments/{id}
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Enrollment ID is required")
		return
	}

	var enrollment models.Enrollment

	// Decode request body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&enrollment); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update enrollment
	updated, err := h.repo.Update(id, &enrollment)
	if err != nil {
		if err.Error() == "enrollment not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, updated)
}

// DeleteEnrollment handles DELETE /api/enrollments/{id}
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Enrollment ID is required")
		return
	}

	err := h.repo.Delete(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, SuccessResponse{Message: "Enrollment deleted successfully"})
}
