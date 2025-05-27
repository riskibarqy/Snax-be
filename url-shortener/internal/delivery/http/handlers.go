package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/delivery/http/middleware"
	internalDomain "github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type URLHandler struct {
	urlService          internalDomain.URLService
	analyticsService    internalDomain.AnalyticsService
	tagService          internalDomain.TagService
	customDomainService internalDomain.CustomDomainService
}

func NewURLHandler(
	urlService internalDomain.URLService,
	analyticsService internalDomain.AnalyticsService,
	tagService internalDomain.TagService,
	customDomainService internalDomain.CustomDomainService,
) *URLHandler {
	return &URLHandler{
		urlService:          urlService,
		analyticsService:    analyticsService,
		tagService:          tagService,
		customDomainService: customDomainService,
	}
}

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

func (h *URLHandler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get claims from context
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	url, err := h.urlService.CreateShortURL(req.URL, claims.Subject, req.ExpiresAt)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	response := ShortenResponse{
		ShortCode: url.ShortCode,
		URL:       url.OriginalURL,
		ExpiresAt: url.ExpiresAt,
		CreatedAt: url.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *URLHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")

	url, err := h.urlService.GetURL(shortCode)
	if err != nil {
		switch err.(type) {
		case *internalDomain.ErrURLNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case *internalDomain.ErrURLExpired:
			http.Error(w, err.Error(), http.StatusGone)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Record click asynchronously
	go h.urlService.RecordClick(url.ID)

	http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
}

func (h *URLHandler) HandleListURLs(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	urls, err := h.urlService.ListUserURLs(claims.Subject)
	if err != nil {
		http.Error(w, "Failed to fetch URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

func (h *URLHandler) HandleDeleteURL(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	urlID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	err = h.urlService.DeleteURL(urlID, claims.Subject)
	if err != nil {
		switch err.(type) {
		case *internalDomain.ErrURLNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Failed to delete URL", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Analytics handlers
func (h *URLHandler) HandleGetURLAnalytics(w http.ResponseWriter, r *http.Request) {
	urlID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	analytics, err := h.analyticsService.GetURLAnalytics(urlID)
	if err != nil {
		http.Error(w, "Failed to fetch analytics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// Tag handlers
type AddTagRequest struct {
	Tag string `json:"tag"`
}

func (h *URLHandler) HandleAddTag(w http.ResponseWriter, r *http.Request) {
	urlID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	var req AddTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.tagService.AddTagToURL(urlID, req.Tag)
	if err != nil {
		http.Error(w, "Failed to add tag", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *URLHandler) HandleRemoveTag(w http.ResponseWriter, r *http.Request) {
	urlID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	tagName := chi.URLParam(r, "tag")
	err = h.tagService.RemoveTagFromURL(urlID, tagName)
	if err != nil {
		http.Error(w, "Failed to remove tag", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *URLHandler) HandleGetURLTags(w http.ResponseWriter, r *http.Request) {
	urlID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	tags, err := h.tagService.GetURLTags(urlID)
	if err != nil {
		http.Error(w, "Failed to fetch tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// Custom domain handlers
type RegisterDomainRequest struct {
	Domain string `json:"domain"`
}

func (h *URLHandler) HandleRegisterDomain(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req RegisterDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	domain, err := h.customDomainService.RegisterDomain(req.Domain, claims.Subject)
	if err != nil {
		switch err.(type) {
		case *internalDomain.ErrDomainNotFound:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Failed to register domain", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain)
}

func (h *URLHandler) HandleVerifyDomain(w http.ResponseWriter, r *http.Request) {
	domainID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid domain ID", http.StatusBadRequest)
		return
	}

	err = h.customDomainService.VerifyDomain(domainID)
	if err != nil {
		switch err.(type) {
		case *internalDomain.ErrDomainNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Failed to verify domain", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *URLHandler) HandleListUserDomains(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domains, err := h.customDomainService.GetUserDomains(claims.Subject)
	if err != nil {
		http.Error(w, "Failed to fetch domains", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domains)
}

func (h *URLHandler) HandleDeleteDomain(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domainID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid domain ID", http.StatusBadRequest)
		return
	}

	err = h.customDomainService.DeleteDomain(domainID, claims.Subject)
	if err != nil {
		switch err.(type) {
		case *internalDomain.ErrDomainNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Failed to delete domain", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
