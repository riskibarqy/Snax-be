package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type urlService struct {
	urlRepo domain.URLRepository
}

// NewURLService creates a new URL service
func NewURLService(urlRepo domain.URLRepository) domain.URLService {
	return &urlService{
		urlRepo: urlRepo,
	}
}

// CreateShortURL creates a new shortened URL
func (s *urlService) CreateShortURL(ctx context.Context, originalURL string, userID string, expiresAt *time.Time) (*domain.URL, error) {
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

	if err := s.urlRepo.Create(url); err != nil {
		return nil, err
	}

	return url, nil
}

// GetURL retrieves a URL by its short code
func (s *urlService) GetURL(ctx context.Context, shortCode string) (*domain.URL, error) {
	url, err := s.urlRepo.GetByShortCode(shortCode)
	if err != nil {
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

// ListUserURLs retrieves all URLs for a user
func (s *urlService) ListUserURLs(ctx context.Context, userID string) ([]domain.URL, error) {
	return s.urlRepo.GetByUserID(userID)
}

// DeleteURL deletes a URL
func (s *urlService) DeleteURL(ctx context.Context, id int64, userID string) error {
	return s.urlRepo.Delete(id, userID)
}

// RecordClick increments the click count for a URL
func (s *urlService) RecordClick(urlID int64) error {
	return s.urlRepo.IncrementClickCount(urlID)
}

// generateShortCode generates a random short code
func generateShortCode() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:6], nil
}
