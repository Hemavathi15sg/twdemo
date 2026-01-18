package main

import (
	"fmt"
	"log"
	"net/http"

	"grademanagement-demo/repos"
	"grademanagement-demo/routes"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	
	// Initialize repositories
	gradeRepo := repos.NewInMemoryGradeRepository()
	enrollmentRepo := repos.NewInMemoryEnrollmentRepository(redisClient)
	
	// Setup routes
	routes.SetupGradeRoutes(r, gradeRepo)
	routes.SetupEnrollmentRoutes(r, enrollmentRepo)
	
	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API - Ready for AI delegation!", "status": "healthy"}`)
	}).Methods("GET")

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Ready for Copilot Agent delegation!")
	fmt.Println("📚 Grade Management endpoints available at /api/grades")
	fmt.Println("🎓 Enrollment endpoints available at /api/enrollments (with Redis caching)")
	
	log.Fatal(http.ListenAndServe(port, r))
}