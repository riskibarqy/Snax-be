package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/delivery/http/middleware"
	internalDomain "github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleShorten")
	defer span.End()

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get claims from context
	claims, ok := ctx.Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
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

func (h *URLHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
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

func (h *URLHandler) HandleListURLs(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleListURLs")
	defer span.End()

	claims, ok := ctx.Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
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

func (h *URLHandler) HandleDeleteURL(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleDeleteURL")
	defer span.End()

	claims, ok := ctx.Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
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

// Analytics handlers
func (h *URLHandler) HandleGetURLAnalytics(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleGetURLAnalytics")
	defer span.End()

	urlID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.Int64("url_id", urlID))

	analytics, err := h.analyticsService.GetURLAnalytics(ctx, urlID)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Failed to fetch analytics", http.StatusInternalServerError)
		return
	}

	span.SetAttributes(attribute.Int("analytics_count", len(analytics)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// Tag handlers
type AddTagRequest struct {
	Tag string `json:"tag"`
}

func (h *URLHandler) HandleAddTag(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleAddTag")
	defer span.End()

	urlID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	var req AddTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	span.SetAttributes(
		attribute.Int64("url_id", urlID),
		attribute.String("tag", req.Tag),
	)

	err = h.tagService.AddTagToURL(ctx, urlID, req.Tag)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Failed to add tag", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *URLHandler) HandleRemoveTag(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleRemoveTag")
	defer span.End()

	urlID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	tagName := chi.URLParam(r, "tag")
	span.SetAttributes(
		attribute.Int64("url_id", urlID),
		attribute.String("tag", tagName),
	)

	err = h.tagService.RemoveTagFromURL(ctx, urlID, tagName)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Failed to remove tag", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *URLHandler) HandleGetURLTags(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleGetURLTags")
	defer span.End()

	urlID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.Int64("url_id", urlID))

	tags, err := h.tagService.GetURLTags(ctx, urlID)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Failed to fetch tags", http.StatusInternalServerError)
		return
	}

	span.SetAttributes(attribute.Int("tag_count", len(tags)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// Custom domain handlers
type RegisterDomainRequest struct {
	Domain string `json:"domain"`
}

func (h *URLHandler) HandleRegisterDomain(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleRegisterDomain")
	defer span.End()

	claims, ok := ctx.Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
	if !ok {
		span.SetAttributes(attribute.String("error", "unauthorized"))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req RegisterDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	span.SetAttributes(
		attribute.String("user_id", claims.Subject),
		attribute.String("domain", req.Domain),
	)

	domain, err := h.customDomainService.RegisterDomain(ctx, req.Domain, claims.Subject)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		switch err.(type) {
		case *internalDomain.ErrDomainNotFound:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Failed to register domain", http.StatusInternalServerError)
		}
		return
	}

	span.SetAttributes(attribute.Int64("domain_id", domain.ID))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain)
}

func (h *URLHandler) HandleVerifyDomain(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleVerifyDomain")
	defer span.End()

	domainID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Invalid domain ID", http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.Int64("domain_id", domainID))

	err = h.customDomainService.VerifyDomain(ctx, domainID)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
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
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleListUserDomains")
	defer span.End()

	claims, ok := ctx.Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
	if !ok {
		span.SetAttributes(attribute.String("error", "unauthorized"))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	span.SetAttributes(attribute.String("user_id", claims.Subject))

	domains, err := h.customDomainService.GetUserDomains(ctx, claims.Subject)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Failed to fetch domains", http.StatusInternalServerError)
		return
	}

	span.SetAttributes(attribute.Int("domain_count", len(domains)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domains)
}

func (h *URLHandler) HandleDeleteDomain(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("url-handler").Start(r.Context(), "HandleDeleteDomain")
	defer span.End()

	claims, ok := ctx.Value(middleware.ClaimsContextKey).(*internalDomain.Claims)
	if !ok {
		span.SetAttributes(attribute.String("error", "unauthorized"))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domainID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, "Invalid domain ID", http.StatusBadRequest)
		return
	}

	span.SetAttributes(
		attribute.String("user_id", claims.Subject),
		attribute.Int64("domain_id", domainID),
	)

	err = h.customDomainService.DeleteDomain(ctx, domainID, claims.Subject)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
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
