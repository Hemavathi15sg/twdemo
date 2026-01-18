package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"grademanagement-demo/cache"
	"grademanagement-demo/models"
	"grademanagement-demo/repository"
	"grademanagement-demo/utils"

	"github.com/gorilla/mux"
)

// SuccessResponse for consistent success responses
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// EnrollmentHandler handles HTTP requests for enrollments
type EnrollmentHandler struct {
	repo  *repository.EnrollmentRepository
	cache *cache.EnrollmentCache
}

// NewEnrollmentHandler creates a new handler instance
func NewEnrollmentHandler(repo *repository.EnrollmentRepository, cache *cache.EnrollmentCache) *EnrollmentHandler {
	return &EnrollmentHandler{
		repo:  repo,
		cache: cache,
	}
}

// CreateEnrollment handles POST /api/enrollments
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var input models.EnrollmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	enrollment, err := h.repo.Create(input)
	if err != nil {
		utils.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Cache the newly created enrollment
	if err := h.cache.Set(enrollment); err != nil {
		log.Printf("Failed to cache enrollment: %v", err)
	}

	utils.RespondJSON(w, enrollment, http.StatusCreated)
}

// GetEnrollment handles GET /api/enrollments/{id}
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondError(w, "invalid enrollment ID", http.StatusBadRequest)
		return
	}

	// Try cache first
	enrollment, err := h.cache.GetByID(id)
	if err != nil {
		log.Printf("Cache error: %v", err)
	}

	if enrollment != nil {
		log.Printf("🎯 Cache HIT for enrollment ID %d", id)
		w.Header().Set("X-Cache-Status", "HIT")
		utils.RespondJSON(w, enrollment, http.StatusOK)
		return
	}

	log.Printf("❌ Cache MISS for enrollment ID %d", id)

	// Cache miss - fetch from repository
	enrollment, err = h.repo.GetByID(id)
	if err != nil {
		utils.RespondError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Cache the result for future requests
	if err := h.cache.Set(enrollment); err != nil {
		log.Printf("Failed to cache enrollment: %v", err)
	}

	w.Header().Set("X-Cache-Status", "MISS")
	utils.RespondJSON(w, enrollment, http.StatusOK)
}

// ListEnrollments handles GET /api/enrollments
func (h *EnrollmentHandler) ListEnrollments(w http.ResponseWriter, r *http.Request) {
	enrollments := h.repo.GetAll()
	w.Header().Set("X-Cache-Status", "SKIP")
	utils.RespondJSON(w, enrollments, http.StatusOK)
}

// UpdateEnrollment handles PUT /api/enrollments/{id}
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondError(w, "invalid enrollment ID", http.StatusBadRequest)
		return
	}

	var input models.EnrollmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	enrollment, err := h.repo.Update(id, input)
	if err != nil {
		if err.Error() == "enrollment not found" {
			utils.RespondError(w, err.Error(), http.StatusNotFound)
		} else {
			utils.RespondError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Invalidate cache after update
	if err := h.cache.Delete(id); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	utils.RespondJSON(w, enrollment, http.StatusOK)
}

// DeleteEnrollment handles DELETE /api/enrollments/{id}
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondError(w, "invalid enrollment ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		utils.RespondError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Invalidate cache after delete
	if err := h.cache.Delete(id); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	utils.RespondJSON(w, SuccessResponse{Message: "enrollment deleted successfully"}, http.StatusOK)
}
