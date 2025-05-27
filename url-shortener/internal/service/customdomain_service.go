package service

import (
	"context"
	"time"

	internalDomain "github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type CustomDomainService struct {
	repo internalDomain.CustomDomainRepository
}

// New creates a new custom domain service
func NewCustomDomainService(repo internalDomain.CustomDomainRepository) internalDomain.CustomDomainService {
	return &CustomDomainService{
		repo: repo,
	}
}

// RegisterDomain registers a new custom domain
func (s *CustomDomainService) RegisterDomain(ctx context.Context, domain string, userID string) (*internalDomain.CustomDomain, error) {
	// Check if domain already exists
	existingDomain, err := s.repo.GetByDomain(ctx, domain)
	if err == nil && existingDomain != nil {
		return nil, &internalDomain.ErrDomainAlreadyExists{Domain: domain}
	}

	customDomain := &internalDomain.CustomDomain{
		Domain:    domain,
		UserID:    userID,
		Verified:  false,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, customDomain); err != nil {
		return nil, err
	}

	return customDomain, nil
}

// GetUserDomains retrieves all domains for a user
func (s *CustomDomainService) GetUserDomains(ctx context.Context, userID string) ([]internalDomain.CustomDomain, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// DeleteDomain deletes a custom domain
func (s *CustomDomainService) DeleteDomain(ctx context.Context, id int64, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}

// VerifyDomain marks a domain as verified
func (s *CustomDomainService) VerifyDomain(ctx context.Context, id int64) error {
	return s.repo.VerifyDomain(ctx, id)
}
