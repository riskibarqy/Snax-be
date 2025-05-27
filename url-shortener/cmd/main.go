package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/uptrace/uptrace-go/uptrace"

	"link-shortener/internal/handlers"
	"link-shortener/pkg/common"
	"link-shortener/pkg/db"
	custommiddleware "link-shortener/pkg/middleware"
)

func main() {
	// Load configuration
	config, err := common.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Uptrace
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(config.UptraceDSN),
		uptrace.WithServiceName(config.ServiceName),
		uptrace.WithServiceVersion("1.0.0"),
	)
	defer uptrace.Shutdown(context.Background())

	// Initialize database connection
	dbConn, err := db.Connect(config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Initialize queries
	queries := db.New(dbConn)

	// Initialize rate limiter with Upstash Redis
	rateLimiter := custommiddleware.NewRateLimiter(
		config.RedisURL,
		config.RedisToken,
		100,         // 100 requests
		time.Minute, // per minute
	)

	// Initialize auth middleware
	authMiddleware, err := custommiddleware.NewAuthMiddleware(config.ClerkSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize handlers
	urlHandler := handlers.NewURLHandler(queries)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(rateLimiter.RateLimit)
		r.Use(authMiddleware.OptionalAuth)

		r.Get("/{shortCode}", urlHandler.HandleRedirect)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(rateLimiter.RateLimit)
		r.Use(authMiddleware.Authenticate)

		r.Post("/shorten", urlHandler.HandleShorten)
		r.Get("/urls", urlHandler.HandleListURLs)
		r.Delete("/urls/{id}", urlHandler.HandleDeleteURL)
	})

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.ServicePort),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error during server shutdown: %v\n", err)
		}
	}()

	log.Printf("URL Shortener service starting on port %s\n", config.ServicePort)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Handler implementations will be added in separate files
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// Implementation will be added
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	// Implementation will be added
}

func handleListURLs(w http.ResponseWriter, r *http.Request) {
	// Implementation will be added
}

func handleDeleteURL(w http.ResponseWriter, r *http.Request) {
	// Implementation will be added
}
