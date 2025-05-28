package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/delivery/http/middleware"
	internalDomain "github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// HandleShorten handles the creation of short URLs
func (h *Handler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleShorten")
	defer span.End()

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := ctx.Value(middleware.SessionContextKey).(*internalDomain.Claims)
	if !ok {
		span.SetAttributes(attribute.String("error", "unauthorized"))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	span.SetAttributes(
		attribute.String("user_id", claims.Subject),
		attribute.String("original_url", req.URL),
	)

	url, err := h.urlService.CreateShortURL(ctx, req.URL, claims.Subject, req.ExpiresAt)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	span.SetAttributes(attribute.String("short_code", url.ShortCode))

	response := ShortenResponse{
		ShortCode: url.ShortCode,
		URL:       url.OriginalURL,
		ExpiresAt: url.ExpiresAt,
		CreatedAt: url.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleRedirect handles URL redirection
func (h *Handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleRedirect")
	defer span.End()

	shortCode := chi.URLParam(r, "shortCode")
	span.SetAttributes(attribute.String("short_code", shortCode))

	url, err := h.urlService.GetURL(ctx, shortCode)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
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

	span.SetAttributes(
		attribute.String("original_url", url.OriginalURL),
		attribute.Int64("url_id", url.ID),
	)

	// Record click asynchronously
	go h.urlService.RecordClick(url.ID)

	http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
}

// HandleListURLs handles listing user's URLs
func (h *Handler) HandleListURLs(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleListURLs")
	defer span.End()

	claims, ok := ctx.Value(middleware.SessionContextKey).(*internalDomain.Claims)
	if !ok {
		span.SetAttributes(attribute.String("error", "unauthorized"))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	span.SetAttributes(attribute.String("user_id", claims.Subject))

	urls, err := h.urlService.ListUserURLs(ctx, claims.Subject)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Failed to fetch URLs", http.StatusInternalServerError)
		return
	}

	span.SetAttributes(attribute.Int("url_count", len(urls)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

// HandleDeleteURL handles URL deletion
func (h *Handler) HandleDeleteURL(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleDeleteURL")
	defer span.End()

	claims, ok := ctx.Value(middleware.SessionContextKey).(*internalDomain.Claims)
	if !ok {
		span.SetAttributes(attribute.String("error", "unauthorized"))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	urlID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	span.SetAttributes(
		attribute.String("user_id", claims.Subject),
		attribute.Int64("url_id", urlID),
	)

	err = h.urlService.DeleteURL(ctx, urlID, claims.Subject)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
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
