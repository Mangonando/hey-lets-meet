# hey-lets-meet — Backend

Go REST API with SQLite, session-based auth, and automatic database migrations.

## Requirements

- Go 1.26+
- [golangci-lint](https://golangci-lint.run/usage/install/) (for linting)

## Run

```sh
cd backend
go run ./cmd/api
```

The server starts on `http://localhost:8080`. On first run it creates `hey-lets-meet.db` and applies all pending migrations automatically.

## Test

```sh
cd backend
go test ./cmd/api/...
```

## Lint

```sh
cd backend
golangci-lint run
```

## Adding a migration

Create a new `.sql` file in `migrations/` following the naming convention:

```
004_your_description.sql
```

Files are sorted alphabetically and run in order. Each migration runs exactly once and is recorded in the `schema_migrations` table.

## Auth endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/register` | Register a new user |
| POST | `/auth/login` | Log in |
| POST | `/auth/logout` | Log out |
| GET | `/auth/me` | Get current user (requires auth) |
| GET | `/api/protected` | Example protected route |
| GET | `/health` | Health check |

### Example — register

```sh
curl -i -c cookies.txt -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"pwd123456"}' \
  http://localhost:8080/auth/register
```

### Example — access protected route

```sh
curl -i -b cookies.txt http://localhost:8080/api/protected
```

## Project structure

```
backend/
├── cmd/api/          # Entry point (main.go) and integration tests
├── internal/
│   ├── auth/         # Auth handlers, service, repo, middleware
│   └── db/           # Database connection and migration runner
│   └── httpapi/      # HTTP server and route registration
└── migrations/       # SQL migration files
```
