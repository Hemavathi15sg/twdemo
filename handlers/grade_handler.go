package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"grademanagement-demo/cache"
	"grademanagement-demo/models"
	"grademanagement-demo/repository"

	"github.com/gorilla/mux"
)

// GradeHandler handles HTTP requests for grades
type GradeHandler struct {
	repo  *repository.GradeRepository
	cache *cache.GradeCache
}

// NewGradeHandler creates a new handler instance
func NewGradeHandler(repo *repository.GradeRepository, cache *cache.GradeCache) *GradeHandler {
	return &GradeHandler{
		repo:  repo,
		cache: cache,
	}
}

// CreateGrade handles POST /api/grades
func (h *GradeHandler) CreateGrade(w http.ResponseWriter, r *http.Request) {
	var input models.GradeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	grade, err := h.repo.Create(input)
	if err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Cache the newly created grade
	if err := h.cache.Set(grade); err != nil {
		log.Printf("Failed to cache grade: %v", err)
	}

	respondJSON(w, grade, http.StatusCreated)
}

// GetGrade handles GET /api/grades/{id}
func (h *GradeHandler) GetGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "invalid grade ID", http.StatusBadRequest)
		return
	}

	// Try cache first
	grade, err := h.cache.GetByID(id)
	if err != nil {
		log.Printf("Cache error: %v", err)
	}

	if grade != nil {
		log.Printf("🎯 Cache HIT for grade ID %d", id)
		w.Header().Set("X-Cache-Status", "HIT")
		respondJSON(w, grade, http.StatusOK)
		return
	}

	log.Printf("❌ Cache MISS for grade ID %d", id)

	// Cache miss - fetch from repository
	grade, err = h.repo.GetByID(id)
	if err != nil {
		respondError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Cache the result for future requests
	if err := h.cache.Set(grade); err != nil {
		log.Printf("Failed to cache grade: %v", err)
	}

	w.Header().Set("X-Cache-Status", "MISS")
	respondJSON(w, grade, http.StatusOK)
}

// ListGrades handles GET /api/grades
func (h *GradeHandler) ListGrades(w http.ResponseWriter, r *http.Request) {
	grades := h.repo.GetAll()
	w.Header().Set("X-Cache-Status", "SKIP")
	respondJSON(w, grades, http.StatusOK)
}

// UpdateGrade handles PUT /api/grades/{id}
func (h *GradeHandler) UpdateGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "invalid grade ID", http.StatusBadRequest)
		return
	}

	var input models.GradeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	grade, err := h.repo.Update(id, input)
	if err != nil {
		if err.Error() == "grade not found" {
			respondError(w, err.Error(), http.StatusNotFound)
		} else {
			respondError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Invalidate cache after update
	if err := h.cache.Delete(id); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	respondJSON(w, grade, http.StatusOK)
}

// DeleteGrade handles DELETE /api/grades/{id}
func (h *GradeHandler) DeleteGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "invalid grade ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		respondError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Invalidate cache after delete
	if err := h.cache.Delete(id); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	respondJSON(w, SuccessResponse{Message: "grade deleted successfully"}, http.StatusOK)
}
