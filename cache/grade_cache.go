package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"grademanagement-demo/models"

	"github.com/redis/go-redis/v9"
)

// GradeCache handles Redis caching for grades
type GradeCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewGradeCache creates a new grade cache instance
func NewGradeCache(client *redis.Client) *GradeCache {
	return &GradeCache{
		client: client,
		ttl:    5 * time.Minute, // 5-minute TTL per project standards
	}
}

// Set caches a grade
func (c *GradeCache) Set(grade *models.Grade) error {
	ctx := context.Background()
	key := fmt.Sprintf("grade:%d", grade.ID)

	data, err := json.Marshal(grade)
	if err != nil {
		return fmt.Errorf("failed to marshal grade: %w", err)
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}

// GetByID retrieves a grade from cache
func (c *GradeCache) GetByID(id int) (*models.Grade, error) {
	ctx := context.Background()
	key := fmt.Sprintf("grade:%d", id)

	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}

	var grade models.Grade
	if err := json.Unmarshal([]byte(data), &grade); err != nil {
		return nil, fmt.Errorf("failed to unmarshal grade: %w", err)
	}

	return &grade, nil
}

// Delete removes a grade from cache
func (c *GradeCache) Delete(id int) error {
	ctx := context.Background()
	key := fmt.Sprintf("grade:%d", id)
	return c.client.Del(ctx, key).Err()
}

// GetByStudentAndCourse retrieves a grade by student and course from cache
func (c *GradeCache) GetByStudentAndCourse(studentID, courseID int) (*models.Grade, error) {
	ctx := context.Background()
	key := fmt.Sprintf("grade:student:%d:course:%d", studentID, courseID)

	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}

	var grade models.Grade
	if err := json.Unmarshal([]byte(data), &grade); err != nil {
		return nil, fmt.Errorf("failed to unmarshal grade: %w", err)
	}

	return &grade, nil
}

// SetByStudentAndCourse caches a grade with student-course key
func (c *GradeCache) SetByStudentAndCourse(grade *models.Grade) error {
	ctx := context.Background()
	key := fmt.Sprintf("grade:student:%d:course:%d", grade.StudentID, grade.CourseID)

	data, err := json.Marshal(grade)
	if err != nil {
		return fmt.Errorf("failed to marshal grade: %w", err)
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}
