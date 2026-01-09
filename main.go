package main

import (
	"fmt"
	"log"
	"net/http"

	"grademanagement-demo/repos"
	"grademanagement-demo/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	
	// Initialize repositories
	gradeRepo := repos.NewInMemoryGradeRepository()
	
	// Setup routes
	routes.SetupGradeRoutes(r, gradeRepo)
	
	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API - Ready for AI delegation!", "status": "healthy"}`)
	}).Methods("GET")

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Ready for Copilot Agent delegation!")
	fmt.Println("📚 Grade Management endpoints available at /api/grades")
	
	log.Fatal(http.ListenAndServe(port, r))
}