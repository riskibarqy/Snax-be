package service

import (
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type authService struct {
	client clerk.Client
}

// NewAuthService creates a new auth service
func NewAuthService(secretKey string) (domain.AuthService, error) {
	client, err := clerk.NewClient(secretKey)
	if err != nil {
		return nil, err
	}

	return &authService{
		client: client,
	}, nil
}

func (s *authService) ValidateToken(token string) (*domain.Claims, error) {
	// Verify token using Clerk
	clerkClaims, err := s.client.VerifyToken(token)
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
