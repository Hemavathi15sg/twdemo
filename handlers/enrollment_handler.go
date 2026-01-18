package handlers

import (
	"encoding/json"
	"grademanagement-demo/models"
	"grademanagement-demo/usecases"
	"net/http"

	"github.com/gorilla/mux"
)

// EnrollmentHandler handles HTTP requests for enrollment operations
type EnrollmentHandler struct {
	useCase *usecases.EnrollmentUseCase
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(useCase *usecases.EnrollmentUseCase) *EnrollmentHandler {
	return &EnrollmentHandler{
		useCase: useCase,
	}
}

// CreateEnrollment handles POST /enrollments
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	
	var enrollment models.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	created, err := h.useCase.CreateEnrollment(r.Context(), &enrollment)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, created)
}

// GetEnrollment handles GET /enrollments/{id}
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	enrollment, err := h.useCase.GetEnrollment(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, enrollment)
}

// UpdateEnrollment handles PUT /enrollments/{id}
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	
	vars := mux.Vars(r)
	id := vars["id"]

	var enrollment models.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	enrollment.ID = id
	updated, err := h.useCase.UpdateEnrollment(r.Context(), &enrollment)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, updated)
}

// DeleteEnrollment handles DELETE /enrollments/{id}
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.useCase.DeleteEnrollment(r.Context(), id); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "enrollment deleted successfully"})
}

// ListEnrollments handles GET /enrollments
func (h *EnrollmentHandler) ListEnrollments(w http.ResponseWriter, r *http.Request) {
	studentID := r.URL.Query().Get("student_id")
	courseID := r.URL.Query().Get("course_id")

	enrollments, err := h.useCase.ListEnrollments(r.Context(), studentID, courseID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if enrollments == nil {
		enrollments = []*models.Enrollment{}
	}

	respondWithJSON(w, http.StatusOK, enrollments)
}

// GetEnrollmentsByStudent handles GET /students/{student_id}/enrollments
func (h *EnrollmentHandler) GetEnrollmentsByStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]

	enrollments, err := h.useCase.GetEnrollmentsByStudent(r.Context(), studentID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if enrollments == nil {
		enrollments = []*models.Enrollment{}
	}

	respondWithJSON(w, http.StatusOK, enrollments)
}

// GetEnrollmentsByCourse handles GET /courses/{course_id}/enrollments
func (h *EnrollmentHandler) GetEnrollmentsByCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["course_id"]

	enrollments, err := h.useCase.GetEnrollmentsByCourse(r.Context(), courseID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if enrollments == nil {
		enrollments = []*models.Enrollment{}
	}

	respondWithJSON(w, http.StatusOK, enrollments)
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
