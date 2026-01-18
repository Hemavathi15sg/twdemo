package main

import (
	"fmt"
	"grademanagement-demo/handlers"
	"grademanagement-demo/repositories"
	"grademanagement-demo/usecases"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Initialize enrollment feature
	enrollmentRepo := repositories.NewInMemoryEnrollmentRepository()
	enrollmentUseCase := usecases.NewEnrollmentUseCase(enrollmentRepo)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentUseCase)

	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API - Ready for AI delegation!", "status": "healthy"}`)
	}).Methods("GET")

	// Enrollment routes
	r.HandleFunc("/enrollments", enrollmentHandler.CreateEnrollment).Methods("POST")
	r.HandleFunc("/enrollments", enrollmentHandler.ListEnrollments).Methods("GET")
	r.HandleFunc("/enrollments/{id}", enrollmentHandler.GetEnrollment).Methods("GET")
	r.HandleFunc("/enrollments/{id}", enrollmentHandler.UpdateEnrollment).Methods("PUT")
	r.HandleFunc("/enrollments/{id}", enrollmentHandler.DeleteEnrollment).Methods("DELETE")
	r.HandleFunc("/students/{student_id}/enrollments", enrollmentHandler.GetEnrollmentsByStudent).Methods("GET")
	r.HandleFunc("/courses/{course_id}/enrollments", enrollmentHandler.GetEnrollmentsByCourse).Methods("GET")

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Ready for Copilot Agent delegation!")
	fmt.Println("📚 Enrollment API endpoints:")
	fmt.Println("   POST   /enrollments")
	fmt.Println("   GET    /enrollments")
	fmt.Println("   GET    /enrollments/{id}")
	fmt.Println("   PUT    /enrollments/{id}")
	fmt.Println("   DELETE /enrollments/{id}")
	fmt.Println("   GET    /students/{student_id}/enrollments")
	fmt.Println("   GET    /courses/{course_id}/enrollments")

	log.Fatal(http.ListenAndServe(port, r))
}
