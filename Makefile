.PHONY: setup sqlc test lint run-api run-shortener run-analytics run-auth

setup:
	go mod tidy
	go mod download

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
migrate:
	migrate -database "${NEON_DATABASE_URL}" -path sqlc/schema up

migrate-down:
	migrate -database "${NEON_DATABASE_URL}" -path sqlc/schema down

# Docker builds
docker-build:
	docker build -t link-shortener-api-gateway -f api-gateway/Dockerfile .
	docker build -t link-shortener-url-shortener -f url-shortener/Dockerfile .
	docker build -t link-shortener-analytics -f analytics/Dockerfile .
	docker build -t link-shortener-auth -f auth/Dockerfile . 