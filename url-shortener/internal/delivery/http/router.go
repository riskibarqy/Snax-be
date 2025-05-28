package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	customMiddleware "github.com/riskibarqy/Snax-be/url-shortener/internal/delivery/http/middleware"
)

// SetupRouter configures and returns the router with all endpoints
func SetupRouter(h *Handler, authMiddleware *customMiddleware.AuthMiddleware) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Public routes
	r.Get("/{shortCode}", h.HandleRedirect)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		// URL routes
		r.Post("/api/urls", h.HandleShorten)
		r.Get("/api/urls", h.HandleListURLs)
		r.Delete("/api/urls/{id}", h.HandleDeleteURL)

		// Analytics routes
		r.Get("/api/urls/{id}/analytics", h.HandleGetURLAnalytics)

		// Tag routes
		r.Get("/api/urls/{id}/tags", h.HandleGetURLTags)
		r.Post("/api/urls/{id}/tags", h.HandleAddTag)
		r.Delete("/api/urls/{id}/tags/{tag}", h.HandleRemoveTag)

		// Custom domain routes
		r.Get("/api/domains", h.HandleListUserDomains)
		r.Post("/api/domains", h.HandleRegisterDomain)
		r.Post("/api/domains/{id}/verify", h.HandleVerifyDomain)
		r.Delete("/api/domains/{id}", h.HandleDeleteDomain)
	})

	return r
}
