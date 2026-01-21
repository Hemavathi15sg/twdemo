package main

import (
	"fmt"
	"log"
	"net/http"

	"grademanagement-demo/handlers"
	"grademanagement-demo/repository"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Initialize repository
	enrollmentRepo := repository.NewEnrollmentRepository()

	// Initialize handlers
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo)

	// Setup routes
	r.HandleFunc("/api/enrollments", enrollmentHandler.CreateEnrollment).Methods("POST")
	r.HandleFunc("/api/enrollments", enrollmentHandler.GetAllEnrollments).Methods("GET")
	r.HandleFunc("/api/enrollments/{id}", enrollmentHandler.GetEnrollment).Methods("GET")
	r.HandleFunc("/api/enrollments/{id}", enrollmentHandler.UpdateEnrollment).Methods("PUT")
	r.HandleFunc("/api/enrollments/{id}", enrollmentHandler.DeleteEnrollment).Methods("DELETE")
	r.HandleFunc("/api/enrollments/student/{student_id}", enrollmentHandler.GetEnrollmentsByStudent).Methods("GET")
	r.HandleFunc("/api/enrollments/course/{course_id}", enrollmentHandler.GetEnrollmentsByCourse).Methods("GET")

	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API - Ready for AI delegation!", "status": "healthy"}`)
	}).Methods("GET")

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Ready for Copilot Agent delegation!")
	fmt.Println("🎓 Enrollment endpoints available at /api/enrollments")

	log.Fatal(http.ListenAndServe(port, r))
}
