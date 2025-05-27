SHELL := /bin/bash
include ../.env

.PHONY: migrate-create migrate-up migrate-down migrate-force

# Create a new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Apply all up migrations
migrate-up:
	migrate -path migrations -database "${NEON_DATABASE_URL}" up

# Roll back one migration
migrate-down:
	migrate -path migrations -database "${NEON_DATABASE_URL}" down 1

# Force a specific version
migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path migrations -database "${NEON_DATABASE_URL}" force $$version

# Install required tools
install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest 