package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	customMiddleware "github.com/riskibarqy/Snax-be/url-shortener/internal/delivery/http/middleware"
)

// SetupRouter configures and returns the router with all endpoints
func SetupRouter(h *Handler, authMiddleware *customMiddleware.AuthMiddleware) *chi.Mux {
	r := chi.NewRouter()

	// CORS middleware - configure it properly!
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Cache preflight for 5 minutes
	}))

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Create rate limiter
	rateLimiter := customMiddleware.NewRateLimiter()

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Public routes (/public/...)
	r.Route("/public", func(r chi.Router) {
		// Apply rate limiting to public shortening endpoint
		r.With(rateLimiter.RateLimit).Post("/shorten", h.HandlePublicShorten)

		// URL shortener redirect endpoint (no rate limit)
		r.Get("/r/{shortCode}", h.HandleRedirect)

		// Public metrics endpoint (if needed)
		r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Public Metrics"))
		})
	})

	// Private routes (/private/...)
	r.Route("/private", func(r chi.Router) {
		// Apply authentication middleware to all private routes
		r.Use(authMiddleware.Authenticate)

		// Apply rate limiting to private URL creation
		r.Route("/urls", func(r chi.Router) {
			r.With(rateLimiter.RateLimit).Post("/", h.HandleShorten)
			r.Get("/", h.HandleListURLs)
			r.Delete("/{id}", h.HandleDeleteURL)

			// URL Analytics
			r.Get("/{id}/analytics", h.HandleGetURLAnalytics)

			// URL Tags
			r.Route("/{id}/tags", func(r chi.Router) {
				r.Get("/", h.HandleGetURLTags)
				r.Post("/", h.HandleAddTag)
				r.Delete("/{tag}", h.HandleRemoveTag)
			})
		})

		// Domain Management
		r.Route("/domains", func(r chi.Router) {
			r.Get("/", h.HandleListUserDomains)
			r.Post("/", h.HandleRegisterDomain)
			r.Post("/{id}/verify", h.HandleVerifyDomain)
			r.Delete("/{id}", h.HandleDeleteDomain)
		})
	})

	return r
}
