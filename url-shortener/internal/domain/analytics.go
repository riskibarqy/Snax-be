package domain

import (
	"context"
	"time"
)

// Analytics represents a visit to a shortened URL
type Analytics struct {
	ID          int64     `json:"id"`
	URLID       int64     `json:"url_id"`
	VisitorIP   string    `json:"visitor_ip"`
	UserAgent   string    `json:"user_agent"`
	Referer     string    `json:"referer"`
	Timestamp   time.Time `json:"timestamp"`
	CountryCode string    `json:"country_code"`
	DeviceType  string    `json:"device_type"`
}

// AnalyticsService defines the interface for analytics operations
type AnalyticsService interface {
	RecordVisit(ctx context.Context, urlID int64, visitorIP, userAgent, referer string) error
	GetURLAnalytics(ctx context.Context, urlID int64) ([]Analytics, error)
}

// AnalyticsRepository defines the interface for analytics storage operations
type AnalyticsRepository interface {
	Create(ctx context.Context, analytics *Analytics) error
	GetByURLID(ctx context.Context, urlID int64) ([]Analytics, error)
}
