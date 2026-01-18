package repositories

import (
	"context"
	"grademanagement-demo/models"
	"testing"
)

func TestInMemoryEnrollmentRepository(t *testing.T) {
	repo := NewInMemoryEnrollmentRepository()
	ctx := context.Background()

	t.Run("Create enrollment", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:   "student-123",
			StudentName: "John Doe",
			CourseID:    "course-456",
			CourseName:  "Introduction to Go",
			Status:      models.StatusActive,
			Credits:     3,
		}

		created, err := repo.Create(ctx, enrollment)
		if err != nil {
			t.Fatalf("Failed to create enrollment: %v", err)
		}

		if created.ID == "" {
			t.Error("Expected ID to be generated")
		}
		if created.StudentID != "student-123" {
			t.Errorf("Expected StudentID to be 'student-123', got '%s'", created.StudentID)
		}
		if created.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
		if created.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})

	t.Run("GetByID enrollment", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:   "student-456",
			StudentName: "Jane Smith",
			CourseID:    "course-789",
			CourseName:  "Advanced Go",
			Status:      models.StatusActive,
			Credits:     4,
		}

		created, _ := repo.Create(ctx, enrollment)
		retrieved, err := repo.GetByID(ctx, created.ID)

		if err != nil {
			t.Fatalf("Failed to get enrollment: %v", err)
		}
		if retrieved.ID != created.ID {
			t.Errorf("Expected ID to be '%s', got '%s'", created.ID, retrieved.ID)
		}
	})

	t.Run("GetByID non-existent enrollment", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "non-existent-id")
		if err == nil {
			t.Error("Expected error when getting non-existent enrollment")
		}
	})

	t.Run("Update enrollment", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:   "student-789",
			StudentName: "Bob Johnson",
			CourseID:    "course-101",
			CourseName:  "Go Testing",
			Status:      models.StatusActive,
			Credits:     2,
		}

		created, _ := repo.Create(ctx, enrollment)
		created.Status = models.StatusCompleted
		grade := "B+"
		created.Grade = &grade

		updated, err := repo.Update(ctx, created)
		if err != nil {
			t.Fatalf("Failed to update enrollment: %v", err)
		}
		if updated.Status != models.StatusCompleted {
			t.Errorf("Expected Status to be '%s', got '%s'", models.StatusCompleted, updated.Status)
		}
		if *updated.Grade != "B+" {
			t.Errorf("Expected Grade to be 'B+', got '%s'", *updated.Grade)
		}
	})

	t.Run("Delete enrollment", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:   "student-999",
			StudentName: "Alice Williams",
			CourseID:    "course-202",
			CourseName:  "Go Patterns",
			Status:      models.StatusActive,
			Credits:     3,
		}

		created, _ := repo.Create(ctx, enrollment)
		err := repo.Delete(ctx, created.ID)

		if err != nil {
			t.Fatalf("Failed to delete enrollment: %v", err)
		}

		_, err = repo.GetByID(ctx, created.ID)
		if err == nil {
			t.Error("Expected error when getting deleted enrollment")
		}
	})

	t.Run("List all enrollments", func(t *testing.T) {
		repo := NewInMemoryEnrollmentRepository()
		
		enrollment1 := &models.Enrollment{
			StudentID:   "student-111",
			StudentName: "Student One",
			CourseID:    "course-111",
			CourseName:  "Course One",
			Status:      models.StatusActive,
			Credits:     3,
		}
		enrollment2 := &models.Enrollment{
			StudentID:   "student-222",
			StudentName: "Student Two",
			CourseID:    "course-222",
			CourseName:  "Course Two",
			Status:      models.StatusActive,
			Credits:     4,
		}

		repo.Create(ctx, enrollment1)
		repo.Create(ctx, enrollment2)

		enrollments, err := repo.List(ctx, "", "")
		if err != nil {
			t.Fatalf("Failed to list enrollments: %v", err)
		}
		if len(enrollments) != 2 {
			t.Errorf("Expected 2 enrollments, got %d", len(enrollments))
		}
	})

	t.Run("List enrollments by student", func(t *testing.T) {
		repo := NewInMemoryEnrollmentRepository()
		
		enrollment1 := &models.Enrollment{
			StudentID:   "student-333",
			StudentName: "Student Three",
			CourseID:    "course-111",
			CourseName:  "Course One",
			Status:      models.StatusActive,
			Credits:     3,
		}
		enrollment2 := &models.Enrollment{
			StudentID:   "student-333",
			StudentName: "Student Three",
			CourseID:    "course-222",
			CourseName:  "Course Two",
			Status:      models.StatusActive,
			Credits:     4,
		}
		enrollment3 := &models.Enrollment{
			StudentID:   "student-444",
			StudentName: "Student Four",
			CourseID:    "course-333",
			CourseName:  "Course Three",
			Status:      models.StatusActive,
			Credits:     2,
		}

		repo.Create(ctx, enrollment1)
		repo.Create(ctx, enrollment2)
		repo.Create(ctx, enrollment3)

		enrollments, err := repo.GetByStudentID(ctx, "student-333")
		if err != nil {
			t.Fatalf("Failed to get enrollments by student: %v", err)
		}
		if len(enrollments) != 2 {
			t.Errorf("Expected 2 enrollments for student-333, got %d", len(enrollments))
		}
	})

	t.Run("List enrollments by course", func(t *testing.T) {
		repo := NewInMemoryEnrollmentRepository()
		
		enrollment1 := &models.Enrollment{
			StudentID:   "student-555",
			StudentName: "Student Five",
			CourseID:    "course-444",
			CourseName:  "Course Four",
			Status:      models.StatusActive,
			Credits:     3,
		}
		enrollment2 := &models.Enrollment{
			StudentID:   "student-666",
			StudentName: "Student Six",
			CourseID:    "course-444",
			CourseName:  "Course Four",
			Status:      models.StatusActive,
			Credits:     3,
		}

		repo.Create(ctx, enrollment1)
		repo.Create(ctx, enrollment2)

		enrollments, err := repo.GetByCourseID(ctx, "course-444")
		if err != nil {
			t.Fatalf("Failed to get enrollments by course: %v", err)
		}
		if len(enrollments) != 2 {
			t.Errorf("Expected 2 enrollments for course-444, got %d", len(enrollments))
		}
	})
}
