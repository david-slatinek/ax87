package cache

import (
	"api/model"
	"context"
	"github.com/go-redis/redis/v9"
	"os"
	"time"
)

type Cache struct {
	// Redis connection string.
	url string
	// Redis password.
	password string
	// Redis cache client.
	client *redis.Client
}

// Load Cache fields - url, password - from the environment.
func (cache *Cache) Load() {
	cache.url = os.Getenv("REDIS_URL")
	cache.password = os.Getenv("REDIS_PASSWORD")
}

// Create redis.Client and check for status.
func (cache *Cache) Create() error {
	cache.client = redis.NewClient(&redis.Options{
		Addr:            cache.url,
		MaxRetries:      5,
		MinRetryBackoff: time.Millisecond * 100,
		DialTimeout:     time.Millisecond * 100,
	})

	_, err := cache.client.Ping(context.Background()).Result()
	return err
}

// Close the client.
func (cache *Cache) Close() error {
	return cache.client.Close()
}

// Add model.DataResponse to the cache.
func (cache *Cache) Add(dr *model.DataResponse) error {
	return cache.client.Set(context.Background(), dr.DataType, dr, time.Minute*5).Err()
}

// Get value from the cache for the given key.
func (cache *Cache) Get(key string) (string, error) {
	return cache.client.Get(context.Background(), key).Result()
}
