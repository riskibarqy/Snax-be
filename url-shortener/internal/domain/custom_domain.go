package domain

import (
	"context"
	"time"
)

// CustomDomain represents a custom domain for URL shortening
type CustomDomain struct {
	ID        int64     `json:"id"`
	Domain    string    `json:"domain"`
	UserID    string    `json:"user_id"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
}

// CustomDomainRepository defines the interface for custom domain data operations
type CustomDomainRepository interface {
	Create(ctx context.Context, domain *CustomDomain) error
	GetByID(ctx context.Context, id int64) (*CustomDomain, error)
	GetByDomain(ctx context.Context, domain string) (*CustomDomain, error)
	GetByUserID(ctx context.Context, userID string) ([]CustomDomain, error)
	Delete(ctx context.Context, id int64, userID string) error
	VerifyDomain(ctx context.Context, id int64) error
}

// CustomDomainService defines the interface for custom domain business logic
type CustomDomainService interface {
	RegisterDomain(ctx context.Context, domain string, userID string) (*CustomDomain, error)
	GetUserDomains(ctx context.Context, userID string) ([]CustomDomain, error)
	DeleteDomain(ctx context.Context, id int64, userID string) error
	VerifyDomain(ctx context.Context, id int64) error
}

// ErrDomainNotFound is returned when a custom domain is not found
type ErrDomainNotFound struct {
	Domain string
}

func (e *ErrDomainNotFound) Error() string {
	return "Domain not found: " + e.Domain
}

// ErrDomainNotVerified is returned when a custom domain is not verified
type ErrDomainNotVerified struct {
	Domain string
}

func (e *ErrDomainNotVerified) Error() string {
	return "Domain not verified: " + e.Domain
}
