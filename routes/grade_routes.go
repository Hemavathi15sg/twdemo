package routes

import (
	"grademanagement-demo/handlers"
	"grademanagement-demo/repos"

	"github.com/gorilla/mux"
)

// SetupGradeRoutes configures all grade-related routes
func SetupGradeRoutes(router *mux.Router, gradeRepo repos.GradeRepository) {
	gradeHandler := handlers.NewGradeHandler(gradeRepo)

	// Grade CRUD endpoints
	router.HandleFunc("/api/grades", gradeHandler.CreateGrade).Methods("POST")
	router.HandleFunc("/api/grades", gradeHandler.GetAllGrades).Methods("GET")
	router.HandleFunc("/api/grades/{id}", gradeHandler.GetGrade).Methods("GET")
	router.HandleFunc("/api/grades/{id}", gradeHandler.UpdateGrade).Methods("PUT")
	router.HandleFunc("/api/grades/{id}", gradeHandler.DeleteGrade).Methods("DELETE")

	// Additional query endpoints
	router.HandleFunc("/api/students/{student_id}/grades", gradeHandler.GetGradesByStudent).Methods("GET")
	router.HandleFunc("/api/courses/{course_id}/grades", gradeHandler.GetGradesByCourse).Methods("GET")
}
