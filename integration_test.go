//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"grademanagement-demo/cache"
	"grademanagement-demo/handlers"
	"grademanagement-demo/models"
	"grademanagement-demo/repository"

	"github.com/alicebob/miniredis/v2"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

// setupTestServer creates a test server with mock Redis for both enrollment and grade endpoints
func setupTestServer(t *testing.T) (*httptest.Server, *miniredis.Miniredis, func()) {
	// Create mock Redis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}

	// Create Redis client pointing to mock server
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Initialize enrollment components
	enrollmentRepo := repository.NewEnrollmentRepository()
	enrollmentCache := cache.NewEnrollmentCache(redisClient)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo, enrollmentCache)

	// Initialize grade components
	gradeRepo := repository.NewGradeRepository()
	gradeCache := cache.NewGradeCache(redisClient)
	gradeHandler := handlers.NewGradeHandler(gradeRepo, gradeCache)

	// Setup router
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	
	// Enrollment routes
	api.HandleFunc("/enrollments", enrollmentHandler.CreateEnrollment).Methods("POST")
	api.HandleFunc("/enrollments", enrollmentHandler.ListEnrollments).Methods("GET")
	api.HandleFunc("/enrollments/{id}", enrollmentHandler.GetEnrollment).Methods("GET")
	api.HandleFunc("/enrollments/{id}", enrollmentHandler.UpdateEnrollment).Methods("PUT")
	api.HandleFunc("/enrollments/{id}", enrollmentHandler.DeleteEnrollment).Methods("DELETE")
	
	// Grade routes
	api.HandleFunc("/grades", gradeHandler.CreateGrade).Methods("POST")
	api.HandleFunc("/grades", gradeHandler.ListGrades).Methods("GET")
	api.HandleFunc("/grades/{id}", gradeHandler.GetGrade).Methods("GET")
	api.HandleFunc("/grades/{id}", gradeHandler.UpdateGrade).Methods("PUT")
	api.HandleFunc("/grades/{id}", gradeHandler.DeleteGrade).Methods("DELETE")

	// Create test server
	ts := httptest.NewServer(r)

	cleanup := func() {
		ts.Close()
		mr.Close()
		redisClient.Close()
	}

	return ts, mr, cleanup
}

// TestCompleteCRUDWorkflow tests the entire CRUD lifecycle
func TestCompleteCRUDWorkflow(t *testing.T) {
	ts, _, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Create enrollment
	enrollment := models.EnrollmentInput{
		StudentID: 123,
		CourseID:  456,
		Status:    "active",
	}
	body, _ := json.Marshal(enrollment)
	resp, err := http.Post(ts.URL+"/api/enrollments", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Create request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var created models.Enrollment
	json.NewDecoder(resp.Body).Decode(&created)
	if created.ID == 0 {
		t.Error("Created enrollment has no ID")
	}
	if created.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", created.Status)
	}

	// 2. Get enrollment (should be cached)
	resp, err = http.Get(fmt.Sprintf("%s/api/enrollments/%d", ts.URL, created.ID))
	if err != nil {
		t.Fatalf("Get request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	cacheStatus := resp.Header.Get("X-Cache-Status")
	if cacheStatus != "HIT" && cacheStatus != "MISS" {
		t.Errorf("Expected X-Cache-Status header, got '%s'", cacheStatus)
	}

	// 3. Update enrollment
	update := models.EnrollmentInput{
		Status: "completed",
	}
	body, _ = json.Marshal(update)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/enrollments/%d", ts.URL, created.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Update request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var updated models.Enrollment
	json.NewDecoder(resp.Body).Decode(&updated)
	if updated.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", updated.Status)
	}

	// 4. Delete enrollment
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/api/enrollments/%d", ts.URL, created.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Delete request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// 5. Verify deletion (should return 404)
	resp, _ = http.Get(fmt.Sprintf("%s/api/enrollments/%d", ts.URL, created.ID))
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 after deletion, got %d", resp.StatusCode)
	}
}

// TestCachePerformance validates cached responses are fast
func TestCachePerformance(t *testing.T) {
	ts, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Create enrollment
	enrollment := models.EnrollmentInput{
		StudentID: 999,
		CourseID:  888,
		Status:    "pending",
	}
	body, _ := json.Marshal(enrollment)
	resp, _ := http.Post(ts.URL+"/api/enrollments", "application/json", bytes.NewBuffer(body))
	var created models.Enrollment
	json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()

	// First request (cache miss)
	start := time.Now()
	resp, _ = http.Get(fmt.Sprintf("%s/api/enrollments/%d", ts.URL, created.ID))
	firstDuration := time.Since(start)
	resp.Body.Close()

	// Second request (cache hit)
	start = time.Now()
	resp, _ = http.Get(fmt.Sprintf("%s/api/enrollments/%d", ts.URL, created.ID))
	cachedDuration := time.Since(start)
	cacheStatus := resp.Header.Get("X-Cache-Status")
	resp.Body.Close()

	if cacheStatus != "HIT" {
		t.Errorf("Expected cache HIT, got '%s'", cacheStatus)
	}

	// Cached response should be significantly faster
	if cachedDuration > 100*time.Millisecond {
		t.Errorf("Cached response too slow: %v (expected < 100ms)", cachedDuration)
	}

	t.Logf("Performance: First request: %v, Cached request: %v ⚡", firstDuration, cachedDuration)
}

// TestCacheInvalidation verifies cache is invalidated on update/delete
func TestCacheInvalidation(t *testing.T) {
	ts, mr, cleanup := setupTestServer(t)
	defer cleanup()

	// Create enrollment
	enrollment := models.EnrollmentInput{
		StudentID: 555,
		CourseID:  666,
		Status:    "active",
	}
	body, _ := json.Marshal(enrollment)
	resp, _ := http.Post(ts.URL+"/api/enrollments", "application/json", bytes.NewBuffer(body))
	var created models.Enrollment
	json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()

	// Cache the enrollment
	resp, _ = http.Get(fmt.Sprintf("%s/api/enrollments/%d", ts.URL, created.ID))
	resp.Body.Close()

	// Verify it's in cache
	cacheKey := fmt.Sprintf("enrollment:%d", created.ID)
	if !mr.Exists(cacheKey) {
		t.Error("Enrollment not cached after GET")
	}

	// Update enrollment (should invalidate cache)
	update := models.EnrollmentInput{
		Status: "completed",
	}
	body, _ = json.Marshal(update)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/enrollments/%d", ts.URL, created.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)

	// Cache should be invalidated
	time.Sleep(50 * time.Millisecond) // Small delay for cache invalidation
	if mr.Exists(cacheKey) {
		t.Error("Cache not invalidated after UPDATE")
	}

	// Get again to re-cache
	resp, _ = http.Get(fmt.Sprintf("%s/api/enrollments/%d", ts.URL, created.ID))
	resp.Body.Close()

	// Delete enrollment (should invalidate cache)
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/api/enrollments/%d", ts.URL, created.ID), nil)
	http.DefaultClient.Do(req)

	// Cache should be invalidated
	time.Sleep(50 * time.Millisecond)
	if mr.Exists(cacheKey) {
		t.Error("Cache not invalidated after DELETE")
	}
}

// TestValidationErrors tests all validation error scenarios
func TestValidationErrors(t *testing.T) {
	ts, _, cleanup := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name       string
		input      models.EnrollmentInput
		wantStatus int
		wantError  string
	}{
		{
			name:       "Invalid status",
			input:      models.EnrollmentInput{StudentID: 1, CourseID: 1, Status: "invalid"},
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid status",
		},
		{
			name:       "Missing student_id",
			input:      models.EnrollmentInput{CourseID: 1, Status: "active"},
			wantStatus: http.StatusBadRequest,
			wantError:  "student_id must be a positive integer",
		},
		{
			name:       "Missing course_id",
			input:      models.EnrollmentInput{StudentID: 1, Status: "active"},
			wantStatus: http.StatusBadRequest,
			wantError:  "course_id must be a positive integer",
		},
		{
			name:       "Invalid date format",
			input:      models.EnrollmentInput{StudentID: 1, CourseID: 1, Status: "active", EnrollmentDate: "invalid-date"},
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid enrollment_date format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			resp, _ := http.Post(ts.URL+"/api/enrollments", "application/json", bytes.NewBuffer(body))
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			var errResp map[string]string
			json.NewDecoder(resp.Body).Decode(&errResp)
			if errResp["error"] == "" {
				t.Error("Expected error message in response")
			}
		})
	}
}

// TestNotFoundErrors tests 404 scenarios
func TestNotFoundErrors(t *testing.T) {
	ts, _, cleanup := setupTestServer(t)
	defer cleanup()

	// GET non-existent enrollment
	resp, _ := http.Get(ts.URL + "/api/enrollments/999999")
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 for GET, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// UPDATE non-existent enrollment
	update := models.EnrollmentInput{Status: "active"}
	body, _ := json.Marshal(update)
	req, _ := http.NewRequest("PUT", ts.URL+"/api/enrollments/999999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 for PUT, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// DELETE non-existent enrollment
	req, _ = http.NewRequest("DELETE", ts.URL+"/api/enrollments/999999", nil)
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 for DELETE, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

// TestResponseSchemaValidation validates response structures
func TestResponseSchemaValidation(t *testing.T) {
	ts, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Create and validate response schema
	enrollment := models.EnrollmentInput{
		StudentID: 111,
		CourseID:  222,
		Status:    "pending",
	}
	body, _ := json.Marshal(enrollment)
	resp, _ := http.Post(ts.URL+"/api/enrollments", "application/json", bytes.NewBuffer(body))
	defer resp.Body.Close()

	var created models.Enrollment
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Validate all required fields exist
	if created.ID == 0 {
		t.Error("Missing ID field")
	}
	if created.StudentID == 0 {
		t.Error("Missing student_id field")
	}
	if created.CourseID == 0 {
		t.Error("Missing course_id field")
	}
	if created.Status == "" {
		t.Error("Missing status field")
	}
	if created.EnrollmentDate.IsZero() {
		t.Error("Missing enrollment_date field")
	}
	if created.CreatedAt.IsZero() {
		t.Error("Missing created_at field")
	}
	if created.UpdatedAt.IsZero() {
		t.Error("Missing updated_at field")
	}
}

// === GRADE TESTS ===

// TestGradeCompleteCRUDWorkflow tests the entire CRUD lifecycle for grades
func TestGradeCompleteCRUDWorkflow(t *testing.T) {
	ts, _, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Create grade
	grade := models.GradeInput{
		StudentID: 123,
		CourseID:  456,
		Grade:     "A",
		Status:    "submitted",
	}
	body, _ := json.Marshal(grade)
	resp, err := http.Post(ts.URL+"/api/grades", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Create request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var created models.Grade
	json.NewDecoder(resp.Body).Decode(&created)
	if created.ID == 0 {
		t.Error("Created grade has no ID")
	}
	if created.Grade != "A" {
		t.Errorf("Expected grade 'A', got '%s'", created.Grade)
	}
	if created.Status != "submitted" {
		t.Errorf("Expected status 'submitted', got '%s'", created.Status)
	}

	// 2. Get grade (should be cached)
	resp, err = http.Get(fmt.Sprintf("%s/api/grades/%d", ts.URL, created.ID))
	if err != nil {
		t.Fatalf("Get request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	cacheStatus := resp.Header.Get("X-Cache-Status")
	if cacheStatus != "HIT" && cacheStatus != "MISS" {
		t.Errorf("Expected X-Cache-Status header, got '%s'", cacheStatus)
	}

	// 3. Update grade
	update := models.GradeInput{
		Grade:  "A+",
		Status: "final",
	}
	body, _ = json.Marshal(update)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/grades/%d", ts.URL, created.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Update request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var updated models.Grade
	json.NewDecoder(resp.Body).Decode(&updated)
	if updated.Grade != "A+" {
		t.Errorf("Expected grade 'A+', got '%s'", updated.Grade)
	}
	if updated.Status != "final" {
		t.Errorf("Expected status 'final', got '%s'", updated.Status)
	}

	// 4. Delete grade
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/api/grades/%d", ts.URL, created.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Delete request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// 5. Verify deletion (should return 404)
	resp, _ = http.Get(fmt.Sprintf("%s/api/grades/%d", ts.URL, created.ID))
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 after deletion, got %d", resp.StatusCode)
	}
}

// TestGradeCachePerformance validates cached responses are fast
func TestGradeCachePerformance(t *testing.T) {
	ts, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Create grade
	grade := models.GradeInput{
		StudentID: 999,
		CourseID:  888,
		Grade:     "B+",
		Status:    "draft",
	}
	body, _ := json.Marshal(grade)
	resp, _ := http.Post(ts.URL+"/api/grades", "application/json", bytes.NewBuffer(body))
	var created models.Grade
	json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()

	// First request (cache miss)
	start := time.Now()
	resp, _ = http.Get(fmt.Sprintf("%s/api/grades/%d", ts.URL, created.ID))
	firstDuration := time.Since(start)
	resp.Body.Close()

	// Second request (cache hit)
	start = time.Now()
	resp, _ = http.Get(fmt.Sprintf("%s/api/grades/%d", ts.URL, created.ID))
	cachedDuration := time.Since(start)
	cacheStatus := resp.Header.Get("X-Cache-Status")
	resp.Body.Close()

	if cacheStatus != "HIT" {
		t.Errorf("Expected cache HIT, got '%s'", cacheStatus)
	}

	// Cached response should be significantly faster
	if cachedDuration > 100*time.Millisecond {
		t.Errorf("Cached response too slow: %v (expected < 100ms)", cachedDuration)
	}

	t.Logf("Performance: First request: %v, Cached request: %v ⚡", firstDuration, cachedDuration)
}

// TestGradeCacheInvalidation verifies cache is invalidated on update/delete
func TestGradeCacheInvalidation(t *testing.T) {
	ts, mr, cleanup := setupTestServer(t)
	defer cleanup()

	// Create grade
	grade := models.GradeInput{
		StudentID: 555,
		CourseID:  666,
		Grade:     "C",
		Status:    "submitted",
	}
	body, _ := json.Marshal(grade)
	resp, _ := http.Post(ts.URL+"/api/grades", "application/json", bytes.NewBuffer(body))
	var created models.Grade
	json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()

	// Cache the grade
	resp, _ = http.Get(fmt.Sprintf("%s/api/grades/%d", ts.URL, created.ID))
	resp.Body.Close()

	// Verify it's in cache
	cacheKey := fmt.Sprintf("grade:%d", created.ID)
	if !mr.Exists(cacheKey) {
		t.Error("Grade not cached after GET")
	}

	// Update grade (should invalidate cache)
	update := models.GradeInput{
		Grade:  "B",
		Status: "final",
	}
	body, _ = json.Marshal(update)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/grades/%d", ts.URL, created.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)

	// Cache should be invalidated
	time.Sleep(50 * time.Millisecond) // Small delay for cache invalidation
	if mr.Exists(cacheKey) {
		t.Error("Cache not invalidated after UPDATE")
	}

	// Get again to re-cache
	resp, _ = http.Get(fmt.Sprintf("%s/api/grades/%d", ts.URL, created.ID))
	resp.Body.Close()

	// Delete grade (should invalidate cache)
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/api/grades/%d", ts.URL, created.ID), nil)
	http.DefaultClient.Do(req)

	// Cache should be invalidated
	time.Sleep(50 * time.Millisecond)
	if mr.Exists(cacheKey) {
		t.Error("Cache not invalidated after DELETE")
	}
}

// TestGradeValidationErrors tests all validation error scenarios
func TestGradeValidationErrors(t *testing.T) {
	ts, _, cleanup := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name       string
		input      models.GradeInput
		wantStatus int
		wantError  string
	}{
		{
			name:       "Invalid grade value",
			input:      models.GradeInput{StudentID: 1, CourseID: 1, Grade: "Z", Status: "draft"},
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid grade",
		},
		{
			name:       "Invalid status",
			input:      models.GradeInput{StudentID: 1, CourseID: 1, Grade: "A", Status: "invalid"},
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid status",
		},
		{
			name:       "Missing student_id",
			input:      models.GradeInput{CourseID: 1, Grade: "A", Status: "draft"},
			wantStatus: http.StatusBadRequest,
			wantError:  "student_id must be a positive integer",
		},
		{
			name:       "Missing course_id",
			input:      models.GradeInput{StudentID: 1, Grade: "A", Status: "draft"},
			wantStatus: http.StatusBadRequest,
			wantError:  "course_id must be a positive integer",
		},
		{
			name:       "Invalid date format",
			input:      models.GradeInput{StudentID: 1, CourseID: 1, Grade: "A", Status: "draft", GradeDate: "invalid-date"},
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid grade_date format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			resp, _ := http.Post(ts.URL+"/api/grades", "application/json", bytes.NewBuffer(body))
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			var errResp map[string]string
			json.NewDecoder(resp.Body).Decode(&errResp)
			if errResp["error"] == "" {
				t.Error("Expected error message in response")
			}
		})
	}
}

// TestGradeNotFoundErrors tests 404 scenarios
func TestGradeNotFoundErrors(t *testing.T) {
	ts, _, cleanup := setupTestServer(t)
	defer cleanup()

	// GET non-existent grade
	resp, _ := http.Get(ts.URL + "/api/grades/999999")
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 for GET, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// UPDATE non-existent grade
	update := models.GradeInput{Grade: "A", Status: "final"}
	body, _ := json.Marshal(update)
	req, _ := http.NewRequest("PUT", ts.URL+"/api/grades/999999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 for PUT, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// DELETE non-existent grade
	req, _ = http.NewRequest("DELETE", ts.URL+"/api/grades/999999", nil)
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 for DELETE, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

// TestGradeResponseSchemaValidation validates response structures
func TestGradeResponseSchemaValidation(t *testing.T) {
	ts, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Create and validate response schema
	grade := models.GradeInput{
		StudentID: 111,
		CourseID:  222,
		Grade:     "B-",
		Status:    "draft",
	}
	body, _ := json.Marshal(grade)
	resp, _ := http.Post(ts.URL+"/api/grades", "application/json", bytes.NewBuffer(body))
	defer resp.Body.Close()

	var created models.Grade
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Validate all required fields exist
	if created.ID == 0 {
		t.Error("Missing ID field")
	}
	if created.StudentID == 0 {
		t.Error("Missing student_id field")
	}
	if created.CourseID == 0 {
		t.Error("Missing course_id field")
	}
	if created.Grade == "" {
		t.Error("Missing grade field")
	}
	if created.Status == "" {
		t.Error("Missing status field")
	}
	if created.GradeDate.IsZero() {
		t.Error("Missing grade_date field")
	}
	if created.CreatedAt.IsZero() {
		t.Error("Missing created_at field")
	}
	if created.UpdatedAt.IsZero() {
		t.Error("Missing updated_at field")
	}
}
