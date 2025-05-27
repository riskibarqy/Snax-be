package domain

import (
	"context"
	"time"
)

// URL represents a shortened URL
type URL struct {
	ID          int64      `json:"id"`
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	UserID      string     `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	ClickCount  int        `json:"click_count"`
	IsActive    bool       `json:"is_active"`
}

// URLRepository defines the interface for URL data operations
type URLRepository interface {
	Create(url *URL) error
	GetByShortCode(shortCode string) (*URL, error)
	GetByUserID(userID string) ([]URL, error)
	Delete(id int64, userID string) error
	IncrementClickCount(id int64) error
}

// URLService defines the interface for URL business logic
type URLService interface {
	CreateShortURL(ctx context.Context, originalURL string, userID string, expiresAt *time.Time) (*URL, error)
	GetURL(ctx context.Context, shortCode string) (*URL, error)
	ListUserURLs(ctx context.Context, userID string) ([]URL, error)
	DeleteURL(ctx context.Context, id int64, userID string) error
	RecordClick(urlID int64) error
}

// ErrURLNotFound is returned when a URL is not found
type ErrURLNotFound struct {
	ShortCode string
}

func (e *ErrURLNotFound) Error() string {
	return "URL not found: " + e.ShortCode
}

// ErrURLExpired is returned when a URL has expired
type ErrURLExpired struct {
	ShortCode string
}

func (e *ErrURLExpired) Error() string {
	return "URL has expired: " + e.ShortCode
}
