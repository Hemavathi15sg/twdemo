package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"grademanagement-demo/cache"
	"grademanagement-demo/models"
	"grademanagement-demo/repository"

	"github.com/gorilla/mux"
)

// GradeHandler handles HTTP requests for grade calculations
type GradeHandler struct {
	repo  repository.GradeRepositoryInterface
	cache *cache.GradeCache
}

// NewGradeHandler creates a new handler instance
func NewGradeHandler(repo repository.GradeRepositoryInterface, cache *cache.GradeCache) *GradeHandler {
	return &GradeHandler{
		repo:  repo,
		cache: cache,
	}
}

// CalculateGrade handles POST /api/grades/calculate
// Performance requirement: <200ms for 100 students
// Validates against Figma design tokens
func (h *GradeHandler) CalculateGrade(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	var input models.GradeCalculationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := models.ValidateGradeInput(input); err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate weighted average
	weightedAverage, err := models.CalculateWeightedAverage(input.Assignments)
	if err != nil {
		respondError(w, "failed to calculate weighted average", http.StatusInternalServerError)
		return
	}

	// Apply curve if requested
	numericGrade := weightedAverage
	if input.ApplyCurve {
		numericGrade = models.ApplyGradeCurve(weightedAverage, input.CurveAmount)
	}

	// Convert to letter grade using Figma design tokens
	letterGrade, gradeColor, gradeStatus := models.ConvertToLetterGrade(numericGrade)

	// Create grade entity
	grade := &models.Grade{
		StudentID:       input.StudentID,
		CourseID:        input.CourseID,
		NumericGrade:    numericGrade,
		LetterGrade:     letterGrade,
		GradeColor:      gradeColor,
		GradeStatus:     gradeStatus,
		WeightedAverage: weightedAverage,
		CurveApplied:    input.ApplyCurve,
		CurveAmount:     input.CurveAmount,
	}

	// Save to repository
	savedGrade, err := h.repo.Create(grade)
	if err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cache the grade
	if err := h.cache.Set(savedGrade); err != nil {
		log.Printf("Failed to cache grade: %v", err)
	}
	if err := h.cache.SetByStudentAndCourse(savedGrade); err != nil {
		log.Printf("Failed to cache grade by student/course: %v", err)
	}

	// Calculate response time
	responseTime := time.Since(startTime)
	w.Header().Set("X-Response-Time", responseTime.String())

	// Log performance
	if responseTime.Milliseconds() > 200 {
		log.Printf("⚠️  Performance warning: Grade calculation took %v (target: <200ms)", responseTime)
	} else {
		log.Printf("✅ Performance OK: Grade calculation took %v", responseTime)
	}

	respondJSON(w, savedGrade, http.StatusCreated)
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

	// Cache miss - query repository
	log.Printf("🔍 Cache MISS for grade ID %d", id)
	w.Header().Set("X-Cache-Status", "MISS")

	grade, err = h.repo.GetByID(id)
	if err != nil {
		respondError(w, "grade not found", http.StatusNotFound)
		return
	}

	// Update cache
	if err := h.cache.Set(grade); err != nil {
		log.Printf("Failed to cache grade: %v", err)
	}

	respondJSON(w, grade, http.StatusOK)
}

// GetGradeByStudentAndCourse handles GET /api/grades/student/{studentId}/course/{courseId}
func (h *GradeHandler) GetGradeByStudentAndCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID, err := strconv.Atoi(vars["studentId"])
	if err != nil {
		respondError(w, "invalid student ID", http.StatusBadRequest)
		return
	}

	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		respondError(w, "invalid course ID", http.StatusBadRequest)
		return
	}

	// Try cache first
	grade, err := h.cache.GetByStudentAndCourse(studentID, courseID)
	if err != nil {
		log.Printf("Cache error: %v", err)
	}

	if grade != nil {
		log.Printf("🎯 Cache HIT for student %d, course %d", studentID, courseID)
		w.Header().Set("X-Cache-Status", "HIT")
		respondJSON(w, grade, http.StatusOK)
		return
	}

	// Cache miss - query repository
	log.Printf("🔍 Cache MISS for student %d, course %d", studentID, courseID)
	w.Header().Set("X-Cache-Status", "MISS")

	grade, err = h.repo.GetByStudentAndCourse(studentID, courseID)
	if err != nil {
		respondError(w, "grade not found", http.StatusNotFound)
		return
	}

	// Update cache
	if err := h.cache.Set(grade); err != nil {
		log.Printf("Failed to cache grade: %v", err)
	}
	if err := h.cache.SetByStudentAndCourse(grade); err != nil {
		log.Printf("Failed to cache grade by student/course: %v", err)
	}

	respondJSON(w, grade, http.StatusOK)
}

// GetAllGrades handles GET /api/grades
func (h *GradeHandler) GetAllGrades(w http.ResponseWriter, r *http.Request) {
	grades := h.repo.GetAll()
	respondJSON(w, grades, http.StatusOK)
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
