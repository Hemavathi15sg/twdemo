package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"grademanagement-demo/handlers"
	"grademanagement-demo/repositories"
	"grademanagement-demo/usecases"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API - Ready for AI delegation!", "status": "healthy"}`)
	}).Methods("GET")

	// Wire enrollment feature under /api prefix
	repo := repositories.NewInMemoryEnrollmentRepository()
	svc := usecases.NewEnrollmentService(repo)
	h := handlers.NewEnrollmentHTTPHandler(svc)
	h.RegisterRoutes(r)

	p := os.Getenv("PORT")
	if p == "" {
		p = "8080"
	}
	port := ":" + p
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Ready for Copilot Agent delegation!")

	log.Fatal(http.ListenAndServe(port, r))
}
