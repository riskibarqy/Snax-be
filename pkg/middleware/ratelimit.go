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
	client := redis.NewClient(redisURL, redisToken)
	return &RateLimiter{
		redisClient: client,
		requests:    requests,
		duration:    duration,
	}
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client identifier (IP or user ID)
		identifier := r.RemoteAddr
		if userID := r.Header.Get("X-User-ID"); userID != "" {
			identifier = userID
		}

		// Create a unique key for this rate limit window
		key := fmt.Sprintf("ratelimit:%s", identifier)

		// Use Redis to track request count
		count, err := rl.redisClient.Incr(context.Background(), key).Int()
		if err != nil {
			http.Error(w, "Rate limiting error", http.StatusInternalServerError)
			return
		}

		// Set expiry on first request
		if count == 1 {
			rl.redisClient.Expire(context.Background(), key, int(rl.duration.Seconds()))
		}

		// Check if rate limit exceeded
		if count > rl.requests {
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.requests))
			w.Header().Set("X-RateLimit-Remaining", "0")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.requests))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.requests-count))

		next.ServeHTTP(w, r)
	})
}
