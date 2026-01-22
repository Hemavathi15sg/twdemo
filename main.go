package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"grademanagement-demo/cache"
	"grademanagement-demo/handlers"
	"grademanagement-demo/repository"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize Redis client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	fmt.Println("✅ Connected to Redis successfully")

	r := mux.NewRouter()

	// Initialize repositories, caches, and handlers
	enrollmentRepo := repository.NewEnrollmentRepository()
	enrollmentCache := cache.NewEnrollmentCache(redisClient)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo, enrollmentCache)

	gradeRepo := repository.NewGradeRepository()
	gradeCache := cache.NewGradeCache(redisClient)
	gradeHandler := handlers.NewGradeHandler(gradeRepo, gradeCache)

	// API routes with /api prefix
	api := r.PathPrefix("/api").Subrouter()
	
	// Enrollment endpoints
	api.HandleFunc("/enrollments", enrollmentHandler.CreateEnrollment).Methods("POST")
	api.HandleFunc("/enrollments", enrollmentHandler.ListEnrollments).Methods("GET")
	api.HandleFunc("/enrollments/{id}", enrollmentHandler.GetEnrollment).Methods("GET")
	api.HandleFunc("/enrollments/{id}", enrollmentHandler.UpdateEnrollment).Methods("PUT")
	api.HandleFunc("/enrollments/{id}", enrollmentHandler.DeleteEnrollment).Methods("DELETE")

	// Grade calculation endpoints - TEC-31
	api.HandleFunc("/grades/calculate", gradeHandler.CalculateGrade).Methods("POST")
	api.HandleFunc("/grades", gradeHandler.GetAllGrades).Methods("GET")
	api.HandleFunc("/grades/{id}", gradeHandler.GetGrade).Methods("GET")
	api.HandleFunc("/grades/student/{studentId}/course/{courseId}", gradeHandler.GetGradeByStudentAndCourse).Methods("GET")

	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API with Redis Caching", "status": "healthy"}`)
	}).Methods("GET")

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Enrollment API available at /api/enrollments")
	fmt.Println("🎓 Grade Calculation API available at /api/grades/calculate")
	fmt.Println("   - POST /api/grades/calculate (TEC-31)")
	fmt.Println("   - GET  /api/grades")
	fmt.Println("   - GET  /api/grades/{id}")
	fmt.Println("   - GET  /api/grades/student/{studentId}/course/{courseId}")
	fmt.Println("⚡ Redis caching enabled with 5-minute TTL")
	fmt.Println("🎨 Figma design tokens enforced for grade display")
	fmt.Println("⏱️  Performance target: <200ms per request")

	log.Fatal(http.ListenAndServe(port, r))
}
