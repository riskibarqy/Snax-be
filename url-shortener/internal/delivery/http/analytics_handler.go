package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// HandleGetURLAnalytics handles retrieving analytics for a URL
func (h *Handler) HandleGetURLAnalytics(w http.ResponseWriter, r *http.Request) {
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
