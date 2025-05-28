package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/config"
)

type RateLimiter struct {
	// Requests per minute limits
	AuthUserLimit   int
	GuestUserLimit  int
	ExpirationInSec int
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		AuthUserLimit:   10, // 10 requests per minute for authenticated users
		GuestUserLimit:  5,  // 5 requests per minute for guests
		ExpirationInSec: 60, // Reset counter every minute (60 seconds)
	}
}

// getIP extracts the IP address from various headers and falls back to RemoteAddr
func getIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		// Get the original client IP (first one)
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header next
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	// Fall back to RemoteAddr, but clean it
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If we can't split, just return RemoteAddr as is
		return r.RemoteAddr
	}
	return ip
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientID string

		// Get client identifier (IP for guests, user ID for authenticated users)
		if userID := r.Header.Get("X-User-ID"); userID != "" {
			clientID = "user:" + userID // Use user ID for authenticated users
		} else {
			clientID = "ip:" + getIP(r) // Use IP address for guests
		}

		// Create Redis key with timestamp to ensure per-minute window
		timestamp := time.Now().Unix() / 60 // Get current minute
		key := fmt.Sprintf("ratelimit:%s:%d", clientID, timestamp)

		// Determine rate limit based on authentication
		limit := rl.GuestUserLimit
		if r.Header.Get("X-User-ID") != "" {
			limit = rl.AuthUserLimit
		}

		// Use Redis MULTI to ensure atomicity
		ctx := context.Background()
		pipe := config.RedisClient.TxPipeline()

		// Increment first (returns new value)
		incr := pipe.Incr(ctx, key)
		// Set expiration (only if key is new)
		pipe.Expire(ctx, key, time.Duration(rl.ExpirationInSec)*time.Second)

		// Execute transaction
		_, err := pipe.Exec(ctx)
		if err != nil {
			http.Error(w, "Rate limiting error", http.StatusInternalServerError)
			return
		}

		// Get the new count after increment
		count := incr.Val()

		// Check if rate limit exceeded
		if count > int64(limit) {
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt((timestamp+1)*60, 10))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Set rate limit headers
		remaining := limit - int(count)
		if remaining < 0 {
			remaining = 0
		}
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt((timestamp+1)*60, 10))

		next.ServeHTTP(w, r)
	})
}
