package common

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

	// QStash
	QStashToken          string
	QStashSigningKey     string
	QStashNextSigningKey string
	QStashEndpoint       string

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

	// Load environment file
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		return nil, fmt.Errorf("error loading %s: %v", envFile, err)
	}

	config := &Config{
		// Database
		DatabaseURL: os.Getenv("NEON_DATABASE_URL"),

		// Redis
		RedisURL:   os.Getenv("UPSTASH_REDIS_URL"),
		RedisToken: os.Getenv("UPSTASH_REDIS_TOKEN"),

		// QStash
		QStashToken:          os.Getenv("QSTASH_TOKEN"),
		QStashSigningKey:     os.Getenv("QSTASH_SIGNING_KEY"),
		QStashNextSigningKey: os.Getenv("QSTASH_NEXT_SIGNING_KEY"),
		QStashEndpoint:       os.Getenv("QSTASH_ENDPOINT"),

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

	return config, nil
}
