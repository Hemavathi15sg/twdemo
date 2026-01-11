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

	// Initialize repository, cache, and handler
	repo := repository.NewEnrollmentRepository()
	enrollmentCache := cache.NewEnrollmentCache(redisClient)
	handler := handlers.NewEnrollmentHandler(repo, enrollmentCache)

	// API routes with /api prefix
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/enrollments", handler.CreateEnrollment).Methods("POST")
	api.HandleFunc("/enrollments", handler.ListEnrollments).Methods("GET")
	api.HandleFunc("/enrollments/{id}", handler.GetEnrollment).Methods("GET")
	api.HandleFunc("/enrollments/{id}", handler.UpdateEnrollment).Methods("PUT")
	api.HandleFunc("/enrollments/{id}", handler.DeleteEnrollment).Methods("DELETE")

	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API with Redis Caching", "status": "healthy"}`)
	}).Methods("GET")

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Enrollment API available at /api/enrollments")
	fmt.Println("⚡ Redis caching enabled with 5-minute TTL")

	log.Fatal(http.ListenAndServe(port, r))
}
