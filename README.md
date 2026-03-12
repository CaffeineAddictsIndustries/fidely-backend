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

### Install Go dependencies

```bash
go mod download
```

### Install golang-migrate

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Getting Started

### 1. Set up the database

#### Option A: Local PostgreSQL

Create the application role and database:

```bash
sudo -u postgres psql -d postgres -c "CREATE ROLE fidely LOGIN PASSWORD 'fidely';"
sudo -u postgres createdb -O fidely fidely
```

If they already exist, you can just reset the password:

```bash
sudo -u postgres psql -d postgres -c "ALTER ROLE fidely WITH LOGIN PASSWORD 'fidely';"
```

Verify connectivity:

```bash
psql "postgres://fidely:fidely@localhost:5432/fidely?sslmode=disable" -c "SELECT current_user, current_database();"
```

#### Option B: Docker

```bash
docker compose up -d
```

If using Docker from WSL, `host.docker.internal` may be required for `DATABASE_URL`.

### 2. Configure environment

The app reads environment variables from the shell; it does not load `.env` automatically.

For local PostgreSQL:

```bash
export DATABASE_URL="postgres://fidely:fidely@localhost:5432/fidely?sslmode=disable"
export SERVER_PORT=8080
```

For Docker from WSL:

```bash
export DATABASE_URL="postgres://fidely:fidely@host.docker.internal:5432/fidely?sslmode=disable"
```

### 3. Run migrations

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

### 4. Start the project

```bash
go run cmd/api/main.go
```

Server runs on `http://localhost:8080`

### 5. Test the server

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
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

## Project Setup Summary

```bash
go mod download
export DATABASE_URL="postgres://fidely:fidely@localhost:5432/fidely?sslmode=disable"
export SERVER_PORT=8080
migrate -path migrations -database "$DATABASE_URL" up
go run cmd/api/main.go
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
