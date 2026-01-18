package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"grademanagement-demo/models"
	"grademanagement-demo/repositories"
	"grademanagement-demo/usecases"
)

// EnrollmentHandler handles HTTP requests for enrollments
type EnrollmentHandler struct {
	service *usecases.EnrollmentService
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(service *usecases.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{service: service}
}

// RegisterRoutes registers enrollment routes
func (h *EnrollmentHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/enrollments", h.CreateEnrollment).Methods("POST")
	router.HandleFunc("/api/enrollments", h.ListEnrollments).Methods("GET")
	router.HandleFunc("/api/enrollments/{id}", h.GetEnrollment).Methods("GET")
	router.HandleFunc("/api/enrollments/{id}", h.UpdateEnrollment).Methods("PUT")
	router.HandleFunc("/api/enrollments/{id}", h.DeleteEnrollment).Methods("DELETE")
}

// CreateEnrollment handles POST /api/enrollments
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var req usecases.CreateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Invalid request body"})
		return
	}

	enrollment, err := h.service.CreateEnrollment(r.Context(), &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toEnrollmentResponse(enrollment))
}

// GetEnrollment handles GET /api/enrollments/{id}
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Invalid enrollment ID"})
		return
	}

	enrollment, err := h.service.GetEnrollmentByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "Enrollment not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toEnrollmentResponse(enrollment))
}

// ListEnrollments handles GET /api/enrollments
func (h *EnrollmentHandler) ListEnrollments(w http.ResponseWriter, r *http.Request) {
	enrollments, err := h.service.ListEnrollments(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  toEnrollmentResponses(enrollments),
		"count": len(enrollments),
	})
}

// UpdateEnrollment handles PUT /api/enrollments/{id}
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Invalid enrollment ID"})
		return
	}

	var req usecases.UpdateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Invalid request body"})
		return
	}

	enrollment, err := h.service.UpdateEnrollment(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "Enrollment not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toEnrollmentResponse(enrollment))
}

// DeleteEnrollment handles DELETE /api/enrollments/{id}
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Invalid enrollment ID"})
		return
	}

	err = h.service.DeleteEnrollment(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "Enrollment not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Response DTOs
type enrollmentResponse struct {
	ID             int64  `json:"id"`
	StudentID      int64  `json:"student_id"`
	CourseID       int64  `json:"course_id"`
	EnrollmentDate string `json:"enrollment_date"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// toEnrollmentResponse converts model to response DTO
func toEnrollmentResponse(e *models.Enrollment) enrollmentResponse {
	return enrollmentResponse{
		ID:             e.ID,
		StudentID:      e.StudentID,
		CourseID:       e.CourseID,
		EnrollmentDate: e.EnrollmentDate.Format("2006-01-02T15:04:05Z"),
		Status:         e.Status,
		CreatedAt:      e.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      e.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// toEnrollmentResponses converts slice of models to response DTOs
func toEnrollmentResponses(enrollments []*models.Enrollment) []enrollmentResponse {
	responses := make([]enrollmentResponse, len(enrollments))
	for i, e := range enrollments {
		responses[i] = toEnrollmentResponse(e)
	}
	return responses
}
