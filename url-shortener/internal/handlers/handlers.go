package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riskibarqy/Snax-be/pkg/middleware"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/db"
)

type URLHandler struct {
	queries *db.Queries
}

func NewURLHandler(queries *db.Queries) *URLHandler {
	return &URLHandler{
		queries: queries,
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

	// Get user ID from context
	claims, ok := r.Context().Value(middleware.SessionContextKey).(*middleware.SessionClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.Subject

	// Generate short code (implementation needed)
	shortCode := generateShortCode()

	// Create URL in database
	url, err := h.queries.CreateURL(r.Context(), db.CreateURLParams{
		OriginalURL: req.URL,
		ShortCode:   shortCode,
		UserID:      userID,
		ExpiresAt:   req.ExpiresAt,
	})
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

	// Get URL from database
	url, err := h.queries.GetURLByShortCode(r.Context(), shortCode)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Increment click count
	_, err = h.queries.IncrementClickCount(r.Context(), url.ID)
	if err != nil {
		// Log error but don't fail the request
		// TODO: Implement proper logging
	}

	// Redirect to original URL
	http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
}

func (h *URLHandler) HandleListURLs(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	claims, ok := r.Context().Value(middleware.SessionContextKey).(*middleware.SessionClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.Subject

	// Get URLs from database
	urls, err := h.queries.GetURLsByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

func (h *URLHandler) HandleDeleteURL(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	claims, ok := r.Context().Value(middleware.SessionContextKey).(*middleware.SessionClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.Subject

	// Get URL ID from path
	urlID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	// Deactivate URL
	_, err = h.queries.DeactivateURL(r.Context(), db.DeactivateURLParams{
		ID:     urlID,
		UserID: userID,
	})
	if err != nil {
		http.Error(w, "Failed to delete URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function to generate short code
func generateShortCode() string {
	// TODO: Implement a proper short code generation algorithm
	// This is a placeholder implementation
	return time.Now().Format("20060102150405")
}
