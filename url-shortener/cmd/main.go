package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	httphandler "github.com/riskibarqy/Snax-be/url-shortener/internal/delivery/http"
	authmiddleware "github.com/riskibarqy/Snax-be/url-shortener/internal/delivery/http/middleware"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/repository/postgres"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Get database URL from environment
	dbURL := os.Getenv("NEON_DATABASE_URL")
	if dbURL == "" {
		log.Fatal("NEON_DATABASE_URL environment variable is required")
	}

	// Get Clerk secret key from environment
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable is required")
	}

	// Initialize auth service
	authService, err := service.NewAuthService(clerkSecretKey)
	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}

	// Initialize auth middleware
	authMiddleware := authmiddleware.NewAuthMiddleware(authService)

	// Connect to database
	ctx := context.Background()
	db, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close(ctx)

	// Initialize repository
	urlRepo := postgres.NewURLRepository(db)

	// Initialize service
	urlService := service.NewURLService(urlRepo)

	// Initialize handler
	urlHandler := httphandler.NewURLHandler(urlService)

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)

	// Public routes
	r.Get("/{shortCode}", urlHandler.HandleRedirect)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		r.Post("/shorten", urlHandler.HandleShorten)
		r.Get("/urls", urlHandler.HandleListURLs)
		r.Delete("/urls/{id}", urlHandler.HandleDeleteURL)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
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
