package service

import (
	"strings"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type analyticsService struct {
	repo domain.AnalyticsRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(repo domain.AnalyticsRepository) domain.AnalyticsService {
	return &analyticsService{
		repo: repo,
	}
}

func (s *analyticsService) RecordVisit(urlID int64, visitorIP, userAgent, referer string) error {
	analytics := &domain.Analytics{
		URLID:      urlID,
		VisitorIP:  visitorIP,
		UserAgent:  userAgent,
		Referer:    referer,
		DeviceType: detectDeviceType(userAgent),
		// TODO: Implement country code detection using GeoIP service
		CountryCode: "UN",
	}

	return s.repo.Create(analytics)
}

func (s *analyticsService) GetURLAnalytics(urlID int64) ([]domain.Analytics, error) {
	return s.repo.GetByURLID(urlID)
}

// Helper function to detect device type from user agent
func detectDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone"):
		return "mobile"
	case strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad"):
		return "tablet"
	default:
		return "desktop"
	}
}
