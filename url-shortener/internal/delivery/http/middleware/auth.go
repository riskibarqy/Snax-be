package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
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

		// Validate token
		claims, err := am.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), SessionContextKey, claims)
		r.Header.Set("X-User-ID", claims.Subject)

		next.ServeHTTP(w, r.WithContext(ctx))
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

func (am *AuthMiddleware) ValidateToken(token string) (*domain.Claims, error) {
	// Verify token using Clerk
	clerkClaims, err := am.client.VerifyToken(token)
	if err != nil {
		return nil, &domain.ErrInvalidToken{Message: "Invalid token"}
	}

	// Convert Clerk claims to our domain claims
	claims := &domain.Claims{
		Subject:   clerkClaims.Subject,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Default to 24 hours
		IssuedAt:  time.Now().Unix(),
		TokenID:   clerkClaims.ID,
	}

	return claims, nil
}
