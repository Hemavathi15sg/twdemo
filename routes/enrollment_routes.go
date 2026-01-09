package routes

import (
	"grademanagement-demo/handlers"
	"grademanagement-demo/repository"

	"github.com/gorilla/mux"
)

// SetupEnrollmentRoutes sets up all enrollment routes with /api prefix
func SetupEnrollmentRoutes(router *mux.Router) {
	// Initialize repository and handler
	enrollmentRepo := repository.NewInMemoryEnrollmentRepository()
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo)

	// Create subrouter with /api prefix
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Register enrollment routes
	apiRouter.HandleFunc("/enrollments", enrollmentHandler.CreateEnrollment).Methods("POST")
	apiRouter.HandleFunc("/enrollments", enrollmentHandler.ListEnrollments).Methods("GET")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.GetEnrollment).Methods("GET")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.UpdateEnrollment).Methods("PUT")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.DeleteEnrollment).Methods("DELETE")
}
