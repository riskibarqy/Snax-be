package domain

// Claims represents the standard JWT claims
type Claims struct {
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	TokenID   string `json:"jti"`
}

// AuthService defines the interface for authentication operations
type AuthService interface {
	ValidateToken(token string) (*Claims, error)
}

// ErrInvalidToken is returned when a token is invalid
type ErrInvalidToken struct {
	Message string
}

func (e *ErrInvalidToken) Error() string {
	return e.Message
}
