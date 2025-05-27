package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/riskibarqy/Snax-be/pkg/middleware"
)

type URLHandler struct {
	db *pgx.Conn
}

func NewURLHandler(db *pgx.Conn) *URLHandler {
	return &URLHandler{
		db: db,
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

type URL struct {
	ID          int64      `json:"id"`
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	UserID      string     `json:"user_id"`
	ClickCount  int64      `json:"click_count"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	IsActive    bool       `json:"is_active"`
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

	// Generate short code
	shortCode := generateShortCode()

	// Create URL in database
	var url URL
	err := h.db.QueryRow(r.Context(),
		`INSERT INTO urls (short_code, original_url, user_id, expires_at, created_at, is_active)
		VALUES ($1, $2, $3, $4, NOW(), true)
		RETURNING id, short_code, original_url, user_id, click_count, expires_at, created_at, is_active`,
		shortCode, req.URL, userID, req.ExpiresAt,
	).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.UserID, &url.ClickCount, &url.ExpiresAt, &url.CreatedAt, &url.IsActive)

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
	var url URL
	err := h.db.QueryRow(r.Context(),
		`SELECT id, short_code, original_url, user_id, click_count, expires_at, created_at, is_active
		FROM urls WHERE short_code = $1 AND is_active = true`,
		shortCode,
	).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.UserID, &url.ClickCount, &url.ExpiresAt, &url.CreatedAt, &url.IsActive)

	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Check if URL has expired
	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		http.Error(w, "URL has expired", http.StatusGone)
		return
	}

	// Increment click count
	_, err = h.db.Exec(r.Context(),
		`UPDATE urls SET click_count = click_count + 1 WHERE id = $1`,
		url.ID,
	)
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
	rows, err := h.db.Query(r.Context(),
		`SELECT id, short_code, original_url, user_id, click_count, expires_at, created_at, is_active
		FROM urls WHERE user_id = $1 AND is_active = true ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		http.Error(w, "Failed to fetch URLs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		err := rows.Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.UserID, &url.ClickCount, &url.ExpiresAt, &url.CreatedAt, &url.IsActive)
		if err != nil {
			http.Error(w, "Failed to fetch URLs", http.StatusInternalServerError)
			return
		}
		urls = append(urls, url)
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
	result, err := h.db.Exec(r.Context(),
		`UPDATE urls SET is_active = false WHERE id = $1 AND user_id = $2`,
		urlID, userID,
	)
	if err != nil {
		http.Error(w, "Failed to delete URL", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function to generate short code
func generateShortCode() string {
	const (
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		length  = 6
	)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
