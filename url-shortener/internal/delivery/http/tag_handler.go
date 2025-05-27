package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// HandleAddTag handles adding a tag to a URL
func (h *Handler) HandleAddTag(w http.ResponseWriter, r *http.Request) {
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

// HandleRemoveTag handles removing a tag from a URL
func (h *Handler) HandleRemoveTag(w http.ResponseWriter, r *http.Request) {
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

// HandleGetURLTags handles retrieving all tags for a URL
func (h *Handler) HandleGetURLTags(w http.ResponseWriter, r *http.Request) {
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
