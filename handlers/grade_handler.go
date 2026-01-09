package handlers

import (
	"encoding/json"
	"grademanagement-demo/models"
	"grademanagement-demo/repos"
	"net/http"

	"github.com/gorilla/mux"
)

// GradeHandler handles HTTP requests for grade operations
type GradeHandler struct {
	repo repos.GradeRepository
}

// NewGradeHandler creates a new grade handler
func NewGradeHandler(repo repos.GradeRepository) *GradeHandler {
	return &GradeHandler{
		repo: repo,
	}
}

// CreateGrade handles POST requests to create a new grade
func (h *GradeHandler) CreateGrade(w http.ResponseWriter, r *http.Request) {
	var req models.CreateGradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.StudentID == "" || req.CourseID == "" || req.Grade == "" {
		http.Error(w, "Missing required fields: student_id, course_id, or grade", http.StatusBadRequest)
		return
	}

	grade := &models.Grade{
		StudentID:    req.StudentID,
		CourseID:     req.CourseID,
		Grade:        req.Grade,
		Score:        req.Score,
		Semester:     req.Semester,
		AcademicYear: req.AcademicYear,
	}

	if err := h.repo.Create(grade); err != nil {
		http.Error(w, "Failed to create grade", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(grade)
}

// GetGrade handles GET requests to retrieve a specific grade by ID
func (h *GradeHandler) GetGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	grade, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Grade not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grade)
}

// GetAllGrades handles GET requests to retrieve all grades
func (h *GradeHandler) GetAllGrades(w http.ResponseWriter, r *http.Request) {
	grades, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Failed to retrieve grades", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

// GetGradesByStudent handles GET requests to retrieve grades for a specific student
func (h *GradeHandler) GetGradesByStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]

	grades, err := h.repo.GetByStudentID(studentID)
	if err != nil {
		http.Error(w, "Failed to retrieve grades", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

// GetGradesByCourse handles GET requests to retrieve grades for a specific course
func (h *GradeHandler) GetGradesByCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["course_id"]

	grades, err := h.repo.GetByCourseID(courseID)
	if err != nil {
		http.Error(w, "Failed to retrieve grades", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

// UpdateGrade handles PUT requests to update an existing grade
func (h *GradeHandler) UpdateGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.UpdateGradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(id, &req); err != nil {
		http.Error(w, "Grade not found", http.StatusNotFound)
		return
	}

	// Retrieve the updated grade to return
	grade, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to retrieve updated grade", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grade)
}

// DeleteGrade handles DELETE requests to remove a grade
func (h *GradeHandler) DeleteGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "Grade not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Grade deleted successfully"})
}
