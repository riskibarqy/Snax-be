package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	redisClient *redis.Client
	requests    int
	duration    time.Duration
}

func NewRateLimiter(redisURL string, redisToken string, requests int, duration time.Duration) *RateLimiter {
	// opts := &redis.Options{
	// 	Addr:     redisURL,
	// 	Password: redisToken, // Redis password/token
	// 	DB:       0,          // Default DB
	// }

	// client := redis.NewClient(opts)

	opt, _ := redis.ParseURL(fmt.Sprintf("rediss://default:%s@%s", redisToken, redisURL))

	client := redis.NewClient(opt)

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	return &RateLimiter{
		redisClient: client,
		requests:    requests,
		duration:    duration,
	}
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get client identifier (IP or user ID)
		identifier := r.RemoteAddr
		if userID := r.Header.Get("X-User-ID"); userID != "" {
			identifier = userID
		}

		// Create a unique key for this rate limit window
		key := fmt.Sprintf("ratelimit:%s", identifier)

		// Use Redis to track request count
		count, err := rl.redisClient.Incr(ctx, key).Result()
		if err != nil {
			http.Error(w, "Rate limiting error", http.StatusInternalServerError)
			return
		}

		// Set expiry on first request
		if count == 1 {
			rl.redisClient.Expire(ctx, key, rl.duration)
		}

		// Check if rate limit exceeded
		if count > int64(rl.requests) {
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.requests))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(rl.duration).Unix()))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.requests))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.requests-int(count)))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(rl.duration).Unix()))

		next.ServeHTTP(w, r)
	})
}
