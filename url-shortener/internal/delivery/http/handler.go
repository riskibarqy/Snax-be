package http

import (
	"encoding/json"
	"net/http"
	"net/url"

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

// HandlePublicShorten handles URL shortening requests from unauthenticated users
func (h *Handler) HandlePublicShorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Validate URL
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	// Create short URL without user association
	shortURL, err := h.urlService.CreateShortURL(r.Context(), req.URL, "", nil)

	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	response := struct {
		ShortURL string `json:"shortUrl"`
	}{
		ShortURL: shortURL.ShortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
