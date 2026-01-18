package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"grademanagement-demo/handlers"
	"grademanagement-demo/repositories"
)

func main() {
	// Initialize repository
	enrollmentRepo := repositories.NewInMemoryEnrollmentRepository()
	
	// Initialize handlers
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo)
	
	r := mux.NewRouter()
	
	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API - Ready for AI delegation!", "status": "healthy"}`)
	}).Methods("GET")

	// API routes with /api prefix
	api := r.PathPrefix("/api").Subrouter()
	
	// Enrollment routes
	api.HandleFunc("/enrollments", enrollmentHandler.Create).Methods("POST")
	api.HandleFunc("/enrollments", enrollmentHandler.List).Methods("GET")
	api.HandleFunc("/enrollments/{id}", enrollmentHandler.GetByID).Methods("GET")
	api.HandleFunc("/enrollments/{id}", enrollmentHandler.Update).Methods("PUT")
	api.HandleFunc("/enrollments/{id}", enrollmentHandler.Delete).Methods("DELETE")

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Ready for Copilot Agent delegation!")
	fmt.Println("📚 Enrollment API endpoints:")
	fmt.Println("   POST   /api/enrollments")
	fmt.Println("   GET    /api/enrollments")
	fmt.Println("   GET    /api/enrollments/{id}")
	fmt.Println("   PUT    /api/enrollments/{id}")
	fmt.Println("   DELETE /api/enrollments/{id}")
	
	log.Fatal(http.ListenAndServe(port, r))
}