package domain

import "time"

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
	Create(domain *CustomDomain) error
	GetByDomain(domain string) (*CustomDomain, error)
	GetByUserID(userID string) ([]CustomDomain, error)
	UpdateVerificationStatus(id int64, verified bool) error
	Delete(id int64, userID string) error
}

// CustomDomainService defines the interface for custom domain business logic
type CustomDomainService interface {
	RegisterDomain(domain, userID string) (*CustomDomain, error)
	VerifyDomain(id int64) error
	GetUserDomains(userID string) ([]CustomDomain, error)
	DeleteDomain(id int64, userID string) error
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
