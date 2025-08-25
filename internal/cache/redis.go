package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// l      = logger.New().Prefix("Cache Logger")
	client *redis.Client
	ctx    = context.Background()
	ttl    = 5 * time.Minute
)

// InitializeRedis initializes Redis connection using env vars.
// Required: REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, REDIS_DB
func InitializeRedis() error {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	pass := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	if host == "" || port == "" || dbStr == "" {
		return fmt.Errorf("missing required Redis env vars (REDIS_HOST, REDIS_PORT, REDIS_DB)")
	}

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
		DB:       db,
	})

	// Test connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}
