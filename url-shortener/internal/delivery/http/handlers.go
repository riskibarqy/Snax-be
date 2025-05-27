package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/delivery/http/middleware"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type URLHandler struct {
	urlService domain.URLService
}

func NewURLHandler(urlService domain.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
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
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*domain.Claims)
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
		case *domain.ErrURLNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case *domain.ErrURLExpired:
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
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*domain.Claims)
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
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*domain.Claims)
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
		case *domain.ErrURLNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Failed to delete URL", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
