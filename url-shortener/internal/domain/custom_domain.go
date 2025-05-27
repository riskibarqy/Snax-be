package domain

import (
	"context"
	"fmt"
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

// CustomDomainService defines the interface for custom domain operations
type CustomDomainService interface {
	RegisterDomain(ctx context.Context, domain string, userID string) (*CustomDomain, error)
	GetUserDomains(ctx context.Context, userID string) ([]CustomDomain, error)
	DeleteDomain(ctx context.Context, id int64, userID string) error
	VerifyDomain(ctx context.Context, id int64) error
}

// CustomDomainRepository defines the interface for custom domain storage operations
type CustomDomainRepository interface {
	Create(ctx context.Context, domain *CustomDomain) error
	GetByID(ctx context.Context, id int64) (*CustomDomain, error)
	GetByDomain(ctx context.Context, domain string) (*CustomDomain, error)
	GetByUserID(ctx context.Context, userID string) ([]CustomDomain, error)
	Delete(ctx context.Context, id int64, userID string) error
	VerifyDomain(ctx context.Context, id int64) error
}

// ErrDomainNotFound is returned when a custom domain is not found
type ErrDomainNotFound struct {
	Domain string
}

func (e *ErrDomainNotFound) Error() string {
	return fmt.Sprintf("Custom domain %s not found", e.Domain)
}

// ErrDomainAlreadyExists is returned when a custom domain already exists
type ErrDomainAlreadyExists struct {
	Domain string
}

func (e *ErrDomainAlreadyExists) Error() string {
	return fmt.Sprintf("Custom domain %s already exists", e.Domain)
}
