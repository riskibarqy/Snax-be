# Link Shortener Microservices

A modern link shortener application built with Go, following DDD principles and microservices architecture.

## Architecture

The application consists of the following microservices:

1. **API Gateway Service**: Main entry point for all requests
2. **URL Shortener Service**: Core service for shortening and redirecting URLs
3. **Analytics Service**: Handles tracking and analytics
4. **Auth Service**: Manages authentication and authorization

## Technologies

- **Language**: Go 1.21+
- **Database**: Neon (PostgreSQL)
- **Cache**: Upstash Redis
- **Message Queue**: Upstash Kafka
- **Authentication**: Clerk
- **Monitoring**: Uptrace
- **Deployment**: fly.io

## Project Structure

```
.
├── api-gateway/        # API Gateway Service
├── url-shortener/      # URL Shortener Service
├── analytics/          # Analytics Service
├── auth/              # Auth Service
├── pkg/               # Shared packages
├── scripts/           # Utility scripts
└── sqlc/              # SQL queries and generated code
```

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and fill in the required values
3. Run `make setup` to install dependencies
4. Run `make migrate` to run database migrations
5. Run `make run` to start all services

## Environment Variables

```env
# Database
NEON_DATABASE_URL=

# Redis
UPSTASH_REDIS_URL=
UPSTASH_REDIS_TOKEN=

# Kafka
UPSTASH_KAFKA_URL=
UPSTASH_KAFKA_USERNAME=
UPSTASH_KAFKA_PASSWORD=

# Clerk
CLERK_SECRET_KEY=

# Uptrace
UPTRACE_DSN=
```

## Development

1. Generate SQLC code: `make sqlc`
2. Run tests: `make test`
3. Run linter: `make lint`

## Deployment

The application is configured for deployment on fly.io. Each service has its own `fly.toml` configuration file.

To deploy:

```bash
flyctl deploy
```

## Rate Limiting

The application implements Redis-based rate limiting with the following rules:
- 100 requests per minute for authenticated users
- 20 requests per minute for unauthenticated users

## License

MIT 