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
	gradeCachePrefix = "grade:"
	gradeCacheTTL    = 5 * time.Minute
)

// GradeCache handles caching for grade data
type GradeCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewGradeCache creates a new cache instance
func NewGradeCache(client *redis.Client) *GradeCache {
	return &GradeCache{
		client: client,
		ctx:    context.Background(),
	}
}

// GetByID retrieves a grade from cache by ID
func (c *GradeCache) GetByID(id int) (*models.Grade, error) {
	key := fmt.Sprintf("%s%d", gradeCachePrefix, id)
	
	data, err := c.client.Get(c.ctx, key).Bytes()
	if err == redis.Nil {
		// Cache miss
		return nil, nil
	}
	if err != nil {
		log.Printf("Redis get error: %v", err)
		return nil, err
	}

	var grade models.Grade
	if err := json.Unmarshal(data, &grade); err != nil {
		log.Printf("Cache unmarshal error: %v", err)
		return nil, err
	}

	return &grade, nil
}

// Set stores a grade in cache with TTL
func (c *GradeCache) Set(grade *models.Grade) error {
	key := fmt.Sprintf("%s%d", gradeCachePrefix, grade.ID)
	
	data, err := json.Marshal(grade)
	if err != nil {
		log.Printf("Cache marshal error: %v", err)
		return err
	}

	if err := c.client.Set(c.ctx, key, data, gradeCacheTTL).Err(); err != nil {
		log.Printf("Redis set error: %v", err)
		return err
	}

	log.Printf("✅ Cached grade ID %d (TTL: %v)", grade.ID, gradeCacheTTL)
	return nil
}

// Delete removes a grade from cache
func (c *GradeCache) Delete(id int) error {
	key := fmt.Sprintf("%s%d", gradeCachePrefix, id)
	
	if err := c.client.Del(c.ctx, key).Err(); err != nil {
		log.Printf("Redis delete error: %v", err)
		return err
	}

	log.Printf("🗑️  Invalidated cache for grade ID %d", id)
	return nil
}
