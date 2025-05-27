package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

// Custom context key for Clerk session
type contextKey string

const SessionContextKey contextKey = "clerk-session"

// SessionClaims represents the claims from a verified token
type SessionClaims struct {
	Subject string
	// Add other claims as needed
}

type AuthMiddleware struct {
	client clerk.Client
}

func NewAuthMiddleware(secretKey string) (*AuthMiddleware, error) {
	client, err := clerk.NewClient(secretKey)
	if err != nil {
		return nil, err
	}

	return &AuthMiddleware{
		client: client,
	}, nil
}

func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		fmt.Println(tokenString, "==")

		// Verify session token
		claims, err := am.client.VerifyToken(tokenString)
		fmt.Println(err, "x=x=x=")
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Create session claims
		sessionClaims := &SessionClaims{
			Subject: claims.Subject,
		}

		// Add session to context
		ctx := context.WithValue(r.Context(), SessionContextKey, sessionClaims)
		r = r.WithContext(ctx)

		// Add user ID to header for rate limiting
		r.Header.Set("X-User-ID", sessionClaims.Subject)

		next.ServeHTTP(w, r)
	})
}

// OptionalAuth middleware for routes that can be accessed without authentication
func (am *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := am.client.VerifyToken(tokenString)
			if err == nil {
				sessionClaims := &SessionClaims{
					Subject: claims.Subject,
				}
				ctx := context.WithValue(r.Context(), SessionContextKey, sessionClaims)
				r = r.WithContext(ctx)
				r.Header.Set("X-User-ID", sessionClaims.Subject)
			}
		}
		next.ServeHTTP(w, r)
	})
}
