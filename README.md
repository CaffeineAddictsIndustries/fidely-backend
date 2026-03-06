# Fidely Backend

Go backend for the Fidely loyalty card platform.

## Tech Stack

- **Go 1.25+**
- **Echo** — HTTP framework
- **pgx** — PostgreSQL driver
- **golang-migrate** — Database migrations

## Prerequisites

- Go 1.25+
- PostgreSQL 14+ (local install) **or** Docker Desktop
- golang-migrate CLI

### Install golang-migrate

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Getting Started

### 1. Configure environment

Use `.env.example` values (or export variables directly):

```bash
export DATABASE_URL="postgres://fidely:fidely@localhost:5432/fidely?sslmode=disable"
export SERVER_PORT=8080
```

### 2. Start PostgreSQL

#### Option A: Local PostgreSQL

Make sure a database named `fidely` exists and your `DATABASE_URL` credentials are valid.

#### Option B: Docker

```bash
docker compose up -d
```

If using Docker from WSL, `host.docker.internal` may be required:

```bash
export DATABASE_URL="postgres://fidely:fidely@host.docker.internal:5432/fidely?sslmode=disable"
```

### 3. Run Migrations

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

### 4. Start Server

```bash
go run cmd/api/main.go
```

Server runs on `http://localhost:8080`

### 5. Test Health Check

```bash
curl http://localhost:8080/health
```

## Database

- **Host:** `localhost` (default local setup)
- **Port:** 5432
- **User:** fidely
- **Password:** fidely
- **Database:** fidely

### Migration Commands

```bash
# Apply all migrations
migrate -path migrations -database "$DATABASE_URL" up

# Rollback last migration
migrate -path migrations -database "$DATABASE_URL" down 1

# Rollback all migrations
migrate -path migrations -database "$DATABASE_URL" down

# Check current version
migrate -path migrations -database "$DATABASE_URL" version
```

## Project Structure

```
fidely-backend/
├── cmd/api/
│   └── main.go                 # Entrypoint
├── internal/
│   ├── config/                 # Environment config
│   ├── db/                     # Database connection
│   ├── model/                  # Data models
│   ├── repository/             # Database queries (scaffold)
│   ├── service/                # Business logic (scaffold)
│   └── handler/                # HTTP handlers (scaffold)
├── migrations/                 # SQL migrations
├── docker-compose.yml
├── go.mod
└── go.sum
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | (required) |
| `SERVER_PORT` | HTTP server port | 8080 |

## Current Scope

- Infrastructure bootstrap is ready: server startup, DB connection, migrations, and health endpoint.
- Business endpoints and authentication are not implemented yet.
