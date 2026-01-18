package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"grademanagement-demo/models"

	"github.com/redis/go-redis/v9"
)

const (
	enrollmentCachePrefix = "enrollment:"
	enrollmentCacheTTL    = 5 * time.Minute
)

// EnrollmentCache handles caching for enrollment data
type EnrollmentCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewEnrollmentCache creates a new cache instance
func NewEnrollmentCache(client *redis.Client) *EnrollmentCache {
	return &EnrollmentCache{
		client: client,
		ctx:    context.Background(),
	}
}

// GetByID retrieves an enrollment from cache by ID
func (c *EnrollmentCache) GetByID(id int) (*models.Enrollment, error) {
	key := fmt.Sprintf("%s%d", enrollmentCachePrefix, id)

	data, err := c.client.Get(c.ctx, key).Bytes()
	if err == redis.Nil {
		// Cache miss
		return nil, nil
	}
	if err != nil {
		log.Printf("Redis get error: %v", err)
		return nil, err
	}

	var enrollment models.Enrollment
	if err := json.Unmarshal(data, &enrollment); err != nil {
		log.Printf("Cache unmarshal error: %v", err)
		return nil, err
	}

	return &enrollment, nil
}

// Set stores an enrollment in cache with TTL
func (c *EnrollmentCache) Set(enrollment *models.Enrollment) error {
	key := fmt.Sprintf("%s%d", enrollmentCachePrefix, enrollment.ID)

	data, err := json.Marshal(enrollment)
	if err != nil {
		log.Printf("Cache marshal error: %v", err)
		return err
	}

	if err := c.client.Set(c.ctx, key, data, enrollmentCacheTTL).Err(); err != nil {
		log.Printf("Redis set error: %v", err)
		return err
	}

	log.Printf("✅ Cached enrollment ID %d (TTL: %v)", enrollment.ID, enrollmentCacheTTL)
	return nil
}

// Delete removes an enrollment from cache
func (c *EnrollmentCache) Delete(id int) error {
	key := fmt.Sprintf("%s%d", enrollmentCachePrefix, id)

	if err := c.client.Del(c.ctx, key).Err(); err != nil {
		log.Printf("Redis delete error: %v", err)
		return err
	}

	log.Printf("🗑️  Invalidated cache for enrollment ID %d", id)
	return nil
}
