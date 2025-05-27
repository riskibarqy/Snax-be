package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type URLService struct {
	repo domain.URLRepository
}

// New creates a new URL service
func NewURLService(repo domain.URLRepository) domain.URLService {
	return &URLService{
		repo: repo,
	}
}

// generateShortCode generates a random short code for URLs
func generateShortCode() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:8], nil
}

// CreateShortURL creates a new shortened URL
func (s *URLService) CreateShortURL(ctx context.Context, originalURL string, userID string, expiresAt *time.Time) (*domain.URL, error) {
	shortCode, err := generateShortCode()
	if err != nil {
		return nil, err
	}

	url := &domain.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		UserID:      userID,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
		IsActive:    true,
	}

	if err := s.repo.Create(url); err != nil {
		return nil, err
	}

	return url, nil
}

// GetURL retrieves a URL by its short code
func (s *URLService) GetURL(ctx context.Context, shortCode string) (*domain.URL, error) {
	url, err := s.repo.GetByShortCode(shortCode)
	if err != nil {
		return nil, err
	}

	if url == nil {
		return nil, &domain.ErrURLNotFound{ShortCode: shortCode}
	}

	if !url.IsActive {
		return nil, &domain.ErrURLNotFound{ShortCode: shortCode}
	}

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return nil, &domain.ErrURLExpired{ShortCode: shortCode}
	}

	return url, nil
}

// ListUserURLs retrieves all URLs for a given user
func (s *URLService) ListUserURLs(ctx context.Context, userID string) ([]domain.URL, error) {
	return s.repo.GetByUserID(userID)
}

// DeleteURL deletes a URL by its ID and user ID
func (s *URLService) DeleteURL(ctx context.Context, id int64, userID string) error {
	return s.repo.Delete(id, userID)
}

// RecordClick increments the click count for a URL
func (s *URLService) RecordClick(urlID int64) error {
	return s.repo.IncrementClickCount(urlID)
}
