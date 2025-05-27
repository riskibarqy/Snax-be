package service

import (
	"context"
	"time"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type AnalyticsService struct {
	repo domain.AnalyticsRepository
}

// New creates a new analytics service
func NewAnalyticsService(repo domain.AnalyticsRepository) domain.AnalyticsService {
	return &AnalyticsService{
		repo: repo,
	}
}

// RecordVisit records a new visit to a URL
func (s *AnalyticsService) RecordVisit(ctx context.Context, urlID int64, visitorIP, userAgent, referer string) error {
	analytics := &domain.Analytics{
		URLID:       urlID,
		VisitorIP:   visitorIP,
		UserAgent:   userAgent,
		Referer:     referer,
		Timestamp:   time.Now(),
		CountryCode: "Unknown", // This could be enhanced with IP geolocation
		DeviceType:  "Unknown", // This could be enhanced with user agent parsing
	}

	return s.repo.Create(ctx, analytics)
}

// GetURLAnalytics retrieves analytics for a specific URL
func (s *AnalyticsService) GetURLAnalytics(ctx context.Context, urlID int64) ([]domain.Analytics, error) {
	return s.repo.GetByURLID(ctx, urlID)
}
