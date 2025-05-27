package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database URL
	dbURL := os.Getenv("NEON_DATABASE_URL")
	if dbURL == "" {
		log.Fatal("NEON_DATABASE_URL is not set")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Successfully connected to database")

	// Read migration file
	migration, err := os.ReadFile("internal/db/schema.sql")
	if err != nil {
		log.Fatalf("Error reading migration file: %v", err)
	}

	// Execute migration
	_, err = db.Exec(string(migration))
	if err != nil {
		log.Fatalf("Error executing migration: %v", err)
	}

	fmt.Println("Migration completed successfully")
}
