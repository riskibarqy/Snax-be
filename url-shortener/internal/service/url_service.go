package service

import (
	"math/rand"
	"time"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type urlService struct {
	repo domain.URLRepository
}

// NewURLService creates a new URL service
func NewURLService(repo domain.URLRepository) domain.URLService {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) CreateShortURL(originalURL, userID string, expiresAt *time.Time) (*domain.URL, error) {
	url := &domain.URL{
		ShortCode:   generateShortCode(),
		OriginalURL: originalURL,
		UserID:      userID,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
		IsActive:    true,
	}

	if err := s.repo.Create(url); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *urlService) GetURL(shortCode string) (*domain.URL, error) {
	url, err := s.repo.GetByShortCode(shortCode)
	if err != nil {
		return nil, &domain.ErrURLNotFound{ShortCode: shortCode}
	}

	if !url.IsActive {
		return nil, &domain.ErrURLNotFound{ShortCode: shortCode}
	}

	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		return nil, &domain.ErrURLExpired{ShortCode: shortCode}
	}

	return url, nil
}

func (s *urlService) ListUserURLs(userID string) ([]domain.URL, error) {
	return s.repo.GetByUserID(userID)
}

func (s *urlService) DeleteURL(id int64, userID string) error {
	return s.repo.Deactivate(id, userID)
}

func (s *urlService) RecordClick(id int64) error {
	_, err := s.repo.IncrementClickCount(id)
	return err
}

// Helper function to generate short code
func generateShortCode() string {
	const (
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		length  = 6
	)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
