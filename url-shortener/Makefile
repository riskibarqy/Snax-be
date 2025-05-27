SHELL := /bin/bash
include ../.env

.PHONY: migrate-create migrate-up migrate-down migrate-force sync-schema

# Create a new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Apply all up migrations
migrate-up:
	migrate -path migrations -database "${NEON_DATABASE_URL}" up
	$(MAKE) sync-schema

# Roll back one migration
migrate-down:
	migrate -path migrations -database "${NEON_DATABASE_URL}" down 1
	$(MAKE) sync-schema

# Force a specific version
migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path migrations -database "${NEON_DATABASE_URL}" force $$version
	$(MAKE) sync-schema

# Sync migrations to sqlc schema
sync-schema:
	@echo "Syncing migrations to sqlc schema..."
	@mkdir -p ../sqlc/schema
	@echo "-- This file is auto-generated from migrations. DO NOT EDIT DIRECTLY." > ../sqlc/schema/schema.sql
	@echo "-- Last updated: $$(date)" >> ../sqlc/schema/schema.sql
	@echo "" >> ../sqlc/schema/schema.sql
	@for f in migrations/*.up.sql; do \
		echo "-- Including migration: $$(basename $$f)" >> ../sqlc/schema/schema.sql; \
		echo "" >> ../sqlc/schema/schema.sql; \
		cat $$f >> ../sqlc/schema/schema.sql; \
		echo "" >> ../sqlc/schema/schema.sql; \
		echo "" >> ../sqlc/schema/schema.sql; \
	done
	@echo "Schema sync completed successfully!"

# Install required tools
install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest 