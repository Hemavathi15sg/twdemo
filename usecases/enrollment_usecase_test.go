package usecases

import (
	"context"
	"grademanagement-demo/models"
	"grademanagement-demo/repositories"
	"testing"
)

func TestEnrollmentUseCase(t *testing.T) {
	repo := repositories.NewInMemoryEnrollmentRepository()
	useCase := NewEnrollmentUseCase(repo)
	ctx := context.Background()

	t.Run("CreateEnrollment with valid data", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:   "student-123",
			StudentName: "John Doe",
			CourseID:    "course-456",
			CourseName:  "Introduction to Go",
			Credits:     3,
		}

		created, err := useCase.CreateEnrollment(ctx, enrollment)
		if err != nil {
			t.Fatalf("Failed to create enrollment: %v", err)
		}
		if created.Status != models.StatusActive {
			t.Errorf("Expected default status to be '%s', got '%s'", models.StatusActive, created.Status)
		}
	})

	t.Run("CreateEnrollment with missing student_id", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentName: "Jane Smith",
			CourseID:    "course-789",
			CourseName:  "Advanced Go",
			Credits:     4,
		}

		_, err := useCase.CreateEnrollment(ctx, enrollment)
		if err == nil {
			t.Error("Expected error when student_id is missing")
		}
	})

	t.Run("CreateEnrollment with missing student_name", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:  "student-456",
			CourseID:   "course-789",
			CourseName: "Advanced Go",
			Credits:    4,
		}

		_, err := useCase.CreateEnrollment(ctx, enrollment)
		if err == nil {
			t.Error("Expected error when student_name is missing")
		}
	})

	t.Run("CreateEnrollment with invalid credits", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:   "student-789",
			StudentName: "Bob Johnson",
			CourseID:    "course-101",
			CourseName:  "Go Testing",
			Credits:     15,
		}

		_, err := useCase.CreateEnrollment(ctx, enrollment)
		if err == nil {
			t.Error("Expected error when credits > 10")
		}
	})

	t.Run("CreateEnrollment with invalid status", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:   "student-999",
			StudentName: "Alice Williams",
			CourseID:    "course-202",
			CourseName:  "Go Patterns",
			Status:      "invalid-status",
			Credits:     3,
		}

		_, err := useCase.CreateEnrollment(ctx, enrollment)
		if err == nil {
			t.Error("Expected error when status is invalid")
		}
	})

	t.Run("GetEnrollment by ID", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:   "student-111",
			StudentName: "Student One",
			CourseID:    "course-111",
			CourseName:  "Course One",
			Credits:     3,
		}

		created, _ := useCase.CreateEnrollment(ctx, enrollment)
		retrieved, err := useCase.GetEnrollment(ctx, created.ID)

		if err != nil {
			t.Fatalf("Failed to get enrollment: %v", err)
		}
		if retrieved.ID != created.ID {
			t.Errorf("Expected ID to be '%s', got '%s'", created.ID, retrieved.ID)
		}
	})

	t.Run("GetEnrollment with empty ID", func(t *testing.T) {
		_, err := useCase.GetEnrollment(ctx, "")
		if err == nil {
			t.Error("Expected error when ID is empty")
		}
	})

	t.Run("UpdateEnrollment", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:   "student-222",
			StudentName: "Student Two",
			CourseID:    "course-222",
			CourseName:  "Course Two",
			Credits:     4,
		}

		created, _ := useCase.CreateEnrollment(ctx, enrollment)
		
		update := &models.Enrollment{
			ID:      created.ID,
			Status:  models.StatusCompleted,
			Credits: 4,
		}
		grade := "A"
		update.Grade = &grade

		updated, err := useCase.UpdateEnrollment(ctx, update)
		if err != nil {
			t.Fatalf("Failed to update enrollment: %v", err)
		}
		if updated.Status != models.StatusCompleted {
			t.Errorf("Expected Status to be '%s', got '%s'", models.StatusCompleted, updated.Status)
		}
		if *updated.Grade != "A" {
			t.Errorf("Expected Grade to be 'A', got '%s'", *updated.Grade)
		}
	})

	t.Run("UpdateEnrollment with invalid status", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:   "student-333",
			StudentName: "Student Three",
			CourseID:    "course-333",
			CourseName:  "Course Three",
			Credits:     2,
		}

		created, _ := useCase.CreateEnrollment(ctx, enrollment)
		
		update := &models.Enrollment{
			ID:      created.ID,
			Status:  "invalid-status",
			Credits: 2,
		}

		_, err := useCase.UpdateEnrollment(ctx, update)
		if err == nil {
			t.Error("Expected error when updating with invalid status")
		}
	})

	t.Run("DeleteEnrollment", func(t *testing.T) {
		enrollment := &models.Enrollment{
			StudentID:   "student-444",
			StudentName: "Student Four",
			CourseID:    "course-444",
			CourseName:  "Course Four",
			Credits:     3,
		}

		created, _ := useCase.CreateEnrollment(ctx, enrollment)
		err := useCase.DeleteEnrollment(ctx, created.ID)

		if err != nil {
			t.Fatalf("Failed to delete enrollment: %v", err)
		}

		_, err = useCase.GetEnrollment(ctx, created.ID)
		if err == nil {
			t.Error("Expected error when getting deleted enrollment")
		}
	})

	t.Run("ListEnrollments", func(t *testing.T) {
		repo := repositories.NewInMemoryEnrollmentRepository()
		useCase := NewEnrollmentUseCase(repo)

		enrollment1 := &models.Enrollment{
			StudentID:   "student-555",
			StudentName: "Student Five",
			CourseID:    "course-555",
			CourseName:  "Course Five",
			Credits:     3,
		}
		enrollment2 := &models.Enrollment{
			StudentID:   "student-666",
			StudentName: "Student Six",
			CourseID:    "course-666",
			CourseName:  "Course Six",
			Credits:     4,
		}

		useCase.CreateEnrollment(ctx, enrollment1)
		useCase.CreateEnrollment(ctx, enrollment2)

		enrollments, err := useCase.ListEnrollments(ctx, "", "")
		if err != nil {
			t.Fatalf("Failed to list enrollments: %v", err)
		}
		if len(enrollments) != 2 {
			t.Errorf("Expected 2 enrollments, got %d", len(enrollments))
		}
	})

	t.Run("GetEnrollmentsByStudent", func(t *testing.T) {
		repo := repositories.NewInMemoryEnrollmentRepository()
		useCase := NewEnrollmentUseCase(repo)

		enrollment1 := &models.Enrollment{
			StudentID:   "student-777",
			StudentName: "Student Seven",
			CourseID:    "course-111",
			CourseName:  "Course One",
			Credits:     3,
		}
		enrollment2 := &models.Enrollment{
			StudentID:   "student-777",
			StudentName: "Student Seven",
			CourseID:    "course-222",
			CourseName:  "Course Two",
			Credits:     4,
		}

		useCase.CreateEnrollment(ctx, enrollment1)
		useCase.CreateEnrollment(ctx, enrollment2)

		enrollments, err := useCase.GetEnrollmentsByStudent(ctx, "student-777")
		if err != nil {
			t.Fatalf("Failed to get enrollments by student: %v", err)
		}
		if len(enrollments) != 2 {
			t.Errorf("Expected 2 enrollments for student-777, got %d", len(enrollments))
		}
	})

	t.Run("GetEnrollmentsByCourse", func(t *testing.T) {
		repo := repositories.NewInMemoryEnrollmentRepository()
		useCase := NewEnrollmentUseCase(repo)

		enrollment1 := &models.Enrollment{
			StudentID:   "student-888",
			StudentName: "Student Eight",
			CourseID:    "course-999",
			CourseName:  "Course Nine",
			Credits:     3,
		}
		enrollment2 := &models.Enrollment{
			StudentID:   "student-999",
			StudentName: "Student Nine",
			CourseID:    "course-999",
			CourseName:  "Course Nine",
			Credits:     3,
		}

		useCase.CreateEnrollment(ctx, enrollment1)
		useCase.CreateEnrollment(ctx, enrollment2)

		enrollments, err := useCase.GetEnrollmentsByCourse(ctx, "course-999")
		if err != nil {
			t.Fatalf("Failed to get enrollments by course: %v", err)
		}
		if len(enrollments) != 2 {
			t.Errorf("Expected 2 enrollments for course-999, got %d", len(enrollments))
		}
	})
}
