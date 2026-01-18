package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"grademanagement-demo/models"
	"grademanagement-demo/repositories"
	"grademanagement-demo/usecases"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type EnrollmentHTTPHandler struct {
	svc *usecases.EnrollmentService
}

func NewEnrollmentHTTPHandler(svc *usecases.EnrollmentService) *EnrollmentHTTPHandler {
	return &EnrollmentHTTPHandler{svc: svc}
}

func (h *EnrollmentHTTPHandler) RegisterRoutes(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/enrollments", h.create).Methods(http.MethodPost)
	api.HandleFunc("/enrollments", h.list).Methods(http.MethodGet)
	api.HandleFunc("/enrollments/{id}", h.get).Methods(http.MethodGet)
	api.HandleFunc("/enrollments/{id}", h.update).Methods(http.MethodPut)
	api.HandleFunc("/enrollments/{id}", h.delete).Methods(http.MethodDelete)
}

type enrollmentRequest struct {
	StudentID      string  `json:"student_id"`
	CourseID       string  `json:"course_id"`
	EnrollmentDate *string `json:"enrollment_date,omitempty"`
	Status         string  `json:"status"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type enrollmentResponse struct {
	ID             string    `json:"id"`
	StudentID      string    `json:"student_id"`
	CourseID       string    `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (h *EnrollmentHTTPHandler) create(w http.ResponseWriter, r *http.Request) {
	var req enrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid JSON request"})
		return
	}
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid student_id"})
		return
	}
	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid course_id"})
		return
	}
	var enrollDate time.Time
	if req.EnrollmentDate != nil && *req.EnrollmentDate != "" {
		t, perr := time.Parse(time.RFC3339, *req.EnrollmentDate)
		if perr != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid enrollment_date (use RFC3339)"})
			return
		}
		enrollDate = t
	}
	e := &models.Enrollment{
		StudentID:      studentID,
		CourseID:       courseID,
		EnrollmentDate: enrollDate,
		Status:         req.Status,
	}
	res, err := h.svc.Create(r.Context(), e)
	if err != nil {
		if err == usecases.ErrValidation {
			writeJSON(w, http.StatusBadRequest, errorResponse{Message: "validation failed"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Message: "internal server error"})
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(res))
}

func (h *EnrollmentHTTPHandler) get(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid id"})
		return
	}
	res, err := h.svc.GetByID(context.Background(), id)
	if err != nil {
		if err == repositories.ErrNotFound {
			writeJSON(w, http.StatusNotFound, errorResponse{Message: "enrollment not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Message: "internal server error"})
		return
	}
	writeJSON(w, http.StatusOK, toResponse(res))
}

func (h *EnrollmentHTTPHandler) list(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Message: "internal server error"})
		return
	}
	out := make([]enrollmentResponse, 0, len(res))
	for _, e := range res {
		out = append(out, toResponse(e))
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *EnrollmentHTTPHandler) update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid id"})
		return
	}
	var req enrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid JSON request"})
		return
	}
	studentID, sErr := uuid.Parse(req.StudentID)
	if sErr != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid student_id"})
		return
	}
	courseID, cErr := uuid.Parse(req.CourseID)
	if cErr != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid course_id"})
		return
	}
	var enrollDate time.Time
	if req.EnrollmentDate != nil && *req.EnrollmentDate != "" {
		t, perr := time.Parse(time.RFC3339, *req.EnrollmentDate)
		if perr != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid enrollment_date (use RFC3339)"})
			return
		}
		enrollDate = t
	}
	e := &models.Enrollment{
		ID:             id,
		StudentID:      studentID,
		CourseID:       courseID,
		EnrollmentDate: enrollDate,
		Status:         req.Status,
	}
	res, err := h.svc.Update(r.Context(), e)
	if err != nil {
		if err == usecases.ErrValidation {
			writeJSON(w, http.StatusBadRequest, errorResponse{Message: "validation failed"})
			return
		}
		if err == repositories.ErrNotFound {
			writeJSON(w, http.StatusNotFound, errorResponse{Message: "enrollment not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Message: "internal server error"})
		return
	}
	writeJSON(w, http.StatusOK, toResponse(res))
}

func (h *EnrollmentHTTPHandler) delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid id"})
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		if err == repositories.ErrNotFound {
			writeJSON(w, http.StatusNotFound, errorResponse{Message: "enrollment not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Message: "internal server error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toResponse(e *models.Enrollment) enrollmentResponse {
	return enrollmentResponse{
		ID:             e.ID.String(),
		StudentID:      e.StudentID.String(),
		CourseID:       e.CourseID.String(),
		EnrollmentDate: e.EnrollmentDate,
		Status:         e.Status,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
