package domain

import "time"

// Analytics represents a URL visit analytics entry
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

// AnalyticsRepository defines the interface for analytics data operations
type AnalyticsRepository interface {
	Create(analytics *Analytics) error
	GetByURLID(urlID int64) ([]Analytics, error)
}

// AnalyticsService defines the interface for analytics business logic
type AnalyticsService interface {
	RecordVisit(urlID int64, visitorIP, userAgent, referer string) error
	GetURLAnalytics(urlID int64) ([]Analytics, error)
}
