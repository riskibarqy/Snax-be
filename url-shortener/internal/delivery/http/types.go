package http

import "time"

// URL-related types
type ShortenRequest struct {
	URL       string     `json:"url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type ShortenResponse struct {
	ShortCode string     `json:"short_code"`
	URL       string     `json:"url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// Tag-related types
type AddTagRequest struct {
	Tag string `json:"tag"`
}

// Custom domain-related types
type RegisterDomainRequest struct {
	Domain string `json:"domain"`
}
