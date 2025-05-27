package http

import (
	internalDomain "github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

// Handler contains all the dependencies for HTTP handlers
type Handler struct {
	urlService          internalDomain.URLService
	analyticsService    internalDomain.AnalyticsService
	tagService          internalDomain.TagService
	customDomainService internalDomain.CustomDomainService
}

// NewHandler creates a new Handler instance
func NewHandler(
	urlService internalDomain.URLService,
	analyticsService internalDomain.AnalyticsService,
	tagService internalDomain.TagService,
	customDomainService internalDomain.CustomDomainService,
) *Handler {
	return &Handler{
		urlService:          urlService,
		analyticsService:    analyticsService,
		tagService:          tagService,
		customDomainService: customDomainService,
	}
}
