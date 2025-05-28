package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DatabaseURL string

	// Redis
	RedisURL   string
	RedisToken string

	// Clerk
	ClerkSecretKey string

	// Uptrace
	UptraceDSN string

	// Service specific
	ServicePort string
	ServiceName string
}

func LoadConfig() (*Config, error) {
	// Get environment
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	// Try to load environment file, but don't fail if it doesn't exist
	envFile := fmt.Sprintf(".env.%s", env)
	_ = godotenv.Load(envFile)

	config := &Config{
		// Database
		DatabaseURL: os.Getenv("NEON_DATABASE_URL"),

		// Redis
		RedisURL:   os.Getenv("UPSTASH_REDIS_URL"),
		RedisToken: os.Getenv("UPSTASH_REDIS_TOKEN"),

		// Clerk
		ClerkSecretKey: os.Getenv("CLERK_SECRET_KEY"),

		// Uptrace
		UptraceDSN: os.Getenv("UPTRACE_DSN"),

		// Service specific
		ServicePort: os.Getenv("PORT"),
		ServiceName: os.Getenv("SERVICE_NAME"),
	}

	// Set default port if not specified
	if config.ServicePort == "" {
		config.ServicePort = "8080"
	}

	// Validate required configurations
	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("NEON_DATABASE_URL is required")
	}

	if config.ClerkSecretKey == "" {
		return nil, fmt.Errorf("CLERK_SECRET_KEY is required")
	}

	return config, nil
}
