# hey-lets-meet — Backend

Go REST API with SQLite, session-based auth, and automatic database migrations.

## Requirements

- Go 1.26+
- [golangci-lint](https://golangci-lint.run/usage/install/) (for linting)

## Environment

Create `backend/.env`:

```
ORS_API_KEY=your_key_here
```

Get a free key at [openrouteservice.org](https://openrouteservice.org). Without it the server starts with mock providers.

## Run

```sh
cd backend
export $(cat .env | xargs) && go run ./cmd/api
```

The server starts on `http://localhost:8080`. On first run it creates `data/app.db` and applies all pending migrations automatically.

If `ORS_API_KEY` is not set, mock geocoding and routing are used instead.

## Test

```sh
cd backend
go test ./...
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

## Meetpoints endpoint

Suggests a fair meeting point between two locations based on walking time.
Requires authentication.

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/meetpoints/suggest` | Suggest a meetpoint between two addresses |

### Example

```sh
# 1. Login first
curl -s -c cookies.txt -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"pwd123456"}' \
  http://localhost:8080/auth/login

# 2. Suggest a meetpoint
curl -s -b cookies.txt -H "Content-Type: application/json" \
  -d '{"originA":"Alexanderplatz","originB":"Hermannplatz"}' \
  http://localhost:8080/api/meetpoints/suggest
```

**Response:**
```json
{
  "origins": {
    "a": {"address": "Alexanderplatz", "point": {"lat": 52.521918, "lng": 13.413215}},
    "b": {"address": "Hermannplatz",   "point": {"lat": 52.486355, "lng": 13.424318}}
  },
  "best": {
    "point": {"lat": 52.504, "lng": 13.418},
    "etaASeconds": 1454,
    "etaBSeconds": 1451,
    "maxEtaSeconds": 1454,
    "diffSeconds": 3,
    "distanceAMeters": 2013,
    "distanceBMeters": 2013
  },
  "alternatives": [...],
  "debug": {"midpoint": {"lat": 52.504, "lng": 13.418}}
}
```

> Currently uses mock geocoding and walking-only routing. Real geocoding and transport modes coming soon.

## Project structure

```
backend/
├── cmd/api/          # Entry point (main.go) and integration tests
├── internal/
│   ├── auth/         # Auth handlers, service, repo, middleware
│   ├── db/           # Database connection and migration runner
│   ├── httpapi/      # HTTP server, route registration, response helpers
│   └── meetpoints/   # Meetpoint suggestion logic, mocks, and HTTP handler
└── migrations/       # SQL migration files
```
