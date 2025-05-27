package domain

import (
	"time"
)

// URL represents the core URL entity
type URL struct {
	ID          int64      `json:"id"`
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	UserID      string     `json:"user_id"`
	ClickCount  int64      `json:"click_count"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	IsActive    bool       `json:"is_active"`
}

// URLRepository defines the interface for URL data operations
type URLRepository interface {
	Create(url *URL) error
	GetByShortCode(shortCode string) (*URL, error)
	GetByUserID(userID string) ([]URL, error)
	Deactivate(id int64, userID string) error
	IncrementClickCount(id int64) (int64, error)
}

// URLService defines the interface for URL business logic
type URLService interface {
	CreateShortURL(originalURL, userID string, expiresAt *time.Time) (*URL, error)
	GetURL(shortCode string) (*URL, error)
	ListUserURLs(userID string) ([]URL, error)
	DeleteURL(id int64, userID string) error
	RecordClick(id int64) error
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
