package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(redisURL, redisToken string) error {
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	// Strip any protocol prefix if present
	redisURL = strings.TrimPrefix(redisURL, "redis://")
	redisURL = strings.TrimPrefix(redisURL, "rediss://")

	// Add default port if not present
	if !strings.Contains(redisURL, ":") {
		redisURL = redisURL + ":6379"
	}

	// Build the Redis options
	opt := &redis.Options{
		Addr: redisURL,
	}

	// If we have a token, set it as password and enable TLS
	if redisToken != "" {
		opt.Password = redisToken
		opt.TLSConfig = &tls.Config{} // Enable TLS for secure connections
	}

	RedisClient = redis.NewClient(opt)

	// Test the connection
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}
