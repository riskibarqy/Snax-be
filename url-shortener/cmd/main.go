package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riskibarqy/Snax-be/pkg/common"
	"github.com/riskibarqy/Snax-be/pkg/telemetry"
	httphandler "github.com/riskibarqy/Snax-be/url-shortener/internal/delivery/http"
	authmiddleware "github.com/riskibarqy/Snax-be/url-shortener/internal/delivery/http/middleware"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/repository/postgres"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/service"
)

func main() {
	// Load configuration
	config, err := common.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Uptrace
	if config.UptraceDSN != "" {
		if err := telemetry.InitUptrace(config.UptraceDSN); err != nil {
			log.Printf("Failed to initialize Uptrace: %v", err)
		}
		defer telemetry.Shutdown(context.Background())
	}

	// Initialize auth service
	authService, err := service.NewAuthService(config.ClerkSecretKey)
	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}

	// Initialize auth middleware
	authMiddleware := authmiddleware.NewAuthMiddleware(authService)

	// Connect to database
	ctx := context.Background()
	db, err := pgx.Connect(ctx, config.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close(ctx)

	// Initialize repositories
	urlRepo := postgres.NewURLRepository(db)
	analyticsRepo := postgres.NewAnalyticsRepository(db)
	tagRepo := postgres.NewTagRepository(db)
	customDomainRepo := postgres.NewCustomDomainRepository(db)

	// Initialize services
	urlService := service.NewURLService(urlRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	tagService := service.NewTagService(tagRepo)
	customDomainService := service.NewCustomDomainService(customDomainRepo)

	// Initialize handler
	handler := httphandler.NewHandler(urlService, analyticsService, tagService, customDomainService)

	// Setup router using the router.go configuration
	router := httphandler.SetupRouter(handler, authMiddleware)

	// Set up graceful shutdown
	srv := &http.Server{
		Addr:    ":" + config.ServicePort,
		Handler: router,
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)
	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start the service listening for requests.
	go func() {
		log.Printf("Server starting on port %s", config.ServicePort)
		serverErrors <- srv.ListenAndServe()
	}()

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case <-shutdown:
		log.Println("Starting shutdown")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown did not complete in %v: %v", 5*time.Second, err)
			if err := srv.Close(); err != nil {
				log.Printf("Error killing server: %v", err)
			}
		}
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
