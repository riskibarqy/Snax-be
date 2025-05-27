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

// HandleRegisterDomain handles registering a new custom domain
func (h *Handler) HandleRegisterDomain(w http.ResponseWriter, r *http.Request) {
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

// HandleVerifyDomain handles domain verification
func (h *Handler) HandleVerifyDomain(w http.ResponseWriter, r *http.Request) {
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

// HandleListUserDomains handles listing user's custom domains
func (h *Handler) HandleListUserDomains(w http.ResponseWriter, r *http.Request) {
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

// HandleDeleteDomain handles deleting a custom domain
func (h *Handler) HandleDeleteDomain(w http.ResponseWriter, r *http.Request) {
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
