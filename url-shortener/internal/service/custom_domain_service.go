package service

import (
	"strings"
	"time"

	internalDomain "github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type customDomainService struct {
	repo internalDomain.CustomDomainRepository
}

// NewCustomDomainService creates a new custom domain service
func NewCustomDomainService(repo internalDomain.CustomDomainRepository) internalDomain.CustomDomainService {
	return &customDomainService{
		repo: repo,
	}
}

func (s *customDomainService) RegisterDomain(domain, userID string) (*internalDomain.CustomDomain, error) {
	// Normalize domain
	domain = strings.ToLower(strings.TrimSpace(domain))

	// Check if domain already exists
	if _, err := s.repo.GetByDomain(domain); err == nil {
		return nil, &internalDomain.ErrDomainNotFound{Domain: domain}
	}

	// Create new domain
	d := &internalDomain.CustomDomain{
		Domain:    domain,
		UserID:    userID,
		Verified:  false,
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *customDomainService) VerifyDomain(id int64) error {
	return s.repo.UpdateVerificationStatus(id, true)
}

func (s *customDomainService) GetUserDomains(userID string) ([]internalDomain.CustomDomain, error) {
	return s.repo.GetByUserID(userID)
}

func (s *customDomainService) DeleteDomain(id int64, userID string) error {
	return s.repo.Delete(id, userID)
}
