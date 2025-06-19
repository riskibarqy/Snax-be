package http

import (
	"encoding/json"
	"net/http"
	"net/url"

	internalDomain "github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/utils"
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

	// Decode and validate body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		msg := "Invalid request body"
		if req.URL == "" {
			msg = "URL is required"
		}
		utils.RespondError(w, http.StatusBadRequest, msg, err)
		return
	}

	// Validate URL format
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid URL format", err)
		return
	}

	// Try creating the short URL
	shortURL, err := h.urlService.CreateShortURL(r.Context(), req.URL, "", nil)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create short URL", err)
		return
	}

	// Success
	utils.JSONResponse(utils.JSONOpts{
		W:      w,
		Status: http.StatusOK,
		Data: map[string]any{
			"shortUrl": shortURL.ShortCode,
		},
	})
}
