package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"grademanagement-demo/models"
	"grademanagement-demo/repositories"
)

// EnrollmentHandler handles HTTP requests for enrollment operations
type EnrollmentHandler struct {
	repo repositories.EnrollmentRepository
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(repo repositories.EnrollmentRepository) *EnrollmentHandler {
	return &EnrollmentHandler{
		repo: repo,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// sendJSON sends a JSON response
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// sendError sends an error response
func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, ErrorResponse{Error: message})
}

// Create handles POST /api/enrollments
func (h *EnrollmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEnrollmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.StudentID == 0 {
		sendError(w, http.StatusBadRequest, "student_id is required")
		return
	}
	if req.CourseID == 0 {
		sendError(w, http.StatusBadRequest, "course_id is required")
		return
	}
	if req.Status == "" {
		sendError(w, http.StatusBadRequest, "status is required")
		return
	}

	// Validate status
	if !models.IsValidStatus(req.Status) {
		sendError(w, http.StatusBadRequest, "status must be pending, active, or completed")
		return
	}

	enrollment := &models.Enrollment{
		StudentID: req.StudentID,
		CourseID:  req.CourseID,
		Status:    models.EnrollmentStatus(req.Status),
	}

	created, err := h.repo.Create(r.Context(), enrollment)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create enrollment")
		return
	}

	sendJSON(w, http.StatusCreated, created)
}

// GetByID handles GET /api/enrollments/{id}
func (h *EnrollmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid enrollment ID")
		return
	}

	enrollment, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if err == repositories.ErrEnrollmentNotFound {
			sendError(w, http.StatusNotFound, "Enrollment not found")
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to retrieve enrollment")
		return
	}

	sendJSON(w, http.StatusOK, enrollment)
}

// List handles GET /api/enrollments
func (h *EnrollmentHandler) List(w http.ResponseWriter, r *http.Request) {
	enrollments, err := h.repo.List(r.Context())
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to retrieve enrollments")
		return
	}

	sendJSON(w, http.StatusOK, enrollments)
}

// Update handles PUT /api/enrollments/{id}
func (h *EnrollmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid enrollment ID")
		return
	}

	var req models.UpdateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate status if provided
	if req.Status != nil && *req.Status != "" {
		if !models.IsValidStatus(*req.Status) {
			sendError(w, http.StatusBadRequest, "status must be pending, active, or completed")
			return
		}
	}

	enrollment := &models.Enrollment{}
	if req.Status != nil {
		enrollment.Status = models.EnrollmentStatus(*req.Status)
	}

	updated, err := h.repo.Update(r.Context(), id, enrollment)
	if err != nil {
		if err == repositories.ErrEnrollmentNotFound {
			sendError(w, http.StatusNotFound, "Enrollment not found")
			return
		}
		if err == repositories.ErrInvalidStatus {
			sendError(w, http.StatusBadRequest, "Invalid status value")
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update enrollment")
		return
	}

	sendJSON(w, http.StatusOK, updated)
}

// Delete handles DELETE /api/enrollments/{id}
func (h *EnrollmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid enrollment ID")
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		if err == repositories.ErrEnrollmentNotFound {
			sendError(w, http.StatusNotFound, "Enrollment not found")
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to delete enrollment")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
