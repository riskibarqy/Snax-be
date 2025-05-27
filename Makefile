SHELL := /bin/bash
include .env

.PHONY: setup sqlc test lint migrate-create migrate-up migrate-down migrate-force sync-schema run-api run-shortener run-analytics run-auth docker-build install-tools

# Setup and tools
setup:
	go mod tidy
	go mod download

install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Development
sqlc:
	sqlc generate

test:
	go test -v ./...

lint:
	golangci-lint run

# Run services
run-api:
	SERVICE_NAME=api-gateway go run ./api-gateway/cmd/main.go

run-shortener:
	SERVICE_NAME=url-shortener go run ./url-shortener/cmd/main.go

run-analytics:
	SERVICE_NAME=analytics go run ./analytics/cmd/main.go

run-auth:
	SERVICE_NAME=auth go run ./auth/cmd/main.go

# Database migrations
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir url-shortener/migrations -seq $$name

migrate-up:
	migrate -path url-shortener/migrations -database "${NEON_DATABASE_URL}" up
	$(MAKE) sync-schema

migrate-down:
	migrate -path url-shortener/migrations -database "${NEON_DATABASE_URL}" down 1
	$(MAKE) sync-schema

migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path url-shortener/migrations -database "${NEON_DATABASE_URL}" force $$version
	$(MAKE) sync-schema

# Schema management
sync-schema:
	@echo "Syncing migrations to sqlc schema..."
	@mkdir -p sqlc/schema
	@echo "-- This file is auto-generated from migrations. DO NOT EDIT DIRECTLY." > sqlc/schema/schema.sql
	@echo "-- Last updated: $$(date)" >> sqlc/schema/schema.sql
	@echo "" >> sqlc/schema/schema.sql
	@for f in url-shortener/migrations/*.up.sql; do \
		echo "-- Including migration: $$(basename $$f)" >> sqlc/schema/schema.sql; \
		echo "" >> sqlc/schema/schema.sql; \
		cat $$f >> sqlc/schema/schema.sql; \
		echo "" >> sqlc/schema/schema.sql; \
		echo "" >> sqlc/schema/schema.sql; \
	done
	@echo "Schema sync completed successfully!"

# Docker
docker-build:
	docker build -t link-shortener-api-gateway -f api-gateway/Dockerfile .
	docker build -t link-shortener-url-shortener -f url-shortener/Dockerfile .
	docker build -t link-shortener-analytics -f analytics/Dockerfile .
	docker build -t link-shortener-auth -f auth/Dockerfile . 