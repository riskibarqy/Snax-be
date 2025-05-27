# URL Shortener Service

A modern URL shortening service built with Go, featuring custom domains, analytics, and tag management.

## Features

- URL shortening with custom expiration
- Custom domain support
- Click analytics and tracking
- URL tagging and categorization
- JWT-based authentication
- OpenTelemetry integration with Uptrace

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations
- [sqlc](https://sqlc.dev/) for type-safe SQL
- [golangci-lint](https://golangci-lint.run/) for linting

## Setup

1. Install required tools:
```bash
make install-tools
```

2. Set up environment variables:
   Create a `.env.development` file with:
```env
NEON_DATABASE_URL=postgresql://username:password@localhost:5432/dbname?sslmode=disable
CLERK_SECRET_KEY=your_clerk_secret
UPTRACE_DSN=your_uptrace_dsn  # Optional
```

3. Initialize the database:
```bash
make migrate-up
```

4. Start the service:
```bash
make run-shortener
```

## Project Structure

```
.
├── url-shortener/
│   ├── cmd/                 # Application entry point
│   ├── internal/
│   │   ├── domain/         # Domain models and interfaces
│   │   ├── repository/     # Data access layer
│   │   ├── service/        # Business logic
│   │   └── delivery/       # HTTP handlers and middleware
│   └── migrations/         # Database migrations
├── pkg/
│   ├── common/             # Shared utilities
│   ├── db/                 # Generated SQL code
│   └── telemetry/         # OpenTelemetry setup
└── sqlc/                   # SQL queries and schema
```

## API Endpoints

### Public Endpoints
- `GET /{shortCode}` - Redirect to original URL

### Protected Endpoints (Requires Authentication)
- `POST /api/urls` - Create short URL
- `GET /api/urls` - List user's URLs
- `DELETE /api/urls/{id}` - Delete URL
- `GET /api/urls/{id}/analytics` - Get URL analytics
- `GET /api/urls/{id}/tags` - Get URL tags
- `POST /api/urls/{id}/tags` - Add tag to URL
- `DELETE /api/urls/{id}/tags/{tag}` - Remove tag from URL
- `GET /api/domains` - List user's custom domains
- `POST /api/domains` - Register new domain
- `POST /api/domains/{id}/verify` - Verify domain ownership
- `DELETE /api/domains/{id}` - Delete custom domain

## Development

### Database Management
- Create new migration:
```bash
make migrate-create
```

- Apply migrations:
```bash
make migrate-up
```

- Rollback last migration:
```bash
make migrate-down
```

### Code Generation
- Generate SQL code:
```bash
make sqlc
```

### Testing and Linting
- Run tests:
```bash
make test
```

- Run linter:
```bash
make lint
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 